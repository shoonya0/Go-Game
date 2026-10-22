package core

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// HitboxKind says what a debug box represents; it picks the overlay colour.
type HitboxKind int

const (
	HitboxBody   HitboxKind = iota // collision body: collides with tiles and receives hits
	HitboxAttack                   // area that deals damage
)

// Hitbox is one box an entity reports to the debug overlay. Active marks an
// attack box that is striking right now; inactive ones are drawn fainter.
type Hitbox struct {
	Box    AABB
	Kind   HitboxKind
	Active bool
}

// HitboxProvider is anything that can show its boxes in the debug overlay.
// Enemies implement it; any other entity (e.g. the player) opts in by adding a
// Hitboxes method, with no changes here.
type HitboxProvider interface {
	Hitboxes() []Hitbox
}

// Debug overlay colours.
var (
	hitboxBodyColor     = color.RGBA{R: 40, G: 230, B: 80, A: 255}
	hitboxAttackColor   = color.RGBA{R: 240, G: 50, B: 50, A: 255}
	hitboxIdleColor     = color.NRGBA{R: 240, G: 50, B: 50, A: 110} // NRGBA: not premultiplied
	hitboxActiveFill    = color.NRGBA{R: 240, G: 50, B: 50, A: 70}
	hitboxStrokeWidth   = float32(2)
	hitboxIdleLineWidth = float32(1)
)

// DrawHitboxes outlines every provider's boxes, converting world coordinates to
// screen coordinates with the camera. It is a debug aid: call it only while the
// overlay is switched on.
func DrawHitboxes(screen *ebiten.Image, cam Camera, providers ...HitboxProvider) {
	for _, p := range providers {
		for _, hb := range p.Hitboxes() {
			drawHitbox(screen, cam, hb)
		}
	}
}

func drawHitbox(screen *ebiten.Image, cam Camera, hb Hitbox) {
	x := float32(hb.Box.X - cam.Pos.X)
	y := float32(hb.Box.Y - cam.Pos.Y)
	w, h := float32(hb.Box.Width), float32(hb.Box.Height)

	switch {
	case hb.Kind == HitboxBody:
		vector.StrokeRect(screen, x, y, w, h, hitboxStrokeWidth, hitboxBodyColor, false)
	case hb.Active:
		vector.FillRect(screen, x, y, w, h, hitboxActiveFill, false)
		vector.StrokeRect(screen, x, y, w, h, hitboxStrokeWidth, hitboxAttackColor, false)
	default:
		vector.StrokeRect(screen, x, y, w, h, hitboxIdleLineWidth, hitboxIdleColor, false)
	}
}
