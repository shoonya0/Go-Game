// Package assets embeds the game's image files directly into the binary so the
// game is self-contained and runs anywhere, including WebAssembly where there is
// no host filesystem to read from.
package assets

import "embed"

// FS holds every asset but the level backgrounds (see Backgrounds). Files in
// a level folder are embedded by pattern: a new kind of file there needs its
// own pattern here.
//
//go:embed enemy player UI levels/*.png levels/*/*.json
var FS embed.FS
