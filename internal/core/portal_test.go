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
	if p.PortalLook(level.Portals[0]) != PortalCheckpoint || p.PortalLook(level.Portals[1]) == PortalCheckpoint {
		t.Fatalf("after start: looks %v, %v; want only the start portal golden",
			p.PortalLook(level.Portals[0]), p.PortalLook(level.Portals[1]))
	}

	p.Pos = level.Portals[1].StandPos()
	p.updatePortals(level.Portals, &InputState{Interact: true})
	if p.PortalLook(level.Portals[0]) == PortalCheckpoint || p.PortalLook(level.Portals[1]) != PortalCheckpoint {
		t.Errorf("after saving the second: looks %v, %v; want only the second golden",
			p.PortalLook(level.Portals[0]), p.PortalLook(level.Portals[1]))
	}
}

// twoPortals is a level with a portal at each end of a rock floor, and a
// player who started at the first one.
func twoPortals(t *testing.T) (Level, *PlayerRuntime) {
	t.Helper()
	WorldInit()
	level := BuildLevel([][]int{
		{9, 0, 0, 0, 0, 9},
		{4, 4, 4, 4, 4, 4},
	})
	p := InitPlayer(nil)
	p.StartAt(level)
	return level, p
}

// TestPortalSavesOnlyOnInteract: standing in a portal saves nothing until E is
// pressed; E away from every portal saves nothing either.
func TestPortalSavesOnlyOnInteract(t *testing.T) {
	level, p := twoPortals(t)
	start := p.RespawnPos

	p.Pos = Position{X: 2.5 * LevelTileWidth, Y: level.Portals[1].StandPos().Y} // between the portals
	p.updatePortals(level.Portals, &InputState{Interact: true})
	if p.RespawnPos != start {
		t.Fatalf("E away from portals moved the respawn to %v", p.RespawnPos)
	}

	p.Pos = level.Portals[1].StandPos()
	p.updatePortals(level.Portals, &InputState{})
	if p.RespawnPos != start {
		t.Fatalf("standing in a portal saved it without E (respawn %v)", p.RespawnPos)
	}

	p.updatePortals(level.Portals, &InputState{Interact: true})
	if p.RespawnPos != level.Portals[1].StandPos() || p.PortalLook(level.Portals[1]) != PortalCheckpoint {
		t.Errorf("E in the second portal: respawn %v, want its StandPos %v", p.RespawnPos, level.Portals[1].StandPos())
	}
	if p.savedNote <= 0 {
		t.Error("saving should show the saved note")
	}
}

// TestTapOnPortalSaves: on touch screens, a tap on the portal the player
// stands in saves it; a tap anywhere else does not. Taps are screen points.
func TestTapOnPortalSaves(t *testing.T) {
	level, p := twoPortals(t)
	start := p.RespawnPos
	p.Pos = level.Portals[1].StandPos()
	p.Camera.Pos = Position{X: 100, Y: -40}
	b := level.Portals[1].Bounds()
	onPortal := Position{X: centerX(b) - p.Camera.Pos.X, Y: b.Y + b.Height/2 - p.Camera.Pos.Y}

	p.updatePortals(level.Portals, &InputState{Tapped: true, TapAt: Position{X: onPortal.X + 4*LevelTileWidth, Y: onPortal.Y}})
	if p.RespawnPos != start {
		t.Fatalf("a tap beside the portal saved it (respawn %v)", p.RespawnPos)
	}
	p.Pos = level.Portals[0].StandPos() // standing in the other portal
	p.updatePortals(level.Portals, &InputState{Tapped: true, TapAt: onPortal})
	if p.RespawnPos != start {
		t.Fatalf("a tap on a portal the player isn't in saved it (respawn %v)", p.RespawnPos)
	}

	p.Pos = level.Portals[1].StandPos()
	p.updatePortals(level.Portals, &InputState{Tapped: true, TapAt: onPortal})
	if p.RespawnPos != level.Portals[1].StandPos() || p.PortalLook(level.Portals[1]) != PortalCheckpoint {
		t.Errorf("tap on the portal: respawn %v, want its StandPos %v", p.RespawnPos, level.Portals[1].StandPos())
	}
}

