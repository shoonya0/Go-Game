package core

import (
	"bytes"
	"image"
	_ "image/jpeg" // register JPEG decoder (some level backgrounds are .jpg)
	_ "image/png"  // register PNG decoder
	"io/fs"
	"log"
	"path"
	"strings"

	"player/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

// embedded lists the embedded file systems, searched in order. Backgrounds is
// empty in the web build (levels decision D4).
var embedded = []fs.FS{assets.FS, assets.Backgrounds}

// findAsset returns the raw bytes of the embedded asset with p's base file
// name, or false when no embedded file has that name.
func findAsset(p string) ([]byte, bool) {
	base := path.Base(p)
	for _, fsys := range embedded {
		// Fast path: an asset embedded at the FS root.
		if data, err := fs.ReadFile(fsys, base); err == nil {
			return data, true
		}

		// Fall back to searching subfolders (e.g. enemy/) by base name, so callers
		// keep addressing assets by file name only.
		var found []byte
		_ = fs.WalkDir(fsys, ".", func(pth string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			if path.Base(pth) == base {
				if data, e := fs.ReadFile(fsys, pth); e == nil {
					found = data
					return fs.SkipAll
				}
			}
			return nil
		})
		if found != nil {
			return found, true
		}
	}
	return nil, false
}

// readAsset returns the raw bytes of an embedded asset. The caller may pass a
// full path (e.g. "../assets/player/GideonGraves.png" or "../assets/enemy/Mushroom/Mushroom-Idle.png");
// only the base file name is used, so the same call works on native and
// WebAssembly builds and regardless of which asset subfolder holds the file.
func readAsset(p string) []byte {
	data, ok := findAsset(p)
	if !ok {
		log.Fatalf("load asset %q: not found in embedded FS", p)
	}
	return data
}

// HasAsset reports whether asset p is embedded in this build. A level
// background isn't in the web build: fetch it from AssetURL instead.
func HasAsset(p string) bool {
	_, ok := findAsset(p)
	return ok
}

// AssetURL is where the web server serves asset p ("../assets/levels/x.png"
// → "levels/x.png"), relative to the game's page.
func AssetURL(p string) string {
	return strings.TrimPrefix(path.Clean(p), "../assets/")
}

// DecodeImage decodes downloaded image bytes. It doesn't touch the GPU, so it
// can run off the game loop; ebiten.NewImageFromImage uploads the result.
func DecodeImage(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	return img, err
}

// LoadImage decodes an embedded image asset and uploads it to a GPU-backed
// *ebiten.Image ready for drawing.
func LoadImage(p string) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(readAsset(p)))
	if err != nil {
		log.Fatalf("decode asset %q: %v", p, err)
	}
	return ebiten.NewImageFromImage(img)
}
