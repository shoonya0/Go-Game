package core

import "math"

// Mover is how a species gets around: it turns the brain's intent into
// velocity, then moves the enemy and resolves collisions. Steering and moving
// are separate calls so the brain sets this tick's velocity when it decides,
// and the move happens after, for living and dead enemies alike.
type Mover interface {
	// Steer sets the enemy's velocity for this tick's intent.
	Steer(e *EnemyRuntime, in moveIntent)
	// Move integrates the velocity and pushes the enemy out of solid tiles.
	Move(e *EnemyRuntime, qt *DynamicQuadtree)
	// Grounded reports whether the enemy rested on ground after its last move.
	Grounded(e *EnemyRuntime) bool
	// CanChase reports whether the enemy can get to the player's box from
	// where it is, so chasing them makes sense. qt may be nil (no terrain).
	CanChase(e *EnemyRuntime, target AABB, qt *DynamicQuadtree) bool
}

// moveMode is what the brain wants the body to do this tick.
type moveMode int

const (
	moveHold   moveMode = iota // stand still (attacking, hurt, dead)
	movePatrol                 // pace around the spawn point
	moveChase                  // go toward the player
)

// moveIntent is the brain's request to the Mover.
type moveIntent struct {
	mode moveMode
	dx   float64 // chase: horizontal distance to the player (positive = right)
	dy   float64 // chase: vertical distance to the player (positive = down)
}

// groundMover walks on tiles under gravity, like the player.
type groundMover struct{}

func (groundMover) Steer(e *EnemyRuntime, in moveIntent) {
	st := e.bp.Stats
	switch in.mode {
	case movePatrol:
		turnAtPathEnd(e, st.PatrolTiles*LevelTileWidth)
		e.VelX = e.patrolDir * st.WalkSpeed
	case moveChase:
		if in.dx > 0 {
			e.VelX = st.ChaseSpeed
		} else {
			e.VelX = -st.ChaseSpeed
		}
	default:
		e.VelX = 0
	}
}

// turnAtPathEnd turns a patrolling enemy around at each end of its path,
// rangePx either side of spawn.
func turnAtPathEnd(e *EnemyRuntime, rangePx float64) {
	if e.Pos.X > e.spawnX+rangePx {
		e.patrolDir = -1
	} else if e.Pos.X < e.spawnX-rangePx {
		e.patrolDir = 1
	}
}

func (m groundMover) Move(e *EnemyRuntime, qt *DynamicQuadtree) {
	m.applyGravity(e)
	m.turnAtLedge(e, qt)
	m.resolveHorizontal(e, qt)
	e.onGround = m.resolveVertical(e, qt)
}

func (groundMover) Grounded(e *EnemyRuntime) bool { return e.onGround }

// CanChase: a walker reaches only a player on ground it can walk to (enemy_reach.go).
func (groundMover) CanChase(e *EnemyRuntime, target AABB, qt *DynamicQuadtree) bool {
	return walkReaches(e, target, qt)
}

// turnAtLedge stops a patrolling enemy whose next step would put its leading
// edge past ground (solid or water) and reverses its patrol, so it paces its platform
// instead of walking off. Chasing enemies still follow the player off edges;
// one that lost the player outside its patrol path can't walk back to it past
// the ledge, so it starts a new path where it stands rather than turn toward
// its spawn and stop again every tick.
func (groundMover) turnAtLedge(e *EnemyRuntime, qt *DynamicQuadtree) {
	if qt == nil || !e.onGround || e.State != EnemyPatrol || e.VelX == 0 {
		return
	}
	probeX := e.Pos.X + e.VelX*deltaTime() // leading edge after this step
	if e.VelX > 0 {
		probeX += e.bp.BodyWidth - 1
	}
	probe := AABB{X: probeX, Y: e.Pos.Y + e.bp.BodyHeight, Width: 1, Height: 1}
	for _, obj := range qt.Retrieve(probe) {
		if plat, ok := obj.(*Platform); ok && probe.Intersects(plat.GetBounds()) {
			return // ground ahead; water holds it up too
		}
	}
	e.VelX = 0
	e.patrolDir = -e.patrolDir
	if math.Abs(e.Pos.X-e.spawnX) > e.bp.Stats.PatrolTiles*LevelTileWidth {
		e.spawnX = e.Pos.X
	}
}

