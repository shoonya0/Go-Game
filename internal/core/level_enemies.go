package core

import (
	"encoding/json"
	"fmt"
)

// levelCell is one cell of a JSON level grid. A number is a tile code (see
// tileTypeFromCode and CheckpointCode); a string is an enemy marker, the
// one-character code an EnemyBlueprint registers as its Marker (e.g. "!").
type levelCell struct {
	tile   int
	marker string
}

// UnmarshalJSON accepts either a JSON number or a JSON string.
func (c *levelCell) UnmarshalJSON(b []byte) error {
	if len(b) > 0 && b[0] == '"' {
		return json.Unmarshal(b, &c.marker)
	}
	return json.Unmarshal(b, &c.tile)
}

// EnemySpawn is where an enemy of a given kind starts in a level.
type EnemySpawn struct {
	Kind EnemyKind
	Pos  Position
}

// kindForMarker returns the enemy kind whose blueprint registers marker.
func kindForMarker(marker string) (EnemyKind, bool) {
	for kind, bp := range enemyBlueprints {
		if bp.Marker == marker {
			return kind, true
		}
	}
	return 0, false
}

// splitGrid separates JSON cells into the tile-code grid BuildLevel consumes
// (enemy cells become empty) and the enemy spawns. An unknown marker is an
// error, so a typo in the level file fails at load instead of silently
// dropping an enemy.
func splitGrid(cells [][]levelCell) ([][]int, []EnemySpawn, error) {
	grid := make([][]int, len(cells))
	var spawns []EnemySpawn
	for y, row := range cells {
		grid[y] = make([]int, len(row))
		for x, c := range row {
			if c.marker == "" {
				grid[y][x] = c.tile
				continue
			}
			kind, ok := kindForMarker(c.marker)
			if !ok {
				return nil, nil, fmt.Errorf("unknown enemy marker %q at row %d, col %d", c.marker, y, x)
			}
			spawns = append(spawns, EnemySpawn{Kind: kind, Pos: spawnPos(blueprintFor(kind), x, y)})
		}
	}
	return grid, spawns, nil
}

// spawnPos places an enemy's collision box centred in grid cell (x, y) with its
// feet on the cell's bottom edge, so a marker sitting on a platform stands on it.
func spawnPos(bp *EnemyBlueprint, x, y int) Position {
	return Position{
		X: float64(x*LevelTileWidth) + (LevelTileWidth-bp.BodyWidth)/2,
		Y: float64((y+1)*LevelTileHeight) - bp.BodyHeight,
	}
}

// SpawnEnemies creates one enemy per spawn.
func SpawnEnemies(spawns []EnemySpawn) []*EnemyRuntime {
	enemies := make([]*EnemyRuntime, 0, len(spawns))
	for _, s := range spawns {
		enemies = append(enemies, NewEnemy(s.Kind, s.Pos))
	}
	return enemies
}
