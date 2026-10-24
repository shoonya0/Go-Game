# Portals: state

## Now (updated 2026-10-24 at `f3cb97f`)

- **Status:** ⏳ planned; docs started, no code yet
- **Current task:** 1 (portal sprite at code-9 cells) in [tasks.md](./tasks.md)
- **Stopped at:** spec and decisions written from the user's answers; branch `feature/portals`
  made from `main` plus the enemy-types docs commit `f3cb97f`
- **Open questions / blockers:** D3 (E saves and opens the menu in one press) and D5 (start at
  the first 9, top-left in level 1) are my calls, to confirm in play
- **Don't forget:** `portal.png` is a 7×2 grid, but 3194 / 7 is not whole (456 px cells plus
  2 px). Assets are found by base name: `portal.png` is unique today. Code 9 is non-solid
  (`tileTypeFromCode` maps it to Empty)

## Recent sessions (newest first, keep 5)

### 2026-10-24: feature started

- Changed: docs only (`_docs/features/portals/`), registry row, Active work
- Verified: `node scripts/docs-check.mjs`
- Docs touched: README, spec, decisions, tasks, state
