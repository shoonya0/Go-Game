# Portals: state

## Now (updated 2026-10-25 at `6a3446e`)

- **Status:** 🟡 in progress; increments 1–2 code complete, increment 3 task 7 done. Still waiting
  for the user's look in game: tasks 1 and 3 (gold, ground, E/tap prompt) and enemies coming back
  after a death. Next: task 8 (travel menu) in [tasks.md](./tasks.md)
- **Stopped at:** `World` (levels.go) holds the current level (platforms, quadtree, portals,
  enemies); `Levels` names each level (file, key, background, size); `SwitchLevel(name)` rebuilds
  it and gives a level visited before its enemies back as left (task 9 turns that into
  `RespawnEnemies`, so killed bosses stay dead). main.go uses `g.world`; `syncBackground` reloads
  the background when `world.Def.Background` changes: call it after travel. Nothing switches yet
- **Open questions / blockers:** D3 (E saves and opens the menu in one press) to confirm in play.
  Task 8: the respawn point is a bare `Position`, so a saved portal needs its level name too.
  `portal.png` is 3.6 MB but drawn 150 px tall (asset change, ask the user).
  `_asset_backups/` (untracked) is the user's own folder: never commit it
- **Don't forget:** level sizes come from `LevelDef` (6000×2040 for level 1), not the grid (36×~100
  cells = 6060×2160): computing them would shift the camera. The portal sheet's cells drift, so
  frames are cropped to their art. Code 9 is non-solid. The game has no font: UI symbols are
  vector icons (`touch_icons.go`)

## Recent sessions (newest first, keep 5)

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

### 2026-10-25: golden saved portal

- Changed: `savedPortalTint` + `Draw(…, respawn, …)` (portal.go), `Portal.SavedBy`, `updatePortals` saves `StandPos`; user's `portal.png` fix committed
- Verified: `TestOnlySavedPortalIsGolden` failed first (build); `verify.mjs --full` 8/8; gold previewed offline
- Docs touched: decisions (D6), tasks, state
