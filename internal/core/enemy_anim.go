package core

import "github.com/hajimehoshi/ebiten/v2"

// UpdateEnemyAnimation advances the active animation by one tick. It restarts the
// clip when the state changed this tick, loops looping clips, and resolves the
// follow-up state when a one-shot clip finishes: attack/hurt return to patrol,
// while death holds on its final frame and marks the enemy removable. It mirrors
// PlayerRuntime.UpdateAnimation so both actors animate identically.
func (e *EnemyRuntime) UpdateEnemyAnimation() {
	anim := e.Animations[int(e.State)]
	if anim == nil {
		return
	}

	// Restart the clip whenever the state changed this tick.
	if e.PrevState != e.State {
		e.CurrAnimFrame = anim.AnimStartFrame
		anim.FrameTimer = 0
	}

	if anim.AnimationSpeed <= 0 {
		return // guard against divide-by-zero on a misconfigured blueprint
	}
	timePerFrame := 1.0 / anim.AnimationSpeed
	anim.FrameTimer += deltaTime()

	for anim.FrameTimer >= timePerFrame {
		anim.FrameTimer -= timePerFrame
		e.CurrAnimFrame++

		if e.CurrAnimFrame < anim.AnimStartFrame+anim.TotalFrames {
			continue
		}

		if anim.Looping {
			e.CurrAnimFrame = anim.AnimStartFrame
			continue
		}

		// A one-shot clip finished: resolve the follow-up.
		if e.State.IsDead() {
			e.CurrAnimFrame = anim.AnimStartFrame + anim.TotalFrames - 1 // hold last frame
			e.deathDone = true
			anim.FrameTimer = 0
			return
		}
		e.setState(EnemyPatrol)
		e.CurrAnimFrame = e.Animations[int(EnemyPatrol)].AnimStartFrame
	}
}

// DrawEnemyAnimation draws the current frame through the shared sprite path,
// positioned relative to the supplied camera (the player owns the only camera).
func (e *EnemyRuntime) DrawEnemyAnimation(screen *ebiten.Image, cam Camera) {
	anim := e.Animations[int(e.State)]
	if anim == nil {
		return
	}
	drawSpriteFrame(screen, e.img, anim, e.CurrAnimFrame, e.GetBounds(),
		e.FlipX, e.Scale, cam.Pos, e.bp.OriginHeight)
}
