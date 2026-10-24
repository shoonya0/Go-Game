# Enemy types: state

## Now (updated 2026-10-24 at `477ca01`)

- **Status:** 🟡 in progress; increments 1-2 merged (PR #18 `2bc4469`, PR #19 `9124414`);
  increment 3 (ground roster + guard): tasks 12-14, 16 done, 15/19/20 dropped; committed locally, not pushed
- **Current task:** 21 (place one of each, web smoke check, close increment 3) in [tasks.md](./tasks.md),
  on `feature/enemy-types` from `main` at `9124414`
- **Stopped at:** task 16 done: White Golem strips renamed `WhiteGolem_*` and committed
  (`7dededa`), `TestAssetBaseNamesUnique` (`ce20974`), blueprint + `%` on the platform at
  columns 44-49 of level 1 (`36fc396`), `TestLevelEnemiesKeepTheirGround` (`d21c103`).
  `verify.mjs --full` 8/8. User checked it in game: "working".
  Then removed the succubus (user: animation not good): tests now run on the test-only
  `testEnemyBlueprint` (`4917db1`, `a4a5b53`), blueprint + `!` + sheet gone (`0f42350`)
- **Traps:** behaviour tests use `testEnemyBlueprint` (enemy_test.go), not a game enemy: change the roster freely. ask `bp.Has(clip)` for "can it?", not `bp.Clips[...]`; `Has` ignores fallbacks, `resolveClip` follows them. Player hits must go through `TakeHitFrom` (blockable); `TakeDamage` is unblockable.
  `TestBlueprintsValid` skips `EnemyGuard` in the every-state check (optional) and requires
  `GuardChance` in [0,1] with a guard clip. No strip pack has a `guard` clip. The mushroom
  art faces left (`FacesLeft`). The ledge probe acts only in `EnemyPatrol`: chasing enemies
  still run off edges. A new Attacker must implement `Ready` (a ranged one: its fire timer).
  Death clips of both golems, Enemy3 and the demon end on a blank frame (the body vanishes)
- **Open questions / blockers:** Keep or delete `tools/battler2strip` (its goblin mapping
  points at a deleted sheet). No strip pack for a ranged enemy (task 25). My call, not in
  D8: an enemy mid-swing can't guard (user not asked yet). Hook noise: `docs-check.mjs:292`
  counts untracked assets as changed code (issue proposed, not filed)
- **Don't forget:** assets changed on 2026-10-23. Removed: the desert set, `Bat with VFX/`,
  `slime/`. Skeleton is split into `Skeleton_White/` and `Skeleton_Yellow/`. Added:
  17 RPG Maker MV side-view battlers (deleted on 2026-10-24, see decisions). The golems still
  share `Golem_1_*.png` (D4); the user confirmed they are two different enemies. The enemy
  folders are untracked, but `//go:embed enemy` embeds them anyway; commit each with its code

## Recent sessions (newest first, keep 5)

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

### 2026-10-24: battler sheets deleted

- Changed: `art-src/battlers/` removed (18 sheets + layout guide); tasks 15, 19, 20 dropped; scope and registry updated
- Verified: measured every sheet (all 9×6, 3 frames per motion; walk rows differ 1-3 %); looked at dog and bat
- Next: talk with the user about the enemy plan

