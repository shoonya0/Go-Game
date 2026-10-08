# 07: Verification and quality gates

Closing the loop is the strongest lever there is. This file covers what to check, how
hard each gate should block, how to test AI features and UIs, how to keep the codebase
clean over time, and how to report results honestly.

---

## 1. The verification ladder

Climb only as high as the change needs. Each rung up is slower and catches a different
class of problem.

| Rung | Tool (examples) | Catches | Cost | When |
|---|---|---|---|---|
| 1 Typecheck | `tsc --noEmit`, pyright, `cargo check` | Invented APIs, wrong shapes, null bugs | Seconds | Every change |
| 2 Lint + architecture rules | ESLint/Biome/ruff + dependency-cruiser / boundaries | Bug patterns, layer violations | Seconds | Every change |
| 3 Unit tests | Vitest/Jest/node:test/pytest | Logic errors | Seconds to a minute | Every change |
| 4 Integration tests | Real DB in a temp dir, recorded HTTP fixtures | Seams between modules | Minutes | Changes that cross modules |
| 5 E2E tests | Playwright (`_electron` for Electron) | Broken user flows | Minutes | Flows touched |
| 6 Agent UI check | agent-browser snapshot, click, screenshot; console/errors | Visual and interaction regressions, runtime errors | Minutes | UI changes |
| 7 Build / packaged artifact | Production build, installer, container | Path, ASAR, environment and signing issues | Minutes to an hour | Before release, and after build-config changes |
| 8 Physical / live | A human on real devices and operating systems | OS integration, permissions, hardware | Hours | Before release, and after platform-sensitive changes |

**One command per speed tier** (`&&` works in both cmd.exe and POSIX shells):

```json
{
  "scripts": {
    "typecheck": "tsc --noEmit",
    "lint": "eslint . && depcruise src --config .dependency-cruiser.cjs",
    "test": "vitest run",
    "verify": "npm run typecheck && npm run lint && npm run test",
    "verify:full": "npm run verify && npm run test:e2e && npm run build"
  }
}
```

Avoid Unix-only syntax in package scripts: `VAR=x cmd`, `rm -rf`, single-quoted globs.
On Windows, npm runs scripts through cmd.exe, which handles none of these. Use
`cross-env`, `rimraf` or Node scripts instead. Natively's `build-smoke.yml` comments
record exactly this kind of break shipping green on a macOS-only CI matrix.

---

## 2. Make check output useful to the agent

- **Short output with the failure first.** Use reporters that print failures plus a
  summary, not thousands of lines.
- **Error messages that say how to fix the problem.** OpenAI's custom linters included
  remediation steps in the message itself, so the agent could repair violations without
  a human. Example rule:

```js
// .dependency-cruiser.cjs
module.exports = {
  forbidden: [
    {
      name: 'ui-must-not-import-db',
      comment:
        'UI may not import the data layer. Call a function in src/services/ instead. ' +
        'See docs/01-architecture.md#layers.',
      severity: 'error',
      from: { path: '^src/ui' },
      to: { path: '^src/db' },
    },
    { name: 'no-circular', severity: 'error', from: {}, to: { circular: true } },
  ],
};
// Run: npx depcruise src --config .dependency-cruiser.cjs --output-type err-long
// (err-long includes each rule's comment in the output) (verify)
```

- **Deterministic tests.** Fix flaky tests right away. A flaky test teaches the agent to
  ignore red, which is worse than having no test.

---

## 3. Choose how hard each gate blocks

| Gate | Strength | Use for |
|---|---|---|
| "…then run verify" in the prompt | Weak (advice) | Interactive work you're watching |
| `/goal <condition>` | Medium. A separate evaluator re-checks after every turn. | A session-long target ("all tests in X pass") |
| **Stop hook** running `verify` | Strong and deterministic. Exit 2 blocks the agent from stopping. | Any session you might walk away from |
| **Review subagent** / `/code-review` | Semantic. A fresh model tries to refute the work. | Before commit |
| **Pre-commit** | Strong, local | Fast checks: format, lint, secrets |
| **CI required checks** | Strongest, shared | `verify:full` on push and PR, across the OS matrix |

Template Stop hook: [`templates/hooks/stop-verify.mjs`](./templates/hooks/stop-verify.mjs).
It runs only if files changed, and it lets the agent stop after one blocked retry, so it
can't trap a session in a loop.

---

## 4. Test baselines and pre-existing failures

- Record failing tests **at a known commit** in `docs/current-state.md`.
- Judge an increment by **"no new failures compared with the baseline"**. Compare by
  test name, not by count.
- If you doubt the baseline, rebuild it from a clean HEAD (`git stash -u` → rebuild →
  run → `git stash pop` → rebuild).
