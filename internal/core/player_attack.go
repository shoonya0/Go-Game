package core

import "slices"

// Player attack damage per hit, in enemy health bars (EnemyHealthPerBar HP). A
// succubus has 2 bars.
const (
	halfBar = EnemyHealthPerBar / 2

	WeakAttackDamage      = halfBar     // J, on the ground or in the air: 0.5 bar
	StrongAttackDamage    = halfBar     // I on the ground: 0.5 bar
	StrongAirAttackDamage = 3 * halfBar // I in the air: 1.5 bars
	SpecialAttack1Damage1 = halfBar     // K, first hit: 0.5 bar
	SpecialAttack1Damage2 = 2 * halfBar // K, second hit: 1 bar (1.5 in all)
	SpecialAttack2Damage1 = 2 * halfBar // L, first hit: 1 bar
	SpecialAttack2Damage2 = 4 * halfBar // L, second hit: 2 bars (3 in all)
	SpecialAttack4Damage1 = 2 * halfBar // O, burst (part 1): 1 bar
	SpecialAttack4Damage2 = 4 * halfBar // O, sword thrust (part 2): 2 bars (3 in all)
)

// strikeHit is one hit of an attack: it deals Damage on animation frames
// FirstFrame..LastFrame, at most once per enemy. Reach, when non-zero, replaces
// the attack's Reach for this hit (O's burst is shorter than its thrust).
type strikeHit struct {
	FirstFrame, LastFrame int
	Damage                int
	Reach                 float64
}

// playerStrike is one attack's damage box: it reaches Reach px forward from the
// body's centre and rises Height px from the feet. Most attacks hit once; K, L
// and O hit twice, once per half of the combo.
type playerStrike struct {
	Reach, Height float64
	Hits          []strikeHit
}

// playerStrikes holds the damage box of every attack that has one. Reach, height
// and frames are measured from the sword and slashes in GideonGraves.png. U hurts
// through its ground wave instead.
var playerStrikes = map[PlayerStateType]playerStrike{
	PlayerStateWeakAttack:        {Reach: 78, Height: bodyHeight, Hits: []strikeHit{{2, 4, WeakAttackDamage, 0}}},
	PlayerStateStrongAttack:      {Reach: 72, Height: 96, Hits: []strikeHit{{2, 5, StrongAttackDamage, 0}}},
	PlayerStateWeakAttackInAir:   {Reach: 100, Height: 102, Hits: []strikeHit{{2, 4, WeakAttackDamage, 0}}},
	PlayerStateStrongAttackInAir: {Reach: 77, Height: 120, Hits: []strikeHit{{2, 7, StrongAirAttackDamage, 0}}},
	PlayerStateSpecialAttack1: {Reach: 70, Height: 100, Hits: []strikeHit{
		{3, 4, SpecialAttack1Damage1, 0}, // first crescent, before the combo gate
		{5, 7, SpecialAttack1Damage2, 0}, // second crescent
	}},
	PlayerStateSpecialAttack2: {Reach: 95, Height: 100, Hits: []strikeHit{
		{3, 4, SpecialAttack2Damage1, 0}, // first blue wave, before the combo gate
		{8, 9, SpecialAttack2Damage2, 0}, // the blue disc
	}},
	PlayerStateSpecialAttack4: {Reach: 95, Height: bodyHeight, Hits: []strikeHit{
		{1, 4, SpecialAttack4Damage1, 70},  // energy burst from the fist, before the combo gate
		{12, 14, SpecialAttack4Damage2, 0}, // sword thrust after the step forward
	}},
}

// hitOn returns the index of the hit active on animation frame f, if any.
func (st playerStrike) hitOn(f int) (int, bool) {
	for i, h := range st.Hits {
		if f >= h.FirstFrame && f <= h.LastFrame {
			return i, true
		}
	}
	return 0, false
}

// currentStrike returns the damage box spec of the attack being performed, if any.
func (player *PlayerRuntime) currentStrike() (playerStrike, bool) {
	st, ok := playerStrikes[player.State.CurrentState]
	return st, ok
}

// reachOn returns how far st reaches on animation frame f: the active hit's own
// Reach if it has one, otherwise the attack's.
func (st playerStrike) reachOn(f int) float64 {
	if i, ok := st.hitOn(f); ok && st.Hits[i].Reach > 0 {
		return st.Hits[i].Reach
	}
	return st.Reach
}

// strikeBox places st in the world in front of the player, on the side it faces,
// sized for the current animation frame.
func (player *PlayerRuntime) strikeBox(st playerStrike) AABB {
	b := player.GetBounds()
	reach := st.reachOn(player.CurrAnimFrame)
	x := centerX(b)
	if player.FlipX {
		x -= reach
	}
	return AABB{X: x, Y: b.Y + b.Height - st.Height, Width: reach, Height: st.Height}
}

// AttackBox returns the damage box of the current attack and whether it is
// striking this frame. ok is false when the player isn't attacking.
func (player *PlayerRuntime) AttackBox() (box AABB, striking, ok bool) {
	st, ok := player.currentStrike()
	if !ok {
		return AABB{}, false, false
	}
	_, striking = st.hitOn(player.CurrAnimFrame)
	return player.strikeBox(st), striking, true
}

