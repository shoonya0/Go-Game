package main

import (
	"errors"
	"image"
	"image/color"
	"log"
	"runtime"
	"time"

	"player/internal/core"
	"player/internal/system"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	playerSpriteSheetPath = "../assets/player/GideonGraves.png"
	loadingTextSize       = 16 // "Loading..." while the web build downloads a background
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

	background     *ebiten.Image // nil while the web build downloads it
	backgroundPath string        // the asset background was loaded from
	backgrounds    chan downloadedBackground
	retryFailed    string    // background whose download failed, tried again at retryAt
	retryAt        time.Time // so a server that is down isn't asked 60 times a second
	tileset        *ebiten.Image
	portals        *core.PortalSprite
	travel         *core.TravelMenu // opened by holding E (or a touch) at a portal

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
		backgrounds: make(chan downloadedBackground, 1),
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
	player.UpdateCamera(screenWidth, screenHeight, world.Level.Width, world.Level.Height) // the menus show the player behind them
	g.world, g.player, g.respawns = world, player, 0
	g.travel = core.NewTravelMenu(screenWidth, screenHeight)
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

// toTitle leaves the game for the title screen, Continue ready at the run's
// progress. That progress is stored now but not read back, so a write that
// fails (browser storage blocked) doesn't lose the run.
func (g *Game) toTitle() {
	s := g.player.Snapshot()
	g.storeProgress(s)
	err := g.startRun(&s)
	if err != nil {
		log.Printf("back to title: %v", err)
		_ = g.startRun(nil)
	}
	g.screens.OpenTitle(err == nil)
}

// saveProgress stores the player's progress when it changed since the last write.
func (g *Game) saveProgress() {
	if s, due := g.autosave.Due(g.player); due {
		g.storeProgress(s)
	}
}

// storeProgress writes save s.
func (g *Game) storeProgress(s core.SaveData) {
	data, err := core.EncodeSave(s)
	if err == nil {
		err = system.StoreRecord(core.SaveKey, data)
	}
	if err != nil {
		log.Printf("save: %v", err)
	}
	g.autosave.Written(s) // also on failure: retrying every tick would only spam the log
}

// backgroundRetry is how long after a failed download the web build tries
// again.
const backgroundRetry = 2 * time.Second

// downloadedBackground is a level background the web build downloaded.
type downloadedBackground struct {
	path string
	img  image.Image
	err  error
}

// syncBackground loads the current level's background when it changed. An
// embedded one (desktop) loads at once; the web build downloads it off the
// game loop and draws no background until it arrives (levels decision D4).
func (g *Game) syncBackground() {
	for done := false; !done; {
		select {
		case d := <-g.backgrounds:
			switch {
			case d.err != nil:
				log.Printf("background: %v", d.err)
				if d.path == g.backgroundPath { // still needed: try again in a moment
					g.backgroundPath, g.retryFailed, g.retryAt = "", d.path, time.Now().Add(backgroundRetry)
				}
			case d.path == g.backgroundPath: // not one the player already left
				g.setBackground(ebiten.NewImageFromImage(d.img))
			}
		default:
			done = true
		}
	}

	path := g.world.Def.Background
	if path == g.backgroundPath || (path == g.retryFailed && time.Now().Before(g.retryAt)) {
		return
	}
	g.backgroundPath = path
	if core.HasAsset(path) {
		g.setBackground(core.LoadImage(path))
		return
	}
	g.setBackground(nil)
	go func() {
		d := downloadedBackground{path: path}
		var data []byte
		if data, d.err = system.FetchFile(core.AssetURL(path)); d.err == nil {
			d.img, d.err = core.DecodeImage(data)
		}
		g.backgrounds <- d
	}()
}

// followPlayer brings the camera to the player at once after a teleport, so
// the next frame isn't drawn from where the player was.
func (g *Game) followPlayer() {
	g.player.UpdateCamera(screenWidth, screenHeight, g.world.Level.Width, g.world.Level.Height)
}

// setBackground replaces the background, freeing the old one's GPU memory at
// once (a background is tens of MB decoded).
func (g *Game) setBackground(img *ebiten.Image) {
	if g.background != nil {
		g.background.Deallocate()
	}
	g.background = img
}

// Update advances the simulation one tick (runs at ~60 TPS).
func (g *Game) Update() error {
	g.syncBackground() // picks up a background the web build finished downloading
	system.HandleInput(&g.input)
	g.touch.Hidden = g.screens.Current != core.ScreenNone
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
			g.toTitle()
		}
		// ActionContinue: the save's game is already set up behind the title.
		if g.screens.Current == core.ScreenNone {
			g.travel.WaitRelease() // the E or touch that left the menus can't open the travel menu
		}
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
			g.followPlayer()
			g.syncBackground()
			g.saveProgress() // the checkpoint moved: store it now, not on the next played tick
		}
		g.portals.Update()
		return nil
	}

	core.UpdatePlayer(g.player, &g.input, g.world.Quadtree, g.world.Level.Portals)
	g.player.UpdateAnimation()
	g.player.HitEnemies(g.world.Enemies)
	g.player.UpdateGroundWaves(g.world.Quadtree, g.world.Enemies)
	g.player.UpdateCamera(screenWidth, screenHeight, g.world.Level.Width, g.world.Level.Height)

	if g.player.Respawns != g.respawns {
		g.respawns = g.player.Respawns
		g.world.Respawned(g.player)
		g.followPlayer()
		g.syncBackground()
	}
	if g.player.LevelDone() { // the last unused portal: on to the next level (levels D5)
		if opened, err := g.world.OpenNextLevel(g.player); err != nil {
			log.Printf("next level: %v", err)
		} else if opened {
			g.followPlayer()
			g.syncBackground()
			g.travel.WaitRelease() // an E still held mustn't open the menu at the arrival portal
		}
	}
	g.updateEnemies()
	g.portals.Update()
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

// Draw renders one frame: background, level, then player.
func (g *Game) Draw(screen *ebiten.Image) {
	g.player.DrawParallaxBackground(screen, g.background, screenWidth, screenHeight, g.world.Level.Width)
	if g.background == nil { // the web build is still downloading it
		core.DrawTextCentered(screen, "Loading...", loadingTextSize, screenWidth/2, screenHeight/2, color.White)
	}
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
