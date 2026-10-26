# UI buttons: tasks

## Increment 1: menus and travel arrows from the sheet (decisions: `decisions.md`)

- [ ] 1. Sprite catalog + 3/9-slice drawing (`ui_skin.go`), touch rects moved there → check: `go test ./internal/core -run TestUISkin -v`
- [ ] 2. Menu items as planks (arrow plank when picked), board behind the menu → check: `TestMenu*`, web build screenshot
- [ ] 3. Icons left of the items → check: `TestMenuLayout`, screenshot
- [ ] 4. Volume as the segmented bar with speaker / mute icon → check: `TestSettings*`, screenshot
- [ ] 5. Travel menu scroll arrows from the sheet → check: `TestTravel*`, screenshot
- [ ] 6. Update docs: state.md "Now", current-state → check: `node scripts/docs-check.mjs`

## Later (not planned yet)

- [ ] Crop the sheet to the sprites used (3.1 MB in the WASM; current-state known issue 9)
