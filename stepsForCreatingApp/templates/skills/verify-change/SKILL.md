---
name: verify-change
description: Run the project's verification ladder (typecheck, lint, unit, and when relevant integration, e2e, agent UI check, build) and report evidence in a fixed format. Use before claiming any change is done, after implementing a plan step, or when the user asks "does it work?".
---

# Verify change

Prove the current change works. Report **evidence, not claims**.

## 1. Scope

- `git status` and `git diff --stat`: what changed?
- Choose the rungs needed. Rungs 1 to 3 always run. Add the others when the change
  touches them:

| Rung | Command (edit for this project) | Run when |
|---|---|---|
| 1-3 Fast | `npm run verify` (typecheck + lint + unit) | Always |
| 4 Integration | `npm run test:integration` | Cross-module, DB or IO changes |
| 5 E2E | `npm run test:e2e` | A user flow changed |
| 6 Agent UI | `npm run dev:agent`, then agent-browser (below) | Any UI change |
| 7 Build | `npm run build` (or the packaged build) | Build config, deps, assets, paths |

## 2. Run

- Run each chosen rung. On failure, find the **root cause** and fix it, then rerun
  from rung 1.
- Never delete, skip or weaken a test to get green. If a test is wrong, stop and
  explain why.
- Compare failures against the **baseline** in `docs/current-state.md` by test name.
  A pre-existing failure is reported as such, with evidence. It is never silently
  ignored.

## 3. Agent UI check (rung 6)

```text
agent-browser skills get core          # once per session (and "electron" for Electron apps)
agent-browser --cdp <port> tab         # port from agent-browser.json; never --auto-connect
agent-browser snapshot -i
# perform the flow from the spec: click @eN / fill @eN "…"; re-snapshot after every UI change
agent-browser console
agent-browser errors
agent-browser screenshot docs/evidence/<feature>-<step>.png
```

Compare against the spec or design and list the differences. CDP covers the renderer
only. Native OS behavior needs a human check, so list it under "Requires physical
verification".

## 4. Report (exact format)

```markdown
### Verification
| Rung | Command | Result | Notes |
|---|---|---|---|
| Fast | `npm run verify` | ✅ pass (41 s) | 412/412 unit |
| E2E | `npm run test:e2e` | ✅ pass | 18/18 |
| UI | agent-browser | ✅ | flow X ok; console clean; docs/evidence/x-1.png |

- New failures vs baseline @<sha>: none
- Not verified: <list, or "nothing">
- Requires physical verification on <OS/device>: <list, or "nothing">
```

If any rung fails and can't be fixed in scope, stop and report it. Don't claim done.
