package core

import "github.com/hajimehoshi/ebiten/v2"

// Screen is the menu screen on top of the game, or ScreenNone while playing.
type Screen int

const (
	ScreenNone       Screen = iota // playing: no menu
	ScreenTitle                    // the main screen the game opens on
	ScreenConfirmNew               // New Game with a save: overwrite it? (main-menu D4)
	ScreenSettings                 // the volume (D5)
	ScreenPause                    // Esc while playing: the game paused underneath (D7)
)

// MenuAction is what a menu choice asks the game to do.
type MenuAction int

const (
	ActionNone            MenuAction = iota
	ActionNewGame                    // wipe the save, start in the start portal
	ActionContinue                   // start from the save
	ActionQuit                       // close the game (desktop only, D6)
	ActionSettingsChanged            // Settings changed: store them
	ActionToTitle                    // leave the game for the title screen
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
	Current  Screen
	Settings Settings

	title, confirm, settings, pause *Menu
	back                            Screen // the screen the settings return to
}

// NewScreens lays the menus out for a w×h screen, showing Quit only when
// canQuit. They open on the title screen.
func NewScreens(w, h int, canQuit bool) *Screens {
	labels := []string{labelNewGame, labelContinue, labelSettings}
	if canQuit {
		labels = append(labels, labelQuit)
	}
	return &Screens{
		Current:  ScreenTitle,
		title:    NewMenu(w, h, "PIRATE ADVENTURE", labels...),
		confirm:  NewMenu(w, h, "OVERWRITE SAVE?", "No", "Yes"),
		settings: NewMenu(w, h, "SETTINGS", volumeLabel(DefaultSettings().Volume), "Back"),
		pause:    NewMenu(w, h, "PAUSED", "Resume", labelSettings, "Main Menu"),
		Settings: DefaultSettings(),
	}
}

// OpenPause pauses the game under the pause menu, its pick on Resume.
func (s *Screens) OpenPause() {
	s.Current, s.pause.Selected = ScreenPause, 0
}

// SetSettings shows s on the settings screen.
func (s *Screens) SetSettings(st Settings) {
	s.Settings = st
	s.settings.Items[0].Label = volumeLabel(st.Volume)
}

// openSettings shows the settings screen, returning to back when it closes.
func (s *Screens) openSettings(back Screen) {
	s.Current, s.back, s.settings.Selected = ScreenSettings, back, 0
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
		case labelSettings:
			s.openSettings(ScreenTitle)
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
	case ScreenSettings:
		return s.updateSettings(in)
	case ScreenPause:
		i, back := s.pause.Update(in)
		switch {
		case back || in.Menu || i == 0:
			s.Current = ScreenNone
		case i == 1:
			s.openSettings(ScreenPause)
		case i == 2:
			return ActionToTitle
		}
	}
	return ActionNone
}

// updateSettings runs the settings screen: left/right (or A/D) change the
// volume while it is picked; a tap or click on the volume's left half turns
// it down, on its right half up. Back or Esc returns.
func (s *Screens) updateSettings(in *InputState) MenuAction {
	i, back := s.settings.Update(in)
	step := 0
	switch {
	case back || i == 1:
		s.Current = s.back
		return ActionNone
	case i == 0 && (in.Tapped || in.Clicked):
		at := in.TapAt
		if in.Clicked {
			at = in.Pointer
		}
		r := s.settings.ItemRect(0)
		step = 1
		if int(at.X) < (r.Min.X+r.Max.X)/2 {
			step = -1
		}
	case s.settings.Selected == 0 && in.LeftJustPressed:
		step = -1
	case s.settings.Selected == 0 && in.RightJustPressed:
		step = 1
	}
	v := min(MaxVolume, max(0, s.Settings.Volume+step))
	if v == s.Settings.Volume {
		return ActionNone
	}
	s.SetSettings(Settings{Volume: v})
	return ActionSettingsChanged
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
	case ScreenSettings:
		s.settings.Draw(screen)
	case ScreenPause:
		s.pause.Draw(screen)
	}
}
