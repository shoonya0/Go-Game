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

	// ---------------- run ----------------
	inputState.RunJustPressed = ebiten.IsKeyPressed(ebiten.KeyShiftLeft)

	// ---------------- smug face ----------------
	inputState.SmugFace = ebiten.IsKeyPressed(ebiten.KeyF)

	// ---------------- menu ----------------
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		inputState.Menu = !inputState.Menu
	}

	// ---------------- skills ----------------
	inputState.Skills.WeakAttack = ebiten.IsKeyPressed(ebiten.KeyI)
	inputState.Skills.StrongAttack = ebiten.IsKeyPressed(ebiten.KeyJ)
	inputState.Skills.Defense = ebiten.IsKeyPressed(ebiten.KeyR)
	inputState.Skills.UsePotion = ebiten.IsKeyPressed(ebiten.KeyQ)
	inputState.Skills.SpecialAttack1 = ebiten.IsKeyPressed(ebiten.KeyK)
	inputState.Skills.SpecialAttack2 = ebiten.IsKeyPressed(ebiten.KeyL)
	inputState.Skills.SpecialAttack3 = ebiten.IsKeyPressed(ebiten.KeyU)
	inputState.Skills.SpecialAttack4 = ebiten.IsKeyPressed(ebiten.KeyO)
}
