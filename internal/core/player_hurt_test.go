package core

import (
	"math"
	"testing"
)

// tickKnockback runs horizontal movement for n ticks with no input or world.
func tickKnockback(p *PlayerRuntime, n int) {
	for i := 0; i < n; i++ {
		p.applyHorizontalMovement(&InputState{})
		p.resolveHorizontal(nil)
	}
}

func TestHeavyHitKnocksPlayerAwayFromSource(t *testing.T) {
	for _, c := range []struct {
		name string
		srcX float64
		want float64 // signed push distance
	}{
		{"source on the right", 1000, -KnockbackDistance},
		{"source on the left", 0, KnockbackDistance},
	} {
		p := staminaPlayer()
		p.Pos.X = 500
		p.ApplyDamageFrom(HeavyHitDamage, c.srcX)
		if !p.State.IsDamaged() {
			t.Fatalf("%s: state=%d, want knockdown", c.name, p.State.CurrentState)
		}
		tickKnockback(p, 2*int(tps())) // push ends, then the player stands still
		if got := p.Pos.X - 500; math.Abs(got-c.want) > 1e-6 {
			t.Errorf("%s: pushed %v px, want %v", c.name, got, c.want)
		}
	}
}

func TestHeavyHitWithoutSourcePushesBackward(t *testing.T) {
	p := staminaPlayer()
	p.FlipX = true // facing left: pushed right
	p.Pos.X = 500
	p.ApplyDamage(HeavyHitDamage)
	tickKnockback(p, 2*int(tps()))
	if got := p.Pos.X - 500; math.Abs(got-KnockbackDistance) > 1e-6 {
		t.Errorf("pushed %v px, want %v", got, KnockbackDistance)
	}
}

func TestLightHitDoesNotKnockDown(t *testing.T) {
	p := staminaPlayer()
	p.ApplyDamageFrom(HeavyHitDamage-1, 1000)
	tickKnockback(p, int(tps()))
	if p.State.IsDamaged() || p.Pos.X != 0 {
		t.Errorf("light hit: state=%d x=%v, want idle and unmoved", p.State.CurrentState, p.Pos.X)
	}
}

func TestHeavyHitIsTwoHearts(t *testing.T) {
	if HeavyHitDamage != 2*MaxPlayerHealth/hudSegments {
		t.Fatalf("HeavyHitDamage=%d, want two hearts (%d)", HeavyHitDamage, 2*MaxPlayerHealth/hudSegments)
	}
}

func TestKnockdownReturnsToIdle(t *testing.T) {
	p := staminaPlayer()
	p.Animations = InitPlayerAnimations()
	p.ApplyDamage(HeavyHitDamage)
	for i := 0; i < 3*int(tps()) && p.State.IsDamaged(); i++ {
		p.UpdateAnimation()
		p.PreviousState = p.State
	}
	if !p.State.IsIdle() {
		t.Errorf("after knockdown: state=%d, want idle", p.State.CurrentState)
	}
}

// withDeath makes startDeath pick the kneeling death (kneel) or the other one.
func withDeath(t *testing.T, kneel bool) {
	old := deathKneels
	deathKneels = func() bool { return kneel }
	t.Cleanup(func() { deathKneels = old })
}

func TestDeathPlaysAnimationThenRespawns(t *testing.T) {
	for _, c := range []struct {
		kneel bool
		want  PlayerStateType
	}{{false, PlayerStateDead}, {true, PlayerStateDeadKneel}} {
		withDeath(t, c.kneel)
		p := staminaPlayer()
		p.Animations = InitPlayerAnimations()
		p.RespawnPos = Position{X: 7, Y: 9}

		p.ApplyDamage(p.Combat.Health) // lethal, and heavy: no knockdown when dying
		p.updateCombat(nil, nil, true)
		if p.State.CurrentState != c.want {
			t.Fatalf("kneel=%v: state=%d, want %d", c.kneel, p.State.CurrentState, c.want)
		}

		p.ApplyDamage(HeavyHitDamage) // ignored while dying
		if p.State.CurrentState != c.want {
			t.Fatalf("kneel=%v: damage interrupted the death (state=%d)", c.kneel, p.State.CurrentState)
		}

		for i := 0; i < 5*int(tps()) && p.State.IsDying(); i++ {
			p.UpdateAnimation()
			p.PreviousState = p.State
		}
		if !p.State.IsIdle() || p.Pos != p.RespawnPos || p.Combat.Health != p.Combat.MaxHealth {
			t.Errorf("kneel=%v: after death state=%d pos=%v health=%d, want respawned", c.kneel, p.State.CurrentState, p.Pos, p.Combat.Health)
		}
	}
}

func TestOutOfStaminaKneelsInStepWithBar(t *testing.T) {
	p := staminaPlayer()
	p.Animations = InitPlayerAnimations()
	anim := p.Animations[Exhausted]

	p.Combat.Power = AttackPowerCost
	if !p.trySpendPower(AttackPowerCost) { // the last bar spent on an attack
		t.Fatal("couldn't spend the last bar")
	}
	p.updateStamina()
	p.updateState(&InputState{}, true, true)
	if !p.State.IsExhausted() {
		t.Fatalf("at 0 stamina: state=%d, want exhausted", p.State.CurrentState)
	}

	for _, c := range []struct {
		power float64
		frame int
	}{
		{0, 0},
		{RunRecoverPower / 2, anim.TotalFrames / 2},
		{RunRecoverPower - 0.01, anim.TotalFrames - 1},
	} {
		p.Combat.Power = c.power
		p.UpdateAnimation()
		if p.CurrAnimFrame != c.frame {
			t.Errorf("power %v: frame %d, want %d", c.power, p.CurrAnimFrame, c.frame)
		}
	}

	// Moving is blocked while kneeling.
	in := &InputState{}
	in.Direction.LeftRight = 1
	p.applyHorizontalMovement(in)
	if p.Physics.VelX != 0 {
		t.Errorf("moved while exhausted: VelX=%v", p.Physics.VelX)
	}

	p.Combat.Power = RunRecoverPower - 0.01
	p.updateStamina() // regenerates past one bar
	p.updateState(&InputState{}, true, true)
	if !p.State.IsIdle() {
		t.Errorf("after one bar refilled: state=%d, want idle", p.State.CurrentState)
	}
}
