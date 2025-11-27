package system

import (
	"player/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func HandleInput(inputState *core.InputState) {
	// poll input from keyboard and mouse
	// player movement
	w := ebiten.IsKeyPressed(ebiten.KeyW)
	s := ebiten.IsKeyPressed(ebiten.KeyS)
	a := ebiten.IsKeyPressed(ebiten.KeyA)
	d := ebiten.IsKeyPressed(ebiten.KeyD)

	// ---------------- up down ----------------
	switch {
	case w && !s:
		inputState.Direction.UpDown = 1
	case s && !w:
		inputState.Direction.UpDown = -1
	default:
		inputState.Direction.UpDown = 0
	}

	// ---------------- left right ----------------
	switch {
	case a && !d:
		inputState.Direction.LeftRight = -1
	case d && !a:
		inputState.Direction.LeftRight = 1
	default:
		inputState.Direction.LeftRight = 0
	}

	// ---------------- jump ----------------
	inputState.JumpJustPressed = ebiten.IsKeyPressed(ebiten.KeySpace)

	// ---------------- dash ----------------
	inputState.DashJustPressed = ebiten.IsKeyPressed(ebiten.KeyShiftLeft)

	// ---------------- menu ----------------
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		inputState.Menu = !inputState.Menu
	}

	// ---------------- skills ----------------
	inputState.Skills.WeakAttack = ebiten.IsKeyPressed(ebiten.KeyJ)
	inputState.Skills.StrongAttack = ebiten.IsKeyPressed(ebiten.KeyK)
	inputState.Skills.Defense = ebiten.IsKeyPressed(ebiten.KeyL)
	inputState.Skills.UsePotion = ebiten.IsKeyPressed(ebiten.KeyI)
}
