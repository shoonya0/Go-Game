package core

import (
	"math"
	"testing"
)

// Pool layout for the water tests: banks of grass on either side of a pool
// two tiles deep, whose surface is level with the banks.
const (
	poolLeft  = 600.0
	poolRight = 900.0
)

// waterDepthFloor is where a body in the pool stands: as deep as the player,
// waterSinkTiles below the surface.
const waterDepthFloor = waveGroundY + waterSinkTiles*LevelTileHeight

// poolWorld is waveWorld with the pool cut into its floor.
func poolWorld() *DynamicQuadtree {
	qt := NewDynamicQuadtree(AABB{X: 0, Y: 0, Width: 4000, Height: 2000})
	grass, water := Tile{TileType: Grass}, Tile{TileType: Water}
	for x := 0.0; x < 1500; x += LevelTileWidth {
		if x < poolLeft || x >= poolRight {
			for row := 0.0; row < 3; row++ {
				qt.Insert(&Platform{X: x, Y: waveGroundY + row*LevelTileHeight, Width: LevelTileWidth, Height: LevelTileHeight, TileInfo: grass})
			}
			continue
		}
		qt.Insert(&Platform{X: x, Y: waveGroundY, Width: LevelTileWidth, Height: LevelTileHeight, TileInfo: water})
		qt.Insert(&Platform{X: x, Y: waveGroundY + LevelTileHeight, Width: LevelTileWidth, Height: LevelTileHeight, TileInfo: water})
		qt.Insert(&Platform{X: x, Y: waveGroundY + 2*LevelTileHeight, Width: LevelTileWidth, Height: LevelTileHeight, TileInfo: grass})
	}
	return qt
}

// TestGroundEnemyStandsInWaterLikePlayer: an enemy that falls into water stands
// at the player's depth, not on the bottom of the pool (#30).
func TestGroundEnemyStandsInWaterLikePlayer(t *testing.T) {
	qt := poolWorld()
	e := newTestEnemy(Position{X: 730, Y: waveGroundY - 300})
	p := farPlayer(e)

	for i := 0; i < 3*int(tps()) && !e.onGround; i++ {
		UpdateEnemy(e, p, qt)
	}
	if feet := e.Pos.Y + e.bp.BodyHeight; math.Abs(feet-waterDepthFloor) > 0.01 || !e.onGround {
		t.Errorf("feet at y=%v onGround=%v, want standing in the water at y=%v", feet, e.onGround, waterDepthFloor)
	}
}

// TestGroundEnemyPatrolsThroughWater: a patrol path across a pool wades
// through it and climbs out onto both banks (#30).
func TestGroundEnemyPatrolsThroughWater(t *testing.T) {
	qt := poolWorld()
	e := newTestEnemy(Position{X: 730, Y: waterDepthFloor - testEnemyBlueprint.BodyHeight})
	p := farPlayer(e) // patrol path 730 ± 4 tiles reaches both banks

	var onRight, onLeft bool
	for i := 0; i < 12*int(tps()); i++ {
		UpdateEnemy(e, p, qt)
		feet := e.Pos.Y + e.bp.BodyHeight
		if feet > waterDepthFloor+0.01 {
			t.Fatalf("tick %d: feet at y=%v, sank below the wading depth %v", i, feet, waterDepthFloor)
		}
		onBank := math.Abs(feet-waveGroundY) < 0.01
		onRight = onRight || onBank && e.Pos.X >= poolRight
		onLeft = onLeft || onBank && e.Pos.X+e.bp.BodyWidth <= poolLeft
	}
	if !onRight || !onLeft {
		t.Errorf("climbed out on the right bank=%v, on the left bank=%v; want both", onRight, onLeft)
	}
}
