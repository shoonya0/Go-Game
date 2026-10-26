# UI buttons: state

## Now (updated 2026-10-26 at `8621c11`)

- **Status:** 🟡 in progress: tasks 1–4 done, 5 coded and tested but not seen on screen
- **Current task:** 5 in [tasks.md](./tasks.md): see the travel menu's wooden arrows (needs 5+ active portals)
- **Stopped at:** title, settings, pause and overwrite screens checked in the web build
  (`_docs/evidence/ui-buttons-*.png`, untracked like the other evidence)
- **Open questions / blockers:** scope D2 was my call; the user can veto parts
- **Don't forget:** in agent-browser, key presses stopped reaching the canvas after a few
  screens; mouse clicks (move, down, up) kept working. The settings volume item keeps its
  text label (`volumeLabel`) for width and tests but draws only the bar. Main-menu and
  portals docs-check "stale state" warnings come from this feature touching `menu*.go` and
  `travel_menu.go`; their state didn't change

## Recent sessions (newest first, keep 5)

### 2026-10-26: menus from the sheet

- Changed: `ui_skin.go` (new), `menu.go`, `menu_screens.go`, `travel_menu.go`, `touch_controls.go`, tests
- Verified: `node scripts/verify.mjs --full` → 8/8 pass; web build screenshots
- Docs touched: tasks, state, current-state

### 2026-10-26: started

- Changed: feature docs only
- Verified: not verified
- Docs touched: README, decisions, tasks, state, registry, current-state
