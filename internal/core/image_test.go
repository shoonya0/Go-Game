package core

import (
	"io/fs"
	"path"
	"testing"

	"player/assets"
)

// TestAssetBaseNamesUnique checks no two embedded assets share a file name.
// readAsset finds a file by base name alone, so a second file with the same
// name in another folder would load in place of the first: the White and
// Orange golems once both shipped Golem_1_*.png.
func TestAssetBaseNamesUnique(t *testing.T) {
	seen := map[string]string{}
	err := fs.WalkDir(assets.FS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		base := path.Base(p)
		if first, dup := seen[base]; dup {
			t.Errorf("%s and %s share the base name %s", first, p, base)
		}
		seen[base] = p
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
