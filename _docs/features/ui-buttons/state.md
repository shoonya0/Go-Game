# UI buttons: state

## Now (updated 2026-10-26 at `8c93d3b`)

- **Status:** ✅ done, merged to `main` in PR #27 (`8c93d3b`), branch deleted. Tasks 1–6 done
- **Current task:** none. Leftovers in [tasks.md](./tasks.md) "Later" (native window check, crop the sheet)
- **Stopped at:** feature complete; the user checked it in game; web build screenshots in
  `_docs/evidence/ui-buttons-*.png` (untracked like the other evidence)
- **Open questions / blockers:** scope D2 was my call; the user can veto parts
- **Don't forget:** in agent-browser, key presses stopped reaching the canvas after a few
  screens; mouse clicks (move, down, up) kept working. The settings volume item keeps its
  text label (`volumeLabel`) for width and tests but draws only the bar. Main-menu and
  portals docs-check "stale state" warnings come from this feature touching `menu*.go` and
  `travel_menu.go`; their state didn't change. PR #27's ubuntu `verify` failed on exactly those
  warnings (`docs-check --strict` in CI): a commit message with `[docs-skip]` avoids it

## Recent sessions (newest first, keep 5)

### 2026-10-26: PR

- Changed: docs only; branch pushed, PR opened
- Verified: `node scripts/verify.mjs --full` → 8/8 pass; the user checked it in game
- Docs touched: tasks, state, current-state

### 2026-10-26: menus from the sheet

- Changed: `ui_skin.go` (new), `menu.go`, `menu_screens.go`, `travel_menu.go`, `touch_controls.go`, tests
- Verified: `node scripts/verify.mjs --full` → 8/8 pass; web build screenshots
- Docs touched: tasks, state, current-state

### 2026-10-26: started

- Changed: feature docs only
- Verified: not verified
- Docs touched: README, decisions, tasks, state, registry, current-state
