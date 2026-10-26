package core

import "testing"

func TestEnemyAttackBoxFacesForward(t *testing.T) {
	e := newTestEnemy(Position{X: 100, Y: 50})
	body := e.GetBounds()
	reach := e.bp.Stats.AttackRange

	e.FlipX = false // facing right
	right := e.AttackBox()
	if right.X != centerX(body) || right.Width != reach || right.Y != body.Y || right.Height != body.Height {
		t.Errorf("facing right: box = %+v, want from centre %v out %v at body height", right, centerX(body), reach)
	}

	e.FlipX = true // facing left
	left := e.AttackBox()
	if left.X+left.Width != centerX(body) || left.Width != reach {
		t.Errorf("facing left: box = %+v, want ending at centre %v", left, centerX(body))
	}
}

// The attack used to check horizontal distance only, so a player standing on a
// platform right above the enemy still got hit.
func TestEnemyDoesNotHitPlayerAbove(t *testing.T) {
	e := newTestEnemy(Position{X: 0, Y: 500})
	above := &PlayerRuntime{Pos: Position{X: 30, Y: 500 - 300}, Combat: Combat{Health: 100, MaxHealth: 100}}

	e.decideBehavior(above, nil)
	if e.State == EnemyAttack || above.Combat.Health != 100 {
		t.Errorf("player above: state=%d health=%d, want no attack", e.State, above.Combat.Health)
	}
}

func TestEnemyAttacksPlayerBehind(t *testing.T) {
	e := newTestEnemy(Position{X: 200})
	e.FlipX = false // facing right, player is to the left
	behind := &PlayerRuntime{Pos: Position{X: 170}, Combat: Combat{Health: 100, MaxHealth: 100}}

	e.decideBehavior(behind, nil)
	if e.State != EnemyAttack || !e.FlipX {
		t.Errorf("player behind: state=%d flipX=%v, want attack after turning left", e.State, e.FlipX)
	}
}

func TestEnemyHitboxes(t *testing.T) {
	e := newTestEnemy(Position{X: 0})

	boxes := e.Hitboxes()
	if len(boxes) != 2 || boxes[0].Kind != HitboxBody || boxes[1].Kind != HitboxAttack {
		t.Fatalf("alive enemy hitboxes = %+v, want body + attack", boxes)
	}
	if boxes[0].Box != e.GetBounds() || boxes[1].Box != e.AttackBox() {
		t.Errorf("hitboxes don't match GetBounds/AttackBox")
	}
	if boxes[1].Active {
		t.Error("attack box active while patrolling")
	}

	e.setState(EnemyAttack)
	if !e.Hitboxes()[1].Active {
		t.Error("attack box not active while attacking")
	}

	e.setState(EnemyDead)
	if got := e.Hitboxes(); len(got) != 1 || got[0].Kind != HitboxBody {
		t.Errorf("dead enemy hitboxes = %+v, want body only", got)
	}
}

// attackingPlayer returns a player at pos, mid-way into state on animation frame.
func attackingPlayer(pos Position, state PlayerStateType, frame int) *PlayerRuntime {
	p := staminaPlayer()
	p.Pos = pos
	p.State.CurrentState = state
	p.PreviousState.CurrentState = state
	p.CurrAnimFrame = frame
	return p
}

func TestPlayerAttackBoxFacesForward(t *testing.T) {
	p := attackingPlayer(Position{X: 100, Y: 50}, PlayerStateWeakAttack, 2)
	body := p.GetBounds()
	reach := playerStrikes[PlayerStateWeakAttack].Reach

	right, striking, ok := p.AttackBox()
	if !ok || !striking || right.X != centerX(body) || right.Width != reach || right.Y+right.Height != body.Y+body.Height {
		t.Errorf("facing right: box=%+v striking=%v ok=%v, want from centre %v out %v, feet aligned", right, striking, ok, centerX(body), reach)
	}

	p.FlipX = true
	if left, _, _ := p.AttackBox(); left.X+left.Width != centerX(body) {
		t.Errorf("facing left: box = %+v, want ending at centre %v", left, centerX(body))
	}

	p.State.CurrentState = PlayerStateIdle
	if _, _, ok := p.AttackBox(); ok {
		t.Error("idle player has an attack box")
	}
}

