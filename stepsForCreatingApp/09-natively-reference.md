# 09: Reference case, Natively (this repository)

Natively is a production Electron app for macOS and Windows, built mostly with AI
agents. This file maps each practice in this playbook to how Natively does it, with
paths, so you can copy what works. It also lists the gaps found while researching.

Inspected on 2026-10-06 at commit `fcbfdd2` (`master`). This is a read-only review;
nothing in the repo was changed.

---

## What Natively does, mapped to the playbook

| Practice (playbook §) | How Natively does it | Where |
|---|---|---|
| Rules file (01 §3, 03) | `CLAUDE.md`, 686 lines: graph-first workflow, a **Cross-Platform Development Contract**, CDP UI testing rules, stop conditions, and a **required completion report** with fixed validation categories | `CLAUDE.md` |
| Code graph (02 Tier 1) | `code-review-graph` as a project MCP server; graph `update` after every Edit/Write and `status` at session start (user-level hooks); `detect-changes --brief` in pre-commit. About 37.7k nodes and 318k edges, including `TESTED_BY` edges. | `.mcp.json`, `~/.claude/settings.json`, `.husky/pre-commit` |
| Product docs in the repo (01 §4) | Read-in-order docs: vision & scope → architecture → security/stealth model → engines → **phased roadmap** with 🎯 goal / 📦 deliverables / ✅ exit criterion / ⚠️ risk per phase, plus a **global definition of done**. A second track (career copilot, docs 05–11) follows the same shape. | `_docs/README.md`, `_docs/00–11` |
| Honest ceiling (04 stage 1) | Explicit list of what can't be done (Task Manager hiding, defeating proctoring software) | `_docs/README.md` |
| Status tracking (04 stage 6) | One status doc per phase, split into **Increments** with problem, design, tests and commit hashes; a one-page `current-state.md` snapshot tied to a commit SHA | `_docs/phase-*-status.md`, `_docs/current-state.md` |
| Test plans (07 §5) | Test levels table, a fixtures layout, and the **golden rule for LLM outputs** (assert invariants, not text) | `_docs/test/README.md` |
| Specs and plans (06 §4–5) | superpowers-style design specs and plans with **Global Constraints** and **Review Focus** (edge cases that no happy path exercises), each assigned to a task | `docs/superpowers/specs/`, `docs/superpowers/plans/` |
| Investigation handoff (06 variants) | `00-START-HERE` → root causes ordered by pipeline stage with `file:line` → fix plan → experiments; "no production code changed" | `docs/retrieval-handoff/` |
| Agent UI testing (07 §6) | `npm run dev:agent`: Node launcher, free CDP port (never 9222/9229), gitignored `agent-browser.json`, worktree-local `userData`, single-instance lock skipped only when not packaged; agent-browser for snapshots | `scripts/dev-agent.mjs`, `scripts/devAgentSupport.mjs`, `CLAUDE.md` |
| Fast local gates (07 §3) | husky pre-commit: native-architecture gate (**fails closed**), react-doctor (advisory), graph detect-changes (advisory) | `.husky/pre-commit` |
| CI (07 §3) | `build-smoke.yml`: macOS + Windows matrix with `fail-fast: false`, on PRs, pushes and a weekly schedule; a fork-safe "core only" path when the private submodule can't be fetched. `react-doctor.yml` on PRs. | `.github/workflows/` |
| Evals for AI features (07 §5) | `test:golden`, `test:privacy`, `test:auto-answer:judge`, many `benchmark:*` scripts, tiered `ci:tier1..4` | `package.json` |
| Tests at scale | About 1,339 `*.test.mjs` files under `electron/` and `src/`; script tests under `scripts/__tests__/` | |
| Memory and lessons (01 §1, §8) | Auto-memory "feedback" notes: graph and docs first, compare against a clean-HEAD baseline, the CRLF trap, the Electron-ABI test runner | `~/.claude/projects/<repo>/memory/` |
| Skills library (03) | The levnikolaevich `claude-code-skills` collection is vendored for reference | `_docs/claude-code-skills/` |

---

## Lessons recorded in the repo's own history

These came from real failures, which is why they're worth more than generic advice:

1. **Isolate agent instances.** A second agent attached to a contested debug port and
   "spent four rounds reading somebody else's renderer"; a plain `electron .` was
   refused by the single-instance lock. Fix: the `dev:agent` launcher. (Source:
   `scripts/dev-agent.mjs` header.)
2. **Make sure CI runs at all.** A `schedule:` key placed on the job instead of the
   workflow made GitHub reject the file, so every run completed with **zero jobs**.
   (Source: `build-smoke.yml` comments.)
