# Templates

Files to copy into a project. They are templates, so adapt names, paths and commands.

| Template | Copy to | Purpose |
|---|---|---|
| [`CLAUDE.template.md`](./CLAUDE.template.md) | `./CLAUDE.md` | A short map-style rules file of about 100 to 150 lines. It's named `.template.md` here so Claude Code doesn't auto-load it as a nested CLAUDE.md. |
| [`claude-settings.example.json`](./claude-settings.example.json) | `./.claude/settings.json` | Project hooks: protect files, guard shell commands, format on edit, graph update, Stop → verify |
| [`hooks/protect-files.mjs`](./hooks/protect-files.mjs) | `./.claude/hooks/` | `PreToolUse` (Edit/Write): blocks `.env*`, lockfiles, generated output, migrations |
| [`hooks/guard-shell.mjs`](./hooks/guard-shell.mjs) | `./.claude/hooks/` | `PreToolUse` (Bash/PowerShell): blocks a short list of destructive commands |
| [`hooks/format-on-edit.mjs`](./hooks/format-on-edit.mjs) | `./.claude/hooks/` | `PostToolUse`: runs the project's Prettier on the edited file, if Prettier is installed |
| [`hooks/graph-update.mjs`](./hooks/graph-update.mjs) | `./.claude/hooks/` | `PostToolUse`: `code-review-graph update`, if the tool is installed |
| [`hooks/stop-verify.mjs`](./hooks/stop-verify.mjs) | `./.claude/hooks/` | `Stop`: runs `npm run verify` when the working tree changed; blocks the stop on failure, up to 3 times |
| [`skills/verify-change/SKILL.md`](./skills/verify-change/SKILL.md) | `./.claude/skills/verify-change/` | `/verify-change`: the verification ladder plus an evidence report |
| [`skills/close-increment/SKILL.md`](./skills/close-increment/SKILL.md) | `./.claude/skills/close-increment/` | `/close-increment`: status docs, lessons, completion report, atomic commit |
| [`agents/code-reviewer.md`](./agents/code-reviewer.md) | `./.claude/agents/` | Fresh-context reviewer subagent |
| [`scripts/verify.mjs`](./scripts/verify.mjs) | `./scripts/verify.mjs` | One cross-platform `verify` entry point for non-npm stacks (Python, Go, Rust …); `--full` adds the slow checks |
| [`ci/verify.yml`](./ci/verify.yml) | `./.github/workflows/verify.yml` | CI on push to the default branch and on PRs, with an OS matrix; set `DEFAULT_BRANCH` and the matrix first |
| [`docs/*.md`](./docs/) | `./docs/…` | Vision, architecture, ADR, roadmap, phase status, current state, feature spec, test plan |

## Why the hooks are Node scripts

Hook commands in shell form (no `args`) run through `sh -c` on macOS and Linux, and
through Git Bash (or PowerShell, if Git Bash is missing) on Windows. The same one-liner
can behave differently on each.

The settings template therefore uses the **exec form**, `"command": "node"` plus
`"args": [...]`. Claude Code spawns it directly **with no shell on any OS**, and
substitutes `${CLAUDE_PROJECT_DIR}` as a plain string. The scripts use only Node
built-ins and `node:path`, so they behave the same on Windows and macOS.

## Install

```text
mkdir .claude\hooks .claude\skills .claude\agents         (Windows)
mkdir -p .claude/hooks .claude/skills .claude/agents      (macOS)
```

Copy the files, rename `claude-settings.example.json` to `.claude/settings.json`, and
adjust:

- `protect-files.mjs` → the `PROTECTED` list (e.g. remove `migrations` if the agent
  should write migrations)
- `stop-verify.mjs` → set `CLAUDE_VERIFY_CMD` in settings `env` if your fast check isn't
  `npm run verify`
- run `/hooks` in Claude Code to confirm that they're registered

Test a hook by hand by piping sample input into it:

```bash
echo '{"tool_name":"Edit","tool_input":{"file_path":".env"}}' | node .claude/hooks/protect-files.mjs; echo "exit=$?"
```

Expected output: a "Blocked: …" message and `exit=2`.
