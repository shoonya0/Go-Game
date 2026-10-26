---
kind: feature
scope: [internal/core/levels*.go, internal/core/world*.go, assets/levels]
status: active
---
# Levels: docs index

Eight themed levels, each with its own background (day and night), its own set of tiles and
its own enemies. deepForest is today's level 1; the other seven are new.
Up: [project docs index](../../README.md) · Code: `internal/core/levels.go`, `internal/core/world.go`, `assets/levels/`

## Read in this order (stop when you have enough)

| Need | File |
|---|---|
| Where work stopped | [state.md](./state.md) |
| What's done and what's next | [tasks.md](./tasks.md) |
| What each level must use (tiles, background, enemies) | [spec.md](./spec.md) |

## Entry points (where to start reading code; not a full map)

| Path | Why start here |
|---|---|
| `internal/core/levels.go` | `LevelDef`, `Levels` (levels by name), `World.SwitchLevel` / `Travel` |
| `internal/core/world.go` | `tileRowOrder` (grid code → tile type), `LoadLevel`, `BuildLevel`, level size constants |
| `internal/core/level_enemies.go` | Enemy markers in the grid and how they spawn |
| `assets/embed.go` | `//go:embed levels` embeds every file under `assets/levels/` |

## Depends on / used by

- Uses: [portals](../portals/README.md) for travel between levels, [enemy-types](../enemy-types/README.md) for the enemies placed in each level
