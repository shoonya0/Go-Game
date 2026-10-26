# Architecture overview

> Snapshot 2026-10-07 at `91400c1`; enemy rows rechecked 2026-10-23 at `59449eb` (branch `feature/enemy-types`). Map for humans and agents.
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
| Input | `internal/system/Input.go`, `touch.go`, `touch_js.go` / `touch_other.go` | Polls keyboard, mouse and touch points → `core.InputState`; `IsTouchDevice` asks the browser (`navigator.maxTouchPoints`), false on desktop builds | no (writes caller's struct) |
| Storage | `internal/system/save_other.go` / `save_js.go` | Keeps named records (`save`, `settings`): a JSON file per key in the user's config folder on desktop, `localStorage` on the web | on disk / in the browser |
| Menus | `internal/core/menu.go`, `menu_screens.go`, `text.go` | A menu widget (keys, tap, mouse) and the screens over the game: title, overwrite confirm, settings, pause; `Screens.Update` returns an action for `cmd`; text in Press Start 2P (`text/v2`) | `Game.screens` |
| Save / settings | `internal/core/save.go`, `settings.go` | Snapshot of the player's portals and checkpoint, `World.Continue` back into it, `Autosave` (write on change); the volume setting. (De)serialize only: storage is `internal/system` | `Game.autosave` |
| Touch controls | `internal/core/touch_controls.go`, `touch_icons.go` | On-screen buttons for touch screens: layout, hit-testing, merging into `InputState` (a touch starting on no button is a world tap), drawing (sprites from `UIButtons.png`, symbols drawn with `ebiten/vector`) | `Game.touch` |
| Player | `internal/core/player.go`, `player_runtime.go`, `playerState.go`, `player_attack.go`, `player_hurt.go`, `ground_wave.go` | Movement, gravity/jump, AABB collision resolution, state enum; `HitEnemies` damages enemies inside the attack's damage box, once per hit window; `player_hurt.go` handles heavy-hit knockback, the out-of-stamina kneel and the two deaths; `UpdateGroundWaves` runs U's ground wave along the floor; `attack_step.go` moves the player during some attacks (table `attackSteps`) | `PlayerRuntime` |
| Combat | `internal/core/combat.go`, `skills.go` | Damage, power regen, hazards (water/larva), fall damage, saving a portal (E / tap), respawn | inside `PlayerRuntime` |
| Enemies | `internal/core/enemy*.go` | Blueprint registry (each blueprint lists its animations as `EnemyClip`s under standard clip names with fallbacks, `enemy_clips.go`; every clip names its own sheet, cached by path), clips may hold frames and add draw-time motion (`enemy_motion.go`); a shared brain (patrol/chase/attack, `enemy_ai.go`) drives each species' `Mover` (`enemy_move.go`, today `groundMover`) and `Attacker` (`enemy_attack.go`, today `meleeAttack`, which lands on a clip's `HitFrame`); `AttackBox()` is the area that hits the player; `enemy_health.go` regenerates health and draws the health bars above the sprite | `EnemyRuntime` |
| Animation | `internal/core/animation.go`, `enemy_anim.go` | Sprite-sheet frame tables, frame advance, drawing with camera | per-runtime |
| World / level | `internal/core/world.go`, `levels.go`, `level_enemies.go` | Loads embedded `level_1.json`, builds platforms/tiles/portals (code 9) and enemy spawns (bare grid markers like `!` = `EnemyBlueprint.Marker`), draws level & water; `RespawnEnemies` revives non-boss enemies after the player respawns; `Levels` names each level and `World.SwitchLevel` swaps the current one | `Game.world` (`World`) |
| Portals | `internal/core/portal.go`, `travel_menu.go` | Portal sprite (frames cropped to their art) in three looks (inactive still, active looping, checkpoint gold), the use prompt (icons, no font); the travel menu listing active portals by level (pauses the game) | `Game.portals`, `Game.travel` |
| Spatial index | `internal/core/quadtree.go`, `aabb.go` | `DynamicQuadtree` for collision broad-phase; `Collider` interface | quadtree |
| HUD | `internal/core/hud.go` | Health hearts / stamina (power) bars, 10 sprite segments each | no |
| Debug overlay | `internal/core/debug_draw.go` | F3 toggles outlines of body (green) and attack (red) boxes for any `HitboxProvider` (enemies and the player) | `Game.showHitboxes` |
| Assets | `assets/embed.go`, `backgrounds_*.go`, `internal/core/image.go` | `go:embed` FS; `readAsset` resolves by **base file name** (subfolders searched); files live in `player/`, `enemy/`, `UI/`, `levels/`, and each folder (in `levels/`, each file pattern) must be on the `//go:embed` line. Level backgrounds are a separate FS, `assets.Backgrounds`, empty in the web build: `cmd` downloads the current level's with `system.FetchFile` off the game loop. Raw source art that is baked offline goes in `art-src/`, outside the embed (empty now: the battler sheets were deleted) | embedded, read-only |
| Time | `internal/core/time.go` | `deltaTime()` = 1/TPS, shared by physics and animation | no |
| Web server | `server/main.go` | Serves `$WEB_ROOT` on `$PORT`, `.wasm` MIME, no-cache for HTML; level backgrounds from `$LEVELS_ROOT` under `/levels/` | no |
| Battler baker | `tools/battler2strip/` | Offline, not in the game build: bakes an RPG Maker battler sheet + mapping JSON into right-facing strips and prints a blueprint snippet | no |

