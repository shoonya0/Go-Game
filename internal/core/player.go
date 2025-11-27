package core

import (
	"fmt"
)

// ------------------------ physics constants ------------------------
const (
	AccX         = 10
	AccY         = 10
	DecX         = 30
	DecY         = 10
	MaxSpeed     = 300
	JumpForce    = 5
	GravityScale = 0.5
	TerminalVelY = 10
	CoyoteMs     = 1
	AirJumpsLeft = 1
)

func InitPlayer() PlayerRuntime {
	return PlayerRuntime{
		Pos: Position{
			X: 0,
			Y: 0,
		},
		Physics: Physics{
			VelX:         0,
			VelY:         0,
			AccX:         AccX,
			AccY:         AccY,
			DecX:         DecX,
			DecY:         DecY,
			MaxSpeed:     MaxSpeed,
			JumpForce:    JumpForce,
			GravityScale: GravityScale,
			TerminalVelY: TerminalVelY,
			CoyoteMs:     CoyoteMs,
			AirJumpsLeft: AirJumpsLeft,
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

func approach(current, target, maxDelta float64) float64 {
	if current < target {
		current += maxDelta
		if current > target {
			return target
		}
		return current
	}
	if current > target {
		current -= maxDelta
		if current < target {
			return target
		}
		return current
	}
	return current
}

func UpdatePlayer(player *PlayerRuntime, inputState *InputState) {
	// update player physics
	inputX := float64(inputState.Direction.LeftRight)
	targetVX := inputX * player.Physics.MaxSpeed

	step := player.Physics.AccX
	if inputX == 0 {
		step = player.Physics.DecX // friction when no input
	} else if inputX < 0 && player.Physics.VelX > 0 {
		player.Physics.VelX = 0
	} else if inputX > 0 && player.Physics.VelX < 0 {
		player.Physics.VelX = 0
	}

	player.Physics.VelX = approach(player.Physics.VelX, targetVX, step)

	// if player is in free fall / jumping then only we can alter the y velocity

	// inputY := float64(inputState.Direction.UpDown)
	// targetVY := inputY * player.Physics.MaxSpeed
	// step = player.Physics.AccY
	// if inputY == 0 {
	// 	step = player.Physics.DecY // friction when no input
	// } else if inputY < 0 && player.Physics.VelY > 0 {
	// 	player.Physics.VelY = 0
	// } else if inputY > 0 && player.Physics.VelY < 0 {
	// 	player.Physics.VelY = 0
	// }
	// player.Physics.VelY = approach(player.Physics.VelY, targetVY, step)

	fmt.Println("VelX:", player.Physics.VelX, "VelY:", player.Physics.VelY)
}
