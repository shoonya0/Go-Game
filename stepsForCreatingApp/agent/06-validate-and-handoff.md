# Phase 5: Validate and hand off (V1–V7)

Goal: prove the harness works **in a fresh Claude Code session**, measure its context
cost, write the record, and hand over with honest labels for what was and wasn't
verified.

---

## V1. Restart (the user does this)

Hooks, project settings, plugins and MCP servers load when a session starts. Stop
protocol, header `Restart`:

```text
⏸ ACTION NEEDED: restart Claude Code so the new harness loads (step V1)
1. Exit this session (Ctrl+C twice, or /exit).
2. In <TARGET>, run:  claude --continue
3. If asked, approve the project MCP server "code-review-graph" and review the hooks.
4. Type: continue setup
```

After the restart, re-read `docs/ai-harness-setup.md` (00-START §0.4) and continue at V2.

## V2. Registration checks (the user types the slash commands; you read the results)

Ask the user to run and paste, or just confirm:

| Command | Expect |
|---|---|
| `/hooks` | PreToolUse (protect, guard), PostToolUse (format, graph), Stop (verify, unless deferred) |
| `/mcp` | `code-review-graph` connected |
| `/context` | CLAUDE.md token count; skills and plugin descriptions; MCP tools (deferred) |

You can check some of this yourself: `claude mcp list`, `claude plugin list --json`.

## V3. Live probes (you run these; each is harmless)

1. **Protect hook:** try to Write `.env.harness-probe` with the content `probe`. Expect it
   to be **blocked**. If the write succeeds, the hook isn't active: delete the file,
   report the failure, and go back to H4.
2. **Guard hook:** run `git reset --hard --dry-run-probe`. Expect it to be blocked
   before execution. (If the hook isn't active, git rejects the unknown flag, so nothing
   happens either way.)
3. **Format hook (if installed):** edit a tracked source file by adding a harmless blank
   line, confirm it's formatted, then revert the edit exactly.
4. **Graph MCP:** call `list_graph_stats` (or `semantic_search_nodes` for a known
   function) and show the node count.
5. **Skill:** invoke `/verify-change` and show the evidence table it produces.
6. **Stop hook (if registered):** it ran at the end of your last turn if files had
   changed. Check that no "Verification failed" block is pending. If verify is red, fix
   it or record why it's deferred.
7. **UI (if H14 was done):** `agent-browser snapshot -i` against the running dev server.

## V4. Context budget

From the `/context` output (V2), report:

- CLAUDE.md tokens. Above about 3k tokens (roughly 200 lines), propose moving content
  into rules or skills.
- MCP servers whose deferred tools are large and **unused by this project**. For
  example, a huge API-client MCP server that the project never needs: suggest disabling
  it for this project (`/mcp`) or removing it.
- The number of model-invocable skills. Suggest `disable-model-invocation: true` for
  side-effect skills.

## V5. Final setup log

Complete `docs/ai-harness-setup.md`:

- every step's status (DONE / SKIPPED / DEFERRED / BLOCKED), with evidence;
- tool versions; the plugins with their token costs;
- the decisions the user made (mode, platforms, depth, methodology, protect list);
- **deferred items, with the condition that re-enables each** (e.g. "Stop hook: register
  once `npm run verify` is green; baseline failures: …").

## V6. Completion report (show it in chat; also append it to the log)

```markdown
### Change summary
- Files added / changed: …
- Harness pieces installed: …

### Validation
- Covered by automated checks: hook stdin tests (H3), verify run (H2), JSON/YAML parse (H4/H11)
- Verified live in Claude Code on <OS>: V3 probes 1–7 (list results)
- Not verified on <other OS>: hooks and scripts were written cross-platform but not executed there
- Requires user action: <push + check CI jobs > 0>, <`gh auth login`>, …

### Commands executed
- …

### Remaining risks / next steps
- …
```

Use these labels honestly. "Works on Windows and macOS" is allowed only if both were
actually run.

## V7. Commit (only with approval) and first task

1. Ask (header `Commit`): "Commit the harness as one commit (Recommended)" /
   "Commit in groups (docs / harness / CI)" / "Don't commit; I'll review first".
2. If committing: stage **only** the files this setup created or changed. List them with
   `git status --porcelain` and leave out any unrelated user changes. Suggested message:
   `chore: add AI development harness (CLAUDE.md, hooks, skills, verify, CI)`, plus the
   attribution trailer the session requires. **Never push** unless asked. If asked to
   push, afterward run `gh run list --workflow verify.yml` and confirm there are jobs.
3. Suggest the **first real task**: small, low-risk, done with the playbook loop
   (`PLAYBOOK/06-feature-loop.md`). For example, a bug with a clear reproduction, or the
   next roadmap increment for NEW projects. Offer to start it in a **fresh session**
   (`/clear`), so the setup conversation doesn't fill its context.

Setup complete.
