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

	player   *core.PlayerRuntime
	world    *core.World // the current level: platforms, portals, enemies
	respawns int         // player respawns already handled (enemies revived)

	background     *ebiten.Image
	backgroundPath string // the asset background was loaded from
	tileset        *ebiten.Image
	portals        *core.PortalSprite
	travel         *core.TravelMenu // opened by holding E (or a touch) at a portal
	view           *ebiten.Image    // scratch screen for drawing a portal's thumbnail

	showHitboxes bool // debug overlay, toggled with F3
}

// NewGame loads assets, builds the level and returns a ready-to-run Game.
func NewGame() *Game {
	core.WorldInit()

	world, err := core.NewWorld(core.StartLevel)
	if err != nil {
		log.Fatal(err)
	}

	touch := core.NewTouchControls(screenWidth, screenHeight)
	touch.Enabled = system.IsTouchDevice()

	player := core.InitPlayer(core.LoadImage(playerSpriteSheetPath))
	player.StartAt(world.Level)

	g := &Game{
		state:   core.ModeMenu,
		player:  player,
		world:   world,
		tileset: core.LoadImage(core.Tileset),
		touch:   touch,
		portals: core.NewPortalSprite(),
		travel:  core.NewTravelMenu(screenWidth, screenHeight),
		view:    ebiten.NewImage(screenWidth, screenHeight),
	}
	g.syncBackground()
	return g
}

// syncBackground loads the current level's background when it changed.
func (g *Game) syncBackground() {
	if path := g.world.Def.Background; path != g.backgroundPath {
		g.background = core.LoadImage(path)
		g.backgroundPath = path
	}
}

// Update advances the simulation one tick (runs at ~60 TPS).
func (g *Game) Update() error {
	system.HandleInput(&g.input)
	system.HandleTouch(&g.input, g.touch)
	if g.input.ToggleHitboxes {
		g.showHitboxes = !g.showHitboxes
	}
	if g.travel.Open { // the game is paused while the travel menu is open
		if ref, ok := g.travel.Update(&g.input); ok {
			if err := g.world.Travel(g.player, ref); err != nil {
				log.Printf("travel: %v", err)
			}
			g.syncBackground()
		}
		g.portals.Update()
		return nil
	}

	core.UpdatePlayer(g.player, &g.input, g.world.Quadtree, g.world.Level.Portals)
	g.player.UpdateAnimation()
	g.player.HitEnemies(g.world.Enemies)
	g.player.UpdateGroundWaves(g.world.Quadtree, g.world.Enemies)
	g.player.UpdateCamera(screenWidth, screenHeight, g.world.Def.Width, g.world.Def.Height)

	if g.player.Respawns != g.respawns {
		g.respawns = g.player.Respawns
		g.world.Respawned(g.player)
		g.syncBackground()
	}
	g.updateEnemies()
	g.portals.Update()
	g.captureThumbs()
	if g.player.TravelAsked() {
		g.travel.OpenFor(g.player)
	}

	return nil
}

// updateEnemies advances every live enemy and culls those that have finished
// dying.
func (g *Game) updateEnemies() {
	live := g.world.Enemies[:0]
	for _, e := range g.world.Enemies {
		core.UpdateEnemy(e, g.player, g.world.Quadtree)
		e.UpdateEnemyAnimation()
		if !e.Removable() {
			live = append(live, e)
		}
	}
	g.world.Enemies = live
}

// captureThumbs gives each portal that just became active its thumbnail for
// the travel menu: the level around it, drawn once into the scratch screen.
func (g *Game) captureThumbs() {
	for _, c := range g.world.Level.Portals {
		ref := c.Ref()
		if g.player.PortalLook(c) == core.PortalInactive || g.travel.HasThumb(ref) {
			continue
		}
		cam := g.player.Camera
		b := c.Bounds()
		g.player.Camera.Pos = core.Position{X: b.X + b.Width/2 - screenWidth/2, Y: b.Y + b.Height/2 - screenHeight/2}
		g.view.Clear()
		g.player.DrawParallaxBackground(g.view, g.background, screenWidth, screenHeight, g.world.Def.Width)
		g.player.DrawLevel(g.view, g.world.Quadtree, screenWidth, screenHeight, g.tileset)
		g.portals.Draw(g.view, []core.Portal{c}, g.player)
		g.player.DrawWater(g.view, g.world.Quadtree, screenWidth, screenHeight, g.tileset)
		g.player.Camera = cam
		g.travel.SetThumb(ref, g.view)
	}
}

// Draw renders one frame: background, level, then player.
func (g *Game) Draw(screen *ebiten.Image) {
	g.player.DrawParallaxBackground(screen, g.background, screenWidth, screenHeight, g.world.Def.Width)
	g.player.DrawLevel(screen, g.world.Quadtree, screenWidth, screenHeight, g.tileset)
	g.portals.Draw(screen, g.world.Level.Portals, g.player)
	for _, e := range g.world.Enemies {
		e.DrawEnemyAnimation(screen, g.player.Camera)
		e.DrawHealthBar(screen, g.player.Camera)
	}
	g.player.DrawPlayerAnimation(screen)
	g.player.DrawGroundWaves(screen)
	// Water is drawn after the player so the player appears submerged in it.
	g.player.DrawWater(screen, g.world.Quadtree, screenWidth, screenHeight, g.tileset)
	g.portals.DrawPrompt(screen, g.world.Level.Portals, g.player, g.touch.Enabled)
	if g.showHitboxes {
		g.drawHitboxes(screen)
	}
	if g.touch.Enabled {
		g.player.DrawHUDAt(screen, core.TouchHUDLeft)
	} else {
		g.player.DrawHUD(screen)
	}
	g.travel.Draw(screen)
	g.touch.Draw(screen)
}

// drawHitboxes outlines the body and attack boxes of the enemies and the player.
func (g *Game) drawHitboxes(screen *ebiten.Image) {
	providers := make([]core.HitboxProvider, 0, len(g.world.Enemies)+len(g.world.Level.Portals)+1)
	for _, p := range g.world.Level.Portals {
		providers = append(providers, p)
	}
	for _, e := range g.world.Enemies {
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
