# Levels: tasks

## Increment 1: the eight levels' rules in code (spec: [spec.md](./spec.md))

- [x] 1. Give the backgrounds unique base names (`<level>_day.png`, `<level>_night.*`) and point level 1 at
      its day background → check: `go test ./internal/core -run TestAssetBaseNamesUnique`
- [ ] 2. Decide how to keep the WASM small with ~110 MB of backgrounds (user) → check: decision in `decisions.md`
- [ ] 3. Each `LevelDef` lists its tiles; a test checks every level grid against it (and larva not next to grass in boulderBastion) → check: `go test ./internal/core -run TestLevelTiles`
- [ ] 4. Day and night sub-levels: a `LevelDef` per sub-level, deepForest night as the first new one, with a portal so the travel menu reaches it → check: `go test ./internal/core -run TestSwitchLevel`; web build screenshot
- [ ] 5. Update docs: state.md "Now", spec.md, current-state.md

## Later (not planned yet)

- [ ] Grids for the other 14 sub-levels (levels 2–8, day and night), each with a portal
- [ ] Enemy lists per level, once the user sets them (D3)
