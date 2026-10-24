package core

import "math"

// diveAttack smashes down onto a player below. It winds up in place
// (attack_1), falls straight down at FallSpeed (attack_loop) and lands
// (attack_end), hitting the player if they are in its Box; the end clip then
// hands back to the brain. Once started the dive plays out. A fall that finds
// no ground within MaxDrop px gives up in the air.
type diveAttack struct {
	FallSpeed float64 // px/s
	MaxDrop   float64 // px
}

// InRange: the player is below, their center within AttackRange of its own.
func (diveAttack) InRange(e *EnemyRuntime, target AABB) bool {
	b := e.GetBounds()
	return math.Abs(centerX(target)-centerX(b)) <= e.bp.Stats.AttackRange && target.Y >= b.Y+b.Height
}

// Box is the impact area: AttackRange either side of its center, at body height.
func (diveAttack) Box(e *EnemyRuntime) AABB {
	b := e.GetBounds()
	r := e.bp.Stats.AttackRange
	return AABB{X: centerX(b) - r, Y: b.Y, Width: 2 * r, Height: b.Height}
}

func (diveAttack) Ready(e *EnemyRuntime) bool { return e.attackTimer <= 0 }

func (diveAttack) Start(e *EnemyRuntime) { e.swing = swingState{} }

func (diveAttack) Committed(*EnemyRuntime) bool { return true }

// Update moves the dive through its phases. The brain has already steered the
// body to hold still; only the fall, and resting on the ground after it, move it.
func (d diveAttack) Update(e *EnemyRuntime, target *PlayerRuntime) {
	switch e.swing.phase {
	case phaseStart:
		// On the first tick clipTime still belongs to the previous clip.
		if e.PrevState == EnemyAttack && e.clipTime >= e.currentClip().duration() {
			d.next(e, phaseLoop)
		}
	case phaseLoop:
		landed := e.bp.Mover.Grounded(e)
		if !landed && e.clipTime*d.FallSpeed < d.MaxDrop {
			e.VelY = d.FallSpeed
			return
		}
		d.next(e, phaseEnd)
		e.attackTimer = e.bp.Stats.AttackCooldown
		if landed && d.Box(e).Intersects(target.GetBounds()) {
			target.ApplyDamageFrom(e.bp.Stats.Damage, centerX(e.GetBounds()))
			e.swing.landed = true
		}
		fallthrough
	case phaseEnd:
		if e.bp.Mover.Grounded(e) {
			e.VelY = d.FallSpeed // keep pressing on the ground, so it counts as grounded
		}
	}
}

// next moves the dive to phase p and starts that phase's clip.
func (diveAttack) next(e *EnemyRuntime, p attackPhase) {
	e.swing.phase = p
	e.restartAnimation()
}
