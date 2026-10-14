# AI harness setup log

This is the setup agent's record and resume point. If the setup is interrupted, start
the agent again with the same prompt and it continues from the first step that isn't
DONE or SKIPPED.

| Field | Value |
|---|---|
| Started | 2026-10-07 |
| Target | `D:\projects\Go-Game` |
| Playbook | `D:\projects\Go-Game\stepsForCreatingApp` |
| OS / shell | `win32` / Bash (Git Bash) via Claude Code |
| Mode | EXISTING+AI (merge) |
| Stack | Go 1.24 (module `player`), Ebitengine v2.9; local toolchain go1.26.7 |
| Platforms (CI matrix) | Web (WASM), Windows desktop → `ubuntu-latest`, `windows-latest` |
| Depth | Full |
| Methodology plugin | none (playbook loop) |
| Hosting | GitHub (`shoonya0/Go-Game`) |
| Setup branch | `chore/ai-harness` (from `main` @ `91400c1`) |
| Verify command | `node scripts/verify.mjs` (full: `--full`) |

Status values: `TODO` · `DONE` · `SKIPPED` (with consequence) · `DEFERRED` (with
re-enable condition) · `BLOCKED` (with reason)

## Tools

| Tool | Needed? | Status | Version / evidence | Who installed |
|---|---|---|---|---|
| git | yes | DONE | 2.55.0.windows.3 | pre-existing |
| node ≥ 20 | yes | DONE | v24.19.0 | pre-existing |
| gh (+ auth) | yes | DONE | 2.102.0, logged in as shoonya0 (not on session PATH until restart) | user |
| uv / pipx / pip | no | SKIPPED | code-review-graph already installed via pip/python | — |
| code-review-graph | yes | DONE | graph rebuilt: 21 files, 153 nodes, 584 edges | pre-existing |
| agent-browser (+ Chrome) | yes (web) | DONE (present) | 0.38.1 — but `open` hangs on the game page, see H14 | pre-existing |
| stack toolchain | yes | DONE | go1.26.7 windows/amd64 | pre-existing |
| golangci-lint | yes | DONE | 2.14.0 in `~/go/bin` (`go install …/v2/cmd/golangci-lint@latest`) | agent (approved) |
| gopls | for gopls-lsp | DONE (present) | `~/go/bin/gopls.exe`; needs `%USERPROFILE%\go\bin` on PATH | pre-existing |

## Steps

| ID | Step | Status | Evidence (command → result) / decision |
|---|---|---|---|
| P1 | Target confirmed | DONE | `D:\projects\Go-Game` |
| P2 | Platform + Node | DONE | `win32 x64 node v24.19.0` |
| P3 | Project inspected | DONE | inspect-project.mjs → Go, AI artifacts present, CI `ci.yml`, no pre-commit, no env files |
| P4 | Mode chosen | DONE | EXISTING+AI (merge) |
| P5 | Scope answers | DONE | Web + Windows · Full · no methodology · GitHub |
| P6 | Setup log created | DONE | this file |
| T1–T8 | Tools (see table) | DONE | check-tools.mjs; gh installed by user; golangci-lint by agent |
| E1 | Git + safety (branch) | DONE | `git switch -c chore/ai-harness`; user WIP (6 enemy*.go) left untracked, never staged |
| E2 | Secrets check | DONE | `git ls-files \| grep …` → none; no `.env` in history |
| E3 | Graph build | DONE | `code-review-graph build` → 21 files, 153 nodes, 584 edges (untracked WIP not indexed) |
| E4 | Architecture tour | DONE | 3 communities (core-player 122, cmd-game 7, server-env 4), 22 flows, 16 large nodes |
| E5 | architecture-overview + current-state | DONE | `docs/architecture-overview.md`, `docs/current-state.md`, `docs/README.md` |
| E6 | Checks found | DONE | from Makefile + ci.yml: vet, build, wasm build; tests: `go test ./...`; user approved running all |
| E7 | Baseline at `91400c1` | DONE | all green **with** WIP; committed HEAD alone fails (see below) |
| E8 | Verify command | DONE | `scripts/verify.mjs` (Node, no shell): gofmt → vet → lint new-from-merge-base → test; `--full` + native/wasm build + advisory lint |
| E9 | Characterization tests | SKIPPED | not requested; existing 9 tests cover level building + enemy AI |
| H1 | .gitignore + .claude tracking | DONE | removed `.claude/` and `CLAUDE.md` ignores (user approved); added `.claude/*` + `!` exceptions; `git check-ignore` → local ignored, shared tracked; `.mcp.json` stays ignored (absolute python path) |
| H2 | Verify command | DONE | `node scripts/verify.mjs` → 4/4 ✔; `--full` → 7/7 ✔ |
| H3 | Hook scripts | DONE | protect-files (+go.mod, web/game.wasm, wasm_exec.js, bin/, _asset_backups/; migrations rule dropped), guard-shell (+ `git push … main` = Render deploy), format-on-edit (gofmt), stop-verify; 15/15 stdin tests match expected exit codes. graph-update.mjs SKIPPED: user + project settings already run `code-review-graph update` |
| H4 | .claude/settings.json | DONE | merged (user approved): env, permissions, SessionStart current-state, PreToolUse protect/guard, PostToolUse format, Stop verify; existing graph hooks kept (user chose keep duplicates); JSON parses |
| H5 | Skills + reviewer agent | DONE | `.claude/skills/verify-change`, `.claude/skills/close-increment`, `.claude/agents/code-reviewer.md` (Go/WASM-specific checks) |
| H6 | CLAUDE.md | DONE | project section prepended (user approved), graph section verbatim; 105 lines |
| H7 | Docs skeleton | DONE | `docs/README.md` index; all links resolve |
| H8 | Formatter + linter | DONE | gofmt (hook + verify) + golangci-lint v2 `.golangci.yml` (standard linters); baseline 3 staticcheck issues (non-blocking via new-from-merge-base) |
| H9 | Architecture rules | DONE | depguard rules core/system/cmd/server; scratch violations in `cmd/` and `server/` produced fix-naming messages and turned verify red; scratch files deleted |
| H10 | Pre-commit | DONE | `.githooks/pre-commit` (sh → Node: gofmt on staged .go + go vet); `git config core.hooksPath .githooks`; passes on repo, blocks unformatted staged file in scratch repo |
| H11 | CI | DONE (not yet run) | `.github/workflows/verify.yml`, ubuntu + windows, `go-version-file: go.mod`, autocrlf off, golangci-lint v2.14.0; YAML reviewed by eye (no parser installed); existing `ci.yml` untouched |
| H12 | Plugins | DONE | `gopls-lsp`, `claude-md-management` installed `--scope project`, enabled. `context7`, `security-guidance` installed but `enabled: false` |
| H13 | code-review-graph MCP | DONE (pre-existing) | `.mcp.json` + CLAUDE.md graph section already present; no `install` re-run |
| H14 | Agent UI access | PARTIAL | wasm built, `go run ./server` served `200 application/wasm`; server log shows browser fetched `game.wasm`; `agent-browser open` hung >2 min (WASM render loop) → no screenshot. Fallback documented in verify-change |
| H15 | Optional extras | DONE / USER | PR template written; claude-md-management installed; `/install-github-app` → user runs it |
| V1 | Restart | TODO | user restarts Claude Code (`claude --continue`) |
| V2 | Registration checks | TODO | `/hooks`, `/mcp`, `/context` |
| V3 | Live probes | TODO | protect `.env.harness-probe`, guard, format, graph MCP, `/verify-change` |
| V4 | Context budget | TODO | |
| V5 | Log finalized | TODO | |
| V6 | Completion report | DONE (in chat, pre-restart) | |
| V7 | Commit / first task | DONE | grouped commits on `chore/ai-harness`; not pushed |

