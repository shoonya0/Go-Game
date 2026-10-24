# Enemy types: state

## Now (updated 2026-10-24 at `136582f`)

- **Status:** 🟡 in progress; increment 1 merged to `main` (PR #18, `2bc4469`)
- **Current task:** 11. Where the raw battler sheets live, out of the `//go:embed`
  (see [tasks.md](./tasks.md)): needs the user's answer first, then WASM size before/after
- **Stopped at:** task 10 done: goblin mapping committed, baked into `assets/enemy/Goblin/`
  (untracked, 9 strips, 423 KB; re-bake: the command in `tools/battler2strip/main.go`).
  Checked on a contact sheet: all clips face right, attacks hit on frame 4
- **Traps:** the tool keeps its own copy of the clip names (it can't import the game); keep it in sync with `enemy_clips.go`. `TestBlueprintsValid` needs speeds > 0, so the static training dummy (task 20)
  will need it relaxed for a `static` mover. The mushroom art faces left (`FacesLeft`). No ledge probe yet (task 12).
  A new Attacker must implement `Ready` (a ranged one: its fire timer)
- **Open questions / blockers:** task 11: where the raw battler sheets live once baked
  (the user decides; they must leave the embed). Asked 2026-10-24.
  Hook noise: `docs-check.mjs:292` counts untracked assets as changed code (issue proposed, not filed)
- **Don't forget:** assets changed on 2026-10-23. Removed: the desert set, `Bat with VFX/`,
  `slime/`. Skeleton is split into `Skeleton_White/` and `Skeleton_Yellow/`. Added:
  `Free Sprites/`, 17 RPG Maker MV side-view battlers (9×6 grid, 3 frames per motion, no walk
  cycle; 9 sheets have non-integer cell heights; facing varies per sheet). The golems still
  share `Golem_1_*.png` (D4); the user confirmed they are two different enemies. The enemy
  folders are untracked, but `//go:embed enemy` embeds them anyway; commit each with its code

## Recent sessions (newest first, keep 5)

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

### 2026-10-23: task 7, blueprint validation

- Changed: `internal/core/enemy_test.go` (`TestBlueprintsValid`)
- Verified: `verify.mjs --full` 8/8; mutation (mushroom without Attacker) fails it
- Docs touched: tasks (7 ticked), state

### 2026-10-23: task 6, Attacker + meleeAttack

- Changed: new `enemy_attack.go` + test; `enemy.go` (`HitFrame`, `Attacker`, `swing`), `enemy_ai.go`, `enemy_clips.go` (`clipFor` picks attack_2), `enemy_types.go`
- Verified: `verify.mjs --full` 8/8; reviewer measured missed swings when cooldown > clip (mushroom −40%), fixed test-first (`TestMeleeSwingLandsAfterCooldown`)
- Docs touched: design (Attacker done), tasks (6 ticked), state
