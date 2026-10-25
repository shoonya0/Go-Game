package core

import (
	"bytes"
	"image"
	"log"
	"math"

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

	portalSavedSeconds = 1.5 // how long the "saved" check shows
	portalPromptSize   = 48  // side of the prompt button, px
	portalPromptGap    = 6   // between the portal's top and the prompt
)

// Prompt symbols, drawn like the touch-button icons: E to use the portal on a
// keyboard, a ring to tap on touch screens, and a check once it is saved.
var (
	promptKeyE = []iconShape{
		lineShape(0.09, [2]float32{0.64, 0.27}, [2]float32{0.38, 0.27}, [2]float32{0.38, 0.73}, [2]float32{0.64, 0.73}),
		lineShape(0.09, [2]float32{0.38, 0.5}, [2]float32{0.58, 0.5}),
	}
	promptTap = []iconShape{
		lineShape(0.06, append(starPoints(16, 0.3, 0.3, 0.5, 0.5), [2]float32{0.5, 0.2})...),
		fillShape(starPoints(8, 0.11, 0.11, 0.5, 0.5)...),
	}
	promptCheck = []iconShape{
		lineShape(0.11, [2]float32{0.27, 0.52}, [2]float32{0.44, 0.69}, [2]float32{0.74, 0.33}),
	}
)

// tapArea is where a tap uses the portal: its doorway plus half a tile around
// it, since a fingertip is less exact than a mouse, reaching up over the prompt.
func (c Portal) tapArea() AABB {
	b := c.Bounds()
	top := min(b.Y-LevelTileHeight/2, b.Y-portalLift-portalPromptGap-portalPromptSize)
	return AABB{X: b.X - LevelTileWidth/2, Y: top, Width: b.Width + LevelTileWidth, Height: b.Y + b.Height + LevelTileHeight/2 - top}
}

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

	button, keyE, tap, check *ebiten.Image // the prompt over the portal
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
	return &PortalSprite{
		img: ebiten.NewImageFromImage(src), frames: frames, scale: portalHeight / float64(tallest),
		button: LoadImage(UIButtonsSheet).SubImage(uiRoundUpRect).(*ebiten.Image),
		keyE:   newTouchIcon(promptKeyE),
		tap:    newTouchIcon(promptTap),
		check:  newTouchIcon(promptCheck),
	}
}

// Update advances the animation one tick.
func (s *PortalSprite) Update() { s.t += deltaTime() }

// Draw draws each portal centered on its cell with its base on the ground, as
// player sees it (PortalLook): an inactive one still on the first frame, an
// active one looping, the checkpoint looping in gold.
func (s *PortalSprite) Draw(screen *ebiten.Image, portals []Portal, player *PlayerRuntime) {
	cam := player.Camera.Pos
	loop := int(s.t*portalFPS) % len(s.frames)
	for _, p := range portals {
		look := player.PortalLook(p)
		r := s.frames[loop]
		if look == PortalInactive {
			r = s.frames[0]
		}
		frame := s.img.SubImage(r).(*ebiten.Image)
		w, h := float64(r.Dx())*s.scale, float64(r.Dy())*s.scale
		b := p.Bounds()
		var geo ebiten.GeoM
		geo.Scale(s.scale, s.scale)
		geo.Translate(centerX(b)-w/2-cam.X, b.Y+b.Height-h-portalLift-cam.Y)
		if look == PortalCheckpoint {
			colorm.DrawImage(screen, frame, savedPortalTint, &colorm.DrawImageOptions{GeoM: geo})
			continue
		}
		screen.DrawImage(frame, &ebiten.DrawImageOptions{GeoM: geo})
	}
}

// DrawPrompt draws, over the portal the player stands in, a pulsing button
// showing how to use it (E, or a tap ring on touch screens). For a moment after
// a save, a check over the saved portal takes its place.
func (s *PortalSprite) DrawPrompt(screen *ebiten.Image, portals []Portal, p *PlayerRuntime, touch bool) {
	icon, at, alpha := s.keyE, -1, float32(0.75+0.25*math.Sin(s.t*2*math.Pi))
	if touch {
		icon = s.tap
	}
	if p.savedNote > 0 {
		for i, cp := range portals {
			if p.PortalLook(cp) == PortalCheckpoint {
				icon, at, alpha = s.check, i, float32(min(1, p.savedNote/0.4)) // fades out at the end
			}
		}
	} else {
		at = p.atPortal
	}
	if at < 0 || at >= len(portals) {
		return
	}
	b := portals[at].Bounds()
	x := int(centerX(b) - portalPromptSize/2 - p.Camera.Pos.X)
	y := int(b.Y - portalLift - portalPromptGap - portalPromptSize - p.Camera.Pos.Y)
	r := image.Rect(x, y, x+portalPromptSize, y+portalPromptSize)
	drawScaled(screen, s.button, r, nil, alpha)
	drawScaled(screen, icon, r.Inset(portalPromptSize*touchIconInset/100), touchIconColor, alpha)
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
