package core

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// HUD sprite sheets. Each holds three states of a single segment: full, half
// and empty.
const (
	HealthSheet  = "../assets/UI/health.png"
	StaminaSheet = "../assets/UI/stamina.png"
)

// HUD layout constants. Health and stamina (the Combat.Power pool) are each
// shown as a row of hudSegments sprites.
const (
	hudSegments = 10
	hudMargin   = 16
	hudRowGap   = 8

	hudHeartSize  = 28 // drawn heart width and height
	hudHeartGap   = 4
	hudStaminaW   = 30 // drawn stamina segment width
	hudStaminaH   = 10 // drawn stamina segment height
	hudStaminaGap = 2
)

// segmentState is the fill level of one HUD segment, indexing the sprite rects.
type segmentState int

const (
	segmentFull segmentState = iota
	segmentHalf
	segmentEmpty
)

// Source rects of the full/half/empty states inside each sheet, measured from
// the non-transparent pixels. Hearts sit side by side; stamina bars are stacked.
var (
	heartRects = [3]image.Rectangle{
		image.Rect(44, 36, 236, 228),
		image.Rect(271, 36, 463, 228),
		image.Rect(498, 36, 690, 228),
	}
	staminaRects = [3]image.Rectangle{
		image.Rect(12, 21, 651, 115),
		image.Rect(26, 152, 646, 253),
		image.Rect(25, 272, 644, 362),
	}
)

// hudSprite holds the three sub-images of one sheet, indexed by segmentState.
type hudSprite [3]*ebiten.Image

var hudHearts, hudStamina *hudSprite

func loadHUDSprite(sheetPath string, rects [3]image.Rectangle) *hudSprite {
	sheet := LoadImage(sheetPath)
	var s hudSprite
	for i, r := range rects {
		s[i] = sheet.SubImage(r).(*ebiten.Image)
	}
	return &s
}

// segmentStates splits ratio (clamped to [0, 1]) into n segments. Each segment
// holds two half-steps and partial values round up to the next half, so any
// non-zero amount still shows at least a half segment.
func segmentStates(ratio float64, n int) []segmentState {
	if ratio < 0 {
		ratio = 0
	} else if ratio > 1 {
		ratio = 1
	}
	return halfStepStates(int(math.Ceil(ratio*float64(2*n)-1e-9)), n)
}

// halfStepStates fills n segments with halves half-steps: two per full segment,
// one for the half segment after them, the rest empty.
func halfStepStates(halves, n int) []segmentState {
	states := make([]segmentState, n)
	for i := range states {
		switch {
		case halves >= 2*(i+1):
			states[i] = segmentFull
		case halves == 2*i+1:
			states[i] = segmentHalf
		default:
			states[i] = segmentEmpty
		}
	}
	return states
}

// DrawHUD renders the player's health hearts and stamina bars in the top-left
// corner.
func (player *PlayerRuntime) DrawHUD(screen *ebiten.Image) {
	player.DrawHUDAt(screen, hudMargin)
}

// DrawHUDAt renders the HUD with its left edge at x, leaving room for anything
// drawn in the corner (the touch menu button).
func (player *PlayerRuntime) DrawHUDAt(screen *ebiten.Image, x float64) {
	if hudHearts == nil {
		hudHearts = loadHUDSprite(HealthSheet, heartRects)
		hudStamina = loadHUDSprite(StaminaSheet, staminaRects)
	}

	healthRatio := 0.0
	if player.Combat.MaxHealth > 0 {
		healthRatio = float64(player.Combat.Health) / float64(player.Combat.MaxHealth)
	}
	staminaRatio := 0.0
	if player.Combat.MaxPower > 0 {
		staminaRatio = player.Combat.Power / player.Combat.MaxPower
	}

	drawSegments(screen, hudHearts, segmentStates(healthRatio, hudSegments),
		x, hudMargin, hudHeartSize, hudHeartSize, hudHeartGap)
	drawSegments(screen, hudStamina, segmentStates(staminaRatio, hudSegments),
		x, hudMargin+hudHeartSize+hudRowGap, hudStaminaW, hudStaminaH, hudStaminaGap)
}

// drawSegments draws one sprite per state in a row starting at (x, y), each
// scaled to w×h.
func drawSegments(screen *ebiten.Image, sprite *hudSprite, states []segmentState, x, y, w, h, gap float64) {
	for i, st := range states {
		img := sprite[st]
		b := img.Bounds()
		op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
		op.GeoM.Scale(w/float64(b.Dx()), h/float64(b.Dy()))
		op.GeoM.Translate(x+float64(i)*(w+gap), y)
		screen.DrawImage(img, op)
	}
}
