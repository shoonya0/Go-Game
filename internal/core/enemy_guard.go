package core

import "math/rand/v2"

// rollGuard returns a number in [0, 1) that a guard check compares against
// GuardChance. A variable so tests can make the roll deterministic.
var rollGuard = func() float64 { return rand.Float64() }

// TakeHitFrom is TakeDamage for a player hit coming from world x srcX. An enemy
// that can guard (D8) may block a hit from the side it faces: it takes no
// damage and plays its guard clip instead.
func (e *EnemyRuntime) TakeHitFrom(dmg int, srcX float64) {
	if !e.blocks(srcX) {
		e.TakeDamage(dmg)
		return
	}
	e.VelX = 0
	e.setState(EnemyGuard)
	e.restartAnimation() // see TakeDamage: hits land outside UpdateEnemy
}

// blocks reports whether the enemy guards a hit from srcX: its species has a
// guard clip and a GuardChance, the hit comes from the front, and the roll
// passes. A dead enemy or one mid-swing can't guard.
func (e *EnemyRuntime) blocks(srcX float64) bool {
	if e.State.IsDead() || e.State.IsAttacking() || e.bp.Stats.GuardChance <= 0 {
		return false
	}
	if !e.bp.Has(ClipGuard) {
		return false
	}
	fromRight := srcX > centerX(e.GetBounds())
	if fromRight == e.FlipX { // FlipX means facing left
		return false
	}
	return rollGuard() < e.bp.Stats.GuardChance
}
