# Main menu: state

## Now (updated 2026-10-25 at `aaa849c`)

- **Status:** 🟡 in progress
- **Current task:** 9 (pause menu) in [tasks.md](./tasks.md)
- **Stopped at:** task 8 done: Settings screen (volume 0–10, own record `settings`), checked
  on the web build (volume kept after reload). Increments 1–2 done. Branch `feature/main-menu`
- **Open questions / blockers:** D4 and D7 are my calls; the user can veto
- **Don't forget:** the title screen builds the save's game behind it, so a save that fails
  `Continue` counts as no save (Continue grayed). Focusing the canvas in agent-browser with a
  click hits the menu item in the middle: focus it by script instead. `internal/core` can't touch files or `syscall/js`: storage goes in
  `internal/system` behind build tags (like `touch_js.go`). The travel menu only captures
  thumbnails for the current level's portals, so after Continue other levels' rows have none
  until visited

## Recent sessions (newest first, keep 5)

### 2026-10-25: feature started

- Changed: docs only (`_docs/features/main-menu/`)
- Verified: `node scripts/docs-check.mjs`
- Docs touched: spec, decisions, tasks, state, registry, Active work
