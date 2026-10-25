package main

import (
	"errors"
	"log"
	"runtime"

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
	screens *core.Screens // the menus over the game; ScreenNone while playing
	input   core.InputState
	touch   *core.TouchControls // on-screen buttons for touch screens

	player      *core.PlayerRuntime
	playerSheet *ebiten.Image
	world       *core.World // the current level: platforms, portals, enemies
	respawns    int         // player respawns already handled (enemies revived)

	background     *ebiten.Image
	backgroundPath string // the asset background was loaded from
	tileset        *ebiten.Image
	portals        *core.PortalSprite
	travel         *core.TravelMenu // opened by holding E (or a touch) at a portal
	view           *ebiten.Image    // scratch screen for drawing a portal's thumbnail

	autosave core.Autosave // writes the progress when it changes

	showHitboxes bool // debug overlay, toggled with F3
}

// NewGame loads assets and returns a Game on its title screen.
func NewGame() *Game {
	core.WorldInit()

	touch := core.NewTouchControls(screenWidth, screenHeight)
	touch.Enabled = system.IsTouchDevice()

	g := &Game{
		screens:     core.NewScreens(screenWidth, screenHeight, runtime.GOOS != "js"), // a page can't close its tab (D6)
		playerSheet: core.LoadImage(playerSpriteSheetPath),
		tileset:     core.LoadImage(core.Tileset),
		touch:       touch,
		portals:     core.NewPortalSprite(),
		view:        ebiten.NewImage(screenWidth, screenHeight),
	}
	g.screens.SetSettings(loadSettings())
	g.openTitle()
	return g
}

// loadSettings reads the stored settings; none or unreadable ones give the
// defaults.
func loadSettings() core.Settings {
	data, ok := system.LoadRecord(core.SettingsKey)
	if !ok {
		return core.DefaultSettings()
	}
	s, err := core.DecodeSettings(data)
	if err != nil {
		log.Printf("settings ignored: %v", err)
	}
	return s
}

// saveSettings stores the settings shown on the settings screen.
func (g *Game) saveSettings() {
	data, err := core.EncodeSettings(g.screens.Settings)
	if err == nil {
		err = system.StoreRecord(core.SettingsKey, data)
	}
	if err != nil {
		log.Printf("save settings: %v", err)
	}
}

// openTitle shows the title screen over the stored save's game, ready for
// Continue, or over a new game when there is no save it can use.
func (g *Game) openTitle() {
	s, err := loadSave()
	if err == nil {
		err = g.startRun(&s)
	}
	if err != nil {
		if !errors.Is(err, errNoSave) {
			log.Printf("save ignored: %v", err)
		}
		_ = g.startRun(nil)
	}
	g.screens.OpenTitle(err == nil)
}

var errNoSave = errors.New("no save")

// loadSave reads the stored save.
func loadSave() (core.SaveData, error) {
	data, ok := system.LoadRecord(core.SaveKey)
	if !ok {
		return core.SaveData{}, errNoSave
	}
	return core.DecodeSave(data)
}

// startRun sets up a fresh run: a new game standing in the start portal, or,
// given a save, standing in its checkpoint with its portals active.
func (g *Game) startRun(save *core.SaveData) error {
	world, err := core.NewWorld(core.StartLevel)
	if err != nil {
		log.Fatal(err)
	}
	player := core.InitPlayer(g.playerSheet)
	player.StartAt(world.Level)
	if save != nil {
		if err := world.Continue(player, *save); err != nil {
			return err
		}
	}
	player.UpdateCamera(screenWidth, screenHeight, world.Def.Width, world.Def.Height) // the menus show the player behind them
	g.world, g.player, g.respawns = world, player, 0
	g.travel = core.NewTravelMenu(screenWidth, screenHeight)
	g.travel.WaitRelease() // the E or tap that picked the menu item can't open it
	g.autosave = core.Autosave{}
	g.syncBackground()
	return nil
}

// newGame wipes the save and starts in the start portal.
func (g *Game) newGame() {
	if err := system.DeleteRecord(core.SaveKey); err != nil {
		log.Printf("wipe save: %v", err)
	}
	_ = g.startRun(nil)
}

// saveProgress stores the player's progress when it changed since the last write.
func (g *Game) saveProgress() {
	s, due := g.autosave.Due(g.player)
	if !due {
		return
	}
	data, err := core.EncodeSave(s)
	if err == nil {
		err = system.StoreRecord(core.SaveKey, data)
	}
	if err != nil {
		log.Printf("save: %v", err)
	}
	g.autosave.Written(s) // also on failure: retrying every tick would only spam the log
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
	if g.screens.Current != core.ScreenNone { // the game waits behind the menus
		switch g.screens.Update(&g.input) {
		case core.ActionNewGame:
			g.newGame()
		case core.ActionQuit:
			return ebiten.Termination
		case core.ActionSettingsChanged:
			g.saveSettings()
		case core.ActionToTitle:
			g.openTitle() // the progress is already saved (autosave)
		}
		// ActionContinue: the save's game is already set up behind the title.
		g.portals.Update()
		return nil
	}
	if g.input.Menu && !g.travel.Open { // Esc while the travel menu is open only closes it
		g.screens.OpenPause()
		return nil
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
	g.travel.MaybeOpen(g.player, &g.input)
	g.saveProgress()

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
	if g.screens.Current != core.ScreenNone {
		g.screens.Draw(screen)
		return
	}
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
