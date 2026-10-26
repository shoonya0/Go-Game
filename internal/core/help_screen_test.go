package core

import (
	"image"
	"strings"
	"testing"
)

// rectCenter is the middle of r as a screen point.
func rectCenter(r image.Rectangle) Position {
	return Position{X: float64(r.Min.X+r.Max.X) / 2, Y: float64(r.Min.Y+r.Max.Y) / 2}
}

// openHelp opens How to Play from the title's settings.
func openHelp(t *testing.T) *Screens {
	t.Helper()
	s := NewScreens(1360, 768, true)
	s.OpenTitle(false)
	choose(s, labelSettings)
	s.Update(&InputState{Tapped: true, TapAt: center(s.settings, 1)})
	if s.Current != ScreenHelp {
		t.Fatalf("How to Play: screen %v, want the help screen", s.Current)
	}
	return s
}

// TestHelpOpensFromSettingsAndGoesBack: How to Play opens from the settings,
// on its first page; Esc and its Back button return to the settings.
func TestHelpOpensFromSettingsAndGoesBack(t *testing.T) {
	s := openHelp(t)
	if s.help.Page != 0 {
		t.Errorf("opened on page %d, want 0", s.help.Page)
	}
	s.Update(&InputState{Back: true})
	if s.Current != ScreenSettings {
		t.Errorf("Esc: screen %v, want the settings", s.Current)
	}

	s.Update(&InputState{Tapped: true, TapAt: center(s.settings, 1)})
	b := s.help.BackRect()
	s.Update(&InputState{Clicked: true, Pointer: rectCenter(b)})
	if s.Current != ScreenSettings {
		t.Errorf("Back button: screen %v, want the settings", s.Current)
	}
}

// TestHelpPages: right/left (or D/A) and the arrow buttons turn the pages,
// stopping at the first and the last.
func TestHelpPages(t *testing.T) {
	s := openHelp(t)
	last := len(helpPages) - 1
	for range last + 2 {
		s.Update(&InputState{RightJustPressed: true})
	}
	if s.help.Page != last {
		t.Errorf("right past the end: page %d, want %d", s.help.Page, last)
	}
	s.Update(&InputState{LeftJustPressed: true})
	if s.help.Page != last-1 {
		t.Errorf("left: page %d, want %d", s.help.Page, last-1)
	}
	l := s.help.ArrowRect(-1)
	s.Update(&InputState{Tapped: true, TapAt: rectCenter(l)})
	if s.help.Page != last-2 {
		t.Errorf("tap on the left arrow: page %d, want %d", s.help.Page, last-2)
	}
	r := s.help.ArrowRect(1)
	s.Update(&InputState{Tapped: true, TapAt: rectCenter(r)})
	if s.help.Page != last-1 {
		t.Errorf("tap on the right arrow: page %d, want %d", s.help.Page, last-1)
	}
}

// TestHelpCoversTheControls: the pages explain every key the game reads.
func TestHelpCoversTheControls(t *testing.T) {
	keys := map[string]bool{}
	for _, p := range helpPages {
		for _, r := range p.Rows {
			for _, k := range r.Keys {
				keys[k] = true
			}
		}
	}
	for _, k := range []string{"W", "A", "S", "D", "Space", "Shift", "E", "F", "Esc", "F3", "J", "I", "K", "L", "U", "O"} {
		if !keys[k] {
			t.Errorf("no row explains %s", k)
		}
	}
}

// TestHelpFitsTheBoard: every page's rows fit inside the board, so no text
// runs off it.
func TestHelpFitsTheBoard(t *testing.T) {
	h := NewHelpScreen(1360, 768)
	board := h.PanelRect()
	for i, p := range helpPages {
		h.Page = i
		for j, r := range p.Rows {
			if box := h.rowRect(j, r); !box.In(board) {
				t.Errorf("page %q row %d (%s) at %v runs off the board %v", p.Name, j, strings.Join(r.Keys, "+"), box, board)
			}
		}
	}
	if !h.BackRect().In(board) {
		t.Errorf("Back button %v runs off the board %v", h.BackRect(), board)
	}
}

// TestHelpPowersUseTheGameNumbers: the powers page reads its damage and
// costs from the game's constants, so retuning updates the page.
func TestHelpPowersUseTheGameNumbers(t *testing.T) {
	var u string
	for _, p := range helpPages {
		for _, r := range p.Rows {
			if len(r.Keys) == 1 && r.Keys[0] == "U" {
				u = r.Text
			}
		}
	}
	if want := bars(GroundWaveDamage); !strings.Contains(u, want) {
		t.Errorf("U row %q doesn't give its damage %q", u, want)
	}
}
