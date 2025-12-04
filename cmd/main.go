package main

import (
	"image/color"
	"log"
	"player/internal/core"
	"player/internal/system"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	state core.GameState
	input core.InputState

	spriteSheet *ebiten.Image // player image
	// ------ Entities ------
	player *core.PlayerRuntime

	Level           []core.Platform
	DynamicQuadtree *core.DynamicQuadtree

	// Meta Data
	score     int
	tickCount int
	isDebug   bool
}

// var frame = 0

// run 60 TPS
// run automatically every frame
func (g *Game) Update() error {
	// poll input -> call another function to handle input
	system.HandleInput(&g.input)

	core.UpdatePlayer(g.player, &g.input, g.DynamicQuadtree)

	g.player.UpdateAnimation()

	// update player animation
	// check for collisions
	// update camera position
	// render game state
	return nil
}

// run automatically every frame
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw platforms
	for _, p := range g.Level {
		vector.DrawFilledRect(screen, float32(p.X), float32(p.Y), float32(p.Width), float32(p.Height), color.RGBA{100, 100, 100, 255}, false)
	}

	g.player.DrawAnimation(screen)

	// draw UI
}

// run automatically every frame
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
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
			// Set initial position above ground
			p.Pos.X = 100
			p.Pos.Y = 100
			return &p
		}(),
		spriteSheet: img,
		score:       0,
		tickCount:   0,
		isDebug:     false,
		Level: []core.Platform{
			// Floor
			{X: 0, Y: 400, Width: 640, Height: 80},
			// Floating Platform
			{X: 200, Y: 300, Width: 150, Height: 20},
			{X: 400, Y: 200, Width: 100, Height: 20},
		},
		DynamicQuadtree: core.NewDynamicQuadtree(core.AABB{X: 0, Y: 0, Width: core.Level_1_Width, Height: core.Level_1_Height}),
	}

	// Add platforms to Quadtree
	for i := range game.Level {
		game.DynamicQuadtree.Insert(&game.Level[i])
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Pirate Adventure")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
