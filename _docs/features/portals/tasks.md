# Portals: tasks

## Increment 1: see a portal, start at it, save at it (spec: [spec.md](./spec.md))

- [ ] 1. Portal sprite: 14-frame loop from `portal.png` (7×2 grid, ~456×672 cells), scaled to
      about 2 tiles tall, drawn at every code-9 cell with its base on the cell floor
      → check: `TestBuildLevelPortals`, `TestPortalSheetFrames`, `TestPortalStandsOnItsCell`, manual look
- [x] 2. A new game starts at the first portal (D5), saved → check: `TestPlayerStartsAtFirstPortal`
- [ ] 3. Save with E, or a tap on the portal (touch), while overlapping it; touching alone no
      longer saves; a prompt shows at the portal and a short "saved" note after
      → check: `TestPortalSavesOnlyOnInteract`, `TestTapOnPortalSaves`, manual look
- [x] 3a. The saved portal is drawn golden (D6) → check: `TestOnlySavedPortalIsGolden`, manual look
- [x] 4. Death respawns at the saved portal → check: `TestRespawnAtSavedPortal`

## Increment 2: enemies come back

- [x] 5. `Boss` flag on enemy blueprints (no boss yet) → check: `TestBlueprintsValid`
- [x] 6. Death respawn resets every non-boss enemy of the level to its spawn; killed bosses
      stay dead → check: `TestRespawnRevivesEnemies`, `TestKilledBossStaysDead`

## Increment 3: travel

- [x] 7. Levels by name: the game holds the current level and can switch to another
      (platforms, portals, enemies, background) → check: `TestSwitchLevel`
- [ ] 8a. Portal looks (D7): inactive = one still frame, active = loop, checkpoint = gold;
      E/tap activates and makes it the checkpoint; active portals remembered by level + cell
      → check: `TestPortalActivatesOnInteract`, `TestPortalLooks`, manual look
- [ ] 8b. Hold E 0.5 s, or a long press on the portal, opens the travel menu (D8); up/down
      pick, E/Enter travel, Esc or a tap outside closes; the game pauses while it is open
      → check: `TestHoldEOpensTravelMenu`, `TestTravelMenuListsActivePortals`
- [ ] 8c. Menu draws a column of thumbnails (D9), captured when a portal is activated; the
      checkpoint framed gold → check: manual look (PC and touch)
- [ ] 8d. Picking a portal travels there (switching level if needed); death respawns in the
      checkpoint's level → check: `TestTravelMovesPlayer`, `TestRespawnInCheckpointLevel`
- [ ] 9. Entering a level revives its non-boss enemies → check: `TestEnteringLevelRevivesEnemies`
- [ ] 10. `verify.mjs --full`, web smoke check, `/code-review`, docs update, PR (ask before merge)

## Later (not planned yet)

- [ ] Keep saved portals across game restarts (file on Windows, browser storage on the web)
- [ ] A second level to travel to
- [ ] `portal.png` is 3.6 MB (was 5.6) but drawn 150 px tall: a 3× smaller copy would cut the web download
- [ ] `portal.png` has an AI sparkle mark in its bottom-right corner (inside the last frame's
      cell): crop it out of the asset
