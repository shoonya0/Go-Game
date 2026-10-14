---
name: close-increment
description: Close out a finished increment by updating _docs/status/increments.md, the feature's tasks.md and state.md, and _docs/current-state.md, recording lessons, writing the completion report, and making an atomic commit. Invoke manually after /verify-change and review pass.
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

1. `_docs/status/increments.md`: add a section, newest first:
   - **Problem**: one or two sentences.
   - **Change**: files and modules, and behavior before vs after.
   - **Tests**: what was added, with test names.
   - **Evidence**: the verification table from `/verify-change`.
   - **Commit**: write `<pending>`; fill in the hash in the next increment's commit.

   If the file doesn't exist yet, create it and turn its plain-text row in
   `_docs/README.md` into a link (docs-check flags an unlinked doc as an orphan).
2. **Feature docs** (if the work belongs to a feature with its own docs folder): tick the
   finished items in its `tasks.md`, update each leaf doc whose facts changed, then the
   feature's `state.md` (overwrite "Now"; prepend one 3-line session entry linking this
   increment; keep the newest 5), then its registry row in `_docs/README.md` only if
   status or summary changed. `/feature-docs update` describes the same steps.
3. `_docs/current-state.md`: always update `## Active work` (feature, its `state.md`,
   next step). Update the rest only if the project's state changed (what a user can do,
   known issues, open decisions, baseline, next task). Update the snapshot SHA and date.
4. If a doc you read during this work turned out to be **wrong** (including
   `_docs/architecture-overview.md` or `DEPLOY.md`), fix it now and name it in the report.
5. Run `node scripts/docs-check.mjs`. Fix every ERROR, and every "stale state" or
   "too big" warning about the files this increment touched.

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
- Commit on the feature branch, never on `main`. Message: `<type>(<scope>): <what>`, a
  blank line, then *why*, then `Fixes #N` if an issue exists. No AI attribution
  trailers; the author is the GitHub user (CLAUDE.md "Git").
- Verify: `git log -1 --format='%an <%ae>%n%B'`, `git show --stat HEAD`, `git status --short`.
- Don't push or open a PR unless the user asked to. Never push `main` (Render auto-deploys).
