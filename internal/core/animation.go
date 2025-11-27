package core

// ---------------- player state ----------------
type PlayerState struct {
	IsIdle     bool
	IsMoving   bool
	IsJumping  bool
	IsFalling  bool
	IsGrounded bool
	IsInAir    bool
	IsDamaged  bool
	IsDead     bool
	IsMenuOpen bool
	// ---------------- skills ----------------
	WeakAttack   bool
	StrongAttack bool
	Defense      bool
	UsePotion    bool
}

// ---------------- animation ----------------
type Animation struct {
	State        PlayerState
	TotalFrames  int  // total frames of the animation
	CurrentFrame int  // current frame of the animation
	FrameWidth   int  // width of the frame
	FrameHeight  int  // height of the frame
	FlipX        bool // whether the animation is flipped horizontally
}

func InitAnimation() map[string]Animation {
	animations := make(map[string]Animation)
	animations["idle"] = Animation{
		State: PlayerState{
			IsIdle:    true,
			IsMoving:  false,
			IsJumping: false,
			IsFalling: false,
		},
		TotalFrames:  6,
		CurrentFrame: 0,
		FrameWidth:   40,
		FrameHeight:  80,
		FlipX:        false,
	}
	animations["moving"] = Animation{
		State: PlayerState{
			IsIdle:    false,
			IsMoving:  true,
			IsJumping: false,
			IsFalling: false,
		},
	}
	animations["jumping"] = Animation{
		State: PlayerState{
			IsIdle:    false,
			IsMoving:  false,
			IsJumping: true,
			IsFalling: false,
		},
	}
	animations["falling"] = Animation{
		State: PlayerState{
			IsIdle:    false,
			IsMoving:  false,
			IsJumping: false,
			IsFalling: true,
		},
	}
	animations["damaged"] = Animation{
		State: PlayerState{
			IsIdle:    false,
			IsMoving:  false,
			IsJumping: false,
			IsFalling: false,
		},
	}
	animations["dead"] = Animation{
		State: PlayerState{
			IsIdle:    false,
			IsMoving:  false,
			IsJumping: false,
			IsFalling: false,
		},
	}
	animations["menu"] = Animation{
		State: PlayerState{
			IsIdle:    false,
			IsMoving:  false,
			IsJumping: false,
			IsFalling: false,
		},
		Frame: 0,
		FlipX: false,
	}
	return animations
}

func UpdateAnimation(player *PlayerRuntime) {
	// update player animation

}