// TestRespawnAtSavedPortal: after dying, the player is back standing in the
// last saved portal, not the start one, with full health and power.
func TestRespawnAtSavedPortal(t *testing.T) {
	withDeath(t, false)
	level, p := twoPortals(t)
	p.Pos = level.Portals[1].StandPos()
	p.updatePortals(level.Portals, &InputState{Interact: true})

	p.Pos = Position{X: 2.5 * LevelTileWidth, Y: 0} // wander off, then die there
	p.Combat.Power = 0
	p.ApplyDamage(p.Combat.Health)
	p.updateCombat(nil, true)
	for i := 0; i < 5*int(tps()) && p.State.IsDying(); i++ {
		p.UpdateAnimation()
		p.PreviousState = p.State
	}

	if p.State.IsDying() || p.Respawns != 1 {
		t.Fatalf("dying=%v respawns=%d after 5 s, want respawned once", p.State.IsDying(), p.Respawns)
	}
	if want := level.Portals[1].StandPos(); p.Pos != want {
		t.Errorf("respawned at %v, want standing in the saved portal at %v", p.Pos, want)
	}
	if p.Combat.Health != p.Combat.MaxHealth || p.Combat.Power != p.Combat.MaxPower {
		t.Errorf("health %d/%d power %v/%v, want both full", p.Combat.Health, p.Combat.MaxHealth, p.Combat.Power, p.Combat.MaxPower)
	}
}

// TestPortalActivatesOnInteract: E in a portal activates it and makes it the
// checkpoint; the old checkpoint stays active. Active portals are listed once
// each, in the order they were activated, by level and cell.
func TestPortalActivatesOnInteract(t *testing.T) {
	level, p := twoPortals(t)
	start := level.Portals[0].Ref()
	if len(p.ActivePortals) != 1 || p.ActivePortals[0] != start || p.Checkpoint != start {
		t.Fatalf("after start: active %v, checkpoint %v; want only the start portal %v", p.ActivePortals, p.Checkpoint, start)
	}

	p.Pos = level.Portals[1].StandPos()
	p.updatePortals(level.Portals, &InputState{Interact: true})
	p.updatePortals(level.Portals, &InputState{Interact: true}) // a second press adds nothing
	second := level.Portals[1].Ref()
	if len(p.ActivePortals) != 2 || p.ActivePortals[0] != start || p.ActivePortals[1] != second {
		t.Errorf("active %v, want [%v %v]", p.ActivePortals, start, second)
	}
	if p.Checkpoint != second {
		t.Errorf("checkpoint %v, want the second portal %v", p.Checkpoint, second)
	}
}

// TestPortalLooks: a portal never used is inactive, one used is active, and
// the checkpoint is golden. The same cell in another level is another portal.
func TestPortalLooks(t *testing.T) {
	WorldInit()
	level := BuildLevel([][]int{
		{9, 0, 9, 0, 9},
		{4, 4, 4, 4, 4},
	})
	p := InitPlayer(nil)
	p.StartAt(level)
	p.Pos = level.Portals[1].StandPos()
	p.updatePortals(level.Portals, &InputState{Interact: true})

	want := []PortalLook{PortalActive, PortalCheckpoint, PortalInactive}
	for i, c := range level.Portals {
		if got := p.PortalLook(c); got != want[i] {
			t.Errorf("portal %d: look %v, want %v", i, got, want[i])
		}
	}
	elsewhere := level.Portals[0]
	elsewhere.Level = "otherLevel"
	if got := p.PortalLook(elsewhere); got != PortalInactive {
		t.Errorf("same cell in another level: look %v, want inactive", got)
	}
}
