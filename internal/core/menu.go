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
	menuItemPadY  = 22
	menuItemMinW  = 260 // px: a short label (No, Yes) still gets a full plank
	menuItemGap   = 18  // px between item boxes
	menuTitleGap  = 64  // px between the title and the first item
	menuIconGap   = 14  // px between an item's icon and its plank
	menuPanelPad  = 72  // px of board around the title and the items

	menuBoardScale = 0.5  // screen px per sheet px of the board's frame
	menuTextRise   = 0.45 // the plank's gold face is centered this far down the plank
)

var (
	menuShade         = color.RGBA{0, 0, 0, 0xb0} // over the game behind
	menuTitleColor    = color.RGBA{0xff, 0xc8, 0x3c, 0xff}
	menuTitleShadow   = color.RGBA{0x2a, 0x10, 0x08, 0xff}
	menuItemColor     = color.RGBA{0x4a, 0x1c, 0x0c, 0xff} // dark brown on the gold face
	menuDisabledColor = color.RGBA{0x70, 0x6a, 0x60, 0xff}
	menuIdleTint      = color.RGBA{0xd8, 0xd0, 0xc8, 0xff} // items not picked are a bit darker
	menuDisabledTint  = color.RGBA{0x60, 0x58, 0x50, 0xff}
)

// MenuItem is one line of a menu. A disabled item is drawn gray and can't be
// picked (Continue without a save).
type MenuItem struct {
	Label    string
	Disabled bool
	Icon     image.Rectangle // sprite on UIButtonsSheet left of the plank; empty for none
	Bar      bool            // drawn as a bar filled to Fill (0..1) instead of the label
	Fill     float64
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

// SetIcons gives the items their icons, in order (see MenuItem.Icon). Icons
// past the last item are ignored: Quit is left out where the game can't close.
func (m *Menu) SetIcons(icons ...image.Rectangle) {
	for i := range min(len(icons), len(m.Items)) {
		m.Items[i].Icon = icons[i]
	}
}

// measure sizes the item boxes to the widest label, or bar.
func (m *Menu) measure() {
	w := float64(menuItemMinW)
	for _, it := range m.Items {
		w = max(w, textWidth(it.Label, menuItemSize)+2*menuItemPadX)
		if it.Bar {
			w = max(w, float64(m.barRect(image.Rect(0, 0, 0, menuBoxH())).Dx())+plankSides(menuBoxH()))
		}
	}
	m.boxW = int(w)
}

func menuBoxH() int { return menuItemSize + 2*menuItemPadY }

// plankSides is how wide a plank of height h draws its two ends: the arrow
// plank's, the wider, so a bar or label fits either plank.
func plankSides(h int) float64 {
	return float64((uiPlankArrowInsets.left+uiPlankArrowInsets.right)*h) / float64(uiPlankArrowRect.Dy())
}

// iconSpace is the room the icons take left of the planks: none when no item
// has an icon.
func (m *Menu) iconSpace() int {
	for _, it := range m.Items {
		if !it.Icon.Empty() {
			return menuBoxH() + menuIconGap
		}
	}
	return 0
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

// ItemRect returns the screen box of item i, its plank. Every box is as wide
// as the widest label; the title and the column (icons and planks) are
// centered together.
func (m *Menu) ItemRect(i int) image.Rectangle {
	boxW, boxH := m.boxW, menuBoxH()
	n := len(m.Items)
	column := n*boxH + (n-1)*menuItemGap
	top := (m.screenH-column)/2 + (menuTitleSize+menuTitleGap)/2
	x := (m.screenW-boxW-m.iconSpace())/2 + m.iconSpace()
	y := top + i*(boxH+menuItemGap)
	return image.Rect(x, y, x+boxW, y+boxH)
}

// IconRect returns the square left of item i's plank where its icon goes.
func (m *Menu) IconRect(i int) image.Rectangle {
	r := m.ItemRect(i)
	return image.Rect(r.Min.X-menuIconGap-r.Dy(), r.Min.Y, r.Min.X-menuIconGap, r.Max.Y)
}

// titleCenterY is the screen y the title's letters are centered on.
func (m *Menu) titleCenterY() int {
	return m.ItemRect(0).Min.Y - menuTitleGap - menuTitleSize/2
}

// PanelRect returns the board behind the menu: the title, the icons and the
// planks with menuPanelPad around them.
func (m *Menu) PanelRect() image.Rectangle {
	tw := int(textWidth(m.Title, menuTitleSize))
	ty := m.titleCenterY()
	r := image.Rect((m.screenW-tw)/2, ty-menuTitleSize/2, (m.screenW+tw)/2, ty+menuTitleSize/2)
	for i := range m.Items {
		r = r.Union(m.ItemRect(i))
		if !m.Items[i].Icon.Empty() {
			r = r.Union(m.IconRect(i))
		}
	}
	return r.Inset(-menuPanelPad)
}

// barRect is where a bar goes on plank r: on its gold face, between the ends.
func (m *Menu) barRect(r image.Rectangle) image.Rectangle {
	s := float64(r.Dy()) / float64(uiPlankRect.Dy())
	faceTop, faceBottom := int(42*s), int(172*s) // the face's rows on the sheet's plank
	h := faceBottom - faceTop
	w := h * uiBarFullRect.Dx() / uiBarFullRect.Dy()
	cx := (r.Min.X + r.Max.X) / 2
	return image.Rect(cx-w/2, r.Min.Y+faceTop, cx-w/2+w, r.Min.Y+faceBottom)
}

func textWidth(s string, size float64) float64 {
	w, _ := TextSize(s, size)
	return w
}

// Draw shades the screen and draws the menu on the wooden board: the title,
// then each item as a plank with its icon at the left. The pick is the arrow
// plank, the other items a bit darker, disabled ones dark with gray text.
func (m *Menu) Draw(screen *ebiten.Image) {
	vector.FillRect(screen, 0, 0, float32(m.screenW), float32(m.screenH), menuShade, false)
	drawSliced(screen, uiBoardRect, uiBoardInsets, m.PanelRect(), menuBoardScale, nil, 1)
	cx, ty := float64(m.screenW)/2, float64(m.titleCenterY())
	DrawTextCentered(screen, m.Title, menuTitleSize, cx+3, ty+3, menuTitleShadow)
	DrawTextCentered(screen, m.Title, menuTitleSize, cx, ty, menuTitleColor)
	for i, it := range m.Items {
		r := m.ItemRect(i)
		plank, insets := uiPlankRect, uiPlankInsets
		var tint color.Color = menuIdleTint
		clr := menuItemColor
		switch {
		case it.Disabled:
			tint, clr = menuDisabledTint, menuDisabledColor
		case i == m.Selected:
			plank, insets, tint = uiPlankArrowRect, uiPlankArrowInsets, nil
		}
		drawPlank(screen, plank, insets, r, tint)
		if !it.Icon.Empty() {
			drawRotated(screen, uiSprite(it.Icon), m.IconRect(i), 0, tint, 1)
		}
		if it.Bar {
			b := m.barRect(r)
			drawScaled(screen, uiSprite(uiBarEmptyRect), b, tint, 1)
			fill := barFillRect(it.Fill)
			b.Max.X = b.Min.X + fill.Dx()*b.Dx()/uiBarFullRect.Dx()
			drawScaled(screen, uiSprite(fill), b, tint, 1)
			continue
		}
		DrawTextCentered(screen, it.Label, menuItemSize, float64(r.Min.X+r.Max.X)/2,
			float64(r.Min.Y)+menuTextRise*float64(r.Dy()), clr)
	}
}
