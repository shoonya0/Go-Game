package system

import (
	"player/internal/core"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func HandleInput(inputState *core.InputState) {
	// poll input from keyboard and mouse
	// player movement

	num := inpututil.AppendPressedKeys([]ebiten.Key{})

	// ---------------- up down ----------------
	x := 0
	if slices.Contains(num, ebiten.KeyW) {
		inputState.Direction.UpDown = 1
		x++
	}
	if slices.Contains(num, ebiten.KeyS) {
		inputState.Direction.UpDown = -1
		x++
	}
	if x&1 == 0 {
		inputState.Direction.UpDown = 0
	}

	x = 0
	if slices.Contains(num, ebiten.KeyA) {
		inputState.Direction.LeftRight = -1
		x++
	}
	if slices.Contains(num, ebiten.KeyD) {
		inputState.Direction.LeftRight = 1
		x++
	}
	if x&1 == 0 {
		inputState.Direction.LeftRight = 0
	}

	// ---------------- jump ----------------
	if slices.Contains(num, ebiten.KeySpace) {
		inputState.JumpJustPressed = true
	} else {
		inputState.JumpJustPressed = false
	}

	// ---------------- dash ----------------
	if slices.Contains(num, ebiten.KeyShiftLeft) {
		inputState.DashJustPressed = true
	} else {
		inputState.DashJustPressed = false
	}

	// ---------------- menu ----------------
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		inputState.Menu = !inputState.Menu
	}

	// ---------------- skills ----------------
	if slices.Contains(num, ebiten.KeyJ) {
		inputState.Skills.WeakAttack = true
	} else {
		inputState.Skills.WeakAttack = false
	}

	if slices.Contains(num, ebiten.KeyK) {
		inputState.Skills.StrongAttack = true
	} else {
		inputState.Skills.StrongAttack = false
	}
	if slices.Contains(num, ebiten.KeyL) {
		inputState.Skills.Defense = true
	} else {
		inputState.Skills.Defense = false
	}
	if slices.Contains(num, ebiten.KeyI) {
		inputState.Skills.UsePotion = true
	} else {
		inputState.Skills.UsePotion = false
	}
}
