# 01: Core principles

These are the ideas behind every step in 04 to 07. Each one says what it means,
why it works, and how to apply it. The last section covers where experienced
practitioners disagree.

---

## 1. Agent = model + harness

**Idea.** The model is fixed. You can't change how smart it is. You *can* change
everything around it, and that surrounding setup is the **harness**:

- **Guides** act before the agent does anything: CLAUDE.md, docs, skills, path rules,
  architecture docs, plans, and the code graph.
- **Sensors** act afterward: type checker, linter, tests, UI snapshots, architecture
  rules, review agents and CI.
- **Computational** guides and sensors are deterministic and cheap: linters and tests.
- **Inferential** ones use an AI to judge meaning, such as an AI review. They cost more
  and their output varies.

**Why.** OpenAI's team built about 1M lines with Codex and said that what was hard was
not getting code written but "designing environments, feedback loops, and control
systems." Martin Fowler calls this *harness engineering*.

**Apply.** When the agent produces something wrong, first ask **"what was missing?"**
It could be a doc, a tool, a guardrail, a test or an example. Add it, then rerun. Fix
the output by hand only as a last resort.

## 2. Close the loop (the biggest lever)

**Idea.** The agent stops when the work *looks* done. Without a check it can run,
"looks done" is the only signal, so **you** become the test suite.

**Why.** Boris Cherny, who created Claude Code, calls this "probably the most important
tip" and reports 2 to 3× better output. The official best-practices page leads with
it. Made-up APIs and imaginary fixes fail the build or the tests right away, and the
agent then fixes them itself.

**Apply.**
- Every prompt names a check: "...then run `npm run verify` and fix failures."
- For UI work, take a screenshot or snapshot, compare it to the target, and list the
  differences.
- For bugs, write a failing test that reproduces the bug *first*.
- Make the gate harder when you plan to walk away. A prompt is weakest, then `/goal`
  (re-checked every turn), then a **Stop hook** (deterministic), then a review subagent
  in a fresh context, then a required CI check.
- Ask for **evidence**, not claims: the command that ran, its output, and a screenshot.

## 3. Context is the scarce resource

**Idea.** The context window holds every message, every file read and every command
output. Quality drops as it fills: the agent starts forgetting rules and making more
mistakes.

**Apply.**
- Keep the always-loaded file short. Official guidance is **under about 200 lines**.
  OpenAI kept theirs to **about 100 lines** as a table of contents.
- Use **progressive disclosure**: the entry file links to `docs/`; skills load only
  when relevant; `.claude/rules/` with `paths:` load only for matching files.
- Run `/clear` between unrelated tasks. After **two failed corrections**, `/clear` and
  write a better prompt that includes what you learned.
- Send wide exploration to **subagents**. They return a summary, not 50 files.
- Use a **code graph** (code-review-graph) or **LSP** to find symbols, callers and
  impact instead of reading whole files.
- Use `/btw` for side questions, `/context` to see what's loaded, and `/compact <focus>`
  when you need to keep going.

## 4. If it isn't in the repo, it doesn't exist

**Idea.** The agent can't see Slack, meetings, your head or last week's chat. Anything
it needs must be a **versioned file in the repo**: product spec, architecture,
decisions (ADRs), phase plans, status, test plans and lessons.

**Apply.** Use a fixed doc skeleton (see `templates/docs/`). Update status docs in the
**same commit** as the code. Keep docs honest: when a doc and the code disagree, the
code wins, so fix the doc. Natively's memory notes record a doc number (24,000) that
was stale against the code (120,000).

## 5. Plan before code, when it pays

**Idea.** Separate exploring, planning, implementing and committing. Plan mode
(`Shift+Tab`, or `claude --permission-mode plan`) lets the agent read but not edit.

**Apply.**
- Plan when you're unsure of the approach, the change spans several files, or you
  don't know the code.
- **Skip** planning when you could describe the diff in one sentence.
- For bigger features, **let the agent interview you** (AskUserQuestion) and write
  `SPEC.md`, then implement in a **fresh session**.
- Have a second agent review the plan "as a staff engineer". Cherny does this.
- A good spec names the files and interfaces involved, states what's out of scope, and
  **ends with an end-to-end verification step**.

## 6. Hooks and CI are law; CLAUDE.md is advice

**Idea.** The agent treats CLAUDE.md and skills as context it *should* follow. Hooks
run **every time**, without exception.

**Apply.** If a rule must always hold, such as "never edit `.env`", "format after
every edit" or "tests pass before stopping", turn it into a hook (`PreToolUse`,
`PostToolUse`, `Stop`) or a required CI check. If the agent keeps breaking a
CLAUDE.md rule, the file is probably too long, or the rule should be a hook.

