# Levels: state

## Now (updated 2026-10-26 at `3c19bb2` + the 13 grids)

- **Status:** 🟡 in progress; all 16 sub-levels have grids and are chained in `LevelOrder`
- **Current task:** the user plays the levels (deepForest night, aquaticAbyss day, the 13 new ones)
- **Stopped at:** 13 generated grids (D8) with their backgrounds; `subLevel` / `boulderBastion`
  build their `LevelDef`s (`levels.go`). All level tests pass for all 16; `verify.mjs --full` 8/8.
  Not played in game or checked in the web build yet
- **Open questions / blockers:** (user) are the layouts fun, which levels need hand-made
  features; enemy lists per level (D3). In the web check an enemy at deepForest night Portal 2
  hit the player on arrival (new grids keep enemies 6+ columns from portals)
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

### 2026-10-26: travel menu lists portals by level (task 10, D7)

- Changed: `travel_menu.go` (board, names, planks, scrolling), `LevelDef.Title/Time`, `portalNumber` (`levels.go`); `cmd/main.go` thumbnail capture removed; `TestTravelMenuGroupsByLevel`
- Verified: `node scripts/verify.mjs --full` → 8/8; throwaway web build: list drawn, S moved the pick, Enter travelled to night Portal 2 (`_docs/evidence/levels-travel-list.png`)
- Docs touched: levels D7, tasks, state; portals D9, spec, state; main-menu state, tasks; ui-buttons state; architecture-overview; current-state

### 2026-10-26: docs updated (task 8)

- Changed: levels tasks, state; `_docs/current-state.md` (levels capability, test baseline, next task)
- Verified: `node scripts/verify.mjs --full` → 8/8; `go test ./...` 200/200
- Docs touched: levels tasks, state; current-state

## Known issues

- Local builds embed untracked backgrounds that CI doesn't (see "Don't forget")
