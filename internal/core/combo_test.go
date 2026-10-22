package core

import "testing"

// comboPlayer returns a player mid-way through a K/L swing on the given frame,
// with its animations loaded so UpdateAnimation can advance it.
func comboPlayer(state PlayerStateType, frame int) *PlayerRuntime {
	p := attackingPlayer(Position{X: 100, Y: 100}, state, frame)
	p.Animations = InitPlayerAnimations()
	return p
}

// playFrame advances the player's animation by exactly one frame.
func playFrame(p *PlayerRuntime) {
	anim := p.Animations[p.State.GetPlayerState()]
	anim.FrameTimer = 1.0/anim.AnimationSpeed - deltaTime()/2
	p.UpdateAnimation()
}

func TestComboStopsAtGateWithoutSecondPress(t *testing.T) {
	for _, st := range []PlayerStateType{PlayerStateSpecialAttack1, PlayerStateSpecialAttack2} {
		p := comboPlayer(st, comboGateFrame-1)
		playFrame(p)
		if !p.State.IsIdle() || p.CurrAnimFrame != 0 {
			t.Errorf("state %d: after gate got state=%d frame=%d, want idle frame 0", st, p.State.CurrentState, p.CurrAnimFrame)
		}
	}
}

func TestComboContinuesAfterSecondPress(t *testing.T) {
	cases := []struct {
		st PlayerStateType
		in InputState
	}{
		{PlayerStateSpecialAttack1, InputState{Skills: Skills{SpecialAttack1: true}}},
		{PlayerStateSpecialAttack2, InputState{Skills: Skills{SpecialAttack2: true}}},
	}
	for _, c := range cases {
		p := comboPlayer(c.st, 2)
		p.queueCombo(&c.in)
		p.CurrAnimFrame = comboGateFrame - 1
		playFrame(p)
		if p.State.CurrentState != c.st || p.CurrAnimFrame != comboGateFrame {
			t.Errorf("state %d: got state=%d frame=%d, want to keep swinging on frame %d", c.st, p.State.CurrentState, p.CurrAnimFrame, comboGateFrame)
		}
	}
}

func TestComboIgnoresStartingPressAndOtherKey(t *testing.T) {
	// The press that starts the swing must not count as the second press.
	p := comboPlayer(PlayerStateSpecialAttack1, 0)
	p.PreviousState.CurrentState = PlayerStateIdle
	p.queueCombo(&InputState{Skills: Skills{SpecialAttack1: true}})
	if p.comboQueued {
		t.Error("starting press queued the combo")
	}

	// L doesn't continue a K swing.
	p = comboPlayer(PlayerStateSpecialAttack1, 2)
	p.queueCombo(&InputState{Skills: Skills{SpecialAttack2: true}})
	if p.comboQueued {
		t.Error("L press queued the K combo")
	}

	// A press after the gate is too late.
	p = comboPlayer(PlayerStateSpecialAttack1, comboGateFrame)
	p.queueCombo(&InputState{Skills: Skills{SpecialAttack1: true}})
	if p.comboQueued {
		t.Error("press after the gate queued the combo")
	}
}

func TestComboQueueClearedOnNewSwing(t *testing.T) {
	p := comboPlayer(PlayerStateSpecialAttack1, 0)
	p.comboQueued = true
	p.PreviousState.CurrentState = PlayerStateIdle
	p.queueCombo(&InputState{})
	if p.comboQueued {
		t.Error("combo stayed queued into a new swing")
	}
}

// O is split like K/L, but later: the first press plays up to sprite 9.
func TestOStopsAtItsGateWithoutSecondPress(t *testing.T) {
	p := comboPlayer(PlayerStateSpecialAttack4, specialAttack4GateFrame-1)
	playFrame(p)
	if !p.State.IsIdle() || p.CurrAnimFrame != 0 {
		t.Errorf("after O gate got state=%d frame=%d, want idle frame 0", p.State.CurrentState, p.CurrAnimFrame)
	}

	// K's earlier gate must not cut O.
	p = comboPlayer(PlayerStateSpecialAttack4, comboGateFrame-1)
	playFrame(p)
	if !p.State.IsSpecialAttack4() {
		t.Errorf("O was cut at the K/L gate (frame %d)", comboGateFrame)
	}
}

