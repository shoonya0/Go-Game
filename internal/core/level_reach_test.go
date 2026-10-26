package core

import (
	"encoding/json"
	"fmt"
	"testing"
)

// The reach test checks every level grid against the player's real jump:
// a simulation finds which climbs and gaps a jump makes, then a search over
// the grid's standing spots finds what each portal can reach.

// hopKey is a move from one standing spot to another: up tiles higher
// (negative: a drop), across gap empty columns, starting in water or not.
type hopKey struct {
	up, gap int
	water   bool
}

var hopCache = map[hopKey]bool{}

// canHop reports whether the player can make move k: it runs right across
// a ledge (or wades through a pool) and jumps from some spot before its
// edge, with and without sprinting, and must end standing on the far ledge.
func canHop(k hopKey) bool {
	if ok, seen := hopCache[k]; seen {
		return ok
	}
	ok := false
	for _, run := range []bool{false, true} {
		for back := 0.0; back <= 2*LevelTileWidth && !ok; back += 6 {
			ok = hopOnce(k, run, back)
		}
	}
	hopCache[k] = ok
	return ok
}

// hopOnce builds k's two ledges and tries one jump, taken back px before
// the edge of the first.
func hopOnce(k hopKey, run bool, back float64) bool {
	const cols, from = 24, 6 // the first ledge is columns 0..from-1
	top := 30                // its top row; the far ledge's is top-k.up
	far := top - k.up
	rows := max(top, far) + 4
	grid := make([][]int, rows)
	for r := range grid {
		grid[r] = make([]int, cols)
	}
	fill := func(c0, c1, r0 int, code int) {
		for c := c0; c < c1; c++ {
			for r := r0; r < rows; r++ {
				grid[r][c] = code
			}
		}
	}
	if k.water { // a pool two tiles deep, its surface at the top row
		fill(0, from, top, 1)
		fill(0, from, top+2, 4)
	} else {
		fill(0, from, top, 4)
	}
	fill(from+k.gap, cols, far, 4)

	level := BuildLevel(grid)
	qt := NewDynamicQuadtree(AABB{Width: cols * LevelTileWidth, Height: float64(rows * LevelTileHeight)})
	for i := range level.Platforms {
		qt.Insert(&level.Platforms[i])
	}
	takeoff := float64(from*LevelTileWidth) - bodyWidth - back
	p := InitPlayer(nil) // two tiles of run-up before the takeoff
	p.placeAt(Position{X: takeoff - 2*LevelTileWidth, Y: float64(top*LevelTileHeight) - bodyHeight})
	for range 30 { // land (in water: sink to wading depth)
		UpdatePlayer(p, &InputState{}, qt, nil)
	}
	jumped, still := false, 0
	for range int(6 * tps()) {
		in := InputState{Direction: Direction{LeftRight: 1}, RunJustPressed: run} // RunJustPressed is the held run key
		in.JumpJustPressed = !jumped && p.Pos.X >= takeoff                        // until the jump starts
		UpdatePlayer(p, &in, qt, nil)
		jumped = jumped || p.Physics.VelY < 0
		if jumped && p.Physics.VelY == 0 {
			still++
		} else {
			still = 0
		}
		// Landed: no fall speed two ticks running (a head bump stops it for
		// one). The state can't tell: without a sheet the jump clip never ends.
		if still >= 2 {
			return p.Pos.X+bodyWidth > float64((from+k.gap)*LevelTileWidth) && p.Pos.Y+bodyHeight <= float64(far*LevelTileHeight)+1
		}
	}
	return false
}

// spot is a cell the player can stand on: a solid tile, or a water surface,
// with room for the body above.
type spot struct{ r, c int }

func reachSpots(grid [][]int) (spots map[spot]bool, wet func(spot) bool, free func(r, c int) bool) {
	code := func(r, c int) int {
		if r < 0 || r >= len(grid) || c < 0 || c >= len(grid[r]) {
			return 0
		}
		return grid[r][c]
	}
	solid := func(r, c int) bool {
		t := tileTypeFromCode(code(r, c))
		return t != Empty && t != Water
	}
	water := func(r, c int) bool { return tileTypeFromCode(code(r, c)) == Water }
	free = func(r, c int) bool { return !solid(r, c) && !water(r, c) }
	spots = map[spot]bool{}
	for r := range grid {
		for c := range grid[r] {
			if (solid(r, c) || water(r, c)) && free(r-1, c) && free(r-2, c) {
				spots[spot{r, c}] = true
			}
		}
	}
	return spots, func(s spot) bool { return water(s.r, s.c) }, free
}

// reachable returns the spots the player can get to from start.
func reachable(grid [][]int, start spot) map[spot]bool {
	spots, wet, free := reachSpots(grid)
	seen := map[spot]bool{start: true}
	queue := []spot{start}
	for len(queue) > 0 {
		a := queue[0]
		queue = queue[1:]
		for b := range spots {
			d := b.c - a.c
			if seen[b] || d == 0 || d < -6 || d > 6 {
				continue
			}
			step := 1
			if d < 0 {
				step, d = -1, -d
			}
			clear := true // the air between must be free for the body
			for c := a.c + step; c != b.c && clear; c += step {
				for r := min(a.r, b.r) - 2; r < max(a.r, b.r); r++ {
					clear = clear && free(r, c)
				}
			}
			if clear && canHop(hopKey{up: a.r - b.r, gap: d - 1, water: wet(a)}) {
				seen[b] = true
				queue = append(queue, b)
			}
		}
	}
	return seen
}

// levelGrid reads def's grid of tile codes, enemy markers left out.
func levelGrid(t *testing.T, def LevelDef) [][]int {
	t.Helper()
	var bundle levelBundle
	if err := json.Unmarshal(quoteBareMarkers(readAsset(def.File)), &bundle); err != nil {
		t.Fatal(err)
	}
	grid, _, err := splitGrid(bundle[def.Key].Grid)
	if err != nil {
		t.Fatal(err)
	}
	return grid
}

// TestLevelPortalsReachEachOther: in every level, from each portal the player
// can walk and jump to every other portal of that level.
func TestLevelPortalsReachEachOther(t *testing.T) {
	WorldInit()
	for name, def := range Levels {
		t.Run(name, func(t *testing.T) {
			grid := levelGrid(t, def)
			var portals []spot
			for _, c := range LoadLevel(def.File, def.Key).Portals {
				portals = append(portals, spot{int(c.Pos.Y)/LevelTileHeight + 1, int(c.Pos.X) / LevelTileWidth})
			}
			for i, a := range portals {
				got := reachable(grid, a)
				for j, b := range portals {
					if !got[b] {
						t.Errorf("Portal %d (row %d, col %d) can't reach Portal %d (row %d, col %d)", i+1, a.r-1, a.c, j+1, b.r-1, b.c)
					}
				}
			}
		})
	}
}

// TestHopRules pins the jump rules level grids are drawn to (levels state):
// climbs of 3 tiles over gaps of up to 2, 2 over 3, flat over 4; never 4 up.
func TestHopRules(t *testing.T) {
	WorldInit()
	for _, c := range []struct {
		k    hopKey
		want bool
	}{
		{hopKey{3, 2, false}, true}, {hopKey{2, 3, false}, true}, {hopKey{0, 4, false}, true},
		{hopKey{4, 0, false}, false}, {hopKey{1, 1, true}, true},
	} {
		if got := canHop(c.k); got != c.want {
			t.Errorf("%s: %v, want %v", fmt.Sprintf("up %d, gap %d, water %v", c.k.up, c.k.gap, c.k.water), got, c.want)
		}
	}
}
