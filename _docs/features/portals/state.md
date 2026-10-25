# Portals: state

## Now (updated 2026-10-25 at `aab21cb`)

- **Status:** 🟡 in progress; increments 1 and 2 code complete (tasks 2–6). Waiting for the user's
  look in game: tasks 1 and 3 (gold, ground, `portalLift` = -10, E/tap prompt) and enemies coming
  back after a death. Next: increment 3, task 7 (levels by name) in [tasks.md](./tasks.md)
- **Stopped at:** `RespawnEnemies` (level_enemies.go) on `PlayerRuntime.Respawns` changing (main.go):
  fresh non-boss enemies at their spawns, living bosses kept, killed bosses stay dead. No boss kind yet.
  Before that: `updatePortals(portals, in)` (moved out of `updateCombat` into `UpdatePlayer`) saves
  only on `InputState.Interact` (E) or a `Tapped` screen point inside `Portal.tapArea` (doorway + half
  a tile). `TouchControls.Started` turns a just-pressed touch off every button into a one-frame tap.
  `PortalSprite.DrawPrompt`: icons, not text (user's choice 2026-10-25): gold button with E or a tap
  ring over the portal the player stands in, a fading check after a save
- **Open questions / blockers:** D3 (E saves and opens the menu in one press) to confirm in play.
  `portal.png` is 3.6 MB but drawn 150 px tall (asset change, ask the user).
  `_asset_backups/` (untracked) is the user's own folder: never commit it
- **Don't forget:** the portal sheet's cells drift, so frames are cropped to their art, not to the
  grid. Code 9 is non-solid. The game has no font: UI symbols are vector icons (`touch_icons.go`)

## Recent sessions (newest first, keep 5)

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

### 2026-10-24: user's look in game; floating portal fixed

- Changed: `Portal.FloorDrop` (world.go), test `TestPortalStandsOnTheGround`; probed the sheet (rock missing in 7 frames' art)
- Verified: new test failed first (base at 60, ground at 86.6), then `verify.mjs` 5/5
- Docs touched: state
