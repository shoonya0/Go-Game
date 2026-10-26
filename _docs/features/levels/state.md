# Levels: state

## Now (updated 2026-10-26 at `c475613`)

- **Status:** 🟡 in progress
- **Current task:** 5. each level's tiles checked by a test ([tasks.md](./tasks.md))
- **Stopped at:** task 4 done and committed (code-reviewer: nothing blocking; its two findings,
  a loading text and a retry after a failed download, are fixed). The Dockerfile change is
  untested (Docker daemon off; tasks "Later")
- **Open questions / blockers:** none. D1–D5 settled (the user accepted the D5 rules)
- **Don't forget:** `//go:embed levels` embeds untracked files too, so local builds carry all
  ~110 MB of backgrounds in the exe (the web build none) while CI and Render carry only
  `deepForest_day.png`. A new kind of file in a level folder needs its own pattern in
  `assets/embed.go` (`levels/*.png levels/*/*.json`). The old
  `LBackground_1.png` is now `moltenCrater_day.png` (same bytes). Saves store the level name
  (`deepForest1`): renaming it for the day/night split breaks old saves' checkpoint

## Recent sessions (newest first, keep 5)

### 2026-10-26: backgrounds downloaded per level on the web (task 4)

- Changed: `assets/embed.go` + `backgrounds_*.go` (split embed), `internal/core/image.go` (`HasAsset`, `AssetURL`, `DecodeImage`), `internal/system/fetch_*.go`, `cmd/main.go` (`syncBackground` async), `server/main.go` (`/levels/`), `Dockerfile`, `DEPLOY.md`
- Verified: `node scripts/verify.mjs --full` → 8/8; agent-browser: background downloaded and drawn, "Loading..." + 2 s retries with the files missing, recovered without reload (`_docs/evidence/levels-bg-*.png`); `game.wasm` 33.1 → 27.6 MB. Docker build not run (daemon off)
- Docs touched: levels state; DEPLOY.md; architecture-overview Assets / Web server rows

### 2026-10-26: JPEG decoder (task 3)

- Changed: `internal/core/image.go` imports `image/jpeg`; `TestBackgroundsDecode` in `image_test.go`
- Verified: `node scripts/verify.mjs` → 5/5 pass; the test decodes all 16 backgrounds locally
- Docs touched: levels tasks, state

### 2026-10-26: the user's level rules (D1–D5)

- Changed: docs only: decisions.md D1–D5, spec rules, tasks 2–8
- Verified: `node scripts/docs-check.mjs` → 0 errors
- Docs touched: levels decisions, spec, tasks, state; enemy-types D10; current-state

### 2026-10-26: unique background names (task 1)

- Changed: `assets/levels/*/background/*` renamed; `Background_1` (`internal/core/world.go`), `levels_test.go`
- Verified: `node scripts/verify.mjs` → passed (5 steps); `TestAssetBaseNamesUnique`, `TestAssetPathsAreEmbedded` pass
- Docs touched: levels tasks, state; current-state

### 2026-10-26: feature started

- Changed: branch renamed `docs/ui-buttons-merged` → `feature/levels`; docs folder created
- Verified: `go test ./internal/core -run TestAssetBaseNamesUnique` → FAIL (day/night name clashes)
- Docs touched: levels README, spec, tasks, state; registry; current-state Active work

## Known issues

- Local builds embed untracked backgrounds that CI doesn't (see "Don't forget")
