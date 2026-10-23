package core

import "math"

// Draw-time motion (decisions D7): a few frames per motion look stiff, so a
// clip can add a bob, a lean and a hit shake when it is drawn. Only the sprite
// moves; the collision box stays where the physics put it.

const (
	hitShakeHz       = 30   // side-to-side shakes per second
	hitShakeDuration = 0.25 // seconds the shake takes to die out
)

// motionFX is the movement a clip adds at draw time. The zero value adds none.
type motionFX struct {
	BobPx      float64 // how far the sprite rises per bounce (screen px)
	BobHz      float64 // bounces per second
	LeanDeg    float64 // forward tilt about the feet (degrees)
	HitShakePx float64 // side-to-side shake when the clip starts, fading out (screen px)
}

// spritePose is how far a drawn sprite is moved from its resting place.
type spritePose struct {
	DX, DY float64 // offset in screen px (DY < 0 is up)
	Rot    float64 // rotation about the bottom centre (radians, clockwise)
}

// pose is the motion t seconds into the clip, for a sprite facing left or
// right. Leaning forward tilts toward the facing side.
func (m motionFX) pose(t float64, facingLeft bool) spritePose {
	var p spritePose
	if m.BobPx != 0 && m.BobHz > 0 {
		p.DY = -m.BobPx * math.Abs(math.Sin(math.Pi*m.BobHz*t))
	}
	if m.LeanDeg != 0 {
		p.Rot = m.LeanDeg * math.Pi / 180
		if facingLeft {
			p.Rot = -p.Rot
		}
	}
	if m.HitShakePx != 0 && t < hitShakeDuration {
		fade := 1 - t/hitShakeDuration
		p.DX = m.HitShakePx * fade * math.Sin(2*math.Pi*hitShakeHz*t)
	}
	return p
}
