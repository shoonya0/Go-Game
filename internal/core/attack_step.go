package core

// AttackStepDistance is how far an attack steps the player: three quarters of
// a tile.
const AttackStepDistance = 0.75 * LevelTileWidth

// attackStep moves the player during part of an attack: Distance px the way it
// faces (negative steps back) over Frames animation frames from FirstFrame,
// then holds it still until the attack ends. A Combo step only happens in the
// second part of a combo, after the second press.
type attackStep struct {
	FirstFrame, Frames int
	Distance           float64
	Combo              bool
}

// attackSteps holds the attacks that move the player, timed to the sprites in
// GideonGraves.png: J steps back as it slashes, I and U lunge forward into the
// swing, L's second part backs off before the disc and O's steps in before the
// thrust.
var attackSteps = map[PlayerStateType]attackStep{
	PlayerStateWeakAttack:     {1, 2, -AttackStepDistance, false},
	PlayerStateStrongAttack:   {1, 2, AttackStepDistance, false},
	PlayerStateSpecialAttack3: {2, 2, AttackStepDistance, false},
	PlayerStateSpecialAttack2: {comboGateFrame, 2, -AttackStepDistance, true},
	PlayerStateSpecialAttack4: {specialAttack4GateFrame, 3, AttackStepDistance, true},
}

// applyAttackStep drives VelX through the current attack's step, through the
// normal collision like applyKnockback. It reports whether it took over
// horizontal movement this tick.
func (player *PlayerRuntime) applyAttackStep() bool {
	step, ok := attackSteps[player.State.CurrentState]
	if !ok || (step.Combo && !player.comboQueued) {
		return false
	}
	anim := player.Animations[player.State.GetPlayerState()]
	f := player.CurrAnimFrame - step.FirstFrame
	if anim == nil || f < 0 {
		return false
	}
	player.running = false
	if f >= step.Frames {
		player.Physics.VelX = 0
		return true
	}
	dir := 1.0
	if player.FlipX {
		dir = -1
	}
	player.Physics.VelX = dir * step.Distance * anim.AnimationSpeed / float64(step.Frames)
	return true
}
