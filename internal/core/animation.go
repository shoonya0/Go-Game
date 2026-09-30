package core

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// frameOriginHeight is the sprite-sheet row height used to locate a frame's
// vertical origin (SpriteSheetYPosition is expressed in these units).
const frameOriginHeight = frameHeight_minimum

// ---------------- states ----------------
const (
	Idle int = iota
	Moving
	Running
	Jumping
	Falling
	Landing
	SmugFace
	WeakAttack
	StrongAttack
	SpecialAttack1
	SpecialAttack2
	SpecialAttack3
	SpecialAttack4
	WeakAttackInAir
	StrongAttackInAir
	// not used yet
	Damaged
	Dead
	MenuOpen
	Defense
	UsePotion
)

// ---------------- animation ----------------
type Animation struct {
	CurrentState         PlayerStateType // Use the enum instead of embedded PlayerState
	SpriteSheetYPosition int             // the y position of the sprite sheet in the image
	TotalFrames          int
	AnimStartFrame       int
	FrameWidth           int
	FrameHeight          int
	FrameTimer           float64 // time in seconds that the current frame has been displayed for ie: if FrameTimer is x, then the current frame is displayed for x/AnimationSpeed seconds
	AnimationSpeed       float64 // no of frames to display per second in seconds
	Looping              bool    // true if the animation should loop
}

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
)

func InitPlayerAnimations() map[int]*Animation {
	animations := make(map[int]*Animation)
	animations[Idle] = &Animation{
		CurrentState:         PlayerStateIdle,
		SpriteSheetYPosition: 0,
		TotalFrames:          4,
		AnimStartFrame:       0,
		FrameWidth:           frameWidth_minimum,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       6,
		Looping:              true,
	}
	animations[Moving] = &Animation{
		CurrentState:         PlayerStateMoving,
		SpriteSheetYPosition: 2,
		TotalFrames:          5,
		AnimStartFrame:       0,
		FrameWidth:           frameWidth_minimum,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       6,
		Looping:              true,
	}
	animations[Running] = &Animation{
		CurrentState:         PlayerStateRunning,
		SpriteSheetYPosition: 4,
		TotalFrames:          8,
		AnimStartFrame:       0,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       10,
		Looping:              true,
	}
	animations[Jumping] = &Animation{
		CurrentState:         PlayerStateJumping,
		SpriteSheetYPosition: 10,
		TotalFrames:          3,
		AnimStartFrame:       0,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       10,
		Looping:              false,
	}
	animations[Falling] = &Animation{
		CurrentState:         PlayerStateFalling,
		SpriteSheetYPosition: 10,
		TotalFrames:          1,
		AnimStartFrame:       3,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       5,
		Looping:              true,
	}
	animations[Landing] = &Animation{
		CurrentState:         PlayerStateLanding,
		SpriteSheetYPosition: 10,
		TotalFrames:          3,
		AnimStartFrame:       7,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0.5,
		AnimationSpeed:       10,
		Looping:              false,
	}

	animations[SmugFace] = &Animation{
		CurrentState:         PlayerStateSmugFace,
		SpriteSheetYPosition: 8,
		TotalFrames:          11,
		AnimStartFrame:       0,
		FrameWidth:           frameWidth_minimum,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       7,
		Looping:              false,
	}

	animations[WeakAttack] = &Animation{
		CurrentState:         PlayerStateWeakAttack,
		SpriteSheetYPosition: 19,
		TotalFrames:          8,
		AnimStartFrame:       0,
		FrameWidth:           frameWidth_large,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       12,
		Looping:              false,
	}

	animations[StrongAttack] = &Animation{
		CurrentState:         PlayerStateStrongAttack,
		SpriteSheetYPosition: 21,
		TotalFrames:          8,
		AnimStartFrame:       0,
		FrameWidth:           frameWidth_large,
		FrameHeight:          frameHeight_medium,
		FrameTimer:           0,
		AnimationSpeed:       10,
		Looping:              false,
	}

	animations[SpecialAttack1] = &Animation{
		CurrentState:         PlayerStateSpecialAttack1,
		SpriteSheetYPosition: 24,
		TotalFrames:          12,
		AnimStartFrame:       0,
		FrameWidth:           frameWidth_maximum,
		FrameHeight:          frameHeight_medium,
		FrameTimer:           0,
		AnimationSpeed:       8,
		Looping:              false,
	}

	animations[SpecialAttack2] = &Animation{
		CurrentState:         PlayerStateSpecialAttack2,
		SpriteSheetYPosition: 27,
		TotalFrames:          13,
		AnimStartFrame:       0,
		FrameWidth:           frameWidth_maximum,
		FrameHeight:          frameHeight_medium,
		FrameTimer:           0,
		AnimationSpeed:       8,
		Looping:              false,
	}

	animations[SpecialAttack3] = &Animation{
		CurrentState:         PlayerStateSpecialAttack3,
		SpriteSheetYPosition: 30,
		TotalFrames:          9,
		AnimStartFrame:       0,
		FrameWidth:           frameWidth_maximum,
		FrameHeight:          frameHeight_medium,
		FrameTimer:           0,
		AnimationSpeed:       8,
		Looping:              false,
	}

	animations[SpecialAttack4] = &Animation{
		CurrentState:         PlayerStateSpecialAttack4,
		SpriteSheetYPosition: 39,
		TotalFrames:          17,
		AnimStartFrame:       0,
		FrameWidth:           frameWidth_maximum,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       12,
		Looping:              false,
	}

	animations[WeakAttackInAir] = &Animation{
		CurrentState:         PlayerStateWeakAttackInAir,
		SpriteSheetYPosition: 36,
		TotalFrames:          6,
		AnimStartFrame:       0,
		FrameWidth:           frameWidth_maximum,
		FrameHeight:          frameHeight_medium,
		FrameTimer:           0,
		AnimationSpeed:       10,
		Looping:              false,
	}

	animations[StrongAttackInAir] = &Animation{
		CurrentState:         PlayerStateStrongAttackInAir,
		SpriteSheetYPosition: 33,
		TotalFrames:          10,
		AnimStartFrame:       0,
		FrameWidth:           frameWidth_maximum,
		FrameHeight:          frameHeight_medium,
		FrameTimer:           0,
		AnimationSpeed:       8,
		Looping:              false,
	}

	// animation yet to make
	animations[Damaged] = &Animation{
		CurrentState:         PlayerStateDamaged,
		SpriteSheetYPosition: 7,
		TotalFrames:          6,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       1,
		Looping:              false,
	}
	animations[Dead] = &Animation{
		CurrentState:         PlayerStateDead,
		SpriteSheetYPosition: 8,
		TotalFrames:          6,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       1,
		Looping:              false,
	}

	animations[Defense] = &Animation{
		CurrentState:         PlayerStateDefense,
		SpriteSheetYPosition: 10,
		TotalFrames:          6,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       1,
		Looping:              false,
	}
	animations[UsePotion] = &Animation{
		CurrentState:         PlayerStateUsePotion,
		SpriteSheetYPosition: 11,
		TotalFrames:          6,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       1,
		Looping:              false,
	}

	return animations
}

