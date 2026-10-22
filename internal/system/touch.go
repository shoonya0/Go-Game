package system

import (
	"image"

	"player/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
)

// HandleTouch feeds this frame's touch points to the on-screen controls and
// merges the pressed buttons into inputState. Call it after HandleInput.
func HandleTouch(inputState *core.InputState, touch *core.TouchControls) {
	ids := ebiten.AppendTouchIDs(nil)
	points := make([]image.Point, 0, len(ids))
	for _, id := range ids {
		x, y := ebiten.TouchPosition(id)
		points = append(points, image.Pt(x, y))
	}
	touch.Update(points)
	touch.Apply(inputState)
}
