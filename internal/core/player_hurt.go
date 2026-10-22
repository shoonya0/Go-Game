package core

import "math/rand/v2"

// Heavy hits, exhaustion and death.
const (
	HeavyHitDamage    = 2 * MaxPlayerHealth / hudSegments // a hit that takes 2 hearts or more knocks the player down
	KnockbackDistance = 1.5 * LevelTileWidth              // px a heavy hit pushes the player away from its source
	KnockbackSeconds  = 0.4                               // how long the push lasts (the tumbling frames of the knockdown)
)

// deathKneels picks which death animation plays; tests replace it to choose one.
var deathKneels = func() bool { return rand.IntN(2) == 0 }

// ApplyDamageFrom is ApplyDamage for a hit coming from world x srcX: a heavy hit
// knocks the player away from srcX.
func (p *PlayerRuntime) ApplyDamageFrom(dmg int, srcX float64) {
	dir := p.backward()
	if c := centerX(p.GetBounds()); c != srcX {
		dir = 1
		if c < srcX {
			dir = -1
		}
	}
	p.hurt(dmg, dir)
}

// backward is the push direction for a hit without a source: opposite the way
// the player faces.
func (p *PlayerRuntime) backward() float64 {
	if p.FlipX {
		return 1
	}
	return -1
}

// hurt applies dmg and, for a heavy hit the player survives, starts the
// knockdown and a push of KnockbackDistance in direction dir (-1 or 1). A dying
// player takes no more damage.
func (p *PlayerRuntime) hurt(dmg int, dir float64) {
	if dmg <= 0 || p.State.IsDying() {
		return
	}
	p.damageFlash = DamageFlashSeconds
	p.Combat.Health = max(p.Combat.Health-dmg, 0)
	if dmg < HeavyHitDamage || p.Combat.Health == 0 {
		return
	}
	p.State.SetPlayerState(int(PlayerStateDamaged))
	p.knockTicks = int(KnockbackSeconds*tps() + 0.5)
	p.knockDir = dir
}

// knockbackSpeed is the push velocity (px/s) that covers KnockbackDistance in
// knockTicks ticks.
func knockbackSpeed() float64 {
	ticks := int(KnockbackSeconds*tps() + 0.5)
	return KnockbackDistance / (float64(ticks) * deltaTime())
}

// applyKnockback drives VelX during a knockdown: the push while it lasts, then
// standing still until the animation ends. It reports whether it took over
// horizontal movement this tick.
func (p *PlayerRuntime) applyKnockback() bool {
	switch {
	case p.knockTicks > 0:
		p.knockTicks--
		p.Physics.VelX = p.knockDir * knockbackSpeed()
	case p.State.IsDamaged() || p.State.IsDying():
		p.Physics.VelX = 0
	default:
		return false
	}
	p.running = false
	return true
}

// startDeath plays one of the two death animations, picked at random. The
// player respawns when it ends (UpdateAnimation).
func (p *PlayerRuntime) startDeath() {
	if deathKneels() {
		p.State.SetPlayerState(int(PlayerStateDeadKneel))
	} else {
		p.State.SetPlayerState(int(PlayerStateDead))
	}
}

// exhaustedFrame keeps the out-of-stamina animation in step with the stamina
// bar: its first frame at empty, its last once one bar (RunRecoverPower) refills.
func (p *PlayerRuntime) exhaustedFrame(anim *Animation) int {
	f := int(p.Combat.Power / RunRecoverPower * float64(anim.TotalFrames))
	return anim.AnimStartFrame + min(max(f, 0), anim.TotalFrames-1)
}
