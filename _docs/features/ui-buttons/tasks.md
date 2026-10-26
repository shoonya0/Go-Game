# UI buttons: tasks

## Increment 1: menus and travel arrows from the sheet (decisions: `decisions.md`)

- [x] 1. Sprite catalog + 3/9-slice drawing (`ui_skin.go`), touch rects moved there → check: `go test ./internal/core -run TestUISkin -v`
- [x] 2. Menu items as planks (arrow plank when picked), board behind the menu → check: `TestMenu*`, web build screenshot
- [x] 3. Icons left of the items → check: `TestMenuLayout`, screenshot
- [x] 4. Volume as the segmented bar with speaker / mute icon → check: `TestSettings*`, screenshot
- [x] 5. Travel menu scroll arrows from the sheet → check: `TestTravel*`, screenshot (tests pass at `8621c11`; the arrows show only with 5+ active portals, the user checked it in game)
- [x] 6. Update docs: state.md "Now", current-state → check: `node scripts/docs-check.mjs`

## Later (not planned yet)

- [ ] Native window check by hand (only the web build was checked)
- [ ] Crop the sheet to the sprites used (3.1 MB in the WASM; current-state known issue 9)
