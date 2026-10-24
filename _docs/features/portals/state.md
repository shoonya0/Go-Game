# Portals: state

## Now (updated 2026-10-24 at `def2257`)

- **Status:** 🟡 in progress; increment 1
- **Current task:** 1 (portal sprite) built, waiting for the user's look in game; next 2 (start at the
  first portal) in [tasks.md](./tasks.md)
- **Stopped at:** `Checkpoint` renamed `Portal` (`7fcf4e2`); `portal.go` draws the 14-frame loop at every
  code 9, cropped per frame by `artFrames`; a portal's area is its 1.5×2.5-tile doorway (F3 shows it).
  Touching it still saves (until task 3)
- **Open questions / blockers:** D3 (E saves and opens the menu in one press) and D5 (start at
  the first 9, top-left in level 1) are my calls, to confirm in play. `portal.png` is 5.6 MB but
  drawn 150 px tall: a 3× smaller copy would cut the web download (asset change, ask the user)
- **Don't forget:** the portal sheet's cells drift (a few px per frame; row 2 sits 33 px higher),
  so frames must be cropped to their art, not to the grid. `respawn` puts the player's top-left at
  `Portal.Pos` (the cell's top-left), which task 4 must fix. Code 9 is non-solid

## Recent sessions (newest first, keep 5)

### 2026-10-24: task 1, portal sprite

- Changed: `Checkpoint` → `Portal` (8 files, one rename); `portal.go` (`PortalSprite`, `artFrames`), `Portal.Bounds` doorway, `cmd/main.go` draws it
- Verified: new tests failed first (build); `verify.mjs --full` 8/8; frame scan 0.17 s with the PNG decode
- Docs touched: state

### 2026-10-24: feature started

- Changed: docs only (`_docs/features/portals/`), registry row, Active work
- Verified: `node scripts/docs-check.mjs`
- Docs touched: README, spec, decisions, tasks, state
