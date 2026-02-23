package main

import (
	"fmt"
	"log"

	"player/enemy"
	"player/internal/core"
	"player/internal/system"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// playerSpriteSheetPath = "../assets/NewGideonGraves.png"
	playerSpriteSheetPath = "../assets/GideonGraves.png"
	screenWidth           = 1360
	screenHeight          = 768
	// screenWidth  = 1920
	// screenHeight = 1080
)

type Game struct {
	state core.GameState
	input core.InputState

	// ------ Entities ------
	player *core.PlayerRuntime

	ParallelEnemyManager *enemy.ParallelEnemyManager // enemy manager

	Background      *ebiten.Image
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
			if g.Level[i].TileInfo.TileType == core.EnemyBasic {
				// we are now registering the enemy basic tile to the quadtree
				continue
			}
			g.DynamicQuadtree.Insert(&g.Level[i])
		}
		g.ParallelEnemyManager.AddEnemyToLevel(g.Level)
		fmt.Println("Level loaded")
		doOnce = true
	}
}

// run 60 TPS
// run automatically every frame
func (g *Game) Update() error {

	g.loadLevel()

	// Update enemies
	g.ParallelEnemyManager.Update(g.player, g.DynamicQuadtree)

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
	// draw background
	g.player.DrawParallaxBackground(screen, g.Background, float64(screenWidth), float64(screenHeight))

	// draw level
	g.player.DrawLevel(screen, g.DynamicQuadtree, float64(screenWidth), float64(screenHeight), g.Tileset)

	// draw player animation
	g.player.DrawPlayerAnimation(screen)
	// draw UI
}

// run automatically every frame
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
	// return 640, 480
}

// run Once
func main() {
	backGroundData := core.LoadImage(core.Background_1)
	levelData := core.LoadImage(core.Level_1)
	tileData := core.LoadImage(core.Tileset)
	img := core.LoadImage(playerSpriteSheetPath)

	var parallelEnemyManager = enemy.DefaultParallelConfig()
	fmt.Println("Parallel Enemy Manager will create ", parallelEnemyManager.WorkerCount, "workers")

	game := &Game{
		state: core.ModeMenu,
		player: func() *core.PlayerRuntime {
			p := core.InitPlayer(img)
			return &p
		}(),

		ParallelEnemyManager: &parallelEnemyManager,

		Background:      backGroundData,
		LevelData:       levelData,
		Tileset:         tileData,
		score:           0,
		tickCount:       0,
		isDebug:         false,
		Level:           []core.Platform{},
		DynamicQuadtree: core.NewDynamicQuadtree(core.AABB{X: 0, Y: 0, Width: float64(core.Level_1_Width), Height: float64(core.Level_1_Height)}),
	}

	// ebiten.SetWindowSize(640, 480) // 640, 480
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pirate Adventure")

	// need to setup the brain for later

	if err := ebiten.RunGame(game); err != nil {
		if game.ParallelEnemyManager != nil {
			game.ParallelEnemyManager.Shutdown()
		}
		log.Fatal(err)
	}
}
