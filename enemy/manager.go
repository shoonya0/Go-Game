package enemy

// EnemyManager coordinates all enemies, parties, and shared learning
type EnemyManager struct {
	Enemies      []*EnemyRuntime
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
