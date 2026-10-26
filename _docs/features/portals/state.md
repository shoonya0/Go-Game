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

### 2026-10-26: touched by levels (tasks 6, 7)

- Changed: saving a level's last unused portal opens the next level (`combat.go` `savePortal`, `levels.go` `OpenNextLevel`) ([levels](../levels/state.md)); this feature's behaviour unchanged
- Verified: `node scripts/verify.mjs --full` → 8/8
- Docs touched: this entry only

### 2026-10-26: `LevelDef` gained tile lists (levels task 5)

- Changed: `levels.go` `LevelDef.Tiles`, `KeepApart` ([levels](../levels/state.md)); portals unchanged
- Verified: `node scripts/verify.mjs` → 5/5
- Docs touched: this entry only

### 2026-10-26: touched by levels (task 4)

- Changed: `cmd/main.go` `captureThumbs` waits until the background is loaded; `world.go` draws no background while it downloads ([levels](../levels/state.md)); this feature's behaviour unchanged
- Verified: `node scripts/verify.mjs --full` → 8/8
- Docs touched: this entry only

### 2026-10-26: level 1 background renamed by levels

- Changed: `Background_1` in `world.go` now `deepForest_day.png` ([levels](../levels/state.md)); portals unchanged
- Verified: `node scripts/verify.mjs` → 5/5 pass
- Docs touched: this entry only

### 2026-10-26: reskinned by ui-buttons

- Changed: the travel menu's scroll arrows now drawn from `assets/UI/UIButtons.png` ([ui-buttons](../ui-buttons/state.md), PR #27); behaviour unchanged
- Verified: `node scripts/verify.mjs --full` → 8/8 pass (on the ui-buttons branch)
- Docs touched: this entry only
