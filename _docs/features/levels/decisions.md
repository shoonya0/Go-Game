# Levels: decisions

## D1. A day and a night sub-level per level, no switching (2026-10-26, user's call)

- **Decision:** each of the eight levels is two sub-levels, day and night, each with its own
  grid and its own background (`<level>_day.*` / `<level>_night.*`). Nothing switches
  between them over time. 8 levels → 16 sub-levels.
- **Why:** the user wants two places to explore per theme, not a clock.

## D2. Portals connect every sub-level (2026-10-26, user's call)

- **Decision:** the player roams freely between sub-levels by portal, like a connected graph:
  any used portal can be reached from any other. There is no fixed order.
- **How:** this already works. The travel menu lists every portal the player has used, in
  any level (`travel_menu.go`, from `p.ActivePortals`). Each sub-level needs at least one
  portal (code 9), and the player must reach it on foot once before it can be travelled to.
  How a new player first gets into a sub-level they've never been to is open
  (see [state.md](./state.md)).

## D3. Enemies can go in any level for now (2026-10-26, user's call)

- **Decision:** there is no enemy list per level yet. Any enemy marker can be placed in any
  sub-level. This replaces "enemies according to their type" in [spec.md](./spec.md) until
  the user sets per-level lists.

## D4. Full-quality backgrounds, loaded only for the sub-level played (2026-10-26, user's call)

- **Decision:** backgrounds are not shrunk or recompressed. Only the background of the
  sub-level the player is in gets loaded.
- **How:** decoding is already per level (`cmd/main.go`, loads when `Def.Background`
  changes). The download is the problem: `//go:embed levels` puts every background in
  `game.wasm`, so the web build stops embedding backgrounds and fetches the one it needs from
  the server when the player enters a sub-level (a few MB each). The native exe keeps
  embedding them. The fetch must not freeze the game loop (draw a loading state until it
  arrives), and the OS/network code goes in `internal/system`, not `internal/core`.
- **Rejected:** scaling down or lossy re-encoding (the user: it lowers the quality).

## D5. Finding every portal of a sub-level opens the next one (2026-10-26, user's call)

- **Decision:** when the player saves (E / tap) the **last unused** portal of a sub-level,
  they are teleported at once to the entry portal of the **next** sub-level. That portal
  becomes used (so it's in the travel menu) and the checkpoint. From there the player can
  travel back to any used portal.
- **Rules (my call, the user can veto):** it triggers only on the save that completes the
  set, once; arriving by teleport doesn't trigger it. So every sub-level needs **at least
  two portals**, or a sub-level entered by teleport would be a dead end. The next sub-level's
  entry portal is its first portal in reading order. The order is the [spec](./spec.md)
  table, day before night: deepForest day → deepForest night → aquaticAbyss day → … →
  rustyRefinery night, which is last and opens nothing.

## D6. Every level is walled in by an invisible border (2026-10-26, user's call)

- **Decision:** each level is closed on all four sides by walls that block like rock but are
  never drawn, so the player (and enemies) can't leave it.
- **How (my call):** the engine adds them, not the grid: `World` puts four solid walls of the
  `Border` tile type just outside the level's edges when it loads a level (`levelBorder`,
  `internal/core/world.go`), and drawing skips that type. No grid needs a rock frame, so a
  level whose tiles don't include rock (rustyRefinery) is walled in too.

## D7. The travel menu is a list by level, not thumbnails (2026-10-26, user's call)

- **Decision:** the travel menu is a wooden board like the main menu: each level's name
  ("Deep Forest"), under it each sub-level's name ("Deep Forest Day"), then a plank per
  active portal, "Portal 1", "Portal 2", ... The pick is the arrow plank, the checkpoint has
  a check at its left; a list taller than the screen scrolls.
- **How (my call):** a portal's number is its place among all its sub-level's portals in
  reading order (1 is the entry portal), so it stays the same however many are used. Levels
  appear in `LevelOrder`. Names come from `LevelDef.Title` and `Time`.
- **Why:** the user asked for it; names read better than small pictures, and no thumbnails
  need capturing (the web build skipped them while a background downloaded).
- **Replaces:** [portals D9](../portals/decisions.md) (rows were thumbnails).
