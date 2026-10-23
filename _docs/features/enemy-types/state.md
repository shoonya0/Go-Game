# Enemy types: state

## Now (updated 2026-10-23 at `19b6be7`)

- **Status:** 🟡 in progress
- **Current task:** 3. Standard clip names + fallbacks (see [tasks.md](./tasks.md))
- **Stopped at:** task 2 done: `EnemyClip` per state with its own sheet (`internal/core/enemy.go`),
  sheets cached by path, succubus walk fixed to 7 frames. Not yet looked at in a running game
- **Open questions / blockers:** task 11: where the raw battler sheets live once baked
  (the user decides; they must leave the embed)
- **Don't forget:** assets changed on 2026-10-23. Removed: the desert set, `Bat with VFX/`,
  `slime/`. Skeleton is split into `Skeleton_White/` and `Skeleton_Yellow/`. Added:
  `Free Sprites/`, 17 RPG Maker MV side-view battlers (9×6 grid, 3 frames per motion, no walk
  cycle; 9 sheets have non-integer cell heights; facing varies per sheet). The golems still
  share `Golem_1_*.png` (D4); the user confirmed they are two different enemies. The enemy
  folders are untracked, but `//go:embed enemy` embeds them anyway; commit each with its code

## Recent sessions (newest first, keep 5)

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

### 2026-10-23: check the Free Sprites battler format

- Changed: nothing in code; scanned the sheets with a scratch Go tool (frame gaps and grid overlay)
- Verified: all 17 `Free Sprites/` sheets are 9 cols × 6 rows; goblin and rat match `imagePattern.png`
- Docs touched: state.md (open question: how to handle the two sprite formats)

### 2026-10-23: plan the feature

- Changed: docs only (`_docs/features/enemy-types/`, registry, Active work, rules pointer)
- Verified: `go build ./...` ok with the untracked folders; asset base-name scan found the golem collision
- Docs touched: README, spec, design, decisions, tasks, state

## Known issues

- Duplicate asset base names `Golem_1_*.png` (`assets/enemy/White_Golem/`, `assets/enemy/Orange_Golem/`), fixed by task 10
