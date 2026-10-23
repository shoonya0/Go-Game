---
kind: feature
scope: [internal/core/enemy*.go, internal/core/level_enemies*.go, assets/enemy]
status: active
---
# Enemy types: docs index

Adds the enemy sprites in `assets/enemy/` as playable enemies, grouped by how they move
(ground, flying, later water) and how they attack (melee, ranged, dive), so new enemies are
mostly data. Bosses (`Striker/`, `Orange_Golem/`) are out of scope for now.
Up: [project docs index](../../README.md) · Code: `internal/core/enemy*.go`, `internal/core/level_enemies.go`

## Read in this order (stop when you have enough)

| Need | File |
|---|---|
| Where work stopped | [state.md](./state.md) |
| What's done and what's next | [tasks.md](./tasks.md) |
| What it must do, and the enemy roster (sheets, frame sizes, type) | [spec.md](./spec.md) |
| Why it's built this way | [decisions.md](./decisions.md) |
| The planned split into movers and attackers (not in the code yet) | [design.md](./design.md) |

## Entry points (where to start reading code; not a full map)

| Path | Why start here |
|---|---|
| `internal/core/enemy.go` | `EnemyKind`, `EnemyBlueprint`, `EnemyRuntime`: what a species is |
| `internal/core/enemy_types.go` | Blueprint registry; the succubus is the only entry today |
| `internal/core/enemy_ai.go` | Shared state machine plus ground physics; this is what gets split |
| `internal/core/enemy_anim.go` | Clip playback and drawing; one sheet per enemy today |
| `internal/core/level_enemies.go` | Level-grid markers → spawns |

## Depends on / used by

- Uses: player combat (`HitEnemies`, `ApplyDamageFrom`), the quadtree, `drawSpriteFrame`, `readAsset`
- Used by: `cmd/main.go` (`updateEnemies`, draw loop, hitbox overlay)
