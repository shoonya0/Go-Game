package enemy

import (
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
