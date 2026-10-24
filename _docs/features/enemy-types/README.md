---
kind: feature
scope: [internal/core/enemy*.go, internal/core/level_enemies*.go, assets/enemy, tools/battler2strip]
status: active
---
# Enemy types: docs index

Adds the enemy sprites in `assets/enemy/` as playable enemies, grouped by how they move
(ground, flying, later water) and how they attack (melee, ranged, dive), so new enemies are
mostly data. Bosses (`Striker/`, `Orange_Golem/`) are out of scope for now.
Up: [project docs index](../../README.md) · Code: `internal/core/enemy*.go`, `internal/core/level_enemies.go`,
`tools/battler2strip/` (offline baking)

## Read in this order (stop when you have enough)

| Need | File |
|---|---|
| Where work stopped | [state.md](./state.md) |
| What's done and what's next | [tasks.md](./tasks.md) |
| What it must do, and the enemy roster (sheets, frame sizes, type) | [spec.md](./spec.md) |
| Why it's built this way | [decisions.md](./decisions.md) |
| Why the parts are shaped as they are, and the parts still to build | [design.md](./design.md) |

## Entry points (where to start reading code; not a full map)

| Path | Why start here |
|---|---|
| `internal/core/enemy.go` | `EnemyKind`, `EnemyBlueprint`, `EnemyClip`, `EnemyRuntime`: what a species is |
| `internal/core/enemy_types.go` | Blueprint registry: skeleton (white/yellow), mushroom, White Golem |
| `internal/core/enemy_ai.go` | Shared brain: patrol → chase → attack, driving the parts below |
| `internal/core/enemy_move.go`, `enemy_attack.go` | `Mover` (`groundMover`) and `Attacker` (`meleeAttack`) |
| `internal/core/enemy_clips.go`, `enemy_anim.go`, `enemy_motion.go` | Clip names + fallbacks, playback, draw-time motion |
| `internal/core/level_enemies.go` | Level-grid markers → spawns |
| `tools/battler2strip/` | Bakes a battler sheet + `mappings/*.json` into strips (offline) |

## Depends on / used by

- Uses: player combat (`HitEnemies`, `ApplyDamageFrom`), the quadtree, `drawSpriteFrame`, `readAsset`
- Used by: `cmd/main.go` (`updateEnemies`, draw loop, hitbox overlay)
