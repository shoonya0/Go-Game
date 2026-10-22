# Go-Game (Pirate Adventure): current state (snapshot 2026-10-22)

Snapshot at `2abec37` on `feature/level-enemy-markers` (`main` is at PR #12, `c59d6b0`). One page. Details live in each increment's entry in
`status/increments.md` and each feature's `state.md`. Update this page when an increment lands.
Labels: **VERIFIED** = command run or code read · **ASSUMED** = not checked.

## Active work

<!-- Read by the SessionStart hook (.claude/hooks/session-context.mjs): keep this section
first, short, and with a link to exactly one feature state.md. Updated by /close-increment
and /feature-docs update. -->

- **Feature:** none yet (set by `/feature-docs start`)
- **Branch:** `feature/o-combo` (from `main` at PR #13, issue #14): O is a two-press combo like K/L. The first press plays sprites 0-8 and costs 1 HUD segment (`SpecialAttack4PowerCost1`); a second O before frame 9 (`specialAttack4GateFrame`) finishes all 17 frames for 2 more segments (`SpecialAttack4PowerCost2`), and is ignored without that stamina. O reads as a single press. O now hits: burst on frames 1-4 (1 bar, 70 px reach), sword thrust on frames 12-14 (2 bars, 95 px; `strikeHit.Reach` lets a hit override the attack reach); attacks step the player 45 px (`attack_step.go` table): J back on frames 1-2, I forward 1-2, U forward 2-3, K and L forward 22.5 px (half step) on 5-6 and O forward 9-11, those three only after the second press; then the player holds still. Verify `--full` passed. Also open: `feature/mobile-settings` (touch controls, its O button already acts once per press)
- **Next step:** try O in-game (`go run ./cmd`): one press stops after sprite 9, two presses finish it, stamina bar drops 1 then 2 segments; then push + PR for both branches (`_docs/current-state.md` will conflict between them: keep both). Open: the full knockdown (`KnockedOut`) still draws from row 17
- **Updated:** 2026-10-22

<!-- Once a feature is active, the first line becomes, for example:
- **Feature:** [enemies](./features/enemies/state.md) · Increment 2
-->

## What a user can do today

- Run a 2D side-scrolling platformer natively (Windows) or in the browser (WASM). VERIFIED (CI builds both)
- Move (WASD), run (Shift; drains 5 stamina/s, exhausted at 0 until 20 regenerates), jump (costs 5 stamina; at 0 stamina the player kneels until 2 bars refill), attack (J/I cost 10 stamina; U/O specials 34); sprite HUD with 10 hearts + 10 stamina segments; water/larva hazards, fall damage, checkpoints, respawn. VERIFIED (code read, tests)
- Fight one enemy type (succubus) with patrol/chase/attack AI; it hits only when your body is inside its attack box, placed from the level file by grid marker; it shows 2 health bars above its head and regenerates 0.1 bar/s. Your attacks hurt it (J/I 0.5 bar, jump + I 1.5, K 0.5 + 1, L 1 + 2; K/L stop halfway unless pressed again); U sends a ground wave along the floor that hurts it too (O deals no damage yet). Taking damage makes the player blink red for 1 s; a hit of 2+ hearts knocks them down 1.5 tiles back, 3+ hearts plays the full knockdown and 2 tiles; at 0 HP one of two death animations plays before respawn. F3 shows body/attack hitboxes for both. VERIFIED (`internal/core/enemy*.go` tracked, tests pass)
- `GameState` (menu/playing/game over) is declared and set to `ModeMenu`, but never read: there is no menu. VERIFIED (`cmd/main.go:45`)

## How to build, run, test

- Fast check: `node scripts/verify.mjs` · Full (adds native + WASM build): `node scripts/verify.mjs --full`
- Docs check: `node scripts/docs-check.mjs` (part of the fast check)
- Run native: `go run ./cmd` · Web: build wasm then `go run ./server` → http://localhost:8080 (see `DEPLOY.md`)
- `make` targets exist (`Makefile`) but `make` may be missing on Windows; use the commands above.

## Test baseline (at `8eff08b`, Windows, go1.26.7)

| Command | Result | Failing tests / issues (by name) |
|---|---|---|
| `node scripts/verify.mjs` | ✅ pass (gofmt, vet, new lint issues, tests) | — |
| `go test ./...` | ✅ 36/36 pass, all in `internal/core` | — |
| `golangci-lint run ./...` (advisory) | 1 issue (was 3; the hud.go SA1019s went away with the sprite HUD) | SA9003 `internal/core/quadtree.go:217` |
| GitHub CI on `main` (`ci` + `verify`, ubuntu + windows) | ✅ green since 2026-10-08 | — |

Tests: `TestTileTypeFromCode`, `TestBuildLevelPlatformCount`, `TestBuildLevelCheckpoints`,
`TestBuildLevelTopBottomDetection`, `TestBuildLevelWorldCoordinates`, `TestBlueprintRegistered`,
`TestEnemyTakeDamage`, `TestEnemyPatrolTurnsAround`, `TestEnemyDetectsChasesAndAttacks`,
`TestSegmentStates`, `TestSegmentStatesHealthSteps`, `TestRunningDrainsHalfSegmentPerSecond`,
`TestRunHeldWithoutMovingRegenerates`, `TestExhaustionBlocksRunningUntilRecovered`, `TestJumpCostsHalfSegment`,
`TestNormalAttacksCostOneSegment`, `TestAirAttackCostsOneSegment`, `TestEnemyAttackBoxFacesForward`, `TestEnemyDoesNotHitPlayerAbove`,
`TestEnemyAttacksPlayerBehind`, `TestEnemyHitboxes`, `TestPlayerAttackBoxFacesForward`,
`TestPlayerHitsEnemyOncePerSwing`, `TestPlayerMissesEnemyBehindOrAbove`, `TestPlayerHitboxes`,
`TestEnemyHitRestartsHurtClip` (all on `main`);
this branch adds `TestDamageFlashBlinksRedForOneSecond`, `TestZeroDamageDoesNotFlash`,
`TestRespawnClearsDamageFlash`, `TestGroundWaveCastOncePerSwing`, `TestGroundWaveRunsTheWayThePlayerFaces`,
`TestGroundWaveFadesAfterItsFrames`, `TestGroundWaveNeedsGround`, `TestGroundWaveStopsAtWall`,
`TestGroundWaveHitsEachEnemyOnce`, `TestPlayerHitboxesIncludeGroundWaves`.
Compare new runs **by test name**.

## Known issues and risks

1. `game.exe` is tracked in git although `.gitignore` has `*.exe`. VERIFIED (`git ls-files`)
2. `go.mod` says `go 1.24.3` and CI uses 1.24; local toolchain is 1.26.7. Language features newer than 1.24 would pass locally and fail CI. VERIFIED
3. `readAsset` resolves by base file name — same-named assets in different folders silently collide (`internal/core/image.go:20`). VERIFIED
4. Large files: `player.go` 381 lines, `animation.go` 365 (`InitPlayerAnimations` is 222 lines), `world.go` 334. VERIFIED (`wc -l`)
5. ~~`Layout` returns the outside window size, so the logical resolution changes with the window/browser size (`cmd/main.go:100`).~~ **FIXED** (`d6e602f`) — `Layout` now returns the fixed `screenWidth × screenHeight` constants; Ebitengine auto-scales to any device.
6. No tests for player physics (movement/collision), quadtree or the server; combat/stamina rules and HUD segment math are covered. VERIFIED (test list above)
7. `DEPLOY.md` "Continuous integration" section may drift from the workflows. ASSUMED
8. ~~`game.wasm` served without `Cache-Control` header — browser caches the old binary across deploys, preventing fixes from reaching users.~~ **FIXED** (`f63b7e7`) — server now sends `no-cache` for `/game.wasm` alongside `index.html`.

## Open decisions

| # | Decision needed | Options | Owner | Blocking? |
|---|---|---|---|---|
| 1 | Untrack `game.exe` | `git rm --cached game.exe` / keep | user | no |

## Next task

- Not chosen yet. Candidates from the issues above: untrack `game.exe`; tests for player physics or combat.
