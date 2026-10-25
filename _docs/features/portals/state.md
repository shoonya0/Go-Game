# Portals: state

## Now (updated 2026-10-25 at `00454fe`)

- **Status:** 🟡 in progress; increment 1
- **Current task:** the saved portal is golden (D6, `00454fe`); the user fixed the blinking stone in
  `portal.png` (`09eb174`). Waiting on the user's look in game (gold, ground, stone); then tick
  task 1 and start 3 (save with E or a tap) in [tasks.md](./tasks.md)
- **Stopped at:** `Portal.SavedBy(respawn)` picks the golden one; touching a portal now saves its
  `StandPos` (as the start does), so death already respawns the player standing in it; task 4 still
  needs `TestRespawnAtSavedPortal`. Loop plays at 5 fps (user's call)
- **Open questions / blockers:** D3 (E saves and opens the menu in one press) to confirm in play.
  `portal.png` is 3.6 MB but drawn 150 px tall (asset change, ask the user).
  `_asset_backups/` (untracked) is the user's own folder: never commit it
- **Don't forget:** the portal sheet's cells drift (a few px per frame; row 2 sits 33 px higher),
  so frames must be cropped to their art, not to the grid. Code 9 is non-solid. The art's base is
  a stone disc with a rounded bottom: if it still looks like it floats, sink it a few px

## Recent sessions (newest first, keep 5)

### 2026-10-25: golden saved portal

- Changed: `savedPortalTint` + `Draw(…, respawn, …)` (portal.go), `Portal.SavedBy`, `updatePortals` saves `StandPos`; user's `portal.png` fix committed
- Verified: `TestOnlySavedPortalIsGolden` failed first (build); `verify.mjs --full` 8/8; gold previewed offline
- Docs touched: decisions (D6), tasks, state

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
