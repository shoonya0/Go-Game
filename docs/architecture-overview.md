# Architecture overview

> Snapshot 2026-10-07 at `91400c1` (+ untracked enemy WIP). Map for humans and agents.
> Rows marked **target** describe intended rules the code does not fully enforce yet.

Go 1.24 module `player` · Ebitengine v2.9 · one codebase, two targets: native desktop
(Windows) and WebAssembly (browser), plus a tiny static server for deployment (Render).

## Overview diagram

```text
 cmd/main.go (Game: ebiten.Game)
   │  Update(): system.HandleInput → core.UpdatePlayer → anim/camera → updateEnemies
   │  Draw():   parallax bg → level tiles → enemies → player → water → HUD
   ▼
 internal/system ──► internal/core ──► assets (go:embed FS)
   (keyboard → InputState)   (simulation + rendering)

 server/main.go  — standalone static file server for ./web (no game imports)
 web/index.html  — boots game.wasm via wasm_exec.js
```

## Modules

| Module | Path | Responsibility | Owns state? |
|---|---|---|---|
| Entry / game loop | `cmd/main.go` | Wires subsystems; implements `ebiten.Game` (Update/Draw/Layout) | `Game` struct: player, enemies, level, quadtree |
| Input | `internal/system/Input.go` | Polls Ebitengine keyboard → `core.InputState` | no (writes caller's struct) |
| Player | `internal/core/player.go`, `player_runtime.go`, `playerState.go` | Movement, gravity/jump, AABB collision resolution, state enum | `PlayerRuntime` |
| Combat | `internal/core/combat.go`, `skills.go` | Damage, power regen, hazards (water/larva), fall damage, checkpoints, respawn | inside `PlayerRuntime` |
| Enemies | `internal/core/enemy*.go` (**untracked WIP**) | Blueprint registry, AI (patrol/chase/attack), animation, damage | `EnemyRuntime` |
| Animation | `internal/core/animation.go`, `enemy_anim.go` | Sprite-sheet frame tables, frame advance, drawing with camera | per-runtime |
| World / level | `internal/core/world.go` | Loads embedded `level_1.json`, builds platforms/tiles/checkpoints, draws level & water | `Level` |
| Spatial index | `internal/core/quadtree.go`, `aabb.go` | `DynamicQuadtree` for collision broad-phase; `Collider` interface | quadtree |
| HUD | `internal/core/hud.go` | Health / power bars | no |
| Assets | `assets/embed.go`, `internal/core/image.go` | `go:embed` FS; `readAsset` resolves by **base file name** (subfolders searched) | embedded, read-only |
| Time | `internal/core/time.go` | `deltaTime()` = 1/TPS, shared by physics and animation | no |
| Web server | `server/main.go` | Serves `$WEB_ROOT` on `$PORT`, `.wasm` MIME, no-cache for HTML | no |

## Layers and allowed imports (target, enforced by golangci-lint depguard)

| From \ To | cmd | internal/system | internal/core | assets | ebiten |
|---|---|---|---|---|---|
| **cmd** | — | ✅ | ✅ | ❌ (go through core) | ✅ |
| **internal/system** | ❌ | — | ✅ (types only) | ❌ | ✅ |
| **internal/core** | ❌ | ❌ | — | ✅ | ✅ |
| **server** | ❌ | ❌ | ❌ | ❌ | ❌ (stdlib only) |

`server` must stay CGO-free and stdlib-only: the Dockerfile builds it with `CGO_ENABLED=0`.

## Main flows

1. **Startup**: `main` → `NewGame` → `core.WorldInit` → `core.LoadLevel(LevelFile)` →
   insert platforms into `DynamicQuadtree` → `InitPlayer(LoadImage(...))` → `NewEnemy(...)` (`cmd/main.go:33`).
2. **Tick (~60 TPS)**: `Game.Update` → `system.HandleInput` → `core.UpdatePlayer`
   (movement → gravity/jump → resolveHorizontal/Vertical via quadtree → `updateCombat`) →
   `UpdateAnimation` → `UpdateCamera` → `updateEnemies` (AI, anim, cull `Removable`) (`cmd/main.go:58`).
3. **Frame**: `Game.Draw` — draw order matters: water after player so the player looks submerged (`cmd/main.go:87`).
4. **Web build/deploy**: `GOOS=js GOARCH=wasm go build ./cmd` → `web/game.wasm` + `wasm_exec.js`
   → Docker multi-stage → distroless `server` on Render (`Dockerfile`, `render.yaml`).

## Data and state

- All runtime state is in memory, owned by `Game` in `cmd/main.go`. No persistence, no save files.
- Level data: `assets/level_1.json` (embedded). Assets are addressed by base file name only, so
  **two assets with the same file name in different folders collide** (`internal/core/image.go:20`).
- Single-threaded: no goroutines/channels in the current code.

## Secrets and security

- No secrets. Server reads only `PORT` and `WEB_ROOT` env vars. Serves static files only.

## Errors, logging

- Asset load failure → `log.Fatalf` (crash fast). Server logs each request with duration.

## Platform notes

| Capability | Windows native | WASM (browser) | Linux (CI) |
|---|---|---|---|
| Build | `go build ./cmd` | `GOOS=js GOARCH=wasm go build ./cmd` | needs `libgl1-mesa-dev xorg-dev libasound2-dev` |
| Run | `go run ./cmd` | `make serve` / `go run ./server` → :8080 | — |

## Testing seams

- Tests live in `internal/core` (package-internal). Logic that touches only data
  (level building, tile codes, enemy AI/damage) is testable without a window.
- No interface seam for time or input yet: `deltaTime()` reads `ebiten.TPS()` directly.

## Decisions

No ADRs yet. Add them under `docs/adr/` (template: `stepsForCreatingApp/templates/docs/adr.md`).
