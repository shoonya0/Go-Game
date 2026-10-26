package core

import "math"

// Ground reach (#34): a walking enemy chases only a player it can walk to. It
// follows the ground from its own feet toward the player, the way it would
// walk: no gaps, no walls, no step higher than maxStepUp, and no larva unless
// it already stands in some. The player must stand on that ground at the end,
// or be in the air above it.

const (
	reachStep        = LevelTileWidth / 3  // px between ground samples along the walk
	playerFloorReach = 4 * LevelTileHeight // how far below a jumping player's feet their floor is looked for
)

// groundIn is where a walker's feet come to rest inside probe: the top of the
// highest solid tile there, or waterSinkTiles under the highest water tile,
// and that tile's type. A tile that rises above the probe (a wall) counts at
// its top. ok is false when the probe holds no ground.
func groundIn(qt *DynamicQuadtree, probe AABB) (top float64, tile TileType, ok bool) {
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
		if t < top {
			top, tile = t, plat.TileInfo.TileType
		}
	}
	return top, tile, !math.IsInf(top, 1)
}

// feetProbe is the strip under a box's feet: its whole width, 2 px deep, so a
// body held up by part of its width still finds its ground.
func feetProbe(b AABB) AABB { return AABB{X: b.X, Y: b.Y + b.Height, Width: b.Width, Height: 2} }

// walkReaches reports whether e can walk from where it stands to the player's
// box. Without a world it can always reach them.
func walkReaches(e *EnemyRuntime, target AABB, qt *DynamicQuadtree) bool {
	if qt == nil {
		return true
	}
	eb := e.GetBounds()
	floor, standsOn, ok := groundIn(qt, feetProbe(eb))
	if !ok {
		return false // in the air: it can't walk anywhere
	}
	to := centerX(target)
	dir := math.Copysign(1, to-centerX(eb))
	for x := centerX(eb) + dir*reachStep; (to-x)*dir > 0; x += dir * reachStep {
		column := AABB{X: x, Y: floor - e.bp.BodyHeight, Width: 1, Height: e.bp.BodyHeight + maxStepUp + 1}
		top, tile, ok := groundIn(qt, column)
		if !ok || math.Abs(top-floor) > maxStepUp || (tile == Larva && standsOn != Larva) {
			return false // a gap, a wall, a ledge or larva on the way
		}
		floor = top
	}
	under := feetProbe(target)
	under.Height = playerFloorReach
	playerFloor, _, ok := groundIn(qt, under)
	return ok && math.Abs(playerFloor-floor) <= maxStepUp
}
