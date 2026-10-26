# Main menu: state

## Now (updated 2026-10-27 at `56e3e22`, V1.0)

- **Status:** 🟢 increment 4 (How to Play) merged to `main` (PR #36); the game is tagged v1.0, the last push for now
- **Current task:** none; candidates in tasks.md (real audio, reuse the World's levels)
- **Stopped at:** the user checked How to Play in game; "V1.0" means the version on the board and a `v1.0` git tag
- **Open questions / blockers:** the rules page says nothing about enemies guarding or reach (enemy-types D11, now merged): add a line in the next version
- **Don't forget:** Main Menu builds the title from the run in memory (`toTitle`), never from
  storage, so blocked storage loses nothing until the tab closes. The travel menu's
  `WaitRelease` is armed in one place: whenever `Update` leaves the menus. The title screen builds the save's game behind it, so a save that fails
  `Continue` counts as no save (Continue grayed). In agent-browser, focus the canvas by script
  (`canvas.tabIndex=0; canvas.focus()`): a click lands on the middle menu item. `text/v2`
  made the WASM ~5 MB bigger (28 → 33 MB, shaping library; the unused M+ font is dropped by
  the linker).

## Recent sessions (newest first, keep 5)

### 2026-10-27: How to Play merged, v1.0

- Changed: nothing in code; PR #36 merged, README added on `chore/v1.0`
- Verified: `verify.mjs --full` 8/8; PR #36 CI green (build, verify on Ubuntu and Windows)
- Docs touched: tasks 12, this state

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
