package core

// AABB represents an Axis-Aligned Bounding Box
type AABB struct {
	X, Y, Width, Height float64
}

// Intersects checks if two AABBs intersect
func (a AABB) Intersects(b AABB) bool {
	return a.X < b.X+b.Width &&
		a.X+a.Width > b.X &&
		a.Y < b.Y+b.Height &&
		a.Y+a.Height > b.Y
}

// Collider is an interface for any object that has a bounding box
type Collider interface {
	// here GetBounds returns the bounding box of the collider
	GetBounds() AABB
}
