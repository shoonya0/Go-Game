package core

import (
	"image"
	"testing"
)

func choose(s *Screens, label string) MenuAction {
	m := s.title
	if s.Current == ScreenConfirmNew {
		m = s.confirm
	}
	for i, it := range m.Items {
		if it.Label == label {
			return s.Update(&InputState{Tapped: true, TapAt: center(m, i)})
		}
	}
	return ActionNone
}

// TestMenuScreensNoSave: without a save Continue is grayed and New Game
// starts right away, closing the menus.
func TestMenuScreensNoSave(t *testing.T) {
	s := NewScreens(1360, 768, true)
	s.OpenTitle(false)
	if a := choose(s, labelContinue); a != ActionNone || s.Current != ScreenTitle {
		t.Errorf("Continue without a save: %v on screen %v, want nothing", a, s.Current)
	}
	if got, _ := s.title.Update(&InputState{Confirm: true}); got != 0 {
		t.Errorf("Enter on the fresh title: chose %d, want New Game (0)", got)
	}
	if a := choose(s, labelNewGame); a != ActionNewGame || s.Current != ScreenNone {
		t.Errorf("New Game: %v on screen %v, want new game, menus closed", a, s.Current)
	}
}

// TestMenuScreensConfirmOverwrite: with a save, the pick starts on Continue;
// New Game asks first: No (or Esc) goes back, Yes starts a new game.
func TestMenuScreensConfirmOverwrite(t *testing.T) {
	s := NewScreens(1360, 768, true)
	s.OpenTitle(true)
	if s.title.Selected != 1 {
		t.Errorf("with a save the pick is %d, want Continue (1)", s.title.Selected)
	}
	if a := choose(s, labelNewGame); a != ActionNone || s.Current != ScreenConfirmNew {
		t.Fatalf("New Game with a save: %v on screen %v, want the confirm screen", a, s.Current)
	}
	if got := s.Update(&InputState{Confirm: true}); got != ActionNone || s.Current != ScreenTitle {
		t.Errorf("Enter on the confirm screen: %v on %v, want No picked first: back to the title", got, s.Current)
	}
	choose(s, labelNewGame)
	if a := s.Update(&InputState{Back: true}); a != ActionNone || s.Current != ScreenTitle {
		t.Errorf("Esc on the confirm screen: %v on %v, want back to the title", a, s.Current)
	}
	choose(s, labelNewGame)
	if a := choose(s, "Yes"); a != ActionNewGame || s.Current != ScreenNone {
		t.Errorf("Yes: %v on %v, want a new game", a, s.Current)
	}
}

// TestMenuScreensContinueAndQuit: Continue starts from the save; Quit is
// there only when the game can close, and asks to quit.
func TestMenuScreensContinueAndQuit(t *testing.T) {
	s := NewScreens(1360, 768, true)
	s.OpenTitle(true)
	if a := choose(s, labelQuit); a != ActionQuit {
		t.Errorf("Quit: %v, want quit", a)
	}
	if a := choose(s, labelContinue); a != ActionContinue || s.Current != ScreenNone {
		t.Errorf("Continue: %v on %v, want continue, menus closed", a, s.Current)
	}
	web := NewScreens(1360, 768, false)
	for _, it := range web.title.Items {
		if it.Label == labelQuit {
			t.Error("web title screen lists Quit, want it hidden")
		}
	}
}

// TestMenuScreensPause: the pause menu resumes on Resume, Esc or the menu
// button; its Settings return to it; Main Menu asks for the title screen.
func TestMenuScreensPause(t *testing.T) {
	s := NewScreens(1360, 768, true)
	for name, in := range map[string]InputState{
		"Resume":      {Confirm: true},
		"Esc":         {Back: true, Menu: true},
		"menu button": {Menu: true},
	} {
		s.OpenPause()
		if a := s.Update(&in); a != ActionNone || s.Current != ScreenNone {
			t.Errorf("%s: %v on screen %v, want playing again", name, a, s.Current)
		}
	}

	s.OpenPause()
	s.Update(&InputState{DownJustPressed: true})
	s.Update(&InputState{Confirm: true})
	if s.Current != ScreenSettings {
		t.Fatalf("Settings: screen %v, want the settings screen", s.Current)
	}
	s.Update(&InputState{Back: true})
	if s.Current != ScreenPause {
		t.Errorf("Esc in settings opened from pause: screen %v, want the pause menu", s.Current)
	}
	s.Update(&InputState{DownJustPressed: true})
	if a := s.Update(&InputState{Confirm: true}); a != ActionToTitle {
		t.Errorf("Main Menu: %v, want to the title", a)
	}
}

// TestMenuScreensRestingCursorKeepsDefault: a mouse resting where the last
// screen's pick was doesn't move the next screen's default pick. Settings
// and Yes share a row, so a resting cursor once made Enter wipe the save.
func TestMenuScreensRestingCursorKeepsDefault(t *testing.T) {
	s := NewScreens(1360, 768, true)
	s.OpenTitle(true)
	rest := center(s.title, 2) // on Settings, the row Yes takes on the confirm screen
	if yes := s.confirm.ItemRect(1); !(image.Point{X: int(rest.X), Y: int(rest.Y)}).In(yes) {
		t.Fatalf("test setup: cursor %v is not on Yes %v", rest, yes)
	}
	s.Update(&InputState{Pointer: Position{X: 1, Y: 1}})
	s.Update(&InputState{Pointer: rest})                      // the cursor moves onto Settings
	s.Update(&InputState{Pointer: rest, UpJustPressed: true}) // keys up to Continue
	s.Update(&InputState{Pointer: rest, UpJustPressed: true}) // and New Game
	s.Update(&InputState{Pointer: rest, Confirm: true})       // the confirm screen opens
	if a := s.Update(&InputState{Pointer: rest, Confirm: true}); a != ActionNone || s.Current != ScreenTitle {
		t.Errorf("Enter on the confirm screen with the cursor resting on Yes: %v on %v; want No: back to the title", a, s.Current)
	}

	s.Update(&InputState{Pointer: rest, Confirm: true}) // the pause menu, opened over a resting cursor:
	s.OpenPause()
	s.Update(&InputState{Pointer: center(s.pause, 2)}) // resting on Main Menu when it opens
	if a := s.Update(&InputState{Pointer: center(s.pause, 2), Confirm: true}); a != ActionNone || s.Current != ScreenNone {
		t.Errorf("Enter on the pause menu with the cursor resting on Main Menu: %v on %v; want Resume", a, s.Current)
	}

	s.OpenPause()
	s.Update(&InputState{Pointer: Position{X: 1, Y: 1}})
	s.Update(&InputState{Pointer: center(s.pause, 2)}) // a moving cursor still picks
	if s.pause.Selected != 2 {
		t.Errorf("cursor moved onto Main Menu: pick %d, want 2", s.pause.Selected)
	}
}
