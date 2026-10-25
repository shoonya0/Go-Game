# Main menu: decisions

## D1. Press Start 2P through text/v2 (2026-10-25, user)

- **Decision:** draw menu text with Ebitengine's `text/v2` and the Press Start 2P font that
  ships in the ebiten module (`examples/resources/fonts`). `text/v2` adds
  `github.com/go-text/typesetting` as a new indirect dependency (via `go get`).
- **Rejected:** ebiten's debug bitmap font (plain look); image buttons (art for every word).

## D2. A save holds portals and settings only (2026-10-25)

- **Decision:** active portals, checkpoint, settings. No position, health or enemies:
  Continue starts in the checkpoint like a respawn, matching "spawn at the last portal".
- **Why:** an Ori-style open world where the portal is the save point; fewer fields to keep
  compatible.

## D3. Save on change, not on quit (2026-10-25)

- **Decision:** write the save when the checkpoint or the active portals change (save at a
  portal, travel) and when a setting changes.
- **Why:** a browser tab can be closed without any event the game sees.

## D4. New Game asks before wiping a save (2026-10-25, my call)

- **Decision:** with a save present, New Game shows Yes / No first.
- **Why:** one wrong Enter would otherwise lose the save. Change it if it feels slow.

## D5. Volume only, no audio yet (2026-10-25, user)

- **Decision:** Settings stores a volume; no audio package until there are sound files.
- **Why:** `ebiten/audio` adds `oto` as a dependency, and nothing would play yet.

## D6. Quit only on desktop (2026-10-25)

- **Decision:** Quit returns `ebiten.Termination` on desktop and is hidden on the web.
- **Why:** a page can't close its own tab.

## D7. Esc opens a pause menu (2026-10-25, my call)

- **Decision:** Resume / Settings / Main menu, the game paused underneath.
- **Why:** without it the settings and the title screen can't be reached once playing.
  The user asked only for the main screen; drop this if unwanted.
