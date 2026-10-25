# Portals: state

## Now (updated 2026-10-25 at `f1401fb`)

- **Status:** 🟢 tasks 1–9 done, the user checked the game ("everything fine"). D10 and D11 built: the
  portal travelled to becomes the checkpoint, and every teleport revives the level's enemies. Left: task 10's push + PR, and ask before merge
- **Stopped at:** `World.Travel` calls `activatePortal` on the destination and `RespawnEnemies`
  (`TestTravelMakesDestinationCheckpoint`, `TestTravelRevivesEnemies`); `updatePortals` ignores input while dying
- **Open questions / blockers:** only one level exists: cross-level travel is tested with
  `testSecondLevel` only. `_asset_backups/` is the user's: never commit it
- **Don't forget:** Esc both closes the menu (`Back`) and toggles `InputState.Menu` (no menu
  reads it yet). A thumbnail is taken once, so it shows the portal as it looked then. Level
  sizes come from `LevelDef`, not the grid. Code 9 is non-solid. No font: UI symbols are vector
  icons (`touch_icons.go`)

## Recent sessions (newest first, keep 5)

### 2026-10-25: D11, every teleport revives enemies; level edits

- Changed: `World.Travel` revives the level's non-bosses; D11, spec; user's level edits committed (portal moved, Smasher added)
- Verified: `TestTravelRevivesEnemies` failed first (6 of 7), then passes; `verify.mjs --full` 8/8
- Docs touched: decisions, spec, tasks, state, current-state

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
