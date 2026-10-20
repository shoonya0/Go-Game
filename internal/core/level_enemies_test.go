package core

import (
	"encoding/json"
	"testing"
)

func TestSplitGridEnemyMarkers(t *testing.T) {
	var cells [][]levelCell
	if err := json.Unmarshal([]byte(`[[0,"!",9],[2,2,2]]`), &cells); err != nil {
		t.Fatal(err)
	}
	grid, spawns, err := splitGrid(cells)
	if err != nil {
		t.Fatal(err)
	}

	if grid[0][1] != 0 || grid[0][2] != CheckpointCode || grid[1][0] != 2 {
		t.Errorf("tile grid = %v, want marker cell empty and other codes kept", grid)
	}
	if len(spawns) != 1 || spawns[0].Kind != KindSuccubus {
		t.Fatalf("spawns = %+v, want one succubus", spawns)
	}

	bp := blueprintFor(KindSuccubus)
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
	in := `{"lvl!": {"grid": [[0,!,@],[9,"#",-1]]}}`
	want := `{"lvl!": {"grid": [[0,"!","@"],[9,"#",-1]]}}`
	if got := string(quoteBareMarkers([]byte(in))); got != want {
		t.Fatalf("quoteBareMarkers:\n got %s\nwant %s", got, want)
	}
}

func TestEnemyMarkersUniqueSingleChar(t *testing.T) {
	seen := map[string]EnemyKind{}
	for kind, bp := range enemyBlueprints {
		if len(bp.Marker) != 1 || !isMarkerChar(bp.Marker[0]) {
			t.Errorf("%s: marker %q must be one ASCII punctuation character usable bare in the level", bp.Name, bp.Marker)
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
