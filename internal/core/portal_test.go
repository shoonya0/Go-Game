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

// TestPlayerStartsAtFirstPortal: a new game puts the player standing in the
// level's first portal in reading order (top row first), and saves it.
func TestPlayerStartsAtFirstPortal(t *testing.T) {
	WorldInit()
	grid := [][]int{
		{0, 0, 0, 9},
		{9, 0, 0, 1},
		{1, 1, 1, 1},
	}
	level := BuildLevel(grid)
	p := InitPlayer(nil)
	p.StartAt(level)

	first := level.Portals[0]
	if first.Pos != (Position{X: 3 * LevelTileWidth, Y: 0}) {
		t.Fatalf("first portal at %v, want the top-row one", first.Pos)
	}
	body := p.GetBounds()
	if centerX(body) != centerX(first.Bounds()) || body.Y+body.Height != first.Bounds().Y+first.Bounds().Height {
		t.Errorf("player body %+v, want centered in portal %+v, feet on its floor", body, first.Bounds())
	}
	if p.RespawnPos != p.Pos {
		t.Errorf("respawn %v, want the start position %v (start portal saved)", p.RespawnPos, p.Pos)
	}
}

// TestPortalStandsOnTheGround: a portal's base, and the player standing in it,
// sit on the solid top of the tile below, which for grass, sand and water is
// nudged below the cell line (TopTileVisualOffset).
func TestPortalStandsOnTheGround(t *testing.T) {
	WorldInit()
	for _, below := range []TileType{Grass, Sand, Water, Rock} {
		code := 0
		for i, tt := range tileRowOrder {
			if tt == below {
				code = i + 1
			}
		}
		level := BuildLevel([][]int{{9}, {code}})
		ground := level.Platforms[0].Y
		portal := level.Portals[0]
		if b := portal.Bounds(); b.Y+b.Height != ground {
			t.Errorf("%s: portal base at y=%v, want the ground's top %v", below, b.Y+b.Height, ground)
		}
		if feet := portal.StandPos().Y + bodyHeight; feet != ground {
			t.Errorf("%s: player feet at y=%v, want the ground's top %v", below, feet, ground)
		}
	}
}

// TestOnlySavedPortalIsGolden: the saved portal (drawn golden) is the start
// portal at first, then whichever portal the player saved last.
func TestOnlySavedPortalIsGolden(t *testing.T) {
	WorldInit()
	level := BuildLevel([][]int{
		{9, 0, 0, 9},
		{4, 4, 4, 4},
	})
	p := InitPlayer(nil)
	p.StartAt(level)
	if !level.Portals[0].SavedBy(p.RespawnPos) || level.Portals[1].SavedBy(p.RespawnPos) {
		t.Fatalf("after start: saved = %v, %v; want only the start portal",
			level.Portals[0].SavedBy(p.RespawnPos), level.Portals[1].SavedBy(p.RespawnPos))
	}

	p.Pos = level.Portals[1].StandPos()
	p.updatePortals(level.Portals)
	if level.Portals[0].SavedBy(p.RespawnPos) || !level.Portals[1].SavedBy(p.RespawnPos) {
		t.Errorf("after saving the second: saved = %v, %v; want only the second",
			level.Portals[0].SavedBy(p.RespawnPos), level.Portals[1].SavedBy(p.RespawnPos))
	}
}
