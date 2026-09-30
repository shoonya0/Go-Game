package core

import "github.com/hajimehoshi/ebiten/v2"

// Position is a point in world space.
type Position struct{ X, Y float64 }

// Camera describes the visible viewport into the world.
type Camera struct {
	Pos  Position
	Zoom float64
}

// Physics holds the player's movement state and tuning parameters.
type Physics struct {
	VelX, VelY   float64 // current velocity
	AccX, AccY   float64 // acceleration
	DecX         float64 // deceleration (horizontal friction)
	MaxSpeed     float64 // maximum walking speed
	MaxRunSpeed  float64 // maximum running speed
	MaxFallSpeed float64 // maximum falling speed
	JumpForce    float64 // upward impulse applied on jump
	GravityScale float64 // gravity applied each tick
	TerminalVelY float64 // maximum vertical speed
	CoyoteMs     float64 // remaining coyote time (jump grace after leaving ground)
	AirJumpsLeft int     // remaining mid-air jumps
}

// PlayerRuntime is the complete, mutable state of the player for one session.
type PlayerRuntime struct {
	img *ebiten.Image

	// State and animation.
	State         PlayerState
	PreviousState PlayerState
	Animations    map[int]*Animation
	CurrAnimFrame int // current frame index within the active animation

	// Spatial state.
	FlipX   bool    // true when the player faces left
	Scale   float64 // render scale of the sprite
	Pos     Position
	Physics Physics
	Camera  Camera

	// Progression.
	Combat       Combat
	CheckpointID string
}
