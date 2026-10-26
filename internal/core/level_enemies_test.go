package core

import (
	"encoding/json"
	"testing"
)

func TestSplitGridEnemyMarkers(t *testing.T) {
	var cells [][]levelCell
	if err := json.Unmarshal([]byte(`[[0,"@",9],[2,2,2]]`), &cells); err != nil {
		t.Fatal(err)
	}
	grid, spawns, err := splitGrid(cells)
	if err != nil {
		t.Fatal(err)
	}

	if grid[0][1] != 0 || grid[0][2] != PortalCode || grid[1][0] != 2 {
		t.Errorf("tile grid = %v, want marker cell empty and other codes kept", grid)
	}
	if len(spawns) != 1 || spawns[0].Kind != KindSkeletonWhite {
		t.Fatalf("spawns = %+v, want one white skeleton", spawns)
	}

	bp := blueprintFor(KindSkeletonWhite)
	pos := spawns[0].Pos
	if feet := pos.Y + bp.BodyHeight; feet != LevelTileHeight {
		t.Errorf("enemy feet at y=%v, want on top of the row below (y=%v)", feet, float64(LevelTileHeight))
	}
	if mid := pos.X + bp.BodyWidth/2; mid != 1.5*LevelTileWidth {
		t.Errorf("enemy centre at x=%v, want centred in column 1 (x=%v)", mid, 1.5*LevelTileWidth)
	}
}

func TestSplitGridUnknownMarker(t *testing.T) {
	var cells [][]levelCell
	if err := json.Unmarshal([]byte(`[[0,"?"]]`), &cells); err != nil {
		t.Fatal(err)
	}
	if _, _, err := splitGrid(cells); err == nil {
		t.Fatal("splitGrid accepted an unknown marker")
	}
}

func TestQuoteBareMarkers(t *testing.T) {
	in := `{"lvl!": {"grid": [[0,!,@],[9,"#",-1],[!@, $$ ,2]]}}`
	want := `{"lvl!": {"grid": [[0,"!","@"],[9,"#",-1],["!@", "$$" ,2]]}}`
	if got := string(quoteBareMarkers([]byte(in))); got != want {
		t.Fatalf("quoteBareMarkers:\n got %s\nwant %s", got, want)
	}
}

// TestTwoCharMarkerIsOneCell checks a bare two-character marker reaches the
// grid as one marker cell, not two cells or a parse error.
func TestTwoCharMarkerIsOneCell(t *testing.T) {
	var cells [][]levelCell
	if err := json.Unmarshal(quoteBareMarkers([]byte(`[[0,!@,0]]`)), &cells); err != nil {
		t.Fatal(err)
	}
	if len(cells[0]) != 3 || cells[0][1].marker != "!@" {
		t.Fatalf("cells = %+v, want 3 cells with marker \"!@\" in the middle", cells)
	}
}

func TestEnemyMarkersUniqueAndBare(t *testing.T) {
	seen := map[string]EnemyKind{}
	for kind, bp := range enemyBlueprints {
		bare := len(bp.Marker) >= 1 && len(bp.Marker) <= maxMarkerLen
		for i := 0; i < len(bp.Marker); i++ {
			bare = bare && isMarkerChar(bp.Marker[i])
		}
		if !bare {
			t.Errorf("%s: marker %q must be 1-%d ASCII punctuation characters usable bare in the level", bp.Name, bp.Marker, maxMarkerLen)
		}
		if other, dup := seen[bp.Marker]; dup {
			t.Errorf("marker %q used by kinds %d and %d", bp.Marker, other, kind)
		}
		seen[bp.Marker] = kind
	}
}

// Enemies now come only from the level file, so the shipped level must parse
// and place at least one.
func TestLevelFileSpawnsEnemies(t *testing.T) {
	WorldInit()
	level := LoadLevel(LevelFile, DefaultLevelName)
	if len(level.Enemies) == 0 {
		t.Fatal("level_1.json has no enemy markers")
	}
}

// TestLevelEnemiesKeepTheirGround spawns each level as the game does and lets
// every enemy patrol for 20 seconds with the player far away: each enemy must
// still stand on the ground it spawned on, not fall off a platform or through
// the floor. A flyer has no ground: it must still fly at the height it spawned
// at. Level 1 must also hold every registered kind.
func TestLevelEnemiesKeepTheirGround(t *testing.T) {
	WorldInit()
	for name, def := range Levels {
		t.Run(name, func(t *testing.T) { levelEnemiesKeepTheirGround(t, name, def) })
	}
}

func levelEnemiesKeepTheirGround(t *testing.T, name string, def LevelDef) {
	level := LoadLevel(def.File, def.Key)
	qt := NewDynamicQuadtree(AABB{Width: level.Width, Height: level.Height})
	for i := range level.Platforms {
		qt.Insert(&level.Platforms[i])
	}

	if name == StartLevel {
		placed := map[EnemyKind]bool{}
		for _, s := range level.Enemies {
			placed[s.Kind] = true
		}
		for kind, bp := range enemyBlueprints {
			if !placed[kind] {
				t.Errorf("%s (marker %q) is not placed in level_1.json", bp.Name, bp.Marker)
			}
		}
	}

	for i, e := range SpawnEnemies(level.Enemies) {
		start := level.Enemies[i].Pos
		p := farPlayer(e)
		for range int(20 * tps()) { // the first ticks drop it onto its ground
			UpdateEnemy(e, p, qt)
		}
		if _, flies := e.bp.Mover.(flyingMover); flies {
			if e.Pos.Y != start.Y || e.onGround {
				t.Errorf("%s from (%v, %v): after 20 s at (%v, %v) onGround=%v, want still flying at its spawn height",
					e.bp.Name, start.X, start.Y, e.Pos.X, e.Pos.Y, e.onGround)
			}
			continue
		}
		if !e.onGround || e.Pos.Y > start.Y+LevelTileHeight {
			t.Errorf("%s from (%v, %v): after 20 s at (%v, %v) onGround=%v, want still on the ground it spawned on",
				e.bp.Name, start.X, start.Y, e.Pos.X, e.Pos.Y, e.onGround)
		}
	}
}
