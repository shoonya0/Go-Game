package core

// enemyBlueprints is the registry of enemy species. To add an enemy: give it an
// EnemyKind (enemy.go), drop its sprite sheets in assets/enemy/, and add one entry
// here describing its box, stats, clips (each names its sheet) and level-grid Marker (a unique
// ASCII punctuation character, written bare in assets/levels/deepForest/level_1.json). No other
// code changes.
var enemyBlueprints = map[EnemyKind]*EnemyBlueprint{
	KindSuccubus:       &succubusBlueprint,
	KindSkeletonWhite:  skeletonBlueprint("White", "@"),
	KindSkeletonYellow: skeletonBlueprint("Yellow", "#"),
	KindMushroom:       &mushroomBlueprint,
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
	Clips: map[ClipName]EnemyClip{ // no run clip: chasing falls back to alert
		ClipIdle:    succubusClip(0, 0, 4, 6, true),
		ClipWalk:    succubusClip(1, 1, 7, 10, true),                                                 // column 0 is a front-facing stand, not a step
		ClipAlert:   succubusClip(2, 0, 3, 8, true).moving(motionFX{BobPx: 3, BobHz: 3, LeanDeg: 5}), // chase: a standing pose, so lean in
		ClipAttack1: succubusClip(3, 0, 6, 12, false).held(1, 1, 1, 2),                               // frame 3 is the big slash
		ClipHurt:    succubusClip(4, 0, 2, 8, false).moving(motionFX{HitShakePx: 4}),
		ClipDeath:   succubusClip(5, 0, 6, 8, false),
	},
}

// succubusClip is a clip on one row of the succubus grid.
func succubusClip(row, start, frames int, fps float64, loop bool) EnemyClip {
	return EnemyClip{
		Sheet: succubusSheet, Row: row, Start: start, Frames: frames,
		FrameW: succubusFrameWidth, FrameH: succubusFrameHeight, FPS: fps, Loop: loop,
	}
}

// ---------------- skeleton ----------------

// The skeleton ships as one 96×64 strip per animation, facing right, in two
// colours with the same frames: assets/enemy/Skeleton_<Colour>/Skeleton_01_<Colour>_<Anim>.png.
// The body is about 46px tall in the cell, feet on the bottom row.
const (
	skeletonFrameWidth  = 96
	skeletonFrameHeight = 64
)

// skeletonBlueprint is the skeleton in one colour. The colours differ only in
// their sheets, name and marker.
func skeletonBlueprint(colour, marker string) *EnemyBlueprint {
	clip := func(anim string, frames int, fps float64, loop bool) EnemyClip {
		return EnemyClip{
			Sheet:  "../assets/enemy/Skeleton_" + colour + "/Skeleton_01_" + colour + "_" + anim + ".png",
			Frames: frames, FrameW: skeletonFrameWidth, FrameH: skeletonFrameHeight, FPS: fps, Loop: loop,
		}
	}
	return &EnemyBlueprint{
		Name:       "Skeleton " + colour,
		Marker:     marker,
		BodyWidth:  30,
		BodyHeight: 88,
		Scale:      2, // 46px of art -> ~92px on screen
		Stats: EnemyStats{
			MaxHealth:      40,
			Damage:         10,
			WalkSpeed:      50,
			ChaseSpeed:     110,
			DetectRange:    380,
			AttackRange:    70,
			AttackCooldown: 1.2,
			PatrolTiles:    2.5,
		},
		Clips: map[ClipName]EnemyClip{ // no run clip: chasing walks faster
			ClipIdle:    clip("Idle", 8, 8, true),
			ClipWalk:    clip("Walk", 10, 12, true),
			ClipAttack1: clip("Attack1", 10, 14, false).held(1, 1, 1, 1, 1, 2, 1.5), // frames 5-6: thrust at full reach
			ClipAttack2: clip("Attack2", 9, 14, false).held(1, 1, 1, 1, 1.5, 2),     // frames 4-5: second thrust
			ClipHurt:    clip("Hurt", 5, 12, false).moving(motionFX{HitShakePx: 3}),
			ClipDeath:   clip("Die", 13, 12, false),
		},
	}
}

// ---------------- mushroom ----------------

// The mushroom ships as one 80×64 strip per animation, facing left (it is
// knocked back to the right when hit and slams to the left). The body is about
// 33px tall, feet on the bottom row. Its stun clips wait for a player stun state.
const (
	mushroomFrameWidth  = 80
	mushroomFrameHeight = 64
)

var mushroomBlueprint = EnemyBlueprint{
	Name:       "Mushroom",
	Marker:     "&",
	BodyWidth:  44,
	BodyHeight: 64,
	Scale:      2, // 33px of art -> ~66px on screen
	FacesLeft:  true,
	Stats: EnemyStats{
		MaxHealth:      40,
		Damage:         8,
		WalkSpeed:      40,
		ChaseSpeed:     130,
		DetectRange:    320,
		AttackRange:    56,
		AttackCooldown: 1.0,
		PatrolTiles:    2.5,
	},
	Clips: map[ClipName]EnemyClip{ // no walk art: patrol plays the run cycle slower
		ClipIdle:    mushroomClip("Idle", 7, 8, true),
		ClipWalk:    mushroomClip("Run", 8, 8, true),
		ClipRun:     mushroomClip("Run", 8, 14, true).moving(motionFX{LeanDeg: 4}),
		ClipAttack1: mushroomClip("Attack", 10, 14, false).held(1, 1, 1, 1, 1, 2.5), // frame 5: the slam lands
		ClipHurt:    mushroomClip("Hit", 5, 12, false).moving(motionFX{HitShakePx: 3}),
		ClipDeath:   mushroomClip("Die", 15, 12, false),
	},
}

// mushroomClip is one of the mushroom's strips.
func mushroomClip(anim string, frames int, fps float64, loop bool) EnemyClip {
	return EnemyClip{
		Sheet:  "../assets/enemy/Mushroom/Mushroom-" + anim + ".png",
		Frames: frames, FrameW: mushroomFrameWidth, FrameH: mushroomFrameHeight, FPS: fps, Loop: loop,
	}
}
