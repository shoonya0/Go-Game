# Phase 3b: EXISTING project (E1–E9)

Goal: make an existing codebase **readable** to the agent (graph, docs) and
**verifiable** (a known check command and a recorded baseline) before any harness piece
changes behavior. Background reading: `PLAYBOOK/05-existing-app-step-by-step.md`.

The order matters: safety → secrets → map → write down → checks → baseline. Everything
up to E7 is read-only apart from the docs you write.

---

## Merge mode (EXISTING+AI): read this first if AI config already exists

The project already has some of: `CLAUDE.md`, `AGENTS.md`, `.claude/`, `.mcp.json`,
`.cursor*`, `.github/copilot-instructions.md`.

- Treat existing instructions as **authoritative**. Read them fully before E3. Use their
  commands and conventions; don't contradict them.
- Never delete or rewrite existing rules. In phase 4, **add** only what's missing, and
  show each addition as a diff for approval.
- **AGENTS.md without CLAUDE.md:** Claude Code reads AGENTS.md natively. If you also need
  a CLAUDE.md, it must start with `@AGENTS.md`. Once a CLAUDE.md exists, Claude reads
  only CLAUDE.md by default.
- **Cursor/Copilot rules:** list them in the setup log. Offer to fold the useful ones into
  CLAUDE.md, but don't delete the originals.
- **User-level hooks:** read `~/.claude/settings.json`. If it already has
  code-review-graph or format hooks, don't add duplicate project hooks in H3. Note it
  instead.
- If `CLAUDE.md` is over about 200 lines, *propose* (don't do) moving sections into
  `.claude/rules/*.md` with `paths:` frontmatter (playbook 03 §1).

---

## E1. Safety

1. **Not a git repo?** STOP. Ask (header `Git`): "git init + commit the current state
   (Recommended, so every change can be undone)" / "Stop here". Never proceed without
   version control.
2. **Dirty tree** (`git status --porcelain` not empty)? STOP and show the files. Ask the
   user to commit or stash *their* work themselves, or approve that you continue with it
   left as is. Never stash or commit their changes for them unasked.
3. **Setup branch.** Ask (header `Branch`): "Create `chore/ai-harness` from <current>
   (Recommended)" / "Use the current branch". On yes: `git switch -c chore/ai-harness`.
4. Note the default branch (from `inspect-project.mjs`: `git.defaultBranch`). CI
   triggers in H11 must match it.

## E2. Secrets check (read-only)

```bash
git ls-files | grep -iE '(^|/)\.env($|\.)|\.pem$|\.p12$|id_rsa|credentials' 
git log --all --oneline -- .env .env.local .env.production | head
```

PowerShell variant for the first line:
`git ls-files | Select-String -Pattern '(^|/)\.env($|\.)|\.pem$|\.p12$|id_rsa|credentials'`

- Tracked secret files, or secrets in history? STOP and report the findings (paths and
  commits, **never the secret values**). Recommend rotating the keys. Rewriting history
  is the user's decision and outside this setup; don't do it.
- `.env.example` / `.env.sample` are fine.

## E3. Build the code graph (needs T4)

The `install` step changes config, so it waits for H13. `build` only writes the local
cache, but still ask, because the first build can take minutes on big repos:

```bash
code-review-graph build --repo "<TARGET>"
code-review-graph status --repo "<TARGET>"
```

If code-review-graph was skipped, use narrow Glob/Grep for E4 and say so in the log.

## E4. Architecture tour (read-only)

Use the graph CLI (the MCP tools exist only after H13 and a restart):

```bash
code-review-graph architecture --repo "<TARGET>"
code-review-graph communities  --repo "<TARGET>"
code-review-graph large-functions --repo "<TARGET>"
code-review-graph flows --repo "<TARGET>"
```

(Run `code-review-graph <cmd> --help` if a flag differs.) Then read **only** the entry
points, the build and CI config, and the top hub files. Find out: what the app does,
its modules and dependencies, the 3 to 5 main flows, where state and secrets live, how
it's built, run and tested, and the 10 riskiest areas (`file:line`).

For a large legacy codebase, also suggest (don't install yet) the `code-modernization`
plugin's `/modernize` assessment (playbook 03 §3B).

## E5. Write the state down → docs (show each one, get approval)

Using the templates in `PLAYBOOK/templates/docs/`:

- `docs/architecture-overview.md`: module map, flows, and the layers or import
  directions you **want** (marked as "target" where the code doesn't follow them yet).
- `docs/current-state.md`: what works, what's broken, known issues, how to
  build/run/test, open decisions. Mark every claim **VERIFIED** (you ran it or read the
  code) or **ASSUMED**.
- `docs/README.md`: index. If `docs/` already holds user documentation, add an "AI /
  engineering docs" section instead of replacing anything.

## E6. Find the checks (read-only)

From package scripts, Makefile/justfile, CI workflows, README and CONTRIBUTING, list
every typecheck, lint, test, build and e2e command. Show the list and ask (header `Run
checks`): "Run all now" / "Run only the safe ones (I'll mark which)" / "Don't run;
I'll tell you the baseline". Some suites hit real services. Ask which ones need
env vars or network.

## E7. Baseline (runs project commands, so only with E6 approval)

Run the approved commands. Record in `docs/current-state.md` → **Test baseline at
`<sha>`**: each command, its duration, pass/fail, and **the names of the failing
tests**. From now on, compare by test name, never by count.

## E8. Verify command (a small additive change, so ask first)

Propose ONE fast `verify` (typecheck + lint + unit, under about 2 minutes) and ONE
`verify:full`:

- **JS/TS:** add `"verify"` / `"verify:full"` scripts to `package.json`, composing the
  existing scripts with `&&` (works under cmd.exe and sh). Show the diff and get
  approval.
- **Other stacks:** copy `PLAYBOOK/templates/scripts/verify.mjs` to `scripts/verify.mjs`
  and set its STEPS from `reference/stack-matrix.md`.

If the baseline has failures, `verify` will be red. Record that **the Stop hook (H3)
stays off until `verify` is green** (or until the user approves a verify that excludes
the named baseline failures).

## E9. Characterization tests (optional)

Ask (header `Safety net`): "Name a first area you'll change, and I'll pin its current
behavior with tests (Recommended)" / "Skip for now". If chosen, write tests that record
**current** behavior, including odd behavior marked with TODO notes, and run them
green on unchanged code. Use the graph (`code-review-graph query` / `impact`) to find
existing coverage first.

✅ **Phase 3b checkpoint:** summarize the branch, secrets result, graph stats, docs
written, baseline and verify command. Ask to continue, then open `05-harness.md` in
**incremental** mode: lint rules start as warnings, and the Stop hook waits for green.
