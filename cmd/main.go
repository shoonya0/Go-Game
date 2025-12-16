package main

import (
	"fmt"
	"log"
	"player/internal/core"
	"player/internal/system"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	playerSpriteSheetPath = "../assets/GideonGraves.png"
	screenWidth           = 1360
	screenHeight          = 768
	// screenWidth  = 1920
	// screenHeight = 1080
)

type Game struct {
	state core.GameState
	input core.InputState

	spriteSheet *ebiten.Image // player image
	// ------ Entities ------
	player *core.PlayerRuntime

	LevelData       *ebiten.Image
	Tileset         *ebiten.Image
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
		core.WorldInit()
		g.Level = g.player.LoadLevel(g.LevelData)
		for i := range g.Level {
			g.DynamicQuadtree.Insert(&g.Level[i])
		}
		fmt.Println("Level loaded")
		doOnce = true
	}
}

// run 60 TPS
// run automatically every frame
func (g *Game) Update() error {

	g.loadLevel()

	// poll input -> call another function to handle input
	system.HandleInput(&g.input)

	core.UpdatePlayer(g.player, &g.input, g.DynamicQuadtree)

	g.player.UpdateAnimation()

	// update camera position
	g.player.UpdateCamera(float64(screenWidth), float64(screenHeight), float64(core.Level_1_Width), float64(core.Level_1_Height))

	// render game state
	return nil
}

// run automatically every frame
func (g *Game) Draw(screen *ebiten.Image) {
	g.player.DrawLevel(screen, g.DynamicQuadtree, float64(screenWidth), float64(screenHeight), g.Tileset)

	g.player.DrawPlayerAnimation(screen)
	// draw UI
}

// run automatically every frame
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func LoadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

// run Once
func main() {
	img := LoadImage(playerSpriteSheetPath)
	levelData := LoadImage(core.Level_1)
	TileData := LoadImage(core.Tileset)

	game := &Game{
		state: core.ModeMenu,
		player: func() *core.PlayerRuntime {
			p := core.InitPlayer()
			return &p
		}(),
		spriteSheet:     img,
		LevelData:       levelData,
		Tileset:         TileData,
		score:           0,
		tickCount:       0,
		isDebug:         false,
		Level:           []core.Platform{},
		DynamicQuadtree: core.NewDynamicQuadtree(core.AABB{X: 0, Y: 0, Width: float64(core.Level_1_Width), Height: float64(core.Level_1_Height)}),
	}

	// ebiten.SetWindowSize(640, 480) // 640, 480
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pirate Adventure")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
