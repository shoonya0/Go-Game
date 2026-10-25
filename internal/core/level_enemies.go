package core

import (
	"encoding/json"
	"fmt"
	"strings"
)

// quoteBareMarkers lets level files write enemy markers unquoted ([0,!,0] or
// [0,!@,0]) by wrapping each run of bare marker characters outside a JSON
// string in quotes, turning the file into valid JSON for encoding/json.
func quoteBareMarkers(data []byte) []byte {
	out := make([]byte, 0, len(data)+64)
	inString := false
	for i := 0; i < len(data); i++ {
		c := data[i]
		switch {
		case inString:
			out = append(out, c)
			if c == '\\' && i+1 < len(data) {
				i++
				out = append(out, data[i])
			} else if c == '"' {
				inString = false
			}
		case c == '"':
			inString = true
			out = append(out, c)
		case isMarkerChar(c):
			end := i + 1
			for end < len(data) && isMarkerChar(data[end]) {
				end++
			}
			out = append(out, '"')
			out = append(out, data[i:end]...)
			out = append(out, '"')
			i = end - 1
		default:
			out = append(out, c)
		}
	}
	return out
}

// maxMarkerLen is the longest enemy marker a blueprint may register: one
// punctuation character each is too few for the roster (D6).
const maxMarkerLen = 2

// isMarkerChar reports whether c may be used as a bare enemy marker: printable
// ASCII punctuation that has no meaning in JSON. Letters are excluded because
// they collide with true/false/null and number exponents.
func isMarkerChar(c byte) bool {
	if c <= ' ' || c >= 0x7f {
		return false
	}
	if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
		return false
	}
	return !strings.ContainsRune(`,[]{}:"-+.\`, rune(c))
}

// levelCell is one cell of a JSON level grid. A number is a tile code (see
// tileTypeFromCode and PortalCode); a string is an enemy marker, the
// code an EnemyBlueprint registers as its Marker (e.g. ! or !@).
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

// RespawnEnemies returns a level's enemies after the player respawns: a fresh
// enemy at every non-boss spawn, killed or not, and each living boss as it is.
// A killed boss is no longer in enemies, so it stays dead.
func RespawnEnemies(spawns []EnemySpawn, enemies []*EnemyRuntime) []*EnemyRuntime {
	revived := make([]*EnemyRuntime, 0, len(spawns))
	for _, e := range enemies {
		if e.bp.Boss {
			revived = append(revived, e)
		}
	}
	for _, s := range spawns {
		if !blueprintFor(s.Kind).Boss {
			revived = append(revived, NewEnemy(s.Kind, s.Pos))
		}
	}
	return revived
}

// SpawnEnemies creates one enemy per spawn.
func SpawnEnemies(spawns []EnemySpawn) []*EnemyRuntime {
	enemies := make([]*EnemyRuntime, 0, len(spawns))
	for _, s := range spawns {
		enemies = append(enemies, NewEnemy(s.Kind, s.Pos))
	}
	return enemies
}