// applyGravity accelerates the enemy downward, capped at the shared fall speed.
func (groundMover) applyGravity(e *EnemyRuntime) {
	e.VelY += GravityScale * physicsUnitsPerTick()
	if e.VelY > MaxFallSpeed {
		e.VelY = MaxFallSpeed
	}
}

// resolveHorizontal integrates X and pushes the enemy out of solid tiles. Hitting
// a wall reverses the patrol direction so patrolling enemies turn around at walls.
// A grounded enemy steps up onto a tile no higher than maxStepUp above its
// feet, so it can wade out of water onto the bank.
func (groundMover) resolveHorizontal(e *EnemyRuntime, qt *DynamicQuadtree) {
	e.Pos.X += e.VelX * deltaTime()
	if qt == nil {
		return
	}

	for _, obj := range qt.Retrieve(e.GetBounds()) {
		plat, ok := obj.(*Platform)
		if !ok || plat.TileInfo.TileType == Water {
			continue // water is not solid
		}
		bounds := plat.GetBounds()
		if !e.GetBounds().Intersects(bounds) {
			continue
		}
		if e.onGround && stepUp(e, bounds, qt) {
			continue
		}
		if e.VelX > 0 {
			e.Pos.X = bounds.X - e.bp.BodyWidth
		} else if e.VelX < 0 {
			e.Pos.X = bounds.X + bounds.Width
		}
		e.VelX = 0
		e.patrolDir = -e.patrolDir // bounce off the wall while patrolling
	}
}

// maxStepUp is the highest ledge (px above its feet) a walking enemy climbs
// without jumping: the depth it wades at, so it can leave the water.
const maxStepUp = waterSinkTiles*LevelTileHeight + 1

// stepUp lifts e onto the top of the tile it walked into, if that top is at
// most maxStepUp above its feet and nothing solid is in the way up there.
func stepUp(e *EnemyRuntime, tile AABB, qt *DynamicQuadtree) bool {
	if rise := e.Pos.Y + e.bp.BodyHeight - tile.Y; rise <= 0 || rise > maxStepUp {
		return false
	}
	up := e.GetBounds()
	up.Y = tile.Y - e.bp.BodyHeight
	for _, obj := range qt.Retrieve(up) {
		if plat, ok := obj.(*Platform); ok && plat.TileInfo.TileType != Water && up.Intersects(plat.GetBounds()) {
			return false
		}
	}
	e.Pos.Y = up.Y
	return true
}

// resolveVertical integrates Y, resolves floor/ceiling collisions and reports
// whether the enemy is standing on ground. Like the player, an enemy in water
// stands waterSinkTiles below the highest water tile it touches.
func (groundMover) resolveVertical(e *EnemyRuntime, qt *DynamicQuadtree) (onGround bool) {
	e.Pos.Y += e.VelY * deltaTime()
	if qt == nil {
		return false
	}

	body := e.GetBounds()
	inWater, surfaceY := false, 0.0
	for _, obj := range qt.Retrieve(body) {
		plat, ok := obj.(*Platform)
		if !ok {
			continue
		}
		bounds := plat.GetBounds()
		if !body.Intersects(bounds) {
			continue
		}
		if plat.TileInfo.TileType == Water {
			if !inWater || bounds.Y < surfaceY {
				surfaceY = bounds.Y
			}
			inWater = true
			continue
		}
		if e.VelY > 0 { // landing
			e.Pos.Y = bounds.Y - e.bp.BodyHeight
			e.VelY = 0
			onGround = true
		} else if e.VelY < 0 { // head bonk
			e.Pos.Y = bounds.Y + bounds.Height
			e.VelY = 0
		}
		body = e.GetBounds()
	}
	if floorY := surfaceY + waterSinkTiles*LevelTileHeight; inWater && e.Pos.Y+e.bp.BodyHeight >= floorY {
		e.Pos.Y = floorY - e.bp.BodyHeight
		e.VelY = 0
		onGround = true
	}
	return onGround
}
