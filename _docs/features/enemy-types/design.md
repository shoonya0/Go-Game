# Enemy types: design (planned, not in the code yet)

A proposal for increment 1. Once the code exists, the code wins: shrink this file to what
the code doesn't show (the reasons stay in [decisions.md](./decisions.md)).

## Shape

```text
EnemyBlueprint (data, one per species)
  ├─ Clips    map[EnemyState]EnemyClip // sheet path + frame size + count + speed + loop
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

**Mover: done (task 5),** `internal/core/enemy_move.go`. It differs from the first plan: it has
`Steer` (intent → velocity, called by the brain as it decides) and `Move` (integrate + collide)
instead of a single `Move(e, intent, qt)`. That way the brain's tests can still check
velocity right after `decideBehavior`, and a dead enemy just steers to `moveHold` and keeps
falling.

```go
// Attacker decides when the enemy can strike and what a strike does.
type Attacker interface {
	InRange(e *EnemyRuntime, target AABB) bool
	Update(e *EnemyRuntime, target *PlayerRuntime, fx *EnemyEffects)
	Box(e *EnemyRuntime) AABB // for the F3 overlay
}
```

| Part | Today's code it comes from | Notes |
|---|---|---|
| `groundMover` | done (task 5) | ledge probe still to add (task 12) |
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

**Done (tasks 2–4):** `EnemyClip` (`internal/core/enemy.go`), sheets cached by path,
per-frame `Hold`, and draw-time motion (D7, `internal/core/enemy_motion.go`). Motion values
are screen pixels and move only the sprite; the collision box stays put. Standard clip names
and their fallbacks are listed in [spec.md](./spec.md).

## Battler baking (offline, `tools/battler2strip`)

Input: a sheet plus a mapping file per enemy (JSON next to the tool, not embedded):
`{"cellCols":9,"cellRows":6,"facesLeft":true,"drawHeight":120,"clips":{"attack_1":[[1,0,0],[1,0,1],…]}}`
where each entry is `[motionBlock, row, frame]`. Output: `assets/enemy/<Enemy>/<Enemy>_<clip>.png`
plus a printed Go snippet with frame sizes and counts, to paste into the blueprint. Cell
size = sheet size ÷ 9 and ÷ 6 as floats, each cell's rect rounded, so non-integer heights work.

## Projectiles

`EnemyEffects` holds live projectiles. `Game` owns it, `updateEnemies` advances it, and
`Draw` draws it after enemies. Same pattern as the player's ground waves (`ground_wave.go`).
This is the only `cmd/main.go` change the feature needs.

## New states

- `EnemyGuard` (D8): a one-shot entered from `TakeDamage` when the blow comes from the
  facing side and a `GuardChance` roll passes. It deals no damage and hands back like hurt.
  The roll goes through a seeded `rand.Rand` held by the enemy, so tests are deterministic.
- `EnemyRun` is not a state: chase plays `run` and patrol plays `walk`. Clip choice is a
  function of state plus speed.
- `EnemySleep` / `EnemyWake` only if a sleeping enemy is added (none in the current roster).
