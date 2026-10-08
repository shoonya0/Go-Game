# 03: Skills, plugins, subagents and hooks

Claude Code has seven ways to extend it. Picking the right one for each need matters
more than how many you install.

---

## 1. Which mechanism for which need

| Mechanism | Loads | Use it for | Example |
|---|---|---|---|
| **CLAUDE.md** (or AGENTS.md) | Every session, in full | Facts and rules needed **every** session: build/test commands, conventions, the project map | "Run `npm run verify` before saying done." |
| **`.claude/rules/*.md`** with `paths:` | Only when matching files are touched | Rules for one language or folder | `paths: ["src/api/**"]` → API rules |
| **Skill** (`.claude/skills/<name>/SKILL.md`) | Description every session; body on demand | Reusable procedures and reference knowledge; `/name` commands | `/verify-change`, `/close-increment`, an API style guide |
| **Subagent** (`.claude/agents/<name>.md`) | Fresh, isolated context | Wide exploration, review in a fresh context, parallel side tasks | `security-reviewer`, `code-reviewer` |
| **Hook** (`settings.json`) | Runs on lifecycle events; costs no context unless it prints output | Anything that must happen **every time** | Format after edit; block `.env` edits; run tests before stopping |
| **MCP server** | Tool names at start; schemas when first used | Systems the agent can't reach through a CLI | Context7 docs, the code graph, Sentry, Figma |
| **Plugin** | A bundle of the above | Reusing one setup across repos; installing from marketplaces | `superpowers`, `code-review`, `context7` |

**Rules of thumb (from the official docs):**

- Hooks are law. CLAUDE.md and skills are guidance. Subagents are isolated workers.
- Keep CLAUDE.md under about 200 lines. Move detail into skills, rules or `docs/`.
- Mark skills that have side effects (deploy, commit, publish) with
  `disable-model-invocation: true`. Then they run only when you type the command, and
  their description costs no context.
- The CLAUDE.md files at different levels (managed, user, project, local) **add
  together**. Skills, subagents and MCP servers **override by name**. Hooks from every
  source **merge and all run**.

