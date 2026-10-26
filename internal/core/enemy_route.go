package core

import "math"

// Flyer routes and line of sight (#35). A flyer whose straight flight to the
// player is blocked by tiles searches a grid of half-tile cells around itself
// for the way that gets nearest the player, and heads for the farthest point
// on it it can fly to in a straight line. No enemy strikes through a tile.

const (
	routeCell   = LevelTileWidth / 2 // grid cell (px); tiles cover whole cells across
	routeRadius = 24                 // cells searched either side of the flyer
	sweepStep   = 10.0               // px between samples along a straight flight or sight line
	lookAhead   = 8                  // route cells checked for a straight flight
)

// solidsIn lists the boxes of the solid (non-water) tiles touching area.
func solidsIn(qt *DynamicQuadtree, area AABB) []AABB {
	var solids []AABB
	for _, obj := range qt.Retrieve(area) {
		if plat, ok := obj.(*Platform); ok && plat.TileInfo.TileType != Water {
			if b := plat.GetBounds(); b.Intersects(area) {
				solids = append(solids, b)
			}
		}
	}
	return solids
}

// sweepSteps is how many samples a straight line of length dist takes.
func sweepSteps(dist float64) int { return max(1, int(math.Ceil(dist/sweepStep))) }

// clearSight reports whether the straight line between the centers of a and b
// passes through no solid tile. Without a world nothing blocks it.
func clearSight(qt *DynamicQuadtree, a, b AABB) bool {
	if qt == nil {
		return true
	}
	ax, ay, bx, by := centerX(a), centerY(a), centerX(b), centerY(b)
	solids := solidsIn(qt, AABB{X: min(ax, bx), Y: min(ay, by), Width: math.Abs(bx-ax) + 1, Height: math.Abs(by-ay) + 1})
	n := sweepSteps(math.Hypot(bx-ax, by-ay))
	for i := 0; i <= n; i++ {
		f := float64(i) / float64(n)
		x, y := ax+(bx-ax)*f, ay+(by-ay)*f
		for _, s := range solids {
			if s.Contains(x, y) {
				return false
			}
		}
	}
	return true
}

// sweepClear reports whether body can fly (dx, dy) in a straight line without
// touching any of solids.
func sweepClear(body AABB, dx, dy float64, solids []AABB) bool {
	n := sweepSteps(math.Hypot(dx, dy))
	for i := 1; i <= n; i++ {
		f := float64(i) / float64(n)
		b := body
		b.X += dx * f
		b.Y += dy * f
		for _, s := range solids {
			if b.Intersects(s) {
				return false
			}
		}
	}
	return true
}

// routeGrid marks which cells of a square around the flyer hold a solid tile.
type routeGrid struct {
	ox, oy  float64 // world position of cell (0, 0), on the tile grid
	n       int     // cells per side
	blocked []bool
	spanW   int // cells a body covers across
	spanH   int // and down
	padX    float64
	padY    float64 // body offset inside its cells (centred)
}

func newRouteGrid(body AABB, solids []AABB) *routeGrid {
	g := &routeGrid{
		ox:    math.Floor(body.X/routeCell)*routeCell - routeRadius*routeCell,
		oy:    math.Floor(body.Y/routeCell)*routeCell - routeRadius*routeCell,
		n:     2*routeRadius + 1,
		spanW: int(math.Ceil(body.Width / routeCell)),
		spanH: int(math.Ceil(body.Height / routeCell)),
	}
	g.padX = (float64(g.spanW)*routeCell - body.Width) / 2
	g.padY = (float64(g.spanH)*routeCell - body.Height) / 2
	g.blocked = make([]bool, g.n*g.n)
	for _, s := range solids {
		c0, c1, r0, r1 := g.cells(s)
		for r := max(r0, 0); r <= min(r1, g.n-1); r++ {
			for c := max(c0, 0); c <= min(c1, g.n-1); c++ {
				g.blocked[r*g.n+c] = true
			}
		}
	}
	return g
}

// cells is the range of cells box overlaps.
func (g *routeGrid) cells(b AABB) (c0, c1, r0, r1 int) {
	c0 = int(math.Floor((b.X - g.ox) / routeCell))
	c1 = int(math.Ceil((b.X+b.Width-g.ox)/routeCell)) - 1
	r0 = int(math.Floor((b.Y - g.oy) / routeCell))
	r1 = int(math.Ceil((b.Y+b.Height-g.oy)/routeCell)) - 1
	return
}

