# Go-Game (Pirate Adventure): docs index

Linked from CLAUDE.md. This page is a table of contents: one line per doc or feature, no
detail. Read top-down and stop as soon as you have what the task needs.

<!-- Budget: 80 lines. When the registry outgrows it, group features into areas
(kind: area) and give each area its own index. scripts/docs-check.mjs reports every
feature index that isn't reachable from this page. -->

## Start here

| Need | Read |
|---|---|
| Where the project stands, active work, test baseline | [current-state.md](./current-state.md) |
| Modules, layers, allowed imports, main flows | [architecture-overview.md](./architecture-overview.md) |
| How the AI harness (hooks, skills, verify, CI) was set up, and what's deferred | [ai-harness-setup.md](./ai-harness-setup.md) |
| Go backend skills distilled from this engine (reference; not loaded as Claude skills) | [rulesAndSkills/](./rulesAndSkills/README.md) |
| Increment log, written by `/close-increment` (created on first use) | `status/increments.md` |
| Why a cross-cutting decision was made (none yet) | `adr/` |

## Feature registry

One row per feature that has its own docs folder. Status: ⏳ planned · 🟡 in progress ·
✅ done · ⚠️ blocked · 💤 dormant.

| Feature | Docs index | Code | Status | One line |
|---|---|---|---|---|
| enemy-types | [index](./features/enemy-types/README.md) | `internal/core/enemy*.go`, `assets/enemy/`, `tools/battler2strip/` | 🟡 | Ground/flying enemies from the strip packs in `assets/enemy/`, built from movement + attack parts (low priority now) |
| portals | [index](./features/portals/README.md) | `internal/core/portal*.go`, `world.go`, `combat.go` | ✅ | Portals: start, save point (E / tap), travel menu; non-boss enemies revive on death, on entering a level and on every teleport |
| main-menu | [index](./features/main-menu/README.md) | `cmd/main.go`, `internal/core/menu*.go`, `save*.go`, `internal/system/save*.go` | ✅ | Title screen (New Game, Continue, Settings, Quit), pause menu, save across restarts, volume setting |
| ui-buttons | [index](./features/ui-buttons/README.md) | `internal/core/ui_skin*.go`, `menu.go`, `travel_menu.go` | ✅ | Buttons, menus and the volume bar drawn from the wooden sheet `assets/UI/UIButtons.png` |
| levels | [index](./features/levels/README.md) | `internal/core/levels*.go`, `world*.go`, `assets/levels/` | 🟡 | Eight themed levels, each a day and a night sub-level with its own tiles and background, joined by portals (increment 1 merged) |

<!-- Example row. Most features here live in the flat internal/core package, so their docs
go in _docs/features/<feature>/ and the index lists the feature's files in "scope:".
| enemies | [index](./features/enemies/README.md) | `internal/core/enemy*.go` | 🟡 | Enemy blueprints, AI and damage |
-->

## Other docs

| Doc | What it holds |
|---|---|
| [DEPLOY.md](../DEPLOY.md) | Build targets, WebAssembly build, Render deployment |
| [CLAUDE.md](../CLAUDE.md) | Agent entry file: map, commands, workflow, rules |

## Rules for these docs

- One fact, one home. Other docs link to it instead of repeating it.
- Document only what the code can't tell: intent (spec), reasons (decisions), progress
  (tasks, state) and traps. Structure and call flow come from the code graph.
- Feature docs live next to the feature's code (`<folder>/_docs/`), or in
  `_docs/features/<feature>/` when the feature has no single folder.
- Over budget (index 80, state 120, other docs 200 lines)? Split it into a folder with
  its own `README.md` index. Run `node scripts/docs-check.mjs` after any docs change.
