package core

import (
	"bytes"
	"image"
	"testing"
)

// TestEnemyClipsFitTheirSheets checks every clip's frames lie inside its sheet.
// A frame past the edge draws nothing, so the enemy blinks out for that frame:
// the succubus patrol clip once ran one column past its 8-column row.
func TestEnemyClipsFitTheirSheets(t *testing.T) {
	for _, bp := range enemyBlueprints {
		for st, c := range bp.Clips {
			cfg, _, err := image.DecodeConfig(bytes.NewReader(readAsset(c.Sheet)))
			if err != nil {
				t.Errorf("%s state %d: decode %q: %v", bp.Name, st, c.Sheet, err)
				continue
			}
			if c.Frames <= 0 || c.FrameW <= 0 || c.FrameH <= 0 || c.FPS <= 0 {
				t.Errorf("%s state %d: frames=%d size=%dx%d fps=%v, want all positive", bp.Name, st, c.Frames, c.FrameW, c.FrameH, c.FPS)
				continue
			}
			if right := (c.Start + c.Frames) * c.FrameW; right > cfg.Width {
				t.Errorf("%s state %d: frames end at x=%d, past the %dpx-wide sheet %q", bp.Name, st, right, cfg.Width, c.Sheet)
			}
			if bottom := (c.Row + 1) * c.FrameH; bottom > cfg.Height {
				t.Errorf("%s state %d: row ends at y=%d, past the %dpx-high sheet %q", bp.Name, st, bottom, cfg.Height, c.Sheet)
			}
		}
	}
}
