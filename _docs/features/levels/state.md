# Levels: state

## Now (updated 2026-10-26 at `edb1ba6`)

- **Status:** ⏳ planned
- **Current task:** 1. unique background names ([tasks.md](./tasks.md))
- **Stopped at:** docs only: the user's per-level tiles are in [spec.md](./spec.md); no code changed
- **Open questions / blockers:** (user) which enemies go in which level; how day/night is
  picked (per visit, a timer, the clock); how to shrink the backgrounds for the WASM (task 2)
- **Don't forget:** the new backgrounds are all named `day.*` / `night.*`, so
  `TestAssetBaseNamesUnique` fails and `readAsset` would load the wrong one (lookup is by base
  name). The old `LBackground_*.png` were deleted on disk, but `Background_1`
  (`internal/core/world.go`) still points at `LBackground_3.png`: the game can't load level 1
  until task 1. The new backgrounds total ~110 MB and `//go:embed levels` embeds all of them
  (today's `web/game.wasm` is ~33 MB)

## Recent sessions (newest first, keep 5)

### 2026-10-26: feature started

- Changed: branch renamed `docs/ui-buttons-merged` → `feature/levels`; docs folder created
- Verified: `go test ./internal/core -run TestAssetBaseNamesUnique` → FAIL (day/night name clashes)
- Docs touched: levels README, spec, tasks, state; registry; current-state Active work

## Known issues

- Background files clash by base name and level 1's background path is gone (see "Don't forget")
