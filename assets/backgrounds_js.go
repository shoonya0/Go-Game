//go:build js

package assets

import "embed"

// Backgrounds is empty on the web: the game downloads each level's background
// from the server when the player enters the level (levels decision D4).
var Backgrounds embed.FS
