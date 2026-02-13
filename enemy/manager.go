package enemy

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// EnemyManager coordinates all enemies, parties, and shared learning
type EnemyManager struct {
	img          *ebiten.Image
	Animations   map[int]*Animation
	Enemies      []EnemyRuntime
	PartyManager *PartyManager
	NextID       int

	// Spawn configuration
	MaxEnemies      int
	spawnCoolDown   float64
	CurrentCoolDown float64

	// Statistics
	TotalKills       int
	TotalDeaths      int
	TotalDamageDealt float64
	TotalDamageTaken float64
}
