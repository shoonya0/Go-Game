# Levels: tasks

## Increment 1: the eight levels' rules in code (spec: [spec.md](./spec.md))

- [x] 1. Give the backgrounds unique base names (`<level>_day.png`, `<level>_night.*`) and point level 1 at
      its day background → check: `go test ./internal/core -run TestAssetBaseNamesUnique`
- [x] 2. Decide how to keep the WASM small with ~110 MB of backgrounds → D4 (load per sub-level, full quality)
- [x] 3. Register the JPEG decoder explicitly (it only worked because the font library imports `image/jpeg`) → check: `TestBackgroundsDecode`
- [x] 4. Web build fetches the current sub-level's background instead of embedding all (D4): native keeps the embed; "Loading..." while it downloads, retried every 2 s on failure; server and Dockerfile serve the files → check: `game.wasm` 27.6 MB; web build screenshots `_docs/evidence/levels-bg-*.png`
- [x] 5. Each `LevelDef` lists its tiles (`Tiles`, `KeepApart` for larva/grass in boulderBastion); a test checks every level grid; deepForest's larva pit became water → check: `go test ./internal/core -run TestLevelTiles`
- [ ] 6. Day and night sub-levels in order (D1, D5): a `LevelDef` per sub-level with its next one; deepForest night as the first new grid, with 2+ portals → check: `go test ./internal/core -run TestSwitchLevel`
- [ ] 7. Saving the last unused portal teleports to the next sub-level's entry portal, once (D5) → check: `go test ./internal/core -run TestLastPortal`; web build
- [ ] 8. Update docs: state.md "Now", spec.md, current-state.md

## Later (not planned yet)

- [ ] Grids for the other 14 sub-levels (levels 2–8, day and night), each with a portal
- [ ] Enemy lists per level, once the user sets them (D3)
- [ ] Look at deepForest's former larva pit (now water) in game
- [ ] Build the Docker image once (`docker build .`) and check `/levels/deepForest/background/deepForest_day.png` is served: not run, Docker daemon was off
