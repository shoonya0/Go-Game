# 05: Taking over an existing app with AI (steps 0 to 8)

An existing codebase is harder for an agent than a new one. Martin Fowler notes that
"the harness is most needed where it is hardest to build", meaning legacy code without
types, tests or clear module boundaries. So the order is: **make it readable** for the
agent, **make it verifiable**, and only then **change it**.

```text
0 Safety → 1 Tools → 2 Map the code → 3 Write down the state → 4 Find/build the checks
  → 5 Lock behavior with tests → 6 Add the harness → 7 First small task → 8 Grow (06 loop)
```

---

## Step 0: Safety first

- Work on a **branch**, never directly on `main`/`master`. Confirm `git status` is
  clean.
- Check that `.env` and other secret files are **gitignored** and not already in
  history (`git log --all -- .env`). Never paste secrets into prompts.
- Keep the agent in Manual or Auto permission mode, with **no** blanket skipping of
  permissions. Set up allowlists (`/permissions`) for safe commands such as test and
  build.
- Note anything irreversible (production databases, deploy scripts, payment keys) and
  block those paths with a `PreToolUse` hook *before* letting the agent explore freely.

---

## Step 1: Install the toolchain

Install [`02-apps-and-tools.md`](./02-apps-and-tools.md) Tier 0 and Tier 1, plus the
LSP plugin for the repo's main language. Then, in the repo:

```bash
code-review-graph install
code-review-graph build                 # full parse; later runs are incremental
code-review-graph status
```

Optional, for very large or unfamiliar repos: `code-review-graph embed` for semantic
search, and `code-review-graph visualize` for an interactive HTML map.

Run the **claude-code-setup** plugin once. It analyzes the codebase and recommends
hooks, skills, MCP servers and subagents that fit it.

---

## Step 2: Map the codebase (read-only)

Start in **plan mode** so nothing gets edited:

```text
Plan mode. Use the code-review-graph MCP tools, not broad file reading:
get_architecture_overview, list_communities, get_hub_nodes (most-connected = riskiest to change),
get_bridge_nodes (chokepoints), list_flows (main execution paths), find_large_functions,
get_knowledge_gaps (untested/undocumented areas), get_suggested_questions.
Then read ONLY the entry points and the top 5 hub files.
Report: what the app does, its modules and how they depend on each other, the 3-5 main flows,
where state and secrets live, how it's built/run/tested, and the 10 riskiest areas, with file:line.
```

Then run `/init` to get a draft CLAUDE.md. **Cut it hard.** Keep only commands, real
conventions and gotchas, as a map shaped like
[`templates/CLAUDE.template.md`](./templates/CLAUDE.template.md). Delete anything the agent can work out
by reading the code.

For legacy or very messy code, also consider:

- `/modernize` (**code-modernization** plugin): assessment, interactive map, business
  rules mined from the code.
- `ln-22-current-architecture-documenter` (architecture-suite).

---

## Step 3: Write down the current state

Ask the agent to produce these, then **you review** them, because they become the
agent's memory of the project:

