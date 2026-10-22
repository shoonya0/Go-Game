package core

import "math"

// ------------------------ combat constants ------------------------
const (
	MaxPlayerHealth = 100 // starting/maximum health (integer)
	MaxPlayerPower  = 100 // starting/maximum special power

	LarvaDamage          = 10  // health lost each time the player enters larva
	WaterDamagePerSecond = 2.0 // health drained per second while submerged in water
	FallSafeTiles        = 5   // falls up to this many tiles do no damage; each extra tile costs 1

	PowerRegenPerSecond = 10.0 // special power regenerated per second
	SpecialPowerCost    = 34.0 // power spent per U special attack (~3 on a full bar)

	// O is a two-press combo (see specialAttack4GateFrame); each part has its own cost.
	SpecialAttack4PowerCost1 = 10.0 // first press: one HUD segment
	SpecialAttack4PowerCost2 = 20.0 // second press, to finish the attack: two HUD segments

	AttackPowerCost = 10.0 // power spent per normal (J / I) attack: one HUD segment
	JumpPowerCost   = 5.0  // power spent per jump: half a HUD segment

	RunPowerDrainPerSecond = 5.0  // power drained per second while running (half a HUD segment)
	RunRecoverPower        = 20.0 // after power runs out, the player kneels and can't run until it refills to this (2 HUD segments)

	DamageFlashSeconds = 1.0 // the player blinks red this long after taking damage
	damageBlinkSeconds = 0.1 // length of each red / normal half of a blink
)

// Combat holds the player's health and special-power pools. Health is an integer
// so damage below one point is never counted; power is continuous so it can
// regenerate smoothly and drive a HUD bar.
type Combat struct {
	Health    int
	MaxHealth int
	Power     float64
	MaxPower  float64
}

// newCombat returns a full health/power pool.
func newCombat() Combat {
	return Combat{
		Health:    MaxPlayerHealth,
		MaxHealth: MaxPlayerHealth,
		Power:     MaxPlayerPower,
		MaxPower:  MaxPlayerPower,
	}
}

// ApplyDamage subtracts dmg from the player's health, floored at zero, and starts
// the red damage flash; a heavy hit also knocks the player backward (see hurt).
// Death is handled by updateCombat on the next tick when health reaches zero, so
// external callers can deal damage without duplicating that logic. Every source
// of damage goes through here or ApplyDamageFrom.
func (p *PlayerRuntime) ApplyDamage(dmg int) {
	p.hurt(dmg, p.backward())
}

// trySpendPower deducts cost and reports true when the player had enough power;
// otherwise it leaves power untouched and reports false.
func (p *PlayerRuntime) trySpendPower(cost float64) bool {
	if p.Combat.Power < cost {
		return false
	}
	p.Combat.Power -= cost
	return true
}

// regeneratePower refills special power over time, up to the maximum.
func (p *PlayerRuntime) regeneratePower() {
	p.Combat.Power += PowerRegenPerSecond * deltaTime()
	if p.Combat.Power > p.Combat.MaxPower {
		p.Combat.Power = p.Combat.MaxPower
	}
}

// updateStamina drains power while the player runs and regenerates it otherwise.
// Power running out (by running or spending) means exhaustion, which blocks
// running and puts the player in the exhausted state until power has refilled
// to RunRecoverPower, so the player can't flicker between run and walk at zero.
func (p *PlayerRuntime) updateStamina() {
	if p.Combat.Power <= 0 {
		p.exhausted = true
	}
	if p.running {
		p.Combat.Power -= RunPowerDrainPerSecond * deltaTime()
		if p.Combat.Power <= 0 {
			p.Combat.Power = 0
			p.exhausted = true
		}
		return
	}
	p.regeneratePower()
	if p.exhausted && p.Combat.Power >= RunRecoverPower {
		p.exhausted = false
	}
}

