# 02: Phased roadmap

Each phase can be shipped and tested on its own and has **one demoable exit criterion**.
Build in order; later phases assume earlier invariants.

Legend: 🎯 goal · 📦 deliverables · ✅ exit criterion · ⚠️ risk · 🧪 test plan

---

## Phase 0: Foundations and harness (no features)

🎯 An empty app that boots, wrapped in a working harness.

📦
- Scaffold (official generator), strict compiler settings
- Scripts: `typecheck`, `lint` (+ architecture rules), `test`, `test:e2e`, `build`, `verify`, `verify:full`
- CLAUDE.md map, `.claude/settings.json` hooks, `/verify-change` and `/close-increment` skills
- Code graph built; agent can drive the UI
- Pre-commit and CI (push + PR, OS matrix if multi-platform)
- Docs skeleton: current-state, status/phase-0

✅ A fresh clone needs one setup command; `verify` is green locally and in CI (more than
0 jobs); the agent can open, snapshot and click the empty app.

⚠️ <e.g. pin runtime/ABI versions early>

---

## Phase 1: <Riskiest / most foundational capability>

🎯
📦
✅
⚠️
🧪 `docs/test/<feature>.md`

---

## Phase 2: …

---

## Milestones

| Milestone | Phases | What the user can feel |
|---|---|---|
| M1 "<it works at all>" | 0–1 | A ship-or-kill gate for the core risk |
| M2 | | |

## Global definition of done (every phase)

- [ ] Works in the **packaged or deployed** artifact, not only in dev mode
- [ ] `verify:full` green in CI on every target platform
- [ ] No new network calls or data collection beyond what's documented
- [ ] Security review run; no secrets in code, logs or fixtures
- [ ] Status doc and current-state updated; completion report written
- [ ] What wasn't physically verified is listed explicitly
