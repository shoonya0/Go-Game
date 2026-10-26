# Levels: state

## Now (updated 2026-10-26 at `402afe5`)

- **Status:** 🟡 in progress; increment 1 tasks 1–10 done
- **Current task:** the user plays deepForest night ([tasks.md](./tasks.md) "Later"); then grids for levels 2–8
- **Stopped at:** the travel menu is a wooden board listing portals by level and sub-level
  ("Deep Forest" › "Deep Forest Day" › "Portal 1", D7), thumbnails removed (`402afe5`). Checked
  by tests and a throwaway web build with every portal active: pick moved, travelled to night
  Portal 2. deepForest night not played through by the user
- **Open questions / blockers:** (user) is the night layout right; who draws levels 2–8. In the
  web check an enemy at night Portal 2 hit the player right after arriving: worth a look
- **Don't forget:** climbs are at most 3 tiles (the user: the jump reaches 3–3.5); with a 3-tile
  climb the gap is at most 2 tiles, 2-tile climb 3, flat 4. The first portal in reading order
  (top row first) is a sub-level's entry: keep it the topmost. A portal's menu number is its
  reading-order place in its sub-level, so moving portals in a grid renumbers them. Water drains
  2 HP/s, larva hits 10 per entry. `//go:embed` takes untracked files (commit a level's
  backgrounds with its grid); a new kind of file in a level folder needs its own pattern in
  `assets/embed.go`. Saves store the level name: `deepForest1` stays the day sub-level's name.
  A new level needs `Title` and `Time` in its `LevelDef` for the menu

## Recent sessions (newest first, keep 5)

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

### 2026-10-26: deepForest night, last portal opens the next level (tasks 6, 7)

- Changed: `deepForest_night.json` + background; `LevelOrder`, `World.OpenNextLevel` (`levels.go`); `savePortal`, `LevelDone` (`combat.go`); `cmd/main.go`
- Verified: `node scripts/verify.mjs --full` → 8/8; `TestLastPortalOpensNextLevel`, `TestLevelPortals`; night level screenshot
- Docs touched: levels tasks, state; current-state; portals, main-menu, enemy-types sessions

## Known issues

- Local builds embed untracked backgrounds that CI doesn't (see "Don't forget")
