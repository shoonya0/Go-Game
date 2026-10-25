# Portals: state

## Now (updated 2026-10-25 at `b0190ce`)

- **Status:** 🟡 in progress; tasks 7, 8a, 8b, 8d, 9 done. 8c (menu thumbnails) and the looks
  wait for the user's look in game. Next: user's look, then task 10 in [tasks.md](./tasks.md)
- **Stopped at:** `SwitchLevel` runs a visited level's left enemies through `RespawnEnemies`
  (non-bosses fresh, bosses as left); `TestEnteringLevelRevivesEnemies`. Travel inside the
  current level does not revive (it isn't entering). Looks by `PlayerRuntime.PortalLook`;
  `TravelAsked()` after a 0.5 s hold opens `TravelMenu`; `World.Travel` / `World.Respawned` (levels.go)
- **Open questions / blockers:** web check: synthetic key events don't reach the game in
  agent-browser, so the menu was not seen there. Only one level exists, so cross-level travel
  is tested with `testSecondLevel` only. `portal.png` is 3.6 MB (ask the user). `_asset_backups/`
  is the user's: never commit it
- **Don't forget:** Esc both closes the menu (`Back`) and toggles `InputState.Menu` (no menu
  reads it yet). A thumbnail is taken once, so it shows the portal as it looked then. Level
  sizes come from `LevelDef`, not the grid. Code 9 is non-solid. No font: UI symbols are vector
  icons (`touch_icons.go`)

## Recent sessions (newest first, keep 5)

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

### 2026-10-25: tasks 5–6, enemies come back

- Changed: `EnemyBlueprint.Boss`, `RespawnEnemies`, `PlayerRuntime.Respawns`, main.go revives on change
- Verified: `TestRespawnRevivesEnemies`, `TestKilledBossStaysDead` failed first (build); `verify.mjs --full` 8/8
- Docs touched: tasks, state, enemy-types state

### 2026-10-25: task 4, respawn at the saved portal

- Changed: `TestRespawnAtSavedPortal` only (behavior came with `00454fe`)
- Verified: passes; fails when the respawn point is mutated back to the cell corner; `verify.mjs` 5/5
- Docs touched: tasks, state
