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
	bp.Mover = flyingMover{}
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

// TestFlyingEnemyHoversOverPlayer: with a HoverHeight the flyer chases to a
// point that far above the player's center, not into them.
func TestFlyingEnemyHoversOverPlayer(t *testing.T) {
	bp := testFlyerBlueprint
	bp.Mover = flyingMover{HoverHeight: 200}
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

// TestStrandedFlyerPatrolsWhereItIs: a flyer whose patrol path lies beyond a
// wall must not pin itself against the wall, turned back toward its path
// every tick; it patrols where it is instead, like a stranded ground enemy.
func TestStrandedFlyerPatrolsWhereItIs(t *testing.T) {
	wall := AABB{X: 400, Y: waveGroundY - 400, Width: LevelTileWidth, Height: 400}
	qt := waveWorld(0, 1200, wall)
	e := newTestFlyer(Position{X: 300, Y: waveGroundY - 300})
	e.spawnX = 1000 // home is past the wall
	p := farPlayer(e)

	for i := 0; i < 4*int(tps()); i++ {
		UpdateEnemy(e, p, qt)
	}
	if right := e.Pos.X + e.bp.BodyWidth; right > wall.X-20 {
		t.Errorf("after 4 s: right edge %v, want it to have flown away from the wall at %v (pinned there)", right, wall.X)
	}
}
