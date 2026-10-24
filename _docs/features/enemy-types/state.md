# Enemy types: state

## Now (updated 2026-10-24 at `f7b861a`)

- **Status:** 🟡 in progress; increments 1-2 merged (PR #18 `2bc4469`, PR #19 `9124414`);
  increment 3 (ground roster + guard): tasks 12-14 done, 15/19/20 dropped; committed locally, not pushed
- **Current task:** next is 16 (White Golem: rename to `WhiteGolem_*`, add it) in [tasks.md](./tasks.md),
  on `feature/enemy-types` from `main` at `9124414`
- **Stopped at:** strip packs are now the only enemy format (D9, user's call): goblin
  reverted (`12de28a`, `3b98b41`), battler sheets deleted (`b450bdd`). Every strip pack
  checked frame by frame: all fine; Striker's frames are 96×96, attack 128×96 (`f7b861a`)
- **Traps:** player hits must go through `TakeHitFrom` (blockable); `TakeDamage` is unblockable.
  `TestBlueprintsValid` skips `EnemyGuard` in the every-state check (optional) and requires
  `GuardChance` in [0,1] with a guard clip. No strip pack has a `guard` clip. The mushroom
  art faces left (`FacesLeft`). The ledge probe acts only in `EnemyPatrol`: chasing enemies
  still run off edges. A new Attacker must implement `Ready` (a ranged one: its fire timer).
  Death clips of both golems, Enemy3 and the demon end on a blank frame (the body vanishes)
- **Open questions / blockers:** keep or delete `tools/battler2strip` (its goblin mapping
  points at a deleted sheet). No strip pack for a ranged enemy (task 25). My call, not in
  D8: an enemy mid-swing can't guard (user not asked yet). Hook noise: `docs-check.mjs:292`
  counts untracked assets as changed code (issue proposed, not filed)
- **Don't forget:** assets changed on 2026-10-23. Removed: the desert set, `Bat with VFX/`,
  `slime/`. Skeleton is split into `Skeleton_White/` and `Skeleton_Yellow/`. Added:
  17 RPG Maker MV side-view battlers (deleted on 2026-10-24, see decisions). The golems still
  share `Golem_1_*.png` (D4); the user confirmed they are two different enemies. The enemy
  folders are untracked, but `//go:embed enemy` embeds them anyway; commit each with its code

## Recent sessions (newest first, keep 5)

### 2026-10-24: strip packs only (D9)

- Changed: checked every strip pack (whole frames, steady feet, nothing cut off); Striker size fixed; spec, decisions (D9), tasks, current-state updated
- Verified: frame scan of `assets/enemy/`; `go test ./...` 115 pass; docs-check clean
- Next: task 16, White Golem

### 2026-10-24: battler sheets deleted

- Changed: `art-src/battlers/` removed (18 sheets + layout guide); tasks 15, 19, 20 dropped; scope and registry updated
- Verified: measured every sheet (all 9×6, 3 frames per motion; walk rows differ 1-3 %); looked at dog and bat
- Next: talk with the user about the enemy plan

### 2026-10-24: goblin removed

- Changed: reverted the goblin (`12de28a`, `3b98b41`); the user turned it down because it had no walk animation
- Verified: `verify.mjs` 5/5 after the revert
- Next: talk with the user about the enemy plan

### 2026-10-24: task 15, goblin (built, look pending)

- Changed: `battler2strip` (`floorStrips`, bake all clips before writing); goblin strips re-baked and committed; `goblinBlueprint`, `KindGoblin`, `$` in level 1
- Verified: `TestFloorStripsMovesFeetToBottom` failed first; `verify.mjs --full` 8/8; feet measured at y=119/120; web screenshots with and without F3
- Docs touched: design (feet shift), tasks (15 note), state, current-state

### 2026-10-24: task 14, guard trait

- Changed: new `enemy_guard.go` (`TakeHitFrom`, `rollGuard`) + tests; `EnemyGuard`, `GuardChance`; swings and the ground wave call `TakeHitFrom`
- Verified: `verify.mjs --full` 8/8; mutations (front check, each caller's side) fail the tests; code-reviewer: no blocking findings (2 suggestions applied)
- Docs touched: design (guard notes), tasks (14 ticked), state, current-state

