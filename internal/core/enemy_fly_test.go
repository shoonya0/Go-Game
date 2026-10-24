package core

import (
	"math"
	"testing"
)

// testFlyerBlueprint is the test enemy with wings: same stats and clips, moved
// by a flyingMover.
var testFlyerBlueprint = func() EnemyBlueprint {
	bp := testEnemyBlueprint
	bp.Name = "Test flyer"
	bp.Marker = "~~"
	bp.Mover = flyingMover{LeashTiles: 3}
	return bp
}()

// newTestFlyer is newTestEnemy with the test flyer's blueprint.
func newTestFlyer(pos Position) *EnemyRuntime {
	e := newTestEnemy(pos)
	e.bp = &testFlyerBlueprint
	return e
}

func TestFlyingEnemyIgnoresGravity(t *testing.T) {
	qt := waveWorld(0, 1200)
	const y = waveGroundY - 400
	e := newTestFlyer(Position{X: 300, Y: y})
	p := farPlayer(e)

	for i := 0; i < 3*int(tps()); i++ {
		UpdateEnemy(e, p, qt)
		if e.Pos.Y != y || e.bp.Mover.Grounded(e) {
			t.Fatalf("tick %d: y=%v grounded=%v, want to stay up at %v", i, e.Pos.Y, e.bp.Mover.Grounded(e), y)
		}
	}
	if e.Pos.X == 300 {
		t.Error("flyer didn't patrol: x still 300")
	}
}

// TestFlyingEnemyLeash: a flyer chases a player it sees no farther than its
// leash around its spawn, then flies back home when the player leaves.
func TestFlyingEnemyLeash(t *testing.T) {
	qt := waveWorld(0, 2000)
	spawn := Position{X: 300, Y: waveGroundY - 300}
	e := newTestFlyer(spawn)
	leash := 3.0 * LevelTileWidth
	p := &PlayerRuntime{Pos: Position{X: 700, Y: waveGroundY - bodyHeight}, Combat: Combat{Health: 100, MaxHealth: 100}}

	for i := 0; i < 4*int(tps()); i++ {
		UpdateEnemy(e, p, qt)
		if math.Abs(e.Pos.X-spawn.X) > leash || math.Abs(e.Pos.Y-spawn.Y) > leash {
			t.Fatalf("tick %d: at (%v, %v), want within %v of spawn (%v, %v)", i, e.Pos.X, e.Pos.Y, leash, spawn.X, spawn.Y)
		}
	}
	if e.State != EnemyAlert || e.Pos.X < spawn.X+leash-1 || e.Pos.Y < spawn.Y+leash-1 {
		t.Errorf("chasing: state=%d at (%v, %v), want EnemyAlert at the leash corner toward the player", e.State, e.Pos.X, e.Pos.Y)
	}

	p.Pos.X = 3000 // out of sight
	for i := 0; i < 5*int(tps()); i++ {
		UpdateEnemy(e, p, qt)
	}
	if e.State != EnemyPatrol || math.Abs(e.Pos.Y-spawn.Y) > 0.01 || math.Abs(e.Pos.X-spawn.X) > leash {
		t.Errorf("after the player left: state=%d at (%v, %v), want patrolling back at height %v near x=%v",
			e.State, e.Pos.X, e.Pos.Y, spawn.Y, spawn.X)
	}
}

// TestFlyingEnemyHoversOverPlayer: with a HoverHeight the flyer chases to a
// point that far above the player's center, not into them.
func TestFlyingEnemyHoversOverPlayer(t *testing.T) {
	bp := testFlyerBlueprint
	bp.Mover = flyingMover{LeashTiles: 5, HoverHeight: 200}
	qt := waveWorld(0, 1200)
	p := &PlayerRuntime{Pos: Position{X: 400, Y: waveGroundY - bodyHeight}, Combat: Combat{Health: 100, MaxHealth: 100}}
	e := newTestFlyer(Position{X: 250, Y: waveGroundY - 500})
	e.bp = &bp

	for i := 0; i < 3*int(tps()); i++ {
		UpdateEnemy(e, p, qt)
	}
	eb, pb := e.GetBounds(), p.GetBounds()
	if math.Abs(centerX(eb)-centerX(pb)) > 1 || math.Abs(centerY(eb)-(centerY(pb)-200)) > 1 {
		t.Errorf("center (%v, %v), want (%v, %v): 200 px above the player", centerX(eb), centerY(eb), centerX(pb), centerY(pb)-200)
	}
}

// TestFlyingEnemyStopsAtWall: tiles are solid to a flyer too.
func TestFlyingEnemyStopsAtWall(t *testing.T) {
	wall := AABB{X: 400, Y: waveGroundY - 400, Width: LevelTileWidth, Height: 400}
	qt := waveWorld(0, 1200, wall)
	e := newTestFlyer(Position{X: 300, Y: waveGroundY - 300})
	p := farPlayer(e)

	for i := 0; i < 3*int(tps()); i++ {
		UpdateEnemy(e, p, qt)
		if right := e.Pos.X + e.bp.BodyWidth; right > wall.X {
			t.Fatalf("tick %d: right edge %v went into the wall at %v", i, right, wall.X)
		}
	}
}

// TestDeadFlyerFalls: a flyer stops flying when it dies and drops to the ground.
func TestDeadFlyerFalls(t *testing.T) {
	qt := waveWorld(0, 1200)
	e := newTestFlyer(Position{X: 300, Y: waveGroundY - 400})
	e.TakeDamage(e.Health)

	for i := 0; i < 3*int(tps()); i++ {
		UpdateEnemy(e, farPlayer(e), qt)
	}
	if want := waveGroundY - e.bp.BodyHeight; e.Pos.Y != want {
		t.Errorf("dead flyer at y=%v, want fallen to %v", e.Pos.Y, want)
	}
}
