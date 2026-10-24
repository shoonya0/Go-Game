package core

import "testing"

// These tests pin down how the first enemy, the succubus, behaved (now the
// test-only testEnemyBlueprint), so the enemy-types refactor (movers,
// attackers, clip tables) can't change it unnoticed.

// newAnimatedTestEnemy is newTestEnemy plus its animation table, for tests that
// tick UpdateEnemyAnimation. It still loads no sprite sheet.
func newAnimatedTestEnemy(pos Position) *EnemyRuntime {
	e := newTestEnemy(pos)
	e.Animations = e.bp.newAnimations()
	return e
}

// farPlayer is a player well outside the enemy's detection range.
func farPlayer(e *EnemyRuntime) *PlayerRuntime {
	return &PlayerRuntime{Pos: Position{X: e.Pos.X + e.bp.Stats.DetectRange + 1000}, Combat: Combat{Health: 100, MaxHealth: 100}}
}

// tickEnemyAnim advances the enemy's animation n ticks.
func tickEnemyAnim(e *EnemyRuntime, n int) {
	for range n {
		e.UpdateEnemyAnimation()
	}
}

// clipTicks is how many ticks the enemy's clip for state st takes to play
// once, held frames included.
func clipTicks(e *EnemyRuntime, st EnemyState) int {
	return int(e.bp.Clips[e.clipFor(st)].duration()*tps()) + 1
}

func TestEnemyLandsOnGround(t *testing.T) {
	qt := waveWorld(0, 1200)
	e := newTestEnemy(Position{X: 300, Y: waveGroundY - 400})
	p := farPlayer(e)

	for i := 0; i < 3*int(tps()); i++ {
		UpdateEnemy(e, p, qt)
	}
	if want := waveGroundY - e.bp.BodyHeight; e.Pos.Y != want || !e.onGround || e.VelY != 0 {
		t.Errorf("after falling: y=%v onGround=%v velY=%v, want y=%v on ground at rest", e.Pos.Y, e.onGround, e.VelY, want)
	}
}

func TestEnemyPatrolBouncesOffWall(t *testing.T) {
	wall := AABB{X: 400, Y: waveGroundY - 2*LevelTileHeight, Width: LevelTileWidth, Height: 2 * LevelTileHeight}
	qt := waveWorld(0, 1200, wall)
	e := newTestEnemy(Position{X: 300, Y: waveGroundY - testEnemyBlueprint.BodyHeight})
	p := farPlayer(e)

	for i := 0; i < 2*int(tps()); i++ {
		UpdateEnemy(e, p, qt)
		if right := e.Pos.X + e.bp.BodyWidth; right > wall.X {
			t.Fatalf("tick %d: enemy right edge %v went into the wall at %v", i, right, wall.X)
		}
	}
	if e.patrolDir != -1 || !e.FlipX {
		t.Errorf("after reaching the wall: patrolDir=%v flipX=%v, want -1 and facing left", e.patrolDir, e.FlipX)
	}
}

func TestGroundEnemyTurnsAtLedge(t *testing.T) {
	const floorEnd = 12 * LevelTileWidth // 720, on a tile edge
	qt := waveWorld(0, floorEnd)
	e := newTestEnemy(Position{X: 640, Y: waveGroundY - testEnemyBlueprint.BodyHeight})
	p := farPlayer(e) // patrol path reaches past the edge: 640 + 4 tiles

	for i := 0; i < 4*int(tps()); i++ {
		UpdateEnemy(e, p, qt)
		if right := e.Pos.X + e.bp.BodyWidth; right > floorEnd || !e.onGround {
			t.Fatalf("tick %d: right edge %v onGround=%v, want to stay on the floor ending at %v", i, right, e.onGround, floorEnd)
		}
	}
	if e.patrolDir != -1 || !e.FlipX || e.Pos.X >= 640 {
		t.Errorf("after reaching the ledge: x=%v patrolDir=%v flipX=%v, want walking back left of 640", e.Pos.X, e.patrolDir, e.FlipX)
	}
}

// TestStrandedEnemyPatrolsWhereItIs covers an enemy that chased the player
// down to another platform and lost them: its patrol path points back to its
// spawn, past a ledge. It must not freeze at the ledge, turning back toward
// the spawn every tick; it patrols where it is instead.
func TestStrandedEnemyPatrolsWhereItIs(t *testing.T) {
	const floorEnd = 12 * LevelTileWidth // 720, on a tile edge
	qt := waveWorld(0, floorEnd)
	e := newTestEnemy(Position{X: 640, Y: waveGroundY - testEnemyBlueprint.BodyHeight})
	e.spawnX = 2000 // spawn is far to the right, past the edge
	p := farPlayer(e)

	for i := 0; i < 4*int(tps()); i++ {
		UpdateEnemy(e, p, qt)
		if right := e.Pos.X + e.bp.BodyWidth; right > floorEnd || !e.onGround {
			t.Fatalf("tick %d: right edge %v onGround=%v, want to stay on the floor ending at %v", i, right, e.onGround, floorEnd)
		}
	}
	if e.Pos.X > 600 {
		t.Errorf("after 4 s: x=%v, want it to have walked back from the ledge (frozen there)", e.Pos.X)
	}
}

