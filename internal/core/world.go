package core

import (
	"encoding/json"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// eg:= grass ,ice ,sand etc.
type TileType string
type TileLvl string
type TileLength int

const (
	// tileset constants
	Tileset                     = "../assets/LevelTile.png"
	Grass            TileType   = "grass"
	Rock             TileType   = "stone"
	Ice              TileType   = "ice"
	Sand             TileType   = "sand"
	Larva            TileType   = "larva"
	Metal            TileType   = "metal"
	Water            TileType   = "water"
	Wood             TileType   = "wood"
	Empty            TileType   = "empty"
	TopTileLen       TileLength = 11
	BottomTileLen    TileLength = 4
	TopTile          TileLvl    = "top"
	BottomTile       TileLvl    = "bottom"
	PixelTileWidth              = 192 // before it was 4
	PixelTileHeight             = 192 // before it was 4
	TotTopTilesWidth            = 2112
	TotBotTilesWidth            = 768
	TopTileXStart               = 0
	BotTileXStart               = TotTopTilesWidth

	// level constants
	LevelFile        = "../assets/level_1.json" // embedded JSON level bundle
	DefaultLevelName = "level0"                 // key to load from the bundle
	Background_1     = "../assets/LBackground_3.png"

	// ---------------- level constants ----------------
	LevelTileWidth      = 60      // before it was 40
	LevelTileHeight     = 60      // before it was 40
	TopTileVisualOffset = 26.5625 // for water ,grass ,sand
	Level_1_Width       = 6000
	Level_1_Height      = 2040
)

type Tile struct {
	X, Y     int      // these are the texture coordinates
	TileLvl  TileLvl  // top or bottom
	TileType TileType // eg:= grass ,ice ,sand etc.
}

// CheckpointCode is the grid value that marks a (non-solid) respawn point.
const CheckpointCode = 9

// ---------------- platform ----------------
type Platform struct {
	X, Y, Width, Height float64 // these are the world coordinates
	TileInfo            Tile    // this is the tile information
	DrawOffsetY         float64 // this is the draw offset for the tile
}

// Checkpoint is a respawn location the player activates by touching it.
type Checkpoint struct {
	Pos Position
}

// Bounds returns the checkpoint's one-tile activation zone.
func (c Checkpoint) Bounds() AABB {
	return AABB{X: c.Pos.X, Y: c.Pos.Y, Width: LevelTileWidth, Height: LevelTileHeight}
}

// Level is the built, ready-to-use result of loading a level: its solid
// platforms and its checkpoints.
type Level struct {
	Platforms   []Platform
	Checkpoints []Checkpoint
}

var Tiles map[TileType][2]Tile

// tileRowOrder lists tile types in the order their rows appear in the tileset,
// top to bottom. The slice index is the tileset row for that type.
var tileRowOrder = []TileType{Water, Grass, Sand, Rock, Metal, Ice, Wood, Larva}

// WorldInit builds the tileset lookup table. Each type maps to its top and
// bottom variants; both share the type's row and differ only by X column.
// It must be called once before LoadLevel.
func WorldInit() {
	Tiles = make(map[TileType][2]Tile, len(tileRowOrder))
	for row, tileType := range tileRowOrder {
		y := PixelTileHeight * row
		Tiles[tileType] = [2]Tile{
			{X: TopTileXStart, Y: y, TileType: tileType, TileLvl: TopTile},
			{X: BotTileXStart, Y: y, TileType: tileType, TileLvl: BottomTile},
		}
	}
}

// GetBounds returns the bounding box of the platform
func (p *Platform) GetBounds() AABB {
	return AABB{
		X:      p.X,
		Y:      p.Y,
		Width:  p.Width,
		Height: p.Height,
	}
}

// tileTypeFromCode maps a JSON grid value to a TileType. 0 (and any value
// outside the valid range) is Empty; codes 1..N select the Nth entry of
// tileRowOrder, so 1=water, 2=grass, 3=sand, 4=rock, 5=metal, 6=ice, 7=wood,
// 8=larva.
func tileTypeFromCode(code int) TileType {
	if code >= 1 && code <= len(tileRowOrder) {
		return tileRowOrder[code-1]
	}
	return Empty
}

// cellAt reads a grid cell, returning 0 (empty) for any out-of-range coordinate
// so callers can probe neighbours without bounds-checking.
func cellAt(grid [][]int, x, y int) int {
	if y < 0 || y >= len(grid) || x < 0 || x >= len(grid[y]) {
		return 0
	}
	return grid[y][x]
}

// tileVariantIndex maps a tile level to its index within a Tiles entry
// ([0] = top, [1] = bottom).
func tileVariantIndex(lvl TileLvl) int {
	if lvl == BottomTile {
		return 1
	}
	return 0
}

// hasVisualOffset reports whether a top tile of this type sits lower in its cell
// (water, grass and sand have a transparent margin above the solid surface).
func hasVisualOffset(t TileType) bool {
	return t == Water || t == Grass || t == Sand
}

func getTileInfo(x, y int, tileType TileType, grid [][]int, prevPlat Platform) Platform {
	// Create the platform with basic world coordinates.
	plat := Platform{
		X:      float64(x * LevelTileWidth),
		Y:      float64(y * LevelTileHeight),
		Width:  LevelTileWidth,
		Height: LevelTileHeight,
		TileInfo: Tile{
			TileType: tileType,
			TileLvl:  TopTile, // default to top; corrected below
		},
	}

	// A tile with a solid neighbour directly above is a bottom (interior) tile.
	if y > 0 && tileTypeFromCode(cellAt(grid, x, y-1)) != Empty {
		plat.TileInfo.TileLvl = BottomTile
	}

	// Look up the sprite-sheet source coordinates for this type/level. Unknown
	// types are absent from Tiles and keep the zero origin.
	if variants, ok := Tiles[tileType]; ok {
		src := variants[tileVariantIndex(plat.TileInfo.TileLvl)]
		plat.TileInfo.X = src.X
		plat.TileInfo.Y = src.Y
	}

	// Top tiles with a transparent margin are nudged down so the solid part lines up.
	if plat.TileInfo.TileLvl == TopTile && hasVisualOffset(tileType) {
		plat.Y += TopTileVisualOffset
		plat.DrawOffsetY = -TopTileVisualOffset
		plat.Height -= TopTileVisualOffset
	}

	if prevPlat.TileInfo.TileType == plat.TileInfo.TileType && prevPlat.TileInfo.TileLvl == plat.TileInfo.TileLvl {
		plat.TileInfo.X = prevPlat.TileInfo.X + PixelTileWidth

		if plat.TileInfo.TileLvl == TopTile {
			if plat.TileInfo.X >= TotTopTilesWidth {
				plat.TileInfo.X = TopTileXStart
			}
		} else {
			if plat.TileInfo.X >= BotTileXStart+4*PixelTileWidth {
				plat.TileInfo.X = BotTileXStart
			}
		}
	}
	return plat
}

// levelBundle is the on-disk JSON schema: a map of level name to its grid.
// Example: {"level0": {"grid": [[0,1,...], ...]}}.
type levelBundle map[string]struct {
	Grid [][]int `json:"grid"`
}

// LoadLevel reads the embedded JSON level bundle at path and builds the named
// level. Only the base file name of path is used, so it works on native and
// WebAssembly builds alike. WorldInit must be called before this.
func LoadLevel(path, name string) Level {
	data := readAsset(path)

	var bundle levelBundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		log.Fatalf("parse level %q: %v", path, err)
	}

	level, ok := bundle[name]
	if !ok {
		log.Fatalf("level %q not found in %q", name, path)
	}
	return BuildLevel(level.Grid)
}

