package core

import (
	"image"

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

	// level 1 constants
	Level_1      = "../assets/Level_1.png"
	Background_1 = "../assets/LBackground_3.png"

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

// ---------------- platform ----------------
type Platform struct {
	X, Y, Width, Height float64 // these are the world coordinates
	TileInfo            Tile    // this is the tile information
	DrawOffsetY         float64 // this is the draw offset for the tile
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

func getTileType(r, g, b uint32) TileType {
	// larva color code -> 255 ,0 ,0
	// grass color code -> 0 ,255 ,0
	// water color code -> 0 ,0 ,255
	// stone color code -> 255 ,255 ,0
	switch {
	case r == 255 && g == 0 && b == 0:
		return Larva
	case r == 0 && g == 255 && b == 0:
		return Grass
	case r == 0 && g == 0 && b == 255:
		return Water
	case r == 255 && g == 255 && b == 0:
		return Rock
	default:
		return Empty
	}
}

func getColor(x, y int, levelData *ebiten.Image) (uint32, uint32, uint32) {
	color := levelData.At(x, y)
	r, g, b, _ := color.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	return r, g, b
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

func getTileInfo(x, y int, tileType TileType, levelData *ebiten.Image, prevPlat Platform) Platform {
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
	if y > 0 && getTileType(getColor(x, y-1, levelData)) != Empty {
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

// LoadLevel decodes the level layout from the color-coded level image, returning
// one Platform per solid pixel. WorldInit must be called before this.
func LoadLevel(levelData *ebiten.Image) []Platform {
	level := []Platform{}

	prevPlat := Platform{TileInfo: Tile{TileType: Empty}}

	for y := 0; y < levelData.Bounds().Dy(); y++ {
		prevPlat.TileInfo.TileType = Empty

		for x := 0; x < levelData.Bounds().Dx(); x++ {
			tileType := getTileType(getColor(x, y, levelData))

			if tileType != Empty {
				level = append(level, getTileInfo(x, y, tileType, levelData, prevPlat))
			}

			if len(level) > 0 {
				prevPlat = level[len(level)-1]
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

func (player *PlayerRuntime) DrawLevel(screen *ebiten.Image, quadtree *DynamicQuadtree, screenWidth, screenHeight float64, tileset *ebiten.Image) {

	// Define the camera viewport
	viewport := AABB{
		X:      player.Camera.Pos.X,
		Y:      player.Camera.Pos.Y,
		Width:  screenWidth,
		Height: screenHeight,
	}

	// Retrieve visible platforms from the Quadtree
	visibleObjects := quadtree.Retrieve(viewport)

	// Draw visible platforms
	for _, obj := range visibleObjects {
		if p, ok := obj.(*Platform); ok {
			// draw the tile image
			op := &ebiten.DrawImageOptions{}
			// Scale the tile image (PixelTileWidth/Height) to fit the platform size (LevelTileWidth/Height)
			scaleX := float64(LevelTileWidth) / float64(PixelTileWidth)
			scaleY := float64(LevelTileHeight) / float64(PixelTileHeight)
			op.GeoM.Scale(scaleX, scaleY)

			// translate is used to position the tile image on the screen
			op.GeoM.Translate(float64(p.X-player.Camera.Pos.X), float64(p.Y+p.DrawOffsetY-player.Camera.Pos.Y))

			// Draw the sub-image from the tileset using coordinates from TileInfo
			screen.DrawImage(tileset.SubImage(image.Rect(int(p.TileInfo.X), int(p.TileInfo.Y), int(p.TileInfo.X)+PixelTileWidth, int(p.TileInfo.Y)+PixelTileHeight)).(*ebiten.Image), op)
		}
	}
}
