package main

import (
	"log"

	"player/internal/core"
	"player/internal/system"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	playerSpriteSheetPath = "../assets/player/GideonGraves.png"
	screenWidth           = 1360
	screenHeight          = 768
)

// Game wires the engine's subsystems together and satisfies ebiten.Game.
type Game struct {
	state core.GameState
	input core.InputState
	touch *core.TouchControls // on-screen buttons for touch screens

	player  *core.PlayerRuntime
	enemies []*core.EnemyRuntime

	background *ebiten.Image
	tileset    *ebiten.Image
	level      core.Level
	quadtree   *core.DynamicQuadtree
	portals    *core.PortalSprite

	showHitboxes bool // debug overlay, toggled with F3
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

	touch := core.NewTouchControls(screenWidth, screenHeight)
	touch.Enabled = system.IsTouchDevice()

	return &Game{
		state:      core.ModeMenu,
		player:     core.InitPlayer(core.LoadImage(playerSpriteSheetPath)),
		enemies:    core.SpawnEnemies(level.Enemies),
		background: core.LoadImage(core.Background_1),
		tileset:    core.LoadImage(core.Tileset),
		level:      level,
		touch:      touch,
		quadtree:   quadtree,
		portals:    core.NewPortalSprite(),
	}
}

// Update advances the simulation one tick (runs at ~60 TPS).
func (g *Game) Update() error {
	system.HandleInput(&g.input)
	system.HandleTouch(&g.input, g.touch)
	if g.input.ToggleHitboxes {
		g.showHitboxes = !g.showHitboxes
	}

	core.UpdatePlayer(g.player, &g.input, g.quadtree, g.level.Portals)
	g.player.UpdateAnimation()
	g.player.HitEnemies(g.enemies)
	g.player.UpdateGroundWaves(g.quadtree, g.enemies)
	g.player.UpdateCamera(screenWidth, screenHeight, float64(core.Level_1_Width), float64(core.Level_1_Height))

	g.updateEnemies()
	g.portals.Update()

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
	g.portals.Draw(screen, g.level.Portals, g.player.Camera.Pos)
	for _, e := range g.enemies {
		e.DrawEnemyAnimation(screen, g.player.Camera)
		e.DrawHealthBar(screen, g.player.Camera)
	}
	g.player.DrawPlayerAnimation(screen)
	g.player.DrawGroundWaves(screen)
	// Water is drawn after the player so the player appears submerged in it.
	g.player.DrawWater(screen, g.quadtree, screenWidth, screenHeight, g.tileset)
	if g.showHitboxes {
		g.drawHitboxes(screen)
	}
	if g.touch.Enabled {
		g.player.DrawHUDAt(screen, core.TouchHUDLeft)
	} else {
		g.player.DrawHUD(screen)
	}
	g.touch.Draw(screen)
}

// drawHitboxes outlines the body and attack boxes of the enemies and the player.
func (g *Game) drawHitboxes(screen *ebiten.Image) {
	providers := make([]core.HitboxProvider, 0, len(g.enemies)+len(g.level.Portals)+1)
	for _, p := range g.level.Portals {
		providers = append(providers, p)
	}
	for _, e := range g.enemies {
		providers = append(providers, e)
	}
	providers = append(providers, g.player)
	core.DrawHitboxes(screen, g.player.Camera, providers...)
}

// Layout maps the outside window size to the game's logical screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pirate Adventure")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
