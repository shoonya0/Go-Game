# Levels: state

## Now (updated 2026-10-26 at `fe512b9`)

- **Status:** 🟡 in progress; increment 1 tasks 1–10 done, level 2 started
- **Current task:** the user plays deepForest night and aquaticAbyss day ([tasks.md](./tasks.md) "Later"); then the other 13 grids
- **Stopped at:** aquaticAbyss day added as a draft (`fe512b9`): a sea along the bottom (rows 29–32)
  with three grass islands, rock ledges down from the entry portal and up to Portal 2 on the
  right; Portal 3 is on the middle island. Water 336 > rock 199 > grass 33 cells. It follows
  deepForest night in `LevelOrder`. A scratch script checked that every portal reaches the
  others (climbs ≤ 3, ≤ 2 out of water); tests pass; web build drew the sea and the scrolling travel menu
- **Open questions / blockers:** (user) are the deepForest night and aquaticAbyss day layouts right;
  enemy lists per level (D3). In the web check an enemy at deepForest night Portal 2 hit the player on arrival
- **Don't forget:** climbs are at most 3 tiles (the user: the jump reaches 3–3.5); with a 3-tile
  climb the gap is at most 2 tiles, 2-tile climb 3, flat 4; out of water a jump reaches ~2 tiles.
  Water isn't a fall-through: the player stands 0.5 tile under the **highest** water tile they
  touch, so no waterfalls or water in row 0 (it held the player above the top border). Falls up
  to 5 tiles are free, into water always. The first portal in reading order (top row first) is
  a sub-level's entry: keep it the topmost; portals are numbered in reading order in the travel
  menu. `//go:embed` takes untracked files (commit a level's backgrounds with its grid). Saves
  store the level name: `deepForest1` stays the day sub-level's name. A new level needs `Title`
  and `Time` in its `LevelDef`

## Recent sessions (newest first, keep 5)

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

### 2026-10-26: invisible border around every level (task 9, D6)

- Changed: `Border` tile type, `levelBorder`, `Level.Width/Height` (`world.go`); `World.border` (`levels.go`); `TestLevelBorderKeepsPlayerIn`
- Verified: `node scripts/verify.mjs --full` → 8/8; the test failed first (player walked out at x = -1.1)
- Docs touched: levels decisions D6, spec, tasks, state; portals session

## Known issues

- Local builds embed untracked backgrounds that CI doesn't (see "Don't forget")
