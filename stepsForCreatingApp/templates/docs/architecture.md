# 01: Architecture

> Keep it under about 250 lines. It is the map the agent must follow. Link deeper docs;
> don't inline them.

## Overview diagram

```text
[ UI ] ──calls──► [ services ] ──uses──► [ data / adapters ] ──► external
   ▲                    │
   └──── events ────────┘
```

## Modules

| Module | Path | Responsibility | Owns state? |
|---|---|---|---|
| | | | |

## Layers and allowed imports (enforced by `npm run lint`)

| From \ To | ui | services | data | platform | shared |
|---|---|---|---|---|---|
| **ui** | ✅ | ✅ | ❌ | ❌ | ✅ |
| **services** | ❌ | ✅ | ✅ | ✅ (via adapter) | ✅ |
| **data** | ❌ | ❌ | ✅ | ❌ | ✅ |

Each ❌ is a lint rule whose message points back to this table.

## Main flows

1. **<Flow name>**: step → step → step (`file` refs)
2. …

## Data and state

- Where data lives (DB, files, memory) and who may write it.
- Migrations policy.

## Secrets and security

- Where keys live (env / OS keychain), and which process can read them.
- Trust boundaries (e.g. renderer is untrusted; validate all IPC input).

## Errors, logging, observability

- Error-handling convention (typed errors? result objects?).
- Structured logging fields; where logs go; how the agent can read them.

## Platform adapters (if multi-platform)

| Capability | darwin | win32 | linux | Interface |
|---|---|---|---|---|
| | `src/platform/darwin/x.ts` | `src/platform/win32/x.ts` | n/a (throws) | `XAdapter` |

## Testing seams

- Where fakes are injected (clock, network, LLM, platform).

## Decisions

See `docs/adr/`. List the most important ADRs here, one line each.
