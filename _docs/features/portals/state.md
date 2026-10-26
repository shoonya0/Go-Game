# Portals: state

## Now (updated 2026-10-25 at `10b430d`)

- **Status:** ✅ done, merged to `main` in PR #22 (`10b430d`), branch deleted. All tasks 1–10 done
- **Stopped at:** feature complete. Leftovers are in [tasks.md](./tasks.md) "Later" (saved portals
  across restarts, a second level, one wrong camera frame after a cross-level travel, the sparkle mark)
- **Open questions / blockers:** none. `_asset_backups/` is the user's: never commit it
- **Don't forget:** Esc both closes the travel menu (`Back`) and sets `InputState.Menu`, which
  pauses (main-menu D7); `cmd` skips the pause while the travel menu is open. `WaitRelease`
  keeps the press that picked a title menu item from opening the travel menu. Level
  sizes come from `LevelDef`, not the grid. Code 9 is non-solid. No font: UI symbols are vector
  icons (`touch_icons.go`)

## Recent sessions (newest first, keep 5)

### 2026-10-26: touched by levels (aquaticAbyss day)

- Changed: `levels.go` registers a new level ([levels](../levels/state.md)); portal code unchanged
- Verified: `node scripts/verify.mjs --full` → 8/8
- Docs touched: this state

### 2026-10-26: touched by levels (travel menu list)

- Changed: the travel menu lists portals by level instead of thumbnails (D9 superseded by levels D7) ([levels](../levels/state.md))
- Verified: `node scripts/verify.mjs --full` → 8/8; web build screenshot `_docs/evidence/levels-travel-list.png`
- Docs touched: this state

### 2026-10-26: touched by levels (task 9)

- Changed: `world.go` / `levels.go`: an invisible `Border` wall around every level ([levels](../levels/state.md)); portals unchanged
- Verified: `node scripts/verify.mjs --full` → 8/8
- Docs touched: this entry only

### 2026-10-26: touched by levels (tasks 6, 7)

- Changed: saving a level's last unused portal opens the next level (`combat.go` `savePortal`, `levels.go` `OpenNextLevel`) ([levels](../levels/state.md)); this feature's behaviour unchanged
- Verified: `node scripts/verify.mjs --full` → 8/8
- Docs touched: this entry only

### 2026-10-26: `LevelDef` gained tile lists (levels task 5)

- Changed: `levels.go` `LevelDef.Tiles`, `KeepApart` ([levels](../levels/state.md)); portals unchanged
- Verified: `node scripts/verify.mjs` → 5/5
- Docs touched: this entry only
