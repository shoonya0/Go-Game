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
)

// strikeHit is one hit of an attack: it deals Damage on animation frames
// FirstFrame..LastFrame, at most once per enemy.
type strikeHit struct {
	FirstFrame, LastFrame int
	Damage                int
}

// playerStrike is one attack's damage box: it reaches Reach px forward from the
// body's centre and rises Height px from the feet. Most attacks hit once; K and
// L hit twice, once per half of the combo.
type playerStrike struct {
	Reach, Height float64
	Hits          []strikeHit
}

// playerStrikes holds the damage box of every attack that has one. Reach, height
// and frames are measured from the sword and slashes in GideonGraves.png. U hurts
// through its ground wave instead; O has no box yet.
var playerStrikes = map[PlayerStateType]playerStrike{
	PlayerStateWeakAttack:        {Reach: 78, Height: bodyHeight, Hits: []strikeHit{{2, 4, WeakAttackDamage}}},
	PlayerStateStrongAttack:      {Reach: 72, Height: 96, Hits: []strikeHit{{2, 5, StrongAttackDamage}}},
	PlayerStateWeakAttackInAir:   {Reach: 100, Height: 102, Hits: []strikeHit{{2, 4, WeakAttackDamage}}},
	PlayerStateStrongAttackInAir: {Reach: 77, Height: 120, Hits: []strikeHit{{2, 7, StrongAirAttackDamage}}},
	PlayerStateSpecialAttack1: {Reach: 70, Height: 100, Hits: []strikeHit{
		{3, 4, SpecialAttack1Damage1}, // first crescent, before the combo gate
		{5, 7, SpecialAttack1Damage2}, // second crescent
	}},
	PlayerStateSpecialAttack2: {Reach: 95, Height: 100, Hits: []strikeHit{
		{3, 4, SpecialAttack2Damage1}, // first blue wave, before the combo gate
		{8, 9, SpecialAttack2Damage2}, // the blue disc
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

// strikeBox places st in the world in front of the player, on the side it faces.
func (player *PlayerRuntime) strikeBox(st playerStrike) AABB {
	b := player.GetBounds()
	x := centerX(b)
	if player.FlipX {
		x -= st.Reach
	}
	return AABB{X: x, Y: b.Y + b.Height - st.Height, Width: st.Reach, Height: st.Height}
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

// comboPressed reports whether this tick's input presses the key of the combo
// attack the player is performing (K for special 1, L for special 2).
func (player *PlayerRuntime) comboPressed(in *InputState) bool {
	switch {
	case player.State.IsSpecialAttack1():
		return in.Skills.SpecialAttack1
	case player.State.IsSpecialAttack2():
		return in.Skills.SpecialAttack2
	}
	return false
}

// queueCombo remembers a second K/L press made before the combo gate. The press
// that starts the swing doesn't count, and a new swing (or none) clears it.
func (player *PlayerRuntime) queueCombo(in *InputState) {
	if player.State.CurrentState != player.PreviousState.CurrentState {
		player.comboQueued = false
		return
	}
	if player.CurrAnimFrame < comboGateFrame && player.comboPressed(in) {
		player.comboQueued = true
	}
}

// cutCombo ends a K/L swing at the combo gate when no second press was queued,
// returning the player to idle. It reports whether it did.
func (player *PlayerRuntime) cutCombo() bool {
	s := &player.State
	inCombo := s.IsSpecialAttack1() || s.IsSpecialAttack2()
	if !inCombo || player.comboQueued || player.CurrAnimFrame != comboGateFrame {
		return false
	}
	s.SetPlayerState(int(PlayerStateIdle))
	player.CurrAnimFrame = 0
	return true
}
