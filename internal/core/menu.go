package core

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// A menu is a title over a centered column of items (main-menu spec).
const (
	menuTitleSize = 44 // px of the title's letters
	menuItemSize  = 26 // px of an item's letters
	menuItemPadX  = 28 // px around an item's text inside its box
	menuItemPadY  = 16
	menuItemGap   = 18 // px between item boxes
	menuTitleGap  = 64 // px between the title and the first item
	menuBorder    = 3
)

var (
	menuShade         = color.RGBA{0, 0, 0, 0xb0} // over the game behind
	menuTitleColor    = color.RGBA{0xff, 0xc8, 0x3c, 0xff}
	menuItemColor     = color.RGBA{0xe8, 0xe0, 0xd0, 0xff}
	menuPickColor     = color.RGBA{0xff, 0xff, 0xff, 0xff}
	menuDisabledColor = color.RGBA{0x70, 0x6a, 0x60, 0xff}
	menuBoxColor      = color.RGBA{0x30, 0x24, 0x1c, 0xe0}
)

// MenuItem is one line of a menu. A disabled item is drawn gray and can't be
// picked (Continue without a save).
type MenuItem struct {
	Label    string
	Disabled bool
}

// Menu is a column of items picked with up/down and Enter or E, a tap or a
// mouse click; hovering the mouse over an item picks it too.
type Menu struct {
	Title    string
	Items    []MenuItem
	Selected int

	screenW, screenH int
	boxW             int      // width of every item box: the widest label's, measured when labels change
	pointer          Position // the mouse cursor last tick, to see it move
	pointerSeen      bool     // false until the first tick after Open: a resting cursor doesn't pick
}

// NewMenu lays a menu out for a logical screen of w×h pixels.
func NewMenu(w, h int, title string, labels ...string) *Menu {
	m := &Menu{Title: title, screenW: w, screenH: h}
	for _, l := range labels {
		m.Items = append(m.Items, MenuItem{Label: l})
	}
	m.measure()
	return m
}

// Open shows the menu with item sel picked. The mouse picks an item only once
// it moves: one resting where the last screen had its pick would otherwise
// move this pick (onto Yes of the overwrite question, say).
func (m *Menu) Open(sel int) {
	m.Selected, m.pointerSeen = sel, false
}

// SetLabel changes item i's text.
func (m *Menu) SetLabel(i int, label string) {
	m.Items[i].Label = label
	m.measure()
}

// measure sizes the item boxes to the widest label.
func (m *Menu) measure() {
	var w float64
	for _, it := range m.Items {
		w = max(w, textWidth(it.Label, menuItemSize))
	}
	m.boxW = int(w) + 2*menuItemPadX
}

// SetDisabled enables or disables item i, moving the pick off it if needed.
func (m *Menu) SetDisabled(i int, disabled bool) {
	m.Items[i].Disabled = disabled
	if m.Items[m.Selected].Disabled {
		m.Selected = m.step(m.Selected, 1)
	}
}

// step returns the next enabled item from i in direction dir (±1), wrapping
// around; i itself when no other is enabled.
func (m *Menu) step(i, dir int) int {
	n := len(m.Items)
	for k := 1; k <= n; k++ {
		j := ((i+dir*k)%n + n) % n
		if !m.Items[j].Disabled {
			return j
		}
	}
	return i
}

// Update applies one tick of input. It returns the item chosen this tick, or
// -1, and whether Esc asked to go back.
func (m *Menu) Update(in *InputState) (chosen int, back bool) {
	moved := m.pointerSeen && in.Pointer != m.pointer
	m.pointer, m.pointerSeen = in.Pointer, true
	if moved {
		if i := m.itemAt(in.Pointer); i >= 0 {
			m.Selected = i
		}
	}
	at, pressed := in.PressAt()
	switch {
	case in.Back:
		return -1, true
	case pressed:
		if i := m.itemAt(at); i >= 0 {
			m.Selected = i
			return i, false
		}
	case in.Confirm || in.Interact:
		if !m.Items[m.Selected].Disabled {
			return m.Selected, false
		}
	case in.UpJustPressed:
		m.Selected = m.step(m.Selected, -1)
	case in.DownJustPressed:
		m.Selected = m.step(m.Selected, 1)
	}
	return -1, false
}

// itemAt returns the enabled item whose box holds screen point p, or -1.
func (m *Menu) itemAt(p Position) int {
	pt := image.Pt(int(p.X), int(p.Y))
	for i, it := range m.Items {
		if !it.Disabled && pt.In(m.ItemRect(i)) {
			return i
		}
	}
	return -1
}

// ItemRect returns the screen box of item i. Every box is as wide as the
// widest label, and the title and the column are centered together.
func (m *Menu) ItemRect(i int) image.Rectangle {
	boxW, boxH := m.boxW, menuItemSize+2*menuItemPadY
	n := len(m.Items)
	column := n*boxH + (n-1)*menuItemGap
	top := (m.screenH-column)/2 + (menuTitleSize+menuTitleGap)/2
	x := (m.screenW - boxW) / 2
	y := top + i*(boxH+menuItemGap)
	return image.Rect(x, y, x+boxW, y+boxH)
}

func textWidth(s string, size float64) float64 {
	w, _ := TextSize(s, size)
	return w
}

// Draw shades the screen and draws the title and the items: the pick framed
// white, disabled items gray.
func (m *Menu) Draw(screen *ebiten.Image) {
	vector.FillRect(screen, 0, 0, float32(m.screenW), float32(m.screenH), menuShade, false)
	first := m.ItemRect(0)
	DrawTextCentered(screen, m.Title, menuTitleSize, float64(m.screenW)/2,
		float64(first.Min.Y-menuTitleGap-menuTitleSize/2), menuTitleColor)
	for i, it := range m.Items {
		r := m.ItemRect(i)
		x, y, w, h := float32(r.Min.X), float32(r.Min.Y), float32(r.Dx()), float32(r.Dy())
		vector.FillRect(screen, x, y, w, h, menuBoxColor, false)
		clr := menuItemColor
		switch {
		case it.Disabled:
			clr = menuDisabledColor
		case i == m.Selected:
			clr = menuPickColor
			vector.StrokeRect(screen, x, y, w, h, menuBorder, menuPickColor, false)
		}
		DrawTextCentered(screen, it.Label, menuItemSize, float64(r.Min.X+r.Max.X)/2, float64(r.Min.Y+r.Max.Y)/2, clr)
	}
}
