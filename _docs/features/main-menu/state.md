# Main menu: state

## Now (updated 2026-10-26 at `6a8cc90`)

- **Status:** ✅ done, merged to `main` in PR #25 (`6a8cc90`), branch deleted. All tasks 1–10 done
- **Current task:** none. Leftovers are in [tasks.md](./tasks.md) "Later" (real audio, the run
  rebuilt on each title visit)
- **Stopped at:** feature complete; web build checked in agent-browser. The native window was
  not checked by hand before the merge
- **Open questions / blockers:** D4 (confirm before wiping) and D7 (pause menu) were my calls;
  the user can still veto
- **Don't forget:** Main Menu builds the title from the run in memory (`toTitle`), never from
  storage, so blocked storage loses nothing until the tab closes. The travel menu's
  `WaitRelease` is armed in one place: whenever `Update` leaves the menus. The title screen builds the save's game behind it, so a save that fails
  `Continue` counts as no save (Continue grayed). In agent-browser, focus the canvas by script
  (`canvas.tabIndex=0; canvas.focus()`): a click lands on the middle menu item. `text/v2`
  made the WASM ~5 MB bigger (28 → 33 MB, shaping library; the unused M+ font is dropped by
  the linker).

## Recent sessions (newest first, keep 5)

### 2026-10-26: touched by levels (PR #28 merged)

- Changed: `cmd/main.go`: the camera follows a teleport at once (`followPlayer`), level size from `World.Level`; old saves open the next level ([levels D5](../levels/decisions.md)); menus unchanged
- Verified: `node scripts/verify.mjs --full` 8/8; PR #28 checks green; the camera change by build only
- Docs touched: this state, tasks (second level done)

### 2026-10-26: touched by levels (travel menu list)

- Changed: `cmd/main.go` no longer captures travel thumbnails; menus unchanged ([levels](../levels/state.md))
- Verified: `node scripts/verify.mjs --full` → 8/8; web build screenshot `_docs/evidence/levels-travel-list.png`
- Docs touched: this state

### 2026-10-26: touched by levels (thumbnail fix)

- Changed: `cmd/main.go` captures portal thumbnails before opening the next level ([levels](../levels/state.md)); menus unchanged
- Verified: `node scripts/verify.mjs --full` → 8/8
- Docs touched: this entry only

### 2026-10-26: touched by levels (tasks 6, 7)

- Changed: `cmd/main.go` opens the next level after the player's update ([levels](../levels/state.md)); this feature's behaviour unchanged
- Verified: `node scripts/verify.mjs --full` → 8/8
- Docs touched: this entry only

### 2026-10-26: touched by levels (task 4)

- Changed: `cmd/main.go`: `syncBackground` now downloads the background on the web and runs every tick ([levels](../levels/state.md)); this feature's behaviour unchanged
- Verified: `node scripts/verify.mjs --full` → 8/8
- Docs touched: this entry only
