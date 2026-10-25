package core

import "testing"

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
