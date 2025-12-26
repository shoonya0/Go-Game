package core

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// ------------------------ physics constants ------------------------
const (
	AccX         = 100
	AccY         = 10
	DecX         = 10
	DecY         = 10
	MaxSpeed     = 200
	MaxRunSpeed  = 500
	MaxFallSpeed = 700
	JumpForce    = 700
	GravityScale = 10
	TerminalVelY = 10
	CoyoteMs     = 1
	AirJumpsLeft = 1
)

func InitPlayer() PlayerRuntime {
	return PlayerRuntime{
		State:         PlayerState{CurrentState: PlayerStateIdle},
		PreviousState: PlayerState{CurrentState: PlayerStateIdle},
		Animations:    InitPlayerAnimations(),
		FlipX:         false,
		Scale:         1.0,
		Camera:        Camera{Zoom: 1.0},
		Pos: Position{
			X: 100,
			Y: 100,
		},
		Physics: Physics{
			VelX:         0,
			VelY:         0,
			AccX:         AccX,
			AccY:         AccY, // acceleration of the object when jumping
			DecX:         DecX,
			MaxSpeed:     MaxSpeed,
			MaxRunSpeed:  MaxRunSpeed,
			MaxFallSpeed: MaxFallSpeed,
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

func reduceLeft(current, target, maxDelta float64) float64 {
	if current < target {
		current += maxDelta
		if current > target {
			return target
		}
		return current
	}
	return current
}

func reduceRight(current, target, maxDelta float64) float64 {
	if current > target {
		current -= maxDelta
		if current < target {
			return target
		}
		return current
	}
	return current
}

// GetBounds returns the bounding box of the player
func (player *PlayerRuntime) GetBounds() AABB {
	return AABB{
		X:      player.Pos.X,
		Y:      player.Pos.Y,
		Width:  30,
		Height: 80,
	}
}

func (player *PlayerRuntime) GetGroundSensor() AABB {
	sensorHeight := 50.0
	return AABB{
		X:      player.Pos.X,
		Y:      player.Pos.Y + 80,
		Width:  30,
		Height: sensorHeight,
	}
}

func UpdatePlayer(player *PlayerRuntime, inputState *InputState, qt *DynamicQuadtree) {
	// Update previous state at the start of the frame
	player.PreviousState = PlayerState{CurrentState: PlayerStateType(player.State.GetPlayerState())}

	// Time Management
	tps := float64(ebiten.TPS())
	if tps <= 0 {
		tps = 60
	}
	dt := 1.0 / tps
	dtUnits := 100.0 / tps // Scaling factor for physics constants

	// Check Permissions
	// Define which states allow movement input.
	// We allow movement in air, idle, run, etc., but block it during "Action" states.
	canMove := player.State.IsIdle() || player.State.IsMoving() ||
		player.State.IsRunning() || player.State.IsJumping() ||
		player.State.IsFalling()

	// Input Processing
	inputX := 0.0
	if canMove {
		inputX = float64(inputState.Direction.LeftRight)
	}

	targetVX := inputX * player.Physics.MaxSpeed
	if inputState.RunJustPressed {
		// Run logic: slightly slower air control if not grounded
		if !player.State.IsGrounded() {
			targetVX = inputX * (player.Physics.MaxRunSpeed / 1.5)
		} else {
			targetVX = inputX * player.Physics.MaxRunSpeed
		}
	}

	// X Physics (Acceleration & Friction)
	accX := player.Physics.AccX * dtUnits
	decX := player.Physics.DecX * dtUnits

	step := accX
	// Apply friction if no input and on the ground
	if inputX == 0 && player.State.IsGrounded() {
		step = decX
	}

	// Handle Flipping
	if inputX < 0 {
		player.FlipX = true
	} else if inputX > 0 {
		player.FlipX = false
	}

	// Apply Velocity Changes
	if math.Abs(player.Physics.VelX) > player.Physics.MaxSpeed && !inputState.RunJustPressed {
		// Decelerate from run speed to normal max speed
		step = decX
		if player.Physics.VelX > 0 {
			player.Physics.VelX = reduceRight(player.Physics.VelX, player.Physics.MaxSpeed, step)
		} else if player.Physics.VelX < 0 {
			player.Physics.VelX = reduceLeft(player.Physics.VelX, -player.Physics.MaxSpeed, step)
		}
	} else {
		player.Physics.VelX = approach(player.Physics.VelX, targetVX, step)
	}

	// Y Physics (Gravity & Jumping)
	player.Physics.VelY += player.Physics.GravityScale * dtUnits

	// Jump Input
	if inputState.JumpJustPressed && player.State.IsGrounded() && canMove {
		player.State.SetPlayerState(int(PlayerStateJumping))
		player.Physics.VelY = -player.Physics.JumpForce // Instant impulse
	}

	// Integration & Collision Resolution

	// Apply X
	player.Pos.X += player.Physics.VelX * dt
	if qt != nil {
		objs := qt.Retrieve(player.GetBounds())
		for _, obj := range objs {
			if _, ok := obj.(*Platform); ok {
				bounds := obj.GetBounds()
				if player.GetBounds().Intersects(bounds) {
					if player.Physics.VelX > 0 { // Moving Right
						player.Pos.X = bounds.X - player.GetBounds().Width
					} else if player.Physics.VelX < 0 { // Moving Left
						player.Pos.X = bounds.X + bounds.Width
					}
					player.Physics.VelX = 0
				}
			}
		}
	}

	// Apply Y
	if player.Physics.VelY < player.Physics.MaxFallSpeed {
		player.Pos.Y += player.Physics.VelY * dt
	} else {
		player.Physics.VelY = player.Physics.MaxFallSpeed
		player.Pos.Y += player.Physics.VelY * dt
	}

	onGround := false
	detectGround := false

	if qt != nil {
		objs := qt.Retrieve(player.GetBounds())
		sensorObjs := player.GetGroundSensor()
		for _, obj := range objs {
			if _, ok := obj.(*Platform); ok {
				bounds := obj.GetBounds()

				// Ground detection (sensor only)
				if sensorObjs.Intersects(bounds) {
					detectGround = true
				}

				// Physics collision (body only)
				if player.GetBounds().Intersects(bounds) {
					if player.Physics.VelY > 0 { // Falling/Landing
						player.Pos.Y = bounds.Y - player.GetBounds().Height
						onGround = true
						player.Physics.VelY = 0
					} else if player.Physics.VelY < 0 { // Bonking head
						player.Pos.Y = bounds.Y + bounds.Height
						// fmt.Println("Bonking head and velY is", player.Physics.VelY)
						player.Physics.VelY = 0
					}
				}
			}
		}
	}

	// State Management
	// Transition states based on the physical results of this frame
	if detectGround {
		if onGround {
			if player.State.IsGrounded() {

				// Smug Face Input
				if inputState.SmugFace {
					fmt.Println("Smug Face Input")
					player.State.SetPlayerState(int(PlayerStateSmugFace))
				} else if inputState.Skills.SpecialAttack1 {
					fmt.Println("Special Attack 1 Input")
					player.State.SetPlayerState(int(PlayerStateSpecialAttack1))
				} else if inputState.Skills.SpecialAttack2 {
					fmt.Println("Special Attack 2 Input")
					player.State.SetPlayerState(int(PlayerStateSpecialAttack2))
				} else if inputState.Skills.SpecialAttack3 {
					fmt.Println("Special Attack 3 Input")
					player.State.SetPlayerState(int(PlayerStateSpecialAttack3))
				} else if inputState.Skills.SpecialAttack4 {
					fmt.Println("Special Attack 4 Input")
					player.State.SetPlayerState(int(PlayerStateSpecialAttack4))
				} else if inputState.Skills.WeakAttack {
					fmt.Println("Weak Attack Input")
					player.State.SetPlayerState(int(PlayerStateWeakAttack))
				} else if inputState.Skills.StrongAttack {
					fmt.Println("Strong Attack Input")
					player.State.SetPlayerState(int(PlayerStateStrongAttack))
				} else if player.Physics.VelX == 0 {
					player.State.SetPlayerState(int(PlayerStateIdle))
				} else {
					if math.Abs(player.Physics.VelX) > player.Physics.MaxSpeed {
						player.State.SetPlayerState(int(PlayerStateRunning))
					} else {
						player.State.SetPlayerState(int(PlayerStateMoving))
					}
				}
			} else if player.State.IsFalling() && player.Physics.VelY >= 0 {
				player.State.SetPlayerState(int(PlayerStateIdle))
			}
		} else if player.State.IsFalling() && player.Physics.VelY >= 0 {
			player.State.SetPlayerState(int(PlayerStateLanding))
		}
	} else {
		// In Air
		if player.Physics.VelY > 0 && !player.State.IsFalling() && !player.State.IsLanding() {
			player.State.SetPlayerState(int(PlayerStateFalling))
		}
	}

	// Update spatial partition
	qt.Update(player)
}

func (player *PlayerRuntime) UpdateCamera(screenWidth, screenHeight, levelWidth, levelHeight float64) {
	minX := player.Pos.X - 2*screenWidth/3
	maxX := player.Pos.X - screenWidth/3
	if player.Camera.Pos.X < minX {
		player.Camera.Pos.X = minX
	} else if player.Camera.Pos.X > maxX {
		player.Camera.Pos.X = maxX
	}

	minY := player.Pos.Y - 2*screenHeight/3
	maxY := player.Pos.Y - screenHeight/3
	if player.Camera.Pos.Y < minY {
		player.Camera.Pos.Y = minY
	} else if player.Camera.Pos.Y > maxY {
		player.Camera.Pos.Y = maxY
	}

	// Clamp camera to level bounds
	if player.Camera.Pos.X < 0 {
		player.Camera.Pos.X = 0
	}
	if player.Camera.Pos.Y < 0 {
		player.Camera.Pos.Y = 0
	}
	if player.Camera.Pos.X > levelWidth-screenWidth {
		player.Camera.Pos.X = levelWidth - screenWidth
	}
	if player.Camera.Pos.Y > levelHeight-screenHeight {
		player.Camera.Pos.Y = levelHeight - screenHeight
	}
}
