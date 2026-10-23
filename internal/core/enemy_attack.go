package core

import "math/rand/v2"

// Attacker is how a species fights: when the player is within reach, and what
// a swing does to them. The brain decides to attack; the Attacker decides
// whether and when the swing lands.
type Attacker interface {
	// InRange reports whether the enemy could strike target if it turned to face it.
	InRange(e *EnemyRuntime, target AABB) bool
	// Box is the area a strike covers with the enemy's current facing.
	Box(e *EnemyRuntime) AABB
	// Ready reports whether a swing started now would land; until then the
	// enemy waits instead of swinging.
	Ready(e *EnemyRuntime) bool
	// Start begins a new swing, as the enemy enters its attack state.
	Start(e *EnemyRuntime)
	// Update runs each tick of a swing while the target is in range.
	Update(e *EnemyRuntime, target *PlayerRuntime)
}

// swingState is the attack in progress. The zero value is a fresh attack_1 swing.
type swingState struct {
	second bool // plays attack_2 instead of attack_1
	landed bool // this swing has hit, so it can't hit again
}

// rollSecondAttack picks attack_2 over attack_1 for a new swing. A variable so
// tests can make the pick deterministic.
var rollSecondAttack = func() bool { return rand.IntN(2) == 0 }

// meleeAttack strikes the area in front of the body, AttackRange px deep. The
// swing lands on its clip's HitFrame (or later, once the cooldown is ready),
// once per swing and at most once per AttackCooldown.
type meleeAttack struct{}

func (meleeAttack) InRange(e *EnemyRuntime, target AABB) bool {
	dx := centerX(target) - centerX(e.GetBounds())
	return meleeBox(e, dx < 0).Intersects(target)
}

func (meleeAttack) Box(e *EnemyRuntime) AABB { return meleeBox(e, e.FlipX) }

// Ready holds a new swing back until the cooldown will be over by its hit
// frame. The cooldown can outlast a whole attack clip, so a swing started any
// sooner would end without landing.
func (meleeAttack) Ready(e *EnemyRuntime) bool {
	soonest := e.bp.Clips[ClipAttack1].hitTime()
	if c, ok := e.bp.Clips[ClipAttack2]; ok {
		soonest = min(soonest, c.hitTime())
	}
	return e.attackTimer <= soonest
}

func (meleeAttack) Start(e *EnemyRuntime) {
	_, hasSecond := e.bp.Clips[ClipAttack2]
	e.swing = swingState{second: hasSecond && rollSecondAttack()}
}

func (meleeAttack) Update(e *EnemyRuntime, target *PlayerRuntime) {
	// Ready lets a swing start only when the cooldown ends by its hit frame; if
	// tick rounding still blocks it there, it lands as soon as it's ready.
	if e.swing.landed || e.swingFrame() < e.currentClip().HitFrame || e.attackTimer > 0 {
		return
	}
	target.ApplyDamageFrom(e.bp.Stats.Damage, centerX(e.GetBounds()))
	e.swing.landed = true
	e.attackTimer = e.bp.Stats.AttackCooldown
}

// meleeBox is the strike area for an explicit facing (left = true): from the
// body's centre out AttackRange px, at body height. The AI tests it toward the
// player before turning to face them.
func meleeBox(e *EnemyRuntime, left bool) AABB {
	b := e.GetBounds()
	reach := e.bp.Stats.AttackRange
	x := centerX(b)
	if left {
		x -= reach
	}
	return AABB{X: x, Y: b.Y, Width: reach, Height: b.Height}
}

// hitTime is how long a swing with this clip takes to reach its HitFrame.
func (c EnemyClip) hitTime() float64 { return c.timeTo(c.HitFrame) }

// swingFrame is how many frames of the attack clip have played. On the tick
// the attack starts, the clip hasn't been rewound yet (UpdateEnemyAnimation
// does that after the AI runs), so it is 0.
func (e *EnemyRuntime) swingFrame() int {
	anim := e.currentAnim()
	if e.PrevState != EnemyAttack || anim == nil {
		return 0
	}
	return e.CurrAnimFrame - anim.AnimStartFrame
}
