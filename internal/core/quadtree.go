package core

const (
	MaxObjects = 10
	MaxLevels  = 5
)

// Quadtree represents a spatial partition for efficient collision detection
type Quadtree struct {
	Level   int
	Bounds  AABB
	Objects []Collider
	Nodes   [4]*Quadtree // 0: NE, 1: NW, 2: SW, 3: SE
}

// NewQuadtree creates a new Quadtree node
func NewQuadtree(level int, bounds AABB) *Quadtree {
	return &Quadtree{
		Level:   level,
		Bounds:  bounds,
		Objects: make([]Collider, 0),
	}
}

// Clear recursively clears the quadtree
func (q *Quadtree) Clear() {
	q.Objects = q.Objects[:0]
	for i := 0; i < 4; i++ {
		if q.Nodes[i] != nil {
			q.Nodes[i].Clear()
			q.Nodes[i] = nil
		}
	}
}

// split divides the node into 4 subnodes
func (q *Quadtree) split() {
	subWidth := q.Bounds.Width / 2
	subHeight := q.Bounds.Height / 2
	x := q.Bounds.X
	y := q.Bounds.Y

	q.Nodes[0] = NewQuadtree(q.Level+1, AABB{x + subWidth, y, subWidth, subHeight})             // NE
	q.Nodes[1] = NewQuadtree(q.Level+1, AABB{x, y, subWidth, subHeight})                        // NW
	q.Nodes[2] = NewQuadtree(q.Level+1, AABB{x, y + subHeight, subWidth, subHeight})            // SW
	q.Nodes[3] = NewQuadtree(q.Level+1, AABB{x + subWidth, y + subHeight, subWidth, subHeight}) // SE
}

// getIndex determines which node the object belongs to. -1 means it doesn't fit in a child node (overlaps split)
func (q *Quadtree) getIndex(rect AABB) int {
	// index is -1 if the object doesn't fit in a child node (overlaps split)
	index := -1
	verticalMidpoint := q.Bounds.X + (q.Bounds.Width / 2)
	horizontalMidpoint := q.Bounds.Y + (q.Bounds.Height / 2)

	topQuadrant := (rect.Y < horizontalMidpoint && rect.Y+rect.Height < horizontalMidpoint)
	bottomQuadrant := (rect.Y > horizontalMidpoint)

	// Check left quadrants
	if rect.X < verticalMidpoint && rect.X+rect.Width < verticalMidpoint {
		if topQuadrant {
			index = 1 // NW
		} else if bottomQuadrant {
			index = 2 // SW
		}
	} else if rect.X > verticalMidpoint { // Check right quadrants
		if topQuadrant {
			index = 0 // NE
		} else if bottomQuadrant {
			index = 3 // SE
		}
	}

	return index
}

// Insert adds an object to the quadtree
func (q *Quadtree) Insert(obj Collider) {
	if q.Nodes[0] != nil {
		index := q.getIndex(obj.GetBounds())
		if index != -1 {
			q.Nodes[index].Insert(obj)
			return
		}
	}

	q.Objects = append(q.Objects, obj)

	if len(q.Objects) > MaxObjects && q.Level < MaxLevels {
		if q.Nodes[0] == nil {
			q.split()
		}

		i := 0
		for i < len(q.Objects) {
			index := q.getIndex(q.Objects[i].GetBounds())
			if index != -1 {
				objToRemove := q.Objects[i]
				// Remove object from current node
				q.Objects = append(q.Objects[:i], q.Objects[i+1:]...)
				// Insert into child
				q.Nodes[index].Insert(objToRemove)
			} else {
				i++
			}
		}
	}
}

// Retrieve appends every object that could overlap rect to returnObjects and
// returns the extended slice. Objects at this node are always candidates; the
// search then descends into the child rect fits in, or, when rect straddles a
// split (index == -1), into every child whose bounds intersect it.
func (q *Quadtree) Retrieve(returnObjects []Collider, rect AABB) []Collider {
	returnObjects = append(returnObjects, q.Objects...)

	if q.Nodes[0] == nil {
		return returnObjects
	}

	if index := q.getIndex(rect); index != -1 {
		return q.Nodes[index].Retrieve(returnObjects, rect)
	}

	for i := 0; i < 4; i++ {
		if q.Nodes[i].Bounds.Intersects(rect) {
			returnObjects = q.Nodes[i].Retrieve(returnObjects, rect)
		}
	}
	return returnObjects
}

