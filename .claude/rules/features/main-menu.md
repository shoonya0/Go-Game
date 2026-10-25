---
paths:
  - cmd/main.go
  - internal/core/menu*.go
  - internal/core/save*.go
  - internal/core/settings*.go
  - internal/core/text*.go
  - internal/system/save*.go
---
# Feature: main-menu

Docs for this feature live in `_docs/features/main-menu/`.

- Before changing code here, read `_docs/features/main-menu/state.md` (where work stopped)
  and `tasks.md` (what's next), then follow its `README.md` only as far as the task needs.
- After any piece of work here, save the state: `/feature-docs update`. Link facts; don't copy them.
- Storage (files, `localStorage`) lives in `internal/system` behind build tags; `internal/core` only (de)serializes.
