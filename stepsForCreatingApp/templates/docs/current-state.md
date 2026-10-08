# <Project>: current state (snapshot YYYY-MM-DD)

Snapshot taken at `<sha>` on `<branch>`, working tree <clean | dirty: …>.
One page. Details live in each phase's status doc. Update this page when a phase or
increment lands.

## At a glance

| Phase | State | Status doc |
|---|---|---|
| 0 Foundations and harness | ✅ / 🟡 / ⏳ / ⚠️ | `status/phase-0.md` |
| 1 … | | |

## What a user can do today

- …

## How to build, run, test

- Setup: `…`
- Dev: `…` · Agent UI instance: `…`
- Fast check: `npm run verify` · Full: `npm run verify:full`

## Test baseline (at `<sha>`)

| Suite | Result | Known failures (by name) |
|---|---|---|
| unit | 412/414 | `foo › handles X` (issue #12), `bar › Y` (flaky, #15) |

Compare new runs against this list **by test name**. Rebuild it from a clean HEAD if
in doubt.

## Open decisions

| # | Decision needed | Options | Owner | Blocking? |
|---|---|---|---|---|
| | | | | |

## Known issues and risks

- …

## Next task

- Phase <N>, Increment <M>: <one line> (spec: `docs/specs/…`)
