package core

import "math"

// Ground reach (#34): a walking enemy chases only a player it can walk to. It
// follows the ground from its own feet toward the player, the way it would
// walk: no gaps, no walls, no step higher than maxStepUp. The player must
// stand on that ground at the end, or be in the air above it.

const (
	reachStep        = LevelTileWidth / 3  // px between ground samples along the walk
	playerFloorReach = 4 * LevelTileHeight // how far below a jumping player's feet their floor is looked for
)

// groundBelow is where a walker's feet come to rest in the column at x between
// y and y+depth: the top of the highest solid tile there, or waterSinkTiles
// under the highest water tile. A tile that rises above y (a wall) counts at
// its top. ok is false when the column holds no ground.
func groundBelow(qt *DynamicQuadtree, x, y, depth float64) (top float64, ok bool) {
	probe := AABB{X: x, Y: y, Width: 1, Height: depth}
	top = math.Inf(1)
	for _, obj := range qt.Retrieve(probe) {
		plat, isPlat := obj.(*Platform)
		if !isPlat {
			continue
		}
		b := plat.GetBounds()
		if !probe.Intersects(b) {
			continue
		}
		t := b.Y
		if plat.TileInfo.TileType == Water {
			t += waterSinkTiles * LevelTileHeight
		}
		top = min(top, t)
	}
	return top, !math.IsInf(top, 1)
}

// walkReaches reports whether e can walk from where it stands to the player's
// box. Without a world it can always reach them.
func walkReaches(e *EnemyRuntime, target AABB, qt *DynamicQuadtree) bool {
	if qt == nil {
		return true
	}
	eb := e.GetBounds()
	floor, ok := groundBelow(qt, centerX(eb), eb.Y+eb.Height, 2)
	if !ok {
		return false // in the air: it can't walk anywhere
	}
	to := centerX(target)
	dir := math.Copysign(1, to-centerX(eb))
	for x := centerX(eb) + dir*reachStep; (to-x)*dir > 0; x += dir * reachStep {
		top, ok := groundBelow(qt, x, floor-e.bp.BodyHeight, e.bp.BodyHeight+maxStepUp+1)
		if !ok || math.Abs(top-floor) > maxStepUp {
			return false // a gap, a wall or a ledge on the way
		}
		floor = top
	}
	playerFloor, ok := groundBelow(qt, to, target.Y+target.Height, playerFloorReach)
	return ok && math.Abs(playerFloor-floor) <= maxStepUp
}
