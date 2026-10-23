package core

import "math"

// UpdateEnemy advances one enemy by a single tick: it decides behaviour from the
// player's position, then lets the species' Mover move it through the world
// (enemy_move.go). It is the shared entry point for every enemy species — the
// per-species differences all live in the blueprint (enemy_types.go).
func UpdateEnemy(e *EnemyRuntime, target *PlayerRuntime, qt *DynamicQuadtree) {
	e.PrevState = e.State
	if e.attackTimer > 0 {
		e.attackTimer -= deltaTime()
	}

	// A dead enemy keeps falling to rest on the ground but takes no actions.
	if e.State.IsDead() {
		e.bp.Mover.Steer(e, moveIntent{mode: moveHold})
		e.bp.Mover.Move(e, qt)
		return
	}

	e.regenerateHealth()
	e.decideBehavior(target)
	e.bp.Mover.Move(e, qt)
}

// decideBehavior runs the enemy's state machine for this tick: attack when the
// player's body is inside the attack box (turned toward them), chase when they
// are within detection range, otherwise patrol. It sets the behavioural state and
// facing, and steers the Mover, which sets the velocity.
func (e *EnemyRuntime) decideBehavior(target *PlayerRuntime) {
	// Let a hurt reaction play out before re-deciding, so hits interrupt cleanly.
	if e.State.IsHurt() {
		e.bp.Mover.Steer(e, moveIntent{mode: moveHold})
		return
	}

	st := e.bp.Stats
	dx := centerX(target.GetBounds()) - centerX(e.GetBounds())
	dist := math.Abs(dx)

	switch {
	case e.bp.Attacker.InRange(e, target.GetBounds()):
		e.attack(target, dx)
	case dist <= st.DetectRange:
		e.chase(dx)
	default:
		e.patrol()
	}
}

// attack faces the player and swings, standing still. A new swing starts each
// time the enemy enters its attack state; the Attacker decides when it lands.
func (e *EnemyRuntime) attack(target *PlayerRuntime, dx float64) {
	e.faceToward(dx)
	e.bp.Mover.Steer(e, moveIntent{mode: moveHold})
	if e.State != EnemyAttack {
		e.bp.Attacker.Start(e)
	}
	e.setState(EnemyAttack)
	e.bp.Attacker.Update(e, target)
}

// chase moves the enemy toward the player at chase speed.
func (e *EnemyRuntime) chase(dx float64) {
	e.setState(EnemyAlert)
	e.faceToward(dx)
	e.bp.Mover.Steer(e, moveIntent{mode: moveChase, dx: dx})
}

// patrol paces the enemy around its spawn point; the Mover decides where the
// path turns.
func (e *EnemyRuntime) patrol() {
	e.bp.Mover.Steer(e, moveIntent{mode: movePatrol})
	e.faceToward(e.patrolDir)
	e.setState(EnemyPatrol)
}

// centerX returns the horizontal center of a box.
func centerX(b AABB) float64 { return b.X + b.Width/2 }
