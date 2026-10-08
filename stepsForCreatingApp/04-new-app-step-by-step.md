# 04: Creating a new app with AI (stages 0 to 8)

Each stage lists its **goal**, the steps, a **prompt** to paste, the **files**
it produces, and a **done-when** check. Don't skip stage 5 (the harness). Teams that
skip it pay for it on every feature afterward.

```text
0 Decide → 1 Define → 2 Choose stack → 3 Architecture → 4 Roadmap
   → 5 Phase 0: scaffold + HARNESS → 6 Build phase by phase (06 loop)
   → 7 Harden + release → 8 Operate + maintain
```

Where these files live in the repo: [`templates/docs/`](./templates/docs/).

---

## Stage 0: Decide whether it's worth building (optional, 30 minutes)

**Goal:** don't spend weeks on an idea that fails a basic check.

```text
I'm considering building <idea> for <users>. Act as a skeptical product advisor.
Interview me (AskUserQuestion) about the problem, who has it, how they solve it today,
why now, distribution, and what would make this fail. Then give a go / no-go / spike-first
recommendation with the 3 riskiest assumptions and the cheapest test for each.
```

Optional skill: `ln-11-opportunity-evaluator` (product-discovery-suite).
**Done when:** you can name the riskiest assumption and how you'll test it.

---

## Stage 1: Define what you're building (vision and scope)

**Goal:** a written, testable description of the product. The agent will reread it
for months, so it has to be precise.

```text
I want to build <one-paragraph description>. Interview me in detail using the
AskUserQuestion tool. Cover users and their jobs-to-be-done, core requirements,
explicit NON-goals, platforms, data and privacy, offline needs, integrations,
constraints (budget, time, licensing), and how we'll know v1 succeeded.
Don't ask obvious questions. Dig into the hard parts I haven't considered.
Then write docs/00-vision-and-scope.md using templates/docs/vision-and-scope.md.
Include an "Honest ceiling" section: what is NOT achievable or is out of scope, and why.
```

**Files:** `docs/00-vision-and-scope.md`.
**Done when:** every requirement can be checked, every non-goal is explicit, and a
stranger could explain the product after reading it.

> Natively example: `_docs/00-vision-and-scope.md` plus an **"Honest ceiling"**
> section in `_docs/README.md` that lists what can't be done (e.g. "Hidden from Task
> Manager's process list: not achievable in user space"). Writing limits down stops
> the agent, and you, from promising or attempting the impossible.

---

## Stage 2: Choose the stack

**Goal:** technology the AI writes **well** and that you can **verify**.

Criteria, in order:

1. **Boring and popular.** OpenAI's harness team deliberately picked "boring
   technologies" because they have stable APIs, compose well, and appear widely in
   training data. Avoid brand-new frameworks for the core.
2. **Strong typing and fast feedback.** TypeScript (strict), Python with type hints and
   pyright, Go, Rust, Kotlin or Swift. The type checker is your cheapest sensor.
3. **Good test tooling**, including e2e for UIs (Playwright, or `_electron` for
   Electron).
4. **Official scaffolders** (`npm create vite@latest`, `create-next-app`,
   `cargo new`, `uv init`), so the base isn't invented.
5. **Platform constraints first.** Desktop, mobile, OS-level features, native modules
   and signing decide a lot (Natively's choices follow from stealth needs: Electron
   plus a Rust native module).

```text
Read docs/00-vision-and-scope.md. Propose 2-3 stack options. For each: why it fits the
requirements, what's risky, how we'd test it (unit / e2e / packaged), how well AI agents
write it, and the deployment/packaging story on every target platform. Recommend one.
Then write docs/adr/0001-stack.md (templates/docs/adr.md).
```

**Spike first if there's a "make-or-break" risk.** Build a throwaway prototype of the
riskiest part before anything else (Natively's milestone M1, "It's invisible", was a
*ship-or-kill gate*).
**Done when:** an ADR records the decision, the alternatives and the consequences.

