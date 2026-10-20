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