**When to add each one** (the docs' "build your setup over time" triggers):

| You notice… | Add |
|---|---|
| The agent gets a convention or command wrong **twice** | A line in CLAUDE.md |
| You type the same starting prompt over and over | A user-invoked **skill** |
| You paste the same multi-step playbook a **third** time | A **skill** |
| You keep copying data from a browser tab into the chat | An **MCP server** (or a CLI) |
| The agent reads many files to find a symbol | An **LSP plugin** or the **code graph** |
| A side task floods the chat with output | A **subagent** |
| Something must happen every time | A **hook** |
| A second repo needs the same setup | Package it as a **plugin** |

---

## 2. Skills: how they work

```text
.claude/skills/verify-change/SKILL.md        ← project skill (commit it)
~/.claude/skills/<name>/SKILL.md      ← personal skill (all projects)
<plugin>/skills/<name>/SKILL.md       ← plugin skill, invoked as /plugin:name
```

```markdown
---
name: fix-issue
description: Fix a GitHub issue end to end with a failing test first
disable-model-invocation: true
---
Analyze and fix the GitHub issue: $ARGUMENTS.
1. `gh issue view $ARGUMENTS` …
```

- **Progressive disclosure.** At startup only `name` and `description` load, a few
  dozen tokens each. The body loads when the skill is used. Files the body links to
  (`references/`, `scripts/`) load only if needed.
- The **description is how the agent decides to use the skill**, so write it like a
  trigger: what it does, *when* to use it, and when not to.
- `$ARGUMENTS` receives whatever follows `/skill-name`.
- `context: fork` runs the skill in an isolated context. A subagent's `skills:` field
  preloads skills into that subagent.
- Write and test skills with the **skill-creator** plugin. It creates the skill, runs
  evals, and tunes the description so the skill triggers reliably.

Templates: [`templates/skills/verify-change/SKILL.md`](./templates/skills/verify-change/SKILL.md)
and [`templates/skills/close-increment/SKILL.md`](./templates/skills/close-increment/SKILL.md).

---

## 3. Recommended plugins and skills (catalog)

Install from the official marketplace with
`/plugin install <name>@claude-plugins-official`, then `/reload-plugins`. Browse with
`/plugin`. Every **plugin** name below was confirmed in this machine's copy of the
official marketplace on 2026-10-06, except where another source is named
(levnikolaevich, Spec Kit, agent-browser). Slash commands marked "built in" ship with
Claude Code.

### A. Methodology: how the work flows

| Plugin / skills | What it gives you | Use when |
|---|---|---|
| **superpowers** (obra) | Brainstorming, then **writing-plans** (tasks of 2 to 5 minutes), then **subagent-driven-development** (a fresh subagent per task, with review), plus TDD (red/green), systematic debugging, verification before completion, git worktrees, and code-review request/receive. Natively used it: see `docs/superpowers/specs/` and `docs/superpowers/plans/`. | You want a strict, opinionated process. It deliberately slows down the start of coding, a cost that pays off on non-trivial work. |
| **feature-dev** | A feature workflow with agents for codebase exploration, architecture design and quality review | Mid-sized features in an existing codebase |
| **levnikolaevich/claude-code-skills** (vendored at `_docs/claude-code-skills/`) | Standalone skills grouped by lifecycle stage: **product discovery** (ln-11 opportunity, ln-12 requirements, ln-13 interaction design), **architecture** (ln-21 to ln-26: baseline, current-architecture documenter, proposal, ADR, diagrams, migration), **delivery planning** (ln-31 plan, ln-32 test strategy, ln-33 plan review), **implementation** (ln-41 surgical change, ln-42 dependency upgrade, ln-43 modernize, ln-44 performance, ln-45 benchmark), **QA** (ln-51 acceptance tests, ln-52 delivery review, ln-53 to ln-57 audits for docs, codebase, tests, architecture and persistence), **delivery** (ln-61 to ln-64), **operations** (ln-71, ln-72) | You want one well-scoped skill per stage with a clear "done" condition. Install with `/plugin marketplace add levnikolaevich/claude-code-skills`, then `/plugin install <suite>@levnikolaevich-skills-marketplace`. |
| **GitHub Spec Kit / BMAD / OpenSpec** | A formal spec-driven pipeline (constitution → spec → plan → tasks → implement) | Larger greenfield projects, or teams that need traceability. Optional. |

> Pick **one** methodology per project. Mixing superpowers, Spec Kit and BMAD gives
> conflicting instructions and wastes context.

### B. Understanding code

| Plugin / skill | Use |
|---|---|
| **code-review-graph** skills: `explore-codebase`, `debug-issue`, `refactor-safely`, `review-changes` | Graph-first exploration, debugging, safe refactors and change review |
| **typescript-lsp** / other `*-lsp` | Symbol navigation and live type errors |
| **context7** | Current library docs |
| **code-modernization** (`/modernize`) | Legacy code: assessment, interactive map, business rules mined from the code, an approved plan, then upgrade or rewrite with proof that behavior didn't change |

### C. Review and quality

| Command / plugin | Use |
|---|---|
| `/code-review [low…max]` (built in) | Reviews the current diff for bugs in a fresh subagent. `ultra` runs a multi-agent review in the cloud, triggered by the user only. |
| `/security-review` (built in) | Security review of pending changes |
| `/simplify` (built in) | Cleanups for reuse, simplification and efficiency, then applies the fixes |
| **code-review**, **pr-review-toolkit** | Multi-agent PR review covering comments, tests, error handling, type design and simplification |
| **code-simplifier** | An agent that refines code for clarity without changing behavior |
| **security-guidance**, **claude-security**, **semgrep**, **sonarqube** | Security gates (see 02, Tier 4) |
| **coderabbit**, **greptile** | External AI review bots |

### D. Maintaining the setup itself

| Plugin | Use |
|---|---|
| **claude-md-management** | Audits CLAUDE.md quality and records lessons from a session into it |
| **claude-code-setup** | Analyzes a codebase and **recommends hooks, skills, MCP servers and subagents**. A good first step on an existing app. |
| **hookify** | Creates hooks from conversation patterns ("stop doing X") or explicit rules |
| **skill-creator** | Creates, evaluates and benchmarks skills |
| **plugin-dev** | Packages your setup as a plugin |
| **session-report** | HTML report of token usage, cache efficiency, subagents and the most expensive prompts |
| **commit-commands** | Commands for commit, push and PR |

### E. Loops and automation

| Command / plugin | Use |
|---|---|
| `/goal <condition>` (built in) | An evaluator re-checks the condition every turn, and the agent keeps working until it holds |
| `/loop [interval] <prompt>` (built in) | Repeats a prompt or command on an interval |
| `/schedule` (built in) | Cloud routines on a cron schedule |
| `/batch <instruction>` (built in) | Splits a large change across 5 to 30 subagents, each in its own worktree |
| **ralph-loop** | Runs the same prompt in a loop until the task completes (the "Ralph Wiggum" technique). Only use it with a strong Stop gate. |

### F. UI and design

| Plugin / skill | Use |
|---|---|
| **frontend-design** | Production-grade UI with real design quality |
| **playground** | Single-file HTML explorers with live controls |
| **agent-browser** skills (`agent-browser skills get core`, `... electron`) | How to drive the browser or Electron renderer correctly, matched to your installed CLI version |
| **figma** | Read designs and tokens |

### G. Output styles for learning

`explanatory-output-style` and `learning-output-style` make the agent explain *why*
as it works. Use them while you build the skills in [`08`](./08-human-skills.md).

---

## 4. Skills to write yourself (project-specific)

Generic skills don't know your project. These six pay off in almost every app:

| Skill | Invoked by | What it does |
|---|---|---|
| `/verify-change` | Agent or you | Runs the verification ladder (07) in order and reports evidence in a fixed format. Template provided. |
| `/close-increment` | You | Updates the phase status doc and `current-state.md`, adds lessons to CLAUDE.md, writes the completion report, and makes an atomic commit. Template provided. |
| `/spec <feature>` | You | Interviews you with AskUserQuestion, then writes `docs/specs/<date>-<feature>.md` from the template |
| `/triage-bug <issue>` | You | Gets the issue, reproduces it with a failing test, finds the root cause with the graph, proposes a fix, and stops for approval |
| `/weekly-cleanup` | You or `/schedule` | Runs knip and jscpd, checks docs against code, lists the top 5 cleanup PRs. Read-only until you approve. |
| `/release` | You | Release checklist: version bump, changelog, packaged-build smoke test, tag. `disable-model-invocation: true`. |

---

## 5. Subagents worth defining

`.claude/agents/<name>.md`:

```markdown
---
name: code-reviewer
description: Reviews the current diff in a fresh context against the spec. Use after implementation, before commit.
tools: Read, Grep, Glob, Bash
model: opus
---
You review diffs. You did not write this code. Inputs: the diff (git diff), the spec path.
Report ONLY: (1) requirement not met, (2) correctness bug, (3) missing test for a listed
edge case, (4) change outside the task's scope. Each finding: file:line, why, suggested fix.
Do not report style preferences. If nothing qualifies, say "No blocking findings".
```

Other useful agents:

- `security-reviewer`: injection, auth, secrets, unsafe deserialization, SSRF.
- `test-writer`: writes tests from the spec *before* implementation, for the
  "one writes tests, another passes them" pattern.
- `docs-drift-checker`: compares `docs/` claims with the code and lists stale
  statements with file references.
- `explorer`: read-only investigation that returns a short report with `file:line`
  references.

Keep each agent to one area of expertise and the smallest tool set it needs.

---

## 6. Hooks worth setting

The full template is [`templates/claude-settings.example.json`](./templates/claude-settings.example.json),
with Node scripts in [`templates/hooks/`](./templates/hooks/).

| Event | Hook | Why |
|---|---|---|
| `SessionStart` | `code-review-graph status`; print the first lines of `docs/current-state.md` | The agent starts each session knowing the state of the project |
| `PostToolUse` (`Edit\|Write`) | Format the edited file; run `code-review-graph update --skip-flows` | Clean diffs and a fresh graph without any prompting |
| `PreToolUse` (`Edit\|Write`) | Block edits to `.env*`, lockfiles, generated code and migrations | Hard guardrails |
| `PreToolUse` (`Bash`) | Block destructive commands (`rm -rf`, `git push --force` to main, `DROP TABLE`) | Hard guardrails |
| `Stop` | Run the fast `verify` when files changed; block (exit 2) on failure | The agent can't finish while checks are red |
| `PreCompact` | Remind the agent to keep the list of modified files and the test commands | Important context survives compaction |
| `InstructionsLoaded` | Log which CLAUDE.md and rules files loaded | Debugging rules that don't seem to apply |

The block mechanics, per the official docs:

- **PreToolUse:** exit code `2` blocks the tool call, and stderr becomes the reason
  shown to the agent. Alternatively, print JSON with
  `hookSpecificOutput.permissionDecision: "deny"`.
- **Stop:** exit code `2`, or JSON `{"decision":"block","reason":"…"}`, keeps the
  agent working.

Hooks run with your user permissions, so review every hook script as you would any
code.

---

## 7. Where Natively stands (for reference)

| Item | Natively today |
|---|---|
| CLAUDE.md | 686 lines (cross-platform contract, graph-first workflow, CDP testing rules, completion report) |
| MCP | `code-review-graph` (project `.mcp.json`) |
| Hooks | **User-level** (`~/.claude/settings.json`): graph `update` after Edit/Write, graph `status` at SessionStart. Not in the repo. |
| Skills | code-review-graph's 4 skills (user-level); superpowers was used for specs and plans; the levnikolaevich suite is vendored as docs |
| Pre-commit | husky: native-arch gate (fails closed), react-doctor (advisory), graph `detect-changes` (advisory) |

What to improve is in [`09-natively-reference.md`](./09-natively-reference.md#gaps-worth-fixing).
