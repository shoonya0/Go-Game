package enemy

import (
	"player/internal/core"
	"sync"
)

type EnemyStateType int
type PartyStatus int

type EnemyState struct {
	Current  EnemyStateType
	Previous EnemyStateType
}

const (
	PartySolo   PartyStatus = iota // Alone
	PartyMember                    // Part of a group but not leader
	PartyLeader                    // Leading a group
)

// ========================================================================
// Q-Learning Reinforcement Learning Brain for Enemy AI
// ========================================================================
// State Vector for Q-Learning Brain Inputs (Inputs to the Neural Network)
type StateVector struct {
	SelfHealthBin    int // 0-4 (critical, low, medium, high, full)
	PlayerHealthBin  int // 0-4
	IQBin            int // 0-2 (low, medium, high)
	StrengthBin      int // 0-2
	PartyStatus      int // 0-2 (solo, member, leader)
	BerserkActive    int // 0-1
	PlayerDistance   int // 0-3 (close, medium, far, very far)
	PlayerDirection  int // 0-1 (left, right)
	OnGround         int // 0-1
	NearbyAllies     int // 0-2 (none, few, many)
	PlatformAhead    int // 0-1 (can continue, gap/wall)
	PlayerAboveBelow int // 0-2 (below, same level, above)
}

// ActionType represents possible actions the AI can take
type ActionType int

const (
	ActionIdle ActionType = iota
	ActionMoveLeft
	ActionMoveRight
	ActionJump
	ActionJumpLeft
	ActionJumpRight
	ActionAttack
	ActionFlee
	ActionEnterBerserk
	ActionExitBerserk
	ActionFormParty
	ActionBetrayAlly
	ActionRest
	ActionHuntPlayer
)

// QLearningBrain implements a Q-learning agent for enemy decision making
type QLearningBrain struct {
	QTable       [][]float64 // Q[state][action] = expected reward
	LearningRate float64     // Alpha: how much new info overrides old (0.1-0.3 typical)
	DiscountRate float64     // Gamma: importance of future rewards (0.9-0.99)
	Epsilon      float64     // Exploration rate (starts high, decays)
	EpsilonMin   float64     // Minimum exploration rate
	EpsilonDecay float64     // How fast epsilon decays

	// Experience tracking
	LastState  StateVector
	LastAction ActionType
	HasHistory bool

	// Statistics
	TotalReward    float64
	EpisodeRewards []float64
	StepCount      int

	mu sync.RWMutex
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

	// Party System
	PartyID      int             // -1 if solo, otherwise the party ID
	PartyMembers []*EnemyRuntime // pointer of party members
	PartyLeader  *EnemyRuntime   // pointer to the party leader(if member)

	// AI Brain
	Brains *QLearningBrain
}
