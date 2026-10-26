# Levels: state

## Now (updated 2026-10-26 at `3c27666`)

- **Status:** 🟡 increment 1 merged to `main` (PR #28); playtest and enemy lists still open. Larva surface fix on `fix/30-enemy-playtest-bugs` (local, not pushed)
- **Current task:** the user plays the levels and says which need hand-made changes
- **Stopped at:** larva top tiles are solid from half a tile down (`3c27666`). It matches the
  user's new `assets/levels/LevelTile.png` (lava in the lower half of the larva top tile),
  which is uncommitted; the committed tileset's larva fills the whole tile
- **Open questions / blockers:** (user) commit the new `LevelTile.png` with the larva fix on
  this branch? Are the layouts fun, which levels need hand-made
  features; the 5 review findings left (tasks.md): issues or not; enemy lists per level (D3). In the web check an enemy at deepForest night Portal 2
  hit the player on arrival (new grids keep enemies 6+ columns from portals; detection is now
  600-760 px, enemy-types D11, so more enemies will notice an arriving player)
- **Don't forget:** a new grid must pass `TestLevelPortalsReachEachOther` (it prints which portal
  can't reach which). The jump climbs 3 tiles (also straight up a wall), never 4; with a 3-tile
  climb the gap is at most 3 tiles, flat 4; out of water a jump reaches ~2 tiles. Water isn't a
  fall-through: the player stands 0.5 tile under the **highest** water tile they touch, so no
  waterfalls or water in row 0. Falls up to 5 tiles are free, into water always. The first portal
  in reading order (top row first) is a sub-level's entry: keep it the topmost; portals are
  numbered in reading order in the travel menu. `//go:embed` takes untracked files (commit a
  level's backgrounds with its grid). Saves store the level name: `deepForest1` stays the day
  sub-level's name. A new level needs `Title` and `Time` in its `LevelDef`

## Recent sessions (newest first, keep 5)

### 2026-10-26: PR #28 merged

- Changed: nothing in code; local `main` fast-forwarded to `500c45f`, `feature/levels` deleted
- Verified: `node scripts/verify.mjs` 5/5 on `main`; PR checks green
- Docs touched: levels state; current-state; README registry; architecture-overview

### 2026-10-26: code review fixes for PR #28

- Changed: `levels.go` (size from the grid, `portalCells`, `nextLevel`, `OpenNextLevel`), `combat.go` `savePortal`, `cmd/main.go` `followPlayer`; 3 tests; levels D5
- Verified: the 3 new tests failed first, then pass; `go test ./...` 207 pass, 0 fail; `node scripts/verify.mjs --full` 8/8
- Docs touched: levels D5, tasks, state; current-state

### 2026-10-26: the other 13 grids, all levels connected

- Changed: 13 `assets/levels/<level>/<level>_<day|night>.json` + their backgrounds; `levels.go` (constants, `subLevel`, `LevelOrder` of 16); levels D8
- Verified: `node scripts/verify.mjs --full` → 8/8; the reach test first failed on aquaticAbyss night (water under the stairs), fixed in the generator
- Docs touched: levels D8, tasks, state; current-state

### 2026-10-26: 3-tile jump, portal reach test, night fixed

- Changed: `JumpForce` 700 → 605 (`player.go`), `player_jump_test.go`; `level_reach_test.go`; `deepForest_night.json` (ledge up a row, 3 stair ledges, a step)
- Verified: `node scripts/verify.mjs --full` → 8/8; `TestJumpReachesThreeTiles` failed first (4.18 tiles, climbed 4); the reach test failed on night with both jumps, passes after the grid fix
- Docs touched: levels tasks, state; current-state

### 2026-10-26: aquaticAbyss day grid (draft)

- Changed: `aquaticAbyss_day.json` + its background; `AquaticAbyssDay` in `Levels` and `LevelOrder` (`levels.go`)
- Verified: `node scripts/verify.mjs --full` → 8/8 (`TestLevelBorderKeepsPlayerIn` failed on a first draft with waterfalls: water in row 0); throwaway web builds: `_docs/evidence/levels-aquatic-day.png`, `levels-travel-scroll.png`
- Docs touched: levels tasks, state; current-state



## Known issues

- All 16 backgrounds are committed (~110 MB): the repo and the Docker image are big