// dynamic quadtree
// findLeafFor descends the tree to find the leaf node where rect would be placed
func (q *Quadtree) findLeafFor(rect AABB) *Quadtree {
	node := q
	for {
		if node.Nodes[0] == nil {
			return node
		}
		index := node.getIndex(rect)
		if index == -1 {
			return node
		}
		node = node.Nodes[index]
	}
}

// remove removes obj from this node only (does not search children/parents)
func (q *Quadtree) remove(obj Collider) bool {
	for i := range q.Objects {
		if q.Objects[i] == obj {
			q.Objects = append(q.Objects[:i], q.Objects[i+1:]...)
			return true
		}
	}
	return false
}

// removeRecursive searches the subtree to remove the object (handles stale index case)
func (q *Quadtree) removeRecursive(obj Collider) bool {
	if q.remove(obj) {
		return true
	}
	if q.Nodes[0] != nil {
		for i := 0; i < 4; i++ {
			if q.Nodes[i].removeRecursive(obj) {
				return true
			}
		}
	}
	return false
}

// ------------------------ Dynamic Quadtree Wrapper ------------------------

// DynamicQuadtree keeps an index of object->node to allow O(log N) updates
type DynamicQuadtree struct {
	Root  *Quadtree
	index map[Collider]*Quadtree
}

// NewDynamicQuadtree initializes a dynamic quadtree with given bounds
func NewDynamicQuadtree(bounds AABB) *DynamicQuadtree {
	return &DynamicQuadtree{
		Root:  NewQuadtree(0, bounds),
		index: make(map[Collider]*Quadtree),
	}
}

// Clear removes all objects and resets the index
func (dq *DynamicQuadtree) Clear() {
	dq.Root.Clear()
	dq.index = make(map[Collider]*Quadtree)
}

// Insert inserts obj and records its leaf node for future updates
func (dq *DynamicQuadtree) Insert(obj Collider) {
	dq.Root.Insert(obj)
	leaf := dq.Root.findLeafFor(obj.GetBounds())
	dq.index[obj] = leaf
}

// Update updates the position of obj by removing it from its current leaf and reinserting
func (dq *DynamicQuadtree) Update(obj Collider) {
	leaf, ok := dq.index[obj]
	if !ok {
		// Not tracked yet, treat as insert
		dq.Insert(obj)
		return
	}

	// Try to remove from the tracked leaf
	// If the node split since we last checked, the object might be in a child now.
	if !leaf.remove(obj) {
		// Fallback: search the subtree of the last known location
		if !leaf.removeRecursive(obj) {
			// If still not found, it might have been removed or logic error.
			// Proceed to re-insert to ensure it exists.
		}
	}

	// Reinsert from root based on new bounds
	dq.Root.Insert(obj)
	newLeaf := dq.Root.findLeafFor(obj.GetBounds())
	dq.index[obj] = newLeaf
}

// UpdateAll updates a slice of objects efficiently
func (dq *DynamicQuadtree) UpdateAll(objs []Collider) {
	for _, o := range objs {
		dq.Update(o)
	}
}

// Remove deletes an object from the tree and index
func (dq *DynamicQuadtree) Remove(obj Collider) bool {
	leaf, ok := dq.index[obj]
	if !ok {
		return false
	}
	deleted := leaf.remove(obj)
	delete(dq.index, obj)
	return deleted
}

// Retrieve returns the potential colliders for rect, allocating a fresh slice.
func (dq *DynamicQuadtree) Retrieve(rect AABB) []Collider {
	return dq.Root.Retrieve(nil, rect)
}

// RetrieveInto appends the potential colliders for rect to buf and returns the
// extended slice. Pass a reused (zero-length) buffer in a hot loop to avoid a
// per-query allocation.
func (dq *DynamicQuadtree) RetrieveInto(buf []Collider, rect AABB) []Collider {
	return dq.Root.Retrieve(buf, rect)
}
