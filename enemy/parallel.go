package enemy

import (
	"fmt"
	"os"
	"runtime"
)

// ParallelEnemyManager extends EnemyManager with parallel processing
type ParallelEnemyManager struct {
	EnemyManager []*EnemyManager

	WorkerCount int
	// ConcurrentBrain *ConcurrentQLearningBrain
}

// DefaultParallelConfig returns sensible defaults
func DefaultParallelConfig() ParallelEnemyManager {
	workerCount := runtime.NumCPU() - 1
	if _, err := os.Stat("/.dockerenv"); err == nil {
		fmt.Println("Running in Docker")
		workerCount = 1
	}

	return ParallelEnemyManager{
		EnemyManager: make([]*EnemyManager, workerCount),
		// this is for direct machine not for docker
		WorkerCount: workerCount,
	}
}

func (p *ParallelEnemyManager) Shutdown() {
	for _, enemyManager := range p.EnemyManager {
		if enemyManager != nil {
			enemyManager.Enemies = nil
			enemyManager.PartyManager = nil
		}
	}
}
