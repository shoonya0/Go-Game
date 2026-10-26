# Go-Game (Pirate Adventure): current state (snapshot 2026-10-26)

Snapshot at `6a8cc90` (`main`; PR #25 main menu merged). One page. Details live in each increment's entry in
`status/increments.md` and each feature's `state.md`. Update this page when an increment lands.
Labels: **VERIFIED** = command run or code read · **ASSUMED** = not checked.

## Active work

<!-- Read by the SessionStart hook (.claude/hooks/session-context.mjs): keep this section
first, short, and with a link to exactly one feature state.md. Updated by /close-increment
and /feature-docs update. -->

- **Feature:** [levels](./features/levels/state.md) · Increment 1: eight themed levels, each a day and a night sub-level joined by portals (its `decisions.md`). Paused: enemy-types (low priority; task 25 waits for ranged art)
- **Branch:** `feature/levels` (from `main` at `8c93d3b`)
- **Next step:** PR #28 is open with the code review fixes; the user plays the levels: deepForest night, aquaticAbyss day and the 13 new generated grids (levels D8). Done: all 16 sub-levels have grids and are chained in `LevelOrder` (2026-10-26); the jump reaches just over 3 tiles, and `TestLevelPortalsReachEachOther` checks every grid's portals reach each other; aquaticAbyss day (level 2's first sub-level); the travel menu lists portals by level and sub-level (D7, no thumbnails); an invisible border walls in every level (D6); day/night sub-levels with deepForest night (tasks 6, 7, D5), tiles per level are tested (task 5; deepForest's larva pit is water now), the web build downloads only the current sub-level's background (task 4, D4); backgrounds renamed `<level>_day.*`/`<level>_night.*`; every file under `assets/` is committed (all 16 backgrounds, the Orange Golem strips and `SkeletonIdle.gif`, not loaded yet); `_asset_backups/` and `_futureAssets/` (Bat, bosses, Mossy tileset, plants, Striker) are the user's, not embedded; never commit them unless asked

<!-- Once a feature is active, the first line becomes, for example:
- **Feature:** [enemies](./features/enemies/state.md) · Increment 2
-->

## What a user can do today

- Run a 2D side-scrolling platformer natively (Windows) or in the browser (WASM). VERIFIED (CI builds both)
- Move (WASD), run (Shift; drains 5 stamina/s, exhausted at 0 until 20 regenerates), jump (lifts the feet just over 3 tiles, so 3-tile steps are climbable and 4 aren't; costs 5 stamina; at 0 stamina the player kneels until 2 bars refill), attack (J/I cost 10 stamina; U 34; O 10, plus 20 for its second part); sprite HUD with 10 hearts + 10 stamina segments; water/larva hazards, fall damage, respawn. VERIFIED (code read, tests)
- Portals (merged, PR #22): an animated stone portal stands on the ground at every code-`9` cell. A new game starts standing in the first one (reading order), already saved. Standing in a portal shows a pulsing gold button (E, or a tap ring on touch screens); E or a tap on it saves it, walking past does nothing. The saved portal is drawn gold, and a check flashes after a save. Portals the player never used stand still in dark gray; used ones swirl. Holding E (or a long press) at a portal for half a second opens a travel menu (pauses the game): a wooden board listing every active portal by level and sub-level ("Deep Forest", "Deep Forest Day", "Portal 1"...), picked with W/S or arrows and E/Enter, or a tap; Esc or a tap outside closes. The portal travelled to becomes the checkpoint, and every teleport brings the level's non-boss enemies back. Dying respawns the player standing in the saved portal with full health and power, and every non-boss enemy of the level comes back at its spawn (a killed boss would stay dead; no boss kind exists yet). VERIFIED (tests; web build; the user checked it in game)
- Fight six enemy types in level 1, each placed by a grid marker (one or two punctuation characters): the skeleton in white (`@`) and yellow (`#`), the mushroom (`&`) the White Golem (`%`, slow and heavy; checked by the user in game) and the Flying Demon (`^`, flies within 4 tiles of its spawn and bites; checked by the user in game) and the Smasher (`*`, hangs over you and dives down; checked by the user in game). All patrol, chase and attack; while patrolling they turn at walls and platform edges (chasing ones can still run off); they hit only when your body is inside their attack box. Each deals damage on the frame the weapon connects; the skeleton alternates two attacks. (The succubus and the goblin were removed: their animation didn't look good.) Between swings an enemy stands idle until its cooldown is ready, so no swing misses. Each shows 2 health bars above its head and regenerates 0.1 bar/s; sprites bob, lean and shake when hit (draw-time motion). Skeleton/mushroom checked by the user in game; hit timing not yet. Your attacks hurt it (J/I 0.5 bar, jump + I 1.5 and its spin hits behind you too, K 0.5 + 1, L 1 + 2, O 1 + 2); K, L and O stop halfway unless pressed again (O after sprite 9); U sends a ground wave along the floor that hurts it too (only enemies on the ground). Attacks move the player (`internal/core/attack_step.go`): J steps 45 px back, I and U 45 px forward, K's second part 15 px and O's 45 px forward; L doesn't move. Taking damage makes the player blink red for 1 s; a hit of 2+ hearts knocks them down 1.5 tiles back, 3+ hearts plays the full knockdown and 2 tiles; at 0 HP one of two death animations plays before respawn. F3 shows body/attack hitboxes for both. VERIFIED (`internal/core/enemy*.go`, tests pass)
- Enemies are built from parts, not one Go type per species: a blueprint (stats, clips from one strip pack) plus a movement part (`groundMover`) and an attack part (`meleeAttack`). What an enemy can do is the clips it has (`bp.Has`); `go test ./internal/core -run TestEnemyCapabilities -v` prints the table. A guard trait exists (blocks hits from the front) but no current enemy has guard art. VERIFIED (`internal/core/enemy_types.go`, `enemy_guard.go`)
- On a touch screen, buttons for every key: arrows bottom-left, attack symbols and a bigger sprint button bottom-right, a menu button top-left (the HUD moves right of it). Idle buttons are half see-through and stay out of the middle third of the screen. VERIFIED (tests; web build with simulated touches; not on a real device)
- Main menu (merged, PR #25): the game opens on a title screen (New Game, Continue, Settings, Quit; Press Start 2P font), picked with W/S + Enter/E, a tap or a mouse click. Progress (used portals, checkpoint) autosaves to `%AppData%\PirateAdventure\save.json` on Windows and `localStorage` on the web; Continue starts standing in the saved portal, New Game asks before wiping a save. Settings holds a saved volume 0–10 (no audio yet). Quit is desktop only. Esc or the touch menu button pauses (Resume, Settings, Main Menu). VERIFIED (tests; web build in agent-browser, `_docs/evidence/main-menu-*.png`; native window not checked by hand)
- Menus drawn from the wooden sheet `assets/UI/UIButtons.png` (merged, PR #27): planks with icons on the wooden board, the pick is the arrow plank, the volume a segmented bar; the travel menu's scroll arrows are wooden. VERIFIED (tests; web build screenshots `_docs/evidence/ui-buttons-*.png`; the user checked it in game)
- Levels (branch `feature/levels`, not merged): level 1 `deepForest` has a day and a night sub-level, each with its own grid and background; saving day's last unused portal teleports to night's entry portal once, and the travel menu then lists it. Every level is walled in by an invisible border. A test checks each level uses only its own tiles. The web build downloads only the current background (`game.wasm` 27.6 MB). All eight levels have both sub-levels, chained by D5 from deepForest day to rustyRefinery night; aquaticAbyss night and levels 3–8 are generated from one pattern (levels D8). VERIFIED (tests incl. portal reach for all 16; web build screenshots `_docs/evidence/levels-*.png` cover levels 1–2 only; the new grids not played)

## How to build, run, test

- Fast check: `node scripts/verify.mjs` · Full (adds native + WASM build): `node scripts/verify.mjs --full`
- Docs check: `node scripts/docs-check.mjs` (part of the fast check)
- Run native: `go run ./cmd` · Web: build wasm then `go run ./server` → http://localhost:8080 (see `DEPLOY.md`)
- `make` targets exist (`Makefile`) but `make` may be missing on Windows; use the commands above.

## Test baseline (at `af9a251` on `feature/levels`; Windows, go1.26.7)

| Command | Result | Failing tests / issues (by name) |
|---|---|---|
| `node scripts/verify.mjs --full` | ✅ pass, all 8 steps | — |
| `go test ./...` | ✅ 207/207 top-level tests pass: 195 in `internal/core`, 2 in `internal/system`, 9 in `tools/battler2strip`, 1 in `server` | — |
| `golangci-lint run ./...` (advisory) | 1 issue | SA9003 `internal/core/quadtree.go:217` |
| GitHub CI | `main` at `6a8cc90` (PR #25): merged by the user | — |

Compare runs **by test name**: `go test ./internal/core -v | grep "^--- "` lists them. Added since the last
baseline: the touch tests (`TestTouch*`), the O combo tests (`TestO*`), `TestAttackSteps`,
`TestAirHeavyAttackHitsBothSides`, `TestAttackDamageInBars`, `TestComboAttacksHitTwicePerSwing`, the
enemy health / level marker tests, and on `feature/enemy-types` the characterisation tests (written against the succubus, now run on the test-only `testEnemyBlueprint`)
(`internal/core/enemy_behavior_test.go`), the clip tests in `internal/core/enemy_clips_test.go`, hold/motion tests
(`enemy_motion_test.go`), melee hit-frame tests (`enemy_attack_test.go`), `TestBlueprintsValid`, and the baking tool's tests
(`tools/battler2strip/bake_test.go`), the ledge, marker and guard tests (`TestGroundEnemyTurnsAtLedge`, `TestTwoCharMarkerIsOneCell`, `enemy_guard_test.go`), the capability tests (`TestBlueprintHasOwnArtOnly`, `TestEnemyCapabilities`), `TestAssetBaseNamesUnique` and `TestLevelEnemiesKeepTheirGround`; on `feature/portals` the portal tests (`internal/core/portal_test.go`: sheet frames, standing on the ground, start, golden saved portal, E/tap saving, respawn), `TestTouchTapOffButtons` and `enemy_revive_test.go`; on `feature/main-menu` `TestSave*` (`save_test.go`), `TestSaveFile*` (`internal/system/save_other_test.go`), `TestTextSize`, `TestMenu*` (`menu_test.go`, `menu_screens_test.go`), `TestSettings*`, `TestTravelMenuWaitRelease`, `TestTouchHiddenButtons`; on `feature/levels` `TestBackgroundsDecode`, `TestLevelTiles`, `TestTilesTouch`, `TestLevelPortals`, `TestLastPortalOpensNextLevel`, `TestLevelBorderKeepsPlayerIn`, `TestTravelMenuGroupsByLevel`, `TestJumpReachesThreeTiles`, `TestHopRules`, `TestLevelPortalsReachEachOther`, `TestOldSaveOpensNextLevel`, `TestOpenNextLevelWithoutPortalKeepsWorld`, `TestLevelSizeFromGrid`.

## Known issues and risks

1. `game.exe` is tracked in git although `.gitignore` has `*.exe`. VERIFIED (`git ls-files`)
2. `go.mod` says `go 1.24.3` and CI uses 1.24; local toolchain is 1.26.7. Language features newer than 1.24 would pass locally and fail CI. VERIFIED
3. `readAsset` resolves by base file name — same-named assets in different folders silently collide (`internal/core/image.go:20`). The White Golem's files were renamed `WhiteGolem_*` (enemy-types task 16), and `TestAssetBaseNamesUnique` now fails on any duplicate among the embedded files. VERIFIED (test)
4. Large files: `animation.go` 434 lines (`InitPlayerAnimations` is most of it), `player.go` 397, `world.go` 343. VERIFIED (`wc -l`)
5. ~~`Layout` returns the outside window size, so the logical resolution changes with the window/browser size (`cmd/main.go:100`).~~ **FIXED** (`d6e602f`) — `Layout` now returns the fixed `screenWidth × screenHeight` constants; Ebitengine auto-scales to any device.
6. No tests for player physics with collision, quadtree, the server or `cmd/main.go` (menu → game wiring, autosave calls; checked only in the web build); combat, stamina, attack steps (without a level), touch input, HUD math, menus, save and storage are covered. VERIFIED (test names)
7. `DEPLOY.md` "Continuous integration" section may drift from the workflows. ASSUMED
8. ~~`game.wasm` served without `Cache-Control` header — browser caches the old binary across deploys, preventing fixes from reaching users.~~ **FIXED** (`f63b7e7`) — server now sends `no-cache` for `/game.wasm` alongside `index.html`.
9. `assets/UI/UIButtons.png` is 3.1 MB and embedded in the WASM; the menus and touch controls use about half its sprites (`internal/core/ui_skin.go`), so cropping to an atlas would shrink the phone download. VERIFIED (`du`)
10. The ground J sprite slashes behind the player, and J now steps back into it. Issue not filed. VERIFIED (sprite sheet)
11. The full knockdown (`KnockedOut`) still draws from row 17 while the short one uses row 16. ASSUMED (not rechecked this session)
12. ~~`//go:embed enemy` embeds untracked enemy art that CI and Render don't have~~ **FIXED**: every file under `assets/` is committed (the Orange Golem strips and `player/SkeletonIdle.gif` are embedded but not loaded yet). `//go:embed` still takes untracked files, so commit new assets with their code. VERIFIED (`git status`)
13. The Stop hook (`.claude/hooks/stop-state.mjs`) and docs-check's "stale state" warning fire while untracked files sit under a feature's `scope:` (`assets/enemy/*`), even right after state.md was committed. Noise until the check ignores untracked files (`assets/` has none now). VERIFIED (this session)
14. ~~Enemies have no ledge probe~~ Patrolling enemies turn at ledges (task 12), and one stranded off its patrol path patrols where it is (`3edfa17`); a chasing enemy can still run off a ledge. VERIFIED (`TestGroundEnemyTurnsAtLedge`)
15. Eight levels' `day`/`night` backgrounds (~110 MB) sit in `assets/levels/*/background/`, all committed and loaded per sub-level (the web build downloads one at a time, D4); they make the repo and the Docker image big. VERIFIED (`git ls-files`)
16. The menu font (`text/v2`) made `web/game.wasm` ~5 MB bigger: 28 → 33 MB with `-ldflags="-s -w"` (the text shaping library; the unused M+ font is dropped by the linker). VERIFIED (`ls -la web/game.wasm`)
17. The main menu was checked in the web build only; the native Windows window was not run by hand. Saves land in `%AppData%\PirateAdventure\` (test runs left one there). VERIFIED (this session)
18. Each title visit rebuilds the run: levels parsed again (code review; main-menu tasks "Later"). VERIFIED (code read)

## Open decisions

| # | Decision needed | Options | Owner | Blocking? |
|---|---|---|---|---|
| 1 | Untrack `game.exe` | `git rm --cached game.exe` / keep | user | no |
| 2 | Keep or delete `tools/battler2strip` (nothing left to bake, enemy-types D9) | keep / delete | user | no |
| 3 | Keep main-menu D4 (New Game asks before wiping a save) and D7 (Esc pause menu), both my calls | keep / drop | user | no |

## Next task

- Levels: the user plays deepForest night, then grids for levels 2–8 ([levels tasks](./features/levels/tasks.md)). Other candidates: real audio once sound files exist (main-menu D5); crop `UIButtons.png`; touch buttons for jump (Space), F, R, Q; file the J sprite issue; untrack `game.exe`; tests for player physics with collision.
