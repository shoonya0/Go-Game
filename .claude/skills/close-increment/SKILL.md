---
name: close-increment
description: Close out a finished increment by updating docs/status/increments.md and docs/current-state.md, recording lessons, writing the completion report, and making an atomic commit. Invoke manually after /verify-change and review pass.
disable-model-invocation: true
---

# Close increment

Run this only after `/verify-change` passed and the `code-reviewer` findings were
addressed. `$ARGUMENTS` may name the increment, e.g. `enemy knockback`.

## 1. Confirm the gate

- `/verify-change` evidence exists from **this** state of the tree. If anything changed since,
  rerun `/verify-change`.
- Every acceptance criterion of the task is marked PASS with evidence.

## 2. Update the docs (same commit as the code)

1. `docs/status/increments.md` (create it if missing): add a section, newest first:
   - **Problem**: one or two sentences.
   - **Change**: files and modules, and behavior before vs after.
   - **Tests**: what was added, with test names.
   - **Evidence**: the verification table from `/verify-change`.
   - **Commit**: write `<pending>`; fill in the hash in the next increment's commit.
2. `docs/current-state.md`: update only if the project's state changed (what a user can
   do, known issues, open decisions, baseline, next task). Update the snapshot SHA and date.
3. If a doc you read during this work turned out to be **wrong** (including
   `docs/architecture-overview.md` or `DEPLOY.md`), fix it now and name it in the report.

## 3. Lessons

Did a mistake happen that will happen again (a wrong command, a missed convention,
the wrong layer, a WASM-only failure)?

- If yes, add **one** concise line to the Rules section of CLAUDE.md, *or* propose a hook
  if it must never happen, *or* a depguard rule in `.golangci.yml` for a layer mistake.
- If no, change nothing. Don't add generic advice.

## 4. Completion report

```markdown
### Change summary
- Files changed / behavior changed / modules affected

### Validation
- Covered by automated tests: …
- Verified by web smoke check: …
- Played natively on Windows: …
- Reviewed but not executed: …
- Requires physical verification: …

### Commands executed
- `…` → result

### Remaining risks
- …
```

## 5. Commit

- Stage **only** the files changed for this increment (`git add <paths>`, never
  `git add -A` blindly). Check with `git diff --cached --stat` and confirm no new `.go`
  file the change needs is left untracked.
- Message: `<type>(<scope>): <what>`, a blank line, then *why*, then any trailers the
  session requires.
- Don't push or open a PR unless the user asked to. Never push `main` (Render auto-deploys).
