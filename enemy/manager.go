package enemy

import "player/internal/core"

// EnemyManager coordinates all enemies, parties, and shared learning
type EnemyManager struct {
	ID      string
	Enemies []EnemyRuntime
	nextID  int // internal counter for generating enemy IDs

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

func (em *EnemyManager) InitEnemyManager(id string) EnemyManager {
	return EnemyManager{
		ID: id,
		// 	img:              core.LoadImage(enemySpriteSheetPath),
		// 	Animations:       InitEnemyAnimations(),
		Enemies:          []EnemyRuntime{},
		nextID:           0,
		MaxEnemies:       10,
		spawnCoolDown:    0,
		CurrentCoolDown:  0,
		TotalKills:       0,
		TotalDeaths:      0,
		TotalDamageDealt: 0,
		TotalDamageTaken: 0,
	}
}

func (em *EnemyManager) Update(player *core.PlayerRuntime, qt *core.DynamicQuadtree) {
	for i := range em.Enemies {
		em.Enemies[i].Update(player, qt)
	}
}

func (em *EnemyManager) UpdateAnimations(animations map[int]Animation) {
	for i := range em.Enemies {
		em.Enemies[i].UpdateEnemyAnimation(&animations)
	}
}
