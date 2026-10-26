package core

import "math"

// flyingMover flies in 2D without gravity and collides with solid tiles like a
// ground enemy. It chases the player wherever they go and, once they are out of
// sight, patrols at the height of the place it guards (its spawn, or where it
// lost the player). A dead flyer falls.
type flyingMover struct {
	HoverHeight float64 // chase: how far (px) above the player's center it aims
}

func (m flyingMover) Steer(e *EnemyRuntime, in moveIntent) {
	if e.State.IsDead() {
		groundMover{}.Steer(e, in) // keep the fall speed gravity built up
		return
	}
	st := e.bp.Stats
	switch in.mode {
	case movePatrol:
		turnAtPathEnd(e, m.pathRange(e))
		e.VelX = e.patrolDir * st.WalkSpeed
		_, e.VelY = towardAt(0, e.spawnY-e.Pos.Y, st.WalkSpeed)
	case moveChase:
		e.VelX, e.VelY = towardAt(in.dx, in.dy-m.HoverHeight, st.ChaseSpeed)
	default:
		e.VelX, e.VelY = 0, 0
	}
}

func (m flyingMover) Move(e *EnemyRuntime, qt *DynamicQuadtree) {
	g := groundMover{}
	if e.State.IsDead() {
		g.Move(e, qt)
		return
	}
	velX := e.VelX
	g.resolveHorizontal(e, qt)
	// A wall between a patrolling flyer and its path would pin it there,
	// turned back toward the path every tick: it starts a new path where it is.
	if e.State == EnemyPatrol && velX != 0 && e.VelX == 0 && math.Abs(e.Pos.X-e.spawnX) > m.pathRange(e) {
		e.spawnX = e.Pos.X
	}
	e.onGround = g.resolveVertical(e, qt)
}

// pathRange is how far (px) either side of the place it guards the flyer
// patrols.
func (m flyingMover) pathRange(e *EnemyRuntime) float64 {
	return e.bp.Stats.PatrolTiles * LevelTileWidth
}

func (flyingMover) Grounded(e *EnemyRuntime) bool { return e.onGround }

// CanChase: a flyer can go anywhere.
func (flyingMover) CanChase(*EnemyRuntime, AABB, *DynamicQuadtree) bool { return true }

// towardAt is the velocity that covers (dx, dy) at speed, slowed on the last
// tick so it stops on the target instead of overshooting and jittering.
func towardAt(dx, dy, speed float64) (vx, vy float64) {
	dist := math.Hypot(dx, dy)
	if dist == 0 {
		return 0, 0
	}
	if step := speed * deltaTime(); dist <= step {
		return dx / deltaTime(), dy / deltaTime()
	}
	return dx / dist * speed, dy / dist * speed
}
