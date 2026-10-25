# Portals: spec

Asked for by the user on 2026-10-24 ("spawning player using portal.png": checkpoint,
passage to other levels, open-world enemy revival). Their answers are in
[decisions.md](./decisions.md).

## What the player sees

- A swirling stone portal stands at every code-`9` cell of a level grid
  (`assets/levels/deepForest/level_1.json`), its base on the cell's floor.
- A new game starts at the level's start portal, already saved.
- A portal the player never used stands still (one frame). Standing at a portal, a prompt
  shows how to use it. Pressing **E** (PC), or tapping the portal on a touch screen while
  near it, **activates** it (it swirls, in its own colors) and makes it the **checkpoint**
  (gold), the respawn point. Walking past a portal does nothing.
- **Holding E** (a long press on touch) for half a second opens a **travel menu**: a column
  of thumbnails, one per active portal in any level. Picking one takes the player there
  and makes it the checkpoint (D10); closing the menu stays put.
- When the player **dies**, they respawn at the saved portal and every non-boss enemy of
  that level is back at its spawn point with full health.
- **Entering a level** through a portal, and any teleport inside a level (D11), brings back
  all its non-boss enemies too.
  Bosses that were killed stay dead.

## Acceptance criteria

1. Every code-9 cell draws the animated portal (14 frames, looping); no other cell does.
2. A new game places the player at the start portal, which is saved.
3. Touching a portal does not save it; E or a tap on it while overlapping activates it and
   makes it the checkpoint. Inactive portals draw one still frame, active ones the loop in
   their own colors, the checkpoint the loop in gold.
4. Death respawns the player at the last saved portal with full health and power.
5. Death respawn, entering a level and any teleport reset every non-boss enemy of that level to its spawn
   (alive, full health, patrolling); a killed boss stays dead.
6. Holding E (or a long press) 0.5 s at a portal opens the travel menu. It lists exactly
   the active portals, as thumbnails in a column, and moves the player to the one picked,
   loading its level when it is another one; the picked portal becomes the checkpoint.
7. Works on Windows and in the browser (keyboard and touch). Native and WASM builds pass.

## Out of scope for now

- Saving progress to disk or browser storage: saved portals last for one play session.
- A second level: travel is built for several levels, tried with level 1's portals.
