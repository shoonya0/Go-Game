package main

import (
	"log"

	"player/internal/core"
	"player/internal/system"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	playerSpriteSheetPath = "../assets/GideonGraves.png"
	screenWidth           = 1360
	screenHeight          = 768
)

// Game wires the engine's subsystems together and satisfies ebiten.Game.
type Game struct {
	state core.GameState
	input core.InputState

	player  *core.PlayerRuntime
	enemies []*core.EnemyRuntime

	background *ebiten.Image
	tileset    *ebiten.Image
	level      core.Level
	quadtree   *core.DynamicQuadtree
}

// NewGame loads assets, builds the level and returns a ready-to-run Game.
func NewGame() *Game {
	core.WorldInit()

	levelBounds := core.AABB{Width: float64(core.Level_1_Width), Height: float64(core.Level_1_Height)}
	quadtree := core.NewDynamicQuadtree(levelBounds)

	level := core.LoadLevel(core.LevelFile, core.DefaultLevelName)
	for i := range level.Platforms {
		quadtree.Insert(&level.Platforms[i])
	}

	return &Game{
		state:  core.ModeMenu,
		player: core.InitPlayer(core.LoadImage(playerSpriteSheetPath)),
		enemies: []*core.EnemyRuntime{
			core.NewEnemy(core.KindSuccubus, core.Position{X: 700, Y: 1000}),
		},
		background: core.LoadImage(core.Background_1),
		tileset:    core.LoadImage(core.Tileset),
		level:      level,
		quadtree:   quadtree,
	}
}

// Update advances the simulation one tick (runs at ~60 TPS).
func (g *Game) Update() error {
	system.HandleInput(&g.input)

	core.UpdatePlayer(g.player, &g.input, g.quadtree, g.level.Checkpoints)
	g.player.UpdateAnimation()
	g.player.UpdateCamera(screenWidth, screenHeight, float64(core.Level_1_Width), float64(core.Level_1_Height))

	g.updateEnemies()

	return nil
}

// updateEnemies advances every live enemy and culls those that have finished
// dying.
func (g *Game) updateEnemies() {
	live := g.enemies[:0]
	for _, e := range g.enemies {
		core.UpdateEnemy(e, g.player, g.quadtree)
		e.UpdateEnemyAnimation()
		if !e.Removable() {
			live = append(live, e)
		}
	}
	g.enemies = live
}

// Draw renders one frame: background, level, then player.
func (g *Game) Draw(screen *ebiten.Image) {
	g.player.DrawParallaxBackground(screen, g.background, screenWidth, screenHeight)
	g.player.DrawLevel(screen, g.quadtree, screenWidth, screenHeight, g.tileset)
	for _, e := range g.enemies {
		e.DrawEnemyAnimation(screen, g.player.Camera)
	}
	g.player.DrawPlayerAnimation(screen)
	// Water is drawn after the player so the player appears submerged in it.
	g.player.DrawWater(screen, g.quadtree, screenWidth, screenHeight, g.tileset)
	g.player.DrawHUD(screen)
}

// Layout maps the outside window size to the game's logical screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pirate Adventure")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
