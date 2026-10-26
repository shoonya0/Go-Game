# Enemy types: state

## Now (updated 2026-10-26 at `8c2093e`)

- **Status:** 🟡 playtest fixes on `fix/30-enemy-playtest-bugs` (issues #30-#35, D11), committed locally, not pushed; task 25 still waits for ranged art
- **Current task:** task 34: the user plays the redone knockdown (#33) and the doubled detection; then push and PR (ask before merge)
- **Stopped at:** playtest: #30-#32, #34, #35 fine. #33 redone (heavy hit turns the player toward the source); DetectRange doubled (600-760 px). `verify.mjs --full` 8/8. Wider detection may wake enemies placed 6+ columns from portals (levels state): check on arrival
- **Traps:** behaviour tests use `testEnemyBlueprint` (enemy_test.go), not a game enemy: change the roster freely; flyer tests use `testFlyerBlueprint` (enemy_fly_test.go). Ask `bp.Has(clip)` for "can it?", not `bp.Clips[...]`; `Has` ignores fallbacks, `resolveClip` follows them. Player hits must go through `TakeHitFrom` (blockable); `TakeDamage` is unblockable.
  `TestBlueprintsValid` skips `EnemyGuard` in the every-state check (optional) and requires
  `GuardChance` in [0,1] with a guard clip. No strip pack has a `guard` clip. The mushroom
  art faces left (`FacesLeft`). `decideBehavior` takes the quadtree: nil in a test means no terrain (every player reachable, every line clear).
  `spawnX`/`spawnY` are the place it guards, moved when it patrols after a chase (`engaged`); `RespawnEnemies` builds fresh enemies at the grid spawns.
  Walkers step up at most `maxStepUp` (half a tile, to leave water) and don't jump. Tests that hit an enemy with a ground wave must put it on the floor
  with `standOnWaveFloor` (sets `onGround`). A new Attacker must implement `Ready` (a ranged one: its fire timer) and `Committed`; a new Mover `CanChase`.
  Detection is still horizontal only (`DetectRange` on dx). Death clips of both golems, Enemy3 and the demon end on a blank frame
- **Open questions / blockers:** guard rules while mid-swing or hurt (decisions D8, open). Keep or delete `tools/battler2strip` (its goblin mapping
  points at a deleted sheet). No strip pack for a ranged enemy (task 25). Hook noise: `docs-check.mjs:292`
  counts untracked assets as changed code (issue proposed, not filed)
- **Don't forget:** the remaining enemy folders are untracked, but `//go:embed enemy` embeds them
  anyway; commit each with its code. Rename the Orange Golem files `OrangeGolem_*` before
  using them (D4)

## Recent sessions (newest first, keep 5)

### 2026-10-26: touched by levels (PR #28 merged)

- Changed: the Orange Golem strips are committed (not loaded yet); `TestLevelEnemiesKeepTheirGround` sizes its quadtree from the grid and runs on all 16 sub-levels ([levels](../levels/state.md)); this feature's behaviour unchanged
- Verified: `node scripts/verify.mjs --full` 8/8
- Docs touched: this entry only

### 2026-10-26: touched by levels (tasks 6, 7)

- Changed: `TestLevelEnemiesKeepTheirGround` runs on every level (deepForest night has six enemies) ([levels](../levels/state.md)); this feature's behaviour unchanged
- Verified: `node scripts/verify.mjs --full` → 8/8
- Docs touched: this entry only

### 2026-10-25: portals work in this scope

- Changed (by the portals feature): `EnemyBlueprint.Boss`, `RespawnEnemies` + `enemy_revive_test.go`
- Verified: `TestRespawnRevivesEnemies`, `TestKilledBossStaysDead` failed first (build); `verify.mjs --full` 8/8
- Docs touched: state

### 2026-10-24: increments 4-5 merged

- Changed: PR #21 merged by the user (`f6a311a`); remote branch gone, local `main` synced, `feature/enemy-types` restarted from it
- Verified: `verify.mjs --full` 8/8 before the push; branch was up to date with `main`
- Next: task 25 when ranged art exists, or a "Later" item (bosses, mushroom stun)

### 2026-10-24: task 27 checks (partial)

- Changed: flyer wall-pin fix in `flyingMover.Move` + `TestStrandedFlyerPatrolsWhereItIs` (review note)
- Verified: new test failed first (pinned at the wall); web build loads, no console errors, demon + F3 boxes draw; `code-reviewer`: no blockers; `verify.mjs` 5/5
- Next: PR (ask before merge); the user confirmed "Smasher working". Review gaps, not bugs: no test for hurt mid-dive
