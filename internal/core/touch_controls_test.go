package core

import (
	"image"
	"testing"
)

func touchCenter(tc *TouchControls, b TouchButton) image.Point {
	r := tc.Rect(b)
	return image.Pt((r.Min.X+r.Max.X)/2, (r.Min.Y+r.Max.Y)/2)
}

func TestTouchLayoutFitsScreenWithoutOverlap(t *testing.T) {
	const w, h = 1360, 768
	tc := NewTouchControls(w, h)
	screen := image.Rect(0, 0, w, h)
	for a := TouchButton(0); a < touchButtonCount; a++ {
		if !tc.Rect(a).In(screen) {
			t.Errorf("button %d at %v is off screen", a, tc.Rect(a))
		}
		for b := a + 1; b < touchButtonCount; b++ {
			if tc.Rect(a).Overlaps(tc.Rect(b)) {
				t.Errorf("buttons %d and %d overlap", a, b)
			}
		}
	}
	// Movement sits bottom-left, attacks bottom-right, menu top-left.
	if c := touchCenter(tc, TouchLeft); c.X > w/2 || c.Y < h/2 {
		t.Errorf("A button at %v, want bottom-left", c)
	}
	if c := touchCenter(tc, TouchSpecial2); c.X < w/2 || c.Y < h/2 {
		t.Errorf("L button at %v, want bottom-right", c)
	}
	if s, o := tc.Rect(TouchSprint), tc.Rect(TouchSpecial4); s.Max.Y > o.Min.Y || s.Max.X != o.Max.X || s.Dx() <= o.Dx() {
		t.Errorf("sprint at %v, want a bigger button above the O button at %v, right edges aligned", s, o)
	}
	if c := touchCenter(tc, TouchMenu); c.X > w/4 || c.Y > h/4 {
		t.Errorf("menu button at %v, want top-left", c)
	}
}

func TestTouchDisabledLeavesKeyboardInput(t *testing.T) {
	tc := NewTouchControls(1360, 768)
	tc.Update(nil)
	in := InputState{Direction: Direction{LeftRight: 1}}
	tc.Apply(&in)
	if tc.Enabled || in.Direction.LeftRight != 1 {
		t.Fatalf("enabled=%v dir=%d, want controls off and keyboard input kept", tc.Enabled, in.Direction.LeftRight)
	}
}

func TestTouchMovementAndHeldAttacks(t *testing.T) {
	tc := NewTouchControls(1360, 768)
	tc.Update([]image.Point{touchCenter(tc, TouchLeft), touchCenter(tc, TouchUp), touchCenter(tc, TouchWeak)})

	var in InputState
	tc.Apply(&in)
	if !tc.Enabled {
		t.Fatal("a touch should enable the controls")
	}
	if in.Direction.LeftRight != -1 || in.Direction.UpDown != 1 {
		t.Errorf("direction = %+v, want left and up", in.Direction)
	}
	if !in.JumpJustPressed {
		t.Error("W should jump")
	}
	if !in.Skills.WeakAttack || in.Skills.StrongAttack {
		t.Errorf("skills = %+v, want only the weak attack (J)", in.Skills)
	}
}

func TestTouchOppositeDirectionsCancel(t *testing.T) {
	tc := NewTouchControls(1360, 768)
	tc.Update([]image.Point{touchCenter(tc, TouchLeft), touchCenter(tc, TouchRight)})
	var in InputState
	tc.Apply(&in)
	if in.Direction.LeftRight != 0 {
		t.Errorf("LeftRight = %d, want 0 with A and D both held", in.Direction.LeftRight)
	}
}

// K, L and O read as just-pressed, like their keys, so holding them does not
// re-trigger the combo; the menu button too, once per press (it pauses).
func TestTouchJustPressedButtons(t *testing.T) {
	tc := NewTouchControls(1360, 768)
	held := []image.Point{touchCenter(tc, TouchSpecial1), touchCenter(tc, TouchSpecial2), touchCenter(tc, TouchSpecial4), touchCenter(tc, TouchMenu)}

	var in InputState
	tc.Update(held)
	tc.Apply(&in)
	if !in.Skills.SpecialAttack1 || !in.Skills.SpecialAttack2 || !in.Skills.SpecialAttack4 || !in.Menu {
		t.Fatalf("first frame: K=%v L=%v O=%v menu=%v, want all true", in.Skills.SpecialAttack1, in.Skills.SpecialAttack2, in.Skills.SpecialAttack4, in.Menu)
	}

	in = InputState{} // HandleInput rewrites every field each frame
	tc.Update(held)
	tc.Apply(&in)
	if in.Skills.SpecialAttack1 || in.Skills.SpecialAttack2 || in.Skills.SpecialAttack4 || in.Menu {
		t.Fatalf("held: K=%v L=%v O=%v menu=%v, want all false", in.Skills.SpecialAttack1, in.Skills.SpecialAttack2, in.Skills.SpecialAttack4, in.Menu)
	}

	in = InputState{}
	tc.Update(nil)
	tc.Apply(&in)
	tc.Update(held[3:])
	tc.Apply(&in)
	if !in.Menu {
		t.Error("second menu press: menu=false, want true (each press pauses or resumes)")
	}
}

