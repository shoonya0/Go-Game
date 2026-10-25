---
paths:
  - internal/core/portal*.go
  - internal/core/world.go
  - internal/core/combat.go
---
# Feature: portals

Docs for this feature live in `_docs/features/portals/`.

- Before changing code here, read `_docs/features/portals/state.md` (where work stopped)
  and `tasks.md` (what's next), then follow its `README.md` only as far as the task needs.
- After any piece of work here, save the state: `/feature-docs update`. Link facts; don't copy them.
- Code 9 in a level grid is a portal (decisions D1); `0` is always an empty cell.
