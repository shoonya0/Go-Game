# Enemy types: state

## Now (updated 2026-10-23 at `fda2323`)

- **Status:** ⏳ planned
- **Current task:** 1. Characterisation tests for the succubus (see [tasks.md](./tasks.md))
- **Stopped at:** planning done; branch `feature/enemy-types` created from `origin/main`, no code yet
- **Open questions / blockers:** none blocking. Increment 1's interface shapes in
  [design.md](./design.md) are a proposal; adjust them while implementing
- **Don't forget:** `White_Golem/` and `Orange_Golem/` share file names (`Golem_1_*.png`), so
  one shadows the other in `readAsset` (D4). The enemy folders are untracked, but
  `//go:embed enemy` embeds whatever is on disk, so the build already contains them. Commit
  each folder with the code that loads it

## Recent sessions (newest first, keep 5)

### 2026-10-23: plan the feature

- Changed: docs only (`_docs/features/enemy-types/`, registry, Active work, rules pointer)
- Verified: `go build ./...` ok with the untracked folders; asset base-name scan found the golem collision
- Docs touched: README, spec, design, decisions, tasks, state

## Known issues

- Duplicate asset base names `Golem_1_*.png` (`assets/enemy/White_Golem/`, `assets/enemy/Orange_Golem/`), fixed by task 10
