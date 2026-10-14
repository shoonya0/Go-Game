---
name: go-spatial-partitioning
description: Build and reason about spatial indexes in Go — a quadtree for O(log N) range queries and viewport culling, with a dynamic object→leaf index so moving objects re-index cheaply. Use when you need fast neighborhood or range lookups over many moving objects, collision detection, viewport/broad-phase culling, geo or 2D spatial queries, or when replacing an O(n^2) all-pairs scan with a logarithmic structure.
---

# Go spatial partitioning (quadtree)

A quadtree recursively subdivides 2D space into four quadrants so that a range query
touches only the objects near the query region, not every object. It turns an
O(n²) all-pairs scan into roughly O(n log n) broad-phase work.

## When to reach for it

- Collision / overlap detection among many objects.
- "What is near this rectangle / point?" queries (viewport culling, proximity, geo-ish
  lookups in a flat plane).
- Any hot loop currently comparing every object against every other object.

## Core structure

From the engine's `internal/core/quadtree.go`:

```go
type Quadtree struct {
    Level   int
    Bounds  AABB
    Objects []Collider
    Nodes   [4]*Quadtree // NE, NW, SW, SE
}
```

- **Insert**: if the node has children and the object fits fully in one, descend; otherwise
  keep it here. Split once `len(Objects) > MaxObjects` and `Level < MaxLevels`.
- **Retrieve(rect)**: collect this node's objects, then descend into the child the rect fits
  in — or, if it straddles a boundary, into every child whose bounds intersect the rect.

Tune `MaxObjects` (e.g. 10) and `MaxLevels` (e.g. 5) to the object density.

## Making it dynamic (cheap updates for moving objects)

A naive quadtree rebuilds every frame. Keep an index from each object to its current leaf
so a move is a local remove + reinsert instead of a full rebuild:

```go
type DynamicQuadtree struct {
    Root  *Quadtree
    index map[Collider]*Quadtree // object -> leaf it lives in
}

func (dq *DynamicQuadtree) Update(obj Collider) {
    leaf := dq.index[obj]        // O(1) lookup
    leaf.remove(obj)             // pull from old leaf
    dq.Root.Insert(obj)          // reinsert by new bounds
    dq.index[obj] = dq.Root.findLeafFor(obj.GetBounds())
}
```

Guard against a stale leaf (the node may have split since): fall back to a recursive remove
from the last-known subtree before reinserting.

## Broad phase + narrow phase

The quadtree is the **broad phase**: it returns *candidates*. Always confirm with an exact
test (the **narrow phase**), e.g. AABB intersection:

```go
func (a AABB) Intersects(b AABB) bool {
    return a.X < b.X+b.Width && a.X+a.Width > b.X &&
           a.Y < b.Y+b.Height && a.Y+a.Height > b.Y
}
```

## Checklist

- [ ] Objects implement a small interface (`GetBounds() AABB`) so the tree is generic.
- [ ] `Retrieve` appends into a caller-supplied slice to avoid per-query allocations.
- [ ] Straddling objects (index == -1) are checked against all intersecting children.
- [ ] Moves go through the dynamic index, not a full rebuild.
- [ ] Broad-phase candidates are always confirmed by an exact narrow-phase test.

## Pitfalls

- Returning only the target quadrant and missing objects that straddle a split boundary.
- Rebuilding the whole tree each tick when a dynamic index would be O(log N) per move.
- Allocating a fresh result slice every query in the hot path.

## Backend relevance

The same instinct powers spatial and geo indexes, R-trees, and the general habit of picking
a data structure that turns a linear or quadratic scan into a logarithmic lookup — the
difference between a query that scales and one that does not.
