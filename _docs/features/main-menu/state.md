# Main menu: state

## Now (updated 2026-10-26 at `6a8cc90`)

- **Status:** ✅ done, merged to `main` in PR #25 (`6a8cc90`), branch deleted. All tasks 1–10 done
- **Current task:** none. Leftovers are in [tasks.md](./tasks.md) "Later" (real audio, the run
  rebuilt on each title visit, a second level)
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
  the linker). The travel menu only has thumbnails for portals of levels visited this run

## Recent sessions (newest first, keep 5)

### 2026-10-26: reskinned by ui-buttons

- Changed: the menus now drawn from `assets/UI/UIButtons.png` ([ui-buttons](../ui-buttons/state.md), PR #27); behaviour unchanged
- Verified: `node scripts/verify.mjs --full` → 8/8 pass (on the ui-buttons branch)
- Docs touched: this entry only

### 2026-10-26: merged

- Changed: nothing in code; PR #25 merged by the user
- Verified: `git log origin/main` → `6a8cc90`
- Docs touched: state, registry, current-state

### 2026-10-25: code review fixes

- Changed: `menu.go`, `menu_screens.go`, `input.go`, `text.go`, `touch_controls.go`, `cmd/main.go`
- Verified: `node scripts/verify.mjs` → pass; web build: blocked storage still offers Continue
- Docs touched: state, tasks ("Later"), current-state

### 2026-10-25: tasks 1–10

- Changed: `internal/core/save*.go`, `settings*.go`, `menu*.go`, `text*.go`, `internal/system/save_*.go`, `cmd/main.go`, input
- Verified: `node scripts/verify.mjs --full` → 8/8 pass; web build in agent-browser (see Now)
- Docs touched: spec, decisions D2, tasks, state, architecture overview, current-state, portals state/tasks

### 2026-10-25: feature started

- Changed: docs only (`_docs/features/main-menu/`)
- Verified: `node scripts/docs-check.mjs`
- Docs touched: spec, decisions, tasks, state, registry, Active work
