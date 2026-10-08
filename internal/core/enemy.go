package core

import "github.com/hajimehoshi/ebiten/v2"

// EnemyKind identifies a concrete enemy species. Adding a new enemy is a matter
// of adding a kind here and a blueprint in enemyBlueprints (enemy_types.go) — no
// changes to the shared update, physics, animation or draw code are needed.
type EnemyKind int

const (
	KindSuccubus EnemyKind = iota
)

// EnemyStats are the tunable gameplay numbers for an enemy species. They live in
// the blueprint (shared by every instance of that kind) rather than on the
// runtime, so all succubi share one set of stats.
type EnemyStats struct {
	MaxHealth      int     // hit points at spawn
	Damage         int     // health removed from the player per landed hit
	WalkSpeed      float64 // patrol speed (world px/s)
	ChaseSpeed     float64 // speed while pursuing the player (world px/s)
	DetectRange    float64 // horizontal distance at which the enemy notices the player
	AttackRange    float64 // horizontal distance at which the enemy can strike
	AttackCooldown float64 // minimum seconds between landed hits
	PatrolTiles    float64 // half-width (in tiles) of the patrol path around spawn
}

// EnemyBlueprint is the immutable description of an enemy species: which sprite
// sheet to use, how its frames are laid out, its collision box, and its stats.
// This is the "different sprites take on the general enemy" part — every field
// that varies per species is data here, while behaviour is shared code.
//
// Sprite sheets differ in frame size and layout, so each blueprint carries its
// own OriginHeight (the sprite-sheet row unit) and its own Animations table; the
// shared drawing code (drawSpriteFrame) is parameterised on both.
type EnemyBlueprint struct {
	Name         string
	SpritePath   string
	OriginHeight int     // sprite-sheet row height used to locate a frame's Y origin
	BodyWidth    float64 // collision box width
	BodyHeight   float64 // collision box height
	Scale        float64 // render scale of the sprite
	Stats        EnemyStats

	// NewAnimations builds a fresh animation table for one instance. It is a
	// factory (not a shared map) because Animation carries a mutable FrameTimer;
	// sharing one map across instances would couple their playback.
	NewAnimations func() map[int]*Animation
}

// EnemyRuntime is the complete, mutable state of one enemy instance. It mirrors
// the shape of PlayerRuntime (sprite + spatial + combat + per-tick bookkeeping)
// but only carries what an AI-driven enemy needs.
type EnemyRuntime struct {
	img *ebiten.Image
	bp  *EnemyBlueprint // the species blueprint (shared, read-only)

	// State and animation.
	State         EnemyState
	PrevState     EnemyState
	Animations    map[int]*Animation
	CurrAnimFrame int

	// Spatial state.
	FlipX      bool // true when the enemy faces left
	Scale      float64
	Pos        Position
	VelX, VelY float64

	// Combat.
	Health int

	// AI / per-tick bookkeeping (not for external use).
	spawnX      float64 // patrol anchor (X at spawn)
	patrolDir   float64 // current patrol direction: -1 or +1
	attackTimer float64 // seconds until the next hit may land
	onGround    bool    // whether the enemy rested on ground last tick
	deathDone   bool    // whether the death animation has finished
}

// enemyImages caches one GPU image per species so spawning many enemies of the
// same kind does not re-upload the sheet. The game loop is single-threaded, so
// no synchronisation is required.
var enemyImages = map[EnemyKind]*ebiten.Image{}

// enemyImage returns the (cached) sprite sheet for a species.
func enemyImage(kind EnemyKind, path string) *ebiten.Image {
	if img, ok := enemyImages[kind]; ok {
		return img
	}
	img := LoadImage(path)
	enemyImages[kind] = img
	return img
}

// NewEnemy builds a ready-to-run enemy of the given kind at pos, wiring its
// blueprint, sprite sheet, animations and starting stats.
func NewEnemy(kind EnemyKind, pos Position) *EnemyRuntime {
	bp := blueprintFor(kind)
	return &EnemyRuntime{
		img:        enemyImage(kind, bp.SpritePath),
		bp:         bp,
		State:      EnemyPatrol,
		PrevState:  EnemyPatrol,
		Animations: bp.NewAnimations(),
		Scale:      bp.Scale,
		Pos:        pos,
		Health:     bp.Stats.MaxHealth,
		spawnX:     pos.X,
		patrolDir:  1,
	}
}

// GetBounds returns the enemy's collision box, satisfying Collider so enemies can
// share the quadtree and AABB collision paths with tiles and the player.
func (e *EnemyRuntime) GetBounds() AABB {
	return AABB{X: e.Pos.X, Y: e.Pos.Y, Width: e.bp.BodyWidth, Height: e.bp.BodyHeight}
}

// Removable reports whether the enemy has finished dying and can be culled from
// the world by the caller.
func (e *EnemyRuntime) Removable() bool { return e.deathDone }

// setState switches the enemy's behavioural state. PrevState is captured once per
// tick in UpdateEnemy, so the animation code can detect the change and restart.
func (e *EnemyRuntime) setState(s EnemyState) { e.State = s }

// faceToward orients the sprite along dx (positive = right). FlipX true means the
// enemy faces left, matching the player's convention.
func (e *EnemyRuntime) faceToward(dx float64) { e.FlipX = dx < 0 }

// TakeDamage applies dmg to the enemy, transitioning to hurt or, if the blow is
// lethal, to the death state. It is a no-op once the enemy is dead.
func (e *EnemyRuntime) TakeDamage(dmg int) {
	if e.State.IsDead() {
		return
	}
	e.Health -= dmg
	if e.Health <= 0 {
		e.Health = 0
		e.VelX = 0
		e.setState(EnemyDead)
		return
	}
	e.setState(EnemyHurt)
}