---

## Stage 3: Architecture

**Goal:** a map the agent must follow, covering modules, boundaries, data flow and
security model.

```text
Using docs/00 and the ADRs, write docs/01-architecture.md (templates/docs/architecture.md):
process/module map, data flow for the 3 main user journeys, layers and the ALLOWED import
directions between them, where secrets live, error-handling and logging conventions,
platform adapters (if multi-platform), and testing seams (where to inject fakes).
Add an ASCII diagram. Keep it under 250 lines.
```

Then have a **second session** review it (fresh context = unbiased):

```text
Review docs/01-architecture.md as a staff engineer. Find: missing boundaries, places where
the agent will be tempted to couple layers, security holes, untestable designs, and anything
over-engineered for v1. Report only issues that matter; propose concrete edits.
```

**Done when:** you can draw the architecture from memory, and the allowed-imports
rules are written down so stage 5 can turn them into lint rules.

---

## Stage 4: Phased roadmap with exit criteria

**Goal:** a sequence of phases you can ship and test independently.

```text
Write docs/02-roadmap.md (templates/docs/roadmap.md). 5-8 phases. Each phase: goal,
deliverables, ONE demoable exit criterion, risks. Order by risk (riskiest/most foundational
first). Phase 0 is always "foundations & harness" with no features. Add milestones and a
"Global definition of done (every phase)" list. Then for each feature write a test plan in
docs/test/<feature>.md (templates/docs/test-plan.md): test IDs, fixtures, pass bar.
```

**Files:** `docs/02-roadmap.md`, `docs/test/*.md`, and one status doc per phase,
started when work on that phase begins.
**Done when:** each phase's exit criterion is something you could demo in two minutes.

> Natively example: `_docs/04-phased-roadmap.md`. Every phase has 🎯 goal, 📦
> deliverables, ✅ exit criterion and ⚠️ risk. Milestones are named after what users
> will feel ("It answers", "It preps me"). The global DoD includes "works in the
> **packaged** build, not just `npm start`".

---

## Stage 5: Phase 0, scaffold plus harness (most important)

**Goal:** an empty app that **boots** and is wrapped in a working harness before any
feature exists.

Do these in order. Each step is a prompt to the agent plus a check.

| # | Step | Check |
|---|---|---|
| 1 | `git init`; `.gitignore` (+ `.env`); `.gitattributes` (`* text=auto eol=lf` so line endings stay the same across Windows and macOS) | `git status` is clean |
| 2 | Scaffold with the **official generator**; strict compiler settings | App starts |
| 3 | Scripts: `dev`, `build`, `typecheck`, `lint`, `format`, `test`, `test:e2e`, and **`verify`** (= typecheck && lint && test) | `npm run verify` passes |
| 4 | Formatter and linter configured; **architecture rules** from stage 3 (dependency-cruiser / eslint-plugin-boundaries) with fix-it messages | Breaking a rule on purpose fails `lint` with a helpful message |
| 5 | One real unit test, one real e2e smoke test | Both run green |
| 6 | **CLAUDE.md**: run `/init`, then cut it down to the template's map shape ([`templates/CLAUDE.template.md`](./templates/CLAUDE.template.md)), at or under 150 lines | `/context` shows it loaded |
| 7 | **`.claude/settings.json`** hooks: format on edit, protect files, Stop → verify ([`templates/`](./templates/README.md)) | Edit a file and it's formatted; try editing `.env` and it's blocked |
| 8 | **`.claude/skills/`**: `verify-change`, `close-increment` | `/verify-change` produces the evidence report |
| 9 | **code-review-graph**: `install` + `build`; add the graph update hook | `code-review-graph status` shows nodes |
| 10 | **Agent UI access**: the dev server on a known port; agent-browser can `open`, `snapshot -i`, `click` | The agent can describe the screen from a snapshot |
| 11 | **Pre-commit** (husky/lefthook): format check + fast lint + secret scan | A bad commit is rejected |
| 12 | **CI** (GitHub Actions) on **push and PR**: `verify` + build (OS matrix if you ship to several platforms) | The first PR shows green checks, and the run has **more than 0 jobs** |
| 13 | **Docs skeleton**: `docs/README.md` index, `docs/current-state.md`, `docs/status/phase-0.md` | Linked from CLAUDE.md |
| 14 | **Secrets and config**: `.env.example`, secrets only in env or the OS keychain, never in code or prompts | Grep the repo for keys: nothing found |

