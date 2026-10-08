# Test plan: <feature>

Phase: `docs/02-roadmap.md` §Phase <N> · Spec: `docs/specs/<spec>.md`

These plans are the source of truth for the test files. Each case has an ID, fixtures
and a pass bar.

## Levels

| Level | Tool | Covers |
|---|---|---|
| Unit | | Pure logic, no network |
| Integration | | Service seams, temp DB, recorded fixtures |
| E2E | | Real user flow |
| Golden / eval (AI features) | | Invariants and bands, never exact model text |
| Privacy / offline | | No sockets opened when the feature is disabled |
| Manual / live | Checklist below | What CI can't check (OS integration, hardware) |

## Cases

| ID | Level | Case | Fixture | Pass bar |
|---|---|---|---|---|
| T-1 | unit | | | |
| T-2 | e2e | | | |

## Fixtures

- Location: `test/fixtures/<feature>/`. Synthetic or public data only; recorded network
  responses, no live calls.

## Manual checklist (not in CI)

- [ ] On <OS/device>: …

## Rule for AI or LLM outputs

Assert structure and invariants (required sections, non-empty citations, scores within
a band, nothing fabricated). For deterministic checks, stub the model with a fixed
responder.
