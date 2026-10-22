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
