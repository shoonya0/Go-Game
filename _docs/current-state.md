# Go-Game (Pirate Adventure): current state (snapshot 2026-10-07)

Snapshot at `91400c1` on `chore/ai-harness` (from `main`). Working tree dirty: 6 untracked
`internal/core/enemy*.go` files (enemy system WIP) + `stepsForCreatingApp/`.
Labels: **VERIFIED** = command run or code read · **ASSUMED** = not checked.

## What a user can do today

- Run a 2D side-scrolling platformer natively (Windows) or in the browser (WASM). VERIFIED (builds)
- Move (WASD), jump, attack; health/power HUD; water/larva hazards, fall damage, checkpoints, respawn. VERIFIED (code read)
- Fight one enemy type (succubus) with patrol/chase/attack AI — **only with the untracked WIP files**. VERIFIED
- `GameState` (menu/playing/game over) is declared and set to `ModeMenu`, but never read: there is no menu. VERIFIED (`cmd/main.go:45`)

## How to build, run, test

- Fast check: `node scripts/verify.mjs` · Full (adds native + WASM build): `node scripts/verify.mjs --full`
- Run native: `go run ./cmd` · Web: build wasm then `go run ./server` → http://localhost:8080 (see `DEPLOY.md`)
- `make` targets exist (`Makefile`) but `make` may be missing on Windows; use the commands above.

## Test baseline (at `91400c1` + untracked enemy WIP, Windows, Go 1.26.7)

| Command | Duration | Result | Failing tests (by name) |
|---|---|---|---|
| `go vet ./...` | 17s | ✅ pass | — |
| `go build ./...` | 5s | ✅ pass | — |
| `go test ./...` | 6s | ✅ pass (9 tests, all in `internal/core`) | — |
| `gofmt -l .` | <1s | ✅ no unformatted files | — |
| `GOOS=js GOARCH=wasm go build ./cmd` | 20s | ✅ pass | — |
| **committed HEAD alone** (`go build ./...` in a clean worktree) | — | ❌ **FAIL** | `cmd/main.go:24: undefined: core.EnemyRuntime`, `NewEnemy`, `KindSuccubus` |

Tests: `TestTileTypeFromCode`, `TestBuildLevelPlatformCount`, `TestBuildLevelCheckpoints`,
`TestBuildLevelTopBottomDetection`, `TestBuildLevelWorldCoordinates` (tracked);
`TestBlueprintRegistered`, `TestEnemyTakeDamage`, `TestEnemyPatrolTurnsAround`,
`TestEnemyDetectsChasesAndAttacks` (untracked `enemy_test.go`).
Compare new runs **by test name**.

## Known issues and risks

1. **`main` does not compile; CI has been red for the last 5 pushes** (2026-09-30 → 10-01).
   Every failure is the same pattern: code committed without the new file it calls
   (`undefined: newCombat` on 09-30, `undefined: core.EnemyRuntime` on 10-01). VERIFIED
   (`gh run list`, `gh run view --log-failed`, clean-worktree build).
   Fix: enemy files committed on `chore/ai-harness`; resolved once that branch is merged.
2. `game.exe` is tracked in git although `.gitignore` has `*.exe`. VERIFIED (`git ls-files`)
3. CI pins Go `1.24`; local toolchain is 1.26.7. Language features newer than 1.24 would pass locally and fail CI. VERIFIED
4. CI runs only on `ubuntu-latest`; Windows native build is not checked in CI. VERIFIED (`.github/workflows/ci.yml`)
5. `readAsset` resolves by base file name — same-named assets in different folders silently collide (`internal/core/image.go:20`). VERIFIED
6. Large files: `player.go` 381 lines, `animation.go` 365 (`InitPlayerAnimations` is 222 lines), `world.go` 334. VERIFIED (graph)
7. `Layout` returns the outside window size, so the logical resolution changes with the window/browser size (`cmd/main.go:100`). VERIFIED (code); visual impact ASSUMED.
8. No tests for player physics, combat, quadtree or the server. VERIFIED
9. `DEPLOY.md` "Continuous integration" section may drift from `ci.yml`. ASSUMED

## Open decisions

| # | Decision needed | Options | Owner | Blocking? |
|---|---|---|---|---|
| 1 | Commit the enemy WIP to fix `main` | commit as-is / finish first | user | yes, for green CI |
| 2 | Untrack `game.exe` | `git rm --cached game.exe` / keep | user | no |
| 3 | Align CI Go version with local | bump CI to `go-version-file: go.mod` / keep 1.24 | user | no |

## Next task

- Commit the enemy files so `main` builds, then confirm CI is green (`gh run list`).
