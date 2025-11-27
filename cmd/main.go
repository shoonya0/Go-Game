package main

import (
	"image"
	"log"
	"player/internal/core"
	"player/internal/system"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// image dimensions
const (
	frameWidth  = 40
	frameHeight = 80
	numFrames   = 6
)

type Game struct {
	state core.GameState
	input core.InputState

	spriteSheet *ebiten.Image // player image
	// ------ Entities ------
	player *core.PlayerRuntime

	// Meta Data
	score     int
	tickCount int
	isDebug   bool
}

var num = "HELL YAAA"
var frame = 0

// run 60 TPS
// run automatically every frame
func (g *Game) Update() error {
	// poll input -> call another function to handle input
	system.HandleInput(&g.input)

	core.UpdatePlayer(g.player, &g.input)

	g.tickCount++
	if g.tickCount >= 60 {
		g.tickCount = 0
	}
	if g.tickCount%8 == 0 {
		frame++
		frame %= 6
	}

	// update player animation
	// check for collisions
	// update camera position
	// render game state
	return nil
}

// run automatically every frame
func (g *Game) Draw(screen *ebiten.Image) {
	// draw game state

	sx := frame * 40
	sy := 80 // 0 = idle

	rect := image.Rect(sx, sy, sx+40, sy+80)

	subImage := g.spriteSheet.SubImage(rect).(*ebiten.Image)

	// op is the options for the image uses the player's position
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.player.Pos.X, g.player.Pos.Y)
	op.GeoM.Scale(1, 1)

	screen.DrawImage(subImage, op)

	// draw UI
	// num = fmt.Sprintf("Direction: %v", g.input.Direction)
	// ebitenutil.DebugPrint(screen, num)

	// num = fmt.Sprintf("Jump Just Pressed: %v", g.input.JumpJustPressed)
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n%s", num))

	// num = fmt.Sprintf("Menu: %v", g.input.Menu)
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n%s", num))

	// num = fmt.Sprintf("Dash Just Pressed: %v", g.input.DashJustPressed)
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n\n%s", num))

	// num = fmt.Sprintf("Weak Attack: %v", g.input.Skills.WeakAttack)
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n\n\n%s", num))

	// num = fmt.Sprintf("Strong Attack: %v", g.input.Skills.StrongAttack)
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n\n\n\n%s", num))

	// num = fmt.Sprintf("Defense: %v", g.input.Skills.Defense)
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n\n\n\n\n%s", num))

	// num = fmt.Sprintf("Use Potion: %v", g.input.Skills.UsePotion)
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n\n\n\n\n\n%s", num))
}

// run automatically every frame
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

// run Once
func main() {

	img, _, err := ebitenutil.NewImageFromFile("GideonGraves.png")
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{
		state: core.ModeMenu,
		player: func() *core.PlayerRuntime {
			p := core.InitPlayer()
			return &p
		}(),
		spriteSheet: img,
		score:       0,
		tickCount:   0,
		isDebug:     false,
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Pirate Adventure")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
