# Main menu: spec

Asked by the user 2026-10-25. Levels (a second level) are a separate, later feature.

## Main screen

- The game opens on a title screen with, top to bottom: **New Game**, **Continue**,
  **Settings**, **Quit**. Text uses the Press Start 2P pixel font ([D1](./decisions.md)).
- Picked with W/S or the arrows and Enter/E; a tap or mouse click on an item picks it too.
- **New Game** wipes the save and starts in the start portal, as a new game does today. When
  a save exists it first asks to confirm (Yes / No) ([D4](./decisions.md)).
- **Continue** starts at the saved checkpoint portal (in its level), standing in it, with
  full health and power, and every portal used before still active. Drawn grayed and
  skipped by the pick when there is no save.
- **Settings** opens the settings screen. **Quit** closes the game on desktop; on the web it
  isn't shown ([D6](./decisions.md)).

## Save

- Holds the active portals, the checkpoint and the settings ([D2](./decisions.md)).
- Written when a portal is saved, after a travel and when a setting changes ([D3](./decisions.md)).
- Windows: a JSON file in the user's config folder. Web: browser `localStorage`.
- A save that can't be read (corrupt, unknown level or portal) counts as no save; it never
  crashes the game.

## Settings

- **Sound volume**, 0 to 10 in steps of 1, changed with A/D or the arrows (or taps on
  − / +), saved. The game has no sound yet, so it plays nothing until audio is added ([D5](./decisions.md)).
- Esc or **Back** returns to the screen that opened it.

## In game

- Esc (or the touch menu button) pauses and opens: **Resume**, **Settings**, **Main menu**
  ([D7](./decisions.md)). Esc again resumes.

## Done when

- Both builds compile, tests cover the menu navigation, save round trip and Continue placing
  the player; the web build was checked in a browser (title screen, Continue after reload).
