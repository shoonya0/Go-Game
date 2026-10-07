package core

import "math"

// UpdateEnemy advances one enemy by a single tick: it decides behaviour from the
// player's position, then integrates physics and resolves collisions against the
// world. It is the shared entry point for every enemy species — the per-species
// differences all live in the blueprint (enemy_types.go).
//
// It reuses the same gravity constants (GravityScale, MaxFallSpeed) and quadtree
// collision approach as the player, so enemies fall and stand on the same tiles.
func UpdateEnemy(e *EnemyRuntime, target *PlayerRuntime, qt *DynamicQuadtree) {
	e.PrevState = e.State
	if e.attackTimer > 0 {
		e.attackTimer -= deltaTime()
	}

	// A dead enemy keeps falling to rest on the ground but takes no actions.
	if e.State.IsDead() {
		e.VelX = 0
		e.applyGravity()
		e.resolveHorizontal(qt)
		e.onGround = e.resolveVertical(qt)
		return
	}

	e.decideBehavior(target)
	e.applyGravity()
	e.resolveHorizontal(qt)
	e.onGround = e.resolveVertical(qt)
}

// decideBehavior runs the enemy's state machine for this tick: attack when the
// player is in reach, chase when they are within detection range, otherwise
// patrol. It sets the behavioural state, horizontal velocity and facing.
func (e *EnemyRuntime) decideBehavior(target *PlayerRuntime) {
	// Let a hurt reaction play out before re-deciding, so hits interrupt cleanly.
	if e.State.IsHurt() {
		e.VelX = 0
		return
	}

	st := e.bp.Stats
	dx := centerX(target.GetBounds()) - centerX(e.GetBounds())
	dist := math.Abs(dx)

	switch {
	case dist <= st.AttackRange:
		e.attack(target, dx)
	case dist <= st.DetectRange:
		e.chase(dx)
	default:
		e.patrol()
	}
}

// attack faces and strikes the player, standing still. Damage is gated by the
// attack cooldown so contact deals one hit per AttackCooldown seconds rather than
// once per frame.
func (e *EnemyRuntime) attack(target *PlayerRuntime, dx float64) {
	e.faceToward(dx)
	e.VelX = 0
	e.setState(EnemyAttack)

	if e.attackTimer <= 0 {
		target.ApplyDamage(e.bp.Stats.Damage)
		e.attackTimer = e.bp.Stats.AttackCooldown
	}
}

// chase moves the enemy toward the player at chase speed.
func (e *EnemyRuntime) chase(dx float64) {
	e.setState(EnemyAlert)
	e.faceToward(dx)
	if dx > 0 {
		e.VelX = e.bp.Stats.ChaseSpeed
	} else {
		e.VelX = -e.bp.Stats.ChaseSpeed
	}
}

// patrol paces the enemy back and forth within PatrolTiles of its spawn point,
// reversing at each end of the path.
func (e *EnemyRuntime) patrol() {
	rangePx := e.bp.Stats.PatrolTiles * LevelTileWidth
	if e.Pos.X > e.spawnX+rangePx {
		e.patrolDir = -1
	} else if e.Pos.X < e.spawnX-rangePx {
		e.patrolDir = 1
	}

	e.VelX = e.patrolDir * e.bp.Stats.WalkSpeed
	e.faceToward(e.patrolDir)
	e.setState(EnemyPatrol)
}

// applyGravity accelerates the enemy downward, capped at the shared fall speed.
func (e *EnemyRuntime) applyGravity() {
	e.VelY += GravityScale * physicsUnitsPerTick()
	if e.VelY > MaxFallSpeed {
		e.VelY = MaxFallSpeed
	}
}

// resolveHorizontal integrates X and pushes the enemy out of solid tiles. Hitting
// a wall reverses the patrol direction so patrolling enemies turn around at walls.
func (e *EnemyRuntime) resolveHorizontal(qt *DynamicQuadtree) {
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
func (e *EnemyRuntime) resolveVertical(qt *DynamicQuadtree) (onGround bool) {
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

// centerX returns the horizontal center of a box.
func centerX(b AABB) float64 { return b.X + b.Width/2 }
