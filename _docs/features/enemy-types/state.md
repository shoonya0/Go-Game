# Enemy types: state

## Now (updated 2026-10-24 at `f6a311a`)

- **Status:** 🟡 in progress; increments 1-5 merged except task 25 (PR #18 `2bc4469`, #19 `9124414`,
  #20 `c2c3b87`, #21 `f6a311a`, merged by the user); fresh `feature/enemy-types` from `main`
- **Current task:** paused, low priority (the user, 2026-10-24: portals first). None ready: task 25 (projectiles, ranged enemy) is blocked on art with a shooting clip; or the "Later" list in [tasks.md](./tasks.md)
- **Stopped at:** tasks 22-23 done; Flying Demon `^` added (row 25, col 14 of level 1, over the first floor).
  The user checked the demon in game (increment 4 done). Smasher `*` added (row 25, col 38): body 56×80,
  `HoverHeight` 160, `FallSpeed` 700 are first guesses. Task 25 waits for ranged art. Task 27 partly done:
  web smoke check passed, `code-reviewer` found no blockers (one flyer wall-pin edge case, fixed)
- **Traps:** behaviour tests use `testEnemyBlueprint` (enemy_test.go), not a game enemy: change the roster freely; flyer tests use `testFlyerBlueprint` (enemy_fly_test.go). Ask `bp.Has(clip)` for "can it?", not `bp.Clips[...]`; `Has` ignores fallbacks, `resolveClip` follows them. Player hits must go through `TakeHitFrom` (blockable); `TakeDamage` is unblockable.
  `TestBlueprintsValid` skips `EnemyGuard` in the every-state check (optional) and requires
  `GuardChance` in [0,1] with a guard clip. No strip pack has a `guard` clip. The mushroom
  art faces left (`FacesLeft`). The ledge probe acts only in `EnemyPatrol`: chasing enemies
  still run off edges. Tests that hit an enemy with a ground wave must put it on the floor
  with `standOnWaveFloor` (sets `onGround`). A new Attacker must implement `Ready` (a ranged one: its fire timer) and `Committed`.
  Detection is horizontal only (`DetectRange` on dx), also for flyers. A flyer's spawn
  cell is its home height. Death clips of both golems, Enemy3 and the demon end on a blank frame
- **Open questions / blockers:** guard rules while mid-swing or hurt (decisions D8, open). Keep or delete `tools/battler2strip` (its goblin mapping
  points at a deleted sheet). No strip pack for a ranged enemy (task 25). Hook noise: `docs-check.mjs:292`
  counts untracked assets as changed code (issue proposed, not filed)
- **Don't forget:** the remaining enemy folders are untracked, but `//go:embed enemy` embeds them
  anyway; commit each with its code. Rename the Orange Golem files `OrangeGolem_*` before
  using them (D4)

## Recent sessions (newest first, keep 5)

### 2026-10-24: increments 4-5 merged

- Changed: PR #21 merged by the user (`f6a311a`); remote branch gone, local `main` synced, `feature/enemy-types` restarted from it
- Verified: `verify.mjs --full` 8/8 before the push; branch was up to date with `main`
- Next: task 25 when ranged art exists, or a "Later" item (bosses, mushroom stun)

### 2026-10-24: task 27 checks (partial)

- Changed: flyer wall-pin fix in `flyingMover.Move` + `TestStrandedFlyerPatrolsWhereItIs` (review note)
- Verified: new test failed first (pinned at the wall); web build loads, no console errors, demon + F3 boxes draw; `code-reviewer`: no blockers; `verify.mjs` 5/5
- Next: PR (ask before merge); the user confirmed "Smasher working". Review gaps, not bugs: no test for hurt mid-dive

### 2026-10-24: task 26, dive attack + Smasher built

- Changed: `attack_loop`/`attack_end` clips, `swingState.phase`, `Attacker.Committed`; `diveAttack` (`enemy_dive.go`); Smasher `*` in level 1
- Verified: dive tests failed first (build); mutations (not committed, no MaxDrop) caught; it hits a player in the level; `verify.mjs --full` 8/8
- Next: user looks at the Smasher in game; then task 25 (needs ranged art) or 27

### 2026-10-24: task 24, Flying Demon built

- Changed: strips renamed `FlyingDemon_*`; `flyingDemonBlueprint`, `KindFlyingDemon`, `^` in level 1; level test checks flyers keep their height
- Verified: level test failed first (not placed); swing test lands every bite; it reaches a player under it in the level; `verify.mjs --full` 8/8
- Next: task 25; the user confirmed "Flying Demon working"

### 2026-10-24: task 23, ground wave skips flyers

- Changed: `GroundWave.update` skips enemies not `Grounded`; test helper `standOnWaveFloor` (sets `onGround`), used by the two wave-hit tests
- Verified: `TestGroundWaveMissesFlyers` failed first (flyer hit); without `onGround` in the helper the 2 old tests fail; `verify.mjs` 5/5
- Next: task 24, Flying Demon
