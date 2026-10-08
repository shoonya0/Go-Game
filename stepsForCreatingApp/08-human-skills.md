# 08: Skills you still need when AI writes the code

"Only AI writes code" doesn't mean you need no skills. The skills **change**. You
move from writing code to directing, checking and judging it. In every case studied,
the humans stayed responsible for intent, review and the checks the code must pass.

Each skill below lists why it matters, the signs you're missing it, the minimum level
to reach, and how to learn it with the agent itself. While learning, turn on the
`explanatory-output-style` or `learning-output-style` plugin, and ask for ASCII
diagrams or HTML explainers. The Claude Code team uses both.

---

## Priority 1: you can't work without these

### 1. Turning ideas into checkable requirements

- **Why:** the agent builds exactly what you describe. Vague intent produces wrong code
  that still looks plausible.
- **Missing it looks like:** "make it better", features without acceptance criteria,
  scope that keeps growing.
- **Minimum:** write a spec of a page or less with acceptance criteria someone else
  could check, explicit non-goals, and an end-to-end verification step.
- **Learn:** have the agent interview you (04 stage 1), then ask it to criticize your
  spec: "which criteria can't be tested?"

### 2. Breaking work down

- **Why:** small, separately verifiable steps are what keep AI work reliable.
- **Missing it looks like:** 2,000-line diffs, "phase 1 = everything".
- **Minimum:** split a feature into increments that each leave the build green and can
  be reviewed in about 15 minutes.
- **Learn:** ask for 3 different ways to split the work and compare them.

### 3. Reading diffs and evidence (not writing code)

- **Why:** you are the final reviewer, and comprehension debt builds up quietly.
- **Minimum:** for any diff, you can say what changed, why, what could break, and
  whether the tests really cover it.
- **Red flags to look for in AI diffs:**
  - tests deleted, skipped, `.only`'d, or their assertions weakened
  - `any`, type casts, `// @ts-ignore`, `eslint-disable`
  - `catch {}` blocks that swallow errors, or fallbacks that hide failures
  - new dependencies, especially ones you haven't heard of
  - changes outside the task's scope; renamed or moved files nobody asked for
  - hardcoded values, URLs, keys or test data in production paths
  - mocks that replace the very thing under test
  - "should work", "production ready" or "all tests pass" with no output shown
  - a big rewrite when a small change was asked for
- **Learn:** "Explain this diff to me as a reviewer: intent, risk per file, what's
  untested."

### 4. Testing literacy

- **Why:** tests are the agent's sensors. Bad tests mean blind sensors.
- **Minimum:** know unit, integration and e2e tests apart; spot tests that only test
  mocks; understand fixtures, golden tests and flaky tests; know "assert invariants,
  not exact text" for AI outputs.
- **Learn:** "For this module, which behaviors are untested? Which existing tests
  would still pass if the code were wrong?" (A mutation-testing mindset.)

### 5. Git fluency

- **Why:** git is your undo, your audit trail and your parallelism.
- **Minimum:** branch, commit, diff, log, restore, revert, stash, cherry-pick, bisect,
  worktree; reading a PR; resolving a simple conflict.
- **Learn:** use a scratch repo and ask the agent to set up scenarios ("make a merge
  conflict for me to resolve").

### 6. Managing context and sessions

- **Why:** context is the scarce resource (01 §3).
- **Minimum:** know when to `/clear`, `/compact` or `/rewind`; what belongs in
  CLAUDE.md versus a skill versus a hook; how to name and resume sessions; how to read
  `/context`.

---

## Priority 2: needed before you ship

### 7. Architecture basics

- **Why:** the agent will couple layers and duplicate logic unless the design is clear
  and enforced.
- **Minimum:** layers and dependency direction, modules and boundaries, where state
  lives, data flow, sync vs async, how errors move. Be able to draw your app's
  architecture from memory.
- **Learn:** "Draw the architecture of this repo as ASCII. Where are the boundary
  violations?" Use the graph tools (`get_architecture_overview`, hub and bridge nodes).

### 8. Debugging and reading evidence

- **Why:** when the agent gets stuck, you supply the evidence: logs, screenshots,
  repro steps.
- **Minimum:** read a stack trace; reproduce reliably; narrow down by halves (bisect);
  read app and CI logs; tell a symptom from its cause.

### 9. Security and privacy basics

- **Why:** a large share of AI-generated code has security flaws (Veracode 2025: about
  45%), and AI invents package names.
- **Minimum:** OWASP Top 10 basics (injection, auth, access control, SSRF, XSS);
  secrets handling; least privilege; dependency hygiene; what personal data you store
  and why.

### 10. Your target platform

- **Why:** packaging, signing, permissions, OS integrations and app-store rules don't
  show up in dev mode.
- **Minimum:** for each target (web, mobile, desktop on Windows/macOS), know how it's
  built, signed, installed, updated and permitted, and what "works in dev" fails to
  prove.

---

## Priority 3: needed to scale

### 11. Building and tuning the harness

Write hooks, skills, subagents, lint rules and CI. Notice when a correction should
become a rule. The plugins `hookify`, `skill-creator`, `claude-md-management` and
`claude-code-setup` help.

### 12. Cost and time judgment

Choose models and effort levels (`/effort`; e.g. Opus for planning and hard bugs, a
faster model for routine or subagent work). Read `session-report` output. Know when a
human decision saves an hour of agent wandering.

### 13. Product judgment and honesty

Decide what *not* to build. Keep an "honest ceiling". Report only what was verified.
This is the part no harness can do for you. Fowler calls it the human's "social
accountability".

---

## A 4-week learning plan (while building something real)

| Week | Focus | Practice |
|---|---|---|
| 1 | Skills 1, 5, 6 | Build a tiny app with the 04 flow; write the spec yourself; commit every increment |
| 2 | Skills 3, 4 | Review every diff with the red-flag list; ask "which tests would still pass if this were wrong?" |
| 3 | Skills 2, 7, 8 | Take over a small open-source repo with the 05 flow; draw its architecture; fix one real bug with a repro test |
| 4 | Skills 9, 10, 11 | Add a Stop hook, a security review and one architecture rule; ship a packaged build and test it on a clean machine |

**Warning (Addy Osmani):** don't build loops *to avoid thinking*. "Build the loop. But
build it like someone who intends to stay the engineer, not just the person who presses
go."
