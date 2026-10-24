# Enemy types: state

## Now (updated 2026-10-24 at `849b30e`)

- **Status:** 🟡 in progress; increments 1-2 merged (PR #18 `2bc4469`, PR #19 `9124414`);
  increment 3 (ground roster + guard): tasks 12-14 done, 15 built; committed locally, not pushed
- **Current task:** 15. Goblin: waiting on the user's native look, then tick it; next 16 (see [tasks.md](./tasks.md)),
  on `feature/enemy-types` from `main` at `9124414`
- **Stopped at:** task 15 built: bake tool puts feet on the bottom row (`1615825`), goblin
  strips re-baked (`d32fe9f`), blueprint + `$` at level 1 column 15 by the pit (`849b30e`).
  Web build checked: stands on the ground, turns at the pit, F3 boxes fit, no console errors
- **Traps:** re-bake every battler after the feet fix (the goblin is the only one so far).
  The goblin's idle, walk and run strips are identical: rows 0 and 1 of its sheet hold the
  same art, not a bake bug; walk and run get their motion from draw-time bob. Player hits must go through `TakeHitFrom` (blockable); `TakeDamage` is unblockable.
  `TestBlueprintsValid` skips `EnemyGuard` in the every-state check (optional) and requires
  `GuardChance` in [0,1] with a guard clip. The tool keeps its own copy of the clip names; keep
  it in sync with `enemy_clips.go`. `TestBlueprintsValid` needs speeds > 0, so the static
  training dummy (task 20) will need it relaxed for a `static` mover. The mushroom art faces
  left (`FacesLeft`). The ledge probe acts only in `EnemyPatrol`: chasing enemies still run off
  edges. A new Attacker must implement `Ready` (a ranged one: its fire timer)
- **Open questions / blockers:** user's native look at the goblin (task 15 check). Goblin
  tuning is a first guess (40 HP, 7 dmg, guard 35 %). My call, not in D8: an enemy mid-swing
  can't guard (user not asked yet). Hook noise: `docs-check.mjs:292` counts untracked assets as changed code (issue
  proposed, not filed)
- **Don't forget:** assets changed on 2026-10-23. Removed: the desert set, `Bat with VFX/`,
  `slime/`. Skeleton is split into `Skeleton_White/` and `Skeleton_Yellow/`. Added:
  17 RPG Maker MV side-view battlers (now in `art-src/battlers/`; 9×6 grid, 3 frames per motion, no walk
  cycle; 9 sheets have non-integer cell heights; facing varies per sheet). The golems still
  share `Golem_1_*.png` (D4); the user confirmed they are two different enemies. The enemy
  folders are untracked, but `//go:embed enemy` embeds them anyway; commit each with its code

## Recent sessions (newest first, keep 5)

### 2026-10-24: task 15, goblin (built, look pending)

- Changed: `battler2strip` (`floorStrips`, bake all clips before writing); goblin strips re-baked and committed; `goblinBlueprint`, `KindGoblin`, `$` in level 1
- Verified: `TestFloorStripsMovesFeetToBottom` failed first; `verify.mjs --full` 8/8; feet measured at y=119/120; web screenshots with and without F3
- Docs touched: design (feet shift), tasks (15 note), state, current-state

### 2026-10-24: task 14, guard trait

- Changed: new `enemy_guard.go` (`TakeHitFrom`, `rollGuard`) + tests; `EnemyGuard`, `GuardChance`; swings and the ground wave call `TakeHitFrom`
- Verified: `verify.mjs --full` 8/8; mutations (front check, each caller's side) fail the tests; code-reviewer: no blocking findings (2 suggestions applied)
- Docs touched: design (guard notes), tasks (14 ticked), state, current-state

### 2026-10-24: task 13, two-character markers

- Changed: `level_enemies.go` (`quoteBareMarkers` quotes runs, `maxMarkerLen`); tests `TestTwoCharMarkerIsOneCell`, `TestEnemyMarkersUniqueAndBare` (renamed); marker comments
- Verified: new tests failed first (`"!""@"`), then pass; `TestLevelFileSpawnsEnemies` passes; `verify.mjs` 5/5
- Docs touched: decisions (D6 done), tasks (13 ticked), state, current-state

### 2026-10-24: task 12, ledge probe

- Changed: `enemy_move.go` (`groundMover.turnAtLedge`); new `TestGroundEnemyTurnsAtLedge`
- Verified: test failed first (right edge 721 > 720 at tick 34), then passes; `verify.mjs` 5/5
- Docs touched: tasks (12 ticked), state, current-state

### 2026-10-24: task 11, raw sheets out of the embed

- Changed: `assets/enemy/Free Sprites/` → `art-src/battlers/` (committed); `goblin.json` + `bake.go` paths; `.dockerignore`
- Verified: WASM 51.2 → 36.2 MB (`go build` js/wasm before/after); goblin re-bake byte-identical; `verify.mjs --full`
- Docs touched: decisions (D7 where raw sheets live), spec, design, tasks (11 ticked), state, current-state, architecture

