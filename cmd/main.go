package main

import (
	"fmt"
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

	LevelData       *ebiten.Image
	Level           []core.Platform
	DynamicQuadtree *core.DynamicQuadtree

	// Meta Data
	score     int
	tickCount int
	isDebug   bool
}

var doOnce = false

func (g *Game) loadLevel() {
	if !doOnce {
		Level := core.LoadLevel(g.LevelData)
		g.Level = Level
		for i := range g.Level {
			g.DynamicQuadtree.Insert(&g.Level[i])
		}
		fmt.Println("Level loaded")
		doOnce = true
	}
}

// var frame = 0

// run 60 TPS
// run automatically every frame
func (g *Game) Update() error {

	g.loadLevel()

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
	// return 640, 480
	return 1360, 768
}

func LoadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

const (
	playerSpriteSheetPath = "../assets/GideonGraves.png"
	screenWidth           = 1360
	screenHeight          = 768
)

// run Once
func main() {
	img := LoadImage(playerSpriteSheetPath)
	levelData := LoadImage(core.Level_1)

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
		LevelData:   levelData,
		score:       0,
		tickCount:   0,
		isDebug:     false,
		Level:       []core.Platform{},
		// DynamicQuadtree: nil,
		DynamicQuadtree: core.NewDynamicQuadtree(core.AABB{X: 0, Y: 0, Width: float64(core.Level_1_Width), Height: float64(core.Level_1_Height)}),
	}

	// Add platforms to Quadtree
	// for i := range game.Level {
	// 	game.DynamicQuadtree.Insert(&game.Level[i])
	// }

	// ebiten.SetWindowSize(640, 480) // 640, 480
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pirate Adventure")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
