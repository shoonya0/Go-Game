package core

import (
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

// The U special attack sends a ground wave along the floor in the direction the
// player faces. It plays the frames after the player's own U frames on the same
// sprite-sheet row, while the player stays on screen and recovers.
const (
	GroundWaveDamage = 30    // per enemy, once per wave
	GroundWaveSpeed  = 320.0 // px per second

	groundWaveFPS        = 10   // frames per second of the wave clip
	groundWaveSpawnFrame = 4    // U frame on which the sword strikes the ground
	groundWaveWidth      = 60.0 // damage box, standing on the ground
	groundWaveHeight     = 40.0
	groundWaveProbe      = 4.0 // depth of the ground check under the leading edge
	groundWaveFrames     = 15  // wave frames after the U frames on the row
)

// GroundWave is one travelling wave. X is its centre, Bottom the ground line it
// runs on.
type GroundWave struct {
	X, Bottom float64
	dir       float64 // +1 right, -1 left
	clip      *Animation
	frame     int
	timer     float64
	struck    []*EnemyRuntime // enemies this wave already hit
	done      bool
}

// newGroundWave places a wave just in front of the player's feet, facing the way
// the player faces. It reports false when there is no ground there to run on.
func (player *PlayerRuntime) newGroundWave(qt *DynamicQuadtree) (*GroundWave, bool) {
	u := player.Animations[SpecialAttack3]
	dir := 1.0
	if player.FlipX {
		dir = -1
	}
	body := player.GetBounds()
	w := &GroundWave{
		X:      centerX(body) + dir*(bodyWidth+groundWaveWidth)/2,
		Bottom: body.Y + body.Height,
		dir:    dir,
		clip: &Animation{
			CurrentState:         PlayerStateSpecialAttack3,
			SpriteSheetYPosition: u.SpriteSheetYPosition,
			AnimStartFrame:       u.AnimStartFrame + u.TotalFrames,
			TotalFrames:          groundWaveFrames,
			FrameWidth:           u.FrameWidth,
			FrameHeight:          u.FrameHeight,
			AnimationSpeed:       groundWaveFPS,
		},
	}
	w.frame = w.clip.AnimStartFrame
	return w, w.hasGround(qt) && !w.blocked(qt)
}

// Box is the wave's damage box.
func (w *GroundWave) Box() AABB {
	return AABB{X: w.X - groundWaveWidth/2, Y: w.Bottom - groundWaveHeight, Width: groundWaveWidth, Height: groundWaveHeight}
}

// hasGround reports whether solid ground lies under the wave's leading edge.
func (w *GroundWave) hasGround(qt *DynamicQuadtree) bool {
	lead := w.X + w.dir*groundWaveWidth/2
	return solidAt(qt, AABB{X: lead - 1, Y: w.Bottom, Width: 2, Height: groundWaveProbe})
}

// blocked reports whether a wall stands in the wave's box (lifted one pixel off
// the floor so the ground itself doesn't count).
func (w *GroundWave) blocked(qt *DynamicQuadtree) bool {
	b := w.Box()
	b.Height--
	return solidAt(qt, b)
}

// update advances the clip, moves the wave and damages enemies in its box. The
// wave ends when its clip finishes, the ground runs out or a wall stops it.
func (w *GroundWave) update(qt *DynamicQuadtree, enemies []*EnemyRuntime) {
	w.timer += deltaTime()
	for w.timer >= 1.0/w.clip.AnimationSpeed {
		w.timer -= 1.0 / w.clip.AnimationSpeed
		w.frame++
		if w.frame >= w.clip.AnimStartFrame+w.clip.TotalFrames {
			w.done = true
			return
		}
	}

	w.X += w.dir * GroundWaveSpeed * deltaTime()
	if !w.hasGround(qt) || w.blocked(qt) {
		w.done = true
		return
	}

	box := w.Box()
	for _, e := range enemies {
		if e.State.IsDead() || slices.Contains(w.struck, e) || !box.Intersects(e.GetBounds()) {
			continue
		}
		e.TakeDamage(GroundWaveDamage)
		w.struck = append(w.struck, e)
	}
}

// solidAt reports whether box overlaps any solid (non-water) platform.
func solidAt(qt *DynamicQuadtree, box AABB) bool {
	if qt == nil {
		return false
	}
	for _, obj := range qt.Retrieve(box) {
		if plat, ok := obj.(*Platform); ok && plat.TileInfo.TileType != Water && box.Intersects(plat.GetBounds()) {
			return true
		}
	}
	return false
}

// UpdateGroundWaves casts a wave when the U swing reaches the ground (once per
// swing) and moves every live wave. Call it after UpdateAnimation.
func (player *PlayerRuntime) UpdateGroundWaves(qt *DynamicQuadtree, enemies []*EnemyRuntime) {
	if player.State.CurrentState != player.PreviousState.CurrentState {
		player.waveCast = false // a new swing (or none) started this tick
	}
	if player.State.IsSpecialAttack3() && !player.waveCast && player.CurrAnimFrame >= groundWaveSpawnFrame {
		player.waveCast = true
		if w, ok := player.newGroundWave(qt); ok {
			player.waves = append(player.waves, w)
		}
	}

	live := player.waves[:0]
	for _, w := range player.waves {
		w.update(qt, enemies)
		if !w.done {
			live = append(live, w)
		}
	}
	player.waves = live
}

// DrawGroundWaves draws every live wave, mirrored when it runs left.
func (player *PlayerRuntime) DrawGroundWaves(screen *ebiten.Image) {
	for _, w := range player.waves {
		drawSpriteFrame(screen, player.img, w.clip, w.frame, w.Box(), w.dir < 0,
			player.Scale, player.Camera.Pos, frameOriginHeight, ebiten.ColorScale{})
	}
}
