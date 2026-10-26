package core

import "testing"

// Walkers keep out of larva (playtest of #34): it hurts the player, and an
// enemy chasing across a larva pit shook at its edge.

// gridWorld builds a quadtree from a level grid, the way the game loads one.
func gridWorld(grid [][]int) *DynamicQuadtree {
	WorldInit()
	level := BuildLevel(grid)
	qt := NewDynamicQuadtree(AABB{X: 0, Y: 0, Width: 4000, Height: 2000})
	for i := range level.Platforms {
		qt.Insert(&level.Platforms[i])
	}
	return qt
}

// larvaPitWorld is a rock floor at row 6 (y=360) with a larva pit in columns
// 10-14 (x=600..900), level with the rock.
func larvaPitWorld() *DynamicQuadtree {
	const rock, larva = 4, 8
	row := func() []int {
		r := make([]int, 25)
		for i := range r {
			r[i] = rock
			if i >= 10 && i < 15 {
				r[i] = larva
			}
		}
		return r
	}
	empty := make([]int, 25)
	return gridWorld([][]int{empty, empty, empty, empty, empty, empty, row(), row()})
}

const (
	larvaFloorY = 6 * LevelTileHeight
	larvaPitX   = 10 * LevelTileWidth
)

// TestWalkerStaysOutOfLarva: a walker that sees the player across a larva pit
// doesn't chase into it or shake at its edge: it stops chasing and patrols its
// own side.
func TestWalkerStaysOutOfLarva(t *testing.T) {
	qt := larvaPitWorld()
	bp := testEnemyBlueprint
	bp.Stats.DetectRange = 760 // a skeleton's
	e := newTestEnemy(Position{X: 7 * LevelTileWidth, Y: larvaFloorY - bp.BodyHeight})
	e.bp = &bp
	p := playerOn(19*LevelTileWidth, larvaFloorY)

	flips := 0
	for i := 0; i < 5*int(tps()); i++ {
		before := e.State
		UpdateEnemy(e, p, qt)
		e.UpdateEnemyAnimation()
		if right := e.Pos.X + bp.BodyWidth; right > larvaPitX {
			t.Fatalf("tick %d: right edge %v went into the larva at %v", i, right, larvaPitX)
		}
		if before != e.State {
			flips++
		}
	}
	if flips > 1 || e.State != EnemyPatrol {
		t.Errorf("state changed %d times, now %d; want it to settle into EnemyPatrol", flips, e.State)
	}
}

// TestWalkerWalksOutOfLarva: a walker standing in larva (spawned or fallen
// there) still chases a player on the rock beside the pit, and climbs out.
func TestWalkerWalksOutOfLarva(t *testing.T) {
	qt := larvaPitWorld()
	larvaTop := float64(larvaFloorY + LarvaTopOffset)
	e := newTestEnemy(Position{X: 12 * LevelTileWidth, Y: larvaTop - testEnemyBlueprint.BodyHeight})
	p := playerOn(17*LevelTileWidth, larvaFloorY)

	for i := 0; i < 4*int(tps()) && e.Pos.Y+e.bp.BodyHeight != larvaFloorY; i++ {
		UpdateEnemy(e, p, qt)
	}
	if feet := e.Pos.Y + e.bp.BodyHeight; feet != larvaFloorY || e.Pos.X+e.bp.BodyWidth <= larvaPitX+5*LevelTileWidth {
		t.Errorf("at x=%v feet y=%v, want out of the pit onto the rock at y=%v", e.Pos.X, feet, larvaFloorY)
	}
}

// TestWalkerOverAGapEdgeStillChases: a walker held up by only part of its
// body (its centre over a gap) still sees the ground it stands on.
func TestWalkerOverAGapEdgeStillChases(t *testing.T) {
	const floorEnd = 600.0
	qt := waveWorld(0, floorEnd)
	e := floorEnemy(floorEnd - testEnemyBlueprint.BodyWidth/3) // centre past the edge
	e.decideBehavior(playerOn(300, waveGroundY), qt)
	if e.State != EnemyAlert {
		t.Errorf("state %d, want EnemyAlert: the player is on its floor", e.State)
	}
}