// BuildLevel converts a grid of tile codes into a Level: one solid Platform per
// non-empty tile cell, plus a Checkpoint per CheckpointCode cell. It is a pure
// function of its input, so it is easy to test and reuse independently of asset
// loading.
func BuildLevel(grid [][]int) Level {
	level := Level{Platforms: []Platform{}}
	prevPlat := Platform{TileInfo: Tile{TileType: Empty}}

	for y := range grid {
		prevPlat.TileInfo.TileType = Empty

		for x := range grid[y] {
			code := grid[y][x]

			if code == CheckpointCode {
				level.Checkpoints = append(level.Checkpoints, Checkpoint{
					Pos: Position{X: float64(x * LevelTileWidth), Y: float64(y * LevelTileHeight)},
				})
			}

			// Checkpoint cells are non-solid (code 9 maps to Empty), so no platform.
			tileType := tileTypeFromCode(code)
			if tileType != Empty {
				level.Platforms = append(level.Platforms, getTileInfo(x, y, tileType, grid, prevPlat))
			}

			if len(level.Platforms) > 0 {
				prevPlat = level.Platforms[len(level.Platforms)-1]
			}
			prevPlat.TileInfo.TileType = tileType
		}
	}
	return level
}

func (player *PlayerRuntime) DrawParallaxBackground(screen *ebiten.Image, background *ebiten.Image, screenWidth, screenHeight float64) {
	bgW := background.Bounds().Dx()
	bgH := background.Bounds().Dy()

	parallaxFactor := 0.3

	baseScale := screenHeight / float64(bgH)
	scale := baseScale * (1 + parallaxFactor)

	scaledW := float64(bgW) * scale

	// Calculate offset based on camera position relative to total level width
	maxLevelWidth := float64(Level_1_Width)
	maxCamX := maxLevelWidth - screenWidth
	if maxCamX < 1 {
		maxCamX = 1
	}

	currentCamX := player.Camera.Pos.X
	if currentCamX < 0 {
		currentCamX = 0
	} else if currentCamX > maxCamX {
		currentCamX = maxCamX
	}

	// Map camera progress (0-1) to background scroll
	// We shift the background from 0 to -(scaledW - screenWidth)
	progression := currentCamX / maxCamX
	maxBgOffset := scaledW - screenWidth
	bgOffset := progression * maxBgOffset

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(-bgOffset, 0)
	screen.DrawImage(background, op)
}