## 7. Small, verified, atomic increments

**Idea.** One increment does one thing. It is specified, implemented, tested,
reviewed, documented and committed. The agent commits **only the files it touched**.

**Why.** A bad change is easy to find (`git bisect`) and easy to revert. Reviews stay
small enough to actually read.

**Apply.** Natively's phase docs are a good model. Each phase is split into
*Increment 1..N*, and each increment has a problem statement, its changes, its tests
and its commit hashes.

## 8. Separate the maker from the checker

**Idea.** An agent grading its own work is biased toward it. A reviewer in a **fresh
context** sees only the diff and the criteria.

**Apply.** Run `/code-review` (it reviews in a fresh subagent), use a writer session
and a separate reviewer session, or have one agent write tests and another write the
code. Tell the reviewer to report **only gaps in correctness or requirements**.
Reviewers asked to find problems always find some, and chasing all of them leads to
over-engineering.

## 9. Fight entropy continuously

**Idea.** AI-written code builds up duplication, dead code and docs that no longer
match the code, and it does so faster than human-written code. A messy codebase makes
the *next* AI output worse.

**Apply.** Steinberger spends about 20% of his time on agent-driven refactoring, using
`jscpd` for duplication, `knip` for dead code and ESLint. OpenAI used to spend every
Friday cleaning up by hand. They replaced it with scheduled cleanup agents that check
the code against "golden principles" and open small fix PRs.

## 10. Stay the engineer

**Idea.** Shipping code you didn't write widens **comprehension debt**. Building loops
so you can stop thinking, which Addy Osmani calls *cognitive surrender*, is the trap.

**Apply.** Read every diff summary and its evidence. Ask the agent to explain any part
you don't understand, with diagrams if needed. Keep an architecture you can draw from
memory. Learn the skills in [`08-human-skills.md`](./08-human-skills.md).

## 11. Report honestly

**Idea.** Keep **reviewed**, **tested automatically**, **tested by hand** and **not
verified** clearly apart.

**Apply.** Require a completion report with fixed categories. Natively's CLAUDE.md
does this with categories such as "Tested physically on Windows", "Covered by
automated macOS branch tests" and "Requires physical macOS verification". A generic
version is in [`07`](./07-verification-and-quality-gates.md#9-evidence-report).

## 12. Assume AI code is insecure until checked

**Idea.** The AI can invent package names, and attackers register those names
("slopsquatting"). A USENIX Security 2025 study found that about 20% of packages
recommended across 576k generated samples didn't exist. Veracode's 2025 report found
security flaws in about 45% of AI-generated code samples.

**Apply.** Check every new dependency yourself: does it exist, is it popular, is it
maintained? Run `/security-review` and a security plugin before you merge. Never let
secrets into the repo or into prompts.

---

## Where the experts disagree (and what to do)

| Topic | Cherny / Anthropic docs | Steinberger (OpenClaw) | OpenAI harness team | Recommendation |
|---|---|---|---|---|
| **Plan mode** | About 80% of sessions start in plan mode; one Claude plans, another reviews | Rarely uses big plan files; "start a discussion" with the agent instead | Execution plans are first-class files in the repo | Plan for multi-file or unclear work. Use a conversation for small work. Save plans for anything spanning more than one session. |
| **Subagents** | Use them to keep the main context clean | "Mostly marketing"; uses separate terminal windows | Agents review other agents' work | Use them for **exploration** and **review**. Run core implementation where you can watch it. |
| **Parallel agents** | 3 to 5 sessions, one git worktree each | 3 to 8 agents in the **same folder**, one dev server | The app can boot per worktree | Start with **one**. Add a second only once Tier 1 (checks and rules) is solid. Use worktrees only if the app can boot per worktree on its own ports. |
| **Rules file size** | Short and pruned; under 200 lines | About 800 lines of "scar tissue" | About 100 lines, used as a map | A **map of 100 to 200 lines**, plus path-scoped rules and `docs/`. Size matters less than whether each rule changes behavior. |
| **MCP vs CLI** | CLIs are the most context-efficient; MCP for systems without a CLI | CLIs over MCP (GitHub MCP cost about 23k tokens) | `gh` CLI | Prefer a **CLI** when one exists. Claude Code now loads MCP tool schemas only when needed (tool search), which shrinks MCP's idle cost. |
| **Specs** | Interview, then SPEC.md, then a fresh session | Stopped writing big specs | Product specs and design docs in `docs/` | Write a spec for each **feature**, and a short conversation for each **tweak**. |

Common ground: verification loops, a living rules file, small commits, writing tests
in the same session as the feature, regular refactoring, and humans focusing on
intent and review.
