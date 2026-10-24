# Go-Game (Pirate Adventure): current state (snapshot 2026-10-24)

Snapshot at `c2c3b87` (`main`, PR #20 merged). One page. Details live in each increment's entry in
`status/increments.md` and each feature's `state.md`. Update this page when an increment lands.
Labels: **VERIFIED** = command run or code read · **ASSUMED** = not checked.

## Active work

<!-- Read by the SessionStart hook (.claude/hooks/session-context.mjs): keep this section
first, short, and with a link to exactly one feature state.md. Updated by /close-increment
and /feature-docs update. -->

- **Feature:** [enemy-types](./features/enemy-types/state.md) · Increments 1–3 merged (PR #18, #19, #20; tasks 1–21, 15/19/20 dropped); increment 4 (flying enemies) in progress: tasks 22-23 done
- **Branch:** `feature/enemy-types`, fresh from `main` at `c2c3b87` (PR #20 merged by the user)
- **Next step:** task 24: the Flying Demon (`^`) is built, waiting for your look in game (F3) in [tasks.md](./features/enemy-types/tasks.md). Untracked assets (`assets/enemy/*` new folders, `assets/levels/portal.png`, `assets/player/SkeletonIdle.gif`) are committed with the code that loads them. `deepForest/background/LBackground_*.png` names will collide with a second level's backgrounds (lookup is by base name)
- **Updated:** 2026-10-24 (increment 3 merged, PR #20)

<!-- Once a feature is active, the first line becomes, for example:
- **Feature:** [enemies](./features/enemies/state.md) · Increment 2
-->

## What a user can do today

- Run a 2D side-scrolling platformer natively (Windows) or in the browser (WASM). VERIFIED (CI builds both)
- Move (WASD), run (Shift; drains 5 stamina/s, exhausted at 0 until 20 regenerates), jump (costs 5 stamina; at 0 stamina the player kneels until 2 bars refill), attack (J/I cost 10 stamina; U 34; O 10, plus 20 for its second part); sprite HUD with 10 hearts + 10 stamina segments; water/larva hazards, fall damage, checkpoints, respawn. VERIFIED (code read, tests)
- Fight five enemy types in level 1, each placed by a grid marker (one or two punctuation characters): the skeleton in white (`@`) and yellow (`#`), the mushroom (`&`) the White Golem (`%`, slow and heavy; checked by the user in game) and the Flying Demon (`^`, flies within 4 tiles of its spawn and bites; not yet checked in game). All patrol, chase and attack; while patrolling they turn at walls and platform edges (chasing ones can still run off); they hit only when your body is inside their attack box. Each deals damage on the frame the weapon connects; the skeleton alternates two attacks. (The succubus and the goblin were removed: their animation didn't look good.) Between swings an enemy stands idle until its cooldown is ready, so no swing misses. Each shows 2 health bars above its head and regenerates 0.1 bar/s; sprites bob, lean and shake when hit (draw-time motion). Skeleton/mushroom checked by the user in game; hit timing not yet. Your attacks hurt it (J/I 0.5 bar, jump + I 1.5 and its spin hits behind you too, K 0.5 + 1, L 1 + 2, O 1 + 2); K, L and O stop halfway unless pressed again (O after sprite 9); U sends a ground wave along the floor that hurts it too (only enemies on the ground). Attacks move the player (`internal/core/attack_step.go`): J steps 45 px back, I and U 45 px forward, K's second part 15 px and O's 45 px forward; L doesn't move. Taking damage makes the player blink red for 1 s; a hit of 2+ hearts knocks them down 1.5 tiles back, 3+ hearts plays the full knockdown and 2 tiles; at 0 HP one of two death animations plays before respawn. F3 shows body/attack hitboxes for both. VERIFIED (`internal/core/enemy*.go`, tests pass)
- Enemies are built from parts, not one Go type per species: a blueprint (stats, clips from one strip pack) plus a movement part (`groundMover`) and an attack part (`meleeAttack`). What an enemy can do is the clips it has (`bp.Has`); `go test ./internal/core -run TestEnemyCapabilities -v` prints the table. A guard trait exists (blocks hits from the front) but no current enemy has guard art. VERIFIED (`internal/core/enemy_types.go`, `enemy_guard.go`)
- On a touch screen, buttons for every key: arrows bottom-left, attack symbols and a bigger sprint button bottom-right, a menu button top-left (the HUD moves right of it). Idle buttons are half see-through and stay out of the middle third of the screen. VERIFIED (tests; web build with simulated touches; not on a real device)
- `GameState` (menu/playing/game over) is declared and set to `ModeMenu`, but never read: there is no menu, so the touch menu button and Escape only toggle `InputState.Menu`. VERIFIED (`cmd/main.go`)

## How to build, run, test

- Fast check: `node scripts/verify.mjs` · Full (adds native + WASM build): `node scripts/verify.mjs --full`
- Docs check: `node scripts/docs-check.mjs` (part of the fast check)
- Run native: `go run ./cmd` · Web: build wasm then `go run ./server` → http://localhost:8080 (see `DEPLOY.md`)
- `make` targets exist (`Makefile`) but `make` may be missing on Windows; use the commands above.

## Test baseline (at `e5264dc`, Windows, go1.26.7)

| Command | Result | Failing tests / issues (by name) |
|---|---|---|
| `node scripts/verify.mjs --full` | ✅ pass, all 8 steps | — |
| `go test ./...` | ✅ 126/126 top-level tests pass: 117 in `internal/core`, 9 in `tools/battler2strip` | — |
| `golangci-lint run ./...` (advisory) | 1 issue | SA9003 `internal/core/quadtree.go:217` |
| GitHub CI | PR #20 and `main` at `c2c3b87`: green (build, verify ubuntu + windows) | — |

Compare runs **by test name**: `go test ./internal/core -v | grep "^--- "` lists them. Added since the last
baseline: the touch tests (`TestTouch*`), the O combo tests (`TestO*`), `TestAttackSteps`,
`TestAirHeavyAttackHitsBothSides`, `TestAttackDamageInBars`, `TestComboAttacksHitTwicePerSwing`, the
enemy health / level marker tests, and on `feature/enemy-types` the characterisation tests (written against the succubus, now run on the test-only `testEnemyBlueprint`)
(`internal/core/enemy_behavior_test.go`), the clip tests in `internal/core/enemy_clips_test.go`, hold/motion tests
(`enemy_motion_test.go`), melee hit-frame tests (`enemy_attack_test.go`), `TestBlueprintsValid`, and the baking tool's tests
(`tools/battler2strip/bake_test.go`), the ledge, marker and guard tests (`TestGroundEnemyTurnsAtLedge`, `TestTwoCharMarkerIsOneCell`, `enemy_guard_test.go`), the capability tests (`TestBlueprintHasOwnArtOnly`, `TestEnemyCapabilities`), `TestAssetBaseNamesUnique` and `TestLevelEnemiesKeepTheirGround`.

## Known issues and risks

1. `game.exe` is tracked in git although `.gitignore` has `*.exe`. VERIFIED (`git ls-files`)
2. `go.mod` says `go 1.24.3` and CI uses 1.24; local toolchain is 1.26.7. Language features newer than 1.24 would pass locally and fail CI. VERIFIED
3. `readAsset` resolves by base file name — same-named assets in different folders silently collide (`internal/core/image.go:20`). The White Golem's files were renamed `WhiteGolem_*` (enemy-types task 16), and `TestAssetBaseNamesUnique` now fails on any duplicate among the embedded files. VERIFIED (test)
4. Large files: `animation.go` 434 lines (`InitPlayerAnimations` is most of it), `player.go` 397, `world.go` 343. VERIFIED (`wc -l`)
5. ~~`Layout` returns the outside window size, so the logical resolution changes with the window/browser size (`cmd/main.go:100`).~~ **FIXED** (`d6e602f`) — `Layout` now returns the fixed `screenWidth × screenHeight` constants; Ebitengine auto-scales to any device.
6. No tests for player physics with collision, quadtree or the server; combat, stamina, attack steps (without a level), touch input and HUD math are covered. VERIFIED (test names)
7. `DEPLOY.md` "Continuous integration" section may drift from the workflows. ASSUMED
8. ~~`game.wasm` served without `Cache-Control` header — browser caches the old binary across deploys, preventing fixes from reaching users.~~ **FIXED** (`f63b7e7`) — server now sends `no-cache` for `/game.wasm` alongside `index.html`.
9. `assets/UI/UIButtons.png` is 3.1 MB and embedded in the WASM, though touch controls use only three of its sprites; cropping it would shrink the phone download. VERIFIED (`du`)
10. The ground J sprite slashes behind the player, and J now steps back into it. Issue not filed. VERIFIED (sprite sheet)
11. The full knockdown (`KnockedOut`) still draws from row 17 while the short one uses row 16. ASSUMED (not rechecked this session)
12. `//go:embed enemy` embeds every file on disk, untracked ones included, so local builds carry untracked enemy art that CI and Render don't (31 untracked files: 5 enemy packs, `levels/portal.png`, `player/SkeletonIdle.gif`; the portal and GIF are planned, keep them). The raw battler sheets are gone: moved out of the embed (task 11), then deleted (enemy-types D9). VERIFIED (`git ls-files`)
13. The Stop hook (`.claude/hooks/stop-state.mjs`) and docs-check's "stale state" warning fire while untracked files sit under a feature's `scope:` (`assets/enemy/*`), even right after state.md was committed. Noise until those assets are committed or the check ignores untracked files. VERIFIED (this session)
14. ~~Enemies have no ledge probe~~ Patrolling enemies turn at ledges (task 12), and one stranded off its patrol path patrols where it is (`3edfa17`); a chasing enemy can still run off a ledge. VERIFIED (`TestGroundEnemyTurnsAtLedge`)
15. `LBackground_1.png` and `LBackground_2.png` (2.6 MB) are embedded but not loaded; only `LBackground_3.png` is. Kept on purpose (user). VERIFIED (code read)

## Open decisions

| # | Decision needed | Options | Owner | Blocking? |
|---|---|---|---|---|
| 1 | Untrack `game.exe` | `git rm --cached game.exe` / keep | user | no |
| 2 | Keep or delete `tools/battler2strip` (nothing left to bake, enemy-types D9) | keep / delete | user | no |

## Next task

- enemy-types task 22: `flyingMover` ([tasks.md](./features/enemy-types/tasks.md)), start of increment 4.
- Other candidates: draw a real menu for the menu button; crop `UIButtons.png`; touch buttons for jump (Space), F, R, Q; file the J sprite issue; untrack `game.exe`; tests for player physics with collision.
