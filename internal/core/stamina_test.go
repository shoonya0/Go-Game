package core

import (
	"math"
	"testing"
)

// staminaPlayer returns a grounded, idle player with full pools.
func staminaPlayer() *PlayerRuntime {
	p := &PlayerRuntime{Combat: newCombat(), State: PlayerState{CurrentState: PlayerStateIdle}}
	p.Physics.MaxSpeed = MaxSpeed
	p.Physics.MaxRunSpeed = MaxRunSpeed
	return p
}

// tickStamina runs movement input and the stamina update for n ticks.
func tickStamina(p *PlayerRuntime, in *InputState, n int) {
	for i := 0; i < n; i++ {
		p.applyHorizontalMovement(in)
		p.updateStamina()
	}
}

func TestRunningDrainsHalfSegmentPerSecond(t *testing.T) {
	p := staminaPlayer()
	in := &InputState{RunJustPressed: true}
	in.Direction.LeftRight = 1

	tickStamina(p, in, int(tps())) // one second

	want := MaxPlayerPower - RunPowerDrainPerSecond
	if math.Abs(p.Combat.Power-want) > 1e-6 {
		t.Fatalf("power after 1s running = %v, want %v", p.Combat.Power, want)
	}
	if got := RunPowerDrainPerSecond; got != MaxPlayerPower/hudSegments/2 {
		t.Fatalf("drain %v/s is not half a HUD segment (%v)", got, MaxPlayerPower/hudSegments/2)
	}
}

func TestRunHeldWithoutMovingRegenerates(t *testing.T) {
	p := staminaPlayer()
	p.Combat.Power = 50
	tickStamina(p, &InputState{RunJustPressed: true}, int(tps()))

	if p.running || p.Combat.Power <= 50 {
		t.Fatalf("standing with run held: running=%v power=%v, want not running and regenerating", p.running, p.Combat.Power)
	}
}

func TestExhaustionBlocksRunningUntilRecovered(t *testing.T) {
	p := staminaPlayer()
	p.Combat.Power = 1
	run := &InputState{RunJustPressed: true}
	run.Direction.LeftRight = 1

	for i := 0; i < int(tps()) && !p.exhausted; i++ { // run until power is gone
		tickStamina(p, run, 1)
	}
	if !p.exhausted || p.Combat.Power != 0 {
		t.Fatalf("after draining: exhausted=%v power=%v, want exhausted at 0", p.exhausted, p.Combat.Power)
	}

	tickStamina(p, run, 1) // still holding run: walks and regenerates
	if p.running {
		t.Fatal("exhausted player is still running")
	}

	walk := &InputState{}
	tickStamina(p, walk, int(tps()*RunRecoverPower/PowerRegenPerSecond)+1)
	if p.exhausted {
		t.Fatalf("still exhausted at power %v, want recovered at %v", p.Combat.Power, RunRecoverPower)
	}
	tickStamina(p, run, 1)
	if !p.running {
		t.Fatal("recovered player cannot run")
	}
}

func TestJumpCostsHalfSegment(t *testing.T) {
	p := staminaPlayer()
	p.Physics.JumpForce = 100
	p.applyGravityAndJump(&InputState{JumpJustPressed: true})

	if !p.State.IsJumping() || p.Combat.Power != MaxPlayerPower-JumpPowerCost {
		t.Fatalf("jump: jumping=%v power=%v, want jumping with power %v", p.State.IsJumping(), p.Combat.Power, MaxPlayerPower-JumpPowerCost)
	}
	if JumpPowerCost != MaxPlayerPower/hudSegments/2 {
		t.Fatalf("jump cost %v is not half a HUD segment", JumpPowerCost)
	}

	low := staminaPlayer()
	low.Combat.Power = JumpPowerCost - 1
	low.applyGravityAndJump(&InputState{JumpJustPressed: true})
	if low.State.IsJumping() || low.Combat.Power != JumpPowerCost-1 {
		t.Fatalf("jump without stamina: jumping=%v power=%v, want blocked and unchanged", low.State.IsJumping(), low.Combat.Power)
	}
}

func TestNormalAttacksCostOneSegment(t *testing.T) {
	if AttackPowerCost != MaxPlayerPower/hudSegments {
		t.Fatalf("attack cost %v is not one HUD segment", AttackPowerCost)
	}
	for _, tc := range []struct {
		name   string
		skills Skills
		want   PlayerStateType
	}{
		{"weak", Skills{WeakAttack: true}, PlayerStateWeakAttack},
		{"strong", Skills{StrongAttack: true}, PlayerStateStrongAttack},
	} {
		p := staminaPlayer()
		in := &InputState{Skills: tc.skills}
		p.updateState(in, true, true)
		if p.State.CurrentState != tc.want || p.Combat.Power != MaxPlayerPower-AttackPowerCost {
			t.Errorf("%s attack: state=%d power=%v, want state %d and power %v", tc.name, p.State.CurrentState, p.Combat.Power, tc.want, MaxPlayerPower-AttackPowerCost)
		}
		p.updateState(in, true, true) // still mid-attack with the key held: no second charge
		if p.Combat.Power != MaxPlayerPower-AttackPowerCost {
			t.Errorf("%s attack charged twice: power=%v", tc.name, p.Combat.Power)
		}

		low := staminaPlayer()
		low.Combat.Power = AttackPowerCost - 1
		low.updateState(in, true, true)
		if low.State.CurrentState == tc.want {
			t.Errorf("%s attack started without enough stamina", tc.name)
		}
	}
}

func TestAirAttackCostsOneSegment(t *testing.T) {
	p := staminaPlayer()
	p.State.CurrentState = PlayerStateFalling
	p.updateAirState(&InputState{Skills: Skills{WeakAttack: true}})
	if !p.State.IsWeakAttackInAir() || p.Combat.Power != MaxPlayerPower-AttackPowerCost {
		t.Fatalf("air attack: state=%d power=%v", p.State.CurrentState, p.Combat.Power)
	}
}
