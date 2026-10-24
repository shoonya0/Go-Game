# Enemy types: state

## Now (updated 2026-10-24 at `0c4c596`)

- **Status:** 🟡 in progress; increments 1-3 merged (PR #18 `2bc4469`, PR #19 `9124414`,
  PR #20 `c2c3b87`); increment 4 started on `feature/enemy-types` (local, not pushed)
- **Current task:** task 27 (PR for increments 4-5) in [tasks.md](./tasks.md); task 25 stays open, blocked on ranged art
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

### 2026-10-24: task 22, flyingMover

- Changed: `spawnY` + `dy` in the move intent, `turnAtPathEnd`; new `flyingMover` (leash box, `HoverHeight`, dead flyers fall)
- Verified: tests failed first (build); mutations (gravity on, leash off) caught; `verify.mjs` 5/5; 125 tests pass
- Next: task 23, ground wave skips flyers
