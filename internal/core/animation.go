package core

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// ---------------- states ----------------
const (
	Idle         = 0
	Moving       = 1
	Running      = 2
	Jumping      = 3
	Falling      = 4
	Grounded     = 5
	InAir        = 6
	Damaged      = 7
	Dead         = 8
	MenuOpen     = 9
	WeakAttack   = 10
	StrongAttack = 11
	Defense      = 12
	UsePotion    = 13
)

// ---------------- animation ----------------
type Animation struct {
	CurrentState         PlayerStateType // Use the enum instead of embedded PlayerState
	SpriteSheetYPosition int             // the y position of the sprite sheet in the image
	TotalFrames          int
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
	// frameWidth_large   = 160
	// frameWidth_maximum = 200

	// frameHeight_minimum = 40  // this frame is currently not used
	frameHeight_small = 80 // this frame is used
	// frameHeight_medium  = 120 // this frame is used
	// frameHeight_large   = 160 // this frame is currently not used
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
		FrameWidth:           frameWidth_minimum,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       6,
		Looping:              true,
	}
	animations[Moving] = &Animation{
		CurrentState:         PlayerStateMoving,
		SpriteSheetYPosition: 1,
		TotalFrames:          5,
		FrameWidth:           frameWidth_minimum,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       6,
		Looping:              true,
	}
	animations[Running] = &Animation{
		CurrentState:         PlayerStateRunning,
		SpriteSheetYPosition: 2,
		TotalFrames:          8,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       10,
		Looping:              true,
	}
	animations[Jumping] = &Animation{
		CurrentState:         PlayerStateJumping,
		SpriteSheetYPosition: 5,
		TotalFrames:          9,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       10,
		Looping:              false,
	}
	// animation yet to make
	animations[InAir] = &Animation{
		CurrentState:         PlayerStateInAir,
		SpriteSheetYPosition: 6,
		TotalFrames:          6,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       1,
		Looping:              true,
	}
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
	animations[WeakAttack] = &Animation{
		CurrentState:         PlayerStateWeakAttack,
		SpriteSheetYPosition: 8,
		TotalFrames:          6,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FrameTimer:           0,
		AnimationSpeed:       1,
		Looping:              false,
	}
	animations[StrongAttack] = &Animation{
		CurrentState:         PlayerStateStrongAttack,
		SpriteSheetYPosition: 9,
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

	timePerFrame := 1.0 / anim.AnimationSpeed // time in seconds to display each frame
	dt := 1.0 / tps                           // time in seconds between each frame
	anim.FrameTimer += dt                     // add the time in seconds between each frame to the frame timer

	// the for loop helps to keep the animation running at the correct speed
	for anim.FrameTimer >= timePerFrame {
		anim.FrameTimer -= timePerFrame

		player.CurrAnimFrame++

		if player.CurrAnimFrame >= anim.TotalFrames {
			player.CurrAnimFrame = 0
			if !anim.Looping {
				// anim.FrameTimer = 0
				// player.State.SetPlayerState(int(PlayerStateIdle))
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

	currState := player.State.GetPlayerState()
	width := player.Animations[currState].FrameWidth
	height := player.Animations[currState].FrameHeight
	frameX := player.CurrAnimFrame * width
	frameY := player.Animations[currState].SpriteSheetYPosition * height

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
