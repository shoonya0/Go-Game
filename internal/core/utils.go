package core

import (
	"bytes"
	"image"
	_ "image/png" // register PNG decoder
	"log"
	"path"

	"player/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

// LoadImage loads an image by name from the embedded asset bundle. The caller
// may pass a full path (e.g. "../assets/GideonGraves.png"); only the base file
// name is used, so the same call works on native and WebAssembly builds.
func LoadImage(p string) *ebiten.Image {
	data, err := assets.FS.ReadFile(path.Base(p))
	if err != nil {
		log.Fatalf("load asset %q: %v", p, err)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("decode asset %q: %v", p, err)
	}

	return ebiten.NewImageFromImage(img)
}
