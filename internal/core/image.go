package core

import (
	"bytes"
	"image"
	_ "image/png" // register PNG decoder
	"io/fs"
	"log"
	"path"

	"player/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

// readAsset returns the raw bytes of an embedded asset. The caller may pass a
// full path (e.g. "../assets/GideonGraves.png" or "../assets/enemy/succubus.png");
// only the base file name is used, so the same call works on native and
// WebAssembly builds and regardless of which asset subfolder holds the file.
func readAsset(p string) []byte {
	base := path.Base(p)

	// Fast path: an asset embedded at the FS root.
	if data, err := assets.FS.ReadFile(base); err == nil {
		return data
	}

	// Fall back to searching subfolders (e.g. enemy/) by base name, so callers
	// keep addressing assets by file name only.
	var found []byte
	_ = fs.WalkDir(assets.FS, ".", func(pth string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if path.Base(pth) == base {
			if data, e := assets.FS.ReadFile(pth); e == nil {
				found = data
				return fs.SkipAll
			}
		}
		return nil
	})
	if found != nil {
		return found
	}

	log.Fatalf("load asset %q: not found in embedded FS", p)
	return nil
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
