# Main menu: state

## Now (updated 2026-10-25 at `0d5bd4b`)

- **Status:** ✅ done on the branch, not pushed: all tasks 1–10 ticked
- **Current task:** none. Next: `/code-review` (touches `cmd`, core, system, ~900 lines), then
  push + PR; merging deploys to Render, so ask the user first
- **Stopped at:** web build checked in agent-browser (title with and without a save, Continue
  after reload, New Game confirm + wipe, settings kept after reload, pause → Main Menu);
  screenshots `_docs/evidence/main-menu-title.png`, `main-menu-pause.png` (untracked, like the
  portals evidence). The native window was not checked by hand
- **Open questions / blockers:** D4 (confirm before wiping) and D7 (pause menu) are my calls;
  the user can veto
- **Don't forget:** the title screen builds the save's game behind it, so a save that fails
  `Continue` counts as no save (Continue grayed). In agent-browser, focus the canvas by script
  (`canvas.tabIndex=0; canvas.focus()`): a click lands on the middle menu item. `text/v2`
  made the WASM ~5 MB bigger (28 → 33 MB, shaping library; the unused M+ font is dropped by
  the linker). The travel menu only has thumbnails for portals of levels visited this run

## Recent sessions (newest first, keep 5)

### 2026-10-25: tasks 1–10

- Changed: `internal/core/save*.go`, `settings*.go`, `menu*.go`, `text*.go`, `internal/system/save_*.go`, `cmd/main.go`, input
- Verified: `node scripts/verify.mjs --full` → 8/8 pass; web build in agent-browser (see Now)
- Docs touched: spec, decisions D2, tasks, state, architecture overview, current-state, portals state/tasks

### 2026-10-25: feature started

- Changed: docs only (`_docs/features/main-menu/`)
- Verified: `node scripts/docs-check.mjs`
- Docs touched: spec, decisions, tasks, state, registry, Active work