func TestEnemyChasesLeft(t *testing.T) {
	e := newTestEnemy(Position{X: 1000})
	p := &PlayerRuntime{Pos: Position{X: 1000 - e.bp.Stats.DetectRange + 10}, Combat: Combat{Health: 100, MaxHealth: 100}}
	e.decideBehavior(p)
	if e.State != EnemyAlert || e.VelX != -e.bp.Stats.ChaseSpeed || !e.FlipX {
		t.Errorf("player left in range: state=%d velX=%v flipX=%v, want EnemyAlert at -%v facing left",
			e.State, e.VelX, e.FlipX, e.bp.Stats.ChaseSpeed)
	}
}

func TestHurtEnemyHoldsStill(t *testing.T) {
	e := newAnimatedTestEnemy(Position{})
	e.VelX = 100
	e.TakeDamage(1)
	near := &PlayerRuntime{Pos: Position{X: 30}, Combat: Combat{Health: 100, MaxHealth: 100}}

	e.decideBehavior(near)
	if e.State != EnemyHurt || e.VelX != 0 || near.Combat.Health != 100 {
		t.Errorf("while hurt: state=%d velX=%v player health=%d, want EnemyHurt, still, no hit",
			e.State, e.VelX, near.Combat.Health)
	}
}

func TestEnemyOneShotClipsReturnToPatrol(t *testing.T) {
	for _, st := range []EnemyState{EnemyAttack, EnemyHurt} {
		e := newAnimatedTestEnemy(Position{})
		e.State, e.PrevState = st, EnemyPatrol // entered this tick
		anim := e.animFor(st)

		tickEnemyAnim(e, 1)
		if e.State != st || e.CurrAnimFrame != anim.AnimStartFrame {
			t.Fatalf("state %d first tick: state=%d frame=%d, want to start the clip", st, e.State, e.CurrAnimFrame)
		}
		e.PrevState = e.State
		tickEnemyAnim(e, clipTicks(e, st))
		if e.State != EnemyPatrol {
			t.Errorf("state %d after its clip: state=%d, want EnemyPatrol", st, e.State)
		}
	}
}

func TestEnemyDeathHoldsLastFrameThenRemovable(t *testing.T) {
	e := newAnimatedTestEnemy(Position{})
	e.TakeDamage(e.Health)
	e.PrevState = e.State
	anim := e.animFor(EnemyDead)
	last := anim.AnimStartFrame + anim.TotalFrames - 1

	tickEnemyAnim(e, clipTicks(e, EnemyDead)/2)
	if e.Removable() {
		t.Fatal("removable halfway through the death clip")
	}
	tickEnemyAnim(e, clipTicks(e, EnemyDead))
	if !e.Removable() || e.CurrAnimFrame != last || e.State != EnemyDead {
		t.Errorf("after death clip: removable=%v frame=%d state=%d, want removable on frame %d, still dead",
			e.Removable(), e.CurrAnimFrame, e.State, last)
	}
	tickEnemyAnim(e, 10)
	if e.CurrAnimFrame != last {
		t.Errorf("death frame moved to %d after finishing, want held on %d", e.CurrAnimFrame, last)
	}
}

func TestDeadEnemyFallsWithoutActing(t *testing.T) {
	qt := waveWorld(0, 1200)
	e := newTestEnemy(Position{X: 300, Y: waveGroundY - 400})
	e.TakeDamage(e.Health)
	near := &PlayerRuntime{Pos: Position{X: 310, Y: waveGroundY - bodyHeight}, Combat: Combat{Health: 100, MaxHealth: 100}}

	for i := 0; i < 3*int(tps()); i++ {
		UpdateEnemy(e, near, qt)
	}
	if want := waveGroundY - e.bp.BodyHeight; e.Pos.Y != want || e.Pos.X != 300 {
		t.Errorf("dead enemy at (%v, %v), want fallen straight down to (300, %v)", e.Pos.X, e.Pos.Y, want)
	}
	if near.Combat.Health != 100 || e.State != EnemyDead {
		t.Errorf("dead enemy acted: player health=%d state=%d", near.Combat.Health, e.State)
	}
}
