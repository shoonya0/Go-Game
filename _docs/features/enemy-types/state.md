# Enemy types: state

## Now (updated 2026-10-24 at `f12857c`)

- **Status:** 🟡 in progress; increment 1 merged (PR #18, `2bc4469`); increment 2 (battler baking, tasks 9-11) merged (PR #19, `9124414`)
- **Current task:** 13. Two-character level markers (D6) (see [tasks.md](./tasks.md)); increment 3
  on `feature/enemy-types` from `main` at `9124414`
- **Stopped at:** task 12 done (`f12857c`): patrolling ground enemies turn at ledges and at
  water. Goblin strips in `assets/enemy/Goblin/` stay untracked until task 15
- **Traps:** the tool keeps its own copy of the clip names (it can't import the game); keep it in sync with `enemy_clips.go`. `TestBlueprintsValid` needs speeds > 0, so the static training dummy (task 20)
  will need it relaxed for a `static` mover. The mushroom art faces left (`FacesLeft`). The ledge probe acts only in `EnemyPatrol`: chasing enemies still run off edges after the player.
  A new Attacker must implement `Ready` (a ranged one: its fire timer)
- **Open questions / blockers:** none for task 13. Hook noise: `docs-check.mjs:292` counts untracked assets as changed code (issue proposed, not filed)
- **Don't forget:** assets changed on 2026-10-23. Removed: the desert set, `Bat with VFX/`,
  `slime/`. Skeleton is split into `Skeleton_White/` and `Skeleton_Yellow/`. Added:
  17 RPG Maker MV side-view battlers (now in `art-src/battlers/`; 9×6 grid, 3 frames per motion, no walk
  cycle; 9 sheets have non-integer cell heights; facing varies per sheet). The golems still
  share `Golem_1_*.png` (D4); the user confirmed they are two different enemies. The enemy
  folders are untracked, but `//go:embed enemy` embeds them anyway; commit each with its code

## Recent sessions (newest first, keep 5)

### 2026-10-24: task 12, ledge probe

- Changed: `enemy_move.go` (`groundMover.turnAtLedge`); new `TestGroundEnemyTurnsAtLedge`
- Verified: test failed first (right edge 721 > 720 at tick 34), then passes; `verify.mjs` 5/5
- Docs touched: tasks (12 ticked), state, current-state

### 2026-10-24: task 11, raw sheets out of the embed

- Changed: `assets/enemy/Free Sprites/` → `art-src/battlers/` (committed); `goblin.json` + `bake.go` paths; `.dockerignore`
- Verified: WASM 51.2 → 36.2 MB (`go build` js/wasm before/after); goblin re-bake byte-identical; `verify.mjs --full`
- Docs touched: decisions (D7 where raw sheets live), spec, design, tasks (11 ticked), state, current-state, architecture

### 2026-10-24: task 10, goblin mapping + contact sheet

- Changed: `tools/battler2strip` (`-contact`, clip-order baking, `TestContactSheetStacksStrips`); new `mappings/goblin.json`
- Verified: `verify.mjs --full` 8/8; goblin contact sheet checked row by row (9 clips)
- Docs touched: tasks (10 ticked), state, current-state

### 2026-10-23: task 9, battler2strip

- Changed: new `tools/battler2strip/` (`bake.go`, `main.go`, `bake_test.go`); `.golangci.yml` (`tools` stdlib-only); `architecture-overview.md` (layers + modules rows)
- Verified: `verify.mjs --full` 8/8; `go test ./tools/...` 7/7; mutations (truncated cell edges, game import in tools) fail; code-reviewer: no blocking findings (design.md mapping example was stale, fixed)
- Docs touched: design (baking section), tasks (9 ticked), state, architecture

### 2026-10-23: task 8, close increment 1

- Changed: `enemy_attack.go` (`Attacker.Ready`, `hitTime`), `enemy_ai.go` (idle until ready), `enemy_clips.go` (`timeTo`); new `TestMeleeEverySwingLands`
- Verified: `verify.mjs --full` 8/8 before and after; the test failed for all 4 species before the fix; code-reviewer: one blocking finding (whiffed swings), fixed
- Docs touched: design (trimmed to what the code doesn't show), tasks (8 ticked), state, current-state

