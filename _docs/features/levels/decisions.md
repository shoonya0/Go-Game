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
