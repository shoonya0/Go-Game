package core

import (
	"bytes"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// uiFont is the menus' pixel font, Press Start 2P (main-menu D1), parsed on
// first use. The fonts package's other font isn't referenced, so the linker
// leaves it out of the binary.
var uiFont *text.GoTextFaceSource

func uiFace(size float64) *text.GoTextFace {
	if uiFont == nil {
		src, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
		if err != nil {
			log.Fatalf("load UI font: %v", err)
		}
		uiFont = src
	}
	return &text.GoTextFace{Source: uiFont, Size: size}
}

// TextSize returns the width and height of s drawn at size px.
func TextSize(s string, size float64) (w, h float64) {
	return text.Measure(s, uiFace(size), 0)
}

// DrawTextCentered draws s at size px, centered on the screen point (cx, cy).
func DrawTextCentered(dst *ebiten.Image, s string, size, cx, cy float64, clr color.Color) {
	op := &text.DrawOptions{}
	op.PrimaryAlign, op.SecondaryAlign = text.AlignCenter, text.AlignCenter
	op.GeoM.Translate(cx, cy)
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(dst, s, uiFace(size), op)
}
