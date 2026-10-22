package core

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// touchIconSize is the side of the square image each button symbol is drawn
// into; Draw scales it onto the button.
const touchIconSize = 128

// iconShape is one part of a button symbol in unit coordinates (0..1, y down).
// It is filled as a polygon, or stroked as a line of that width when width > 0.
type iconShape struct {
	pts   [][2]float32
	width float32
}

func fillShape(pts ...[2]float32) iconShape { return iconShape{pts: pts} }

func lineShape(width float32, pts ...[2]float32) iconShape {
	return iconShape{pts: pts, width: width}
}

// starPoints returns a star of n tips around (cx, cy), alternating between the
// outer and inner radius, starting with a tip straight up.
func starPoints(n int, outer, inner, cx, cy float64) [][2]float32 {
	pts := make([][2]float32, 0, 2*n)
	for i := 0; i < 2*n; i++ {
		r := outer
		if i%2 == 1 {
			r = inner
		}
		a := float64(i)*math.Pi/float64(n) - math.Pi/2
		pts = append(pts, [2]float32{float32(cx + r*math.Cos(a)), float32(cy + r*math.Sin(a))})
	}
	return pts
}

// wavePoints samples one and a half sine waves across the button, for U's
// ground wave.
func wavePoints() [][2]float32 {
	const steps = 24
	pts := make([][2]float32, 0, steps+1)
	for i := 0; i <= steps; i++ {
		t := float64(i) / steps
		pts = append(pts, [2]float32{float32(0.2 + 0.6*t), float32(0.5 - 0.1*math.Sin(t*3*math.Pi))})
	}
	return pts
}

// touchIconShapes are the button symbols: arrows to move, and for attacks a
// thin slash (weak), two thick slashes (heavy), crossed slashes (K combo), a
// burst (L, the strongest combo), a wave over the ground (U) and a sparkle (O).
var touchIconShapes = [touchButtonCount][]iconShape{
	TouchUp:    {fillShape([2]float32{0.5, 0.24}, [2]float32{0.78, 0.68}, [2]float32{0.22, 0.68})},
	TouchDown:  {fillShape([2]float32{0.5, 0.76}, [2]float32{0.78, 0.32}, [2]float32{0.22, 0.32})},
	TouchLeft:  {fillShape([2]float32{0.24, 0.5}, [2]float32{0.68, 0.22}, [2]float32{0.68, 0.78})},
	TouchRight: {fillShape([2]float32{0.76, 0.5}, [2]float32{0.32, 0.22}, [2]float32{0.32, 0.78})},

	TouchWeak: {lineShape(0.07, [2]float32{0.32, 0.72}, [2]float32{0.68, 0.28})},
	TouchStrong: {
		lineShape(0.12, [2]float32{0.24, 0.66}, [2]float32{0.56, 0.26}),
		lineShape(0.12, [2]float32{0.44, 0.76}, [2]float32{0.76, 0.36}),
	},
	TouchSpecial1: {
		lineShape(0.09, [2]float32{0.28, 0.28}, [2]float32{0.72, 0.72}),
		lineShape(0.09, [2]float32{0.72, 0.28}, [2]float32{0.28, 0.72}),
	},
	TouchSpecial2: {fillShape(starPoints(8, 0.3, 0.14, 0.5, 0.5)...)},
	TouchSpecial3: {
		lineShape(0.08, wavePoints()...),
		lineShape(0.06, [2]float32{0.2, 0.7}, [2]float32{0.8, 0.7}),
	},
	TouchSpecial4: {fillShape(starPoints(4, 0.3, 0.09, 0.5, 0.5)...)},
}

// newTouchIcon draws shapes in white into a new image, ready to be tinted.
func newTouchIcon(shapes []iconShape) *ebiten.Image {
	img := ebiten.NewImage(touchIconSize, touchIconSize)
	const s = touchIconSize
	for _, sh := range shapes {
		var p vector.Path
		for i, pt := range sh.pts {
			if i == 0 {
				p.MoveTo(pt[0]*s, pt[1]*s)
			} else {
				p.LineTo(pt[0]*s, pt[1]*s)
			}
		}
		draw := &vector.DrawPathOptions{AntiAlias: true}
		if sh.width > 0 {
			stroke := &vector.StrokeOptions{Width: sh.width * s, LineCap: vector.LineCapRound, LineJoin: vector.LineJoinRound}
			vector.StrokePath(img, &p, stroke, draw)
		} else {
			p.Close()
			vector.FillPath(img, &p, nil, draw)
		}
	}
	return img
}
