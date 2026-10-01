package core

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// HUD layout constants.
const (
	hudMargin    = 16
	hudBarWidth  = 220
	hudBarHeight = 18
	hudBarGap    = 8
)

// HUD bar colors.
var (
	hudHealthColor = color.RGBA{R: 210, G: 50, B: 50, A: 255}
	hudPowerColor  = color.RGBA{R: 60, G: 130, B: 230, A: 255}
	hudTrackColor  = color.RGBA{R: 25, G: 25, B: 30, A: 210}
	hudBorderColor = color.RGBA{R: 240, G: 240, B: 240, A: 255}
)

// DrawHUD renders the player's health and power bars in the top-left corner.
func (player *PlayerRuntime) DrawHUD(screen *ebiten.Image) {
	healthRatio := 0.0
	if player.Combat.MaxHealth > 0 {
		healthRatio = float64(player.Combat.Health) / float64(player.Combat.MaxHealth)
	}
	powerRatio := 0.0
	if player.Combat.MaxPower > 0 {
		powerRatio = player.Combat.Power / player.Combat.MaxPower
	}

	drawBar(screen, hudMargin, hudMargin, healthRatio, hudHealthColor)
	drawBar(screen, hudMargin, hudMargin+hudBarHeight+hudBarGap, powerRatio, hudPowerColor)
}

// drawBar draws a bordered progress bar filled to ratio (clamped to [0, 1]).
func drawBar(screen *ebiten.Image, x, y float32, ratio float64, fill color.RGBA) {
	if ratio < 0 {
		ratio = 0
	} else if ratio > 1 {
		ratio = 1
	}

	vector.DrawFilledRect(screen, x, y, hudBarWidth, hudBarHeight, hudTrackColor, false)
	vector.DrawFilledRect(screen, x, y, float32(hudBarWidth*ratio), hudBarHeight, fill, false)
	vector.StrokeRect(screen, x, y, hudBarWidth, hudBarHeight, 2, hudBorderColor, false)
}
