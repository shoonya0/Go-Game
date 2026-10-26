---
kind: feature
scope: [internal/core/ui_skin*.go, internal/core/menu.go, internal/core/menu_screens.go, internal/core/travel_menu.go, internal/core/touch_controls.go, assets/UI/UIButtons.png]
status: in-progress
---
# UI buttons: docs index

Every button in the game is drawn from the wooden sheet `assets/UI/UIButtons.png`: menu
items are wooden planks (the picked one is the arrow plank), menus sit on the wooden board,
items carry an icon, the volume is the segmented bar, and the travel menu's scroll arrows are
the wooden arrow buttons. Touch controls and the portal prompt already used the sheet.
Up: [project docs index](../../README.md) · Code: `internal/core/ui_skin.go`, `menu.go`, `travel_menu.go`

## Read in this order (stop when you have enough)

| Need | File |
|---|---|
| Where work stopped | [state.md](./state.md) |
| What's done and what's next | [tasks.md](./tasks.md) |
| Why it's built this way | [decisions.md](./decisions.md) |

## Entry points (where to start reading code; not a full map)

| Path | Why start here |
|---|---|
| `internal/core/ui_skin.go` | the sprite rects on the sheet and the slicing that stretches them |
| `internal/core/menu.go` | `Menu.Draw`: planks, board, icons, bar |

## Depends on / used by

- Uses: main-menu ([index](../main-menu/README.md)) for the menus it skins
