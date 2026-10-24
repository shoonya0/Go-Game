# Portals: state

## Now (updated 2026-10-24 at `1c0f72e`)

- **Status:** 🟡 in progress; increment 1
- **Current task:** 2 (start at the first portal) done; 1 still waits for the user's look in game;
  next 3 (save with E or a tap) in [tasks.md](./tasks.md)
- **Stopped at:** `PlayerRuntime.StartAt` puts a new game in level 1's first portal (cell 33,3,
  ground below) via `Portal.StandPos`, saved as `RespawnPos`. Touching any portal still saves
  (until task 3), but at the cell's top-left, not `StandPos`
- **Open questions / blockers:** D3 (E saves and opens the menu in one press) and D5 (start at
  the first 9) are my calls, to confirm in play. `portal.png` is 5.6 MB but drawn 150 px tall:
  a 3× smaller copy would cut the web download (asset change, ask the user)
- **Don't forget:** the portal sheet's cells drift (a few px per frame; row 2 sits 33 px higher),
  so frames must be cropped to their art, not to the grid. `updatePortals` must save
  `cp.StandPos()`, not `cp.Pos` (task 3 or 4). Code 9 is non-solid

## Recent sessions (newest first, keep 5)

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
