# Levels: state

## Now (updated 2026-10-26 at `72c83e5`)

- **Status:** 🟡 in progress
- **Current task:** 6. day and night sub-levels in order, deepForest night first ([tasks.md](./tasks.md))
- **Stopped at:** task 5 done: `LevelDef.Tiles` / `KeepApart`, `TestLevelTiles`. deepForest's
  larva pit (rows 31–35, columns 48–83) is now water (the user's call); not looked at in game yet
- **Open questions / blockers:** none. The deepForest night grid needs content: the user's
  layout, or a first draft from me (ask before task 6)
- **Don't forget:** water drains 2 HP/s and cancels fall damage, larva hits 10 HP per entry, so
  the old pit is gentler now. `//go:embed` takes untracked files: the exe carries all ~110 MB
  of backgrounds (the web build none) while CI and Render get only committed ones; commit a
  level's backgrounds with its grid. A new kind of file in a level folder needs its own pattern
  in `assets/embed.go` (`levels/*.png levels/*/*.json`). Saves store the level name
  (`deepForest1`): renaming it for the day/night split breaks old saves' checkpoint

## Recent sessions (newest first, keep 5)

### 2026-10-26: each level's tiles checked (task 5)

- Changed: `LevelDef.Tiles`, `KeepApart` (`levels.go`); `TestLevelTiles`, `TestTilesTouch`; `level_1.json` larva → water
- Verified: `node scripts/verify.mjs` → 5/5; `TestLevelTiles` failed first (113 larva cells), passes after
- Docs touched: levels tasks, state; current-state

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

## Known issues

- Local builds embed untracked backgrounds that CI doesn't (see "Don't forget")
