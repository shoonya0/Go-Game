package core

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

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

	playerSpriteSheetPath = "../assets/GideonGraves.png"
)

var playerSpriteSheet *ebiten.Image

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

	var err error
	playerSpriteSheet, _, err = ebitenutil.NewImageFromFile(playerSpriteSheetPath)
	if err != nil {
		log.Fatal(err)
	}

	return animations
}

func (player *PlayerRuntime) UpdateAnimation() {
	// // update player animation
	tps := float64(ebiten.TPS())
	if tps <= 0 {
		tps = 60
	}

	// here DT signifies the time in seconds between each frame of the animation means how long the current frame is displayed for
	currState := player.State.GetPlayerState()
	anim := player.Animations[currState]

	if player.PreviousState.GetPlayerState() != player.State.GetPlayerState() {
		player.CurrAnimFrame = anim.AnimStartFrame
	}

	timePerFrame := 1.0 / anim.AnimationSpeed // time (in seconds)to display each frame
	dt := 1.0 / tps                           // time (in seconds) between each frame
	anim.FrameTimer += dt                     // add the time (in seconds) between each frame to the frame timer

	// the for loop helps to keep the animation running at the correct speed
	for anim.FrameTimer >= timePerFrame {
		anim.FrameTimer -= timePerFrame

		player.CurrAnimFrame++

		// fmt.Println("player.CurrAnimFrame", player.CurrAnimFrame, anim.AnimStartFrame+anim.TotalFrames, currState)
		if player.CurrAnimFrame >= anim.AnimStartFrame+anim.TotalFrames {
			if anim.Looping {
				player.CurrAnimFrame = anim.AnimStartFrame
			} else {
				if player.State.IsLanding() || player.State.IsSmugFace() || player.State.IsWeakAttack() ||
					player.State.IsStrongAttack() || player.State.IsSpecialAttack1() || player.State.IsSpecialAttack2() ||
					player.State.IsSpecialAttack3() || player.State.IsSpecialAttack4() {
					player.State.SetPlayerState(int(PlayerStateIdle))
					player.CurrAnimFrame = 0
				} else if player.State.IsWeakAttackInAir() || player.State.IsStrongAttackInAir() || player.State.IsJumping() {
					player.State.SetPlayerState(int(PlayerStateFalling))
				}
			}
		}
	}
}

func (player *PlayerRuntime) DrawPlayerAnimation(screen *ebiten.Image) {

	bounds := player.GetBounds()
	// Adjust bounds for camera
	// drawBoundsX := float32(bounds.X - player.Camera.Pos.X)
	// drawBoundsY := float32(bounds.Y - player.Camera.Pos.Y)

	// vector.StrokeRect(screen, drawBoundsX, drawBoundsY, float32(bounds.Width), float32(bounds.Height), 1, color.White, false)
	// // draw ground sensor
	// groundSensor := player.GetGroundSensor()
	// vector.FillRect(screen, float32(groundSensor.X-player.Camera.Pos.X), float32(groundSensor.Y-player.Camera.Pos.Y), float32(groundSensor.Width), float32(groundSensor.Height), color.RGBA{255, 0, 0, 50}, false)

	currState := player.State.GetPlayerState()
	width := player.Animations[currState].FrameWidth
	height := player.Animations[currState].FrameHeight
	frameX := player.CurrAnimFrame * width
	frameY := player.Animations[currState].SpriteSheetYPosition * frameHeight_minimum

	rect := image.Rect(frameX, frameY, frameX+width, frameY+height)
	subImage := playerSpriteSheet.SubImage(rect).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}

	// Scale factor to resize the player image
	// Change this value to make the player smaller (e.g., 0.5) or larger (e.g., 2.0)

	if player.FlipX {
		// Flip horizontally
		op.GeoM.Scale(-player.Scale, player.Scale)
		// Translate back because flipping moves the image to the left of the axis
		op.GeoM.Translate(float64(width)*player.Scale, 0)
	} else {
		op.GeoM.Scale(player.Scale, player.Scale)
	}

	// Calculate draw position to center the sprite on the collision box
	// Center horizontally: bounds.X + (bounds.Width - spriteWidth) / 2
	drawX := bounds.X + (bounds.Width-float64(width)*player.Scale)/2

	// Align bottom vertically: bounds.Y + (bounds.Height - spriteHeight)
	drawY := bounds.Y + (bounds.Height - float64(height)*player.Scale)

	// Apply Camera Offset
	drawX -= player.Camera.Pos.X
	drawY -= player.Camera.Pos.Y

	// Move to the calculated position
	op.GeoM.Translate(drawX, drawY)

	screen.DrawImage(subImage, op)
}
