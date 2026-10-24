# Enemy types: state

## Now (updated 2026-10-24 at `3edfa17`)

- **Status:** 🟡 in progress; increments 1-2 merged (PR #18 `2bc4469`, PR #19 `9124414`);
  increment 3 (ground roster + guard) done: tasks 12-14, 16, 21 done, 15/19/20 dropped; committed locally, not pushed
- **Current task:** close increment 3: push `feature/enemy-types`, open the PR, ask the user before merging;
  then increment 4 (task 22, `flyingMover`) in [tasks.md](./tasks.md)
- **Stopped at:** task 21 done: web smoke check clean (32.6 MB WASM, no console errors),
  `verify.mjs --full` 8/8, 120 tests. `code-reviewer` found an enemy frozen at a ledge
  outside its patrol path (task 12 bug); fixed with a test (`3edfa17`)
- **Traps:** behaviour tests use `testEnemyBlueprint` (enemy_test.go), not a game enemy: change the roster freely. ask `bp.Has(clip)` for "can it?", not `bp.Clips[...]`; `Has` ignores fallbacks, `resolveClip` follows them. Player hits must go through `TakeHitFrom` (blockable); `TakeDamage` is unblockable.
  `TestBlueprintsValid` skips `EnemyGuard` in the every-state check (optional) and requires
  `GuardChance` in [0,1] with a guard clip. No strip pack has a `guard` clip. The mushroom
  art faces left (`FacesLeft`). The ledge probe acts only in `EnemyPatrol`: chasing enemies
  still run off edges. A new Attacker must implement `Ready` (a ranged one: its fire timer).
  Death clips of both golems, Enemy3 and the demon end on a blank frame (the body vanishes)
- **Open questions / blockers:** guard rules while mid-swing or hurt (decisions D8, open). Keep or delete `tools/battler2strip` (its goblin mapping
  points at a deleted sheet). No strip pack for a ranged enemy (task 25). Hook noise: `docs-check.mjs:292`
  counts untracked assets as changed code (issue proposed, not filed)
- **Don't forget:** assets changed on 2026-10-23. Removed: the desert set, `Bat with VFX/`,
  `slime/`. Skeleton is split into `Skeleton_White/` and `Skeleton_Yellow/`. Added:
  17 RPG Maker MV side-view battlers (deleted on 2026-10-24, see decisions). The golems still
  share `Golem_1_*.png` (D4); the user confirmed they are two different enemies. The enemy
  folders are untracked, but `//go:embed enemy` embeds them anyway; commit each with its code

## Recent sessions (newest first, keep 5)

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

### 2026-10-24: capability check

- Changed: `Has` + `capabilities` in `enemy_clips.go`; guard, second attack and `TestBlueprintsValid` use `Has`
- Verified: new tests failed first (build), mutation (fallback counted) caught; `verify.mjs` 5/5; 117 tests pass
- Next: task 16, White Golem

### 2026-10-24: strip packs only (D9)

- Changed: checked every strip pack (whole frames, steady feet, nothing cut off); Striker size fixed; spec, decisions (D9), tasks, current-state updated
- Verified: frame scan of `assets/enemy/`; `go test ./...` 115 pass; docs-check clean
- Next: task 16, White Golem

