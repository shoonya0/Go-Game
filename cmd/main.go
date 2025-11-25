package main

import (
	"fmt"
	"log"
	"player/internal/core"
	"player/internal/system"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	Input core.InputState
}

var num = "HELL YAAA"

// run 60 TPS
// run automatically every frame
func (g *Game) Update() error {
	// poll input -> call another function to handle input
	system.HandleInput(&g.Input)

	// update player state
	// check for collisions
	// update camera position
	// render game state
	return nil
}

// run automatically every frame
func (g *Game) Draw(screen *ebiten.Image) {
	// draw game state
	// draw UI
	num = fmt.Sprintf("Direction: %v", g.Input.Direction)
	ebitenutil.DebugPrint(screen, num)

	num = fmt.Sprintf("Jump Just Pressed: %v", g.Input.JumpJustPressed)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n%s", num))

	num = fmt.Sprintf("Menu: %v", g.Input.Menu)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n%s", num))

	num = fmt.Sprintf("Dash Just Pressed: %v", g.Input.DashJustPressed)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n\n%s", num))

	num = fmt.Sprintf("Weak Attack: %v", g.Input.Skills.WeakAttack)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n\n\n%s", num))

	num = fmt.Sprintf("Strong Attack: %v", g.Input.Skills.StrongAttack)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n\n\n\n%s", num))

	num = fmt.Sprintf("Defense: %v", g.Input.Skills.Defense)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n\n\n\n\n%s", num))

	num = fmt.Sprintf("Use Potion: %v", g.Input.Skills.UsePotion)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n\n\n\n\n\n%s", num))
}

// run automatically every frame
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

// run Once
func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Pirate Adventure")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
