# Architecture overview

> Snapshot 2026-10-07 at `91400c1`; enemy rows rechecked 2026-10-23 at `ba831b3` (branch `feature/enemy-types`). Map for humans and agents.
> Rows marked **target** describe intended rules the code does not fully enforce yet.

Go 1.24 module `player` · Ebitengine v2.9 · one codebase, two targets: native desktop
(Windows) and WebAssembly (browser), plus a tiny static server for deployment (Render).

## Overview diagram

```text
 cmd/main.go (Game: ebiten.Game)
   │  Update(): system.HandleInput + HandleTouch → core.UpdatePlayer → anim/camera → updateEnemies
   │  Draw():   parallax bg → level tiles → enemies → player → water → HUD → touch buttons
   ▼
 internal/system ──► internal/core ──► assets (go:embed FS)
   (keyboard + touch → InputState)   (simulation + rendering)

 server/main.go  — standalone static file server for ./web (no game imports)
 web/index.html  — boots game.wasm via wasm_exec.js
```

## Modules

| Module | Path | Responsibility | Owns state? |
|---|---|---|---|
| Entry / game loop | `cmd/main.go` | Wires subsystems; implements `ebiten.Game` (Update/Draw/Layout) | `Game` struct: player, enemies, level, quadtree |
| Input | `internal/system/Input.go`, `touch.go`, `touch_js.go` / `touch_other.go` | Polls keyboard and touch points → `core.InputState`; `IsTouchDevice` asks the browser (`navigator.maxTouchPoints`), false on desktop builds | no (writes caller's struct) |
| Touch controls | `internal/core/touch_controls.go`, `touch_icons.go` | On-screen buttons for touch screens: layout, hit-testing, merging into `InputState`, drawing (sprites from `UIButtons.png`, symbols drawn with `ebiten/vector`) | `Game.touch` |
| Player | `internal/core/player.go`, `player_runtime.go`, `playerState.go`, `player_attack.go`, `player_hurt.go`, `ground_wave.go` | Movement, gravity/jump, AABB collision resolution, state enum; `HitEnemies` damages enemies inside the attack's damage box, once per hit window; `player_hurt.go` handles heavy-hit knockback, the out-of-stamina kneel and the two deaths; `UpdateGroundWaves` runs U's ground wave along the floor; `attack_step.go` moves the player during some attacks (table `attackSteps`) | `PlayerRuntime` |
| Combat | `internal/core/combat.go`, `skills.go` | Damage, power regen, hazards (water/larva), fall damage, checkpoints, respawn | inside `PlayerRuntime` |
| Enemies | `internal/core/enemy*.go` | Blueprint registry (each blueprint lists its animations as `EnemyClip`s under standard clip names with fallbacks, `enemy_clips.go`; every clip names its own sheet, cached by path), AI (patrol/chase/attack), animation, damage; `AttackBox()` is the area that hits the player; `enemy_health.go` regenerates health and draws the health bars above the sprite | `EnemyRuntime` |
| Animation | `internal/core/animation.go`, `enemy_anim.go` | Sprite-sheet frame tables, frame advance, drawing with camera | per-runtime |
| World / level | `internal/core/world.go`, `level_enemies.go` | Loads embedded `level_1.json`, builds platforms/tiles/checkpoints and enemy spawns (bare grid markers like `!` = `EnemyBlueprint.Marker`), draws level & water | `Level`, `SpawnEnemies` |
| Spatial index | `internal/core/quadtree.go`, `aabb.go` | `DynamicQuadtree` for collision broad-phase; `Collider` interface | quadtree |
| HUD | `internal/core/hud.go` | Health hearts / stamina (power) bars, 10 sprite segments each | no |
| Debug overlay | `internal/core/debug_draw.go` | F3 toggles outlines of body (green) and attack (red) boxes for any `HitboxProvider` (enemies and the player) | `Game.showHitboxes` |
| Assets | `assets/embed.go`, `internal/core/image.go` | `go:embed` FS; `readAsset` resolves by **base file name** (subfolders searched); files live in `player/`, `enemy/`, `UI/`, `levels/`, and each folder must be on the `//go:embed` line | embedded, read-only |
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
   insert platforms into `DynamicQuadtree` → `InitPlayer(LoadImage(...))` → `SpawnEnemies(level.Enemies)` (markers from the grid) (`cmd/main.go:33`).
2. **Tick (~60 TPS)**: `Game.Update` → `system.HandleInput` → `core.UpdatePlayer`
   (movement → gravity/jump → resolveHorizontal/Vertical via quadtree → `updateCombat`) →
   `UpdateAnimation` → `HitEnemies` → `UpdateGroundWaves` → `UpdateCamera` → `updateEnemies` (AI, anim, cull `Removable`) (`cmd/main.go:58`).
3. **Frame**: `Game.Draw` — draw order matters: water after player so the player looks submerged (`cmd/main.go:87`).
4. **Web build/deploy**: `GOOS=js GOARCH=wasm go build ./cmd` → `web/game.wasm` + `wasm_exec.js`
   → Docker multi-stage → distroless `server` on Render (`Dockerfile`, `render.yaml`).

## Data and state

- All runtime state is in memory, owned by `Game` in `cmd/main.go`. No persistence, no save files.
- Level data: `assets/levels/deepForest/level_1.json` (embedded). Assets are addressed by base file name only, so
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

No ADRs yet. Add them under `_docs/adr/` (template: the app-setup playbook's `templates/docs/adr.md`).
