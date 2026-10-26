package core

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

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
	s := &Screens{
		Current:  ScreenTitle,
		title:    NewMenu(w, h, "PIRATE ADVENTURE", labels...),
		confirm:  NewMenu(w, h, "OVERWRITE SAVE?", "No", "Yes"),
		settings: NewMenu(w, h, "SETTINGS", "", "Back"),
		pause:    NewMenu(w, h, "PAUSED", "Resume", labelSettings, "Main Menu"),
	}
	s.title.SetIcons(uiRestartRect, uiRightRect, uiSettingsRect, uiCrossRect)
	s.confirm.SetIcons(uiCrossRect, uiCheckRect)
	s.settings.SetIcons(image.Rectangle{}, uiLeftRect) // the volume's is set with it
	s.pause.SetIcons(uiRightRect, uiSettingsRect, uiHomeRect)
	s.settings.Items[0].Bar = true
	s.SetSettings(DefaultSettings())
	return s
}

// OpenPause pauses the game under the pause menu, its pick on Resume.
func (s *Screens) OpenPause() {
	s.Current = ScreenPause
	s.pause.Open(0)
}

// SetSettings shows s on the settings screen: the volume as a bar, with the
// speaker crossed out at 0.
func (s *Screens) SetSettings(st Settings) {
	s.Settings = st
	vol := &s.settings.Items[0]
	vol.Fill = float64(st.Volume) / MaxVolume
	vol.Icon = uiSoundRect
	if st.Volume == 0 {
		vol.Icon = uiMuteRect
	}
	s.settings.SetLabel(0, volumeLabel(st.Volume))
}

// openSettings shows the settings screen, returning to back when it closes.
func (s *Screens) openSettings(back Screen) {
	s.Current, s.back = ScreenSettings, back
	s.settings.Open(0)
}

// OpenTitle shows the title screen, with Continue grayed when there is no
// save. The pick starts on Continue when it can be chosen.
func (s *Screens) OpenTitle(hasSave bool) {
	s.Current = ScreenTitle
	s.title.SetDisabled(1, !hasSave)
	sel := 0
	if hasSave {
		sel = 1
	}
	s.title.Open(sel)
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
				s.Current = ScreenConfirmNew
				s.confirm.Open(0)
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
	at, pressed := in.PressAt() // Enter on the volume changes nothing
	step := 0
	switch {
	case back || i == 1:
		s.Current = s.back
		return ActionNone
	case i == 0 && pressed:
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
