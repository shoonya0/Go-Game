package core

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// ------------------------ physics constants ------------------------
const (
	AccX         = 10
	AccY         = 10
	DecX         = 15
	DecY         = 10
	MaxSpeed     = 300
	JumpForce    = 500
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
		Pos: Position{
			X: 0,
			Y: 0,
		},
		Physics: Physics{
			VelX:         0,
			VelY:         0,
			AccX:         AccX,
			AccY:         AccY, // acceleration of the object when jumping
			DecX:         DecX,
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

// GetBounds returns the bounding box of the player
func (player *PlayerRuntime) GetBounds() AABB {
	return AABB{
		X:      player.Pos.X,
		Y:      player.Pos.Y,
		Width:  40,
		Height: 80,
	}
}

func UpdatePlayer(player *PlayerRuntime, inputState *InputState, qt *DynamicQuadtree) {
	// update player physics

	player.PreviousState = PlayerState{CurrentState: PlayerStateType(player.State.GetPlayerState())}

	tps := float64(ebiten.TPS())
	if tps <= 0 {
		tps = 60
	}

	dt := 1.0 / tps

	inputX := float64(inputState.Direction.LeftRight)
	targetVX := inputX * player.Physics.MaxSpeed

	// dtUnits scales the physics constants which seem tuned for a different timeframe (e.g. 1 unit = 1 frame at 60fps?)
	dtUnits := 100.0 / tps

	accX := player.Physics.AccX * dtUnits
	decX := player.Physics.DecX * dtUnits

	step := accX
	if inputX == 0 && player.State.IsGrounded() {
		player.State.SetPlayerState(int(PlayerStateIdle))
		step = decX
	} else {
		if player.State.IsGrounded() {
			if inputX != 0 {
				player.State.SetPlayerState(int(PlayerStateMoving))
			}
		}
	}

	player.Physics.VelX = approach(player.Physics.VelX, targetVX, step)

	jump := inputState.JumpJustPressed

	player.Physics.VelY += player.Physics.GravityScale * dtUnits

	if jump && player.State.IsGrounded() {
		player.State.SetPlayerState(int(PlayerStateJumping))
		player.Physics.VelY = -player.Physics.JumpForce // Instant impulse
	}

	// Apply X Movement
	player.Pos.X += player.Physics.VelX * dt // Apply velocity over time

	// Check X Collisions
	if qt != nil {
		objs := qt.Retrieve(player.GetBounds())
		for _, obj := range objs {
			if _, ok := obj.(*Platform); ok {
				bounds := obj.GetBounds()
				if player.GetBounds().Intersects(bounds) {
					// Resolve X
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

	// Apply Y Movement
	player.Pos.Y += player.Physics.VelY * dt

	// Check Y Collisions
	if qt != nil {
		objs := qt.Retrieve(player.GetBounds())
		for _, obj := range objs {
			if _, ok := obj.(*Platform); ok {
				bounds := obj.GetBounds()
				if player.GetBounds().Intersects(bounds) {
					// Resolve Y
					if player.Physics.VelY > 0 { // Falling
						player.Pos.Y = bounds.Y - player.GetBounds().Height
						player.Physics.VelY = 0
						player.State.SetPlayerState(int(PlayerStateIdle))
					} else if player.Physics.VelY < 0 { // Jumping into ceiling
						player.Pos.Y = bounds.Y + bounds.Height
						player.Physics.VelY = 0
					}
				}
			}
		}
	}

	if player.State.IsGrounded() {
		if player.Physics.VelX == 0 {
			player.State.SetPlayerState(int(PlayerStateIdle))
		} else {
			player.State.SetPlayerState(int(PlayerStateMoving))
		}
	} else {
		if player.Physics.VelY > 0 {
			// player.State.SetPlayerState(int(PlayerStateFalling))
		}
	}
}