## Test baseline (EXISTING)

At `91400c1` + untracked enemy WIP, Windows, go1.26.7:

| Command | Duration | Result | Failing tests (by name) |
|---|---|---|---|
| `go vet ./...` | 17s | pass | — |
| `go build ./...` | 5s | pass | — |
| `go test ./...` | 6s | pass (9 tests) | — |
| `gofmt -l .` | <1s | clean | — |
| wasm build | 20s | pass | — |
| `golangci-lint run ./...` | ~2s | 3 issues | SA1019 ×2 `internal/core/hud.go:49-50`, SA9003 `internal/core/quadtree.go:217` |
| committed HEAD alone, `go build ./...` | — | **FAIL** | `cmd/main.go:24: undefined: core.EnemyRuntime` (+ NewEnemy, KindSuccubus) |
| GitHub CI `ci` on main (last 5 runs) | ~40s | **all failure** | `undefined: core.EnemyRuntime` (10-01), `undefined: newCombat` (09-30) |

## Decisions made by the user

- Mode EXISTING+AI; platforms Web + Windows; Full depth; no methodology plugin; GitHub.
- New branch `chore/ai-harness`; enemy WIP left untracked.
- Run all baseline checks; install golangci-lint (agent); gh (user).
- Track shared harness in git; keep `.mcp.json` + `settings.local.json` ignored.
- Keep duplicate project-level graph hooks.
- gopls-lsp + claude-md-management plugins; pre-commit; web smoke; PR template; `/install-github-app` later.
- `.gitattributes` added. Enemy WIP committed first, on this branch, as its own commit.
- Commits: author shoonya0, no co-author/AI trailers, grouped by task, ≤6 files each. `stepsForCreatingApp/` left untracked.

## Deferred items (and what re-enables them)

- **Web screenshot check:** re-try once agent-browser can attach to a page with a busy WASM loop (e.g. run `open` detached, or a dev-only `?pause` flag in `index.html`).
- **gopls-lsp:** active once `%USERPROFILE%\go\bin` is on PATH and Claude Code restarted.
- **CI run of verify.yml:** after the user pushes this branch; check `gh run list --workflow verify.yml` shows jobs > 0. Linux `go test` with Ebitengine not yet proven headless.
- **Line endings:** resolved — `.gitattributes` forces LF for `*.go` and `.githooks/*` (user approved).

## Completion report

See chat (V6) — to be appended after V3 probes.
