package system

import (
	"image"

	"player/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// HandleTouch feeds this frame's touch points to the on-screen controls and
// merges the pressed buttons into inputState. Call it after HandleInput.
func HandleTouch(inputState *core.InputState, touch *core.TouchControls) {
	touch.Update(touchPoints(ebiten.AppendTouchIDs(nil)))
	touch.Started(touchPoints(inpututil.AppendJustPressedTouchIDs(nil)))
	touch.Apply(inputState)
}

func touchPoints(ids []ebiten.TouchID) []image.Point {
	points := make([]image.Point, 0, len(ids))
	for _, id := range ids {
		x, y := ebiten.TouchPosition(id)
		points = append(points, image.Pt(x, y))
	}
	return points
}
