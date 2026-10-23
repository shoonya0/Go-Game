package core

import "github.com/hajimehoshi/ebiten/v2"

// EnemyKind identifies a concrete enemy species. Adding a new enemy is a matter
// of adding a kind here and a blueprint in enemyBlueprints (enemy_types.go), then
// placing its Marker in the level grid — no
// changes to the shared update, physics, animation or draw code are needed.
type EnemyKind int

const (
	KindSuccubus EnemyKind = iota
	KindSkeletonWhite
	KindSkeletonYellow
	KindMushroom
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

// EnemyClip describes one animation of a species: which sprite sheet holds it
// and where its frames sit. Each clip names its own sheet, so a species whose art
// ships as one PNG strip per animation and a species packed into one grid are
// described the same way.
type EnemyClip struct {
	Sheet          string  // asset path of the sheet holding this clip
	Row            int     // sheet row, in FrameH units (0 for a one-row strip)
	Start, Frames  int     // first column and number of frames
	FrameW, FrameH int     // size of one frame cell
	FPS            float64 // frames shown per second
	Loop           bool
	// Hold stretches single frames: frame i of the clip shows for Hold[i] times
	// as long as the others, so a strike or impact frame stays up. Missing or
	// non-positive entries mean 1.
	Hold     []float64
	Motion   motionFX // movement added at draw time (enemy_motion.go)
	HitFrame int      // attack clips: the frame (from 0) on which the strike lands
}

// EnemyBlueprint is the immutable description of an enemy species: its clips,
// its collision box, and its stats. This is the "different sprites take on the
// general enemy" part — every field that varies per species is data here, while
// behaviour is shared code.
type EnemyBlueprint struct {
	Name       string
	Marker     string // one-character level-grid code that spawns this kind (e.g. !)
	BodyWidth  float64
	BodyHeight float64
	Scale      float64  // render scale of the sprite
	FacesLeft  bool     // the art faces left, so drawing mirrors it the other way
	Mover      Mover    // how it gets around (enemy_move.go)
	Attacker   Attacker // how it fights (enemy_attack.go)
	Stats      EnemyStats
	Clips      map[ClipName]EnemyClip // must hold requiredClips; others fall back (enemy_clips.go)
}

// newAnimations builds a fresh animation table for one instance. Each instance
// needs its own because Animation carries a mutable FrameTimer; sharing one table
// would couple the playback of every enemy of the kind.
func (bp *EnemyBlueprint) newAnimations() map[ClipName]*Animation {
	anims := make(map[ClipName]*Animation, len(bp.Clips))
	for name, c := range bp.Clips {
		anims[name] = &Animation{
			SpriteSheetYPosition: c.Row,
			AnimStartFrame:       c.Start,
			TotalFrames:          c.Frames,
			FrameWidth:           c.FrameW,
			FrameHeight:          c.FrameH,
			AnimationSpeed:       c.FPS,
			Looping:              c.Loop,
		}
	}
	return anims
}

// EnemyRuntime is the complete, mutable state of one enemy instance. It mirrors
// the shape of PlayerRuntime (sprite + spatial + combat + per-tick bookkeeping)
// but only carries what an AI-driven enemy needs.
type EnemyRuntime struct {
	bp *EnemyBlueprint // the species blueprint (shared, read-only)

	// State and animation.
	State         EnemyState
	PrevState     EnemyState
	Animations    map[ClipName]*Animation
	CurrAnimFrame int
	clipTime      float64 // seconds since the current clip (re)started; drives Motion

	// Spatial state.
	FlipX      bool // true when the enemy faces left
	Scale      float64
	Pos        Position
	VelX, VelY float64

	// Combat.
	Health int

	// AI / per-tick bookkeeping (not for external use).
	spawnX      float64    // patrol anchor (X at spawn)
	patrolDir   float64    // current patrol direction: -1 or +1
	attackTimer float64    // seconds until the next hit may land
	swing       swingState // the attack in progress
	regenAccum  float64    // sub-point health regenerated but not yet added
	onGround    bool       // whether the enemy rested on ground last tick
	deathDone   bool       // whether the death animation has finished
}

// enemySheets caches one GPU image per sheet path, so spawning many enemies of
// the same kind (or clips sharing a sheet) does not re-upload it. The game loop
// is single-threaded, so no synchronisation is required.
var enemySheets = map[string]*ebiten.Image{}

// enemySheet returns the (cached) image for a sheet path.
func enemySheet(path string) *ebiten.Image {
	if img, ok := enemySheets[path]; ok {
		return img
	}
	img := LoadImage(path)
	enemySheets[path] = img
	return img
}

// NewEnemy builds a ready-to-run enemy of the given kind at pos, wiring its
// blueprint, animations and starting stats. It loads every clip's sheet now, so
// a missing asset fails at spawn rather than mid-fight.
func NewEnemy(kind EnemyKind, pos Position) *EnemyRuntime {
	bp := blueprintFor(kind)
	for _, c := range bp.Clips {
		enemySheet(c.Sheet)
	}
	return &EnemyRuntime{
		bp:         bp,
		State:      EnemyPatrol,
		PrevState:  EnemyPatrol,
		Animations: bp.newAnimations(),
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

// AttackBox returns the area the enemy's strike covers in the direction it
// faces. The player is hit only while their body overlaps it.
func (e *EnemyRuntime) AttackBox() AABB { return e.bp.Attacker.Box(e) }

// Hitboxes reports the enemy's body and, while alive, its attack box to the
// debug overlay (HitboxProvider).
func (e *EnemyRuntime) Hitboxes() []Hitbox {
	boxes := []Hitbox{{Box: e.GetBounds(), Kind: HitboxBody}}
	if !e.State.IsDead() {
		boxes = append(boxes, Hitbox{Box: e.AttackBox(), Kind: HitboxAttack, Active: e.State == EnemyAttack})
	}
	return boxes
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
	} else {
		e.setState(EnemyHurt)
	}
	// Hits land outside UpdateEnemy, where PrevState is captured, so the
	// animation can't see the change: restart the clip here. This also replays
	// hurt when an already-hurt enemy is hit again.
	e.restartAnimation()
}

// restartAnimation rewinds the current state's clip to its first frame.
func (e *EnemyRuntime) restartAnimation() {
	if anim := e.currentAnim(); anim != nil {
		e.CurrAnimFrame = anim.AnimStartFrame
		anim.FrameTimer = 0
	}
	e.clipTime = 0
}
