package core

import "github.com/hajimehoshi/ebiten/v2"

// UpdateEnemyAnimation advances the active animation by one tick. It restarts the
// clip when the state changed this tick, loops looping clips, and resolves the
// follow-up state when a one-shot clip finishes: attack/hurt return to patrol,
// while death holds on its final frame and marks the enemy removable. It mirrors
// PlayerRuntime.UpdateAnimation, except that a clip may Hold single frames longer.
func (e *EnemyRuntime) UpdateEnemyAnimation() {
	anim := e.currentAnim()
	if anim == nil {
		return
	}
	clip := e.currentClip()

	// Restart the clip whenever the state changed this tick.
	if e.PrevState != e.State {
		e.CurrAnimFrame = anim.AnimStartFrame
		anim.FrameTimer = 0
		e.clipTime = 0
	}

	if anim.AnimationSpeed <= 0 {
		return // guard against divide-by-zero on a misconfigured blueprint
	}
	timePerFrame := 1.0 / anim.AnimationSpeed
	anim.FrameTimer += deltaTime()
	e.clipTime += deltaTime()

	for {
		frameTime := timePerFrame * clip.hold(e.CurrAnimFrame-anim.AnimStartFrame)
		if anim.FrameTimer < frameTime {
			break
		}
		anim.FrameTimer -= frameTime
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
		e.CurrAnimFrame = e.animFor(EnemyPatrol).AnimStartFrame
		e.clipTime = 0
	}
}

// DrawEnemyAnimation draws the current frame through the shared sprite path,
// moved by the clip's draw-time motion and positioned relative to the supplied
// camera (the player owns the only camera).
func (e *EnemyRuntime) DrawEnemyAnimation(screen *ebiten.Image, cam Camera) {
	anim := e.currentAnim()
	if anim == nil {
		return
	}
	clip := e.currentClip()
	drawSpriteFramePosed(screen, enemySheet(clip.Sheet), anim, e.CurrAnimFrame, e.GetBounds(),
		e.FlipX, e.Scale, cam.Pos, anim.FrameHeight, ebiten.ColorScale{}, clip.Motion.pose(e.clipTime, e.FlipX))
}
