package core

// enemyBlueprints is the registry of enemy species. To add an enemy: give it an
// EnemyKind (enemy.go), drop its sprite sheets in assets/enemy/, and add one entry
// here describing its box, stats, clips (each names its sheet) and level-grid Marker (1-2 unique
// ASCII punctuation characters, written bare in assets/levels/deepForest/level_1.json). No other
// code changes.
var enemyBlueprints = map[EnemyKind]*EnemyBlueprint{
	KindSkeletonWhite:  skeletonBlueprint("White", "@"),
	KindSkeletonYellow: skeletonBlueprint("Yellow", "#"),
	KindMushroom:       &mushroomBlueprint,
	KindWhiteGolem:     &whiteGolemBlueprint,
	KindFlyingDemon:    &flyingDemonBlueprint,
	KindSmasher:        &smasherBlueprint,
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
		Mover:      groundMover{},
		Attacker:   meleeAttack{},
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
			ClipAttack1: clip("Attack1", 10, 14, false).held(1, 1, 1, 1, 1, 2, 1.5).hitsOn(5), // frames 5-6: thrust at full reach
			ClipAttack2: clip("Attack2", 9, 14, false).held(1, 1, 1, 1, 1.5, 2).hitsOn(4),     // frames 4-5: second thrust
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
	Mover:      groundMover{},
	Attacker:   meleeAttack{},
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
		ClipAttack1: mushroomClip("Attack", 10, 14, false).held(1, 1, 1, 1, 1, 2.5).hitsOn(5), // frame 5: the slam lands
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

// ---------------- white golem ----------------

// The white golem ships as one 90×64 strip per animation, facing right:
// assets/enemy/White_Golem/WhiteGolem_<anim>.png (renamed from Golem_1_*, which
// the orange golem also uses, D4). The body is about 38px tall, feet on the
// bottom row. A slow, heavy enemy: it winds up for three frames and slams the
// ground on frame 6, raising ice spikes on both sides.
const (
	whiteGolemFrameWidth  = 90
	whiteGolemFrameHeight = 64
)

var whiteGolemBlueprint = EnemyBlueprint{
	Name:       "White Golem",
	Marker:     "%",
	BodyWidth:  56,
	BodyHeight: 74,
	Scale:      2, // 38px of art -> ~76px on screen
	Mover:      groundMover{},
	Attacker:   meleeAttack{},
	Stats: EnemyStats{
		MaxHealth:      80,
		Damage:         15,
		WalkSpeed:      30,
		ChaseSpeed:     70,
		DetectRange:    300,
		AttackRange:    64,
		AttackCooldown: 1.8,
		PatrolTiles:    2,
	},
	Clips: map[ClipName]EnemyClip{ // no run clip: chasing walks faster
		ClipIdle:    whiteGolemClip("idle", 8, 8, true),
		ClipWalk:    whiteGolemClip("walk", 10, 10, true),
		ClipAttack1: whiteGolemClip("attack", 11, 12, false).held(1, 1, 1, 1, 1, 2, 2).hitsOn(6), // frame 5: arm up; frame 6: the slam lands
		ClipHurt:    whiteGolemClip("hurt", 4, 10, false).moving(motionFX{HitShakePx: 2}),
		ClipDeath:   whiteGolemClip("die", 12, 12, false), // the 13th frame is blank: end on the rubble
	},
}

// whiteGolemClip is one of the white golem's strips.
func whiteGolemClip(anim string, frames int, fps float64, loop bool) EnemyClip {
	return EnemyClip{
		Sheet:  "../assets/enemy/White_Golem/WhiteGolem_" + anim + ".png",
		Frames: frames, FrameW: whiteGolemFrameWidth, FrameH: whiteGolemFrameHeight, FPS: fps, Loop: loop,
	}
}

// ---------------- flying demon ----------------

// The flying demon ships as one 81×71 strip per animation, facing left:
// assets/enemy/flying_Demon/FlyingDemon_<anim>.png (renamed from IDLE.png etc.,
// D4). It flies (flyingMover) and bites: the mouth opens on attack frame 3.
// Its art has no walk, so patrol and chase play the wing-beat loop.
const (
	flyingDemonFrameWidth  = 81
	flyingDemonFrameHeight = 71
)

var flyingDemonBlueprint = EnemyBlueprint{
	Name:       "Flying Demon",
	Marker:     "^",
	BodyWidth:  60,
	BodyHeight: 90,
	Scale:      1.5, // 64px of art -> ~96px on screen
	FacesLeft:  true,
	// HoverHeight 20 keeps its feet off the floor while it bites a standing player.
	Mover:    flyingMover{LeashTiles: 4, HoverHeight: 20},
	Attacker: meleeAttack{},
	Stats: EnemyStats{
		MaxHealth:      40,
		Damage:         8,
		WalkSpeed:      50,
		ChaseSpeed:     120,
		DetectRange:    360,
		AttackRange:    50,
		AttackCooldown: 1.2,
		PatrolTiles:    3,
	},
	Clips: map[ClipName]EnemyClip{
		ClipIdle:    flyingDemonClip("idle", 4, 8, true).moving(motionFX{BobPx: 3, BobHz: 2}),
		ClipWalk:    flyingDemonClip("fly", 4, 8, true).moving(motionFX{BobPx: 3, BobHz: 2}),
		ClipRun:     flyingDemonClip("fly", 4, 12, true).moving(motionFX{BobPx: 2, BobHz: 3, LeanDeg: 4}),
		ClipAttack1: flyingDemonClip("attack", 8, 12, false).held(1, 1, 1, 2).hitsOn(3), // frame 3: the bite
		ClipHurt:    flyingDemonClip("hurt", 4, 10, false).moving(motionFX{HitShakePx: 2}),
		ClipDeath:   flyingDemonClip("death", 6, 10, false), // the 7th frame is blank: end on the last scraps
	},
}

// flyingDemonClip is one of the flying demon's strips.
func flyingDemonClip(anim string, frames int, fps float64, loop bool) EnemyClip {
	return EnemyClip{
		Sheet:  "../assets/enemy/flying_Demon/FlyingDemon_" + anim + ".png",
		Frames: frames, FrameW: flyingDemonFrameWidth, FrameH: flyingDemonFrameHeight, FPS: fps, Loop: loop,
	}
}

// ---------------- smasher ----------------

// The smasher (Enemy3) ships as one 64×64 strip per animation, facing left:
// assets/enemy/Enemy3-Movement-In-Animation/Enemy3-<anim>.png. It flies, hangs
// over the player and dives onto them (diveAttack): the smash start is the
// wind-up, the smash loop the fall, and the smash end the impact, squashed on
// its frame 0.
const smasherFrameSize = 64

var smasherBlueprint = EnemyBlueprint{
	Name:       "Smasher",
	Marker:     "*",
	BodyWidth:  56,
	BodyHeight: 80,
	Scale:      1.5, // 52px of art -> ~78px on screen
	FacesLeft:  true,
	Mover:      flyingMover{LeashTiles: 5, HoverHeight: 160},
	Attacker:   diveAttack{FallSpeed: 700, MaxDrop: 300},
	Stats: EnemyStats{
		MaxHealth:      60,
		Damage:         20,
		WalkSpeed:      50,
		ChaseSpeed:     110,
		DetectRange:    360,
		AttackRange:    50,
		AttackCooldown: 2.5,
		PatrolTiles:    2,
	},
	Clips: map[ClipName]EnemyClip{
		ClipIdle:       smasherClip("Idle", 8, 8, true),
		ClipWalk:       smasherClip("Fly", 8, 10, true),
		ClipRun:        smasherClip("Fly", 8, 14, true),
		ClipAttack1:    smasherClip("AttackSmashStart", 12, 12, true), // wind-up: the dive ends it after one play
		ClipAttackLoop: smasherClip("AttackSmashLoop", 3, 12, true),
		ClipAttackEnd:  smasherClip("SmashEnd", 8, 10, false),
		ClipHurt:       smasherClip("Hit", 4, 10, false).moving(motionFX{HitShakePx: 2}),
		ClipDeath:      smasherClip("Die", 16, 12, false), // the 17th frame is blank: end on the last wisp
	},
}

// smasherClip is one of the smasher's strips.
func smasherClip(anim string, frames int, fps float64, loop bool) EnemyClip {
	return EnemyClip{
		Sheet:  "../assets/enemy/Enemy3-Movement-In-Animation/Enemy3-" + anim + ".png",
		Frames: frames, FrameW: smasherFrameSize, FrameH: smasherFrameSize, FPS: fps, Loop: loop,
	}
}
