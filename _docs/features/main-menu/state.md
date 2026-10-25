# Main menu: state

## Now (updated 2026-10-25 at `aaa849c`)

- **Status:** 🟡 in progress
- **Current task:** 8 (settings screen) in [tasks.md](./tasks.md)
- **Stopped at:** increment 2 done: the game opens on the title screen (`core.Screens`,
  `cmd/main.go` `openTitle` / `startRun`). Checked on the web build: title without and with a
  save, Continue after reload, New Game confirm + wipe (screenshots in the session scratchpad,
  not kept). Settings does nothing yet. Branch `feature/main-menu`
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
