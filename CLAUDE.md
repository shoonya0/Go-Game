# Go-Game (Pirate Adventure)

2D side-scrolling platformer built from scratch on Ebitengine. Stack: Go 1.24 (module
`player`), Ebitengine v2.9, `go:embed` assets. Platforms: Windows desktop (native) and
web (WebAssembly, deployed to Render via Docker).

## Read first (map, not manual)

| Need | Read |
|---|---|
| Where the project stands, active work, known issues, test baseline | `_docs/current-state.md` |
| Docs index; which features exist and where their docs live | `_docs/README.md` → Feature registry |
| Modules, layers, allowed imports, main flows | `_docs/architecture-overview.md` |
| Build targets, WASM build, Render deploy | `DEPLOY.md` |
| AI harness setup record | `_docs/ai-harness-setup.md` |

## Commands

- Run native: `go run ./cmd`
- **Fast check, run before saying done: `node scripts/verify.mjs`** (gofmt + vet + new lint issues + unit tests + docs-check)
- Docs changed: `node scripts/docs-check.mjs` (links, registry, budgets, stale feature state)
- Full check: `node scripts/verify.mjs --full` (adds native build, WASM build, full lint report)
- Single test: `go test ./internal/core -run TestName -v`
- Web locally: build wasm (see `DEPLOY.md`, PowerShell section on Windows), then `go run ./server` → :8080
- `golangci-lint` lives in `%USERPROFILE%\go\bin` (`~/go/bin`); verify.mjs finds it there.

## Workflow

- Find code with the code graph (see below) before Grep or reading files.
- Non-trivial change: plan → implement in small steps → `/verify-change` → `code-reviewer`
  agent → `/close-increment`.
- A feature that spans sessions gets a docs folder: `/feature-docs start <feature>` (usually
  `_docs/features/<feature>/`, since most code shares `internal/core`). Read its `state.md`
  and `tasks.md` first, then only what its `README.md` lists for the task. Link facts, don't
  copy them; document only what the code can't tell (intent, reasons, progress, traps).
- Bugs: write a failing test that reproduces the bug first.
- Write tests in the same session as the feature. Never delete, skip or weaken a test
  without explaining why and getting approval.
- Report evidence (commands + output), not claims.

## Git

- Never commit on or push `main` (it auto-deploys to Render). Branch from an up-to-date base:
  `feature/<name>` (matches the feature-docs name), `fix/<issue#>-<name>`, `refactor/`, `docs/`, `chore/`.
- Author = the signed-in GitHub user (`gh api user`): name = login, email =
  `<id>+<login>@users.noreply.github.com`. No AI watermark (Co-Authored-By, "Generated with
  Claude Code", 🤖) in commits, PRs or issues. `.claude/hooks/git-guard.mjs` enforces all three.
- Conventional Commits (`type(scope): summary`, body says why, `Fixes #N`). Stage only the
  files of this increment, by path; review `git diff --staged`.
- After each commit, verify: `git log -1 --format='%an <%ae>%n%n%B'`, `git show --stat HEAD`, `git status --short`.
- Increments are committed locally. Feature done and verified → `/code-review` if non-trivial →
  push → PR (completion report in the body) → CI green → **ask the user before `gh pr merge`**
  (merging to `main` deploys to Render) → delete branch.
- Bug outside the current task → propose a GitHub issue (title, repro, expected/actual); don't fix it here.

## Rules (each one exists because something went wrong)

- IMPORTANT: Always save the current state of the app after every piece of work: run
  `/feature-docs update` (tick `tasks.md`, overwrite "Now" in the feature's `state.md`,
  update "Active work" in `_docs/current-state.md`; with no feature folder, just Active
  work). A Stop hook checks this.
- Every new `.go` file must be committed with the code that uses it: `main` once
  referenced enemy types whose files were never committed, so it didn't build.
- Respect the layers table in `_docs/architecture-overview.md`;
  `.golangci.yml` depguard enforces them. `server/` stays stdlib-only and CGO-free.
- Assets are loaded by **base file name** from the embed FS (`internal/core/image.go`):
  new asset files need unique names, even in subfolders.
- Code must compile for both `GOOS=windows` and `GOOS=js GOARCH=wasm`: no OS-specific
  APIs (files, goroutine-per-frame, syscalls) in `internal/core`.
- CI uses Go 1.24 while local is newer: don't use language/library features newer than `go.mod`.
- New dependency: say why and ask first; change `go.mod` with `go get`, never by hand.
- Draw order in `Game.Draw` is intentional (water after player).

## UI verification

- Native window: manual only. Web build: agent-browser against `go run ./server`
  (http://localhost:8080) — the game renders to a `<canvas>`, so check console errors and
  screenshots rather than DOM elements.

## When compacting

Keep: the list of modified files, the current task, and the verify commands with their latest results.

<!-- code-review-graph MCP tools -->
## MCP Tools: code-review-graph

**This project has a knowledge graph. Start with the code-review-graph
MCP tools to narrow scope, then read the source.** The graph is cheaper than scanning files and
gives you structural context (callers, dependents, test coverage) that file search cannot.

### When to use graph tools FIRST

- **Exploring code**: `semantic_search_nodes_tool` or `query_graph_tool` instead of Grep
- **Understanding impact**: `get_impact_radius_tool` instead of manually tracing imports
- **Code review**: `detect_changes_tool` + `get_review_context_tool` instead of reading entire files
- **Finding relationships**: `query_graph_tool` with callers_of/callees_of/imports_of/tests_for
- **Architecture questions**: `get_architecture_overview_tool` + `list_communities_tool`

### Verify in the source

- Narrow scope with the graph, then read the source. Do not change code from graph output alone.
- For any non-trivial change, read the implementation and the relevant tests before concluding.
- Verify the exact source when touching behavior, database logic, migrations, retries, fallbacks,
  recovery, or compatibility code.
- When the graph and the source disagree, the source wins. The graph may be stale or may not
  model that relationship.
- An empty graph result can mean "not indexed" or "not statically visible", not "does not exist".

### Key Tools

| Tool | Use when |
| ------ | ---------- |
| `detect_changes_tool` | Reviewing code changes — gives risk-scored analysis |
| `get_review_context_tool` | Need source snippets for review — token-efficient |
| `get_impact_radius_tool` | Understanding blast radius of a change |
| `get_affected_flows_tool` | Finding which execution paths are impacted |
| `query_graph_tool` | Tracing callers, callees, imports, tests, dependencies |
| `semantic_search_nodes_tool` | Finding functions/classes by name or keyword |
| `get_architecture_overview_tool` | Understanding high-level codebase structure |
| `refactor_tool` | Planning renames, finding dead code |

### Workflow

1. The graph auto-updates on file changes (via hooks).
2. Use `detect_changes_tool` for code review.
3. Use `get_affected_flows_tool` to understand impact.
4. Use `query_graph_tool` pattern="tests_for" to check coverage.
<!-- /code-review-graph MCP tools -->
