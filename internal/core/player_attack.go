package core

import "slices"

// Player attack damage, per hit. A succubus (40 HP) takes 4 weak or 2 strong hits.
const (
	WeakAttackDamage   = 10 // J, on the ground or in the air
	StrongAttackDamage = 20 // I, on the ground or in the air
)

// playerStrike is one attack's damage box: it reaches Reach px forward from the
// body's centre and rises Height px from the feet. It deals Damage only on
// animation frames FirstFrame..LastFrame, while the sword is out.
type playerStrike struct {
	Reach, Height         float64
	FirstFrame, LastFrame int
	Damage                int
}

// playerStrikes holds the damage box of every attack that has one. Reach, height
// and frames are measured from the sword in GideonGraves.png. Special attacks
// (K/L/U/O) have no box yet.
var playerStrikes = map[PlayerStateType]playerStrike{
	PlayerStateWeakAttack:        {Reach: 78, Height: bodyHeight, FirstFrame: 2, LastFrame: 4, Damage: WeakAttackDamage},
	PlayerStateStrongAttack:      {Reach: 72, Height: 96, FirstFrame: 2, LastFrame: 5, Damage: StrongAttackDamage},
	PlayerStateWeakAttackInAir:   {Reach: 100, Height: 102, FirstFrame: 2, LastFrame: 4, Damage: WeakAttackDamage},
	PlayerStateStrongAttackInAir: {Reach: 77, Height: 120, FirstFrame: 2, LastFrame: 7, Damage: StrongAttackDamage},
}

// activeOn reports whether the strike deals damage on animation frame f.
func (st playerStrike) activeOn(f int) bool { return f >= st.FirstFrame && f <= st.LastFrame }

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
	return player.strikeBox(st), st.activeOn(player.CurrAnimFrame), true
}

// HitEnemies damages every live enemy inside the player's damage box while the
// sword is out, each at most once per swing. Call it after UpdateAnimation so the
// frame matches the state.
func (player *PlayerRuntime) HitEnemies(enemies []*EnemyRuntime) {
	if player.State.CurrentState != player.PreviousState.CurrentState {
		player.struck = player.struck[:0] // a new swing (or none) started this tick
	}
	st, ok := player.currentStrike()
	if !ok || !st.activeOn(player.CurrAnimFrame) {
		return
	}
	box := player.strikeBox(st)
	for _, e := range enemies {
		if e.State.IsDead() || slices.Contains(player.struck, e) || !box.Intersects(e.GetBounds()) {
			continue
		}
		e.TakeDamage(st.Damage)
		player.struck = append(player.struck, e)
	}
}

// Hitboxes reports the player's body and, while attacking, its damage box to the
// debug overlay (HitboxProvider).
func (player *PlayerRuntime) Hitboxes() []Hitbox {
	boxes := []Hitbox{{Box: player.GetBounds(), Kind: HitboxBody}}
	if box, striking, ok := player.AttackBox(); ok {
		boxes = append(boxes, Hitbox{Box: box, Kind: HitboxAttack, Active: striking})
	}
	return boxes
}
