# Phase 4: Install the harness (H1–H15)

Goal: guides (CLAUDE.md, docs, skills) and sensors (verify, hooks, lint, CI, review
agent) installed in TARGET, adapted to its stack, each one tested. Background reading:
`PLAYBOOK/03-skills-plugins-and-extensions.md`, `PLAYBOOK/07-verification-and-quality-gates.md`,
`PLAYBOOK/templates/README.md`.

**Mode differences**

| | NEW | EXISTING | EXISTING+AI (merge) |
|---|---|---|---|
| Existing files | none | create alongside, never replace | **propose diffs only** |
| Lint / architecture rules | error | **warn** first | warn; keep existing configs |
| Stop hook (H3) | on | on only if `verify` is green at baseline | same as EXISTING |
| Reformat whole repo | fine (it's empty) | **never** without a separate, explicit OK | never |

Depth = "Minimal" → do only H1–H7 and H12–H13. Depth = "no CI / no pre-commit" → skip
H10/H11.

Record every step in the setup log. Each step ends with a ✅ check whose output you show.

---

## H1. `.gitignore` and tracking of `.claude/`

1. Ensure these are ignored (append only what's missing; show the diff):
   ```gitignore
   # AI harness: local state
   .code-review-graph/
   agent-browser.json
   .agent/
   .claude/settings.local.json
   CLAUDE.local.md
   .scaffold-tmp/
   ```
2. Ensure the shared harness is **not** ignored. Check:
   `git check-ignore -v .claude/settings.json .claude/hooks/x.mjs .claude/skills/x/SKILL.md`.
   If a broad rule such as `.claude` matches, **ask**, then replace it with:
   ```gitignore
   .claude/*
   !.claude/settings.json
   !.claude/hooks/
   !.claude/skills/
   !.claude/agents/
   !.claude/rules/
   ```
   (Natively's `.gitignore` has exactly this problem; see playbook 09, gap 4.)

✅ `git check-ignore` reports the local files as ignored and the shared ones as not ignored.

## H2. The verify command

It's already done in E8 (EXISTING) or planned in N7 (NEW). For NEW projects, create it now:

- **JS/TS** `package.json` scripts (all work under cmd.exe and sh):
  `typecheck`, `lint`, `test`, `verify` (`npm run typecheck && npm run lint && npm run test`),
  `verify:full` (`npm run verify && npm run test:e2e && npm run build`, keeping only the
  parts that exist).
- **Other stacks:** `scripts/verify.mjs` from `PLAYBOOK/templates/scripts/verify.mjs`,
  with STEPS from `reference/stack-matrix.md`.

NEW projects also need **one real unit test** and, for UI apps, **one e2e smoke test**
(the roadmap's Phase 0 exit criterion). If the scaffold has no test runner, adding one
(Vitest, pytest, …) is a dev-dependency install, so use the Stop protocol.

Set `VERIFY_CMD` (e.g. `npm run verify` or `node scripts/verify.mjs`).
✅ Run it and show the summary (green for NEW; matches the baseline for EXISTING).

## H3. Hook scripts → `TARGET/.claude/hooks/`

Copy from `PLAYBOOK/templates/hooks/`, then adapt:

| Script | Adapt |
|---|---|
| `protect-files.mjs` | `PROTECTED` list: keep `.env*`, lockfiles, `.git`. Add the project's generated folders. **Ask** before protecting `migrations/`, since some projects need the agent to write migrations. |
| `guard-shell.mjs` | Add project-specific dangers (production hosts, deploy commands) if the user names any |
| `format-on-edit.mjs` | JS: works as is if Prettier is the formatter. Other stacks: change the formatter call (`ruff format`, `gofmt -w`, `cargo fmt`, `dotnet format`) or drop the hook. |
| `graph-update.mjs` | **Skip** if `~/.claude/settings.json` already runs `code-review-graph update` (no duplicates), or if code-review-graph was skipped |
| `stop-verify.mjs` | Nothing; it reads `CLAUDE_VERIFY_CMD` from settings `env`. EXISTING with a red baseline: copy it, but **don't register it in H4 yet**. Note "deferred" in the log. |

Test each script **before** registering it. Bash:
```bash
echo '{"tool_name":"Edit","tool_input":{"file_path":".env"}}' | node .claude/hooks/protect-files.mjs; echo "exit=$?"     # expect 2
echo '{"tool_name":"Edit","tool_input":{"file_path":"src/a.ts"}}' | node .claude/hooks/protect-files.mjs; echo "exit=$?" # expect 0
echo '{"tool_name":"Bash","tool_input":{"command":"git reset --hard"}}' | node .claude/hooks/guard-shell.mjs; echo "exit=$?" # expect 2
echo '{"tool_name":"Bash","tool_input":{"command":"npm test"}}' | node .claude/hooks/guard-shell.mjs; echo "exit=$?"         # expect 0
```
PowerShell: `'{"tool_name":"Edit","tool_input":{"file_path":".env"}}' | node .claude/hooks/protect-files.mjs; $LASTEXITCODE`

`stop-verify.mjs` test, run from TARGET:
`echo '{"session_id":"setup-test"}' | node .claude/hooks/stop-verify.mjs; echo "exit=$?"`.
That gives 0 if the tree is clean or verify passes, and 2 with a failure tail if verify
fails. For a non-npm `VERIFY_CMD`, set `CLAUDE_VERIFY_CMD` for this one command
(Bash: `CLAUDE_VERIFY_CMD="node scripts/verify.mjs" node …`; PowerShell:
`$env:CLAUDE_VERIFY_CMD='node scripts/verify.mjs'; …`). Afterward, delete
`<os.tmpdir()>/claude-stop-verify-setup-test.json`.

✅ Every script's results match the expected exit codes.

## H4. `.claude/settings.json` (create or merge)

- **No file yet:** start from `PLAYBOOK/templates/claude-settings.example.json`.
  Remove the entries for skipped hooks. Set `env.CLAUDE_VERIFY_CMD` to `VERIFY_CMD`.
  Adjust `permissions.allow` to the real commands. **Leave out the Stop hook** if it
  was deferred in H3.
- **File exists (merge mode):** load it, add the missing `hooks` entries and `allow`
  rules, never remove or reorder existing ones, show the diff, and ask.
- Keep the **exec form** (`"command": "node", "args": ["${CLAUDE_PROJECT_DIR}/.claude/hooks/…"]`).
  It runs without a shell on Windows and macOS alike.

✅ The JSON parses:
`node -e "JSON.parse(require('fs').readFileSync('.claude/settings.json','utf8'));console.log('ok')"`.
Hooks take effect after the restart in V1.

## H5. Skills and the reviewer agent

Copy:
- `PLAYBOOK/templates/skills/verify-change/SKILL.md` → `.claude/skills/verify-change/SKILL.md`
- `PLAYBOOK/templates/skills/close-increment/SKILL.md` → `.claude/skills/close-increment/SKILL.md`
- `PLAYBOOK/templates/agents/code-reviewer.md` → `.claude/agents/code-reviewer.md`

Adapt: in `verify-change`, replace the commands table with the project's **real**
commands. Remove the agent-UI rung if H14 is skipped. In `close-increment`, point to the
real status-doc path. If a name already exists in `.claude/skills` or `.claude/agents`,
**ask** before choosing a different name. Never overwrite.

✅ The files exist and their frontmatter has `name` and `description`.

## H6. CLAUDE.md

- **No CLAUDE.md (and no AGENTS.md):** write one from `PLAYBOOK/templates/CLAUDE.template.md`
  using **only verified facts**: real commands, the docs map, the workflow lines, and the
  conventions you found. Delete placeholder lines you can't fill. **150 lines or fewer.**
- **AGENTS.md only:** ask. Recommended: keep AGENTS.md as the source and create a
  CLAUDE.md whose first line is `@AGENTS.md`, followed by the Claude-specific additions.
- **CLAUDE.md exists:** propose a diff that adds only the missing pieces (the docs map,
  `VERIFY_CMD`, the workflow lines, the graph line). If the file is over about 200
  lines, propose (don't do) moving sections into `.claude/rules/` with `paths:`.

✅ `wc -l CLAUDE.md` (or `(Get-Content CLAUDE.md).Count`) is within budget, and every
command in it was run successfully in this setup.

## H7. Docs skeleton

NEW: `docs/README.md` (index), `docs/current-state.md`, `docs/status/phase-0.md` (this
setup *is* Phase 0), plus the docs from phase 3a. EXISTING: already done in E5. Add
`docs/status/ai-harness.md`, or use the setup log as the record.
Link the docs index from CLAUDE.md.

✅ Every link in CLAUDE.md and `docs/README.md` resolves.

## H8. Formatter and linter

Check what exists first; reuse the existing tools. If something is missing, use the
Stop protocol for the dev dependencies:

- JS/TS: `prettier`, `eslint`, `@eslint/js`, `typescript-eslint` (or `@biomejs/biome`
  for both). Use a flat `eslint.config.js`. Add the `lint` / `format` scripts.
- Others: see `reference/stack-matrix.md`.
- **Never** run a whole-repo `--write` or format on an EXISTING project without a
  separate, explicit OK. That commit should be on its own. The format hook formats only
  the files the agent edits.

✅ `lint` runs (errors are allowed in EXISTING; record the count).

## H9. Architecture rules

From the layers table in `docs/01-architecture.md` (NEW) or the "target" layers in
`docs/architecture-overview.md` (EXISTING), write import-boundary rules **whose
messages say how to fix the violation**:

- JS/TS: `dependency-cruiser` (dev dependency, so use the Stop protocol). Write
  `.dependency-cruiser.cjs` by hand from the example in playbook 07 §2 instead of the
  interactive `--init`. Add it to `lint`. Set `severity: 'warn'` for EXISTING.
- Others: `import-linter` (Python), `depguard` via golangci-lint (Go), ArchUnit (JVM),
  NetArchTest (.NET).

✅ A deliberately wrong import (in a scratch file you delete afterward) produces a
message that names the fix.

## H10. Pre-commit (fast gates only)

Reuse whatever exists (husky / lefthook / pre-commit) and **merge** into it. Otherwise:

- JS: Stop protocol for `husky` → `npx husky init` → edit `.husky/pre-commit` to run
  fast checks only (e.g. `npm run lint`), plus an advisory
  `code-review-graph detect-changes --brief || true` if it's installed.
- Python: Stop protocol for `pre-commit` → `.pre-commit-config.yaml` with ruff hooks →
  `pre-commit install`.
- Optional secret scanner (gitleaks): Stop protocol; the user installs it.

✅ A trial commit isn't needed. Run the hook script directly and show that it passes.

## H11. CI

Hosting = GitHub: copy `PLAYBOOK/templates/ci/verify.yml` → `.github/workflows/verify.yml`.

- Replace `DEFAULT_BRANCH` with the real default branch (E1 / `inspect-project.mjs`).
- `matrix.os`: the P5 platforms (Windows desktop → `windows-latest`, macOS →
  `macos-latest`, web-only → `ubuntu-latest`).
- If workflows already exist: **propose** adding a job, or check that the existing
  workflow runs `verify:full`. Also check their branch filters against the default
  branch, and report mismatches.
- **Don't push.** Note in the handoff: after the user pushes, run
  `gh run list --workflow verify.yml` and confirm the run has more than 0 jobs.

GitLab / other: write the equivalent pipeline only if the user asks.

✅ YAML is valid: `node -e` with a YAML parser if one is available, otherwise review
by eye that `on:` and `jobs:` are top-level keys.

## H12. Plugins (each is an install, so use the Stop protocol; ask once for the list)

Check availability: `claude plugin list --available --json` (look under `available`).
Propose this set and let the user tick items (AskUserQuestion, multiSelect, header
`Plugins`):

| Plugin | Why |
|---|---|
| `context7@claude-plugins-official` | Current library docs; fewer invented APIs |
| `<lang>-lsp@claude-plugins-official` (see stack matrix) | Symbol navigation and type errors after edits. Its README may require a language-server binary, which is another Stop protocol. |
| the P5 methodology (`superpowers` or `feature-dev`), if chosen | Process |
| `security-guidance@claude-plugins-official` | Security warnings on edits and on Stop |

For each approved plugin, run (or have the user run) `claude plugin install <id> --scope project`.
Project scope makes it shared through `.claude/settings.json`. Then:
`claude plugin list --json` shows it, and `claude plugin details <id>` shows its token
cost. Report that cost. The plugins load after `/reload-plugins` or the V1 restart.

✅ Every chosen plugin shows `enabled: true` in `claude plugin list --json`.

## H13. code-review-graph MCP integration (it changes config, so ask)

1. Preview: `code-review-graph install --platform claude-code --repo "<TARGET>" --dry-run`.
   Show the output. It writes the project `.mcp.json`, injects graph instructions into
   CLAUDE.md, and adds `.code-review-graph/` to `.gitignore`.
2. Recommend `--no-instructions` (CLAUDE.md from H6 already has the graph line and stays
   short) and `--no-hooks` (H3's `graph-update.mjs` is the cross-platform project hook).
   Ask (header `Graph`).
3. Run the approved command. In merge mode, show the `.mcp.json` diff first.
4. `code-review-graph build --repo "<TARGET>"` if E3 didn't run it.
5. Optional semantic search: the embeddings extra is a large download (Stop protocol):
   `uv tool install --force "code-review-graph[embeddings]"`, then
   `code-review-graph embed --repo "<TARGET>"`.

✅ `claude mcp list` lists `code-review-graph`. It shows "Pending approval" until the user
approves it after the restart. `code-review-graph status` shows the node count.

## H14. Agent UI access (UI projects; needs T5)

- **Web:** give the dev server a fixed port (e.g. `--port 5173 --strictPort`). Smoke test,
  with the dev server started only with the user's OK:
  `agent-browser skills get core` → `agent-browser open http://localhost:5173` →
  `agent-browser snapshot -i` → `agent-browser screenshot docs/evidence/setup-smoke.png` →
  `agent-browser close`.
- **Electron / web-renderer desktop:** the agent needs a CDP port in **dev builds only**.
  Propose a follow-up task, done through the normal feature loop and not in this setup,
  to add a Node `dev:agent` launcher modelled on Natively's `scripts/dev-agent.mjs`:
  free port, never 9222/9229, gitignored `agent-browser.json`, worktree-local data
  dir, never enabled when packaged. Record it as a next step.
- **Native UI (SwiftUI, WinUI, mobile):** no agent UI access. Record "UI verification
  is manual" in CLAUDE.md and in `verify-change`.

✅ A snapshot shows the app's interactive elements, or the limitation is recorded.

## H15. Optional extras (offer as one multiSelect; do only what's chosen)

- `/install-github-app` for `@claude` in PRs and automatic review. The **user** runs it,
  since it's interactive and needs repo admin.
- A `/weekly-cleanup` routine (`/schedule`) running knip and jscpd plus a docs-drift check.
- The `claude-md-management` plugin, for monthly CLAUDE.md pruning.
- Splitting CLAUDE.md into `.claude/rules/` with `paths:` (if it's long).
- A `.github/pull_request_template.md` section for the completion report (playbook 07 §9).

✅ **Phase 4 checkpoint:** list H1–H15 as DONE / SKIPPED / DEFERRED, with evidence. Then
open `06-validate-and-handoff.md`.
