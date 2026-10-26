# Levels: state

## Now (updated 2026-10-26 at `21820a4`)

- **Status:** 🟡 in progress; increment 1 tasks 1–9 done
- **Current task:** the user plays deepForest night ([tasks.md](./tasks.md) "Later"); then grids for levels 2–8
- **Stopped at:** docs brought up to date (task 8): spec still matches the code; current-state
  lists the levels capability and a new test baseline (200/200 at `21820a4`, `verify.mjs --full` 8/8).
  deepForest night follows day in `LevelOrder` and opens when day's last unused portal is saved;
  checked by tests and a web build, not played through
- **Open questions / blockers:** (user) is the night layout right; who draws levels 2–8
- **Don't forget:** climbs are at most 3 tiles (the user: the jump reaches 3–3.5); with a 3-tile
  climb the gap is at most 2 tiles, 2-tile climb 3, flat 4. The first portal in reading order
  (top row first) is a sub-level's entry: keep it the topmost. Water drains 2 HP/s, larva hits
  10 per entry. `//go:embed` takes untracked files (commit a level's backgrounds with its grid);
  a new kind of file in a level folder needs its own pattern in `assets/embed.go`. Saves store
  the level name: `deepForest1` stays the day sub-level's name. On the web, a thumbnail is skipped
  while the level's background still downloads; if the player leaves first, it's taken on the
  next visit

## Recent sessions (newest first, keep 5)

### 2026-10-26: docs updated (task 8)

- Changed: levels tasks, state; `_docs/current-state.md` (levels capability, test baseline, next task)
- Verified: `node scripts/verify.mjs --full` → 8/8; `go test ./...` 200/200
- Docs touched: levels tasks, state; current-state

### 2026-10-26: saved portal's thumbnail before the next level

- Changed: `cmd/main.go` captures thumbnails before `OpenNextLevel` (the user saw no picture for the last saved portal)
- Verified: throwaway web build starting at day's last portal + E: `HasThumb` false before the fix, true after; screenshot `_docs/evidence/levels-thumb-after.png`; `verify.mjs --full` 8/8
- Docs touched: levels state; main-menu session

### 2026-10-26: invisible border around every level (task 9, D6)

- Changed: `Border` tile type, `levelBorder`, `Level.Width/Height` (`world.go`); `World.border` (`levels.go`); `TestLevelBorderKeepsPlayerIn`
- Verified: `node scripts/verify.mjs --full` → 8/8; the test failed first (player walked out at x = -1.1)
- Docs touched: levels decisions D6, spec, tasks, state; portals session

### 2026-10-26: deepForest night, last portal opens the next level (tasks 6, 7)

- Changed: `deepForest_night.json` + background; `LevelOrder`, `World.OpenNextLevel` (`levels.go`); `savePortal`, `LevelDone` (`combat.go`); `cmd/main.go`
- Verified: `node scripts/verify.mjs --full` → 8/8; `TestLastPortalOpensNextLevel`, `TestLevelPortals`; night level screenshot
- Docs touched: levels tasks, state; current-state; portals, main-menu, enemy-types sessions

### 2026-10-26: each level's tiles checked (task 5)

- Changed: `LevelDef.Tiles`, `KeepApart` (`levels.go`); `TestLevelTiles`, `TestTilesTouch`; `level_1.json` larva → water
- Verified: `node scripts/verify.mjs` → 5/5; `TestLevelTiles` failed first (113 larva cells), passes after
- Docs touched: levels tasks, state; current-state

## Known issues

- Local builds embed untracked backgrounds that CI doesn't (see "Don't forget")