// DrawLevel draws the solid (non-water) tiles visible in the viewport. It is
// drawn behind the player; water is drawn afterwards by DrawWater so the player
// appears submerged.
func (player *PlayerRuntime) DrawLevel(screen *ebiten.Image, quadtree *DynamicQuadtree, screenWidth, screenHeight float64, tileset *ebiten.Image) {
	player.drawTiles(screen, quadtree, screenWidth, screenHeight, tileset, false)
}

// DrawWater draws only the water tiles visible in the viewport. Call it after
// drawing the player so water renders in front, showing submersion.
func (player *PlayerRuntime) DrawWater(screen *ebiten.Image, quadtree *DynamicQuadtree, screenWidth, screenHeight float64, tileset *ebiten.Image) {
	player.drawTiles(screen, quadtree, screenWidth, screenHeight, tileset, true)
}

// drawTiles renders the visible tiles of one layer: water only when water is
// true, everything else when it is false.
func (player *PlayerRuntime) drawTiles(screen *ebiten.Image, quadtree *DynamicQuadtree, screenWidth, screenHeight float64, tileset *ebiten.Image, water bool) {
	viewport := AABB{
		X:      player.Camera.Pos.X,
		Y:      player.Camera.Pos.Y,
		Width:  screenWidth,
		Height: screenHeight,
	}

	scaleX := float64(LevelTileWidth) / float64(PixelTileWidth)
	scaleY := float64(LevelTileHeight) / float64(PixelTileHeight)

	for _, obj := range quadtree.Retrieve(viewport) {
		p, ok := obj.(*Platform)
		if !ok || (p.TileInfo.TileType == Water) != water {
			continue
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(p.X-player.Camera.Pos.X, p.Y+p.DrawOffsetY-player.Camera.Pos.Y)

		src := image.Rect(int(p.TileInfo.X), int(p.TileInfo.Y), int(p.TileInfo.X)+PixelTileWidth, int(p.TileInfo.Y)+PixelTileHeight)
		screen.DrawImage(tileset.SubImage(src).(*ebiten.Image), op)
	}
}
