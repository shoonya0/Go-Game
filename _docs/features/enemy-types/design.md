# Enemy types: design (planned, not in the code yet)

A proposal for increment 1. Once the code exists, the code wins: shrink this file to what
the code doesn't show (the reasons stay in [decisions.md](./decisions.md)).

## Shape

```text
EnemyBlueprint (data, one per species)
  ├─ Clips    map[EnemyState]ClipDef   // sheet path + frame size + count + speed + loop
  ├─ Mover    Mover                    // ground | flying | (water)
  ├─ Attacker Attacker                 // melee | ranged | dive
  ├─ Stats    EnemyStats               // + per-part tuning lives in the part's own struct
  └─ Traits   StartState, AttackClips […], Variants

UpdateEnemy (shared, unchanged order)
  brain:  decideBehavior → intent {Idle | Patrol | Chase(dx,dy) | Attack}
  move:   bp.Mover.Move(e, intent, qt)        // velocity, gravity or not, tile collision
  attack: bp.Attacker.Update(e, target, fx)   // range test, damage timing, projectiles
```

The brain stays shared: patrol → notice → chase → attack, hurt interrupts, dead is final.
`Attacker.InRange` replaces today's attack-box test so a ranged enemy "attacks" from far
away and a diver attacks from above.

## Parts (Go interfaces, implemented as small structs holding their tuning)

```go
// Mover turns the brain's intent into movement and resolves collisions.
type Mover interface {
	Move(e *EnemyRuntime, in moveIntent, qt *DynamicQuadtree)
	Grounded(e *EnemyRuntime) bool // ground wave and landing logic ask this
}

// Attacker decides when the enemy can strike and what a strike does.
type Attacker interface {
	InRange(e *EnemyRuntime, target AABB) bool
	Update(e *EnemyRuntime, target *PlayerRuntime, fx *EnemyEffects)
	Box(e *EnemyRuntime) AABB // for the F3 overlay
}
```

| Part | Today's code it comes from | Notes |
|---|---|---|
| `groundMover` | `applyGravity`, `resolveHorizontal/Vertical`, `patrol` | adds a ledge probe |
| `flyingMover` | new | no gravity, `HoverHeight`, `Leash` (tiles), steers in X and Y, bobbing optional |
| `waterMover` | later | like flying but clamped to water tiles |
| `meleeAttack` | `attack`, `attackBoxFacing` | `HitFrame` lands damage on a clip frame instead of on entry |
| `rangedAttack` | new | spawns a `Projectile` on `FireFrame`; `KeepDistance` |
| `diveAttack` | new | phases start → loop (falling) → end; uses up to 3 clips |

## Per-instance state

Part-specific runtime state (dive phase, fire timer, wake flag) goes in one small field on
`EnemyRuntime` per part, e.g. `dive diveState`, not in a generic `map[string]any`. That keeps it
typed, and the zero value means "idle".

## Clips and sheets

`ClipDef{Sheet string; FrameW, FrameH, Row, Start, Frames int; FPS float64; Loop bool}`.
Sheets load through one cache keyed by path (replaces `enemyImages` keyed by kind). A
missing optional clip falls back: `Alert → Patrol → Idle`. `drawSpriteFrame` already takes
the image and frame size per call; only the Y origin (`Row * FrameH`) changes source.

## Projectiles

`EnemyEffects` holds live projectiles. `Game` owns it, `updateEnemies` advances it, and
`Draw` draws it after enemies. Same pattern as the player's ground waves (`ground_wave.go`).
This is the only `cmd/main.go` change the feature needs.

## New states

Add `EnemySleep` and `EnemyWake` only when the Bat lands (increment 4). Keep them out of
`IsOneShot` except `EnemyWake`, which hands over to `EnemyAlert`.
