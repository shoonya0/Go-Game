---
paths:
  - internal/core/enemy*.go
  - internal/core/level_enemies*.go
  - assets/enemy/**
---
# Feature: enemy-types

Docs for this feature live in `_docs/features/enemy-types/`.

- Before changing code here, read `_docs/features/enemy-types/state.md` (where work stopped)
  and `tasks.md` (what's next), then follow its `README.md` only as far as the task needs.
- After any piece of work here, save the state: `/feature-docs update`. Link facts; don't copy them.
- A new species is data: a blueprint plus a Mover and an Attacker. Don't add `switch kind` to shared code (decisions D1).
- Every enemy asset needs a unique base file name (`readAsset` looks files up by base name; decisions D4).
