---
kind: feature
scope: [internal/core/portal*.go, internal/core/world.go, internal/core/levels*.go, internal/core/combat.go, assets/levels/portal.png]
status: active
---
# Portals: docs index

Animated portals (`assets/levels/portal.png`) mark where the player starts, saves their
respawn point and travels between levels. Like an open-world game, non-boss enemies come
back when the player dies or enters a level.
Up: [project docs index](../../README.md) · Code: `internal/core/portal*.go` (planned),
`internal/core/world.go` (level grid, code 9), `internal/core/combat.go` (respawn)

## Read in this order (stop when you have enough)

| Need | File |
|---|---|
| Where work stopped | [state.md](./state.md) |
| What's done and what's next | [tasks.md](./tasks.md) |
| What it must do (acceptance criteria) | [spec.md](./spec.md) |
| Why it's built this way (the user's choices) | [decisions.md](./decisions.md) |

## Entry points (where to start reading code; not a full map)

| Path | Why start here |
|---|---|
| `internal/core/world.go` | `CheckpointCode` (9), `Checkpoint`, `BuildLevel`: today's invisible checkpoints, which become portals |
| `internal/core/levels.go` | `Levels` (levels by name) and `World` (the current level, `SwitchLevel`) |
| `internal/core/combat.go` | `updateCheckpoints` (touch to save, replaced by interact) and `respawn` |
| `internal/core/level_enemies.go` | `EnemySpawn`, `SpawnEnemies`: the spawn list enemies are revived from |
| `internal/core/touch_controls.go` | On-screen buttons, for the tap-to-save on touch screens |

## Depends on / used by

- Uses: [enemy-types](../enemy-types/README.md) for the enemy spawn list and blueprints (a `Boss` flag)
