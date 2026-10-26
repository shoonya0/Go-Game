# Levels: state

## Now (updated 2026-10-26 at `4a8307a`)

- **Status:** 🟡 in progress
- **Current task:** 3. register the JPEG decoder, then 4. fetch backgrounds per sub-level on the web ([tasks.md](./tasks.md))
- **Stopped at:** task 1 done: backgrounds renamed `<level>_day.*` / `<level>_night.*`, level 1 loads
  `deepForest_day.png` (the old `LBackground_3.png`, same bytes). Only that file is committed; the
  other levels' backgrounds stay untracked until task 2
- **Open questions / blockers:** (user) the sub-level order and the D5 rules I chose
  (trigger once, 2+ portals per sub-level, first portal is the entry). Settled: D1–D5
- **Don't forget:** `//go:embed levels` embeds untracked files too, so local builds carry all
  ~110 MB of backgrounds while CI and Render carry only `deepForest_day.png`. The old
  `LBackground_1.png` is now `moltenCrater_day.png` (same bytes). Saves store the level name
  (`deepForest1`): renaming it for the day/night split breaks old saves' checkpoint

## Recent sessions (newest first, keep 5)

### 2026-10-26: the user's level rules (D1–D5)

- Changed: docs only: decisions.md D1–D5, spec rules, tasks 2–8
- Verified: `node scripts/docs-check.mjs` → 0 errors
- Docs touched: levels decisions, spec, tasks, state; enemy-types D10; current-state

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
