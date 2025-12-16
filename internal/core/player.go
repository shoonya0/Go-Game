package core

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// ------------------------ physics constants ------------------------
const (
	AccX         = 100
	AccY         = 10
	DecX         = 15
	DecY         = 10
	MaxSpeed     = 1050
	MaxRunSpeed  = 300
	JumpForce    = 1000
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
		Width:  30,
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
		step = decX
	}

	if inputX < 0 {
		player.FlipX = true
	} else if inputX > 0 {
		player.FlipX = false
	}

	player.Physics.VelX = approach(player.Physics.VelX, targetVX, step)

	jump := inputState.JumpJustPressed

	player.Physics.VelY += player.Physics.GravityScale * dtUnits

	if jump && player.State.IsGrounded() {
		player.State.SetPlayerState(int(PlayerStateJumping))
		player.Physics.VelY = -player.Physics.JumpForce // Instant impulse is usually better
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

	// Ground check flag
	onGround := false

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
						onGround = true
					} else if player.Physics.VelY < 0 { // Jumping into ceiling
						player.Pos.Y = bounds.Y + bounds.Height
						player.Physics.VelY = 0
					}
				}
			}
		}
	}

	if onGround {
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
