# Portals: tasks

## Increment 1: see a portal, start at it, save at it (spec: [spec.md](./spec.md))

- [ ] 1. Portal sprite: 14-frame loop from `portal.png` (7×2 grid, ~456×672 cells), scaled to
      about 2 tiles tall, drawn at every code-9 cell with its base on the cell floor
      → check: `TestBuildLevelPortals`, `TestPortalSheetFrames`, `TestPortalStandsOnItsCell`, manual look
- [ ] 2. A new game starts at the first portal (D5), saved → check: `TestPlayerStartsAtFirstPortal`
- [ ] 3. Save with E, or a tap on the portal (touch), while overlapping it; touching alone no
      longer saves; a prompt shows at the portal and a short "saved" note after
      → check: `TestPortalSavesOnlyOnInteract`, `TestTapOnPortalSaves`, manual look
- [ ] 4. Death respawns at the saved portal → check: `TestRespawnAtSavedPortal`

## Increment 2: enemies come back

- [ ] 5. `Boss` flag on enemy blueprints (no boss yet) → check: `TestBlueprintsValid`
- [ ] 6. Death respawn resets every non-boss enemy of the level to its spawn; killed bosses
      stay dead → check: `TestRespawnRevivesEnemies`, `TestKilledBossStaysDead`

## Increment 3: travel

- [ ] 7. Levels by name: the game holds the current level and can switch to another
      (platforms, portals, enemies, background) → check: `TestSwitchLevel`
- [ ] 8. Travel menu: using a portal lists the saved portals (all levels); E/Esc or a tap
      outside closes it; picking one travels there → check: `TestTravelMenuListsSavedPortals`,
      `TestTravelMovesPlayer`, manual look (PC and touch)
- [ ] 9. Entering a level revives its non-boss enemies → check: `TestEnteringLevelRevivesEnemies`
- [ ] 10. `verify.mjs --full`, web smoke check, `/code-review`, docs update, PR (ask before merge)

## Later (not planned yet)

- [ ] Keep saved portals across game restarts (file on Windows, browser storage on the web)
- [ ] A second level to travel to
- [ ] `portal.png` is 5.6 MB but drawn 150 px tall: a 3× smaller copy would cut the web download
- [ ] `portal.png` has an AI sparkle mark in its bottom-right corner (inside the last frame's
      cell): crop it out of the asset
