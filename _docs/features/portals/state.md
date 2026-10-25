# Portals: state

## Now (updated 2026-10-25 at `4976039`)

- **Status:** 🟢 tasks 1–9 done, the user checked the game ("everything fine"). D10 built: the
  portal travelled to becomes the checkpoint. Left: task 10's push + PR, and ask before merge
- **Stopped at:** `World.Travel` calls `activatePortal` on the destination
  (`TestTravelMakesDestinationCheckpoint`); `updatePortals` ignores input while dying
- **Open questions / blockers:** only one level exists: cross-level travel is tested with
  `testSecondLevel` only. `_asset_backups/` is the user's: never commit it
- **Don't forget:** Esc both closes the menu (`Back`) and toggles `InputState.Menu` (no menu
  reads it yet). A thumbnail is taken once, so it shows the portal as it looked then. Level
  sizes come from `LevelDef`, not the grid. Code 9 is non-solid. No font: UI symbols are vector
  icons (`touch_icons.go`)

## Recent sessions (newest first, keep 5)

### 2026-10-25: D10, travel sets the checkpoint; user's look done

- Changed: `World.Travel` makes the destination the checkpoint; D10, spec; tasks 1, 3, 8c ticked (user: "everything fine")
- Verified: `TestTravelMakesDestinationCheckpoint` failed first, then passes; `verify.mjs` 5/5
- Docs touched: decisions, spec, tasks, state, current-state

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
