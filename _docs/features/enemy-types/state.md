# Enemy types: state

## Now (updated 2026-10-24 at `1484faa`)

- **Status:** 🟡 in progress; increments 1-3 merged (PR #18 `2bc4469`, PR #19 `9124414`,
  PR #20 `c2c3b87`); increment 4 started on `feature/enemy-types` (local, not pushed)
- **Current task:** task 23 (ground wave skips flyers) in [tasks.md](./tasks.md)
- **Stopped at:** task 22 done: `flyingMover` (`enemy_fly.go`) with 5 tests. No game enemy flies yet (task 24)
- **Traps:** behaviour tests use `testEnemyBlueprint` (enemy_test.go), not a game enemy: change the roster freely; flyer tests use `testFlyerBlueprint` (enemy_fly_test.go). Ask `bp.Has(clip)` for "can it?", not `bp.Clips[...]`; `Has` ignores fallbacks, `resolveClip` follows them. Player hits must go through `TakeHitFrom` (blockable); `TakeDamage` is unblockable.
  `TestBlueprintsValid` skips `EnemyGuard` in the every-state check (optional) and requires
  `GuardChance` in [0,1] with a guard clip. No strip pack has a `guard` clip. The mushroom
  art faces left (`FacesLeft`). The ledge probe acts only in `EnemyPatrol`: chasing enemies
  still run off edges. A new Attacker must implement `Ready` (a ranged one: its fire timer).
  Detection is horizontal only (`DetectRange` on dx), also for flyers. A flyer's spawn
  cell is its home height. Death clips of both golems, Enemy3 and the demon end on a blank frame
- **Open questions / blockers:** guard rules while mid-swing or hurt (decisions D8, open). Keep or delete `tools/battler2strip` (its goblin mapping
  points at a deleted sheet). No strip pack for a ranged enemy (task 25). Hook noise: `docs-check.mjs:292`
  counts untracked assets as changed code (issue proposed, not filed)
- **Don't forget:** the enemy folders are untracked, but `//go:embed enemy` embeds them
  anyway; commit each with its code. Rename the Orange Golem files `OrangeGolem_*` before
  using them (D4)

## Recent sessions (newest first, keep 5)

### 2026-10-24: task 22, flyingMover

- Changed: `spawnY` + `dy` in the move intent, `turnAtPathEnd`; new `flyingMover` (leash box, `HoverHeight`, dead flyers fall)
- Verified: tests failed first (build); mutations (gravity on, leash off) caught; `verify.mjs` 5/5; 125 tests pass
- Next: task 23, ground wave skips flyers

### 2026-10-24: increment 3 merged

- Changed: PR #20 merged by the user (`c2c3b87`); old branch deleted, `feature/enemy-types` restarted from `main`
- Verified: CI on PR #20 and on `main` green (build, verify ubuntu + windows); merged tree = branch tip `5671d66`
- Next: increment 4, task 22

### 2026-10-24: task 21, increment 3 closed

- Changed: ledge-freeze fix in `turnAtLedge` + `TestStrandedEnemyPatrolsWhereItIs`; D8 open guard questions; tasks 21 ticked
- Verified: new test failed first (frozen at x=674); `verify.mjs --full` 8/8; 120 tests; web build loads, no console errors
- Next: push, PR, user approves merge

### 2026-10-24: succubus removed

- Changed: tests moved to `testEnemyBlueprint`; succubus blueprint, `KindSuccubus`, `!` and `succubus.png` removed; comments
- Verified: same 119 tests pass; each of the 4 commits tested alone in a clean worktree; `verify.mjs --full` 8/8
- Next: user's look at the White Golem, then task 21

### 2026-10-24: task 16, White Golem

- Changed: strips renamed `WhiteGolem_*`; `whiteGolemBlueprint`, `KindWhiteGolem`, `%` in level 1; tests `TestAssetBaseNamesUnique`, `TestLevelEnemiesKeepTheirGround`
- Verified: both tests failed first / on mutation (missing marker, spawn over the pit); `verify.mjs --full` 8/8; 119 tests pass
- Next: task 21, then close increment 3
