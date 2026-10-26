# Levels: state

## Now (updated 2026-10-26 at `4a8307a`)

- **Status:** 🟡 in progress
- **Current task:** 2. keep the WASM small with ~110 MB of backgrounds ([tasks.md](./tasks.md)), the user decides
- **Stopped at:** task 1 done: backgrounds renamed `<level>_day.*` / `<level>_night.*`, level 1 loads
  `deepForest_day.png` (the old `LBackground_3.png`, same bytes). Only that file is committed; the
  other levels' backgrounds stay untracked until task 2
- **Open questions / blockers:** (user) which enemies go in which level; how day/night is
  picked (per visit, a timer, the clock); how to shrink the backgrounds (task 2)
- **Don't forget:** `//go:embed levels` embeds untracked files too, so local builds carry all
  ~110 MB of backgrounds while CI and Render carry only `deepForest_day.png`. The old
  `LBackground_1.png` is now `moltenCrater_day.png` (same bytes)

## Recent sessions (newest first, keep 5)

### 2026-10-26: unique background names (task 1)

- Changed: `assets/levels/*/background/*` renamed; `Background_1` (`internal/core/world.go`), `levels_test.go`
- Verified: `node scripts/verify.mjs` → passed (5 steps); `TestAssetBaseNamesUnique`, `TestAssetPathsAreEmbedded` pass
- Docs touched: levels tasks, state; current-state

### 2026-10-26: feature started

- Changed: branch renamed `docs/ui-buttons-merged` → `feature/levels`; docs folder created
- Verified: `go test ./internal/core -run TestAssetBaseNamesUnique` → FAIL (day/night name clashes)
- Docs touched: levels README, spec, tasks, state; registry; current-state Active work

## Known issues

- Local builds embed untracked backgrounds that CI doesn't (see "Don't forget")
