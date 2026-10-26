package core

import (
	"image"
	"slices"
	"testing"
)

func testMenu() *Menu {
	return NewMenu(1360, 768, "PIRATE ADVENTURE", "New Game", "Continue", "Settings", "Quit")
}

// center returns the middle of item i's box, as a screen point.
func center(m *Menu, i int) Position {
	r := m.ItemRect(i)
	return Position{X: float64(r.Min.X+r.Max.X) / 2, Y: float64(r.Min.Y+r.Max.Y) / 2}
}

// TestMenuKeysSkipDisabled: down/up move the pick, wrap around and jump over
// a disabled item; Enter or E chooses the pick.
func TestMenuKeysSkipDisabled(t *testing.T) {
	m := testMenu()
	m.SetDisabled(1, true) // no save: Continue grayed
	var picks []int
	for range 4 {
		m.Update(&InputState{DownJustPressed: true})
		picks = append(picks, m.Selected)
	}
	if want := []int{2, 3, 0, 2}; !slices.Equal(picks, want) {
		t.Errorf("down ×4 picks %v, want %v (Continue skipped, wrapping)", picks, want)
	}
	m.Update(&InputState{UpJustPressed: true})
	m.Update(&InputState{UpJustPressed: true})
	if m.Selected != 3 {
		t.Errorf("up ×2 from Settings: pick %d, want Quit (3) after wrapping", m.Selected)
	}
	if got, _ := m.Update(&InputState{Confirm: true}); got != 3 {
		t.Errorf("Enter: chose %d, want 3", got)
	}
	if got, _ := m.Update(&InputState{Interact: true}); got != 3 {
		t.Errorf("E: chose %d, want 3", got)
	}
}

// TestMenuDisabledFirstItem: disabling the picked item moves the pick to the
// next enabled one, so Enter never chooses a disabled item.
func TestMenuDisabledFirstItem(t *testing.T) {
	m := NewMenu(1360, 768, "T", "Continue", "New Game")
	m.SetDisabled(0, true)
	if got, _ := m.Update(&InputState{Confirm: true}); got != 1 {
		t.Errorf("Enter: chose %d, want New Game (1)", got)
	}
}

// TestMenuTapAndClick: a tap or a click on an item chooses it; on a disabled
// item or off every item nothing is chosen. Moving the mouse over an item
// picks it without choosing.
func TestMenuTapAndClick(t *testing.T) {
	m := testMenu()
	m.SetDisabled(1, true)
	if got, _ := m.Update(&InputState{Tapped: true, TapAt: center(m, 2)}); got != 2 {
		t.Errorf("tap on Settings: chose %d, want 2", got)
	}
	if got, _ := m.Update(&InputState{Clicked: true, Pointer: center(m, 3)}); got != 3 {
		t.Errorf("click on Quit: chose %d, want 3", got)
	}
	for name, in := range map[string]InputState{
		"tap on disabled":  {Tapped: true, TapAt: center(m, 1)},
		"click off items":  {Clicked: true, Pointer: Position{X: 5, Y: 5}},
		"tap off items":    {Tapped: true, TapAt: Position{X: 5, Y: 760}},
		"mouse moved only": {Pointer: center(m, 0)},
	} {
		if got, _ := m.Update(&in); got != -1 {
			t.Errorf("%s: chose %d, want none", name, got)
		}
	}
	if m.Selected != 0 {
		t.Errorf("mouse over New Game: pick %d, want 0", m.Selected)
	}
	m.Update(&InputState{Pointer: center(m, 0), DownJustPressed: true})
	if m.Selected != 2 {
		t.Errorf("mouse still, down pressed: pick %d, want 2 (a still mouse doesn't hold the pick)", m.Selected)
	}
}

// TestMenuBack: Esc asks to go back and chooses nothing.
func TestMenuBack(t *testing.T) {
	m := testMenu()
	if got, back := m.Update(&InputState{Back: true}); got != -1 || !back {
		t.Errorf("Esc: chose %d, back %v; want none, back", got, back)
	}
}

// TestMenuLayout: items are stacked top to bottom without overlapping, all
// the same width, centered, and below the title.
func TestMenuLayout(t *testing.T) {
	m := testMenu()
	for i := range m.Items {
		r := m.ItemRect(i)
		if r.Min.X+r.Max.X != 1360 {
			t.Errorf("item %d %v not centered on x 680", i, r)
		}
		if i > 0 {
			prev := m.ItemRect(i - 1)
			if r.Min.Y <= prev.Max.Y || r.Dx() != prev.Dx() {
				t.Errorf("item %d %v vs %d %v: want below it, same width", i, r, i-1, prev)
			}
		}
	}
	if top := m.ItemRect(0).Min.Y; top < menuTitleSize+menuTitleGap {
		t.Errorf("first item at y %d leaves no room for the title", top)
	}
	if bottom := m.ItemRect(len(m.Items) - 1).Max.Y; bottom > 768 {
		t.Errorf("last item ends at y %d, below the screen", bottom)
	}
}

// TestMenuLayoutWithIcons: with icons, icons and planks are centered
// together, each icon sits left of its plank, and the board holds the title,
// every icon and plank, and fits on the screen.
func TestMenuLayoutWithIcons(t *testing.T) {
	for _, canQuit := range []bool{true, false} {
		s := NewScreens(1360, 768, canQuit)
		for name, m := range map[string]*Menu{"title": s.title, "confirm": s.confirm, "settings": s.settings, "pause": s.pause} {
			panel := m.PanelRect()
			if !panel.In(image.Rect(0, 0, 1360, 768)) {
				t.Errorf("%s: board %v leaves the 1360×768 screen", name, panel)
			}
			for i := range m.Items {
				r, icon := m.ItemRect(i), m.IconRect(i)
				if icon.Max.X > r.Min.X || icon.Dx() != r.Dy() {
					t.Errorf("%s item %d: icon %v not a square left of plank %v", name, i, icon, r)
				}
				if icon.Min.X+r.Max.X != 1360 {
					t.Errorf("%s item %d: icon and plank %v–%v not centered on x 680", name, i, icon.Min.X, r.Max.X)
				}
				if !r.In(panel) || !icon.In(panel) {
					t.Errorf("%s item %d: plank %v or icon %v outside the board %v", name, i, r, icon, panel)
				}
			}
			if title := m.titleCenterY(); title-menuTitleSize/2 < panel.Min.Y {
				t.Errorf("%s: title at y %d sticks out of the board %v", name, title, panel)
			}
		}
	}
}

// TestMenuBarFitsPlank: the volume bar lies on the plank's face, clear of
// the arrow plank's ends.
func TestMenuBarFitsPlank(t *testing.T) {
	m := NewScreens(1360, 768, true).settings
	r := m.ItemRect(0)
	b := m.barRect(r)
	if !b.In(r) {
		t.Fatalf("bar %v outside its plank %v", b, r)
	}
	side := int(float64(uiPlankArrowInsets.right*r.Dy()) / float64(uiPlankArrowRect.Dy()))
	if b.Max.X > r.Max.X-side {
		t.Errorf("bar ends at x %d, over the arrow plank's point (from x %d)", b.Max.X, r.Max.X-side)
	}
}
