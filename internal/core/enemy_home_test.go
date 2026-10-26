package core

import (
	"math"
	"testing"
)

// An enemy isn't tied to its spawn: once it loses the player, it guards the
// place where it lost them (#31).

// chaseThenLose runs e for chaseSecs against p, then moves p out of sight and
// returns where e was when it lost them.
func chaseThenLose(e *EnemyRuntime, p *PlayerRuntime, qt *DynamicQuadtree, chaseSecs float64) Position {
	for i := 0; i < int(chaseSecs*tps()); i++ {
		UpdateEnemy(e, p, qt)
	}
	p.Pos.X += 5000
	return e.Pos
}

func TestGroundEnemyGuardsWhereItLostPlayer(t *testing.T) {
	qt := waveWorld(0, 3000)
	e := newTestEnemy(Position{X: 300, Y: waveGroundY - testEnemyBlueprint.BodyHeight})
	p := &PlayerRuntime{Pos: Position{X: 700, Y: waveGroundY - bodyHeight}, Combat: Combat{Health: 1000, MaxHealth: 1000}}

	lost := chaseThenLose(e, p, qt, 3)
	path := e.bp.Stats.PatrolTiles * LevelTileWidth
	if lost.X-300 <= path {
		t.Fatalf("chase ended at x=%v, want it farther than its patrol path (%v) from spawn 300", lost.X, path)
	}
	for i := 0; i < 8*int(tps()); i++ {
		UpdateEnemy(e, p, qt)
		if math.Abs(e.Pos.X-lost.X) > path+1 {
			t.Fatalf("tick %d: x=%v, want it to patrol within %v of x=%v where it lost the player", i, e.Pos.X, path, lost.X)
		}
	}
	if e.State != EnemyPatrol {
		t.Errorf("state=%d, want EnemyPatrol", e.State)
	}
}

func TestFlyingEnemyChasesPastItsSpawn(t *testing.T) {
	qt := waveWorld(0, 3000)
	spawn := Position{X: 300, Y: waveGroundY - 300}
	e := newTestFlyer(spawn)
	p := &PlayerRuntime{Pos: Position{X: 700, Y: waveGroundY - bodyHeight}, Combat: Combat{Health: 1000, MaxHealth: 1000}}

	for i := 0; i < 4*int(tps()); i++ {
		UpdateEnemy(e, p, qt)
	}
	if !e.bp.Attacker.InRange(e, p.GetBounds()) {
		t.Errorf("after 4 s at (%v, %v): want it to have reached the player at (%v, %v)", e.Pos.X, e.Pos.Y, p.Pos.X, p.Pos.Y)
	}
}

func TestFlyingEnemyGuardsWhereItLostPlayer(t *testing.T) {
	qt := waveWorld(0, 3000)
	e := newTestFlyer(Position{X: 300, Y: waveGroundY - 300})
	p := &PlayerRuntime{Pos: Position{X: 700, Y: waveGroundY - bodyHeight}, Combat: Combat{Health: 1000, MaxHealth: 1000}}

	lost := chaseThenLose(e, p, qt, 4)
	path := e.bp.Stats.PatrolTiles * LevelTileWidth
	for i := 0; i < 8*int(tps()); i++ {
		UpdateEnemy(e, p, qt)
		if math.Abs(e.Pos.X-lost.X) > path+1 || math.Abs(e.Pos.Y-lost.Y) > 0.01 {
			t.Fatalf("tick %d: at (%v, %v), want it to patrol at y=%v within %v of x=%v", i, e.Pos.X, e.Pos.Y, lost.Y, path, lost.X)
		}
	}
}
