package core

import (
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

// physicsUnitsPerTick scales the frame-independent physics constants (expressed
// per 1/100s) into the current tick. dtUnits in the original code.
func physicsUnitsPerTick() float64 { return 100.0 * deltaTime() }

// player collision box and sensor dimensions.
const (
	bodyWidth       = 30.0
	bodyHeight      = 80.0
	sensorHeight    = 50.0
	footProbeHeight = 6.0 // thin probe at the feet for hazard "standing on" checks
	airRunControl   = 1.5 // running air-control is 1/airRunControl of ground running
	waterSinkTiles  = 0.5 // tiles the player sinks below a water surface before standing
	waterJumpFactor = 0.8 // fraction of jump force retained while submerged in water
)

// InitPlayer builds a player at its default spawn with all subsystems wired up.
func InitPlayer(img *ebiten.Image) *PlayerRuntime {
	return &PlayerRuntime{
		img:           img,
		State:         PlayerState{CurrentState: PlayerStateIdle},
		PreviousState: PlayerState{CurrentState: PlayerStateIdle},
		Animations:    InitPlayerAnimations(),
		FlipX:         false,
		Scale:         1.0,
		Camera:        Camera{Zoom: 1.0},
		Pos:           Position{X: 100, Y: 100},
		RespawnPos:    Position{X: 100, Y: 100},
		Physics: Physics{
			AccX:         AccX,
			AccY:         AccY,
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
		Combat:       newCombat(),
		CheckpointID: "default",
	}
}

// Approach moves current toward target by at most maxDelta, without overshoot.
func Approach(current, target, maxDelta float64) float64 {
	if current < target {
		return math.Min(current+maxDelta, target)
	}
	if current > target {
		return math.Max(current-maxDelta, target)
	}
	return current
}

// ReduceLeft raises current toward target (target below zero) without overshoot.
func ReduceLeft(current, target, maxDelta float64) float64 {
	if current < target {
		return math.Min(current+maxDelta, target)
	}
	return current
}

// ReduceRight lowers current toward target without overshoot.
func ReduceRight(current, target, maxDelta float64) float64 {
	if current > target {
		return math.Max(current-maxDelta, target)
	}
	return current
}

// GetBounds returns the player's collision box.
func (player *PlayerRuntime) GetBounds() AABB {
	return AABB{X: player.Pos.X, Y: player.Pos.Y, Width: bodyWidth, Height: bodyHeight}
}

// GetGroundSensor returns the thin box below the player used to detect ground.
func (player *PlayerRuntime) GetGroundSensor() AABB {
	return AABB{X: player.Pos.X, Y: player.Pos.Y + bodyHeight, Width: bodyWidth, Height: sensorHeight}
}

// GetFootSensor returns a thin box straddling the player's feet, used to detect a
// hazard tile the player is standing on without reaching a full tile below.
func (player *PlayerRuntime) GetFootSensor() AABB {
	return AABB{X: player.Pos.X, Y: player.Pos.Y + bodyHeight - footProbeHeight/2, Width: bodyWidth, Height: footProbeHeight}
}

// canMove reports whether the current state accepts movement input. Movement is
// allowed in the air and on the ground, but blocked during grounded actions.
func (player *PlayerRuntime) canMove() bool {
	s := &player.State
	return s.IsIdle() || s.IsMoving() || s.IsRunning() || s.IsJumping() ||
		s.IsFalling() || s.IsWeakAttackInAir() || s.IsStrongAttackInAir()
}

// UpdatePlayer advances the player one tick: input, physics, collision, state and
// combat (power regen, hazard/fall damage, checkpoint respawn).
func UpdatePlayer(player *PlayerRuntime, in *InputState, qt *DynamicQuadtree, checkpoints []Checkpoint) {
	player.PreviousState = PlayerState{CurrentState: PlayerStateType(player.State.GetPlayerState())}

	player.applyHorizontalMovement(in)
	player.applyGravityAndJump(in)

	player.resolveHorizontal(qt)
	onGround, detectGround := player.resolveVertical(qt)

	player.updateState(in, onGround, detectGround)

	player.regeneratePower()
	player.updateCombat(qt, checkpoints, onGround)

	if qt != nil {
		qt.Update(player)
	}
}

// applyHorizontalMovement updates VelX from input.
func (player *PlayerRuntime) applyHorizontalMovement(in *InputState) {
	inputX := 0.0
	if player.canMove() {
		inputX = float64(in.Direction.LeftRight)
	}

	targetVX := inputX * player.Physics.MaxSpeed
	// Running is disabled while submerged; water keeps you to walking speed.
	if in.RunJustPressed && !player.inWater {
		if player.State.IsGrounded() {
			targetVX = inputX * player.Physics.MaxRunSpeed
		} else {
			targetVX = inputX * (player.Physics.MaxRunSpeed / airRunControl)
		}
	}

	units := physicsUnitsPerTick()
	accX := player.Physics.AccX * units
	decX := player.Physics.DecX * units

	// Flip the sprite to match the movement direction.
	if inputX < 0 {
		player.FlipX = true
	} else if inputX > 0 {
		player.FlipX = false
	}

	// Above walk speed while not holding run: bleed velocity back down to MaxSpeed.
	if math.Abs(player.Physics.VelX) > player.Physics.MaxSpeed && !in.RunJustPressed {
		if player.Physics.VelX > 0 {
			player.Physics.VelX = ReduceRight(player.Physics.VelX, player.Physics.MaxSpeed, decX)
		} else if player.Physics.VelX < 0 {
			player.Physics.VelX = ReduceLeft(player.Physics.VelX, -player.Physics.MaxSpeed, decX)
		}
		return
	}

	// Accelerate toward target; apply friction (decX) when idle on the ground.
	step := accX
	if inputX == 0 && player.State.IsGrounded() {
		step = decX
	}
	player.Physics.VelX = Approach(player.Physics.VelX, targetVX, step)
}

// applyGravityAndJump integrates gravity and applies a jump impulse if requested.
func (player *PlayerRuntime) applyGravityAndJump(in *InputState) {
	player.Physics.VelY += player.Physics.GravityScale * physicsUnitsPerTick()

	if in.JumpJustPressed && player.State.IsGrounded() && player.canMove() {
		player.State.SetPlayerState(int(PlayerStateJumping))
		jumpForce := player.Physics.JumpForce
		if player.inWater {
			jumpForce *= waterJumpFactor // water dampens the jump
		}
		player.Physics.VelY = -jumpForce
	}
}

// resolveHorizontal integrates X and pushes the player out of solid tiles.
func (player *PlayerRuntime) resolveHorizontal(qt *DynamicQuadtree) {
	player.Pos.X += player.Physics.VelX * deltaTime()
	if qt == nil {
		return
	}

	for _, obj := range qt.Retrieve(player.GetBounds()) {
		plat, ok := obj.(*Platform)
		if !ok || plat.TileInfo.TileType == Water {
			continue // water is not solid; the player moves through it
		}
		bounds := plat.GetBounds()
		if !player.GetBounds().Intersects(bounds) {
			continue
		}
		if player.Physics.VelX > 0 { // moving right
			player.Pos.X = bounds.X - bodyWidth
		} else if player.Physics.VelX < 0 { // moving left
			player.Pos.X = bounds.X + bounds.Width
		}
		player.Physics.VelX = 0
	}
}

// resolveVertical integrates Y, resolves floor/ceiling collisions and reports
// whether the player landed (onGround) and whether ground is nearby (detectGround).
func (player *PlayerRuntime) resolveVertical(qt *DynamicQuadtree) (onGround, detectGround bool) {
	if player.Physics.VelY >= player.Physics.MaxFallSpeed {
		player.Physics.VelY = player.Physics.MaxFallSpeed
	}
	player.Pos.Y += player.Physics.VelY * deltaTime()

	if qt == nil {
		return false, false
	}

	sensor := player.GetGroundSensor()
	body := player.GetBounds()
	inWater := false
	waterSurfaceY := 0.0

	for _, obj := range qt.Retrieve(body) {
		plat, ok := obj.(*Platform)
		if !ok {
			continue
		}
		bounds := plat.GetBounds()

		// Water is not a solid floor/wall; instead record the topmost water
		// surface the player is submerged in so it can rest one tile below it.
		if plat.TileInfo.TileType == Water {
			if body.Intersects(bounds) {
				if !inWater || bounds.Y < waterSurfaceY {
					waterSurfaceY = bounds.Y
				}
				inWater = true
			}
			continue
		}

		if sensor.Intersects(bounds) {
			detectGround = true
		}

		if !body.Intersects(bounds) {
			continue
		}
		if player.Physics.VelY > 0 { // landing
			player.Pos.Y = bounds.Y - bodyHeight
			onGround = true
			player.Physics.VelY = 0
		} else if player.Physics.VelY < 0 { // head bonk
			player.Pos.Y = bounds.Y + bounds.Height
			player.Physics.VelY = 0
		}
	}

	// Buoyancy: the player sinks waterSinkTiles below the surface and stands
	// there (on the second water tile), rather than resting on the surface.
	if inWater {
		floorY := waterSurfaceY + waterSinkTiles*LevelTileHeight
		if player.Pos.Y+bodyHeight > floorY {
			player.Pos.Y = floorY - bodyHeight
			player.Physics.VelY = 0
			onGround = true
			detectGround = true
		}
	}

	player.inWater = inWater
	return onGround, detectGround
}

// updateState transitions the player state based on this tick's physics results.
func (player *PlayerRuntime) updateState(in *InputState, onGround, detectGround bool) {
	s := &player.State

	if !detectGround {
		player.updateAirState(in)
		return
	}

	if !onGround {
		// Above ground but not landed: finish a fall or air attack into a land.
		if (s.IsFalling() && player.Physics.VelY >= 0) || s.IsWeakAttackInAir() || s.IsStrongAttackInAir() {
			s.SetPlayerState(int(PlayerStateLanding))
		}
		return
	}

	if !s.IsGrounded() {
		if s.IsFalling() && player.Physics.VelY >= 0 {
			s.SetPlayerState(int(PlayerStateIdle))
		}
		return
	}

	// Grounded: actions take priority, then locomotion from velocity.
	switch {
	case in.SmugFace:
		s.SetPlayerState(int(PlayerStateSmugFace))
	case in.Skills.SpecialAttack1:
		s.SetPlayerState(int(PlayerStateSpecialAttack1))
	case in.Skills.SpecialAttack2:
		s.SetPlayerState(int(PlayerStateSpecialAttack2))
	case in.Skills.SpecialAttack3 && player.trySpendPower(SpecialPowerCost): // U
		s.SetPlayerState(int(PlayerStateSpecialAttack3))
	case in.Skills.SpecialAttack4 && player.trySpendPower(SpecialPowerCost): // O
		s.SetPlayerState(int(PlayerStateSpecialAttack4))
	case in.Skills.WeakAttack:
		s.SetPlayerState(int(PlayerStateWeakAttack))
	case in.Skills.StrongAttack:
		s.SetPlayerState(int(PlayerStateStrongAttack))
	case player.Physics.VelX == 0:
		s.SetPlayerState(int(PlayerStateIdle))
	case math.Abs(player.Physics.VelX) > player.Physics.MaxSpeed:
		s.SetPlayerState(int(PlayerStateRunning))
	default:
		s.SetPlayerState(int(PlayerStateMoving))
	}
}

// updateAirState handles state transitions while the player is airborne.
func (player *PlayerRuntime) updateAirState(in *InputState) {
	s := &player.State
	switch {
	case s.IsFalling():
		if in.Skills.WeakAttack {
			s.SetPlayerState(int(PlayerStateWeakAttackInAir))
		} else if in.Skills.StrongAttack {
			s.SetPlayerState(int(PlayerStateStrongAttackInAir))
		}
	case player.Physics.VelY > 0 && !s.IsLanding() && !s.IsWeakAttackInAir() && !s.IsStrongAttackInAir():
		s.SetPlayerState(int(PlayerStateFalling))
	}
}

// UpdateCamera follows the player and clamps the view to the level bounds.
func (player *PlayerRuntime) UpdateCamera(screenWidth, screenHeight, levelWidth, levelHeight float64) {
	player.Camera.Pos.X = clamp(player.Camera.Pos.X,
		player.Pos.X-2*screenWidth/3, player.Pos.X-screenWidth/3)
	player.Camera.Pos.Y = clamp(player.Camera.Pos.Y,
		player.Pos.Y-2*screenHeight/3, player.Pos.Y-screenHeight/3)

	// Keep the camera inside the level.
	player.Camera.Pos.X = clamp(player.Camera.Pos.X, 0, levelWidth-screenWidth)
	player.Camera.Pos.Y = clamp(player.Camera.Pos.Y, 0, levelHeight-screenHeight)
}

// clamp constrains v to the inclusive [lo, hi] range.
func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
