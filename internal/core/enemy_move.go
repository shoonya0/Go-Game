package core

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
}

// groundMover walks on tiles under gravity, like the player.
type groundMover struct{}

func (groundMover) Steer(e *EnemyRuntime, in moveIntent) {
	st := e.bp.Stats
	switch in.mode {
	case movePatrol:
		// Turn at each end of the path, PatrolTiles either side of spawn.
		rangePx := st.PatrolTiles * LevelTileWidth
		if e.Pos.X > e.spawnX+rangePx {
			e.patrolDir = -1
		} else if e.Pos.X < e.spawnX-rangePx {
			e.patrolDir = 1
		}
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

func (m groundMover) Move(e *EnemyRuntime, qt *DynamicQuadtree) {
	m.applyGravity(e)
	m.resolveHorizontal(e, qt)
	e.onGround = m.resolveVertical(e, qt)
}

func (groundMover) Grounded(e *EnemyRuntime) bool { return e.onGround }

// applyGravity accelerates the enemy downward, capped at the shared fall speed.
func (groundMover) applyGravity(e *EnemyRuntime) {
	e.VelY += GravityScale * physicsUnitsPerTick()
	if e.VelY > MaxFallSpeed {
		e.VelY = MaxFallSpeed
	}
}

// resolveHorizontal integrates X and pushes the enemy out of solid tiles. Hitting
// a wall reverses the patrol direction so patrolling enemies turn around at walls.
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
		if e.VelX > 0 {
			e.Pos.X = bounds.X - e.bp.BodyWidth
		} else if e.VelX < 0 {
			e.Pos.X = bounds.X + bounds.Width
		}
		e.VelX = 0
		e.patrolDir = -e.patrolDir // bounce off the wall while patrolling
	}
}

// resolveVertical integrates Y, resolves floor/ceiling collisions and reports
// whether the enemy is standing on ground.
func (groundMover) resolveVertical(e *EnemyRuntime, qt *DynamicQuadtree) (onGround bool) {
	e.Pos.Y += e.VelY * deltaTime()
	if qt == nil {
		return false
	}

	body := e.GetBounds()
	for _, obj := range qt.Retrieve(body) {
		plat, ok := obj.(*Platform)
		if !ok || plat.TileInfo.TileType == Water {
			continue
		}
		bounds := plat.GetBounds()
		if !body.Intersects(bounds) {
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
	return onGround
}
