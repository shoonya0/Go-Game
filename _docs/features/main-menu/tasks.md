# Main menu: tasks

## Increment 1: save that survives restarts (spec: [spec.md](./spec.md) "Save")

- [x] 1. `core.SaveData` (portals, checkpoint; settings kept apart, D2) with snapshot from the player and
  restore into world + player (checkpoint level, standing in it); bad data → error →
  check: `go test ./internal/core -run TestSave -v`
- [x] 2. Storage in `internal/system`: file on desktop, `localStorage` on the web (build
  tags) → check: `go test ./internal/system -run TestSaveFile -v`, `node scripts/verify.mjs --full`
- [x] 3. Autosave from `cmd/main.go` on portal save / travel (D3) → check: test + manual restart

## Increment 2: title screen

- [x] 4. Font: `go get` for `text/v2` (D1), a small text-drawing helper → check: `node scripts/verify.mjs --full`
- [ ] 5. Main menu screen (items, keys / tap / mouse, Continue grayed without a save); the game
  opens on it → check: `go test ./internal/core -run TestMenu -v`
- [ ] 6. New Game (confirm when a save exists, D4) and Continue → check: tests + `go run ./cmd`
- [ ] 7. Quit (desktop only, D6) → check: `node scripts/verify.mjs --full`

## Increment 3: settings and pause

- [ ] 8. Settings screen: volume 0–10, saved (D5) → check: `go test ./internal/core -run TestSettings -v`
- [ ] 9. Pause menu on Esc / touch menu button (D7) → check: tests
- [ ] 10. Web check (title, Continue after reload) with agent-browser; update docs → check:
  screenshot in `_docs/evidence/`, `node scripts/verify.mjs --full`

## Later (not planned yet)

- [ ] Real audio (music, effects) once sound files exist
- [ ] A second level (own feature)