3. **Run CI on pushes, not only PRs.** A commit broke the first typecheck step and
   "went unnoticed for three days while every subsequent PR inherited a red check it
   did not cause." (Source: `build-smoke.yml`.)
4. **Test every target OS in CI.** Windows-breaking package scripts (POSIX subshells,
   `VAR=value cmd`, single-quoted globs under cmd.exe) "shipped green on the
   macOS-only matrix." (Source: `build-smoke.yml`.)
5. **Read the docs before coding, then check them against the code.** Skipping
   `docs/embedding-architecture.md` let an unbounded-batch bug through, and the same
   doc had a stale number. (Source: memory note.)
6. **Compare against a rebuilt baseline**, not a remembered "N pre-existing failures".
   (Source: memory note.)
7. **Watch for platform text traps.** Python text-mode writes on Windows produce CRLF
   in an `eol=lf` repo; heredocs with quotes failed in the Bash tool, so use the
   editor tools instead. (Source: memory note.)

---

## Gaps worth fixing

These are observations from this research, not changes that were made. Each one is
verified against the repo as of `fcbfdd2`.

| # | Gap | Evidence | Suggested fix |
|---|---|---|---|
| 1 | **CI branch filters don't match the default branch.** `origin/HEAD → master`, but `build-smoke.yml` triggers on push to `main` only, and `react-doctor.yml` on PRs to `main` only. | `git branch -a`; `.github/workflows/*.yml` | Add `master` (or use `branches: ['main', 'master']`), or rename the default branch. Until then, pushes to `master` get no Build Smoke run, and PRs into `master` get no react-doctor run. (Build Smoke's `pull_request` trigger has no branch filter, so PRs still get it.) |
| 2 | **CLAUDE.md is 686 lines.** The official guidance is under about 200, and long files make rules get lost. | `wc -l CLAUDE.md` | Keep a map of 100 to 150 lines. Move the cross-platform contract into `.claude/rules/cross-platform.md` with `paths:` for `electron/**`, `native*/**`, `scripts/**`, the CDP section into a skill, and the completion report into the `/close-increment` skill. |
| 3 | **CLAUDE.md requires linting, but there's no lint command.** ESLint packages are installed, but there's no ESLint config and no `lint` script. | `package.json`; no `eslint.config.*` found | Add a config and a `lint` script (start in warn mode), or drop the requirement. |
| 4 | **Hooks are user-level only, and `.claude/` is gitignored.** Graph update and status live in `~/.claude/settings.json`, and `.gitignore` line 268 ignores the whole `.claude` directory, so clones, teammates and cloud sessions get no project hooks, skills or agents. | `.claude/` holds only `settings.local.json`; `git check-ignore -v .claude/settings.local.json` → `.gitignore:268:.claude` | Change the rule to `.claude/*` plus `!.claude/settings.json`, `!.claude/hooks/`, `!.claude/skills/`, `!.claude/agents/`, `!.claude/rules/`, keeping `settings.local.json` ignored. Then commit a project `.claude/settings.json` with the graph hooks plus format, protect and Stop-verify hooks (templates). |
| 5 | **No project skills.** Recurring procedures (completion report, status-doc update, baseline compare, the Electron-ABI test command) live in prose. | No `.claude/skills/` in the repo | Add `/verify-change`, `/close-increment` and `/baseline` skills. |
| 6 | **No single fast `verify` command.** `test:ci` does a full build plus every suite, which is too slow for a Stop hook. | `package.json` | Add a `verify` script (typecheck:electron + affected unit suites) for hooks, and keep `test:ci` for CI. |
| 7 | **Docs are split across `_docs/` and `docs/`**, and a memory note records an agent missing `docs/`. | Both folders exist | Link both indexes from CLAUDE.md's map, or merge them. |
| 8 | **Architecture rules aren't enforced by tooling** ("keys only in main", "renderer doesn't import electron main"). They're kept by docs and tests only. | No dependency-cruiser or boundaries config | Add import-boundary rules whose error messages explain the fix. |
| 9 | **Graph has no embeddings**, so `semantic_search_nodes` falls back to keyword matching. | `list_graph_stats`: `Embeddings: 0 nodes` | Optional: `pip install "code-review-graph[embeddings]"`, then `code-review-graph embed`. |
| 10 | **The PR template has no slot for AI evidence.** | `.github/PULL_REQUEST_TEMPLATE.md` | Add "Validation categories", "Commands executed" and "Remaining risks" (07 §9). |
| 11 | **Many root-level `*Harness.html` files and an `.xlsx`** clutter the agent's view of the root. | Repo root listing | Move them to `harness/` or `tools/`; the agent spends less effort on noise. |
