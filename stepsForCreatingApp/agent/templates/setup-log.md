# AI harness setup log

This is the setup agent's record and resume point. If the setup is interrupted, start
the agent again with the same prompt and it continues from the first step that isn't
DONE or SKIPPED.

| Field | Value |
|---|---|
| Started | YYYY-MM-DD |
| Target | `<absolute path>` |
| Playbook | `<absolute path to stepsForCreatingApp>` |
| OS / shell | `<win32 | darwin>` / `<Bash (Git Bash) | PowerShell | zsh>` |
| Mode | NEW / EXISTING / EXISTING+AI |
| Stack | |
| Platforms (CI matrix) | |
| Depth | Full / Minimal / Full without CI and pre-commit |
| Methodology plugin | none / superpowers / feature-dev / Spec Kit |
| Hosting | GitHub / GitLab / other |
| Setup branch | |
| Verify command | |

Status values: `TODO` · `DONE` · `SKIPPED` (with consequence) · `DEFERRED` (with
re-enable condition) · `BLOCKED` (with reason)

## Tools

| Tool | Needed? | Status | Version / evidence | Who installed |
|---|---|---|---|---|
| git | yes | TODO | | |
| node ≥ 20 | yes | TODO | | |
| gh (+ auth) | | TODO | | |
| uv / pipx / pip | | TODO | | |
| code-review-graph | | TODO | | |
| agent-browser (+ Chrome) | | TODO | | |
| stack toolchain | | TODO | | |

## Steps

| ID | Step | Status | Evidence (command → result) / decision |
|---|---|---|---|
| P1 | Target confirmed | TODO | |
| P2 | Platform + Node | TODO | |
| P3 | Project inspected | TODO | |
| P4 | Mode chosen | TODO | |
| P5 | Scope answers | TODO | |
| P6 | Setup log created | TODO | |
| T1–T8 | Tools (see table) | TODO | |
| N1 / E1 | Git + safety (branch) | TODO | |
| N2 / E2 | Vision interview / secrets check | TODO | |
| N3 / E3 | Stack ADR / graph build | TODO | |
| N4 / E4 | Architecture doc + review / architecture tour | TODO | |
| N5 / E5 | Roadmap + test plans / architecture-overview + current-state | TODO | |
| N6 / E6 | Scaffold + deps / checks found | TODO | |
| — / E7 | Baseline at `<sha>` | TODO | |
| — / E8 | Verify command | TODO | |
| — / E9 | Characterization tests | TODO | |
| H1 | .gitignore + .claude tracking | TODO | |
| H2 | Verify command | TODO | |
| H3 | Hook scripts copied, adapted, stdin-tested | TODO | |
| H4 | .claude/settings.json | TODO | |
| H5 | Skills + reviewer agent | TODO | |
| H6 | CLAUDE.md | TODO | |
| H7 | Docs skeleton | TODO | |
| H8 | Formatter + linter | TODO | |
| H9 | Architecture rules | TODO | |
| H10 | Pre-commit | TODO | |
| H11 | CI | TODO | |
| H12 | Plugins | TODO | |
| H13 | code-review-graph MCP | TODO | |
| H14 | Agent UI access | TODO | |
| H15 | Optional extras | TODO | |
| V1 | Restart | TODO | |
| V2 | Registration checks | TODO | |
| V3 | Live probes | TODO | |
| V4 | Context budget | TODO | |
| V5 | Log finalized | TODO | |
| V6 | Completion report | TODO | |
| V7 | Commit / first task | TODO | |

## Test baseline (EXISTING)

At `<sha>`:

| Command | Duration | Result | Failing tests (by name) |
|---|---|---|---|
| | | | |

## Decisions made by the user

- …

## Deferred items (and what re-enables them)

- …

## Completion report

(Written in V6.)
