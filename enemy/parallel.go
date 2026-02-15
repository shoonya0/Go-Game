package enemy

import (
	"fmt"
	"os"
	"player/internal/core"
	"runtime"
	"sync"
)

// ParallelEnemyManager extends EnemyManager with parallel processing
type ParallelEnemyManager struct {
	EnemyManager []EnemyManager
	// parrallel processing
	wg          sync.WaitGroup
	mutex       sync.Mutex
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
		// we do not want to copy our mutex and wait group so we make it a pointer
		EnemyManager: make([]EnemyManager, workerCount),
		// this is for direct machine not for docker
		WorkerCount: workerCount,
	}
}

func (p *ParallelEnemyManager) Shutdown() {
	p.EnemyManager = nil
}

func (em *EnemyManager) FindEnemyPosition(enemy EnemyRuntime) core.Position {
	// find position of enemy on the map
	return core.Position{
		X: 100,
		Y: 100,
	}
}

func (em *ParallelEnemyManager) AddEnemyToLevel(level []core.Platform) {
	// find position of enemy on the map
	fmt.Println("Adding enemies to level")
	fmt.Println(em.EnemyManager)
	if em == nil || em.EnemyManager == nil {
		fmt.Println("Enemy Manager is nil")
		return
	}

	for _, platform := range level {
		if platform.TileInfo.TileType == core.EnemyBasic {
			// 	em.EnemyManager[0].Enemies = append(em.EnemyManager[0].Enemies, em.EnemyManager[0].InitEnemy())
			fmt.Println("Enemy Basic Found at: ", platform.X, platform.Y)
		}
	}
}
