# Main menu: state

## Now (updated 2026-10-25 at `aaa849c`)

- **Status:** 🟡 in progress
- **Current task:** 3 (autosave from `cmd/main.go`) in [tasks.md](./tasks.md)
- **Stopped at:** tasks 1–2 done (`internal/core/save.go`, `internal/system/save_*.go`; `TestSave*`, `TestSaveFile*` pass, `verify --full` passes). Branch `feature/main-menu`
- **Open questions / blockers:** D4 and D7 are my calls; the user can veto
- **Don't forget:** `internal/core` can't touch files or `syscall/js`: storage goes in
  `internal/system` behind build tags (like `touch_js.go`). The travel menu only captures
  thumbnails for the current level's portals, so after Continue other levels' rows have none
  until visited

## Recent sessions (newest first, keep 5)

### 2026-10-25: feature started

- Changed: docs only (`_docs/features/main-menu/`)
- Verified: `node scripts/docs-check.mjs`
- Docs touched: spec, decisions, tasks, state, registry, Active work
