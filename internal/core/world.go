package core

type Platform struct {
	X, Y, Width, Height float64
}

func (p *Platform) GetBounds() AABB {
	return AABB{
		X:      p.X,
		Y:      p.Y,
		Width:  p.Width,
		Height: p.Height,
	}
}

const (
	Level_1_Width  = 640
	Level_1_Height = 480
)
