---
kind: feature
scope: [cmd/main.go, internal/core/menu*.go, internal/core/save*.go, internal/core/settings*.go, internal/core/text*.go, internal/system/save*.go]
status: done
---
# Main menu: docs index

A title screen (New Game, Continue, Settings, Quit) and a save that survives restarts, so
Continue brings the player back to the last saved portal. Settings holds the sound volume.
Up: [project docs index](../../README.md) · Code: `cmd/main.go`, `internal/core/menu*.go`, `save*.go`, `internal/system/save*.go`

## Read in this order (stop when you have enough)

| Need | File |
|---|---|
| Where work stopped | [state.md](./state.md) |
| What's done and what's next | [tasks.md](./tasks.md) |
| What it must do (acceptance criteria) | [spec.md](./spec.md) |
| Why it's built this way | [decisions.md](./decisions.md) |

## Entry points (where to start reading code; not a full map)

| Path | Why start here |
|---|---|
| `cmd/main.go` | `Game.state` picks menu or play each tick |
| `internal/core/gamestate.go` | the modes |
| `internal/core/combat.go` | `activatePortal`: what a save remembers |

## Depends on / used by

- Uses: portals ([index](../portals/README.md)) for the checkpoint and active portals that a save stores
