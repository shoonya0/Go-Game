package core

// enemyBlueprints is the registry of enemy species. To add an enemy: give it an
// EnemyKind (enemy.go), drop its sprite sheets in assets/enemy/, and add one entry
// here describing its box, stats, clips (each names its sheet) and level-grid Marker (a unique
// ASCII punctuation character, written bare in assets/levels/deepForest/level_1.json). No other
// code changes.
var enemyBlueprints = map[EnemyKind]*EnemyBlueprint{
	KindSuccubus: &succubusBlueprint,
}

// blueprintFor returns the blueprint for a kind, panicking on an unregistered
// kind — that is a programming error (a kind constant without a blueprint), not a
// runtime condition to handle.
func blueprintFor(kind EnemyKind) *EnemyBlueprint {
	bp, ok := enemyBlueprints[kind]
	if !ok {
		panic("core: no blueprint registered for enemy kind")
	}
	return bp
}

// ---------------- succubus ----------------

// The succubus sheet (assets/enemy/succubus.png) is a clean, uniform grid built
// by repacking the original annotated sheet: one animation per row, every frame
// in a fixed cell, content centered horizontally and bottom-aligned.
//
//	row 0: idle   (4)   row 3: attack (6)
//	row 1: stand + patrol (1+7)   row 4: hurt (2)
//	row 2: alert  (3)   row 5: died   (6)
const (
	succubusSheet       = "../assets/enemy/succubus.png"
	succubusFrameWidth  = 270 // cell width
	succubusFrameHeight = 360 // cell height
)

var succubusBlueprint = EnemyBlueprint{
	Name:       "Succubus",
	Marker:     "!",
	BodyWidth:  46,
	BodyHeight: 100,
	Scale:      0.3, // 360px cell -> ~108px tall on screen
	Stats: EnemyStats{
		MaxHealth:      40,
		Damage:         8,
		WalkSpeed:      60,
		ChaseSpeed:     150,
		DetectRange:    420,
		AttackRange:    72,
		AttackCooldown: 1.0,
		PatrolTiles:    4,
	},
	Clips: map[EnemyState]EnemyClip{
		EnemyIdle:   succubusClip(0, 0, 4, 6, true),
		EnemyPatrol: succubusClip(1, 1, 7, 10, true), // column 0 is a front-facing stand, not a step
		EnemyAlert:  succubusClip(2, 0, 3, 8, true),
		EnemyAttack: succubusClip(3, 0, 6, 12, false),
		EnemyHurt:   succubusClip(4, 0, 2, 8, false),
		EnemyDead:   succubusClip(5, 0, 6, 8, false),
	},
}

// succubusClip is a clip on one row of the succubus grid.
func succubusClip(row, start, frames int, fps float64, loop bool) EnemyClip {
	return EnemyClip{
		Sheet: succubusSheet, Row: row, Start: start, Frames: frames,
		FrameW: succubusFrameWidth, FrameH: succubusFrameHeight, FPS: fps, Loop: loop,
	}
}
