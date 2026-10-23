package core

import "testing"

// These tests pin down how the succubus behaves today, so the enemy-types
// refactor (movers, attackers, clip tables) can't change it unnoticed.

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

// clipTicks is how many ticks a clip takes to play once.
func clipTicks(anim *Animation) int {
	return int(float64(anim.TotalFrames)/anim.AnimationSpeed*tps()) + 1
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
	e := newTestEnemy(Position{X: 300, Y: waveGroundY - blueprintFor(KindSuccubus).BodyHeight})
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
		anim := e.Animations[int(st)]

		tickEnemyAnim(e, 1)
		if e.State != st || e.CurrAnimFrame != anim.AnimStartFrame {
			t.Fatalf("state %d first tick: state=%d frame=%d, want to start the clip", st, e.State, e.CurrAnimFrame)
		}
		e.PrevState = e.State
		tickEnemyAnim(e, clipTicks(anim))
		if e.State != EnemyPatrol {
			t.Errorf("state %d after its clip: state=%d, want EnemyPatrol", st, e.State)
		}
	}
}

func TestEnemyDeathHoldsLastFrameThenRemovable(t *testing.T) {
	e := newAnimatedTestEnemy(Position{})
	e.TakeDamage(e.Health)
	e.PrevState = e.State
	anim := e.Animations[int(EnemyDead)]
	last := anim.AnimStartFrame + anim.TotalFrames - 1

	tickEnemyAnim(e, clipTicks(anim)/2)
	if e.Removable() {
		t.Fatal("removable halfway through the death clip")
	}
	tickEnemyAnim(e, clipTicks(anim))
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
