# Portals: state

## Now (updated 2026-10-25 at `10b430d`)

- **Status:** ✅ done, merged to `main` in PR #22 (`10b430d`), branch deleted. All tasks 1–10 done
- **Stopped at:** feature complete. Leftovers are in [tasks.md](./tasks.md) "Later" (saved portals
  across restarts, a second level, one wrong camera frame after a cross-level travel, the sparkle mark)
- **Open questions / blockers:** none. `_asset_backups/` is the user's: never commit it
- **Don't forget:** Esc both closes the travel menu (`Back`) and sets `InputState.Menu`, which
  pauses (main-menu D7); `cmd` skips the pause while the travel menu is open. `WaitRelease`
  keeps the press that picked a title menu item from opening the travel menu. A thumbnail is taken once, so it shows the portal as it looked then. Level
  sizes come from `LevelDef`, not the grid. Code 9 is non-solid. No font: UI symbols are vector
  icons (`touch_icons.go`)

## Recent sessions (newest first, keep 5)

### 2026-10-26: level 1 background renamed by levels

- Changed: `Background_1` in `world.go` now `deepForest_day.png` ([levels](../levels/state.md)); portals unchanged
- Verified: `node scripts/verify.mjs` → 5/5 pass
- Docs touched: this entry only

### 2026-10-26: reskinned by ui-buttons

- Changed: the travel menu's scroll arrows now drawn from `assets/UI/UIButtons.png` ([ui-buttons](../ui-buttons/state.md), PR #27); behaviour unchanged
- Verified: `node scripts/verify.mjs --full` → 8/8 pass (on the ui-buttons branch)
- Docs touched: this entry only

### 2026-10-25: merged (PR #22)

- Changed: nothing in code; the user merged PR #22, local `main` synced, merged branches deleted
- Verified: `verify.mjs` 5/5 on `main` at `10b430d`
- Docs touched: state, tasks, registry, current-state

### 2026-10-25: D11, every teleport revives enemies; level edits

- Changed: `World.Travel` revives the level's non-bosses; D11, spec; user's level edits committed (portal moved, Smasher added)
- Verified: `TestTravelRevivesEnemies` failed first (6 of 7), then passes; `verify.mjs --full` 8/8
- Docs touched: decisions, spec, tasks, state, current-state

### 2026-10-25: D10, travel sets the checkpoint; user's look done

- Changed: `World.Travel` makes the destination the checkpoint; D10, spec; tasks 1, 3, 8c ticked (user: "everything fine")
- Verified: `TestTravelMakesDestinationCheckpoint` failed first, then passes; `verify.mjs` 5/5
- Docs touched: decisions, spec, tasks, state, current-state
