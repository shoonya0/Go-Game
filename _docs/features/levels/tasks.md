# Levels: tasks

## Increment 1: the eight levels' rules in code (spec: [spec.md](./spec.md))

- [ ] 1. Give the backgrounds unique base names (e.g. `deepForest_day.png`) and point level 1 at
      its day background → check: `go test ./internal/core -run TestAssetBaseNamesUnique`
- [ ] 2. Decide how to keep the WASM small with ~110 MB of backgrounds (user) → check: decision in `decisions.md`
- [ ] 3. Each `LevelDef` lists its tiles; a test checks every level grid against it (and larva not next to grass in boulderBastion) → check: `go test ./internal/core -run TestLevelTiles`
- [ ] 4. Each level's enemy list (after the user picks them); a test checks the grid's markers → check: `go test ./internal/core -run TestLevelEnemies`
- [ ] 5. Day / night background per level (how it switches: user) → check: web build screenshot
- [ ] 6. Update docs: state.md "Now", spec.md, current-state.md

## Later (not planned yet)

- [ ] Grids for levels 2–8, each with a portal pair to reach it
