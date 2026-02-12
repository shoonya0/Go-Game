package enemy

import (
	"player/internal/core"
)

type EnemyState struct {
	Current  EnemyStateType
	Previous EnemyStateType
}

type EnemyRuntime struct {
	// For Identification
	ID   int
	Name string

	// Core Stats (0-100 scale)
	Health       float64 // Curr health
	MaxHealth    float64 // Max health
	IQ           float64 // Strategic depth, affects decision quality
	BaseIQ       float64 // Original IQ (before berserk modifier)
	Strength     float64 // Damage output multiplier
	BaseStrength float64 // Original strength (before berserk modifier)

	// Position and Physics
	Pos     core.Position
	Physics core.Physics
	FlipX   bool    // face direction
	Scale   float64 // sprite scale

	// State Machine
	State       EnemyState
	PartyStatus PartyStatus
	Grounded    bool // true if the enemy is touching the ground

	// Berserk Mode
	BerserkActive   bool    // true if the enemy is in berserk mode
	BerserkDuration float64 // duration of berserk mode in seconds
	BeserkCoolDown  float64 // cooldown time after berserk mode in seconds
}

const (
	EnemyAccX         = 100
	EnemyAccY         = 10
	EnemyDecX         = 10
	EnemyDecY         = 10
	EnemyMaxSpeed     = 200
	EnemyMaxRunSpeed  = 500
	EnemyMaxFallSpeed = 700
	EnemyJumpForce    = 700
	EnemyGravityScale = 10

	DefaultWidth          = 40
	DefaultHeight         = 60
	DefaultDetectionRange = 400 // range to detect player
	DefaultFleeRange      = 250 // range to flee from player
	DefaultAttackRange    = 50  // range to attack player
	DefaultAttackDamage   = 10  // damage dealt to player per attack
	DefaultRegenRate      = 5   // HP per second

	RestMinDuration = 3.0 // minimum seconds to rest
)
