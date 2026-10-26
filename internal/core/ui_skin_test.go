package core

import (
	"bytes"
	"image"
	"math"
	"testing"
)

// TestUISkinRectsOnSheet: every sprite rect lies inside the sheet and its
// slice insets fit inside the sprite.
func TestUISkinRectsOnSheet(t *testing.T) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(readAsset(UIButtonsSheet)))
	if err != nil {
		t.Fatal(err)
	}
	sheet := image.Rect(0, 0, cfg.Width, cfg.Height)
	rects := map[string]image.Rectangle{
		"menu": uiMenuRect, "roundUp": uiRoundUpRect, "roundDown": uiRoundDownRect,
		"plank": uiPlankRect, "plankArrow": uiPlankArrowRect, "board": uiBoardRect,
		"barFull": uiBarFullRect, "barEmpty": uiBarEmptyRect,
		"restart": uiRestartRect, "settings": uiSettingsRect, "home": uiHomeRect,
		"sound": uiSoundRect, "mute": uiMuteRect, "left": uiLeftRect, "right": uiRightRect,
		"cross": uiCrossRect, "check": uiCheckRect,
	}
	for name, r := range rects {
		if r.Empty() || !r.In(sheet) {
			t.Errorf("%s %v is not a sprite inside the %v sheet", name, r, sheet)
		}
	}
	for name, s := range map[string]struct {
		r  image.Rectangle
		in uiInsets
	}{
		"plank": {uiPlankRect, uiPlankInsets}, "plankArrow": {uiPlankArrowRect, uiPlankArrowInsets},
		"board": {uiBoardRect, uiBoardInsets},
	} {
		if s.in.left+s.in.right >= s.r.Dx() || s.in.top+s.in.bottom >= s.r.Dy() {
			t.Errorf("%s insets %+v leave nothing to stretch in %v", name, s.in, s.r)
		}
	}
}

// TestUISkinSlicesCover: the pieces fill dst exactly, without overlapping,
// and the corners keep their shape at the given scale.
func TestUISkinSlicesCover(t *testing.T) {
	for name, c := range map[string]struct {
		src   image.Rectangle
		in    uiInsets
		dst   image.Rectangle
		scale float64
	}{
		"plank": {uiPlankRect, uiPlankInsets, image.Rect(100, 50, 420, 120), 70.0 / 238},
		"board": {uiBoardRect, uiBoardInsets, image.Rect(200, 100, 1000, 640), 0.5},
		"tiny":  {uiPlankArrowRect, uiPlankArrowInsets, image.Rect(0, 0, 30, 40), 40.0 / 240},
	} {
		from, to := sliceRects(c.src, c.in, c.dst, c.scale)
		area := 0
		for i, r := range to {
			area += r.Dx() * r.Dy()
			if !r.In(c.dst) || !from[i].In(c.src) {
				t.Errorf("%s: piece %v→%v outside src %v / dst %v", name, from[i], r, c.src, c.dst)
			}
			for _, o := range to[:i] {
				if r.Overlaps(o) {
					t.Errorf("%s: pieces %v and %v overlap", name, r, o)
				}
			}
		}
		if want := c.dst.Dx() * c.dst.Dy(); area != want {
			t.Errorf("%s: pieces cover %d px², want %d", name, area, want)
		}
		if name == "tiny" {
			continue
		}
		if got, want := float64(to[0].Dx()), float64(from[0].Dx())*c.scale; math.Abs(got-want) > 1 {
			t.Errorf("%s: corner %d px wide, want %.1f (scale %.3f)", name, to[0].Dx(), want, c.scale)
		}
	}
}

// TestUISkinBarFill: an empty bar draws only its left end, a full one the
// whole bar, and the fill grows with the volume.
func TestUISkinBarFill(t *testing.T) {
	if r := barFillRect(0); r.Max.X != uiBarSegLeft {
		t.Errorf("fill 0 ends at x %d, want the first segment's start %d", r.Max.X, uiBarSegLeft)
	}
	if r := barFillRect(1); r != uiBarFullRect {
		t.Errorf("fill 1 is %v, want the whole bar %v", r, uiBarFullRect)
	}
	prev := 0
	for v := range MaxVolume + 1 {
		w := barFillRect(float64(v) / MaxVolume).Dx()
		if w <= prev && v > 0 {
			t.Errorf("volume %d draws %d px of the bar, not more than volume %d's %d", v, w, v-1, prev)
		}
		prev = w
	}
}
