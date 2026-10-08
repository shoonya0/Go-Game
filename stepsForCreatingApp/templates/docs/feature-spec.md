# <Feature>: spec (YYYY-MM-DD)

**Goal:** one sentence describing the user-visible outcome.
**Phase / increment:** Phase <N>, Increment <M> · **Test plan:** `docs/test/<feature>.md`

## Behavior

- When <trigger>, the system <response>.
- UI states: empty / loading / success / error / disabled.

## Acceptance criteria (each checkable)

| ID | Criterion | Verified by |
|---|---|---|
| AC1 | | unit test `…` / e2e `…` / agent UI check / manual on <OS> |

## Out of scope

- …

## Design

- Files and interfaces involved (with paths): …
- New or changed contracts (types, IPC channels, API endpoints, DB schema): …
- Layer placement (per `docs/01-architecture.md`): …

## Global constraints (must NOT change)

- e.g. "Unset settings behave exactly as today."
- e.g. "No new `process.platform` branches."

## Review focus (edge cases no happy path covers; each gets a test)

1. …
2. …

## Risks and open questions

- …

## End-to-end verification (last step of implementation)

Exact steps and commands that prove the feature works, including the agent UI flow and
any physical checks.
