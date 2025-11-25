package core

func InitPlayer() PlayerRuntime {
	return PlayerRuntime{
		Pos: Position{
			X: 0,
			Y: 0,
		},
		Facing: 0,
		Physics: Physics{
			VelX:         0,
			VelY:         0,
			AccX:         0,
			AccY:         0,
			MaxSpeed:     300,
			JumpForce:    5,
			GravityScale: 0.5,
			TerminalVelY: 10,
			CoyoteMs:     1,
			AirJumpsLeft: 1,
		},
		Combat: Combat{
			HP:       100,
			MaxHP:    100,
			InvulnMs: 0,
		},
		Anim: Animation{
			State: PlayerState{
				IsIdle:     true,
				IsMoving:   false,
				IsJumping:  false,
				IsFalling:  false,
				IsGrounded: false,
				IsInAir:    false,
				IsDamaged:  false,
				IsDead:     false,
				IsMenuOpen: false,
			},
			Frame: 0,
			FlipX: false,
		},
		CheckpointID: "default",
	}
}