// updateCombat applies checkpoint pickups, environmental damage and fall damage
// for this tick, then starts a death animation if health has run out.
func (p *PlayerRuntime) updateCombat(qt *DynamicQuadtree, checkpoints []Checkpoint, onGround bool) {
	p.damageFlash = math.Max(p.damageFlash-deltaTime(), 0)
	p.updateCheckpoints(checkpoints)

	inWater, inLarva := p.scanHazards(qt)
	p.applyFallDamage(onGround, inWater)

	// Water drains health continuously while submerged. Accumulate the fractional
	// per-tick amount and subtract whole points so integer health ticks down.
	if inWater {
		p.waterDamageAccum += WaterDamagePerSecond * deltaTime()
		if whole := int(p.waterDamageAccum); whole > 0 {
			p.ApplyDamage(whole)
			p.waterDamageAccum -= float64(whole)
		}
	} else {
		p.waterDamageAccum = 0
	}

	if inLarva {
		if !p.wasInLarva {
			p.ApplyDamage(LarvaDamage)
		}
		p.wasInLarva = true
	} else {
		p.wasInLarva = false
	}

	if p.Combat.Health <= 0 && !p.State.IsDying() {
		p.startDeath()
	}
}

// updateCheckpoints promotes the last checkpoint the player is overlapping to the
// active respawn point.
func (p *PlayerRuntime) updateCheckpoints(checkpoints []Checkpoint) {
	body := p.GetBounds()
	for _, cp := range checkpoints {
		if body.Intersects(cp.Bounds()) {
			p.RespawnPos = cp.Pos
		}
	}
}

// scanHazards reports whether the player is touching water or larva tiles. The
// foot probe catches standing on a hazard; the body catches walking into one.
func (p *PlayerRuntime) scanHazards(qt *DynamicQuadtree) (water, larva bool) {
	if qt == nil {
		return false, false
	}
	body := p.GetBounds()
	foot := p.GetFootSensor()
	scan := AABB{X: p.Pos.X, Y: p.Pos.Y, Width: bodyWidth, Height: bodyHeight + footProbeHeight}

	for _, obj := range qt.Retrieve(scan) {
		plat, ok := obj.(*Platform)
		if !ok {
			continue
		}
		b := plat.GetBounds()
		if !body.Intersects(b) && !foot.Intersects(b) {
			continue
		}
		switch plat.TileInfo.TileType {
		case Water:
			water = true
		case Larva:
			larva = true
		}
	}
	return water, larva
}

// applyFallDamage tracks the fall apex while airborne and, on landing, subtracts
// damage based on the drop measured in tiles. Falls of up to FallSafeTiles do no
// damage; each tile beyond that costs one point. Landing in water is cushioned,
// so it deals no fall damage.
func (p *PlayerRuntime) applyFallDamage(onGround, inWater bool) {
	if onGround {
		if !p.groundedPrev && !inWater {
			tiles := int((p.Pos.Y - p.fallPeakY) / LevelTileHeight)
			if dmg := tiles - FallSafeTiles; dmg > 0 {
				p.ApplyDamage(dmg)
			}
		}
		p.fallPeakY = p.Pos.Y // reset the baseline while grounded
	} else if p.Pos.Y < p.fallPeakY {
		p.fallPeakY = p.Pos.Y // track the highest point reached
	}
	p.groundedPrev = onGround
}

// respawn returns the player to the active checkpoint with full pools.
func (p *PlayerRuntime) respawn() {
	p.Pos = p.RespawnPos
	p.Combat.Health = p.Combat.MaxHealth
	p.Combat.Power = p.Combat.MaxPower
	p.Physics.VelX = 0
	p.Physics.VelY = 0
	p.State.SetPlayerState(int(PlayerStateIdle))
	p.fallPeakY = p.RespawnPos.Y
	p.groundedPrev = false
	p.wasInLarva = false
	p.waterDamageAccum = 0
	p.inWater = false
	p.running = false
	p.exhausted = false
	p.damageFlash = 0
	p.knockTicks = 0
}

// flashingRed reports whether the player is drawn red this frame: the damage
// flash alternates red and normal every damageBlinkSeconds.
func (p *PlayerRuntime) flashingRed() bool {
	if p.damageFlash <= 0 {
		return false
	}
	elapsed := DamageFlashSeconds - p.damageFlash
	return int(elapsed/damageBlinkSeconds)%2 == 0
}
