# Portals: state

## Now (updated 2026-10-25 at `eb2bfa9`)

- **Status:** 🟡 in progress; tasks 1–9 code complete. Task 10: `verify.mjs --full` 8/8, web smoke
  check passed, branch review done (1 blocker, fixed in `129b540`). Left: the user's look in game
  (tasks 1, 3, 8c are manual checks), then push + PR, and ask before merge
- **Stopped at:** `updatePortals` returns early while dying (`TestDyingPlayerCannotUsePortals`).
  `SwitchLevel` revives a visited level's non-bosses (`RespawnEnemies`). `portal.png` halved to 1.0 MB
- **Open questions / blockers:** the user's look in game: three looks, E/tap prompt, hold-E menu and
  thumbnails, enemies back after a death (PC and touch). Web: synthetic key events don't reach the
  game in agent-browser, so the menu was not seen there. Only one level exists: cross-level travel is
  tested with `testSecondLevel` only. `_asset_backups/` is the user's: never commit it
- **Don't forget:** Esc both closes the menu (`Back`) and toggles `InputState.Menu` (no menu
  reads it yet). A thumbnail is taken once, so it shows the portal as it looked then. Level
  sizes come from `LevelDef`, not the grid. Code 9 is non-solid. No font: UI symbols are vector
  icons (`touch_icons.go`)

## Recent sessions (newest first, keep 5)

### 2026-10-25: task 10 checks, branch review, dying fix

- Changed: `updatePortals` ignores input while dying; `AABB.Intersects` comment restored
- Verified: `TestDyingPlayerCannotUsePortals` failed first (4 inputs), then passes; `verify.mjs --full` 8/8; web: loads, console clean, gold portal from the half-size sheet
- Docs touched: tasks, state, current-state

### 2026-10-25: smaller portal sheet

- Changed: `assets/levels/portal.png` halved (Catmull-Rom, scratch tool, no project code)
- Verified: portal tests pass; `verify.mjs --full` 8/8; sheet looked at, detail kept
- Docs touched: tasks, state

### 2026-10-25: task 9, entering a level revives its enemies

- Changed: `SwitchLevel` revives a visited level's non-bosses; `TestSwitchLevel` return check now expects all enemies (D4)
- Verified: `TestEnteringLevelRevivesEnemies` failed first (5 of 6), then passes; `verify.mjs` 5/5
- Docs touched: tasks, state, current-state

### 2026-10-25: user's redesign, looks + hold-E travel menu (8a–8d)

- Changed: D7–D9, spec; `PortalRef`, `PortalLook`, `activatePortal`; input hold ticks, touch hold; `travel_menu.go`; `World.Travel/Respawned`; main wiring
- Verified: new tests failed first (build), then pass; `verify.mjs --full` 8/8; web: loads, gold start portal, no console errors
- Docs touched: decisions, spec, tasks, state, current-state, architecture

### 2026-10-25: task 7, levels by name

- Changed: `levels.go` (`LevelDef`, `Levels`, `World`, `SwitchLevel`), main.go on `g.world`, parallax takes the level width
- Verified: `TestSwitchLevel` passes; `verify.mjs` 5/5; WASM build OK
- Docs touched: tasks, state, README scope, current-state, architecture
