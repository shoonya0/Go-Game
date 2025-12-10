package core

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type TileType string
type CornerType string

const (
	Tileset          = "../assets/Tileset.png"
	Grass   TileType = "grass"
	Stone   TileType = "stone"
	Ice     TileType = "ice"
	Sand    TileType = "sand"
	Larva   TileType = "larva"
	Metal   TileType = "metal"
	Water   TileType = "water"
	Wood    TileType = "wood"
	Empty   TileType = "empty"

	StartTile   CornerType = "start"
	MiddleTile  CornerType = "middle"
	EndTile     CornerType = "end"
	CoveredTile CornerType = "covered"

	PixelTileWidth  = 200 // before it was 4
	PixelTileHeight = 175 // before it was 4
)

const (
	Level_1         = "../assets/Level_1.png"
	LevelTileWidth  = 40 // before it was 4
	LevelTileHeight = 40 // before it was 4
)

type Tile struct {
	X, Y       int
	CornerType CornerType
	TileType   TileType
}

// ---------------- platform ----------------
type Platform struct {
	X, Y, Width, Height float64
	TileInfo            Tile
}

var Tiles *map[TileType][4]*Tile

// initialize the world at the start of the program
func WorldInit() {
	Tiles = new(map[TileType][4]*Tile)   // nil map
	*Tiles = make(map[TileType][4]*Tile) // empty map
	(*Tiles)[Grass] = [4]*Tile{
		{X: 35, Y: 55, TileType: Grass, CornerType: StartTile},
		{X: 262, Y: 55, TileType: Grass, CornerType: MiddleTile},
		{X: 486, Y: 55, TileType: Grass, CornerType: EndTile},
		{X: 1173, Y: 61, TileType: Grass, CornerType: CoveredTile},
	}
	(*Tiles)[Stone] = [4]*Tile{
		{X: 34, Y: 296, TileType: Stone, CornerType: StartTile},
		{X: 262, Y: 288, TileType: Stone, CornerType: MiddleTile},
		{X: 487, Y: 288, TileType: Stone, CornerType: EndTile},
		{X: 937, Y: 285, TileType: Stone, CornerType: CoveredTile},
	}
	(*Tiles)[Ice] = [4]*Tile{
		{X: 0, Y: 0, TileType: Ice, CornerType: StartTile},
		{X: 0, Y: 0, TileType: Ice, CornerType: MiddleTile},
		{X: 0, Y: 0, TileType: Ice, CornerType: EndTile},
		{X: 0, Y: 0, TileType: Ice, CornerType: CoveredTile},
	}
	(*Tiles)[Sand] = [4]*Tile{
		{X: 0, Y: 0, TileType: Sand, CornerType: StartTile},
		{X: 0, Y: 0, TileType: Sand, CornerType: MiddleTile},
		{X: 0, Y: 0, TileType: Sand, CornerType: EndTile},
		{X: 0, Y: 0, TileType: Sand, CornerType: CoveredTile},
	}
	(*Tiles)[Larva] = [4]*Tile{
		{X: 1415, Y: 1025, TileType: Larva, CornerType: StartTile},
		{X: 1642, Y: 1025, TileType: Larva, CornerType: MiddleTile},
		{X: 1868, Y: 1025, TileType: Larva, CornerType: EndTile},
		{X: 2551, Y: 1025, TileType: Larva, CornerType: CoveredTile},
	}
	(*Tiles)[Metal] = [4]*Tile{
		{X: 0, Y: 0, TileType: Metal, CornerType: StartTile},
		{X: 0, Y: 0, TileType: Metal, CornerType: MiddleTile},
		{X: 0, Y: 0, TileType: Metal, CornerType: EndTile},
		{X: 0, Y: 0, TileType: Metal, CornerType: CoveredTile},
	}
	(*Tiles)[Water] = [4]*Tile{
		{X: 1867, Y: 303, TileType: Water, CornerType: StartTile},
		{X: 1867, Y: 303, TileType: Water, CornerType: MiddleTile},
		{X: 1867, Y: 303, TileType: Water, CornerType: EndTile},
		{X: 2550, Y: 280, TileType: Water, CornerType: CoveredTile},
	}
	(*Tiles)[Wood] = [4]*Tile{
		{X: 0, Y: 0, TileType: Wood, CornerType: StartTile},
		{X: 0, Y: 0, TileType: Wood, CornerType: MiddleTile},
		{X: 0, Y: 0, TileType: Wood, CornerType: EndTile},
		{X: 0, Y: 0, TileType: Wood, CornerType: CoveredTile},
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
	Level_1_Width  = 4000
	Level_1_Height = 1440
)

func getTileType(r, g, b uint32) TileType {
	switch {
	case r == 255 && g == 0 && b == 0:
		return Larva
	case r == 0 && g == 255 && b == 0:
		return Grass
	case r == 0 && g == 0 && b == 255:
		return Water
	case r == 255 && g == 255 && b == 0:
		return Stone
	default:
		return Empty
	}
	// larva color code -> 255 ,0 ,0
	// grass color code -> 0 ,255 ,0
	// water color code -> 0 ,0 ,255
	// stone color code -> 255 ,255 ,0
}

func getColor(x, y int, levelData *ebiten.Image) (uint32, uint32, uint32) {
	color := levelData.At(x, y)
	r, g, b, _ := color.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	return r, g, b
}

func getTileInfo(x, y *int, r, g, b *uint32, levelData *ebiten.Image) *Tile {

	TileInfo := &Tile{}
	TileInfo.TileType = getTileType(*r, *g, *b)

	if levelData.Bounds().Max.X == *x {
		TileInfo.CornerType = EndTile
	} else if levelData.Bounds().Min.X == *x && getTileType(*r, *g, *b) != Grass {
		TileInfo.CornerType = StartTile
	} else if levelData.Bounds().Max.Y == *y {
		TileInfo.CornerType = CoveredTile
	} else if levelData.Bounds().Min.Y == *y {
		TileInfo.CornerType = CoveredTile
	} else {
		rPrev, gPrev, bPrev := getColor(*x-1, *y, levelData)
		rNext, gNext, bNext := getColor(*x+1, *y, levelData)
		rAbove, gAbove, bAbove := getColor(*x, *y-1, levelData)

		prevTile := getTileType(rPrev, gPrev, bPrev)
		nextTile := getTileType(rNext, gNext, bNext)
		aboveTile := getTileType(rAbove, gAbove, bAbove)

		switch {
		case aboveTile == TileInfo.TileType:
			TileInfo.CornerType = CoveredTile
		case prevTile != TileInfo.TileType:
			TileInfo.CornerType = StartTile
		case prevTile == TileInfo.TileType && nextTile == TileInfo.TileType && aboveTile == Empty:
			TileInfo.CornerType = MiddleTile
		case nextTile != TileInfo.TileType:
			TileInfo.CornerType = EndTile
		default:
			TileInfo.CornerType = CoveredTile
		}
	}

	// Lookup texture coordinates from the Tiles map
	if tilesArr, ok := (*Tiles)[TileInfo.TileType]; ok {
		for _, t := range tilesArr {
			if t.CornerType == TileInfo.CornerType {
				TileInfo.X = t.X
				TileInfo.Y = t.Y
				break
			}
		}
	}

	return TileInfo
}

// LoadLevel loads the level from the image
func (p *PlayerRuntime) LoadLevel(levelData *ebiten.Image) []Platform {
	level := []Platform{}

	fmt.Println(levelData.Bounds())

	for y := 0; y < levelData.Bounds().Dy(); y++ {
		for x := 0; x < levelData.Bounds().Dx(); x++ {
			r, g, b := getColor(x, y, levelData)
			if r != 255 || g != 255 || b != 255 {
				level = append(level,
					Platform{
						X:        float64(x * LevelTileWidth),
						Y:        float64(y * LevelTileHeight),
						Width:    LevelTileWidth,
						Height:   LevelTileHeight,
						TileInfo: *getTileInfo(&x, &y, &r, &g, &b, levelData),
					})
				// fmt.Println(r, g, b)
			}
		}
		// fmt.Println(y)
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

			op.GeoM.Translate(float64(p.X-player.Camera.Pos.X), float64(p.Y-player.Camera.Pos.Y))

			// Draw the sub-image from the tileset using coordinates from TileInfo
			screen.DrawImage(tileset.SubImage(image.Rect(p.TileInfo.X, p.TileInfo.Y, p.TileInfo.X+PixelTileWidth, p.TileInfo.Y+PixelTileHeight)).(*ebiten.Image), op)
		}
	}
}
