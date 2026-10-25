# Portals: decisions

## D1. Code 9 is the portal (2026-10-24, user)

- **Decision:** the existing checkpoint code 9 marks a portal; the invisible checkpoints
  become portals.
- **Why:** `0` is the empty cell in every grid, so it can't also mean portal. Reusing 9
  keeps one concept (a portal is the checkpoint) and the two 9s already in level 1.

## D2. Save by interacting, not by touching (2026-10-24, user)

- **Decision:** E on PC; on touch screens, tap the portal while near it. Saving sets the
  respawn point and does **not** revive enemies.
- **Why:** passing by shouldn't move the respawn point.

## D3. Travel is a menu of saved portals (2026-10-24, user)

- **Decision:** using a portal opens a list of every saved portal, in any level (fast travel).
  E at a portal saves it and opens the menu in one press (my call; change it if the double
  action feels wrong in play).
- **Rejected:** "next level in order" (needs a level 2 first) and per-portal targets in the
  level file (more setup per level).

## D4. Non-boss enemies revive on death and on entering a level (2026-10-24, user)

- **Decision:** respawning after death and travelling into a level reset its non-boss enemies
  to their spawns. Bosses never revive; a blueprint says whether it is one.
- **Why:** open-world loop: cleared areas fill up again, bosses are one-time.

## D5. The start portal is the first code 9 in reading order (2026-10-24, my call; user confirmed in play 2026-10-24)

- **Decision:** a new game starts at the level's first 9 (top row first, then left to right),
  already saved, instead of the hard-coded (100, 100).
- **Why:** the user asked to spawn the player with a portal; moving a 9 in the grid moves the
  start. Change it if level 1's top-left portal is the wrong start.

## D6. The saved portal is golden (2026-10-25, user)

- **Decision:** the portal holding the respawn point is drawn gold; the others keep the blue art.
- **Why:** the user wants to see which portal is their checkpoint. Done with a color matrix
  (gray, then scaled toward gold), not a second sheet: plain RGB scaling can't turn blue art
  gold, and a tinted copy would double the 3194×1344 texture.
