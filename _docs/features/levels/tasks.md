# Levels: tasks

## Increment 1: the eight levels' rules in code (spec: [spec.md](./spec.md))

- [x] 1. Give the backgrounds unique base names (`<level>_day.png`, `<level>_night.*`) and point level 1 at
      its day background → check: `go test ./internal/core -run TestAssetBaseNamesUnique`
- [x] 2. Decide how to keep the WASM small with ~110 MB of backgrounds → D4 (load per sub-level, full quality)
- [x] 3. Register the JPEG decoder explicitly (it only worked because the font library imports `image/jpeg`) → check: `TestBackgroundsDecode`
- [x] 4. Web build fetches the current sub-level's background instead of embedding all (D4): native keeps the embed; "Loading..." while it downloads, retried every 2 s on failure; server and Dockerfile serve the files → check: `game.wasm` 27.6 MB; web build screenshots `_docs/evidence/levels-bg-*.png`
- [x] 5. Each `LevelDef` lists its tiles (`Tiles`, `KeepApart` for larva/grass in boulderBastion); a test checks every level grid; deepForest's larva pit became water → check: `go test ./internal/core -run TestLevelTiles`
- [x] 6. Day and night sub-levels in order (D1, D5): `LevelOrder`; deepForest night grid (`deepForest_night.json`, 3 portals, climbs ≤ 3 tiles) → check: `TestLevelPortals`, `TestLevelTiles`, `TestLevelEnemiesKeepTheirGround`; web screenshot `_docs/evidence/levels-night-start.png`
- [x] 7. Saving the last unused portal teleports to the next sub-level's entry portal, once (D5) → check: `go test ./internal/core -run TestLastPortal`; in game: not yet (the user plays it)
- [x] 9. Invisible border around every level (D6) → check: `go test ./internal/core -run TestLevelBorderKeepsPlayerIn`
- [x] 10. Travel menu lists portals by level and sub-level instead of thumbnails (D7) → check: `TestTravelMenuGroupsByLevel`; web build screenshot `_docs/evidence/levels-travel-list.png`
- [x] 8. Update docs: state.md "Now", spec.md, current-state.md → check: `node scripts/docs-check.mjs`

## Later (not planned yet)

- [x] aquaticAbyss day grid (`aquaticAbyss_day.json`, 3 portals, draft) → check: `TestLevelTiles`, `TestLevelPortals`, `TestLevelBorderKeepsPlayerIn`, `TestLevelEnemiesKeepTheirGround`; web screenshot `_docs/evidence/levels-aquatic-day.png`
- [ ] Play aquaticAbyss day: is the layout right, is it fun
- [ ] Grids for the other 13 sub-levels (aquaticAbyss night, levels 3–8 day and night), each with at least 2 portals
- [ ] A test that every portal of a sub-level can reach the others within the jump rules (a script did it for aquaticAbyss day)
- [ ] Enemy lists per level, once the user sets them (D3)
- [ ] Play deepForest night: are the climbs right (≤ 3 tiles), is it fun; does saving day's last portal land in it
- [ ] Look at deepForest's former larva pit (now water) in game
- [ ] Build the Docker image once (`docker build .`) and check `/levels/deepForest/background/deepForest_day.png` is served: not run, Docker daemon was off