func TestOFirstPressCostsPartOne(t *testing.T) {
	p := staminaPlayer()
	p.updateState(&InputState{Skills: Skills{SpecialAttack4: true}}, true, true)
	if !p.State.IsSpecialAttack4() || p.Combat.Power != MaxPlayerPower-SpecialAttack4PowerCost1 {
		t.Errorf("state=%d power=%v, want O with power %v", p.State.CurrentState, p.Combat.Power, MaxPlayerPower-SpecialAttack4PowerCost1)
	}

	low := staminaPlayer()
	low.Combat.Power = SpecialAttack4PowerCost1 - 1
	low.updateState(&InputState{Skills: Skills{SpecialAttack4: true}}, true, true)
	if low.State.IsSpecialAttack4() {
		t.Error("O started without enough stamina for part 1")
	}
}

func TestOSecondPressCostsPartTwoAndContinues(t *testing.T) {
	o := &InputState{Skills: Skills{SpecialAttack4: true}}
	p := comboPlayer(PlayerStateSpecialAttack4, 3)
	p.queueCombo(o)
	p.queueCombo(o) // a later press in the same swing is not charged again
	if !p.comboQueued || p.Combat.Power != MaxPlayerPower-SpecialAttack4PowerCost2 {
		t.Fatalf("queued=%v power=%v, want queued and power %v", p.comboQueued, p.Combat.Power, MaxPlayerPower-SpecialAttack4PowerCost2)
	}
	p.CurrAnimFrame = specialAttack4GateFrame - 1
	playFrame(p)
	if !p.State.IsSpecialAttack4() || p.CurrAnimFrame != specialAttack4GateFrame {
		t.Errorf("got state=%d frame=%d, want O to go on to frame %d", p.State.CurrentState, p.CurrAnimFrame, specialAttack4GateFrame)
	}
}

func TestOSecondPressNeedsPartTwoStamina(t *testing.T) {
	p := comboPlayer(PlayerStateSpecialAttack4, 3)
	p.Combat.Power = SpecialAttack4PowerCost2 - 1
	p.queueCombo(&InputState{Skills: Skills{SpecialAttack4: true}})
	if p.comboQueued || p.Combat.Power != SpecialAttack4PowerCost2-1 {
		t.Errorf("queued=%v power=%v, want the press ignored and no power spent", p.comboQueued, p.Combat.Power)
	}
}

// playO runs the per-tick player update (state, combo queue, animation) on a
// grounded player, pressing O on the given ticks, until the O attack ends. It
// returns the last O sprite shown and the power left.
func playO(t *testing.T, pressTicks ...int) (lastFrame int, power float64) {
	t.Helper()
	p := staminaPlayer()
	p.Animations = InitPlayerAnimations()
	pressed := map[int]bool{}
	for _, tick := range pressTicks {
		pressed[tick] = true
	}
	lastFrame = -1
	for tick := 0; tick < 600; tick++ {
		in := InputState{Skills: Skills{SpecialAttack4: pressed[tick]}}
		p.PreviousState = PlayerState{CurrentState: p.State.CurrentState}
		p.updateState(&in, true, true)
		p.queueCombo(&in)
		p.UpdateAnimation()
		if p.State.IsSpecialAttack4() {
			lastFrame = p.CurrAnimFrame
		} else if lastFrame >= 0 {
			return lastFrame, p.Combat.Power
		}
	}
	t.Fatal("O attack never ended")
	return
}

func TestOPlaysInTwoPartsTickByTick(t *testing.T) {
	last, power := playO(t, 0)
	if last != specialAttack4GateFrame-1 || power != MaxPlayerPower-SpecialAttack4PowerCost1 {
		t.Errorf("one press: last sprite %d power %v, want sprite %d and power %v",
			last, power, specialAttack4GateFrame-1, MaxPlayerPower-SpecialAttack4PowerCost1)
	}

	full := InitPlayerAnimations()[SpecialAttack4].TotalFrames - 1
	last, power = playO(t, 0, 20) // second press about a third of a second in
	if want := MaxPlayerPower - SpecialAttack4PowerCost1 - SpecialAttack4PowerCost2; last != full || power != want {
		t.Errorf("two presses: last sprite %d power %v, want sprite %d and power %v", last, power, full, want)
	}
}
