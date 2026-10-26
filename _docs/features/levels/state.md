# Levels: state

## Now (updated 2026-10-26 at `db13caf`)

- **Status:** 🟡 in progress
- **Current task:** 8. update docs, then the user plays deepForest night ([tasks.md](./tasks.md) "Later")
- **Stopped at:** tasks 6 and 7 done: deepForest night (`deepForest_night.json`, key `night`)
  follows day in `LevelOrder`; saving day's last unused portal opens it. Checked by tests and a
  throwaway web build that started in the night level; not played through
- **Open questions / blockers:** (user) is the night layout right; who draws levels 2–8
- **Don't forget:** climbs are at most 3 tiles (the user: the jump reaches 3–3.5); with a 3-tile
  climb the gap is at most 2 tiles, 2-tile climb 3, flat 4. The first portal in reading order
  (top row first) is a sub-level's entry: keep it the topmost. Water drains 2 HP/s, larva hits
  10 per entry. `//go:embed` takes untracked files (commit a level's backgrounds with its grid);
  a new kind of file in a level folder needs its own pattern in `assets/embed.go`. Saves store
  the level name: `deepForest1` stays the day sub-level's name

## Recent sessions (newest first, keep 5)

### 2026-10-26: deepForest night, last portal opens the next level (tasks 6, 7)

- Changed: `deepForest_night.json` + background; `LevelOrder`, `World.OpenNextLevel` (`levels.go`); `savePortal`, `LevelDone` (`combat.go`); `cmd/main.go`
- Verified: `node scripts/verify.mjs --full` → 8/8; `TestLastPortalOpensNextLevel`, `TestLevelPortals`; night level screenshot
- Docs touched: levels tasks, state; current-state; portals, main-menu, enemy-types sessions

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

## Known issues

- Local builds embed untracked backgrounds that CI doesn't (see "Don't forget")
