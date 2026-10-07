---
name: verify-change
description: Run this project's verification ladder (gofmt, vet, lint, unit tests, native + WASM build, and a web smoke check when relevant) and report evidence in a fixed format. Use before claiming any change is done, after implementing a plan step, or when the user asks "does it work?".
---

# Verify change

Prove the current change works. Report **evidence, not claims**.

## 1. Scope

- `git status` and `git diff --stat`: what changed? Make sure every new `.go` file the
  change depends on is listed (untracked files break `main` once committed without them).
- Choose the rungs needed. Rung 1 always runs.

| Rung | Command | Run when |
|---|---|---|
| 1 Fast | `node scripts/verify.mjs` (gofmt + vet + new lint issues + `go test ./...`) | Always |
| 2 Full | `node scripts/verify.mjs --full` (adds native build, WASM build, full lint report) | `cmd/`, `server/`, assets, `go.mod`, anything platform-sensitive |
| 3 Web smoke | build wasm + `go run ./server`, then agent-browser (below) | Rendering, input, asset or `web/` changes |
| 4 Native run | `go run ./cmd` — **user** plays it | Gameplay feel, physics, animation, controls |

## 2. Run

- Run each chosen rung. On failure, find the **root cause** and fix it, then rerun
  from rung 1.
- Never delete, skip or weaken a test to get green. If a test is wrong, stop and
  explain why.
- Compare failures against the **baseline** in `docs/current-state.md` by test name.
  Known lint baseline (not blocking): 2× SA1019 `vector.DrawFilledRect` in
  `internal/core/hud.go`, SA9003 empty branch in `internal/core/quadtree.go`.

## 3. Web smoke check (rung 3)

Ask the user before starting the server the first time in a session.

```text
# Windows PowerShell build (see DEPLOY.md):
$env:GOOS="js"; $env:GOARCH="wasm"; go build -ldflags="-s -w" -o web/game.wasm ./cmd; Remove-Item Env:GOOS, Env:GOARCH
Copy-Item "$(go env GOROOT)/lib/wasm/wasm_exec.js" web/wasm_exec.js
go run ./server                        # background; serves http://localhost:8080

agent-browser skills get core          # once per session
agent-browser open http://localhost:8080
agent-browser wait 3000                # wasm boot
agent-browser console
agent-browser errors
agent-browser screenshot docs/evidence/<feature>-<step>.png
agent-browser close
```

The game draws to a `<canvas>`: there are no DOM elements to snapshot. Judge from the
screenshot and console (no "failed to load", no panics). Gameplay feel needs rung 4.

**Known issue (2026-10-07):** `agent-browser open http://localhost:8080` hung (>2 min,
not killed by `timeout`) on this page, likely because the WASM render loop never lets the
page go idle. Run agent-browser steps in the background, and if `open` hangs, fall back to:
`curl -s -o /dev/null -w "%{http_code} %{content_type}" http://localhost:8080/game.wasm`
(expect `200 application/wasm`) plus the server log showing `GET /game.wasm`, and list
"web rendering" under "Requires physical verification".

## 4. Report (exact format)

```markdown
### Verification
| Rung | Command | Result | Notes |
|---|---|---|---|
| Fast | `node scripts/verify.mjs` | ✅ pass (4 s) | 9/9 tests |
| Full | `node scripts/verify.mjs --full` | ✅ pass | native + wasm build ok |
| Web | agent-browser | ✅ | console clean; docs/evidence/x-1.png |

- New failures vs baseline @<sha>: none
- Not verified: <list, or "nothing">
- Requires physical verification (native window / gameplay feel): <list, or "nothing">
```

If any rung fails and can't be fixed in scope, stop and report it. Don't claim done.
