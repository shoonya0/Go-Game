# 06: The daily loop (feature, bug, refactor)

Every change, in a new app or an existing one, goes through the same loop. Ten
minutes of setup per increment saves hours of cleanup later.

```text
 ┌─► 1 Pick ONE task ─► 2 Fresh context ─► 3 Orient (graph + docs, read-only)
 │        │
 │        ▼
 │   4 Specify (spec / repro) ─► 5 Plan (+ second-opinion review)
 │        │
 │        ▼
 │   6 Implement in small steps, tests in the same session
 │        │
 │        ▼
 │   7 Verify (ladder) ── fail ──► fix ──┐
 │        │ pass                         │
 │        ▼                              │
 │   8 Review in a fresh context ── gaps ┘
 │        │
 │        ▼
 │   9 Record (status doc, lessons) ─► 10 Atomic commit + PR + CI
 └──────────────────────────── 11 /clear ◄─┘
```

---

## The steps

### 1. Pick one task

Take it from the roadmap or status doc, or from an issue. Size it so the diff can be
**reviewed in about 15 minutes**. Split anything larger into increments.

### 2. Fresh context

Use `/clear`, or open a new session and `/rename feature-x`. One task per session. For
long work, resume later with `claude --continue` or `--resume`.

### 3. Orient (read-only)

```text
Task: <one line>. Don't edit anything yet.
1. Use code-review-graph (semantic_search_nodes / query_graph callers_of, callees_of, tests_for,
   get_impact_radius, get_affected_flows) to find the code involved.
2. Read the relevant docs: <docs/...> and docs/current-state.md.
3. Summarize: current behavior, files and functions involved (file:line), callers, existing tests,
   and anything in the docs that contradicts the code.
```

> A Natively memory note records why this step matters. Skipping the graph and
> `docs/` once let a real bug through: an embedding call with one unbounded batch,
> which the existing architecture doc warned against.

### 4. Specify

**Feature:**

```text
Interview me with AskUserQuestion about <feature>: behavior, edge cases, error states,
UI states, data, what's out of scope. Then write docs/specs/<date>-<feature>.md
(templates/docs/feature-spec.md). It must name files/interfaces involved, list acceptance
criteria as checkable statements, list "Review focus" edge cases, and END with an
end-to-end verification step.
```

**Bug:**

```text
Bug: <symptom>, seen when <steps>, expected <X>, got <Y>. Logs/screenshot attached.
Reproduce it with a FAILING test first and show me the failure. Then find the root
cause (not the symptom) and explain it with file:line before proposing a fix.
```

### 5. Plan

Use plan mode (`Shift+Tab`) for anything that touches several files or where the
approach is unclear.

```text
Plan mode. From docs/specs/<spec>.md, write the implementation plan: steps in order, files
to change per step, the test(s) each step adds, the command that verifies each step, global
constraints (things that must NOT change), and risks. Each step must leave the build green.
```

Get a second opinion. Paste the plan into a **new session**:

```text
Review this plan as a staff engineer against docs/01-architecture.md and the spec.
Flag: wrong layer, missing edge cases, untestable steps, scope creep, simpler alternatives.
Only issues that matter.
```

Press `Ctrl+G` to edit the plan yourself before approving it.

### 6. Implement

```text
Implement step 1 of the plan only. Write its tests in this same session. Run
`npm run verify` and show the output. Then stop.
```

Repeat for each step, or approve "steps 1 to 3" once you trust the plan. Rules:

- No new dependency without saying why, and checking that the package exists and is
  maintained.
- Don't edit tests just to make them pass. If a test is wrong, explain why first.
- Stay inside the plan. Anything else goes into a "follow-ups" list.

### 7. Verify

Follow the ladder in [`07`](./07-verification-and-quality-gates.md): typecheck, lint,
unit, integration, e2e, UI snapshot, packaged build. Stop at the highest rung that is
relevant to the change.

```text
Run /verify-change. For UI changes also: start the app (npm run dev:agent or dev server), use
agent-browser to open the screen, snapshot -i, perform the user flow from the spec, check
console and errors, take a screenshot, and compare it with the spec/design. Show evidence.
```

### 8. Review in a fresh context

```text
/code-review
```

Or use a custom subagent:

```text
Use the code-reviewer subagent to review `git diff` against docs/specs/<spec>.md.
Report only unmet requirements, correctness bugs, missing tests for listed edge cases,
and out-of-scope changes.
```