// fits reports whether a body with its top-left cell at (c, r) stays in the
// grid and out of every blocked cell.
func (g *routeGrid) fits(c, r int) bool {
	if c < 0 || r < 0 || c+g.spanW > g.n || r+g.spanH > g.n {
		return false
	}
	for y := r; y < r+g.spanH; y++ {
		for x := c; x < c+g.spanW; x++ {
			if g.blocked[y*g.n+x] {
				return false
			}
		}
	}
	return true
}

// pos is the world position of a body whose top-left cell is (c, r).
func (g *routeGrid) pos(c, r int) (x, y float64) {
	return g.ox + float64(c)*routeCell + g.padX, g.oy + float64(r)*routeCell + g.padY
}

// flyRoute is where a flyer should head to cover (dx, dy) from where it is:
// ok is false when it can fly straight there (or no route gets it closer),
// otherwise (wx, wy) is the offset of the next point on the route.
func flyRoute(e *EnemyRuntime, dx, dy float64, qt *DynamicQuadtree) (wx, wy float64, ok bool) {
	if qt == nil {
		return 0, 0, false
	}
	body := e.GetBounds()
	reach := float64(routeRadius*routeCell + 2*routeCell)
	solids := solidsIn(qt, AABB{X: body.X - reach, Y: body.Y - reach, Width: body.Width + 2*reach, Height: body.Height + 2*reach})
	if sweepClear(body, dx, dy, solids) {
		return 0, 0, false
	}

	g := newRouteGrid(body, solids)
	start := g.n*routeRadius + routeRadius // the flyer's own cell
	tx, ty := body.X+dx, body.Y+dy
	prev := make([]int, len(g.blocked))
	for i := range prev {
		prev[i] = -1
	}
	prev[start] = start
	best, bestDist := start, math.Inf(1)
	queue := []int{start}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		c, r := cur%g.n, cur/g.n
		if x, y := g.pos(c, r); math.Hypot(tx-x, ty-y) < bestDist {
			best, bestDist = cur, math.Hypot(tx-x, ty-y)
		}
		for _, d := range [8][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}} {
			nc, nr := c+d[0], r+d[1]
			// Diagonal steps only past two free sides, so a body can't clip a corner.
			if !g.fits(nc, nr) || (d[0] != 0 && d[1] != 0 && (!g.fits(c+d[0], r) || !g.fits(c, r+d[1]))) {
				continue
			}
			if next := nr*g.n + nc; prev[next] < 0 {
				prev[next] = cur
				queue = append(queue, next)
			}
		}
	}
	if best == start {
		return 0, 0, false
	}

	// Walk back from the best cell, then head for the farthest of the next
	// lookAhead cells it can fly to in a straight line.
	path := []int{}
	for at := best; at != start; at = prev[at] {
		path = append(path, at)
	}
	for _, at := range path[max(0, len(path)-lookAhead):] {
		x, y := g.pos(at%g.n, at/g.n)
		if sweepClear(body, x-body.X, y-body.Y, solids) {
			return x - body.X, y - body.Y, true
		}
	}
	x, y := g.pos(path[len(path)-1]%g.n, path[len(path)-1]/g.n) // the next cell
	return x - body.X, y - body.Y, true
}

// routeReplan is how often (s) a chasing flyer searches its route again; in
// between it keeps heading for the point it found.
const routeReplan = 0.25

// flyerRoute is a flyer's cached route step.
type flyerRoute struct {
	x, y float64 // world position of the point it heads for
	ok   bool    // false: it flies straight at the player
	wait float64 // seconds until the next search
}

// routeStep is flyRoute, searched at most every routeReplan seconds or once
// the flyer reaches the point it was heading for.
func (e *EnemyRuntime) routeStep(dx, dy float64, qt *DynamicQuadtree) (wx, wy float64, ok bool) {
	r := &e.route
	r.wait -= deltaTime()
	if r.wait <= 0 || (r.ok && math.Hypot(r.x-e.Pos.X, r.y-e.Pos.Y) < 1) {
		wx, wy, ok := flyRoute(e, dx, dy, qt)
		*r = flyerRoute{x: e.Pos.X + wx, y: e.Pos.Y + wy, ok: ok, wait: routeReplan}
	}
	if !r.ok {
		return 0, 0, false
	}
	return r.x - e.Pos.X, r.y - e.Pos.Y, true
}
