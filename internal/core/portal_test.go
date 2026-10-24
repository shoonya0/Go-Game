package core

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"testing"
)

// TestArtFramesFindEachFramesArt: on a synthetic 2×1 sheet, each frame is the
// tight box around its own opaque pixels, wherever they sit in the cell.
// Both the fast NRGBA path and the generic one are checked.
func TestArtFramesFindEachFramesArt(t *testing.T) {
	for _, img := range []draw.Image{image.NewNRGBA(image.Rect(0, 0, 20, 10)), image.NewRGBA(image.Rect(0, 0, 20, 10))} {
		img.Set(2, 3, color.NRGBA{A: 255})
		img.Set(5, 8, color.NRGBA{A: 255})
		img.Set(14, 1, color.NRGBA{A: 255})
		img.Set(17, 1, color.NRGBA{A: 10}) // too faint to count

		got := artFrames(img, 2, 1, portalAlphaMin)
		want := []image.Rectangle{image.Rect(2, 3, 6, 9), image.Rect(14, 1, 15, 2)}
		if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
			t.Errorf("%T: artFrames = %v, want %v", img, got, want)
		}
	}
}

// TestPortalSheetFrames: every one of the portal's 14 frames has art, inside
// its own cell of the 7×2 grid.
func TestPortalSheetFrames(t *testing.T) {
	img, _, err := image.Decode(bytes.NewReader(readAsset(portalSheet)))
	if err != nil {
		t.Fatal(err)
	}
	frames := artFrames(img, portalCols, portalRows, portalAlphaMin)
	if len(frames) != portalCols*portalRows {
		t.Fatalf("%d frames, want %d", len(frames), portalCols*portalRows)
	}
	cw, ch := img.Bounds().Dx()/portalCols, img.Bounds().Dy()/portalRows
	for i, r := range frames {
		cell := image.Rect(i%portalCols*cw, i/portalCols*ch, (i%portalCols+1)*cw, (i/portalCols+1)*ch)
		if r.Empty() || !r.In(cell) {
			t.Errorf("frame %d: art %v, want a non-empty box inside its cell %v", i, r, cell)
		}
	}
}

// TestPortalStandsOnItsCell: a portal's area is centered on its grid cell, with
// its base on the cell's floor, portalHeight tall.
func TestPortalStandsOnItsCell(t *testing.T) {
	p := Portal{Pos: Position{X: 3 * LevelTileWidth, Y: 5 * LevelTileHeight}}
	b := p.Bounds()
	if centerX(b) != 3.5*LevelTileWidth || b.Y+b.Height != 6*LevelTileHeight || b.Height != portalHeight {
		t.Errorf("portal bounds %+v, want centered on x=%v, base at y=%v, %v tall",
			b, 3.5*LevelTileWidth, 6*LevelTileHeight, portalHeight)
	}
}
