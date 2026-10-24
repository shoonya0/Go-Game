package main

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

// The cells tile the sheet with no gaps or overlaps, even when its height
// doesn't divide by 6 (the gnoll sheet: 1298/6 = 216.33 px).
func TestCellRectsTileNonIntegerSheet(t *testing.T) {
	bounds := image.Rect(0, 0, 3375, 1298)
	for row := range gridRows {
		for col := range gridCols {
			r := cellRect(bounds, col, row)
			if h := r.Dy(); h != 216 && h != 217 {
				t.Errorf("cell (%d,%d) is %d px tall, want 216 or 217", col, row, h)
			}
			if w := r.Dx(); w != 375 {
				t.Errorf("cell (%d,%d) is %d px wide, want 375", col, row, w)
			}
			if col > 0 && r.Min.X != cellRect(bounds, col-1, row).Max.X {
				t.Errorf("cell (%d,%d) doesn't start where its left neighbour ends", col, row)
			}
			if row > 0 && r.Min.Y != cellRect(bounds, col, row-1).Max.Y {
				t.Errorf("cell (%d,%d) doesn't start where the one above ends", col, row)
			}
		}
	}
	if last := cellRect(bounds, gridCols-1, gridRows-1); last.Max != bounds.Max {
		t.Errorf("last cell ends at %v, want the sheet corner %v", last.Max, bounds.Max)
	}
}

// Shrinking averages the source pixels each output pixel covers, including
// the pixels it only partly covers.
func TestShrinkAverages(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 3, 1))
	for x, v := range []uint8{0, 90, 180} {
		src.SetRGBA(x, 0, color.RGBA{v, v, v, 255})
	}
	got := shrink(src, 2, 1)
	// Output pixel 0 covers source 0 and half of 1: (0 + 45) / 1.5 = 30.
	// Output pixel 1 covers half of 1 and all of 2: (45 + 180) / 1.5 = 150.
	for x, want := range []uint8{30, 150} {
		if c := got.RGBAAt(x, 0); c.R != want || c.A != 255 {
			t.Errorf("pixel %d = %v, want grey %d, opaque", x, c, want)
		}
	}
}

// testSheet is a 9×6 battler sheet whose cells are each one flat colour (red =
// column, green = row), with a white marker pixel in each cell's top-left corner.
func testSheet(w, h int) *image.RGBA {
	sheet := image.NewRGBA(image.Rect(0, 0, w, h))
	for row := range gridRows {
		for col := range gridCols {
			r := cellRect(sheet.Bounds(), col, row)
			for y := r.Min.Y; y < r.Max.Y; y++ {
				for x := r.Min.X; x < r.Max.X; x++ {
					sheet.SetRGBA(x, y, color.RGBA{uint8(col * 20), uint8(row * 40), 0, 255})
				}
			}
			sheet.SetRGBA(r.Min.X, r.Min.Y, color.RGBA{255, 255, 255, 255})
		}
	}
	return sheet
}

func TestBakeClipChainsCellsInOrder(t *testing.T) {
	m := Mapping{Name: "Test", DrawHeight: 10}
	clip := ClipMapping{Frames: [][3]int{{1, 0, 2}, {0, 4, 1}}} // stabbing 3, damage 2
	strip, err := bakeClip(testSheet(90, 60), m, clip)
	if err != nil {
		t.Fatal(err)
	}
	if b := strip.Bounds(); b.Dx() != 20 || b.Dy() != 10 {
		t.Fatalf("strip is %dx%d, want 20x10 (2 frames of 10x10)", b.Dx(), b.Dy())
	}
	// Sample the middle of each frame, away from the marker pixel.
	for i, want := range []color.RGBA{{100, 0, 0, 255}, {20, 160, 0, 255}} {
		if c := strip.RGBAAt(i*10+5, 5); c != want {
			t.Errorf("frame %d = %v, want %v", i, c, want)
		}
	}
}

func TestBakeClipMirrorsLeftFacingArt(t *testing.T) {
	clip := ClipMapping{Frames: [][3]int{{0, 0, 0}}}
	for _, tc := range []struct {
		facesLeft bool
		markerX   int
	}{{false, 0}, {true, 9}} {
		m := Mapping{Name: "Test", DrawHeight: 10, FacesLeft: tc.facesLeft}
		strip, err := bakeClip(testSheet(90, 60), m, clip)
		if err != nil {
			t.Fatal(err)
		}
		if c := strip.RGBAAt(tc.markerX, 0); c.R != 255 || c.G != 255 {
			t.Errorf("facesLeft=%v: marker not at x=%d (got %v)", tc.facesLeft, tc.markerX, c)
		}
	}
}

func TestMappingValidate(t *testing.T) {
	ok := Mapping{Name: "Goblin", Sheet: "goblin.png", DrawHeight: 100,
		Clips: map[string]ClipMapping{"idle": {Frames: [][3]int{{0, 1, 0}}, FPS: 6}}}
	if err := ok.validate(); err != nil {
		t.Fatalf("valid mapping rejected: %v", err)
	}
	for name, mutate := range map[string]func(*Mapping){
		"unknown clip": func(m *Mapping) { m.Clips = map[string]ClipMapping{"jump": m.Clips["idle"]} },
		"block past 2": func(m *Mapping) { m.Clips = map[string]ClipMapping{"idle": {Frames: [][3]int{{3, 0, 0}}, FPS: 6}} },
		"row past 5":   func(m *Mapping) { m.Clips = map[string]ClipMapping{"idle": {Frames: [][3]int{{0, 6, 0}}, FPS: 6}} },
		"frame past 2": func(m *Mapping) { m.Clips = map[string]ClipMapping{"idle": {Frames: [][3]int{{0, 0, 3}}, FPS: 6}} },
		"no frames":    func(m *Mapping) { m.Clips = map[string]ClipMapping{"idle": {FPS: 6}} },
		"no fps":       func(m *Mapping) { m.Clips = map[string]ClipMapping{"idle": {Frames: [][3]int{{0, 0, 0}}}} },
		"hit past the end": func(m *Mapping) {
			m.Clips = map[string]ClipMapping{"idle": {Frames: [][3]int{{0, 0, 0}}, FPS: 6, HitFrame: 1}}
		},
		"no draw height": func(m *Mapping) { m.DrawHeight = 0 },
		"no name":        func(m *Mapping) { m.Name = "" },
	} {
		m := ok
		mutate(&m)
		if err := m.validate(); err == nil {
			t.Errorf("%s: accepted, want an error", name)
		}
	}
}

