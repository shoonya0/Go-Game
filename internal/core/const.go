// internal/core/player_runtime.go
package core

// ========================= player runtime =========================
type Physics struct {
	VelX, VelY   float64 // current velocity of the object
	AccX, AccY   float64 // acceleration of the object
	DecX, DecY   float64 // deceleration of the object
	MaxSpeed     float64 // maximum speed of the object
	JumpForce    float64 // force applied to the object when jumping
	GravityScale float64 // scale of the gravity applied to the object
	TerminalVelY float64 // maximum falling speed
	CoyoteMs     float64 // remaining ms of coyote time (extra jump time after leaving a platform)
	AirJumpsLeft int     // number of air jumps left
}

type Combat struct {
	HP, MaxHP int
	InvulnMs  float64 // remaining ms of i-frames
}

// ---------------- direction ----------------
type Direction struct {
	LeftRight int8
	UpDown    int8
}

// ---------------- skills ----------------
type Skills struct {
	WeakAttack   bool
	StrongAttack bool
	Defense      bool
	UsePotion    bool
}

// ---------------- input state ----------------
type InputState struct {
	Direction       Direction
	JumpJustPressed bool // true if the jump button was just pressed
	JumpHeld        bool // true if the jump button is held down means player is in the air
	DashJustPressed bool // true if the dash button was just pressed
	Menu            bool // true if the menu button was just pressed
	Skills          Skills
}

// ---------------- position ----------------
type Position struct{ X, Y float64 } // position of the player in the world

// ---------------- player runtime ----------------
type PlayerRuntime struct {
	State         PlayerState
	PreviousState PlayerState
	Animations    map[int]*Animation
	CurrAnimFrame int // current frame of the current animation
	Pos           Position
	Physics       Physics
	Combat        Combat
	CheckpointID  string
}

// ---------------- game state ----------------
type GameState int

const (
	ModeMenu GameState = iota
	ModePlaying
	ModeGameOver
)
