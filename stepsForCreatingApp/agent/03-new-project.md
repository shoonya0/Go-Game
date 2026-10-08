# Phase 3a: NEW project (N1–N7)

Goal: turn an idea into written, checkable docs and a scaffolded app that boots, so
phase 4 can wrap it in the harness. Background reading: `PLAYBOOK/04-new-app-step-by-step.md`
stages 0 to 5.

Use **AskUserQuestion** for every interview round (at most 4 questions per call; 2 to 4
options each; the user can always answer "Other"). Write each doc as soon as its
interview ends, show a short summary, and ask for approval before moving on.

---

## N1. Folder and git

- If TARGET doesn't exist: create it (the user already confirmed the path in P1).
- If it isn't a git repo: ask (header `Git`) → "Initialize git (Recommended)" / "Not
  now". On yes:
  ```bash
  git init
  ```
  Then create `.gitignore` (language-appropriate, plus the harness entries listed in
  H1) and `.gitattributes` with `* text=auto eol=lf`. Line endings stay identical
  across Windows and macOS. That's safe in a new repo; never add it to an existing
  repo without asking (renormalizing creates huge diffs).
- Don't commit yet.

## N2. Vision interview → `docs/00-vision-and-scope.md`

Interview in rounds. Round 1: problem, users, the job to be done, why now. Round 2: core
requirements, explicit **non-goals**, platforms, data and privacy. Round 3: offline
needs, integrations, constraints (time, budget, licensing), success criteria. Dig into
the hard parts and skip obvious questions.

Write `docs/00-vision-and-scope.md` from `PLAYBOOK/templates/docs/vision-and-scope.md`.
Every requirement needs a "How we'll verify it" column. Include the **Honest ceiling**
table: what is not achievable or not attempted.

✅ The user approves the doc.

## N3. Stack decision → `docs/adr/0001-stack.md`

Propose 2 to 3 stacks. For each, cover fit with the requirements, risks, how it's
tested (unit, e2e, packaged), how well AI agents write it, and packaging or deployment
on every target platform. Prefer **boring, popular, strongly typed** stacks with
**official scaffolders**. Recommend one and ask (header `Stack`).

Write the ADR from `PLAYBOOK/templates/docs/adr.md`. If there's a make-or-break
technical risk, propose a **timeboxed spike** before continuing. The spike's code is
thrown away; only its findings are kept.

Then run **T6** (phase 2) for this stack's toolchain.

## N4. Architecture → `docs/01-architecture.md`

Write it from `PLAYBOOK/templates/docs/architecture.md`: module map, the 3 main flows,
**layers and allowed imports** (they become lint rules in H9), secrets, errors and
logging, platform adapters if multi-platform, and testing seams. Keep it under about
250 lines.

**Fresh-context review:** use a subagent (Agent tool, general-purpose) with this prompt:
"Review docs/01-architecture.md and docs/00-vision-and-scope.md as a staff engineer.
Report only material issues: missing boundaries, coupling traps, security holes,
untestable designs, over-engineering for v1. Propose concrete edits." Show the findings
to the user and apply the approved edits.

## N5. Roadmap → `docs/02-roadmap.md` and `docs/test/`

Write it from `PLAYBOOK/templates/docs/roadmap.md`: 5 to 8 phases ordered by risk,
**Phase 0 = foundations and harness (no features)**, one demoable exit criterion per
phase, milestones, and the global definition of done. Write one test plan per
early feature from `PLAYBOOK/templates/docs/test-plan.md`.

✅ The user approves the roadmap.

## N6. Scaffold (installs something, so use the Stop protocol)

1. Pick the **official generator** for the stack, with a non-interactive form if one
   exists. Examples (verify the current flags with the generator's `--help`):

   | Stack | Command |
   |---|---|
   | React + Vite (TS) | `npm create vite@latest .scaffold-tmp -- --template react-ts` |
   | Next.js | `npx create-next-app@latest .scaffold-tmp --ts --eslint --app --use-npm` |
   | Electron + Vite | see electron-vite / Electron Forge docs |
   | Python (uv) | `uv init --package .scaffold-tmp` |
   | Go | `go mod init <module>` (no temp folder needed) |
   | Rust | `cargo new .scaffold-tmp` |
   | .NET | `dotnet new <template> -o .scaffold-tmp` |

2. Generators usually refuse non-empty folders, and `docs/` already exists, so
   scaffold into **`TARGET/.scaffold-tmp/`**.
3. Stop protocol. Show the exact command and ask: "I'll run it myself" (needed if the
   generator asks questions), "You run it", or "Different generator".
4. Move the scaffold's contents up into TARGET **without overwriting existing files**.
   If a file conflicts (e.g. `.gitignore`, `README.md`), merge it and show the result.
   Remove `.scaffold-tmp/`. Deleting this temporary folder you created is fine; say you
   did it.
5. Installing dependencies (`npm install`, `uv sync`, …) is an install, so use the Stop
   protocol.
6. Turn on strict compiler settings (TS `"strict": true`, pyright `strict` for new code,
   and so on).

✅ Check: the app builds and starts (`npm run build` / `npm run dev`, or the stack's
equivalent; ask before the first run). Show the output.

## N7. Hand over to the harness

Record N1 to N6 in the setup log, including the stack, the scaffold command and the
dependency install. Set `LANG` and the planned `VERIFY_CMD`.

✅ **Phase 3a checkpoint:** summarize the docs created and the scaffold state. Ask to
continue, then open `05-harness.md`. In NEW mode, phase 4 *is* the roadmap's Phase 0.
