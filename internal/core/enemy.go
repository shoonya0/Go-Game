package core

import "github.com/hajimehoshi/ebiten/v2"

type Enemy interface {
	Collider
	Update(dt float64, player *PlayerRuntime, level []Platform)
	Draw(screen *ebiten.Image)
	IsDead() bool
}
