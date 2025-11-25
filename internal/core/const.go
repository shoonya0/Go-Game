// internal/core/player_runtime.go
package core

// ========================= player runtime =========================
type Physics struct {
	VelX, VelY   float64 // current velocity of the player
	AccX, AccY   float64 // acceleration of the player
	MaxSpeed     float64 // maximum speed of the player
	JumpForce    float64 // force applied to the player when jumping
	GravityScale float64 // scale of the gravity applied to the player
	TerminalVelY float64 // maximum falling speed
	CoyoteMs     float64 // remaining ms of coyote time (extra jump time after leaving a platform)
	AirJumpsLeft int     // number of air jumps left
}

type Combat struct {
	HP, MaxHP int
	InvulnMs  float64 // remaining ms of i-frames
}

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
}

type Animation struct {
	State PlayerState
	Frame int  // current frame of the animation
	FlipX bool // whether the animation is flipped horizontally
}

type Direction struct {
	LeftRight int8
	UpDown    int8
}

type Skills struct {
	WeakAttack   bool
	StrongAttack bool
	Defense      bool
	UsePotion    bool
}

type InputState struct {
	Direction       Direction
	JumpJustPressed bool // true if the jump button was just pressed
	JumpHeld        bool // true if the jump button is held down means player is in the air
	DashJustPressed bool // true if the dash button was just pressed
	Menu            bool // true if the menu button was just pressed
	Skills          Skills
}

type Position struct{ X, Y float64 } // position of the player in the world

type PlayerRuntime struct {
	Pos          Position
	Facing       int // -1 left, +1 right
	Physics      Physics
	Combat       Combat
	Anim         Animation
	CheckpointID string
}
