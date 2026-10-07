package core

// enemyBlueprints is the registry of enemy species. To add an enemy: give it an
// EnemyKind (enemy.go), drop its sprite sheet in assets/enemy/, and add one entry
// here describing its box, stats and frame layout. No other code changes.
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
// in a fixed cell, content centered horizontally and bottom-aligned. So each
// clip starts at column 0 (AnimStartFrame 0) and its row index is its
// SpriteSheetYPosition; OriginHeight equals the cell height.
//
//	row 0: idle   (4)   row 3: attack (6)
//	row 1: patrol (8)   row 4: hurt   (2)
//	row 2: alert  (3)   row 5: died   (6)
const (
	succubusSheetRow    = 360 // cell height = row unit
	succubusFrameWidth  = 270 // cell width
	succubusFrameHeight = 360 // cell height
)

var succubusBlueprint = EnemyBlueprint{
	Name:         "Succubus",
	SpritePath:   "../assets/enemy/succubus.png",
	OriginHeight: succubusSheetRow,
	BodyWidth:    46,
	BodyHeight:   100,
	Scale:        0.3, // 360px cell -> ~108px tall on screen
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
	NewAnimations: newSuccubusAnimations,
}

// newSuccubusAnimations builds a fresh animation table keyed by EnemyState. Each
// enemy instance gets its own table because Animation.FrameTimer is mutable.
func newSuccubusAnimations() map[int]*Animation {
	frame := func(row, start, total int, speed float64, loop bool) *Animation {
		return &Animation{
			SpriteSheetYPosition: row,
			AnimStartFrame:       start,
			TotalFrames:          total,
			FrameWidth:           succubusFrameWidth,
			FrameHeight:          succubusFrameHeight,
			AnimationSpeed:       speed,
			Looping:              loop,
		}
	}

	return map[int]*Animation{
		int(EnemyIdle):   frame(0, 0, 4, 6, true),   // row 0: idle
		int(EnemyPatrol): frame(1, 1, 8, 10, true),  // row 1: patrol walk cycle
		int(EnemyAlert):  frame(2, 0, 3, 8, true),   // row 2: alert
		int(EnemyAttack): frame(3, 0, 6, 12, false), // row 3: attack (one-shot)
		int(EnemyHurt):   frame(4, 0, 2, 8, false),  // row 4: hurt (one-shot)
		int(EnemyDead):   frame(5, 0, 6, 8, false),  // row 5: died (one-shot)
	}
}
