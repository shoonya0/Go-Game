package core

import (
	"image"
	"image/color"
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// The travel menu is a centered column of thumbnails, one per active portal
// (D9). It shows at most travelMenuRows rows at a time, scrolling with the pick.
const (
	travelHoldSeconds = 0.5 // how long E, or a touch on the portal, is held to open the menu

	travelThumbW   = 272 // a thumbnail is a fifth of the 1360×768 screen
	travelThumbH   = 154
	travelRowGap   = 14
	travelMenuRows = 4
	travelBorder   = 4 // px of the frame around a row
	travelArrow    = 44
)

var (
	travelShade      = color.RGBA{0, 0, 0, 0xa0} // over the paused game
	travelRowColor   = color.RGBA{0x30, 0x24, 0x1c, 0xff}
	travelPickColor  = color.RGBA{0xff, 0xff, 0xff, 0xff}
	travelSavedColor = color.RGBA{0xff, 0xc8, 0x3c, 0xff} // the checkpoint's row (D9)
)

// travelHoldTicks is travelHoldSeconds in ticks.
func travelHoldTicks() int { return int(math.Round(travelHoldSeconds * tps())) }

// TravelMenu lists the player's active portals to travel to. While it is open
// the game is paused.
type TravelMenu struct {
	Open     bool
	Rows     []PortalRef // the active portals, in the order they were activated
	Selected int
	Saved    PortalRef // the checkpoint when the menu opened, framed gold

	waitRelease bool // closed while E or a touch was down: that hold can't reopen it

	screenW, screenH int
	thumbs           map[PortalRef]*ebiten.Image // a small picture of the level around each portal
	arrow            *ebiten.Image               // the wooden arrow, turned up or down when rows are scrolled off
}

// NewTravelMenu lays the menu out for a logical screen of w×h pixels.
func NewTravelMenu(w, h int) *TravelMenu {
	return &TravelMenu{screenW: w, screenH: h, thumbs: map[PortalRef]*ebiten.Image{}}
}

// HasThumb reports whether portal ref has its thumbnail.
func (m *TravelMenu) HasThumb(ref PortalRef) bool { return m.thumbs[ref] != nil }

// SetThumb keeps view, a full-screen picture of the level around portal ref,
// shrunk to a thumbnail.
func (m *TravelMenu) SetThumb(ref PortalRef, view *ebiten.Image) {
	thumb := ebiten.NewImage(travelThumbW, travelThumbH)
	drawScaled(thumb, view, thumb.Bounds(), nil, 1)
	m.thumbs[ref] = thumb
}

// OpenFor opens the menu on p's active portals, picking the one p stands in
// (else the checkpoint).
func (m *TravelMenu) OpenFor(p *PlayerRuntime) {
	m.Open = true
	m.Rows = slices.Clone(p.ActivePortals)
	m.Saved = p.Checkpoint
	m.Selected = max(0, slices.Index(m.Rows, p.Checkpoint))
	if at := p.standingIn(); at != nil {
		m.Selected = max(0, slices.Index(m.Rows, *at))
	}
}

// WaitRelease makes a hold of E or a touch already down not open the menu,
// like one left down when it closed: the press that picked a title menu item.
func (m *TravelMenu) WaitRelease() { m.waitRelease = true }

// MaybeOpen opens the menu when p asked for it this tick, unless the hold
// began before the menu last closed (the E that confirmed a pick, a finger
// left on a row): E and the touch must be let go first.
func (m *TravelMenu) MaybeOpen(p *PlayerRuntime, in *InputState) {
	if in.InteractTicks == 0 && in.TouchHoldTicks == 0 {
		m.waitRelease = false
	}
	if p.TravelAsked() && !m.waitRelease {
		m.OpenFor(p)
	}
}

// Update applies one tick of input to the open menu: up/down move the pick, E
// or Enter travels to it, Esc closes. On touch screens a tap on a row travels
// there and a tap anywhere else closes. It returns the portal to travel to.
func (m *TravelMenu) Update(in *InputState) (PortalRef, bool) {
	if !m.Open {
		return PortalRef{}, false
	}
	defer func() { m.waitRelease = !m.Open }()
	switch {
	case in.Back:
		m.Open = false
	case in.Tapped:
		m.Open = false
		first, last := m.VisibleRows()
		for i := first; i < last; i++ {
			if image.Pt(int(in.TapAt.X), int(in.TapAt.Y)).In(m.RowRect(i)) {
				return m.Rows[i], true
			}
		}
	case in.Interact || in.Confirm:
		m.Open = false
		if m.Selected < len(m.Rows) {
			return m.Rows[m.Selected], true
		}
	case in.UpJustPressed:
		m.Selected = max(0, m.Selected-1)
	case in.DownJustPressed:
		m.Selected = min(len(m.Rows)-1, m.Selected+1)
	}
	return PortalRef{}, false
}

// VisibleRows returns the rows on screen, [first, last): up to
// travelMenuRows, keeping the pick among them.
func (m *TravelMenu) VisibleRows() (first, last int) {
	first = min(max(0, m.Selected-travelMenuRows/2), max(0, len(m.Rows)-travelMenuRows))
	return first, min(len(m.Rows), first+travelMenuRows)
}

// RowRect returns the screen area of row i's thumbnail; i must be visible.
func (m *TravelMenu) RowRect(i int) image.Rectangle {
	first, last := m.VisibleRows()
	n := last - first
	top := (m.screenH - n*travelThumbH - (n-1)*travelRowGap) / 2
	x := (m.screenW - travelThumbW) / 2
	y := top + (i-first)*(travelThumbH+travelRowGap)
	return image.Rect(x, y, x+travelThumbW, y+travelThumbH)
}

// Draw shades the paused game and draws the visible rows: each portal's
// thumbnail, the pick framed white, the checkpoint gold, with arrows when
// rows are scrolled off above or below.
func (m *TravelMenu) Draw(screen *ebiten.Image) {
	if !m.Open {
		return
	}
	if m.arrow == nil {
		m.arrow = uiSprite(uiLeftRect)
	}
	vector.FillRect(screen, 0, 0, float32(m.screenW), float32(m.screenH), travelShade, false)

	first, last := m.VisibleRows()
	for i := first; i < last; i++ {
		r := m.RowRect(i)
		x, y, w, h := float32(r.Min.X), float32(r.Min.Y), float32(r.Dx()), float32(r.Dy())
		if thumb := m.thumbs[m.Rows[i]]; thumb != nil {
			drawScaled(screen, thumb, r, nil, 1)
		} else {
			vector.FillRect(screen, x, y, w, h, travelRowColor, false)
		}
		if i != m.Selected {
			vector.FillRect(screen, x, y, w, h, color.RGBA{0, 0, 0, 0x60}, false) // the pick stands out
		}
		const b = travelBorder
		switch {
		case m.Rows[i] == m.Saved:
			vector.StrokeRect(screen, x-b/2, y-b/2, w+b, h+b, b, travelSavedColor, false)
			if i == m.Selected {
				vector.StrokeRect(screen, x-1.5*b, y-1.5*b, w+3*b, h+3*b, b, travelPickColor, false)
			}
		case i == m.Selected:
			vector.StrokeRect(screen, x-b/2, y-b/2, w+b, h+b, b, travelPickColor, false)
		}
	}

	cx := m.screenW / 2
	if first > 0 {
		top := m.RowRect(first).Min.Y - travelRowGap/2 - travelArrow
		drawRotated(screen, m.arrow, image.Rect(cx-travelArrow/2, top, cx+travelArrow/2, top+travelArrow), 1, nil, 1)
	}
	if last < len(m.Rows) {
		top := m.RowRect(last-1).Max.Y + travelRowGap/2
		drawRotated(screen, m.arrow, image.Rect(cx-travelArrow/2, top, cx+travelArrow/2, top+travelArrow), -1, nil, 1)
	}
}
