# Portals: spec

Asked for by the user on 2026-10-24 ("spawning player using portal.png": checkpoint,
passage to other levels, open-world enemy revival). Their answers are in
[decisions.md](./decisions.md).

## What the player sees

- A swirling stone portal stands at every code-`9` cell of a level grid
  (`assets/levels/deepForest/level_1.json`), its base on the cell's floor.
- A new game starts at the level's start portal, already saved.
- Standing at a portal, a prompt shows how to use it. Pressing **E** (PC), or tapping the
  portal on a touch screen while near it, **saves** it: it becomes the respawn point. Walking
  past a portal does nothing.
- Using a portal also opens a **travel menu** listing every portal the player has saved, in
  any level. Picking one takes the player there; closing the menu stays put.
- When the player **dies**, they respawn at the saved portal and every non-boss enemy of
  that level is back at its spawn point with full health.
- **Entering a level** through a portal brings back all its non-boss enemies too.
  Bosses that were killed stay dead.

## Acceptance criteria

1. Every code-9 cell draws the animated portal (14 frames, looping); no other cell does.
2. A new game places the player at the start portal, which is saved.
3. Touching a portal does not save it; E or a tap on it while overlapping does.
4. Death respawns the player at the last saved portal with full health and power.
5. Death respawn and entering a level reset every non-boss enemy of that level to its spawn
   (alive, full health, patrolling); a killed boss stays dead.
6. The travel menu lists exactly the saved portals and moves the player to the one picked,
   loading its level when it is another one.
7. Works on Windows and in the browser (keyboard and touch). Native and WASM builds pass.

## Out of scope for now

- Saving progress to disk or browser storage: saved portals last for one play session.
- A second level: travel is built for several levels, tried with level 1's portals.