Fix the real gaps. Ignore style nitpicks; chasing every finding leads to
over-engineering.

### 9. Record

```text
/close-increment
```

This updates `docs/status/phase-N.md` (what changed, tests, and the commit placeholder)
and `docs/current-state.md` if the state changed. It adds **one** lesson to CLAUDE.md
only if a mistake happened that would happen again, and it writes the completion
report (07 §9).

### 10. Commit, PR and CI

```text
Commit only the files you changed for this increment, with a descriptive message
(what + why). Push the branch and open a PR whose body contains the completion report.
Then watch CI (gh pr checks) and fix any failure.
```

**You** read the PR: the summary, the diff stats, the risky files and the evidence.
Merge only when CI is green.

### 11. Reset

Run `/clear`. Next task.

---

## Variants

| Variant | Differences from the main loop |
|---|---|
| **Tiny change** (you can say the diff in one sentence) | Skip 4 and 5. Prompt: "do X, run verify, show output". Still commit atomically. |
| **Bug fix** | Step 4 = failing repro test; step 5 = root-cause explanation; the fix must turn the test green without changing unrelated tests |
| **Refactor** | Characterization tests first; no behavior change allowed; review with `refactor-safely` (graph skill) and `get_impact_radius` |
| **Spike / experiment** | A throwaway branch or folder, no harness rules, timeboxed; the only output is a short findings doc. Never merge spike code. |
| **UI change** | Paste a screenshot or design; verify with agent-browser snapshots and screenshots; check light/dark, small and large widths, keyboard and accessibility |
| **Dependency upgrade** | One dependency per increment; read the changelog through Context7 or the web; run `verify:full`; skill `ln-42-dependency-upgrader` |
| **Investigation (big unknown bug)** | Write the findings into `docs/<topic>-handoff/00-START-HERE.md`, root causes, a fix plan and experiments, *before* fixing (Natively's `docs/retrieval-handoff/` pattern) |
| **Parallel work** | Separate git worktrees (Claude Code has built-in worktree support; see the docs page "worktrees") and separate ports and data directories; never two agents editing the same files |

---

## When things go wrong

| Symptom | Fix |
|---|---|
| You corrected the same thing twice | `/clear`, then a better first prompt that includes what you learned. A long session full of corrections performs worse than a clean one. |
| Agent says "done" but tests fail | A Stop hook running `verify`; demand evidence in every prompt |
| It invented an API or option | Context7 for current docs; typecheck in `verify`; add the correct usage to a skill |
| It drifts from the plan | Back to plan mode; restate the constraints; smaller steps |
| It edits unrelated files | Narrower prompt with file names; review the diff; `PreToolUse` guard on sensitive paths |
| It keeps breaking a rule in CLAUDE.md | The file is too long, or the rule is vague. Shorten the file, write "IMPORTANT" on that one line, or turn the rule into a hook. |
| Context is filling up | `/compact Focus on <current task>`; send exploration to subagents; `/btw` for side questions |
| You can't tell what it changed | `git diff --stat`, `code-review-graph detect-changes --brief`, the `review-changes` skill |
| You want to undo | `Esc Esc` / `/rewind` (code and/or conversation), or `git restore` / `git revert` |

---

## Prompt library (copy and adapt)

**Start of session:**

```text
Read CLAUDE.md, docs/README.md and docs/current-state.md. Tell me in 5 lines where the
project stands and what the next task in docs/status/phase-<N>.md is. Don't edit.
```

**Explain before change:**

```text
Explain how <flow> works end to end with an ASCII sequence diagram and file:line refs.
```

**Evidence demand:**

```text
Don't tell me it works. Show me: the commands you ran, their output (tail), and a screenshot
for UI. List anything you did NOT verify.
```

**Scope guard:**

```text
Only change files under <path>. If you believe something outside it must change, stop and
ask.
```

**Lesson capture:**

```text
That mistake will happen again. Add ONE concise line to CLAUDE.md (or propose a hook if it
must never happen) so future sessions avoid it.
```

**Tricky problem** (Steinberger's trigger words):

```text
Take your time. Read all related code first. Be comprehensive. Then propose 2 approaches with
tradeoffs before implementing.
```

**Self-review before handing back:**

```text
Before you say done, re-read the spec and your diff. List each acceptance criterion with
PASS/FAIL and the evidence for each.
```
