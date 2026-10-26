# Levels: state

## Now (updated 2026-10-26 at `19a79b3`)

- **Status:** 🟡 in progress; increment 1 tasks 1–10 done, level 2 started
- **Current task:** the user plays deepForest night (changed) and aquaticAbyss day with the lower jump; then the other 13 grids
- **Stopped at:** the jump lifts the feet 3.13 tiles (`JumpForce` 605, was 700 = 4.18 tiles;
  `3ef78d6`). `TestLevelPortalsReachEachOther` simulates the real jump and checks every portal
  reaches every other: deepForest night was one-way (4-tile climb below Portal 1's ledges, 6-tile
  wall right of the pillar), fixed by moving the ledge at row 16 to row 15 and adding three
  stair ledges and a step (`19a79b3`). deepForest day and aquaticAbyss day passed unchanged
- **Open questions / blockers:** (user) are the layouts right with the lower jump; enemy lists per
  level (D3). In the web check an enemy at deepForest night Portal 2 hit the player on arrival
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

### 2026-10-26: saved portal's thumbnail before the next level

- Changed: `cmd/main.go` captures thumbnails before `OpenNextLevel` (the user saw no picture for the last saved portal)
- Verified: throwaway web build starting at day's last portal + E: `HasThumb` false before the fix, true after; screenshot `_docs/evidence/levels-thumb-after.png`; `verify.mjs --full` 8/8
- Docs touched: levels state; main-menu session

## Known issues

- Local builds embed untracked backgrounds that CI doesn't (see "Don't forget")
