package core

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// UIButtonsSheet holds the wooden UI the game's buttons, menus and bars are
// drawn from (ui-buttons feature).
const UIButtonsSheet = "../assets/UI/UIButtons.png"

// Source rects inside UIButtonsSheet, measured from the non-transparent pixels.
var (
	uiMenuRect      = image.Rect(125, 1430, 329, 1641)   // square button with the menu icon
	uiRoundUpRect   = image.Rect(2052, 1732, 2268, 1946) // round gold button
	uiRoundDownRect = image.Rect(1811, 1732, 2027, 1946) // round dark button, used while pressed

	uiPlankRect      = image.Rect(908, 1093, 1561, 1331)  // plain plank: a menu item
	uiPlankArrowRect = image.Rect(1628, 1092, 2279, 1332) // plank pointing right: the picked item
	uiBoardRect      = image.Rect(1148, 162, 2282, 998)   // the board behind a menu
	uiBarFullRect    = image.Rect(138, 179, 1099, 301)    // bar with all 5 segments
	uiBarEmptyRect   = image.Rect(138, 879, 1099, 1001)   // the same bar empty

	uiRestartRect  = image.Rect(365, 1732, 582, 1946) // round buttons with an icon
	uiSettingsRect = image.Rect(606, 1732, 822, 1946)
	uiHomeRect     = image.Rect(1089, 1732, 1304, 1946)
	uiSoundRect    = image.Rect(1329, 1732, 1545, 1946)
	uiMuteRect     = image.Rect(1570, 1732, 1786, 1946)
	uiLeftRect     = image.Rect(688, 2034, 874, 2255) // wooden arrows
	uiRightRect    = image.Rect(901, 2034, 1094, 2255)
	uiCrossRect    = image.Rect(1123, 2047, 1310, 2241)
	uiCheckRect    = image.Rect(1329, 2034, 1518, 2233)
)

// uiInsets are the sheet pixels at each side of a sprite that keep their
// shape when it is stretched: only the part between them stretches.
type uiInsets struct{ left, top, right, bottom int }

var (
	uiPlankInsets      = uiInsets{left: 90, right: 90}
	uiPlankArrowInsets = uiInsets{left: 90, right: 151} // the whole point
	uiBoardInsets      = uiInsets{110, 110, 110, 110}
)

// uiBarSegments is where the bar's segments start and end on the sheet: a
// fill is the full bar cut off between them.
const uiBarSegLeft, uiBarSegRight = 168, 1073

var uiSheet *ebiten.Image

// uiSprite returns sprite r of UIButtonsSheet, loading the sheet on first use.
func uiSprite(r image.Rectangle) *ebiten.Image {
	if uiSheet == nil {
		uiSheet = LoadImage(UIButtonsSheet)
	}
	return uiSheet.SubImage(r).(*ebiten.Image)
}

// sliceRects cuts src into up to 9 pieces at insets in and returns where each
// piece lands in dst: the sides are drawn at scale (screen px per sheet px),
// the middle stretches to fill. Pieces of no size are left out. When dst is
// too small for the sides, they meet in its middle.
func sliceRects(src image.Rectangle, in uiInsets, dst image.Rectangle, scale float64) (from, to []image.Rectangle) {
	cuts := func(lo, hi, a, b int, dlo, dhi int) (s, d [4]int) {
		da := int(math.Round(float64(a) * scale))
		db := int(math.Round(float64(b) * scale))
		if da+db > dhi-dlo { // too small: the sides share what there is
			da = (dhi - dlo) * a / max(1, a+b)
			db = dhi - dlo - da
		}
		return [4]int{lo, lo + a, hi - b, hi}, [4]int{dlo, dlo + da, dhi - db, dhi}
	}
	sx, dx := cuts(src.Min.X, src.Max.X, in.left, in.right, dst.Min.X, dst.Max.X)
	sy, dy := cuts(src.Min.Y, src.Max.Y, in.top, in.bottom, dst.Min.Y, dst.Max.Y)
	for j := range 3 {
		for i := range 3 {
			f := image.Rect(sx[i], sy[j], sx[i+1], sy[j+1])
			t := image.Rect(dx[i], dy[j], dx[i+1], dy[j+1])
			if f.Empty() || t.Empty() {
				continue
			}
			from, to = append(from, f), append(to, t)
		}
	}
	return from, to
}

// drawSliced draws sprite src stretched to dst, its sides kept at scale (see
// sliceRects), tinted by tint when it is not nil.
func drawSliced(screen *ebiten.Image, src image.Rectangle, in uiInsets, dst image.Rectangle, scale float64, tint color.Color, alpha float32) {
	from, to := sliceRects(src, in, dst, scale)
	for k := range from {
		drawScaled(screen, uiSprite(from[k]), to[k], tint, alpha)
	}
}

// drawPlank draws a plank (or the bar) as tall as dst, stretched to its width.
func drawPlank(screen *ebiten.Image, src image.Rectangle, in uiInsets, dst image.Rectangle, tint color.Color) {
	drawSliced(screen, src, in, dst, float64(dst.Dy())/float64(src.Dy()), tint, 1)
}

// barFillRect is the part of the full bar drawn for fill (0..1): the bar's
// left end and the segments up to fill of their length.
func barFillRect(fill float64) image.Rectangle {
	r := uiBarFullRect
	if fill >= 1 {
		return r
	}
	fill = max(0, fill)
	r.Max.X = uiBarSegLeft + int(math.Round(fill*(uiBarSegRight-uiBarSegLeft)))
	return r
}

// drawRotated draws img fitted into the square dst, turned by quarter turns
// clockwise.
func drawRotated(screen, img *ebiten.Image, dst image.Rectangle, quarters int, tint color.Color, alpha float32) {
	b := img.Bounds()
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Translate(-float64(b.Dx())/2, -float64(b.Dy())/2)
	op.GeoM.Rotate(float64(quarters) * math.Pi / 2)
	s := float64(dst.Dx()) / float64(max(b.Dx(), b.Dy()))
	op.GeoM.Scale(s, s)
	op.GeoM.Translate(float64(dst.Min.X+dst.Max.X)/2, float64(dst.Min.Y+dst.Max.Y)/2)
	if tint != nil {
		op.ColorScale.ScaleWithColor(tint)
	}
	op.ColorScale.ScaleAlpha(alpha)
	screen.DrawImage(img, op)
}
