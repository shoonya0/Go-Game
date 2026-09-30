package core

// GameState is the top-level mode of the game (menu, playing, game over).
type GameState int

const (
	ModeMenu GameState = iota
	ModePlaying
	ModeGameOver
)
