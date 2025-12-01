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
	FlipX                bool
	FrameTimer           float64 // time in seconds that the current frame has been displayed for ie: if FrameTimer is x, then the current frame is displayed for x/AnimationSpeed seconds
	AnimationSpeed       float64 // no of frames to display per second in seconds
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
		FlipX:                false,
		FrameTimer:           0,
		AnimationSpeed:       6,
	}
	animations[Moving] = &Animation{
		CurrentState:         PlayerStateMoving,
		SpriteSheetYPosition: 1,
		TotalFrames:          6,
		FrameWidth:           frameWidth_minimum,
		FrameHeight:          frameHeight_small,
		FlipX:                false,
		FrameTimer:           0,
		AnimationSpeed:       6,
	}
	animations[Running] = &Animation{
		CurrentState:         PlayerStateRunning,
		SpriteSheetYPosition: 2,
		TotalFrames:          8,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FlipX:                false,
		FrameTimer:           0,
		AnimationSpeed:       10,
	}
	animations[Jumping] = &Animation{
		CurrentState:         PlayerStateJumping,
		SpriteSheetYPosition: 5,
		TotalFrames:          9,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FlipX:                false,
		FrameTimer:           0,
		AnimationSpeed:       10,
	}
	// animation yet to make
	animations[Grounded] = &Animation{
		CurrentState:         PlayerStateGrounded,
		SpriteSheetYPosition: 6,
		TotalFrames:          6,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FlipX:                false,
		FrameTimer:           0,
		AnimationSpeed:       1,
	}
	// animation yet to make
	animations[InAir] = &Animation{
		CurrentState:         PlayerStateInAir,
		SpriteSheetYPosition: 6,
		TotalFrames:          6,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FlipX:                false,
		FrameTimer:           0,
		AnimationSpeed:       1,
	}
	animations[Damaged] = &Animation{
		CurrentState:         PlayerStateDamaged,
		SpriteSheetYPosition: 7,
		TotalFrames:          6,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FlipX:                false,
		FrameTimer:           0,
		AnimationSpeed:       1,
	}
	animations[Dead] = &Animation{
		CurrentState:         PlayerStateDead,
		SpriteSheetYPosition: 8,
		TotalFrames:          6,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FlipX:                false,
		FrameTimer:           0,
		AnimationSpeed:       1,
	}
	animations[WeakAttack] = &Animation{
		CurrentState:         PlayerStateWeakAttack,
		SpriteSheetYPosition: 8,
		TotalFrames:          6,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FlipX:                false,
		FrameTimer:           0,
		AnimationSpeed:       1,
	}
	animations[StrongAttack] = &Animation{
		CurrentState:         PlayerStateStrongAttack,
		SpriteSheetYPosition: 9,
		TotalFrames:          6,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FlipX:                false,
		FrameTimer:           0,
		AnimationSpeed:       1,
	}
	animations[Defense] = &Animation{
		CurrentState:         PlayerStateDefense,
		SpriteSheetYPosition: 10,
		TotalFrames:          6,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FlipX:                false,
		FrameTimer:           0,
		AnimationSpeed:       1,
	}
	animations[UsePotion] = &Animation{
		CurrentState:         PlayerStateUsePotion,
		SpriteSheetYPosition: 11,
		TotalFrames:          6,
		FrameWidth:           frameWidth_small,
		FrameHeight:          frameHeight_small,
		FlipX:                false,
		FrameTimer:           0,
		AnimationSpeed:       1,
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

	for anim.FrameTimer >= timePerFrame {
		anim.FrameTimer -= timePerFrame

		player.CurrAnimFrame++

		if player.CurrAnimFrame >= anim.TotalFrames {
			player.CurrAnimFrame = 0
		}
	}
}

func (player *PlayerRuntime) DrawAnimation(screen *ebiten.Image) {
	currState := player.State.GetPlayerState()
	width := player.Animations[currState].FrameWidth
	height := player.Animations[currState].FrameHeight
	frameX := player.CurrAnimFrame * width
	frameY := player.Animations[currState].SpriteSheetYPosition * height

	rect := image.Rect(frameX, frameY, frameX+width, frameY+height)
	subImage := playerSpriteSheet.SubImage(rect).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(player.Pos.X, player.Pos.Y)
	op.GeoM.Scale(1, 1)

	screen.DrawImage(subImage, op)
}
