package core

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

// The portal art is one sheet, a 7×2 grid of a looping swirl:
// assets/levels/portal.png. The art drifts inside its cells (a few px per
// frame, and the second row sits higher), so each frame is cropped to its own
// art and drawn standing on its base, which keeps it still while it loops.
const (
	portalSheet            = "../assets/levels/portal.png"
	portalCols, portalRows = 7, 2
	portalFPS              = 5
	portalAlphaMin         = 40 // fainter pixels aren't art (a glow fringe, the sheet's corner mark)

	portalHeight = 2 * LevelTileHeight // drawn height of the art, and of its area
	portalWidth  = 1 * LevelTileWidth  // the doorway the player stands in

	// portalLift raises the drawn art above the ground, in px (negative sinks it
	// in). Only the picture moves: the doorway and where the player stands stay
	// on the ground.
	portalLift = -10.0
)

// savedPortalTint draws the saved portal golden: the blue art is turned gray,
// then scaled toward gold (scaling the blue art's RGB alone can't make gold).
var savedPortalTint = func() colorm.ColorM {
	var cm colorm.ColorM
	cm.ChangeHSV(0, 0, 1)
	cm.Scale(1.5, 1.15, 0.4, 1)
	return cm
}()

// PortalSprite draws every portal of a level with one shared animation.
type PortalSprite struct {
	img    *ebiten.Image
	frames []image.Rectangle // each frame's art, in sheet pixels
	scale  float64
	t      float64 // seconds the animation has run
}

// NewPortalSprite loads the portal sheet and finds its frames.
func NewPortalSprite() *PortalSprite {
	src, _, err := image.Decode(bytes.NewReader(readAsset(portalSheet)))
	if err != nil {
		log.Fatalf("decode %q: %v", portalSheet, err)
	}
	frames := artFrames(src, portalCols, portalRows, portalAlphaMin)
	tallest := 0
	for _, r := range frames {
		tallest = max(tallest, r.Dy())
	}
	return &PortalSprite{img: ebiten.NewImageFromImage(src), frames: frames, scale: portalHeight / float64(tallest)}
}

// Update advances the animation one tick.
func (s *PortalSprite) Update() { s.t += deltaTime() }

// Draw draws the current frame at each portal, centered on its cell with its
// base on the ground; the portal saved at respawn is golden.
func (s *PortalSprite) Draw(screen *ebiten.Image, portals []Portal, respawn, cam Position) {
	r := s.frames[int(s.t*portalFPS)%len(s.frames)]
	frame := s.img.SubImage(r).(*ebiten.Image)
	w, h := float64(r.Dx())*s.scale, float64(r.Dy())*s.scale
	for _, p := range portals {
		b := p.Bounds()
		var geo ebiten.GeoM
		geo.Scale(s.scale, s.scale)
		geo.Translate(centerX(b)-w/2-cam.X, b.Y+b.Height-h-portalLift-cam.Y)
		if p.SavedBy(respawn) {
			colorm.DrawImage(screen, frame, savedPortalTint, &colorm.DrawImageOptions{GeoM: geo})
			continue
		}
		screen.DrawImage(frame, &ebiten.DrawImageOptions{GeoM: geo})
	}
}

// Hitboxes shows the portal's area in the debug overlay (HitboxProvider).
func (p Portal) Hitboxes() []Hitbox { return []Hitbox{{Box: p.Bounds(), Kind: HitboxBody}} }

// artFrames finds each frame's art in a cols×rows grid sheet: the smallest
// rectangle in each cell holding every pixel more opaque than minAlpha.
func artFrames(img image.Image, cols, rows int, minAlpha uint8) []image.Rectangle {
	b := img.Bounds()
	cw, ch := b.Dx()/cols, b.Dy()/rows
	frames := make([]image.Rectangle, 0, cols*rows)
	alpha := func(x, y int) uint8 {
		_, _, _, a := img.At(x, y).RGBA()
		return uint8(a >> 8)
	}
	if n, ok := img.(*image.NRGBA); ok { // PNGs with transparency: read alpha directly, much faster
		alpha = func(x, y int) uint8 { return n.Pix[n.PixOffset(x, y)+3] }
	}
	for row := range rows {
		for col := range cols {
			art := image.Rectangle{}
			for y := b.Min.Y + row*ch; y < b.Min.Y+(row+1)*ch; y++ {
				for x := b.Min.X + col*cw; x < b.Min.X+(col+1)*cw; x++ {
					if alpha(x, y) > minAlpha {
						art = art.Union(image.Rect(x, y, x+1, y+1))
					}
				}
			}
			frames = append(frames, art)
		}
	}
	return frames
}
