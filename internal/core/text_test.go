package core

import "testing"

// TestTextSize: the UI font loads and measures; longer text is wider and a
// bigger size is wider and taller. Press Start 2P is monospaced and square,
// so one letter at 16 px is 16×16.
func TestTextSize(t *testing.T) {
	w1, h1 := TextSize("A", 16)
	if w1 != 16 || h1 != 16 {
		t.Errorf("one letter at 16 px: %v×%v, want 16×16", w1, h1)
	}
	if w, _ := TextSize("AAAA", 16); w != 4*w1 {
		t.Errorf("four letters: width %v, want %v", w, 4*w1)
	}
	if w, h := TextSize("A", 32); w <= w1 || h <= h1 {
		t.Errorf("at 32 px: %v×%v, want bigger than %v×%v", w, h, w1, h1)
	}
}
