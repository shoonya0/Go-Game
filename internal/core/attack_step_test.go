package core

import "testing"

// playAttack runs the per-tick player update on a grounded player with no level
// to collide with, pressing the given skills on the given ticks, until attack
// st ends. It returns how far the player moved.
func playAttack(t *testing.T, st PlayerStateType, flipX bool, presses map[int]Skills) float64 {
	t.Helper()
	p := staminaPlayer()
	p.Animations = InitPlayerAnimations()
	p.FlipX = flipX
	start := p.Pos.X
	started := false
	for tick := 0; tick < 600; tick++ {
		in := InputState{Skills: presses[tick]}
		p.PreviousState = PlayerState{CurrentState: p.State.CurrentState}
		p.applyHorizontalMovement(&in)
		p.resolveHorizontal(nil)
		p.updateState(&in, true, true)
		p.queueCombo(&in)
		p.UpdateAnimation()
		if p.State.CurrentState == st {
			started = true
		} else if started {
			return p.Pos.X - start
		}
	}
	t.Fatalf("attack %d never ended", st)
	return 0
}

func TestAttackSteps(t *testing.T) {
	j, i, u := Skills{WeakAttack: true}, Skills{StrongAttack: true}, Skills{SpecialAttack3: true}
	k, l, o := Skills{SpecialAttack1: true}, Skills{SpecialAttack2: true}, Skills{SpecialAttack4: true}
	for _, c := range []struct {
		name    string
		st      PlayerStateType
		presses map[int]Skills
		want    float64 // px the way the player faces
	}{
		{"J", PlayerStateWeakAttack, map[int]Skills{0: j}, -AttackStepDistance},
		{"I", PlayerStateStrongAttack, map[int]Skills{0: i}, AttackStepDistance},
		{"U", PlayerStateSpecialAttack3, map[int]Skills{0: u}, AttackStepDistance},
		{"L one press", PlayerStateSpecialAttack2, map[int]Skills{0: l}, 0},
		{"L two presses", PlayerStateSpecialAttack2, map[int]Skills{0: l, 20: l}, ShortAttackStepDistance},
		{"O one press", PlayerStateSpecialAttack4, map[int]Skills{0: o}, 0},
		{"O two presses", PlayerStateSpecialAttack4, map[int]Skills{0: o, 20: o}, AttackStepDistance},
		{"K one press", PlayerStateSpecialAttack1, map[int]Skills{0: k}, 0},
		{"K two presses", PlayerStateSpecialAttack1, map[int]Skills{0: k, 10: k}, ShortAttackStepDistance},
	} {
		for _, flip := range []bool{false, true} {
			want := c.want
			if flip {
				want = -want
			}
			const slack = 1.0 // float rounding
			if got := playAttack(t, c.st, flip, c.presses); got < want-slack || got > want+slack {
				t.Errorf("%s facing left=%v: moved %v px, want %v", c.name, flip, got, want)
			}
		}
	}
}
