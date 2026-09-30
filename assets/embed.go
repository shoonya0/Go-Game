// Package assets embeds the game's image files directly into the binary so the
// game is self-contained and runs anywhere, including WebAssembly where there is
// no host filesystem to read from.
package assets

import "embed"

//go:embed *.png
var FS embed.FS
