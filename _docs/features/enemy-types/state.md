# Enemy types: state

## Now (updated 2026-10-24 at `3b98b41`)

- **Status:** 🟡 in progress; increments 1-2 merged (PR #18 `2bc4469`, PR #19 `9124414`);
  increment 3 (ground roster + guard): tasks 12-14 done; committed locally, not pushed
- **Current task:** none. Paused for a talk with the user about the enemy plan (see [tasks.md](./tasks.md)),
  on `feature/enemy-types` from `main` at `9124414`
- **Stopped at:** the user turned the goblin down: its battler sheet has no walk cycle,
  so it slid along in its standing pose. Reverted the blueprint, level marker (`12de28a`)
  and strips (`3b98b41`). The bake tool's feet fix (`1615825`) stays. Tasks 15 and 19
  (battlers) are on hold until the talk
- **Traps:** every battler sheet lacks walk frames, so a walking battler has the same problem.
  Player hits must go through `TakeHitFrom` (blockable); `TakeDamage` is unblockable.
  `TestBlueprintsValid` skips `EnemyGuard` in the every-state check (optional) and requires
  `GuardChance` in [0,1] with a guard clip. The tool keeps its own copy of the clip names; keep
  it in sync with `enemy_clips.go`. `TestBlueprintsValid` needs speeds > 0, so the static
  training dummy (task 20) will need it relaxed for a `static` mover. The mushroom art faces
  left (`FacesLeft`). The ledge probe acts only in `EnemyPatrol`: chasing enemies still run off
  edges. A new Attacker must implement `Ready` (a ranged one: its fire timer)
- **Open questions / blockers:** the user wants to talk about enemies next: how enemies
  without walk art move, or whether to drop the battlers. My call, not in D8: an enemy
  mid-swing can't guard (user not asked yet). Hook noise: `docs-check.mjs:292` counts
  untracked assets as changed code (issue proposed, not filed)
- **Don't forget:** assets changed on 2026-10-23. Removed: the desert set, `Bat with VFX/`,
  `slime/`. Skeleton is split into `Skeleton_White/` and `Skeleton_Yellow/`. Added:
  17 RPG Maker MV side-view battlers (now in `art-src/battlers/`; 9×6 grid, 3 frames per motion, no walk
  cycle; 9 sheets have non-integer cell heights; facing varies per sheet). The golems still
  share `Golem_1_*.png` (D4); the user confirmed they are two different enemies. The enemy
  folders are untracked, but `//go:embed enemy` embeds them anyway; commit each with its code

## Recent sessions (newest first, keep 5)

### 2026-10-24: goblin removed

- Changed: reverted the goblin (`12de28a`, `3b98b41`); the user turned it down because it had no walk animation
- Verified: `verify.mjs` 5/5 after the revert
- Next: talk with the user about the enemy plan

### 2026-10-24: task 15, goblin (built, look pending)

- Changed: `battler2strip` (`floorStrips`, bake all clips before writing); goblin strips re-baked and committed; `goblinBlueprint`, `KindGoblin`, `$` in level 1
- Verified: `TestFloorStripsMovesFeetToBottom` failed first; `verify.mjs --full` 8/8; feet measured at y=119/120; web screenshots with and without F3
- Docs touched: design (feet shift), tasks (15 note), state, current-state

### 2026-10-24: task 14, guard trait

- Changed: new `enemy_guard.go` (`TakeHitFrom`, `rollGuard`) + tests; `EnemyGuard`, `GuardChance`; swings and the ground wave call `TakeHitFrom`
- Verified: `verify.mjs --full` 8/8; mutations (front check, each caller's side) fail the tests; code-reviewer: no blocking findings (2 suggestions applied)
- Docs touched: design (guard notes), tasks (14 ticked), state, current-state

### 2026-10-24: task 13, two-character markers

- Changed: `level_enemies.go` (`quoteBareMarkers` quotes runs, `maxMarkerLen`); tests `TestTwoCharMarkerIsOneCell`, `TestEnemyMarkersUniqueAndBare` (renamed); marker comments
- Verified: new tests failed first (`"!""@"`), then pass; `TestLevelFileSpawnsEnemies` passes; `verify.mjs` 5/5
- Docs touched: decisions (D6 done), tasks (13 ticked), state, current-state

### 2026-10-24: task 12, ledge probe

- Changed: `enemy_move.go` (`groundMover.turnAtLedge`); new `TestGroundEnemyTurnsAtLedge`
- Verified: test failed first (right edge 721 > 720 at tick 34), then passes; `verify.mjs` 5/5
- Docs touched: tasks (12 ticked), state, current-state

