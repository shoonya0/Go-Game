# Building apps with AI writing the code: start here

This folder covers how to build a **new** application, or take over an **existing**
one, when an AI agent writes nearly all of the code. Your job is to set the goal,
build the checks and judge the results.

It is separate from Natively's product work. Natively appears only as a worked example
of an AI-built codebase (see [`09-natively-reference.md`](./09-natively-reference.md)).

> Snapshot date: **2026-10-06**. Tool names, commands and plugin names were checked
> against official docs and this machine's plugin marketplace on that date. They
> change fast. Re-check anything marked *(verify)* before you rely on it.

---

## The method in one picture

```text
                 YOU: intent, acceptance criteria, judgment, final approval
                  │
                  ▼
 GUIDES (before the agent acts)          SENSORS (after the agent acts)
 CLAUDE.md map, docs/ specs,             typecheck, lint, unit/integration tests,
 skills, path rules, code graph,   ──►   e2e + UI snapshots, architecture rules,
 plans with exit criteria       AGENT    review agent in a fresh context, CI
                                writes    │
                                code  ◄───┘ fail → agent fixes → re-run
                  │
                  ▼
 EVIDENCE (test output, screenshots, command log) → you review → atomic commit → CI
                  │
                  ▼
 LESSONS → CLAUDE.md / hook / skill / lint rule, so the same mistake can't recur
```

In short, the agent is the model plus the setup around it. You build and maintain
that setup, which this folder calls the **harness**. The AI writes the code.

## The ten rules

1. **Close the loop.** Every task needs a check the agent can run that returns pass or fail.
   Nothing is "done" without evidence.
2. **Context is the scarce resource.** Keep the entry file short, load detail on demand,
   run `/clear` between tasks, and send exploration to subagents.
3. **If it isn't in the repo, it doesn't exist.** Specs, decisions, status and lessons
   are versioned markdown.
4. **Plan before code** when the change touches several files or the approach is
   uncertain. Skip planning for a change you can describe in one sentence.
5. **Hooks and CI are law; CLAUDE.md is advice.** Anything that must always happen
   becomes a hook or a CI gate.
6. **Small increments.** Each one is specified, verified, reviewed and committed on its own.
7. **Separate the maker from the checker.** Review in a fresh context.
8. **Fix the environment, not just the output.** When the agent fails, add the missing
   doc, tool, rule or test, then rerun.
9. **Fight entropy every week.** Remove dead code and duplication, fix docs that drifted
   from the code, and spend about 20% of the time refactoring.
10. **Stay the engineer.** Read the diffs and evidence. Never claim more verification
    than was actually done.

## Reading order

| # | File | Read it when |
|---|------|--------------|
| 01 | [`01-core-principles.md`](./01-core-principles.md) | First. The ideas behind every step, including where the experts disagree. |
| 02 | [`02-apps-and-tools.md`](./02-apps-and-tools.md) | Before you start. What to install, in tiers, with commands for Windows and macOS. |
| 03 | [`03-skills-plugins-and-extensions.md`](./03-skills-plugins-and-extensions.md) | Setting up Claude Code: skills, plugins, subagents and hooks, and which to use when. |
| 04 | [`04-new-app-step-by-step.md`](./04-new-app-step-by-step.md) | **Creating a new app.** Stages 0 to 8, with prompts. |
| 05 | [`05-existing-app-step-by-step.md`](./05-existing-app-step-by-step.md) | **Taking over an existing app.** Steps 0 to 8, with prompts. |
| 06 | [`06-feature-loop.md`](./06-feature-loop.md) | **Every day.** The per-feature and per-bug loop plus a prompt library. |
| 07 | [`07-verification-and-quality-gates.md`](./07-verification-and-quality-gates.md) | Designing tests, gates, UI checks, LLM-feature evals and cleanup. |
| 08 | [`08-human-skills.md`](./08-human-skills.md) | What **you** still need to know, and how to learn it with AI. |
| 09 | [`09-natively-reference.md`](./09-natively-reference.md) | A real AI-built repo: what it does well, its gaps, and lessons from its history. |
| 10 | [`10-checklists.md`](./10-checklists.md) | Printable checklists for setup, sessions, merges and weekly upkeep. |
| — | [`templates/`](./templates/README.md) | Files to copy into a new project: CLAUDE.md, hooks, skills and doc skeletons. |
| — | [`agent/`](./agent/README.md) | **Let an AI agent do the setup.** Step-by-step instructions it follows for a new project or an existing codebase, stopping to ask before every install, login, overwrite, commit or push. |
| — | [`sources.md`](./sources.md) | Every source, with notes on how reliable each claim is. |
| — | [`previous-research.md`](./previous-research.md) | Earlier research: products built mostly by AI, and the ranked list of techniques. |

## Fastest path

- **Have an agent do it:** open Claude Code in the target project and follow
  [`agent/README.md`](./agent/README.md). It handles both a new project and an existing
  codebase, and it stops for every install (e.g. code-review-graph) to let you decide.
- **New app, by hand:** install the Tier 0 and Tier 1 tools from 02, then follow 04
  stages 1 to 5, copy `templates/`, and use 06 for every feature.
- **Existing app, by hand:** install Tier 0 and Tier 1 from 02, follow 05 steps 0 to 6,
  then use 06 for every change.

## How this was researched

- Official Claude Code docs (best practices, memory, features overview, hooks, GitHub
  Actions), fetched 2026-10-06.
- Practitioner reports: Boris Cherny and the Claude Code team, Peter Steinberger
  (OpenClaw), OpenAI's harness-engineering write-up, Martin Fowler and Addy Osmani.
- Tool READMEs: code-review-graph, agent-browser, GitHub Spec Kit and superpowers.
- The official plugin marketplace as installed on this machine
  (`~/.claude/plugins/marketplaces/claude-plugins-official`).
- A read-through of this repository's own AI setup (CLAUDE.md, `_docs/`, `docs/`,
  hooks, CI, `scripts/dev-agent.mjs` and the memory notes).

Companies and creators report their own numbers ("90% written by AI", "1M lines",
"2 to 3× quality"). No one has audited them. Treat them as direction, not proof.