func TestPlayerHitsEnemyOncePerSwing(t *testing.T) {
	e := newTestEnemy(Position{X: 160, Y: 0})
	p := attackingPlayer(Position{X: 100, Y: 0}, PlayerStateWeakAttack, 0)
	enemies := []*EnemyRuntime{e}
	full := e.Health

	p.HitEnemies(enemies) // wind-up frame: sword not out yet
	if e.Health != full {
		t.Fatalf("hit on wind-up frame: health %d", e.Health)
	}

	for f := 2; f <= 4; f++ { // every striking frame of the same swing
		p.CurrAnimFrame = f
		p.HitEnemies(enemies)
	}
	if e.Health != full-WeakAttackDamage || !e.State.IsHurt() {
		t.Fatalf("after one swing: health=%d state=%d, want %d and hurt", e.Health, e.State, full-WeakAttackDamage)
	}

	p.PreviousState.CurrentState = PlayerStateIdle // next tick starts a new swing
	p.CurrAnimFrame = 2
	p.HitEnemies(enemies)
	if e.Health != full-2*WeakAttackDamage {
		t.Fatalf("second swing: health=%d, want %d", e.Health, full-2*WeakAttackDamage)
	}
}

func TestPlayerMissesEnemyBehindOrAbove(t *testing.T) {
	behind := newTestEnemy(Position{X: 20, Y: 0})    // left of a right-facing player
	above := newTestEnemy(Position{X: 150, Y: -300}) // on a platform above
	p := attackingPlayer(Position{X: 100, Y: 0}, PlayerStateStrongAttack, 3)

	p.HitEnemies([]*EnemyRuntime{behind, above})
	if behind.Health != behind.bp.Stats.MaxHealth || above.Health != above.bp.Stats.MaxHealth {
		t.Errorf("health behind=%d above=%d, want both untouched", behind.Health, above.Health)
	}
}

func TestPlayerHitboxes(t *testing.T) {
	p := attackingPlayer(Position{}, PlayerStateIdle, 0)
	if got := p.Hitboxes(); len(got) != 1 || got[0].Kind != HitboxBody || got[0].Box != p.GetBounds() {
		t.Fatalf("idle player hitboxes = %+v, want body only", got)
	}

	p.State.CurrentState = PlayerStateStrongAttackInAir
	p.CurrAnimFrame = 0
	boxes := p.Hitboxes()
	if len(boxes) != 2 || boxes[1].Kind != HitboxAttack || boxes[1].Active {
		t.Fatalf("winding up: hitboxes = %+v, want body + inactive attack", boxes)
	}
	p.CurrAnimFrame = 5
	if !p.Hitboxes()[1].Active {
		t.Error("attack box not active on a striking frame")
	}
}

// Both actors must keep satisfying the overlay interface.
var (
	_ HitboxProvider = (*EnemyRuntime)(nil)
	_ HitboxProvider = (*PlayerRuntime)(nil)
)

// A player hit lands outside UpdateEnemy, so the hurt clip must restart from
// TakeDamage itself, including when the enemy is already hurt.
func TestEnemyHitRestartsHurtClip(t *testing.T) {
	e := newTestEnemy(Position{})
	e.Animations = e.bp.newAnimations()
	e.TakeDamage(1)
	e.CurrAnimFrame = 1 // part-way through the hurt clip
	e.PrevState = e.State

	e.TakeDamage(1)
	if want := e.animFor(EnemyHurt).AnimStartFrame; e.CurrAnimFrame != want {
		t.Errorf("hurt again: frame %d, want restart at %d", e.CurrAnimFrame, want)
	}
}