// Every button but the menu (whose sprite has its own icon) shows a symbol,
// and the symbol stays inside the button.
func TestTouchIconsCoverButtons(t *testing.T) {
	for b := TouchButton(0); b < touchButtonCount; b++ {
		shapes := touchIconShapes[b]
		if (b == TouchMenu) != (len(shapes) == 0) {
			t.Errorf("button %d has %d icon shapes", b, len(shapes))
		}
		for _, sh := range shapes {
			for _, p := range sh.pts {
				if p[0] < 0.1 || p[0] > 0.9 || p[1] < 0.1 || p[1] > 0.9 {
					t.Errorf("button %d icon point %v leaves the button face", b, p)
				}
			}
		}
	}
}

// Sprint stands in for Left Shift, which both runs and dashes.
func TestTouchSprintRunsAndDashes(t *testing.T) {
	tc := NewTouchControls(1360, 768)
	tc.Update([]image.Point{touchCenter(tc, TouchSprint), touchCenter(tc, TouchRight)})
	var in InputState
	tc.Apply(&in)
	if !in.RunJustPressed || !in.DashJustPressed || in.Direction.LeftRight != 1 {
		t.Errorf("run=%v dash=%v dir=%d, want run and dash while moving right", in.RunJustPressed, in.DashJustPressed, in.Direction.LeftRight)
	}
}

// UpdateCamera keeps the player in the middle third of the screen, so no
// button may reach into it.
func TestTouchButtonsStayOutOfPlayerArea(t *testing.T) {
	const w, h = 1360, 768
	tc := NewTouchControls(w, h)
	middle := image.Rect(w/3, 0, 2*w/3, h)
	for b := TouchButton(0); b < touchButtonCount; b++ {
		if tc.Rect(b).Overlaps(middle) {
			t.Errorf("button %d at %v covers the middle third %v", b, tc.Rect(b), middle)
		}
	}
}

// TestTouchTapOffButtons: a touch that starts on no button is a tap (at its
// screen point) on that frame only, even while a button is held; touches that
// start on a button, or are only held, never are.
func TestTouchTapOffButtons(t *testing.T) {
	tc := NewTouchControls(1360, 768)
	free := image.Pt(680, 300)
	weak := touchCenter(tc, TouchWeak)

	var in InputState
	tc.Update([]image.Point{weak, free})
	tc.Started([]image.Point{free})
	tc.Apply(&in)
	if !in.Tapped || in.TapAt != (Position{X: 680, Y: 300}) || !tc.held[TouchWeak] {
		t.Fatalf("new touch off the buttons beside a held J: tapped=%v at %v, want a tap at (680, 300)", in.Tapped, in.TapAt)
	}

	tc.Update([]image.Point{weak, free})
	tc.Started(nil)
	tc.Apply(&in)
	if in.Tapped {
		t.Error("a held touch tapped again")
	}

	tc.Update([]image.Point{weak})
	tc.Started([]image.Point{weak})
	tc.Apply(&in)
	if in.Tapped {
		t.Error("a button press counted as a tap")
	}
}

// TestTouchHiddenButtons: while a menu hides the buttons, a touch where a
// button would be is a tap at that point and presses nothing (an unseen
// down arrow used to move the menu's pick).
func TestTouchHiddenButtons(t *testing.T) {
	tc := NewTouchControls(1360, 768)
	tc.Enabled, tc.Hidden = true, true
	at := touchCenter(tc, TouchDown)
	tc.Update([]image.Point{at})
	tc.Started([]image.Point{at})
	var in InputState
	tc.Apply(&in)
	if in.DownJustPressed || in.Direction.UpDown != 0 || !in.Tapped || in.TapAt != (Position{X: float64(at.X), Y: float64(at.Y)}) {
		t.Errorf("touch on the hidden down button: down=%v updown=%d tapped=%v at %v; want only a tap at %v",
			in.DownJustPressed, in.Direction.UpDown, in.Tapped, in.TapAt, at)
	}
}
