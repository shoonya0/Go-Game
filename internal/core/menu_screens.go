package core

import "github.com/hajimehoshi/ebiten/v2"

// Screen is the menu screen on top of the game, or ScreenNone while playing.
type Screen int

const (
	ScreenNone       Screen = iota // playing: no menu
	ScreenTitle                    // the main screen the game opens on
	ScreenConfirmNew               // New Game with a save: overwrite it? (main-menu D4)
)

// MenuAction is what a menu choice asks the game to do.
type MenuAction int

const (
	ActionNone     MenuAction = iota
	ActionNewGame             // wipe the save, start in the start portal
	ActionContinue            // start from the save
	ActionQuit                // close the game (desktop only, D6)
)

// The title screen's labels; Quit is left out where the game can't close (D6).
const (
	labelNewGame  = "New Game"
	labelContinue = "Continue"
	labelSettings = "Settings"
	labelQuit     = "Quit"
)

// Screens runs the menus: which one is up, and what a choice in it means.
type Screens struct {
	Current Screen

	title, confirm *Menu
}

// NewScreens lays the menus out for a w×h screen, showing Quit only when
// canQuit. They open on the title screen.
func NewScreens(w, h int, canQuit bool) *Screens {
	labels := []string{labelNewGame, labelContinue, labelSettings}
	if canQuit {
		labels = append(labels, labelQuit)
	}
	return &Screens{
		Current: ScreenTitle,
		title:   NewMenu(w, h, "PIRATE ADVENTURE", labels...),
		confirm: NewMenu(w, h, "OVERWRITE SAVE?", "No", "Yes"),
	}
}

// OpenTitle shows the title screen, with Continue grayed when there is no
// save. The pick starts on Continue when it can be chosen.
func (s *Screens) OpenTitle(hasSave bool) {
	s.Current = ScreenTitle
	s.title.SetDisabled(1, !hasSave)
	s.title.Selected = 0
	if hasSave {
		s.title.Selected = 1
	}
}

// HasSave reports whether the title screen offers Continue.
func (s *Screens) HasSave() bool { return !s.title.Items[1].Disabled }

// Update applies one tick of input to the current menu and returns the
// action its choice asks for.
func (s *Screens) Update(in *InputState) MenuAction {
	switch s.Current {
	case ScreenTitle:
		i, _ := s.title.Update(in)
		if i < 0 {
			return ActionNone
		}
		switch s.title.Items[i].Label {
		case labelNewGame:
			if s.HasSave() {
				s.Current, s.confirm.Selected = ScreenConfirmNew, 0
				return ActionNone
			}
			return s.play(ActionNewGame)
		case labelContinue:
			return s.play(ActionContinue)
		case labelQuit:
			return ActionQuit
		}
	case ScreenConfirmNew:
		i, back := s.confirm.Update(in)
		switch {
		case back || i == 0:
			s.Current = ScreenTitle
		case i == 1:
			return s.play(ActionNewGame)
		}
	}
	return ActionNone
}

// play closes the menus for action a.
func (s *Screens) play(a MenuAction) MenuAction {
	s.Current = ScreenNone
	return a
}

// Draw draws the current menu over the game.
func (s *Screens) Draw(screen *ebiten.Image) {
	switch s.Current {
	case ScreenTitle:
		s.title.Draw(screen)
	case ScreenConfirmNew:
		s.confirm.Draw(screen)
	}
}
