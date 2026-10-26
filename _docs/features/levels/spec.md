# Levels: spec

The user's rules for every level (2026-10-26). Grid codes for the tile types are in
`tileRowOrder` (`internal/core/world.go`).

## Rules for every level

1. A level uses **only its own** tiles and background, as listed below.
2. Every level is two sub-levels, **day** and **night**, each with its own grid and background
   (`assets/levels/<level>/background/<level>_day.*`, `_night.*`). Nothing switches between
   them ([D1](./decisions.md)).
3. Portals connect every sub-level to every other: free roaming, no fixed order ([D2](./decisions.md)).
4. Any enemy can go in any sub-level for now ([D3](./decisions.md)).
5. Tile amounts are how much of the level each tile makes up: most > medium > less > very less.

## Levels

| # | Level (folder in `assets/levels/`) | Tiles |
|---|---|---|
| 1 | `deepForest` | grass (most), water (medium), rock (less) |
| 2 | `aquaticAbyss` | water (most), rock (medium), grass (less) |
| 3 | `boulderBastion` | rock (most), grass (less), larva (very less, **never near grass**) |
| 4 | `cargoYardChaos` | wood (most), grass (medium), water (medium), rock (less) |
| 5 | `frostbitePeaks` | ice (most), water (medium), rock (less) |
| 6 | `mirageCanyon` | sand (most), rock (medium), larva (less) |
| 7 | `moltenCrater` | larva (most), rock (medium) |
| 8 | `rustyRefinery` | metal (most), wood (medium) |

## Acceptance

- A test fails when a level's grid holds a tile code outside its list above.
- A test fails when larva touches grass in boulderBastion.
- Each sub-level loads its own background (day or night); no two embedded assets share a base name.
- Every sub-level has at least one portal, so the travel menu can reach it.
