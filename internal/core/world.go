package core

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// eg:= grass ,ice ,sand etc.
type TileType string
type TileLvl string
type TileLength int

// L -> 703
// B -> 193

// height -> 192

const (
	Tileset          = "../assets/LevelTile.png"
	Grass   TileType = "grass"
	Rock    TileType = "stone"
	Ice     TileType = "ice"
	Sand    TileType = "sand"
	Larva   TileType = "larva"
	Metal   TileType = "metal"
	Water   TileType = "water"
	Wood    TileType = "wood"
	Empty   TileType = "empty"

	TopTileLen    TileLength = 11
	BottomTileLen TileLength = 4

	TopTile    TileLvl = "top"
	BottomTile TileLvl = "bottom"

	PixelTileWidth  = 192 // before it was 4
	PixelTileHeight = 192 // before it was 4

	TotTopTilesWidth = 2112
	TotBotTilesWidth = 768

	TopTileXStart = 0
	BotTileXStart = TotTopTilesWidth
)

const (
	Level_1             = "../assets/Level_1.png"
	LevelTileWidth      = 60      // before it was 40
	LevelTileHeight     = 60      // before it was 40
	TopTileVisualOffset = 26.5625 // for water ,grass ,sand
)

type Tile struct {
	X, Y     int      // these are the texture coordinates
	TileLvl  TileLvl  // top or bottom
	TileType TileType // eg:= grass ,ice ,sand etc.
}

// ---------------- platform ----------------
type Platform struct {
	X, Y, Width, Height float64 // these are the world coordinates
	TileInfo            Tile
	DrawOffsetY         float64
}

var Tiles map[TileType][2]Tile

