# Enemy types: state

## Now (updated 2026-10-23 at `a8077b3`)

- **Status:** 🟡 in progress
- **Current task:** 4. Per-frame `Hold` timing + draw-time motion (see [tasks.md](./tasks.md))
- **Stopped at:** task 4 code and tests done, `verify.mjs --full` 8/8, code-reviewer: no blocking findings,
  committed in `a8077b3`. Not ticked: it still needs the user to look at the succubus in a running
  game. Tuning to judge: chase lean 5° + bob 3px@3Hz, hurt shake 4px, attack frame 3 held ×2
  (values in `enemy_types.go`)
- **Open questions / blockers:** task 11: where the raw battler sheets live once baked
  (the user decides; they must leave the embed)
- **Don't forget:** assets changed on 2026-10-23. Removed: the desert set, `Bat with VFX/`,
  `slime/`. Skeleton is split into `Skeleton_White/` and `Skeleton_Yellow/`. Added:
  `Free Sprites/`, 17 RPG Maker MV side-view battlers (9×6 grid, 3 frames per motion, no walk
  cycle; 9 sheets have non-integer cell heights; facing varies per sheet). The golems still
  share `Golem_1_*.png` (D4); the user confirmed they are two different enemies. The enemy
  folders are untracked, but `//go:embed enemy` embeds them anyway; commit each with its code

## Recent sessions (newest first, keep 5)

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

### 2026-10-23: choose one sprite format + guard

- Changed: docs only. Decided: battlers are baked into strips with chained motions (D7), guard trait (D8)
- Verified: docs-check
- Docs touched: spec (roster, standard clips), decisions D3/D6/D7/D8, design, tasks rewritten (nothing ticked yet)
