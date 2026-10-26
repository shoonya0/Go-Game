# Main menu: state

## Now (updated 2026-10-27 at `41faa1a`)

- **Status:** 🟡 increment 4 (How to Play) on `feature/controls-help`, committed locally, not pushed
- **Current task:** task 12: the user reads the How to Play pages; then push and PR (ask before merge)
- **Stopped at:** How to Play done (`41faa1a`): Settings → How to Play, 3 pages (moves, powers, rules), version v1.0. Checked in the web build (agent-browser, all 3 pages, no console errors); not in the native window
- **Open questions / blockers:** (user) does "V1.0" mean only the version shown, or also a git tag / release? The rules page says nothing about enemies guarding or reach: that lands with `fix/30-enemy-playtest-bugs` (enemy-types D11); add a line once both are merged
- **Don't forget:** Main Menu builds the title from the run in memory (`toTitle`), never from
  storage, so blocked storage loses nothing until the tab closes. The travel menu's
  `WaitRelease` is armed in one place: whenever `Update` leaves the menus. The title screen builds the save's game behind it, so a save that fails
  `Continue` counts as no save (Continue grayed). In agent-browser, focus the canvas by script
  (`canvas.tabIndex=0; canvas.focus()`): a click lands on the middle menu item. `text/v2`
  made the WASM ~5 MB bigger (28 → 33 MB, shaping library; the unused M+ font is dropped by
  the linker).

## Recent sessions (newest first, keep 5)

### 2026-10-27: How to Play (increment 4)

- Changed: `help_screen.go` (new), Settings has How to Play between the volume and Back, `DrawText` (left-aligned)
- Verified: 5 `TestHelp*` tests (the build failed first); web build screenshots; `verify.mjs --full` 8/8
- Docs touched: README scope, spec (Settings), tasks 11-12, state, registry row

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
