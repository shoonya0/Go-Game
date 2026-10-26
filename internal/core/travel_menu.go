package core

import (
	"cmp"
	"fmt"
	"image"
	"image/color"
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// The travel menu is a wooden board listing the active portals under their
// level's and sub-level's names, a plank per portal (levels D7). A list
// taller than the screen scrolls with the pick.
const (
	travelHoldSeconds = 0.5 // how long E, or a touch on the portal, is held to open the menu

	travelHeadSize = 36 // px of a level's name
	travelSubSize  = 24 // px of a sub-level's name
	travelItemSize = 24 // px of a portal's name
	travelHeadLine = 64 // px of the list a level's name takes
	travelSubLine  = 44
	travelPlankW   = 340
	travelPlankH   = 64
	travelRowGap   = 12
	travelBoardPad = 60 // px of board around the list
	travelMargin   = 20 // px between the board and the screen's edges
	travelArrow    = 36
)

var travelSubColor = color.RGBA{0xf4, 0xe4, 0xc4, 0xff} // a sub-level's name, under the gold level name

// travelHoldTicks is travelHoldSeconds in ticks.
func travelHoldTicks() int { return int(math.Round(travelHoldSeconds * tps())) }

// TravelMenu lists the player's active portals to travel to. While it is open
// the game is paused.
type TravelMenu struct {
	Open     bool
	Rows     []PortalRef // the active portals, by level in LevelOrder, then by number
	Selected int
	Saved    PortalRef // the checkpoint when the menu opened, marked with a check

	waitRelease bool // closed while E or a touch was down: that hold can't reopen it

	screenW, screenH int
	arrow            *ebiten.Image // the wooden arrow, turned up or down when the list is scrolled
}

// travelLine is one line of the menu's list: a level's or a sub-level's
// name, or a portal's plank.
type travelLine struct {
	text string
	size float64 // px of the letters
	row  int     // index in Rows of a plank; -1 for a name
	y, h int     // px from the list's top
}

// NewTravelMenu lays the menu out for a logical screen of w×h pixels.
func NewTravelMenu(w, h int) *TravelMenu {
	return &TravelMenu{screenW: w, screenH: h}
}

// OpenFor opens the menu on p's active portals, picking the one p stands in
// (else the checkpoint).
func (m *TravelMenu) OpenFor(p *PlayerRuntime) {
	m.Open = true
	m.Rows = slices.Clone(p.ActivePortals)
	sortTravelRows(m.Rows)
	m.Saved = p.Checkpoint
	m.Selected = max(0, slices.Index(m.Rows, p.Checkpoint))
	if at := p.standingIn(); at != nil {
		m.Selected = max(0, slices.Index(m.Rows, *at))
	}
}

// sortTravelRows orders rows by level, in LevelOrder (levels not in it
// last), then by portal number; portals of levels not in Levels have no
// number and go in reading order.
func sortTravelRows(rows []PortalRef) {
	place := func(level string) int {
		if i := slices.Index(LevelOrder, level); i >= 0 {
			return i
		}
		return len(LevelOrder)
	}
	slices.SortStableFunc(rows, func(a, b PortalRef) int {
		return cmp.Or(
			cmp.Compare(place(a.Level), place(b.Level)),
			cmp.Compare(a.Level, b.Level),
			cmp.Compare(portalNumber(a), portalNumber(b)),
			cmp.Compare(a.Cell.Y, b.Cell.Y),
			cmp.Compare(a.Cell.X, b.Cell.X),
		)
	})
}

// levelNames returns the names the menu shows for a level: its title
// ("Deep Forest") and its sub-level's ("Deep Forest Day"); empty for a level
// not in Levels.
func levelNames(level string) (title, sub string) {
	def, ok := Levels[level]
	if !ok {
		return "", ""
	}
	return def.Title, def.Title + " " + def.Time
}

// lines lays out the list: a level's name where the level changes, the
// sub-level's name where the sub-level does, and a plank per row.
func (m *TravelMenu) lines() []travelLine {
	var out []travelLine
	y := 0
	add := func(text string, size float64, row, h int) {
		out = append(out, travelLine{text: text, size: size, row: row, y: y, h: h})
		y += h
	}
	lastTitle, n := "", 0
	for i, ref := range m.Rows {
		if i == 0 || ref.Level != m.Rows[i-1].Level {
			title, sub := levelNames(ref.Level)
			if title != "" && (i == 0 || title != lastTitle) {
				add(title, travelHeadSize, -1, travelHeadLine)
			}
			if sub != "" {
				add(sub, travelSubSize, -1, travelSubLine)
			}
			lastTitle, n = title, 0
		}
		n++
		num := portalNumber(ref)
		if num == 0 { // not in Levels: count the rows instead
			num = n
		}
		add(fmt.Sprintf("Portal %d", num), travelItemSize, i, travelPlankH+travelRowGap)
	}
	return out
}

// listHeight is how tall lines are, without the gap after the last plank.
func listHeight(lines []travelLine) int {
	if len(lines) == 0 {
		return 0
	}
	return lines[len(lines)-1].y + lines[len(lines)-1].h - travelRowGap
}

// Labels returns the text of every line of the list, top to bottom.
func (m *TravelMenu) Labels() []string {
	var out []string
	for _, l := range m.lines() {
		out = append(out, l.text)
	}
	return out
}

// view returns the screen area the list shows through, and how far the list
// is scrolled: when it is too tall, the pick is kept in the middle.
func (m *TravelMenu) view(lines []travelLine) (r image.Rectangle, scroll int) {
	content := listHeight(lines)
	h := min(content, m.screenH-2*(travelMargin+travelBoardPad))
	w := travelPlankW + 2*(travelPlankH+menuIconGap) // room for the check left of a plank
	for _, l := range lines {
		if l.row < 0 {
			w = max(w, int(textWidth(l.text, l.size)))
		}
	}
	x, y := (m.screenW-w)/2, (m.screenH-h)/2
	for _, l := range lines {
		if l.row == m.Selected {
			scroll = l.y + travelPlankH/2 - h/2
		}
	}
	return image.Rect(x, y, x+w, y+h), min(max(0, scroll), content-h)
}

// lineRect returns where line l is on the screen; it may be scrolled off.
func (m *TravelMenu) lineRect(l travelLine, view image.Rectangle, scroll int) image.Rectangle {
	y := view.Min.Y + l.y - scroll
	if l.row < 0 {
		return image.Rect(view.Min.X, y, view.Max.X, y+l.h)
	}
	x := (m.screenW - travelPlankW) / 2
	return image.Rect(x, y, x+travelPlankW, y+travelPlankH)
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

// VisibleRows returns the rows whose planks are wholly on screen,
// [first, last); the pick is always among them.
func (m *TravelMenu) VisibleRows() (first, last int) {
	lines := m.lines()
	view, scroll := m.view(lines)
	first = -1
	for _, l := range lines {
		if l.row >= 0 && m.lineRect(l, view, scroll).In(view) {
			if first < 0 {
				first = l.row
			}
			last = l.row + 1
		}
	}
	return max(0, first), last
}

// RowRect returns the screen area of row i's plank; i must be visible.
func (m *TravelMenu) RowRect(i int) image.Rectangle {
	lines := m.lines()
	view, scroll := m.view(lines)
	for _, l := range lines {
		if l.row == i {
			return m.lineRect(l, view, scroll)
		}
	}
	return image.Rectangle{}
}

// Draw shades the paused game and draws the list on the wooden board: each
// level's name in gold, its sub-levels' names under it, and a plank per
// portal. The pick is the arrow plank, the checkpoint has a check at its
// left, and arrows show when the list is scrolled.
func (m *TravelMenu) Draw(screen *ebiten.Image) {
	if !m.Open {
		return
	}
	if m.arrow == nil {
		m.arrow = uiSprite(uiLeftRect)
	}
	vector.FillRect(screen, 0, 0, float32(m.screenW), float32(m.screenH), menuShade, false)

	lines := m.lines()
	view, scroll := m.view(lines)
	drawSliced(screen, uiBoardRect, uiBoardInsets, view.Inset(-travelBoardPad), menuBoardScale, nil, 1)
	list := screen.SubImage(view).(*ebiten.Image) // cuts lines scrolled off at its edges
	for _, l := range lines {
		r := m.lineRect(l, view, scroll)
		if !r.Overlaps(view) {
			continue
		}
		cx, cy := float64(r.Min.X+r.Max.X)/2, float64(r.Min.Y+r.Max.Y)/2
		switch {
		case l.row < 0 && l.size == travelHeadSize:
			DrawTextCentered(list, l.text, l.size, cx+3, cy+3, menuTitleShadow)
			DrawTextCentered(list, l.text, l.size, cx, cy, menuTitleColor)
		case l.row < 0:
			DrawTextCentered(list, l.text, l.size, cx+2, cy+2, menuTitleShadow)
			DrawTextCentered(list, l.text, l.size, cx, cy, travelSubColor)
		default:
			plank, insets := uiPlankRect, uiPlankInsets
			var tint color.Color = menuIdleTint
			if l.row == m.Selected {
				plank, insets, tint = uiPlankArrowRect, uiPlankArrowInsets, nil
			}
			drawPlank(list, plank, insets, r, tint)
			if m.Rows[l.row] == m.Saved {
				check := image.Rect(r.Min.X-menuIconGap-r.Dy(), r.Min.Y, r.Min.X-menuIconGap, r.Max.Y)
				drawRotated(list, uiSprite(uiCheckRect), check, 0, nil, 1)
			}
			DrawTextCentered(list, l.text, l.size, cx, float64(r.Min.Y)+menuTextRise*float64(r.Dy()), menuItemColor)
		}
	}

	cx := m.screenW / 2
	if scroll > 0 {
		top := view.Min.Y - travelArrow - 6
		drawRotated(screen, m.arrow, image.Rect(cx-travelArrow/2, top, cx+travelArrow/2, top+travelArrow), 1, nil, 1)
	}
	if scroll+view.Dy() < listHeight(lines) {
		top := view.Max.Y + 6
		drawRotated(screen, m.arrow, image.Rect(cx-travelArrow/2, top, cx+travelArrow/2, top+travelArrow), -1, nil, 1)
	}
}
