//go:build !js

package assets

import "embed"

// Backgrounds holds the level backgrounds. Desktop builds embed them all; the
// web build leaves them out and downloads only the level being played
// (levels decision D4), because they are several MB each.
//
//go:embed levels/*/background
var Backgrounds embed.FS
