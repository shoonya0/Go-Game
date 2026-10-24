# Portals: state

## Now (updated 2026-10-24 at `3f3b728`)

- **Status:** 🟡 in progress; increment 1
- **Current task:** user looked in game: start at the portal confirmed (D5 ✓); floating portal fixed
  (`3f3b728`), to recheck in game. Task 1 stays open for the flickering rock (below); next 3 (save
  with E or a tap) in [tasks.md](./tasks.md)
- **Stopped at:** a portal's base and `StandPos` sit on the ground's solid top: `Portal.FloorDrop`
  (26.6 px under grass, sand, water) set by `BuildLevel`. Touching any portal still saves (until
  task 3), but at the cell's top-left, not `StandPos`
- **Open questions / blockers:** the arch's top-left stone is missing from the source art in
  frames 2–4, 6 and 9–11 (not a crop bug), so it flickers: ask the user to fix the asset, or play
  only the 7 frames that have it. D3 (E saves and opens the menu in one press) to confirm in play.
  `portal.png` is 5.6 MB but drawn 150 px tall (asset change, ask the user)
- **Don't forget:** the portal sheet's cells drift (a few px per frame; row 2 sits 33 px higher),
  so frames must be cropped to their art, not to the grid. `updatePortals` must save
  `cp.StandPos()`, not `cp.Pos` (task 3 or 4). Code 9 is non-solid. The art's base is a stone
  disc with a rounded bottom: if it still looks like it floats, sink it a few px

## Recent sessions (newest first, keep 5)

### 2026-10-24: user's look in game; floating portal fixed

- Changed: `Portal.FloorDrop` (world.go), test `TestPortalStandsOnTheGround`; probed the sheet (rock missing in 7 frames' art)
- Verified: new test failed first (base at 60, ground at 86.6), then `verify.mjs` 5/5
- Docs touched: state

### 2026-10-24: task 2, start at the first portal

- Changed: `Portal.StandPos` (world.go), `PlayerRuntime.StartAt` (player.go), `NewGame` calls it
- Verified: `TestPlayerStartsAtFirstPortal` failed first (build), then `verify.mjs` 5/5
- Docs touched: tasks, state

### 2026-10-24: task 1, portal sprite

- Changed: `Checkpoint` → `Portal` (8 files, one rename); `portal.go` (`PortalSprite`, `artFrames`), `Portal.Bounds` doorway, `cmd/main.go` draws it
- Verified: new tests failed first (build); `verify.mjs --full` 8/8; frame scan 0.17 s with the PNG decode
- Docs touched: state

### 2026-10-24: feature started

- Changed: docs only (`_docs/features/portals/`), registry row, Active work
- Verified: `node scripts/docs-check.mjs`
- Docs touched: README, spec, decisions, tasks, state
