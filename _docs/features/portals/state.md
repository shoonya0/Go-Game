# Portals: state

## Now (updated 2026-10-25 at `32c0bda`)

- **Status:** 🟡 in progress; user's redesign (D7–D9) built: tasks 7, 8a, 8b, 8d done, 8c (menu
  thumbnails) waits for the user's look. Next: user's look in game, then task 9 in [tasks.md](./tasks.md)
- **Stopped at:** looks by `PlayerRuntime.PortalLook` (inactive = frame 0 in dark gray, `inactivePortalTint`; active = loop,
  checkpoint = gold); E/tap `activatePortal` (adds to `ActivePortals`, sets `Checkpoint`).
  `TravelAsked()` fires once when `InputState.InteractTicks` or a `TouchHoldTicks` on the portal
  reaches `travelHoldTicks()` (0.5 s); main opens `TravelMenu` and pauses. Thumbnails: main's
  `captureThumbs` draws the scene around a newly active portal into `g.view`, `SetThumb` shrinks it.
  `World.Travel` / `World.Respawned` (levels.go) move between levels
- **Open questions / blockers:** web check: synthetic key events don't reach the game in
  agent-browser, so the menu was not seen there; the build loads and the start portal is gold.
  Task 9: `SwitchLevel` brings a level's enemies back as left; it should revive non-bosses.
  `portal.png` is 3.6 MB (ask the user). `_asset_backups/` is the user's: never commit it
- **Don't forget:** Esc both closes the menu (`Back`) and toggles `InputState.Menu` (no menu
  reads it yet). A thumbnail is taken once, so it shows the portal as it looked then. Level
  sizes come from `LevelDef`, not the grid. Code 9 is non-solid. No font: UI symbols are vector
  icons (`touch_icons.go`)

## Recent sessions (newest first, keep 5)

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

### 2026-10-25: task 3, save with E or a tap

- Changed: `InputState.Interact/Tapped/TapAt`, E in system/Input.go, touch taps, `updatePortals(…, in)`, `AABB.Contains`, `DrawPrompt` icons
- Verified: `TestPortalSavesOnlyOnInteract`, `TestTapOnPortalSaves`, `TestTouchTapOffButtons` failed first (build); `verify.mjs --full` 8/8; `code-reviewer`: no blockers, 3 notes fixed (just-pressed taps, tap area over the prompt, 2 tests)
- Docs touched: state