// UpdateAnimation advances the active animation by one tick and applies the
// state transition that fires when a non-looping animation reaches its end.
func (player *PlayerRuntime) UpdateAnimation() {
	currState := player.State.GetPlayerState()
	anim := player.Animations[currState]

	// Restart the animation whenever the state changed this frame.
	if player.PreviousState.GetPlayerState() != currState {
		player.CurrAnimFrame = anim.AnimStartFrame
	}

	timePerFrame := 1.0 / anim.AnimationSpeed
	anim.FrameTimer += deltaTime()

	// Advance as many frames as the accumulated time allows.
	for anim.FrameTimer >= timePerFrame {
		anim.FrameTimer -= timePerFrame
		player.CurrAnimFrame++

		if player.CurrAnimFrame < anim.AnimStartFrame+anim.TotalFrames {
			continue
		}

		if anim.Looping {
			player.CurrAnimFrame = anim.AnimStartFrame
			continue
		}

		// Non-looping animation finished: resolve the follow-up state.
		switch {
		case player.State.IsGroundedOneShot():
			player.State.SetPlayerState(int(PlayerStateIdle))
			player.CurrAnimFrame = 0
		case player.State.IsAirOneShot():
			player.State.SetPlayerState(int(PlayerStateFalling))
		}
	}
}

// DrawPlayerAnimation draws the current animation frame, flipped to face the
// movement direction and offset so the sprite is centered on the collision box.
func (player *PlayerRuntime) DrawPlayerAnimation(screen *ebiten.Image) {
	bounds := player.GetBounds()
	anim := player.Animations[player.State.GetPlayerState()]
	width, height := anim.FrameWidth, anim.FrameHeight

	// Select the current frame from the sprite sheet.
	frameX := player.CurrAnimFrame * width
	frameY := anim.SpriteSheetYPosition * frameOriginHeight
	rect := image.Rect(frameX, frameY, frameX+width, frameY+height)
	subImage := player.img.SubImage(rect).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}
	if player.FlipX {
		// Mirror horizontally, then shift back into place.
		op.GeoM.Scale(-player.Scale, player.Scale)
		op.GeoM.Translate(float64(width)*player.Scale, 0)
	} else {
		op.GeoM.Scale(player.Scale, player.Scale)
	}

	// Center the sprite horizontally on the box, align its bottom, then apply camera.
	drawX := bounds.X + (bounds.Width-float64(width)*player.Scale)/2 - player.Camera.Pos.X
	drawY := bounds.Y + (bounds.Height - float64(height)*player.Scale) - player.Camera.Pos.Y
	op.GeoM.Translate(drawX, drawY)

	screen.DrawImage(subImage, op)
}
