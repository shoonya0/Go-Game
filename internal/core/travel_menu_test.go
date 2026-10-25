package core

import (
	"image"
	"testing"
)

// TestHoldEOpensTravelMenu: holding E in a portal asks for the travel menu
// once, on the tick the hold reaches travelHoldSeconds; so does a touch held
// on the portal. A hold away from every portal, or a touch held beside the
// portal, asks nothing.
func TestHoldEOpensTravelMenu(t *testing.T) {
	level, p := twoPortals(t)
	hold := travelHoldTicks()

	p.Pos = Position{X: 2.5 * LevelTileWidth, Y: level.Portals[1].StandPos().Y} // between the portals
	p.updatePortals(level.Portals, &InputState{InteractTicks: hold})
	if p.TravelAsked() {
		t.Fatal("holding E away from portals asked for the menu")
	}

	p.Pos = level.Portals[1].StandPos()
	for ticks, want := range map[int]bool{hold - 1: false, hold: true, hold + 1: false} {
		p.updatePortals(level.Portals, &InputState{InteractTicks: ticks})
		if p.TravelAsked() != want {
			t.Errorf("E held %d ticks (menu at %d): asked=%v, want %v", ticks, hold, p.TravelAsked(), want)
		}
	}

	p.Camera.Pos = Position{X: 100, Y: -40}
	b := level.Portals[1].Bounds()
	onPortal := Position{X: centerX(b) - p.Camera.Pos.X, Y: b.Y + b.Height/2 - p.Camera.Pos.Y}
	p.updatePortals(level.Portals, &InputState{TouchHoldTicks: hold, TapAt: Position{X: onPortal.X + 4*LevelTileWidth, Y: onPortal.Y}})
	if p.TravelAsked() {
		t.Error("a touch held beside the portal asked for the menu")
	}
	p.updatePortals(level.Portals, &InputState{TouchHoldTicks: hold, TapAt: onPortal})
	if !p.TravelAsked() {
		t.Error("a touch held on the portal did not ask for the menu")
	}
}

// TestTouchHoldFollowsTheFinger: a touch started off the buttons counts the
// ticks it stays down, following a finger that drifts a little; lifting it
// ends the hold.
func TestTouchHoldFollowsTheFinger(t *testing.T) {
	tc := NewTouchControls(1360, 768)
	at := image.Pt(680, 300)
	tc.Update([]image.Point{at})
	tc.Started([]image.Point{at})
	for i := 1; i <= 3; i++ {
		at = at.Add(image.Pt(3, 2))
		tc.Update([]image.Point{at})
		tc.Started(nil)
	}
	var in InputState
	tc.Apply(&in)
	if in.TouchHoldTicks != 3 || in.TapAt != (Position{X: float64(at.X), Y: float64(at.Y)}) {
		t.Errorf("held %d ticks at %v, want 3 at %v", in.TouchHoldTicks, in.TapAt, at)
	}

	tc.Update(nil)
	tc.Started(nil)
	in = InputState{}
	tc.Apply(&in)
	if in.TouchHoldTicks != 0 {
		t.Errorf("after lifting the finger: held %d ticks, want 0", in.TouchHoldTicks)
	}
}

// TestTravelMenuListsActivePortals: the menu lists the active portals in the
// order they were activated, starting on the portal the player stands in.
// Down moves the pick, E travels to it; Esc closes without travelling; a tap
// on a row travels there, a tap outside the rows closes.
func TestTravelMenuListsActivePortals(t *testing.T) {
	level, p := twoPortals(t)
	p.Pos = level.Portals[1].StandPos()
	p.updatePortals(level.Portals, &InputState{Interact: true})

	m := NewTravelMenu(1360, 768)
	m.OpenFor(p)
	if !m.Open || len(m.Rows) != 2 || m.Rows[0] != p.ActivePortals[0] || m.Rows[1] != p.ActivePortals[1] {
		t.Fatalf("open=%v rows %v, want open with %v", m.Open, m.Rows, p.ActivePortals)
	}
	if m.Selected != 1 {
		t.Errorf("selected row %d, want 1 (the portal the player stands in)", m.Selected)
	}

	m.Update(&InputState{UpJustPressed: true})
	m.Update(&InputState{UpJustPressed: true}) // stops at the top
	if m.Selected != 0 {
		t.Errorf("after up twice: row %d, want 0", m.Selected)
	}
	m.Update(&InputState{DownJustPressed: true})
	if pick, ok := m.Update(&InputState{Interact: true}); !ok || pick != p.ActivePortals[1] || m.Open {
		t.Errorf("E: pick %v ok=%v open=%v, want %v and the menu closed", pick, ok, m.Open, p.ActivePortals[1])
	}

	m.OpenFor(p)
	if _, ok := m.Update(&InputState{Back: true}); ok || m.Open {
		t.Errorf("Esc: travelled=%v open=%v, want closed without travelling", ok, m.Open)
	}

	m.OpenFor(p)
	r := m.RowRect(0)
	tap := Position{X: float64(r.Min.X + r.Dx()/2), Y: float64(r.Min.Y + r.Dy()/2)}
	if pick, ok := m.Update(&InputState{Tapped: true, TapAt: tap}); !ok || pick != p.ActivePortals[0] || m.Open {
		t.Errorf("tap on row 0: pick %v ok=%v open=%v, want %v and the menu closed", pick, ok, m.Open, p.ActivePortals[0])
	}

	m.OpenFor(p)
	if _, ok := m.Update(&InputState{Tapped: true, TapAt: Position{X: 5, Y: 5}}); ok || m.Open {
		t.Errorf("tap outside: travelled=%v open=%v, want closed without travelling", ok, m.Open)
	}
}

// TestTravelMenuRowsFitTheScreen: however many portals are active, the rows
// shown fit on the screen, and the selected one is always among them.
func TestTravelMenuRowsFitTheScreen(t *testing.T) {
	m := NewTravelMenu(1360, 768)
	screen := image.Rect(0, 0, 1360, 768)
	m.Rows = make([]PortalRef, 9)
	m.Open = true
	for sel := range m.Rows {
		m.Selected = sel
		first, last := m.VisibleRows()
		if sel < first || sel >= last {
			t.Errorf("selected %d not among visible rows [%d, %d)", sel, first, last)
		}
		for i := first; i < last; i++ {
			if !m.RowRect(i).In(screen) {
				t.Errorf("selected %d: row %d at %v is off screen", sel, i, m.RowRect(i))
			}
		}
	}
}
