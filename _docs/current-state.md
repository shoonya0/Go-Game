# Go-Game (Pirate Adventure): current state (snapshot 2026-10-22)

Snapshot at `a2c2e8e` on `main` (PR #16 merged). One page. Details live in each increment's entry in
`status/increments.md` and each feature's `state.md`. Update this page when an increment lands.
Labels: **VERIFIED** = command run or code read · **ASSUMED** = not checked.

## Active work

<!-- Read by the SessionStart hook (.claude/hooks/session-context.mjs): keep this section
first, short, and with a link to exactly one feature state.md. Updated by /close-increment
and /feature-docs update. -->

- **Feature:** none yet (set by `/feature-docs start`)
- **Branch:** `refactor/asset-folders` (PR #17, includes the two post-#15/#16 state refreshes; CI not checked yet): assets grouped into `assets/player/`, `enemy/`, `UI/`, `levels/` (level 1 in `levels/deepForest/`, backgrounds in its `background/`); `assets/embed.go` lists the folders, so **a new asset folder must be added to its `//go:embed` line** (`TestAssetPathsAreEmbedded` in `internal/core/asset_paths_test.go` fails otherwise). PR #15 and #16 are merged into `main`
- **Next step:** CI green on the asset-folders PR, then ask before merging (deploys to Render). Untracked, not in the PR (no code uses them yet): `assets/levels/portal.png`, `assets/enemy/Bat with VFX/` (8 sheets), `assets/enemy/Mushroom/` (7 sheets); commit them with the code that loads them. `deepForest/background/LBackground_*.png` names will collide with a second level's backgrounds (lookup is by base name). Then try the touch controls on a real phone/tablet
- **Updated:** 2026-10-23

<!-- Once a feature is active, the first line becomes, for example:
- **Feature:** [enemies](./features/enemies/state.md) · Increment 2
-->

## What a user can do today

- Run a 2D side-scrolling platformer natively (Windows) or in the browser (WASM). VERIFIED (CI builds both)
- Move (WASD), run (Shift; drains 5 stamina/s, exhausted at 0 until 20 regenerates), jump (costs 5 stamina; at 0 stamina the player kneels until 2 bars refill), attack (J/I cost 10 stamina; U 34; O 10, plus 20 for its second part); sprite HUD with 10 hearts + 10 stamina segments; water/larva hazards, fall damage, checkpoints, respawn. VERIFIED (code read, tests)
- Fight one enemy type (succubus) with patrol/chase/attack AI; it hits only when your body is inside its attack box, placed from the level file by grid marker; it shows 2 health bars above its head and regenerates 0.1 bar/s. Your attacks hurt it (J/I 0.5 bar, jump + I 1.5 and its spin hits behind you too, K 0.5 + 1, L 1 + 2, O 1 + 2); K, L and O stop halfway unless pressed again (O after sprite 9); U sends a ground wave along the floor that hurts it too. Attacks move the player (`internal/core/attack_step.go`): J steps 45 px back, I and U 45 px forward, K's second part 15 px and O's 45 px forward; L doesn't move. Taking damage makes the player blink red for 1 s; a hit of 2+ hearts knocks them down 1.5 tiles back, 3+ hearts plays the full knockdown and 2 tiles; at 0 HP one of two death animations plays before respawn. F3 shows body/attack hitboxes for both. VERIFIED (`internal/core/enemy*.go` tracked, tests pass)
- On a touch screen, buttons for every key: arrows bottom-left, attack symbols and a bigger sprint button bottom-right, a menu button top-left (the HUD moves right of it). Idle buttons are half see-through and stay out of the middle third of the screen. VERIFIED (tests; web build with simulated touches; not on a real device)
- `GameState` (menu/playing/game over) is declared and set to `ModeMenu`, but never read: there is no menu, so the touch menu button and Escape only toggle `InputState.Menu`. VERIFIED (`cmd/main.go`)

## How to build, run, test

- Fast check: `node scripts/verify.mjs` · Full (adds native + WASM build): `node scripts/verify.mjs --full`
- Docs check: `node scripts/docs-check.mjs` (part of the fast check)
- Run native: `go run ./cmd` · Web: build wasm then `go run ./server` → http://localhost:8080 (see `DEPLOY.md`)
- `make` targets exist (`Makefile`) but `make` may be missing on Windows; use the commands above.

## Test baseline (at `a2c2e8e`, Windows, go1.26.7)

| Command | Result | Failing tests / issues (by name) |
|---|---|---|
| `node scripts/verify.mjs --full` | ✅ pass, all 8 steps | — |
| `go test ./...` | ✅ 76/76 pass, all in `internal/core` | — |
| `golangci-lint run ./...` (advisory) | 1 issue | SA9003 `internal/core/quadtree.go:217` |
| GitHub CI | see the checks on PR #16 / `main` | — |

Compare runs **by test name**: `go test ./internal/core -v | grep "^--- "` lists them. Added since the last
baseline: the touch tests (`TestTouch*`), the O combo tests (`TestO*`), `TestAttackSteps`,
`TestAirHeavyAttackHitsBothSides`, `TestAttackDamageInBars`, `TestComboAttacksHitTwicePerSwing` and the
enemy health / level marker tests.

## Known issues and risks

1. `game.exe` is tracked in git although `.gitignore` has `*.exe`. VERIFIED (`git ls-files`)
2. `go.mod` says `go 1.24.3` and CI uses 1.24; local toolchain is 1.26.7. Language features newer than 1.24 would pass locally and fail CI. VERIFIED
3. `readAsset` resolves by base file name — same-named assets in different folders silently collide (`internal/core/image.go:20`). VERIFIED
4. Large files: `animation.go` 434 lines (`InitPlayerAnimations` is most of it), `player.go` 397, `world.go` 343. VERIFIED (`wc -l`)
5. ~~`Layout` returns the outside window size, so the logical resolution changes with the window/browser size (`cmd/main.go:100`).~~ **FIXED** (`d6e602f`) — `Layout` now returns the fixed `screenWidth × screenHeight` constants; Ebitengine auto-scales to any device.
6. No tests for player physics with collision, quadtree or the server; combat, stamina, attack steps (without a level), touch input and HUD math are covered. VERIFIED (test names)
7. `DEPLOY.md` "Continuous integration" section may drift from the workflows. ASSUMED
8. ~~`game.wasm` served without `Cache-Control` header — browser caches the old binary across deploys, preventing fixes from reaching users.~~ **FIXED** (`f63b7e7`) — server now sends `no-cache` for `/game.wasm` alongside `index.html`.
9. `assets/UI/UIButtons.png` is 3.1 MB and embedded in the WASM, though touch controls use only three of its sprites; cropping it would shrink the phone download. VERIFIED (`du`)
10. The ground J sprite slashes behind the player, and J now steps back into it. Issue not filed. VERIFIED (sprite sheet)
11. The full knockdown (`KnockedOut`) still draws from row 17 while the short one uses row 16. ASSUMED (not rechecked this session)

## Open decisions

| # | Decision needed | Options | Owner | Blocking? |
|---|---|---|---|---|
| 1 | Untrack `game.exe` | `git rm --cached game.exe` / keep | user | no |

## Next task

- Not chosen yet. Candidates: draw a real menu for the menu button; crop `UIButtons.png`; touch buttons for jump (Space), F, R, Q; file the J sprite issue; untrack `game.exe`; tests for player physics with collision.
