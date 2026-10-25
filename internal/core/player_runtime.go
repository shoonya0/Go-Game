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
	RespawnPos   Position // where respawn() places the player (active checkpoint)
	Respawns     int      // times the player has respawned; the game revives enemies when it changes

	// Combat bookkeeping (not for external use).
	fallPeakY        float64 // highest point (min Y) reached during the current fall
	groundedPrev     bool    // whether the player was grounded last tick (landing edge)
	wasInLarva       bool    // whether the player was touching larva last tick (damage edge)
	waterDamageAccum float64 // sub-point water damage carried between ticks
	inWater          bool    // whether the player was submerged last tick (movement/jump)
	running          bool    // whether the player is running this tick (drains power)
	exhausted        bool    // power ran out while running; running is blocked until it recovers
	damageFlash      float64 // seconds of red damage flash left
	knockTicks       int     // ticks of heavy-hit push left
	atPortal         int     // index of the portal the player stands in, -1 for none
	savedNote        float64 // seconds the "saved" check stays up
	knockVelX        float64 // push velocity (px/s) while knockTicks last

	struck      []*EnemyRuntime // enemies already hit by the current hit of the swing (HitEnemies)
	strikeIdx   int             // which hit of the swing struck holds
	waves       []*GroundWave   // live U ground waves
	waveCast    bool            // the current U swing already cast its wave
	comboQueued bool            // K/L/O was pressed again before the combo gate (see comboGate)
}
