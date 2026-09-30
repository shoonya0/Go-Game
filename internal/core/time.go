package core

import "github.com/hajimehoshi/ebiten/v2"

// defaultTPS is the fallback tick rate used when Ebitengine reports a
// non-positive value (e.g. before the game loop has started).
const defaultTPS = 60.0

// tps returns the current ticks-per-second, guarded against a non-positive
// reading so callers can always divide by it safely.
func tps() float64 {
	if t := float64(ebiten.TPS()); t > 0 {
		return t
	}
	return defaultTPS
}

// deltaTime returns the duration of a single tick in seconds (1 / TPS). It is
// the shared time step for both physics integration and animation timing.
func deltaTime() float64 {
	return 1.0 / tps()
}