## Layers and allowed imports (target, enforced by golangci-lint depguard)

| From \ To | cmd | internal/system | internal/core | assets | ebiten |
|---|---|---|---|---|---|
| **cmd** | — | ✅ | ✅ | ❌ (go through core) | ✅ |
| **internal/system** | ❌ | — | ✅ (types only) | ❌ | ✅ |
| **internal/core** | ❌ | ❌ | — | ✅ | ✅ |
| **server** | ❌ | ❌ | ❌ | ❌ | ❌ (stdlib only) |
| **tools/** | ❌ | ❌ | ❌ | ❌ | ❌ (stdlib only) |

`server` must stay CGO-free and stdlib-only: the Dockerfile builds it with `CGO_ENABLED=0`.
`tools/` holds offline programs that aren't part of the game build (`tools/battler2strip` bakes
battler sheets into enemy strips, enemy-types D7); they stay stdlib-only too.

## Main flows

1. **Startup**: `main` → `NewGame` → `core.WorldInit` → load settings → `openTitle`: read the save and
   `startRun` (a `World` via `core.NewWorld`, `InitPlayer`, `StartAt`, then `World.Continue` when the save is usable),
   so the title screen shows the game Continue would start. While `Game.screens.Current` isn't `ScreenNone`,
   `Update` runs only the menus and `Draw` puts the menu over the scene (`cmd/main.go`).
2. **Tick (~60 TPS)**: `Game.Update` → `system.HandleInput` → `core.UpdatePlayer`
   (movement → gravity/jump → resolveHorizontal/Vertical via quadtree → `updateCombat`) →
   `UpdateAnimation` → `HitEnemies` → `UpdateGroundWaves` → `UpdateCamera` → `updateEnemies` (AI, anim, cull `Removable`) (`cmd/main.go:58`).
3. **Frame**: `Game.Draw` — draw order matters: water after player so the player looks submerged (`cmd/main.go:87`).
4. **Web build/deploy**: `GOOS=js GOARCH=wasm go build ./cmd` → `web/game.wasm` + `wasm_exec.js`
   → Docker multi-stage → distroless `server` on Render (`Dockerfile`, `render.yaml`).

## Data and state

- Runtime state is in memory, owned by `Game` in `cmd/main.go`. Two records persist (main-menu feature):
  the progress (`save`: used portals + checkpoint, written on change) and the settings (`settings`: volume),
  in `%AppData%\PirateAdventure\*.json` on Windows and `localStorage` (`pirateAdventure.*`) on the web.
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