- Never let the agent "fix" a failing test by deleting it, skipping it, or weakening
  its assertion without an explanation you approve.

---

## 5. Testing features that use LLMs or other AI

The rule from Natively's test strategy (`_docs/test/README.md`):

> Never assert exact model text (non-deterministic, model-dependent). Assert
> **structure and invariants**: sections present, citations non-empty, scores
> within an expected band, required markers present, nothing fabricated. Where a
> deterministic check is needed, stub the LLM with a fixed responder.

| Layer | How |
|---|---|
| Unit | A fixed-responder stub; test prompt assembly, parsing, routing and fallbacks |
| Golden set | Input → expected **properties** (bands, required fields, forbidden content); commit a baseline and track history |
| Judge eval | An LLM grades outputs against a rubric. Use it for trends, not as a hard gate, and keep the judge prompt versioned. |
| Privacy/offline | A network guard that blocks sockets; assert that features which need consent open **no** connection when disabled (Natively: `npm run test:privacy`) |
| Live smoke | Opt-in only (behind an env flag), never in default CI |

---

## 6. UI verification by the agent (CDP)

The pattern Natively uses (CLAUDE.md "Agent UI testing via CDP", `scripts/dev-agent.mjs`):

1. **An isolated instance per worktree**, launched by a **Node script** (no shell
   scripts, so it works on macOS and Windows) that:
   - picks a free debugging port (never 9222 or 9229),
   - writes a gitignored `agent-browser.json` with the port and a session id,
   - uses a worktree-local data or `userData` folder,
   - uses its own dev-server port, and skips the single-instance lock **in dev only**.
2. Load the skills for your CLI version: `agent-browser skills get core` (and `electron`).
3. Then: `agent-browser tab` → `tab t2` → `snapshot -i` → interact →
   **re-snapshot after every UI change** (refs go stale) → `console` / `errors` /
   `network requests` / `screenshot`.
4. **Never** use `--auto-connect`, which can attach to someone else's app. **Never**
   enable the debugging port in packaged builds.
5. Know the limit: CDP reaches the renderer's DOM, console and network. Native
   behavior such as windows, shortcuts, permissions, audio and click-through still
   needs a **human on each OS**.

Why it exists, from the script's own header: a second agent once attached to a shared
debug port and "spent four rounds reading somebody else's renderer".

---

## 7. Security gates

- Before merge: `/security-review` on the branch diff.
- Pre-commit or CI: a secret scanner; dependency audit; the `semgrep` or `sonarqube`
  plugin if you want rule-based static analysis inside the loop.
- For every **new dependency**, confirm that it exists, check its download count,
  maintainer and last release, and record why it was added. The AI invents package
  names (01 §12).
- Keys and secrets live in env vars or the OS keychain, never in code, prompts, logs or
  test fixtures.
- Review every hook and MCP server you install. They run with your permissions.

---

## 8. Entropy control (keeping AI code clean)

| Cadence | Action |
|---|---|
| Every edit | Formatter hook; graph update hook |
| Every commit | Pre-commit: format check, fast lint, secret scan; `code-review-graph detect-changes --brief` (advisory) |
| Every PR | `verify:full` in CI; `/code-review`; `/simplify` if the diff feels heavy |
| Weekly (`/weekly-cleanup` or a `/schedule` routine) | `npx knip` (dead code and dependencies), `npx jscpd src` (duplication), docs-drift check, flaky-test list, one dependency upgrade; open **small** cleanup PRs |
| Monthly | Prune CLAUDE.md (`claude-md-management` / `/doctor`); turn repeated prompts into skills; check the graph's hub nodes and large functions; review installed plugins and MCP servers (remove unused ones) |
| Each release | Packaged-build checks; security deep scan; update the honest-ceiling and current-state docs |

Budget about 20% of agent time for this (Steinberger), or automate it on a schedule
(OpenAI's "garbage collection" agents).

---

## 9. Evidence report

Require this at the end of every increment, in the PR body and the status doc. It is
adapted from Natively's CLAUDE.md "Required completion report".

```markdown
### Change summary
- Files changed: …
- Behavior changed: …
- Modules affected (shared / platform-specific): …

### Validation (use exact categories)
- Covered by automated tests: <suites + counts, e.g. "unit 412/412, e2e 18/18">
- Verified by agent UI check: <flows + screenshot paths>
- Tested physically on <OS/device>: <what>
- Build validated on <OS>: <artifact>
- Reviewed but not executed: <what>
- Requires physical verification on <OS/device>: <what>

### Commands executed
- `npm run verify` → pass (exit 0, 41 s)
- …

### Baseline comparison
- New failures vs baseline @<sha>: none / <list>

### Remaining risks
- …
```

Never write "verified on all platforms" when only one was run. A claim must be backed
by the evidence listed.