// initialize the world at the start of the program
func WorldInit() {
	// Initialize the map
	Tiles = make(map[TileType][2]Tile)
	Tiles[Water] = [2]Tile{
		{X: TopTileXStart, Y: PixelTileHeight * 0, TileType: Water, TileLvl: TopTile},
		{X: BotTileXStart, Y: PixelTileHeight * 0, TileType: Water, TileLvl: BottomTile},
	}
	Tiles[Grass] = [2]Tile{
		{X: TopTileXStart, Y: PixelTileHeight * 1, TileType: Grass, TileLvl: TopTile},
		{X: BotTileXStart, Y: PixelTileHeight * 1, TileType: Grass, TileLvl: BottomTile},
	}
	Tiles[Sand] = [2]Tile{
		{X: TopTileXStart, Y: PixelTileHeight * 2, TileType: Sand, TileLvl: TopTile},
		{X: BotTileXStart, Y: PixelTileHeight * 2, TileType: Sand, TileLvl: BottomTile},
	}
	Tiles[Rock] = [2]Tile{
		{X: TopTileXStart, Y: PixelTileHeight * 3, TileType: Rock, TileLvl: TopTile},
		{X: BotTileXStart, Y: PixelTileHeight * 3, TileType: Rock, TileLvl: BottomTile},
	}
	Tiles[Metal] = [2]Tile{
		{X: TopTileXStart, Y: PixelTileHeight * 4, TileType: Metal, TileLvl: TopTile},
		{X: BotTileXStart, Y: PixelTileHeight * 4, TileType: Metal, TileLvl: BottomTile},
	}
	Tiles[Ice] = [2]Tile{
		{X: TopTileXStart, Y: PixelTileHeight * 5, TileType: Ice, TileLvl: TopTile},
		{X: BotTileXStart, Y: PixelTileHeight * 5, TileType: Ice, TileLvl: BottomTile},
	}
	Tiles[Wood] = [2]Tile{
		{X: TopTileXStart, Y: PixelTileHeight * 6, TileType: Wood, TileLvl: TopTile},
		{X: BotTileXStart, Y: PixelTileHeight * 6, TileType: Wood, TileLvl: BottomTile},
	}
	Tiles[Larva] = [2]Tile{
		{X: TopTileXStart, Y: PixelTileHeight * 7, TileType: Larva, TileLvl: TopTile},
		{X: BotTileXStart, Y: PixelTileHeight * 7, TileType: Larva, TileLvl: BottomTile},
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

// ---------------- level constants ----------------
const (
	Level_1_Width  = 6000
	Level_1_Height = 2040
)

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

func getTileInfo(x, y int, tileType TileType, levelData *ebiten.Image, prevPlat Platform) Platform {
	// Create the platform with basic world coordinates
	plat := Platform{
		X:      float64(x * LevelTileWidth),
		Y:      float64(y * LevelTileHeight),
		Width:  LevelTileWidth,
		Height: LevelTileHeight,
		TileInfo: Tile{
			TileType: tileType,
			TileLvl:  TopTile, // Default to Top
		},
		DrawOffsetY: 0,
	}

	// check for bottom tile
	if y > 0 {
		if getTileType(getColor(x, y-1, levelData)) != Empty {
			plat.TileInfo.TileLvl = BottomTile
		}
	}

	if plat.TileInfo.TileLvl == TopTile {
		switch plat.TileInfo.TileType {
		case Water:
			plat.TileInfo.X = Tiles[Water][0].X
			plat.TileInfo.Y = Tiles[Water][0].Y
		case Grass:
			plat.TileInfo.X = Tiles[Grass][0].X
			plat.TileInfo.Y = Tiles[Grass][0].Y
		case Sand:
			plat.TileInfo.X = Tiles[Sand][0].X
			plat.TileInfo.Y = Tiles[Sand][0].Y
		case Rock:
			plat.TileInfo.X = Tiles[Rock][0].X
			plat.TileInfo.Y = Tiles[Rock][0].Y
		case Metal:
			plat.TileInfo.X = Tiles[Metal][0].X
			plat.TileInfo.Y = Tiles[Metal][0].Y
		case Ice:
			plat.TileInfo.X = Tiles[Ice][0].X
			plat.TileInfo.Y = Tiles[Ice][0].Y
		case Wood:
			plat.TileInfo.X = Tiles[Wood][0].X
			plat.TileInfo.Y = Tiles[Wood][0].Y
		case Larva:
			plat.TileInfo.X = Tiles[Larva][0].X
			plat.TileInfo.Y = Tiles[Larva][0].Y
		default:
			plat.TileInfo.X = 0
			plat.TileInfo.Y = 0
		}

		if plat.TileInfo.TileType == Water || plat.TileInfo.TileType == Grass || plat.TileInfo.TileType == Sand {
			plat.Y += TopTileVisualOffset
			plat.DrawOffsetY = -TopTileVisualOffset
			plat.Height -= TopTileVisualOffset
		}
	} else {
		switch plat.TileInfo.TileType {
		case Water:
			plat.TileInfo.X = Tiles[Water][1].X
			plat.TileInfo.Y = Tiles[Water][1].Y
		case Grass:
			plat.TileInfo.X = Tiles[Grass][1].X
			plat.TileInfo.Y = Tiles[Grass][1].Y
		case Sand:
			plat.TileInfo.X = Tiles[Sand][1].X
			plat.TileInfo.Y = Tiles[Sand][1].Y
		case Rock:
			plat.TileInfo.X = Tiles[Rock][1].X
			plat.TileInfo.Y = Tiles[Rock][1].Y
		case Metal:
			plat.TileInfo.X = Tiles[Metal][1].X
			plat.TileInfo.Y = Tiles[Metal][1].Y
		case Ice:
			plat.TileInfo.X = Tiles[Ice][1].X
			plat.TileInfo.Y = Tiles[Ice][1].Y
		case Wood:
			plat.TileInfo.X = Tiles[Wood][1].X
			plat.TileInfo.Y = Tiles[Wood][1].Y
		case Larva:
			plat.TileInfo.X = Tiles[Larva][1].X
			plat.TileInfo.Y = Tiles[Larva][1].Y
		default:
			plat.TileInfo.X = 0
			plat.TileInfo.Y = 0
		}
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

// LoadLevel loads the level from the image
func (p *PlayerRuntime) LoadLevel(levelData *ebiten.Image) []Platform {
	level := []Platform{}

	var prevPlat Platform = Platform{TileInfo: Tile{TileType: Empty}}

	fmt.Println(levelData.Bounds())

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
			// this is because the y axis is -ive upside and -ive x axis is to the left
			// vector.DrawFilledRect(screen, float32(p.X-player.Camera.Pos.X), float32(p.Y-player.Camera.Pos.Y), float32(p.Width), float32(p.Height), color.RGBA{100, 100, 100, 255}, false)

			op := &ebiten.DrawImageOptions{}
			// Scale the tile image (PixelTileWidth/Height) to fit the platform size (LevelTileWidth/Height)
			scaleX := float64(LevelTileWidth) / float64(PixelTileWidth)
			scaleY := float64(LevelTileHeight) / float64(PixelTileHeight)
			op.GeoM.Scale(scaleX, scaleY)

			// translate is used to position the tile image on the screen
			// op.GeoM.Translate(float64(p.X-player.Camera.Pos.X), float64(p.Y-player.Camera.Pos.Y))
			op.GeoM.Translate(float64(p.X-player.Camera.Pos.X), float64(p.Y+p.DrawOffsetY-player.Camera.Pos.Y))

			// Draw the sub-image from the tileset using coordinates from TileInfo
			screen.DrawImage(tileset.SubImage(image.Rect(int(p.TileInfo.X), int(p.TileInfo.Y), int(p.TileInfo.X)+PixelTileWidth, int(p.TileInfo.Y)+PixelTileHeight)).(*ebiten.Image), op)
		}
	}
}
