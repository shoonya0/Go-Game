# Setup agent: humans start here

This folder holds **instructions for an AI coding agent** (written for Claude Code)
to set up the whole AI-development harness from this playbook in a project. It covers
two cases:

- **A new project:** an empty folder or an idea.
- **An existing codebase where no AI has been used yet.** It also handles codebases
  with some AI setup already, by merging instead of overwriting.

The agent does the reading, writing and checking. **Whenever something has to be
installed or authenticated, or an existing file has to change, it stops and asks
you.** That covers code-review-graph, gh, plugins, npm packages and similar. You
decide whether to run it yourself, let the agent run it, or skip it.

---

## How to start

1. Open a terminal **in the target project folder**. For a new project, open it in the
   parent folder, or create the empty folder first.
2. Start Claude Code: `claude`
3. Paste this prompt, replacing `<PLAYBOOK>` with the absolute path of the
   `stepsForCreatingApp` folder. On this machine that's
   `D:\projects\native\natively-cluely-ai-assistant\stepsForCreatingApp`.

```text
Read <PLAYBOOK>/agent/00-START.md and follow it exactly to set up the AI development
harness for this project. The playbook root is <PLAYBOOK>. Ask me before installing,
authenticating, overwriting, committing or pushing anything.
```

If Claude Code can't read the playbook folder because it's outside the project, run
`/add-dir <PLAYBOOK>` first.

**Optional:** install the launcher skill so that `/setup-ai-harness` works in any
project. See [`skill/setup-ai-harness/SKILL.md`](./skill/setup-ai-harness/SKILL.md).

---

## What it will do

| Phase | File | Summary |
|---|---|---|
| 0 | `00-START.md` | Rules, the stop-and-ask protocol, the progress log, the phase map |
| 1 | `01-preflight.md` | Finds the target, OS and git state; chooses **NEW**, **EXISTING** or **EXISTING+AI** |
| 2 | `02-tools.md` | Checks every tool; **stops for each missing one** with the exact install command for your OS |
| 3a | `03-new-project.md` | Interviews you, writes vision/stack/architecture/roadmap docs, scaffolds the app |
| 3b | `04-existing-project.md` | Safety branch, secrets check, code graph, architecture tour, current-state doc, test baseline |
| 4 | `05-harness.md` | CLAUDE.md, hooks, skills, reviewer agent, verify command, lint/architecture rules, pre-commit, CI, plugins, MCP, UI access |
| 5 | `06-validate-and-handoff.md` | Restart, live checks, context budget, setup log, completion report, optional commit |

Everything it does is recorded in **`docs/ai-harness-setup.md`** in the target
project. If the session ends or the context is cleared, start again with the same
prompt and it resumes from the first unfinished step.

## What it will never do without asking

- Install, upgrade or uninstall anything: system tools, global npm/pip packages,
  project dependencies, plugins, MCP servers, browsers
- Run anything that needs your login (`gh auth login`, `claude` login, OAuth)
- Overwrite or delete an existing file, or change `.gitignore` / `.gitattributes`
- Run your project's test suites for the first time (they might touch real services)
- Commit, push, open PRs, or rewrite git history

## Folder contents

```text
agent/
  README.md                      ← this file (for humans)
  00-START.md … 06-…md           ← the agent's step-by-step instructions
  reference/stack-matrix.md      ← per-language verify/lint/format/architecture tools
  reference/install-commands.md  ← install + verify commands per tool, per OS
  templates/setup-log.md         ← progress log copied into the target project
  scripts/check-tools.mjs        ← read-only tool detection (node check-tools.mjs [--json])
  scripts/inspect-project.mjs    ← read-only project inspection (node inspect-project.mjs <dir> [--json])
  skill/setup-ai-harness/SKILL.md← optional /setup-ai-harness launcher
```

The agent copies files from `../templates/` (hooks, skills, the reviewer agent and doc
skeletons) and adapts them.

## Requirements before you start

- Claude Code installed and logged in (`claude --version`).
- **Node.js 20+.** The hooks and helper scripts are Node. If it's missing, the agent
  can only run in a reduced mode and asks you to install Node first.
- Everything else is checked, and installed only with your OK, during phase 2.

> Note for this machine (checked 2026-10-06): `node`, `npm`, `git`, `python`,
> `code-review-graph` and `agent-browser` are present. **`gh` (GitHub CLI) is
> missing**, and `uv` is missing (optional).