// HitEnemies damages every live enemy inside the player's damage box while a hit
// is active, each at most once per hit. Call it after UpdateAnimation so the
// frame matches the state.
func (player *PlayerRuntime) HitEnemies(enemies []*EnemyRuntime) {
	if player.State.CurrentState != player.PreviousState.CurrentState {
		player.struck = player.struck[:0] // a new swing (or none) started this tick
		player.strikeIdx = 0
	}
	st, ok := player.currentStrike()
	if !ok {
		return
	}
	i, active := st.hitOn(player.CurrAnimFrame)
	if !active {
		return
	}
	if i != player.strikeIdx {
		player.struck = player.struck[:0] // the swing's next hit can land again
		player.strikeIdx = i
	}
	box := player.strikeBox(st)
	for _, e := range enemies {
		if e.State.IsDead() || slices.Contains(player.struck, e) || !box.Intersects(e.GetBounds()) {
			continue
		}
		e.TakeDamage(st.Hits[i].Damage)
		player.struck = append(player.struck, e)
	}
}

// Hitboxes reports the player's body, its damage box while attacking and every
// live ground wave to the debug overlay (HitboxProvider).
func (player *PlayerRuntime) Hitboxes() []Hitbox {
	boxes := []Hitbox{{Box: player.GetBounds(), Kind: HitboxBody}}
	if box, striking, ok := player.AttackBox(); ok {
		boxes = append(boxes, Hitbox{Box: box, Kind: HitboxAttack, Active: striking})
	}
	for _, w := range player.waves {
		boxes = append(boxes, Hitbox{Box: w.Box(), Kind: HitboxAttack, Active: true})
	}
	return boxes
}

// comboGateFrame splits the K and L swings in two. A press plays frames
// 0..comboGateFrame-1; the swing only goes on past that if the same key is
// pressed again before then, otherwise the player returns to idle.
const comboGateFrame = 5

// specialAttack4GateFrame splits O the same way, after its first 9 sprites. The
// second press also pays for part 2 (SpecialAttack4PowerCost2).
const specialAttack4GateFrame = 9

// O's second part opens with a step forward: while the wind-up sprites after
// the gate play (frames 9-11), the player moves SpecialAttack4StepDistance px
// the way it faces, then thrusts.
const (
	SpecialAttack4StepDistance = 0.75 * LevelTileWidth
	specialAttack4StepFrames   = 3
)

// applyComboStep drives VelX during O's second part, through the normal
// collision, like applyKnockback: the step forward, then standing still for the
// thrust. It reports whether it took over horizontal movement this tick.
func (player *PlayerRuntime) applyComboStep() bool {
	if !player.State.IsSpecialAttack4() || !player.comboQueued {
		return false
	}
	anim := player.Animations[player.State.GetPlayerState()]
	f := player.CurrAnimFrame - specialAttack4GateFrame
	if anim == nil || f < 0 {
		return false
	}
	player.running = false
	if f >= specialAttack4StepFrames {
		player.Physics.VelX = 0
		return true
	}
	dir := 1.0
	if player.FlipX {
		dir = -1
	}
	player.Physics.VelX = dir * SpecialAttack4StepDistance * anim.AnimationSpeed / specialAttack4StepFrames
	return true
}

// comboGate returns the frame where the current attack stops unless its key was
// pressed again, and false for attacks that are not combos.
func (player *PlayerRuntime) comboGate() (int, bool) {
	switch s := &player.State; {
	case s.IsSpecialAttack1(), s.IsSpecialAttack2():
		return comboGateFrame, true
	case s.IsSpecialAttack4():
		return specialAttack4GateFrame, true
	}
	return 0, false
}

// comboPressed reports whether this tick's input presses the key of the combo
// attack the player is performing (K for special 1, L for special 2, O for
// special 4).
func (player *PlayerRuntime) comboPressed(in *InputState) bool {
	switch {
	case player.State.IsSpecialAttack1():
		return in.Skills.SpecialAttack1
	case player.State.IsSpecialAttack2():
		return in.Skills.SpecialAttack2
	case player.State.IsSpecialAttack4():
		return in.Skills.SpecialAttack4
	}
	return false
}

// queueCombo remembers a second K/L/O press made before the combo gate. The
// press that starts the swing doesn't count, and a new swing (or none) clears
// it. O's second press only counts if the player can pay for part 2.
func (player *PlayerRuntime) queueCombo(in *InputState) {
	if player.State.CurrentState != player.PreviousState.CurrentState {
		player.comboQueued = false
		return
	}
	gate, ok := player.comboGate()
	if !ok || player.comboQueued || player.CurrAnimFrame >= gate || !player.comboPressed(in) {
		return
	}
	if player.State.IsSpecialAttack4() && !player.trySpendPower(SpecialAttack4PowerCost2) {
		return
	}
	player.comboQueued = true
}

// cutCombo ends a K/L/O attack at its combo gate when no second press was
// queued, returning the player to idle. It reports whether it did.
func (player *PlayerRuntime) cutCombo() bool {
	gate, ok := player.comboGate()
	if !ok || player.comboQueued || player.CurrAnimFrame != gate {
		return false
	}
	player.State.SetPlayerState(int(PlayerStateIdle))
	player.CurrAnimFrame = 0
	return true
}
