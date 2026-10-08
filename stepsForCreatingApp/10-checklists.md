# 10: Checklists

Short, printable lists. Each points back to the file with the details.

---

## A. New project setup (04 stages 1–5)

- [ ] `docs/00-vision-and-scope.md`: requirements, non-goals, honest ceiling, success criteria
- [ ] Stack ADR (`docs/adr/0001-stack.md`); spike done if there's a make-or-break risk
- [ ] `docs/01-architecture.md` with allowed import directions; reviewed in a fresh session
- [ ] `docs/02-roadmap.md`: phases with one demoable exit criterion each, plus a global DoD
- [ ] Test plan per feature in `docs/test/`
- [ ] `git init`, `.gitignore` (with `.env`), `.gitattributes` (`eol=lf`)
- [ ] Scaffolded with the official generator; strict compiler settings
- [ ] Scripts: `typecheck`, `lint`, `test`, `test:e2e`, `build`, **`verify`**, `verify:full`
- [ ] Architecture rules with fix-it messages
- [ ] One unit test and one e2e smoke test, both green
- [ ] CLAUDE.md as a map, at or under 150 lines (`/context` shows it loaded)
- [ ] `.claude/settings.json` hooks: format, protect files, Stop → verify, graph update
- [ ] `.claude/skills/verify-change`, `.claude/skills/close-increment`; `code-reviewer` subagent
- [ ] `code-review-graph install && build`
- [ ] Agent can drive the UI (agent-browser `open`, `snapshot -i`, `click`)
- [ ] Pre-commit: format check, fast lint, secret scan
- [ ] CI on **push and PR**, OS matrix if multi-platform; first run shows **more than 0 jobs**; branch filters match the real default branch
- [ ] `docs/current-state.md` and `docs/status/phase-0.md`, linked from CLAUDE.md

## B. Taking over an existing app (05 steps 0–6)

- [ ] On a branch; working tree clean; `.env` gitignored and absent from history
- [ ] Irreversible paths blocked by a `PreToolUse` hook
- [ ] Graph built; `claude-code-setup` recommendations reviewed
- [ ] Read-only architecture tour done (plan mode, graph tools)
- [ ] CLAUDE.md from `/init`, cut down to a map
- [ ] `docs/current-state.md` with VERIFIED vs ASSUMED claims and the test **baseline @ sha**
- [ ] Single `verify` command; CI confirmed to actually run
- [ ] Characterization tests around the first area to be changed
- [ ] Harness turned on step by step (hooks → Stop verify → CI → lint in warn mode → pre-commit → skills)

## C. Start of every session (06 steps 1–3)

- [ ] One task, sized for a review of about 15 minutes
- [ ] `/clear` or a new named session
- [ ] Agent reads CLAUDE.md, the docs index, and current-state
- [ ] Graph lookup of the code involved; relevant docs read; contradictions noted

## D. Before saying "done" (06 steps 6–8)

- [ ] Spec acceptance criteria each marked PASS/FAIL **with evidence**
- [ ] `verify` green; no new failures compared with the baseline (by test name)
- [ ] UI changes: snapshot, flow, console/errors clean, screenshot compared
- [ ] `/code-review` or a reviewer subagent run; real gaps fixed
- [ ] No deleted or skipped tests, no new `any` or ignores without approval
- [ ] New dependencies checked (exists, reputable, needed)
- [ ] Completion report written with exact validation categories (07 §9)

## E. Before merge (06 step 10)

- [ ] Commit contains only this increment's files; message says what and why
- [ ] PR body = completion report
- [ ] CI green on every OS in the matrix
- [ ] You read: the summary, diff stats, risky files, and evidence
- [ ] Status doc and current-state updated **in the same PR**
- [ ] A lesson added to CLAUDE.md, or a hook added, if a repeat mistake happened

## F. Weekly upkeep (07 §8)

- [ ] `npx knip` (dead code and dependencies), `npx jscpd src` (duplication)
- [ ] Docs-drift check (`docs-drift-checker` subagent)
- [ ] Flaky tests listed and fixed or quarantined with an issue
- [ ] One dependency upgrade (changelog read, `verify:full`)
- [ ] Small cleanup PRs opened, each reviewed like any other change

## G. Monthly upkeep

- [ ] CLAUDE.md pruned (`claude-md-management` / `/doctor`); procedures moved into skills
- [ ] Unused plugins and MCP servers removed (`/context` to see their cost)
- [ ] Graph hub nodes and large functions reviewed for refactoring
- [ ] Hooks, skills and subagents still match how you actually work

## H. Release (04 stage 7)

- [ ] `verify:full` green on every target platform in CI
- [ ] **Packaged or deployed artifact** tested, not dev mode
- [ ] `/security-review` plus a deep scan; dependency audit; secret scan of the full history
- [ ] Signing, notarization and installer checked on each OS that needs them
- [ ] Error tracking wired up and readable by the agent
- [ ] Physical checks done and listed per platform; anything unverified listed under "Remaining risks"
- [ ] Changelog and release notes; honest-ceiling and current-state docs updated
