# <Project name>

<One sentence: what the app is and who it's for.> Stack: <e.g. TypeScript, React, Vite, Node 22>.
Platforms: <web | Windows + macOS | iOS + Android>.

## Read first (map, not manual)

| Need | Read |
|---|---|
| Where the project stands, next task, test baseline | `docs/current-state.md` |
| What we're building and NOT building | `docs/00-vision-and-scope.md` |
| Modules, layers, allowed imports, data flow | `docs/01-architecture.md` |
| Phases and exit criteria | `docs/02-roadmap.md` → `docs/status/phase-<N>.md` |
| Feature specs and test plans | `docs/specs/`, `docs/test/` |
| Why a decision was made | `docs/adr/` |

## Commands

- Install: `npm ci`
- Dev: `npm run dev` (agent UI instance: `npm run dev:agent`)
- **Fast check, run before saying done: `npm run verify`** (typecheck + lint + unit)
- Full check: `npm run verify:full` (adds e2e + build)
- Single test: `npx vitest run path/to/file.test.ts`

## Workflow

- Find code with the code graph (code-review-graph MCP: `semantic_search_nodes`,
  `query_graph` callers_of/tests_for, `get_impact_radius`) before Grep or reading files.
- Non-trivial change: spec → plan → implement in small steps → `/verify-change` → `/code-review`
  → `/close-increment`.
- Bugs: write a failing test that reproduces the bug first.
- Write tests in the same session as the feature. Never delete, skip or weaken a test
  without explaining why and getting approval.
- Commit only the files you changed for this increment.
- Report evidence (commands + output, screenshots), not claims. Use the validation
  categories in `/close-increment`.

## Rules (each one exists because something went wrong)

- IMPORTANT: Secrets only from env or the OS keychain. Never in code, logs, fixtures or prompts.
- New dependency: say why, and check that it exists and is maintained, before installing.
- Respect layer boundaries in `docs/01-architecture.md#layers`; `npm run lint` enforces them.
- <Project-specific gotcha, e.g. "Tests that touch the DB need `npm run test:db`, not plain vitest.">
- <Project-specific convention that differs from defaults.>

## Platform notes (delete if single-platform)

- Platform-specific code lives behind adapters in `src/platform/{darwin,win32}`; never
  scatter `process.platform` checks. Tests inject the platform.
- Package scripts must run under cmd.exe and sh: no `VAR=x cmd`, no `rm -rf`; use
  `cross-env` / `rimraf` / Node scripts.

## When compacting

Keep: the list of modified files, the current task and its spec path, and the test
commands with their latest results.
