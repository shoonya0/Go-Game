package core

import (
	"image"
	"math"
	"slices"
)

// The travel menu is a centered column of thumbnails, one per active portal
// (D9). It shows at most travelMenuRows rows at a time, scrolling with the pick.
const (
	travelHoldSeconds = 0.5 // how long E, or a touch on the portal, is held to open the menu

	travelThumbW   = 272 // a thumbnail is a fifth of the 1360×768 screen
	travelThumbH   = 154
	travelRowGap   = 14
	travelMenuRows = 4
)

// travelHoldTicks is travelHoldSeconds in ticks.
func travelHoldTicks() int { return int(math.Round(travelHoldSeconds * tps())) }

// TravelMenu lists the player's active portals to travel to. While it is open
// the game is paused.
type TravelMenu struct {
	Open     bool
	Rows     []PortalRef // the active portals, in the order they were activated
	Selected int

	screenW, screenH int
}

// NewTravelMenu lays the menu out for a logical screen of w×h pixels.
func NewTravelMenu(w, h int) *TravelMenu { return &TravelMenu{screenW: w, screenH: h} }

// OpenFor opens the menu on p's active portals, picking the one p stands in
// (else the checkpoint).
func (m *TravelMenu) OpenFor(p *PlayerRuntime) {
	m.Open = true
	m.Rows = slices.Clone(p.ActivePortals)
	m.Selected = max(0, slices.Index(m.Rows, p.Checkpoint))
	if at := p.standingIn(); at != nil {
		m.Selected = max(0, slices.Index(m.Rows, *at))
	}
}

// Update applies one tick of input to the open menu: up/down move the pick, E
// or Enter travels to it, Esc closes. On touch screens a tap on a row travels
// there and a tap anywhere else closes. It returns the portal to travel to.
func (m *TravelMenu) Update(in *InputState) (PortalRef, bool) {
	if !m.Open {
		return PortalRef{}, false
	}
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
