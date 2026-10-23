# Enemy types: state

## Now (updated 2026-10-23 at `91fd78b`)

- **Status:** 🟡 in progress
- **Current task:** 17 + 18 waiting for a look in a running game, then 5. Extract `Mover` +
  `groundMover` (see [tasks.md](./tasks.md))
- **Stopped at:** task 4 done. The user skipped the succubus look, so the skeleton (White `@`,
  Yellow `#`) and mushroom (`&`) were added early and placed in level 1 to judge hold timing and
  draw-time motion on real strip art. Tuning lives on the clip lines in `enemy_types.go`
- **Traps:** the mushroom art faces left (`FacesLeft`). No ledge probe yet (task 12): patrols are
  kept at 2.5 tiles so they stop short of pits, but a chase can still run off a ledge
- **Open questions / blockers:** task 11: where the raw battler sheets live once baked
  (the user decides; they must leave the embed)
- **Don't forget:** assets changed on 2026-10-23. Removed: the desert set, `Bat with VFX/`,
  `slime/`. Skeleton is split into `Skeleton_White/` and `Skeleton_Yellow/`. Added:
  `Free Sprites/`, 17 RPG Maker MV side-view battlers (9×6 grid, 3 frames per motion, no walk
  cycle; 9 sheets have non-integer cell heights; facing varies per sheet). The golems still
  share `Golem_1_*.png` (D4); the user confirmed they are two different enemies. The enemy
  folders are untracked, but `//go:embed enemy` embeds them anyway; commit each with its code

## Recent sessions (newest first, keep 5)

### 2026-10-23: skeleton + mushroom pulled ahead (tasks 17, 18)

- Changed: `enemy_types.go` (skeleton builder, mushroom), `enemy.go` (kinds, `FacesLeft`), `enemy_anim.go` (`drawFlipX`), `level_1.json`, assets `Skeleton_*/`, `Mushroom/`
- Verified: `verify.mjs --full` 8/8 after each; new `TestSkeletonVariantsShareBlueprint`, `TestLeftFacingArtMirrorsTheOtherWay`; `TestEnemiesHaveAtLeastTwoWholeBars` caught 30 HP
- Docs touched: tasks (4 ticked, notes on 17/18), spec (roster checked), state

### 2026-10-23: task 4, Hold timing + draw-time motion

- Changed: `enemy.go` (Hold, Motion, clipTime), `enemy_anim.go`, `enemy_clips.go`, new `enemy_motion.go`, `animation.go` (`drawSpriteFramePosed`), succubus tuning
- Verified: `verify.mjs --full` 8/8; new `TestClipHoldTiming`, `TestMotionPose`, `TestEnemyClipTimeRestarts`, `TestClipDurationCountsHolds`
- Docs touched: design (clips section trimmed), state

### 2026-10-23: task 3, standard clip names + fallbacks

- Changed: new `enemy_clips.go`; `enemy.go`, `enemy_anim.go`, `enemy_health.go`, `enemy_types.go` use clip names; tests switched to `animFor`
- Verified: `verify.mjs --full` 8/8; new `TestEnemyMissingClipFallsBack`, `TestSuccubusChasesWithAlertClip`, `TestEveryEnemyStateHasAClip`
- Docs touched: spec (fallback chain), tasks (3 ticked), state

### 2026-10-23: task 2, clips with their own sheet

- Changed: `enemy.go`, `enemy_anim.go`, `enemy_types.go`, new `enemy_clips_test.go`; test call sites renamed
- Verified: `verify.mjs --full` 8/8; `TestEnemyClipsFitTheirSheets` failed on the 8-frame walk, passes at 7
- Docs touched: tasks (2 ticked), design (EnemyClip done), state

### 2026-10-23: task 1, succubus characterisation tests

- Changed: `internal/core/enemy_behavior_test.go` (7 tests)
- Verified: `node scripts/verify.mjs` pass; mutation (one-shot → idle) fails `TestEnemyOneShotClipsReturnToPatrol`
- Docs touched: tasks (1 ticked, patrol-frame note on 2), state