func TestBakeRejectsUpscaling(t *testing.T) {
	m := Mapping{Name: "Test", DrawHeight: 20} // cells are 10 px tall
	if _, err := bakeClip(testSheet(90, 60), m, ClipMapping{Frames: [][3]int{{0, 0, 0}}}); err == nil {
		t.Error("drawHeight above the cell height accepted, want an error")
	}
}

func TestSnippetNamesClipsAndSize(t *testing.T) {
	m := Mapping{Name: "Goblin", Clips: map[string]ClipMapping{
		"idle":     {Frames: make([][3]int, 4), FPS: 6, Loop: true},
		"attack_1": {Frames: make([][3]int, 6), FPS: 12, Hold: []float64{1, 1, 2}, HitFrame: 3},
	}}
	got := snippet(m, "assets/enemy/Goblin", 120, 90)
	for _, want := range []string{
		"goblinFrameWidth  = 120",
		"goblinFrameHeight = 90",
		`ClipIdle: {Sheet: "../assets/enemy/Goblin/Goblin_idle.png", Frames: 4, FrameW: goblinFrameWidth, FrameH: goblinFrameHeight, FPS: 6, Loop: true}`,
		"ClipAttack1: {Sheet: \"../assets/enemy/Goblin/Goblin_attack_1.png\", Frames: 6, FrameW: goblinFrameWidth, FrameH: goblinFrameHeight, FPS: 12, Hold: []float64{1, 1, 2}, HitFrame: 3}",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("snippet is missing %q\n%s", want, got)
		}
	}
	// Clips come out in the game's clip order, so re-runs diff cleanly.
	if strings.Index(got, "ClipIdle") > strings.Index(got, "ClipAttack1") {
		t.Errorf("clips not in standard order:\n%s", got)
	}
}

// The contact sheet stacks the strips top to bottom on an opaque background,
// with a separator line between frames, so every clip can be checked by eye
// in one image.
func TestContactSheetStacksStrips(t *testing.T) {
	a := image.NewRGBA(image.Rect(0, 0, 30, 10)) // 3 frames of 10×10
	b := image.NewRGBA(image.Rect(0, 0, 10, 10)) // 1 frame
	for x := range 10 {
		for y := range 10 {
			b.SetRGBA(x, y, color.RGBA{0, 0, 255, 255})
		}
	}
	sheet := contactSheet([]*image.RGBA{a, b}, 10)
	if got := sheet.Bounds(); got.Dx() != 30 || got.Dy() != 2*(10+contactGap) {
		t.Fatalf("contact sheet is %dx%d, want 30x%d", got.Dx(), got.Dy(), 2*(10+contactGap))
	}
	if c := sheet.RGBAAt(5, 5); c != contactBackground {
		t.Errorf("transparent frame shows %v, want the background %v", c, contactBackground)
	}
	if c := sheet.RGBAAt(10, 5); c != contactLine {
		t.Errorf("frame edge shows %v, want a separator line", c)
	}
	if c := sheet.RGBAAt(5, 10+contactGap+5); c.B != 255 {
		t.Errorf("second strip's frame shows %v, want its blue", c)
	}
}

// TestFloorStripsMovesFeetToBottom checks every strip moves down by the same
// amount, so the lowest opaque pixel of all clips lands on the bottom row and
// the clips stay lined up with each other.
func TestFloorStripsMovesFeetToBottom(t *testing.T) {
	opaque := color.RGBA{0, 0, 0, 255}
	a := image.NewRGBA(image.Rect(0, 0, 2, 6))
	b := image.NewRGBA(image.Rect(0, 0, 2, 6))
	a.SetRGBA(0, 2, opaque)
	b.SetRGBA(1, 4, opaque)
	b.SetRGBA(0, 4, color.RGBA{0, 0, 0, footAlpha - 1}) // too faint to count as a foot

	if got := floorStrips([]*image.RGBA{a, b}); got != 1 {
		t.Errorf("shift = %d, want 1 (lowest foot at y=4 of 6)", got)
	}
	if a.RGBAAt(0, 3) != opaque || a.RGBAAt(0, 2).A != 0 || b.RGBAAt(1, 5) != opaque || b.RGBAAt(1, 0).A != 0 {
		t.Errorf("after floorStrips: a(0,3)=%v a(0,2)=%v b(1,5)=%v b(1,0)=%v, want both moved down 1 px, top row clear",
			a.RGBAAt(0, 3), a.RGBAAt(0, 2), b.RGBAAt(1, 5), b.RGBAAt(1, 0))
	}

	empty := image.NewRGBA(image.Rect(0, 0, 2, 6))
	if got := floorStrips([]*image.RGBA{empty}); got != 0 {
		t.Errorf("shift of an empty strip = %d, want 0", got)
	}
}
