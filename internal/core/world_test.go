package core

import "testing"

func TestTileTypeFromCode(t *testing.T) {
	cases := map[int]TileType{
		0:  Empty, // empty
		1:  Water,
		2:  Grass,
		3:  Sand,
		4:  Rock,
		5:  Metal,
		6:  Ice,
		7:  Wood,
		8:  Larva,
		9:  Empty, // out of range
		-1: Empty, // negative
	}
	for code, want := range cases {
		if got := tileTypeFromCode(code); got != want {
			t.Errorf("tileTypeFromCode(%d) = %q, want %q", code, got, want)
		}
	}
}

func TestBuildLevelPlatformCount(t *testing.T) {
	WorldInit()

	// 3 non-empty cells, 1 empty.
	grid := [][]int{
		{0, 2}, // empty, grass
		{1, 4}, // water, rock
	}
	level := BuildLevel(grid)
	if len(level.Platforms) != 3 {
		t.Fatalf("BuildLevel produced %d platforms, want 3", len(level.Platforms))
	}
}

func TestBuildLevelCheckpoints(t *testing.T) {
	WorldInit()

	// Code 9 is a non-solid checkpoint: it must not become a platform, and it
	// must be recorded at the right world coordinates.
	grid := [][]int{
		{0, 9, 0},
		{1, 1, 1},
	}
	level := BuildLevel(grid)
	if len(level.Platforms) != 3 {
		t.Fatalf("got %d platforms, want 3 (checkpoint must not be solid)", len(level.Platforms))
	}
	if len(level.Checkpoints) != 1 {
		t.Fatalf("got %d checkpoints, want 1", len(level.Checkpoints))
	}
	cp := level.Checkpoints[0]
	if cp.Pos.X != float64(1*LevelTileWidth) || cp.Pos.Y != 0 {
		t.Errorf("checkpoint at (%v, %v), want (%v, 0)", cp.Pos.X, cp.Pos.Y, float64(1*LevelTileWidth))
	}
}

func TestBuildLevelTopBottomDetection(t *testing.T) {
	WorldInit()

	// A single column: grass on top of grass. The lower cell has a solid
	// neighbour above, so it must be a bottom tile; the upper cell is a top tile.
	grid := [][]int{
		{2}, // top grass
		{2}, // bottom grass (solid above)
	}
	level := BuildLevel(grid)
	if len(level.Platforms) != 2 {
		t.Fatalf("got %d platforms, want 2", len(level.Platforms))
	}
	if level.Platforms[0].TileInfo.TileLvl != TopTile {
		t.Errorf("top cell = %q, want %q", level.Platforms[0].TileInfo.TileLvl, TopTile)
	}
	if level.Platforms[1].TileInfo.TileLvl != BottomTile {
		t.Errorf("bottom cell = %q, want %q", level.Platforms[1].TileInfo.TileLvl, BottomTile)
	}
}

func TestBuildLevelWorldCoordinates(t *testing.T) {
	WorldInit()

	// A rock (no visual offset) at grid (x=2, y=1) maps to world (2*W, 1*H).
	grid := [][]int{
		{0, 0, 0},
		{0, 0, 4},
	}
	level := BuildLevel(grid)
	if len(level.Platforms) != 1 {
		t.Fatalf("got %d platforms, want 1", len(level.Platforms))
	}
	if level.Platforms[0].X != float64(2*LevelTileWidth) || level.Platforms[0].Y != float64(1*LevelTileHeight) {
		t.Errorf("platform at (%v, %v), want (%v, %v)",
			level.Platforms[0].X, level.Platforms[0].Y, float64(2*LevelTileWidth), float64(1*LevelTileHeight))
	}
}
