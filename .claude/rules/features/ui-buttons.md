---
paths:
  - internal/core/ui_skin*.go
  - internal/core/menu.go
  - internal/core/menu_screens.go
  - internal/core/travel_menu.go
  - internal/core/touch_controls.go
  - assets/UI/UIButtons.png
---
# Feature: ui-buttons

Docs for this feature live in `_docs/features/ui-buttons/`.

- Before changing code here, read `_docs/features/ui-buttons/state.md` (where work stopped)
  and `tasks.md` (what's next), then follow its `README.md` only as far as the task needs.
- After any piece of work here, save the state: `/feature-docs update`. Link facts; don't copy them.
- Sprite rects in `ui_skin.go` are measured from the sheet's opaque pixels; re-measure if the sheet changes.