| File | Contents |
|---|---|
| `docs/README.md` | Index: read-in-this-order table |
| `docs/architecture-overview.md` | Module map, flows, layers and the import directions you *want* (even if the code doesn't follow them yet) |
| `docs/current-state.md` | What works, what's broken, known issues, how to build/run/test, the **test baseline** (step 4), and open decisions |
| `docs/glossary.md` *(optional)* | Domain terms; the agent confuses domain words surprisingly often |

```text
Write docs/current-state.md (templates/docs/current-state.md) and
docs/architecture-overview.md from what you found. Mark every claim as VERIFIED (you ran it
or read the code) or ASSUMED. List open questions for me at the end.
```

> Docs drift. A Natively memory note records a doc that said 24,000 where the code said
> 120,000. Put the commit SHA on every snapshot, and when the doc and the code disagree,
> **fix the doc**.

---

## Step 4: Find or build the checks, and record a baseline

```text
Find every way this repo is built, type-checked, linted and tested (package scripts,
Makefile, CI workflows, README). Run each one. Report: command, duration, pass/fail, and the
exact list of failing tests. Then propose ONE `verify` command (fast: typecheck+lint+unit)
and ONE `verify:full` (adds integration/e2e/build). Don't fix failures yet.
```

- **Record the baseline:** the list of failing tests at commit `<sha>`, in
  `docs/current-state.md`. From now on, compare against this list.
- **Never trust a remembered count** of "N pre-existing failures". Rerun at a clean
  HEAD (`git stash -u` → rebuild → test → `git stash pop` → rebuild). This lesson is
  recorded in Natively's memory notes.
- Check that **CI actually runs**. Look for jobs > 0 and triggers that match your real
  default branch. Natively's `build-smoke.yml` comments describe a misplaced `schedule:`
  key that made GitHub reject the workflow, so every run had **zero jobs** and no PR was
  actually smoke-tested. (The
  same repo still has triggers for `main` while its default branch is `master`; see 09.)

---

## Step 5: Lock in current behavior before changing it

Before changing an area that has no tests, add **characterization tests**, also
called "golden master" tests. They record what the code does *today*, right or wrong,
so a refactor can't silently change behavior.

```text
Before touching <area>, write characterization tests that pin its CURRENT behavior
(including odd behavior; mark suspected bugs with a TODO and a note, don't fix them).
Use the graph's tests_for / callers_of to find existing coverage first. Run them; all must pass
on the current code.
```

For a bug, the first artifact is always a **failing test that reproduces it**.

---

## Step 6: Add the harness, gradually

Turn on the harness one step at a time, so you don't drown in red:

1. **Hooks:** format on edit, graph update, protected files. These have no downside.
2. **Single `verify` command** and a **Stop hook** that runs it, *only once
   `verify` is green on the baseline* (or have it ignore the recorded baseline
   failures).
3. **CI on push and PR**, using the same `verify`.
4. **Lint and architecture rules in warn mode.** Fix the hot spots, then switch them
   to error.
5. **Pre-commit** for fast gates.
6. **Project skills** (`/verify-change`, `/close-increment`) and a `code-reviewer` subagent.

Put hooks in the **project** `.claude/settings.json` (committed), not only in your
user settings, so every clone, teammate and cloud session gets them.

---

## Step 7: First task, small and reversible

Pick a real but low-risk task: a bug with a clear reproduction, or a small UI change.
Run the full loop from [`06`](./06-feature-loop.md). Afterward, ask yourself:

- Did the agent find the right files quickly? If not, improve the docs index or add
  path rules.
- Did it make a mistake you had to correct? If so, add **one line** to CLAUDE.md, or a
  hook if the rule must never be broken.
- Did verification catch the problems, or did you? If you did, add the missing check.

---

## Step 8: Grow

- Use the 06 loop for every change.
- Over time, refactor the graph's **hub nodes** and **large functions**, each behind
  characterization tests.
- Keep `docs/current-state.md` updated in the same commit as the change.
- Weekly cleanup (07 §8).

---

## Special situations

| Situation | What to do |
|---|---|
| **No tests at all** | Don't start with features. Spend the first sessions on a test harness plus characterization tests for the 3 main flows. |
| **Huge monorepo** | Start Claude in the package folder; put a CLAUDE.md in each package (subfolder CLAUDE.md files load on demand); use `.claude/rules/` with `paths:`; rely on graph and LSP lookups, not reading |
| **Several platforms (e.g. macOS + Windows)** | Find **both** implementations before editing; inject the platform so tests can exercise both branches; CI matrix; report what was physically tested. Natively's CLAUDE.md "Cross-Platform Development Contract" is a strong template. |
| **External services** | Record fixtures (recorded responses) so tests never hit the network; gate live tests behind an env flag |
| **LLM or AI features inside the app** | Assert **invariants, not exact text** (07 §5); golden sets plus judge evals; a fixed responder stub for deterministic tests |
| **No docs and the original authors are gone** | Have the agent mine the business rules from the code (`/modernize`), then confirm each rule with a test before trusting it |
| **Private submodules / premium code** | Document what's absent locally and how builds behave without it; Natively's CI shows a fork-safe "core only" path |
