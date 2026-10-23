# Enemy types: state

## Now (updated 2026-10-23 at `6c45e04`)

- **Status:** 🟡 in progress; increment 1 done (tasks 1-8)
- **Current task:** 9. `tools/battler2strip` (see [tasks.md](./tasks.md)). Before it: push
  `feature/enemy-types` + PR for increment 1 (ask the user before merging: `main` deploys)
- **Stopped at:** task 8 done: full verify 8/8, code review, design.md trimmed. Enemies now
  idle between swings until the cooldown is ready (`Attacker.Ready`); hit timing and the
  idle pause haven't been looked at in a running game yet
- **Traps:** `TestBlueprintsValid` needs speeds > 0, so the static training dummy (task 20)
  will need it relaxed for a `static` mover. The mushroom art faces left (`FacesLeft`). No ledge probe yet (task 12).
  A new Attacker must implement `Ready` (a ranged one: its fire timer)
- **Open questions / blockers:** task 11: where the raw battler sheets live once baked
  (the user decides; they must leave the embed)
- **Don't forget:** assets changed on 2026-10-23. Removed: the desert set, `Bat with VFX/`,
  `slime/`. Skeleton is split into `Skeleton_White/` and `Skeleton_Yellow/`. Added:
  `Free Sprites/`, 17 RPG Maker MV side-view battlers (9×6 grid, 3 frames per motion, no walk
  cycle; 9 sheets have non-integer cell heights; facing varies per sheet). The golems still
  share `Golem_1_*.png` (D4); the user confirmed they are two different enemies. The enemy
  folders are untracked, but `//go:embed enemy` embeds them anyway; commit each with its code

## Recent sessions (newest first, keep 5)

### 2026-10-23: task 8, close increment 1

- Changed: `enemy_attack.go` (`Attacker.Ready`, `hitTime`), `enemy_ai.go` (idle until ready), `enemy_clips.go` (`timeTo`); new `TestMeleeEverySwingLands`
- Verified: `verify.mjs --full` 8/8 before and after; the test failed for all 4 species before the fix; code-reviewer: one blocking finding (whiffed swings), fixed
- Docs touched: design (trimmed to what the code doesn't show), tasks (8 ticked), state, current-state

### 2026-10-23: task 7, blueprint validation

- Changed: `internal/core/enemy_test.go` (`TestBlueprintsValid`)
- Verified: `verify.mjs --full` 8/8; mutation (mushroom without Attacker) fails it
- Docs touched: tasks (7 ticked), state

### 2026-10-23: task 6, Attacker + meleeAttack

- Changed: new `enemy_attack.go` + test; `enemy.go` (`HitFrame`, `Attacker`, `swing`), `enemy_ai.go`, `enemy_clips.go` (`clipFor` picks attack_2), `enemy_types.go`
- Verified: `verify.mjs --full` 8/8; reviewer measured missed swings when cooldown > clip (mushroom −40%), fixed test-first (`TestMeleeSwingLandsAfterCooldown`)
- Docs touched: design (Attacker done), tasks (6 ticked), state

### 2026-10-23: task 5, Mover + groundMover

- Changed: new `enemy_move.go`; `enemy_ai.go` steers the Mover; `enemy.go` (`Mover` field), `enemy_types.go`
- Verified: `verify.mjs --full` 8/8; enemy tests 19/19 with no test file changed; code-reviewer: no blocking findings
- Docs touched: design (Mover done, Steer/Move split), tasks (5, 17, 18 ticked), state

### 2026-10-23: skeleton + mushroom pulled ahead (tasks 17, 18)

- Changed: `enemy_types.go` (skeleton builder, mushroom), `enemy.go` (kinds, `FacesLeft`), `enemy_anim.go` (`drawFlipX`), `level_1.json`, assets `Skeleton_*/`, `Mushroom/`
- Verified: `verify.mjs --full` 8/8 after each; new `TestSkeletonVariantsShareBlueprint`, `TestLeftFacingArtMirrorsTheOtherWay`; `TestEnemiesHaveAtLeastTwoWholeBars` caught 30 HP
- Docs touched: tasks (4 ticked, notes on 17/18), spec (roster checked), state
