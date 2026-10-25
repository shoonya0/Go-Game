package core

import "testing"

// TestSettingsRoundTrip: settings encode and decode to the same volume; bad
// data or a volume out of range falls back to the defaults with an error.
func TestSettingsRoundTrip(t *testing.T) {
	b, err := EncodeSettings(Settings{Volume: 3})
	if err != nil {
		t.Fatal(err)
	}
	if got, err := DecodeSettings(b); err != nil || got.Volume != 3 {
		t.Errorf("decoded %+v (err %v), want volume 3", got, err)
	}
	for _, bad := range []string{"{not json", `{"volume":11}`, `{"volume":-1}`} {
		if got, err := DecodeSettings([]byte(bad)); err == nil || got != DefaultSettings() {
			t.Errorf("%s: %+v (err %v), want the defaults and an error", bad, got, err)
		}
	}
}

// TestSettingsScreenVolume: Settings on the title opens the settings screen;
// right/left change the volume within 0..MaxVolume, each change reported
// once; Esc returns to the title.
func TestSettingsScreenVolume(t *testing.T) {
	s := NewScreens(1360, 768, true)
	s.OpenTitle(false)
	choose(s, labelSettings)
	if s.Current != ScreenSettings {
		t.Fatalf("Settings: screen %v, want the settings screen", s.Current)
	}
	start := s.Settings.Volume
	if a := s.Update(&InputState{RightJustPressed: true}); a != ActionSettingsChanged || s.Settings.Volume != start+1 {
		t.Errorf("right: %v, volume %d; want changed, %d", a, s.Settings.Volume, start+1)
	}
	for range MaxVolume + 2 {
		s.Update(&InputState{RightJustPressed: true})
	}
	if a := s.Update(&InputState{RightJustPressed: true}); a != ActionNone || s.Settings.Volume != MaxVolume {
		t.Errorf("right at the top: %v, volume %d; want nothing, %d", a, s.Settings.Volume, MaxVolume)
	}
	if s.settings.Items[0].Label != volumeLabel(MaxVolume) {
		t.Errorf("label %q, want %q", s.settings.Items[0].Label, volumeLabel(MaxVolume))
	}
	for range MaxVolume + 2 {
		s.Update(&InputState{LeftJustPressed: true})
	}
	if s.Settings.Volume != 0 {
		t.Errorf("left past the bottom: volume %d, want 0", s.Settings.Volume)
	}
	s.Update(&InputState{Back: true})
	if s.Current != ScreenTitle {
		t.Errorf("Esc: screen %v, want the title", s.Current)
	}
}

// TestSettingsScreenTapsAndBack: a tap on the volume's right half turns it
// up, on its left half down; left/right do nothing while Back is picked, and
// choosing Back returns.
func TestSettingsScreenTapsAndBack(t *testing.T) {
	s := NewScreens(1360, 768, true)
	s.openSettings(ScreenTitle)
	r := s.settings.ItemRect(0)
	y := float64(r.Min.Y+r.Max.Y) / 2
	start := s.Settings.Volume
	s.Update(&InputState{Tapped: true, TapAt: Position{X: float64(r.Max.X - 5), Y: y}})
	if s.Settings.Volume != start+1 {
		t.Errorf("tap on the right half: volume %d, want %d", s.Settings.Volume, start+1)
	}
	s.Update(&InputState{Clicked: true, Pointer: Position{X: float64(r.Min.X + 5), Y: y}})
	if s.Settings.Volume != start {
		t.Errorf("click on the left half: volume %d, want %d", s.Settings.Volume, start)
	}
	s.Update(&InputState{DownJustPressed: true})
	if a := s.Update(&InputState{RightJustPressed: true}); a != ActionNone || s.Settings.Volume != start {
		t.Errorf("right on Back: %v, volume %d; want nothing", a, s.Settings.Volume)
	}
	s.Update(&InputState{Confirm: true})
	if s.Current != ScreenTitle {
		t.Errorf("Back: screen %v, want the title", s.Current)
	}
}
