package core

import "testing"

// jumpWorld is 12×12 tiles of flat rock ground (row 11) with, from column 5
// on, a step that is step tiles high.
func jumpWorld(step int) *DynamicQuadtree {
	const n = 12
	grid := make([][]int, n)
	for r := range grid {
		grid[r] = make([]int, n)
	}
	for c := range n {
		grid[n-1][c] = 4
		if c >= 5 {
			for r := n - 1 - step; r < n-1; r++ {
				grid[r][c] = 4
			}
		}
	}
	level := BuildLevel(grid)
	qt := NewDynamicQuadtree(AABB{Width: n * LevelTileWidth, Height: n * LevelTileHeight})
	for i := range level.Platforms {
		qt.Insert(&level.Platforms[i])
	}
	return qt
}

// groundY is the world y of jumpWorld's ground.
const groundY = 11 * LevelTileHeight

// jumpFrom stands a player at x on jumpWorld(step), lets it land, then jumps
// once, holding right when run is set, for two seconds. It returns how high
// the feet rose at most and where the player ended.
func jumpFrom(step int, x float64, right bool) (rise float64, p *PlayerRuntime) {
	WorldInit()
	qt := jumpWorld(step)
	p = InitPlayer(nil)
	p.placeAt(Position{X: x, Y: groundY - bodyHeight})
	for range 10 {
		UpdatePlayer(p, &InputState{}, qt, nil)
	}
	startY, minY := p.Pos.Y, p.Pos.Y
	for i := range int(2 * tps()) {
		in := InputState{JumpJustPressed: i == 0}
		if right {
			in.Direction.LeftRight = 1
		}
		UpdatePlayer(p, &in, qt, nil)
		minY = min(minY, p.Pos.Y)
	}
	return (startY - minY) / LevelTileHeight, p
}

// TestJumpReachesThreeTiles: a jump lifts the feet just over 3 tiles, so the
// player climbs a 3-tile step but not a 4-tile one (the user's rule for every
// level grid).
func TestJumpReachesThreeTiles(t *testing.T) {
	rise, _ := jumpFrom(0, LevelTileWidth, false)
	t.Logf("a jump lifts the feet %.2f tiles", rise)
	if rise < 3 || rise >= 3.25 {
		t.Errorf("a jump lifts the feet %.2f tiles, want 3 to 3.25", rise)
	}
	for step, want := range map[int]bool{3: true, 4: false} {
		_, p := jumpFrom(step, 5*LevelTileWidth-bodyWidth-2, true)
		top := float64(groundY - step*LevelTileHeight)
		if on := p.Pos.Y+bodyHeight <= top+1; on != want {
			t.Errorf("%d-tile step: climbed=%v (feet at %.0f, step top %.0f), want %v", step, on, p.Pos.Y+bodyHeight, top, want)
		}
	}
}
