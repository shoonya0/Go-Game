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
	if s, i := tc.Rect(TouchSprint), tc.Rect(TouchStrong); s.Max.Y > i.Min.Y || s.Min.X != i.Min.X {
		t.Errorf("sprint at %v, want above the I button at %v", s, i)
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

// K and L read as just-pressed, like their keys, so holding them does not
// re-trigger the combo; the menu button toggles once per press.
func TestTouchJustPressedButtons(t *testing.T) {
	tc := NewTouchControls(1360, 768)
	held := []image.Point{touchCenter(tc, TouchSpecial1), touchCenter(tc, TouchSpecial2), touchCenter(tc, TouchMenu)}

	var in InputState
	tc.Update(held)
	tc.Apply(&in)
	if !in.Skills.SpecialAttack1 || !in.Skills.SpecialAttack2 || !in.Menu {
		t.Fatalf("first frame: K=%v L=%v menu=%v, want all true", in.Skills.SpecialAttack1, in.Skills.SpecialAttack2, in.Menu)
	}

	in.Skills = Skills{}
	tc.Update(held)
	tc.Apply(&in)
	if in.Skills.SpecialAttack1 || in.Skills.SpecialAttack2 || !in.Menu {
		t.Fatalf("held: K=%v L=%v menu=%v, want K/L false and menu still open", in.Skills.SpecialAttack1, in.Skills.SpecialAttack2, in.Menu)
	}

	tc.Update(nil)
	tc.Apply(&in)
	tc.Update(held[2:])
	tc.Apply(&in)
	if in.Menu {
		t.Error("second menu press should close the menu")
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