**Phase 0 prompt:**

```text
We are in Phase 0 (docs/02-roadmap.md). No features. Set up the project skeleton and the
harness exactly as listed in stepsForCreatingApp/04 Stage 5 (steps 1-14), one step at a time.
After each step run its check and show me the output. Use only official scaffolders and
well-known packages; for every dependency you add, state its weekly downloads/maintainer and
why it's needed. Stop after each 3 steps for my review.
```

**Done when:** a fresh clone needs one setup command, `verify` is green locally and in
CI, the agent can drive the UI, and the docs index is linked from CLAUDE.md.

> If you'll run agents **in parallel worktrees**, make the app boot per worktree now:
> free ports, a data directory per worktree, and no single-instance lock in dev.
> Natively's `scripts/dev-agent.mjs` is a working example. It picks a free CDP port
> (never 9222/9229), writes a gitignored `agent-browser.json`, uses a worktree-local
> `userData`, and skips the single-instance lock **only when `!app.isPackaged`**.

---

## Stage 6: Build phase by phase

For each phase:

1. Start `docs/status/phase-N.md` (template) with the planned increments.
2. Run the **feature loop** ([`06-feature-loop.md`](./06-feature-loop.md)) for each
   increment. Each increment is one reviewable commit (or a few), with tests and
   updated status.
3. At the end of the phase, run the **exit criterion** demo and record the result,
   date and platform in the status doc.
4. Update `docs/current-state.md` (one-page snapshot) and CLAUDE.md (only lessons
   that will apply again).
5. Weekly: run the cleanup pass (07 §8).

```text
Start Phase <N> from docs/02-roadmap.md. Read the phase, its test plan, and
docs/current-state.md. Propose the increments (each independently verifiable, ordered by
risk), write them into docs/status/phase-<N>.md, and stop for my approval. Don't code yet.
```

---

## Stage 7: Harden and release

| Area | Do |
|---|---|
| **Security** | `/security-review`; the `claude-security` plugin for a deep pass; dependency audit; secret scan on the full history |
| **Packaged build** | Test the **packaged or deployed artifact**, not dev mode. ASAR paths, env vars, signing, permissions and update channels all differ. |
| **Performance** | Measure first (the `ln-44` performance skill or profiling), then optimize |
| **Observability** | Error tracking (Sentry and similar) wired to an MCP server or CLI, so the agent can read real production errors |
| **Signing and release** | Platform signing/notarization and installers; release notes from commits; `/release` skill |
| **Docs** | User-facing README, privacy and security docs, a changelog |

**Done when:** the release checklist in [`10-checklists.md`](./10-checklists.md) passes
on every target platform. Use the platform names exactly: "tested on Windows 11"
is not "tested on macOS".

---

## Stage 8: Operate and maintain

- **Weekly cleanup** (`/weekly-cleanup`, or a `/schedule` routine): dead code,
  duplication, drifted docs, flaky tests, dependency updates (`ln-42`).
- **Triage loop:** production error → issue → `/triage-bug` → failing test → fix → PR.
- **Rules upkeep:** monthly, run `claude-md-management` (or `/doctor`) to prune
  CLAUDE.md; move multi-step procedures into skills and folder-specific rules into
  `.claude/rules/`.
- **Re-verify assumptions** after major dependency or OS updates.
