package enemy

import (
	"image"
	"player/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
)

// image dimensions
const (
	frameWidth_minimum = 40
	frameWidth_small   = 80
	// frameWidth_medium  = 120
	frameWidth_large   = 160
	frameWidth_maximum = 200

	frameHeight_minimum = 40  // this frame is used
	frameHeight_small   = 80  // this frame is used
	frameHeight_medium  = 120 // this frame is  not used
	// frameHeight_large   = 160 // this frame is used
	// frameHeight_maximum = 200 // this frame is currently not used

	enemySpriteSheetPath = "../assets/GideonGraves.png"
)

var enemySpriteSheet *ebiten.Image

// ---------------- animation ----------------
type Animation struct {
	CurrentState         int // Use the enum instead of embedded PlayerState
	SpriteSheetYPosition int // the y position of the sprite sheet in the image
	TotalFrames          int
	AnimStartFrame       int
	FrameWidth           int
	FrameHeight          int
	FrameTimer           float64 // time in seconds that the current frame has been displayed for ie: if FrameTimer is x, then the current frame is displayed for x/AnimationSpeed seconds
	AnimationSpeed       float64 // no of frames to display per second in seconds
	Looping              bool    // true if the animation should loop
}

func InitEnemyAnimations() map[int]Animation {
	animations := make(map[int]Animation)
	animations[int(StateIdle)] = Animation{
		CurrentState:         StateIdle,
		SpriteSheetYPosition: 0,
		TotalFrames:          4,
		AnimStartFrame:       0,
		FrameWidth:           frameWidth_minimum,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       6,
		Looping:              true,
	}
	// walking
	animations[int(StatePatrolling)] = Animation{
		CurrentState:         StatePatrolling,
		SpriteSheetYPosition: 2,
		TotalFrames:          5,
		AnimStartFrame:       0,
		FrameWidth:           frameWidth_minimum,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       6,
		Looping:              true,
	}
	// running
	animations[int(StateHunting)] = Animation{
		CurrentState:         StateHunting,
		SpriteSheetYPosition: 4,
		TotalFrames:          8,
		AnimStartFrame:       0,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       10,
		Looping:              true,
	}
	animations[int(StateJumping)] = Animation{
		CurrentState:         StateJumping,
		SpriteSheetYPosition: 10,
		TotalFrames:          3,
		AnimStartFrame:       0,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       10,
		Looping:              false,
	}
	animations[int(StateFalling)] = Animation{
		CurrentState:         StateFalling,
		SpriteSheetYPosition: 10,
		TotalFrames:          1,
		AnimStartFrame:       3,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       5,
		Looping:              true,
	}
	animations[int(StateLanding)] = Animation{
		CurrentState:         StateLanding,
		SpriteSheetYPosition: 10,
		TotalFrames:          3,
		AnimStartFrame:       7,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0.5,
		AnimationSpeed:       10,
		Looping:              false,
	}

	// animations[SmugFace] = &Animation{
	// 	CurrentState:         PlayerStateSmugFace,
	// 	SpriteSheetYPosition: 8,
	// 	TotalFrames:          11,
	// 	AnimStartFrame:       0,
	// 	FrameWidth:           frameWidth_minimum,
	// 	FrameHeight:          frameHeight_small,
	// 	FrameTimer:           0,
	// 	AnimationSpeed:       7,
	// 	Looping:              false,
	// }

	animations[int(StateAttacking)] = Animation{
		CurrentState:         StateAttacking,
		SpriteSheetYPosition: 19,
		TotalFrames:          8,
		AnimStartFrame:       0,
		FrameWidth:           frameWidth_large,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       12,
		Looping:              false,
	}

	animations[int(StateDefending)] = Animation{
		CurrentState:         StateDefending,
		SpriteSheetYPosition: 21,
		TotalFrames:          8,
		AnimStartFrame:       0,
		FrameWidth:           frameWidth_large,
		FrameHeight:          frameHeight_medium,
		FrameTimer:           0,
		AnimationSpeed:       10,
		Looping:              false,
	}

	// animation yet to make
	animations[int(StateDead)] = Animation{
		CurrentState:         StateDead,
		SpriteSheetYPosition: 7,
		TotalFrames:          6,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       1,
		Looping:              false,
	}

	return animations
}

// UpdateEnemyAnimation advances the animation frame for this enemy based on
// time accumulation. Uses per-enemy FrameTimer (not the shared Animation.FrameTimer)
// because the animation map is shared across all enemies.
// Mirrors PlayerRuntime.UpdateAnimation (internal/core/animation.go:287-328).
func (e *EnemyRuntime) UpdateEnemyAnimation(animations *map[int]Animation) {
	tps := float64(ebiten.TPS())
	if tps <= 0 {
		tps = 60
	}
	dt := 1.0 / tps

	currState := e.State.Current
	anim, ok := (*animations)[currState]
	if !ok {
		anim = (*animations)[StateIdle]
	}

	// On state change: reset to animation's start frame
	if e.State.Previous != e.State.Current {
		e.CurrAnimFrame = anim.AnimStartFrame
		e.FrameTimer = 0
	}

	timePerFrame := 1.0 / anim.AnimationSpeed
	e.FrameTimer += dt

	for e.FrameTimer >= timePerFrame {
		e.FrameTimer -= timePerFrame
		e.CurrAnimFrame++

		if e.CurrAnimFrame >= anim.AnimStartFrame+anim.TotalFrames {
			if anim.Looping {
				e.CurrAnimFrame = anim.AnimStartFrame
			} else {
				// Non-looping end behavior by state
				if e.State.IsEnemyState(StateLanding) ||
					e.State.IsEnemyAttacking() ||
					e.State.IsEnemyDefending() {
					e.State.SetEnemyState(StateIdle)
					e.CurrAnimFrame = 0
					e.FrameTimer = 0
					return
				} else if e.State.IsEnemyJumping() {
					e.State.SetEnemyState(StateFalling)
					e.CurrAnimFrame = anim.AnimStartFrame
					e.FrameTimer = 0
					return
				}
				// Dead / other: hold on last frame
				e.CurrAnimFrame = anim.AnimStartFrame + anim.TotalFrames - 1
			}
		}
	}
}

func (e *EnemyRuntime) DrawEnemyAnimation(screen *ebiten.Image, img *ebiten.Image, animations *map[int]Animation, camera core.Camera) {
	bounds := e.GetBounds()

	// Look up animation for current state, fallback to idle
	currState := e.State.Current
	anim, ok := (*animations)[currState]
	if !ok {
		anim = (*animations)[StateIdle]
	}

	width := anim.FrameWidth
	height := anim.FrameHeight
	frameX := e.CurrAnimFrame * width
	frameY := anim.SpriteSheetYPosition * frameHeight_minimum

	rect := image.Rect(frameX, frameY, frameX+width, frameY+height)
	subImage := img.SubImage(rect).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}

	if e.FlipX {
		// Flip horizontally
		op.GeoM.Scale(-e.Scale, e.Scale)
		// Translate back because flipping moves the image to the left of the axis
		op.GeoM.Translate(float64(width)*e.Scale, 0)
	} else {
		op.GeoM.Scale(e.Scale, e.Scale)
	}

	// Center horizontally on collision box
	drawX := bounds.X + (bounds.Width-float64(width)*e.Scale)/2
	// Align bottom to collision box bottom
	drawY := bounds.Y + (bounds.Height - float64(height)*e.Scale)

	// Apply camera offset
	drawX -= camera.Pos.X
	drawY -= camera.Pos.Y

	op.GeoM.Translate(drawX, drawY)

	screen.DrawImage(subImage, op)
}
