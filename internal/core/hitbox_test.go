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

	e.decideBehavior(above)
	if e.State == EnemyAttack || above.Combat.Health != 100 {
		t.Errorf("player above: state=%d health=%d, want no attack", e.State, above.Combat.Health)
	}
}

func TestEnemyAttacksPlayerBehind(t *testing.T) {
	e := newTestEnemy(Position{X: 200})
	e.FlipX = false // facing right, player is to the left
	behind := &PlayerRuntime{Pos: Position{X: 170}, Combat: Combat{Health: 100, MaxHealth: 100}}

	e.decideBehavior(behind)
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

// EnemyRuntime must keep satisfying the overlay interface.
var _ HitboxProvider = (*EnemyRuntime)(nil)
