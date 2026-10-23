package core

import (
	"io/fs"
	"path"
	"testing"

	"player/assets"
)

// TestAssetPathsAreEmbedded catches an asset that was moved to a folder the
// embed pattern in assets/embed.go doesn't cover: the build still passes, but
// readAsset would log.Fatal at startup.
func TestAssetPathsAreEmbedded(t *testing.T) {
	paths := []string{
		Tileset, LevelFile, Background_1,
		HealthSheet, StaminaSheet, EnemyHealthSheet, UIButtonsSheet,
		"../assets/player/GideonGraves.png", // cmd/main.go playerSpriteSheetPath
	}
	for _, bp := range enemyBlueprints {
		for _, c := range bp.Clips {
			paths = append(paths, c.Sheet)
		}
	}

	embedded := map[string]bool{}
	_ = fs.WalkDir(assets.FS, ".", func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			embedded[path.Base(p)] = true
		}
		return nil
	})

	for _, p := range paths {
		if !embedded[path.Base(p)] {
			t.Errorf("asset %q is not in the embedded FS; check assets/embed.go", p)
		}
	}
}
