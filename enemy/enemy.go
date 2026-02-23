package enemy

import (
	"fmt"
	"math"
	"player/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type EnemyState struct {
	Current  int
	Previous int
}

type EnemyRuntime struct {
	// For Identification
	ID   string
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
	CurrAnimFrame  int
	State          EnemyState
	PartyStatus    PartyStatus
	Grounded       bool    // true if the enemy is touching the ground
	PatrolDir      float64 // +1.0 = right, -1.0 = left; reverses on wall hit
	AttackCooldown float64 // seconds remaining before next attack allowed

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

func (em *EnemyManager) InitEnemy(pos core.Position) EnemyRuntime {
	return EnemyRuntime{
		ID:           em.generateEnemyID(),
		Name:         "Gideon Graves",
		Health:       100,
		MaxHealth:    100,
		IQ:           100,
		BaseIQ:       100,
		Strength:     100,
		BaseStrength: 100,
		Pos: core.Position{
			X: pos.X,
			Y: pos.Y,
		},
		Physics: core.Physics{
			VelX:         0,
			VelY:         0,
			AccX:         EnemyAccX,
			AccY:         EnemyAccY,
			DecX:         EnemyDecX,
			MaxSpeed:     EnemyMaxSpeed,
			MaxRunSpeed:  EnemyMaxRunSpeed,
			MaxFallSpeed: EnemyMaxFallSpeed,
			JumpForce:    EnemyJumpForce,
			GravityScale: EnemyGravityScale,
		},
		State: EnemyState{
			Current:  StateIdle,
			Previous: StateIdle,
		},
		PartyStatus:     PartySolo,
		Grounded:        true,
		PatrolDir:       1.0,
		AttackCooldown:  0,
		BerserkActive:   false,
		BerserkDuration: 0,
		BeserkCoolDown:  0,
	}
}

func (em *EnemyManager) generateEnemyID() string {
	em.nextID++
	if em.ID != "" {
		return fmt.Sprintf("%s-E-%d", em.ID, em.nextID)
	}
	return fmt.Sprintf("E-%d", em.nextID)
}

// GetBounds implements core.Collider so enemies can be tracked in the quadtree.
func (e *EnemyRuntime) GetBounds() core.AABB {
	return core.AABB{
		X:      e.Pos.X,
		Y:      e.Pos.Y,
		Width:  DefaultWidth,
		Height: DefaultHeight,
	}
}

// getGroundSensor returns a thin AABB just below the enemy's feet
// used to detect whether the enemy is near ground before physics resolves.
func (e *EnemyRuntime) getGroundSensor() core.AABB {
	return core.AABB{
		X:      e.Pos.X,
		Y:      e.Pos.Y + DefaultHeight,
		Width:  DefaultWidth,
		Height: 50,
	}
}

// decideAction runs the 3-tier AI and returns a horizontal input direction
// (-1 left, 0 stop, +1 right) and whether the enemy wants to attack this frame.
// This replaces the role of InputState in UpdatePlayer.
func (e *EnemyRuntime) decideAction(player *core.PlayerRuntime) (inputX float64, wantsAttack bool) {
	if e.State.IsEnemyDead() {
		return 0, false
	}

	dx := player.Pos.X - e.Pos.X
	dist := math.Abs(dx)

	// Tier 1: Flee — low health and player is close
	if e.Health < e.MaxHealth*0.30 && dist < DefaultFleeRange && e.State.CanEnemyFlee() {
		e.State.SetEnemyState(StateFleeing)
		if dx > 0 {
			return -1, false // player is to the right → run left
		}
		return 1, false
	}

	// Tier 2: Hunt / Attack — player is within detection range
	if dist < DefaultDetectionRange {
		if dist < DefaultAttackRange && e.AttackCooldown <= 0 && e.State.CanEnemyAttack() {
			return 0, true // stop and attack
		}
		e.State.SetEnemyState(StateHunting)
		if dx > 0 {
			return 1, false
		}
		return -1, false
	}

	// Tier 3: Patrol — oscillate using PatrolDir; wall-reversal handled in Update
	e.State.SetEnemyState(StatePatrolling)
	return e.PatrolDir, false
}

// Update advances the enemy simulation by one frame: AI decisions, physics
// integration, platform collision resolution, and state machine transitions.
// It mirrors the structure of core.UpdatePlayer but replaces InputState with
// the 3-tier AI from decideAction.
func (e *EnemyRuntime) Update(player *core.PlayerRuntime, qt *core.DynamicQuadtree) {
	// 1. Guard: dead enemies don't simulate
	if e.State.IsEnemyDead() {
		return
	}

	// 2. Time management
	tps := float64(ebiten.TPS())
	if tps <= 0 {
		tps = 60
	}

	dt := 1.0 / tps
	dtUnits := 100.0 / tps

	// 3. Tick attack cooldown
	if e.AttackCooldown > 0 {
		e.AttackCooldown -= dt
		if e.AttackCooldown < 0 {
			e.AttackCooldown = 0
		}
	}

	// 4. AI decision (replaces InputState)
	inputX, wantsAttack := e.decideAction(player)

	// 5. Permission gate
	if !e.State.CanEnemyMove() {
		inputX = 0
	}

	// 6. Target velocity
	targetVX := inputX * e.Physics.MaxSpeed
	if e.BerserkActive {
		targetVX = inputX * e.Physics.MaxRunSpeed
	}

	// 7. X physics (acceleration / friction)
	accX := e.Physics.AccX * dtUnits
	decX := e.Physics.DecX * dtUnits

	step := accX
	if inputX == 0 && e.Grounded {
		step = decX // friction when standing still on ground
	}

	if math.Abs(e.Physics.VelX) > e.Physics.MaxSpeed && !e.BerserkActive {
		// Decelerate from berserk/run speed back to normal max
		step = decX
		if e.Physics.VelX > 0 {
			e.Physics.VelX = core.ReduceRight(e.Physics.VelX, e.Physics.MaxSpeed, step)
		} else if e.Physics.VelX < 0 {
			e.Physics.VelX = core.ReduceLeft(e.Physics.VelX, -e.Physics.MaxSpeed, step)
		}
	} else {
		e.Physics.VelX = core.Approach(e.Physics.VelX, targetVX, step)
	}

	// 8. FlipX
	if inputX < 0 {
		e.FlipX = true
	} else if inputX > 0 {
		e.FlipX = false
	}

	// 9. Y physics (gravity)
	e.Physics.VelY += e.Physics.GravityScale * dtUnits

	// 10. Apply X, resolve X collisions
	e.Pos.X += e.Physics.VelX * dt
	if qt != nil {
		objs := qt.Retrieve(e.GetBounds())
		for _, obj := range objs {
			if _, ok := obj.(*core.Platform); ok {
				bounds := obj.GetBounds()
				if e.GetBounds().Intersects(bounds) {
					if e.Physics.VelX > 0 { // moving right
						e.Pos.X = bounds.X - DefaultWidth
					} else if e.Physics.VelX < 0 { // moving left
						e.Pos.X = bounds.X + bounds.Width
					}
					e.Physics.VelX = 0
					// Patrol: reverse direction on wall hit
					if e.State.IsEnemyPatrolling() {
						e.PatrolDir = -e.PatrolDir
					}
				}
			}
		}
	}

	// 11. Apply Y (capped at MaxFallSpeed), resolve Y collisions
	if e.Physics.VelY < e.Physics.MaxFallSpeed {
		e.Pos.Y += e.Physics.VelY * dt
	} else {
		e.Physics.VelY = e.Physics.MaxFallSpeed
		e.Pos.Y += e.Physics.VelY * dt
	}

	onGround := false
	detectGround := false

	if qt != nil {
		objs := qt.Retrieve(e.GetBounds())
		sensor := e.getGroundSensor()
		for _, obj := range objs {
			if _, ok := obj.(*core.Platform); ok {
				bounds := obj.GetBounds()

				// Ground detection (sensor only)
				if sensor.Intersects(bounds) {
					detectGround = true
				}

				// Physics collision (body only)
				if e.GetBounds().Intersects(bounds) {
					if e.Physics.VelY > 0 { // landing
						e.Pos.Y = bounds.Y - DefaultHeight
						onGround = true
						e.Physics.VelY = 0
					} else if e.Physics.VelY < 0 { // head bonk
						e.Pos.Y = bounds.Y + bounds.Height
						e.Physics.VelY = 0
					}
				}
			}
		}
	}

	// 12. Update grounded state
	e.Grounded = onGround

	// 13. Attack
	if wantsAttack {
		e.State.SetEnemyState(StateAttacking)
		e.AttackCooldown = 1.0
	}

	// 14. Physics-driven state transitions
	if detectGround {
		if onGround && e.State.IsEnemyFalling() {
			e.State.SetEnemyState(StateIdle)
		}
	} else {
		// Airborne: transition to falling if not already jumping/falling
		if e.Physics.VelY > 0 && !e.State.IsEnemyJumping() && !e.State.IsEnemyFalling() {
			e.State.SetEnemyState(StateFalling)
		}
	}

	// 15. Update spatial index (auto-inserts on first call)
	if qt != nil {
		qt.Update(e)
	}
}
