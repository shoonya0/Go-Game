package core

import (
	"bytes"
	"image"
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
	for _, fsys := range embedded {
		err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
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
}

// TestBackgroundsDecode checks every embedded level background can be
// decoded: some are JPEGs, and image.Decode only knows formats whose
// decoder is imported.
func TestBackgroundsDecode(t *testing.T) {
	matches, err := fs.Glob(assets.Backgrounds, "levels/*/background/*")
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) == 0 {
		t.Fatal("no level backgrounds embedded")
	}
	for _, p := range matches {
		data, err := assets.Backgrounds.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if _, format, err := image.DecodeConfig(bytes.NewReader(data)); err != nil {
			t.Errorf("%s: %v", p, err)
		} else {
			t.Logf("%s: %s", p, format)
		}
	}
}

// TestBackgroundsOutOfMainFS: the backgrounds are embedded only through
// assets.Backgrounds, which the web build leaves empty, so game.wasm doesn't
// carry several MB per level (levels decision D4).
func TestBackgroundsOutOfMainFS(t *testing.T) {
	matches, err := fs.Glob(assets.FS, "levels/*/background/*")
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) > 0 {
		t.Errorf("backgrounds in assets.FS (they would be in game.wasm): %v", matches)
	}
	if !HasAsset(Background_1) {
		t.Errorf("HasAsset(%q) = false on a desktop build", Background_1)
	}
	if HasAsset("../assets/levels/none/background/none.png") {
		t.Error("HasAsset reports a file that doesn't exist")
	}
}

// TestAssetURL: the web build fetches a background from the path the server
// serves it at, relative to the page.
func TestAssetURL(t *testing.T) {
	got := AssetURL("../assets/levels/deepForest/background/deepForest_day.png")
	if want := "levels/deepForest/background/deepForest_day.png"; got != want {
		t.Errorf("AssetURL = %q, want %q", got, want)
	}
}

// TestDecodeImage decodes downloaded bytes and reports bad data as an error,
// not a crash (a failed download must not end the game).
func TestDecodeImage(t *testing.T) {
	img, err := DecodeImage(readAsset(Background_1))
	if err != nil || img.Bounds().Dx() == 0 {
		t.Fatalf("DecodeImage(background) = %v, %v", img, err)
	}
	if _, err := DecodeImage([]byte("not an image")); err == nil {
		t.Error("DecodeImage accepted bytes that aren't an image")
	}
}
