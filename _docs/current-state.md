# Go-Game (Pirate Adventure): current state (snapshot 2026-10-18)

Snapshot at `f63b7e7` on `main`. One page. Details live in each increment's entry in
`status/increments.md` and each feature's `state.md`. Update this page when an increment lands.
Labels: **VERIFIED** = command run or code read · **ASSUMED** = not checked.

## Active work

<!-- Read by the SessionStart hook (.claude/hooks/session-context.mjs): keep this section
first, short, and with a link to exactly one feature state.md. Updated by /close-increment
and /feature-docs update. -->

- **Feature:** none yet (set by `/feature-docs start`)
- **Branch:** `feature/stamina-costs` (from `main`) — normal attacks J/I cost 10 power (1 stamina segment), jumps cost 5 (half); blocked without enough stamina. Stamina HUD segments narrowed to 30 px. Verify passed. Separate unmerged branch `feature/level-enemy-markers` holds the level enemy markers (`!` in `level_1.json`)
- **Next step:** try the costs in-game, then push + PR for both branches
- **Updated:** 2026-10-20

<!-- Once a feature is active, the first line becomes, for example:
- **Feature:** [enemies](./features/enemies/state.md) · Increment 2
-->

## What a user can do today

- Run a 2D side-scrolling platformer natively (Windows) or in the browser (WASM). VERIFIED (CI builds both)
- Move (WASD), jump, attack; sprite health-heart/stamina HUD (10 segments each, on `feature/hud-sprite-bars`); water/larva hazards, fall damage, checkpoints, respawn. VERIFIED (code read)
- Fight one enemy type (succubus) with patrol/chase/attack AI. VERIFIED (`internal/core/enemy*.go` tracked, tests pass)
- `GameState` (menu/playing/game over) is declared and set to `ModeMenu`, but never read: there is no menu. VERIFIED (`cmd/main.go:45`)

## How to build, run, test

- Fast check: `node scripts/verify.mjs` · Full (adds native + WASM build): `node scripts/verify.mjs --full`
- Docs check: `node scripts/docs-check.mjs` (part of the fast check)
- Run native: `go run ./cmd` · Web: build wasm then `go run ./server` → http://localhost:8080 (see `DEPLOY.md`)
- `make` targets exist (`Makefile`) but `make` may be missing on Windows; use the commands above.

## Test baseline (at `f63b7e7`, Windows, go1.26.7)

| Command | Result | Failing tests / issues (by name) |
|---|---|---|
| `node scripts/verify.mjs` | ✅ pass (gofmt, vet, new lint issues, tests) | — |
| `go test ./...` | ✅ 9/9 pass, all in `internal/core` | — |
| `golangci-lint run ./...` (advisory) | 1 issue (was 3; the hud.go SA1019s went away with the sprite HUD) | SA9003 `internal/core/quadtree.go:217` |
| GitHub CI on `main` (`ci` + `verify`, ubuntu + windows) | ✅ green since 2026-10-08 | — |

Tests (+ `TestSegmentStates`, `TestSegmentStatesHealthSteps`, `TestRunningDrainsHalfSegmentPerSecond`, `TestRunHeldWithoutMovingRegenerates`, `TestExhaustionBlocksRunningUntilRecovered` on the HUD branch): `TestTileTypeFromCode`, `TestBuildLevelPlatformCount`, `TestBuildLevelCheckpoints`,
`TestBuildLevelTopBottomDetection`, `TestBuildLevelWorldCoordinates`, `TestBlueprintRegistered`,
`TestEnemyTakeDamage`, `TestEnemyPatrolTurnsAround`, `TestEnemyDetectsChasesAndAttacks`.
Compare new runs **by test name**.

## Known issues and risks

1. `game.exe` is tracked in git although `.gitignore` has `*.exe`. VERIFIED (`git ls-files`)
2. `go.mod` says `go 1.24.3` and CI uses 1.24; local toolchain is 1.26.7. Language features newer than 1.24 would pass locally and fail CI. VERIFIED
3. `readAsset` resolves by base file name — same-named assets in different folders silently collide (`internal/core/image.go:20`). VERIFIED
4. Large files: `player.go` 381 lines, `animation.go` 365 (`InitPlayerAnimations` is 222 lines), `world.go` 334. VERIFIED (`wc -l`)
5. ~~`Layout` returns the outside window size, so the logical resolution changes with the window/browser size (`cmd/main.go:100`).~~ **FIXED** (`d6e602f`) — `Layout` now returns the fixed `screenWidth × screenHeight` constants; Ebitengine auto-scales to any device.
6. No tests for player physics, combat, quadtree or the server. VERIFIED (test list above)
7. `DEPLOY.md` "Continuous integration" section may drift from the workflows. ASSUMED
8. ~~`game.wasm` served without `Cache-Control` header — browser caches the old binary across deploys, preventing fixes from reaching users.~~ **FIXED** (`f63b7e7`) — server now sends `no-cache` for `/game.wasm` alongside `index.html`.

## Open decisions

| # | Decision needed | Options | Owner | Blocking? |
|---|---|---|---|---|
| 1 | Untrack `game.exe` | `git rm --cached game.exe` / keep | user | no |

## Next task

- Not chosen yet. Candidates from the issues above: untrack `game.exe`; tests for player physics or combat.
