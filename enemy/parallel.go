package enemy

import (
	"fmt"
	"os"
	"player/internal/core"
	"runtime"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

const mxEnInManager = 10

// ParallelEnemyManager extends EnemyManager with parallel processing
type ParallelEnemyManager struct {
	EnemyManager  []EnemyManager
	enemyBasicImg *ebiten.Image
	Animations    map[int]Animation

	PartyManager PartyManager
	// parrallel processing
	wg          sync.WaitGroup // genrally used for waiting for all the workers to finish their work
	mutex       sync.Mutex     // generally used for synchronization of the enemy manager
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

	managers := make([]EnemyManager, workerCount)
	var base EnemyManager
	for i := 0; i < workerCount; i++ {
		managers[i] = base.InitEnemyManager(fmt.Sprintf("EM-%d", i))
	}

	// we do not want to copy our mutex and wait group so we make it a pointer
	return ParallelEnemyManager{
		EnemyManager:  managers,
		enemyBasicImg: core.LoadImage(enemySpriteSheetPath),
		Animations:    InitEnemyAnimations(),
		PartyManager:  InitPartyManager(),
		WorkerCount:   workerCount, // this is for direct machine not for docker

	}
}

func (p *ParallelEnemyManager) Shutdown() {
	p.EnemyManager = nil
}

func (em *ParallelEnemyManager) AddEnemyToLevel(level []core.Platform) {
	// find position of enemy on the map
	fmt.Println("Adding enemies to level")
	if em == nil {
		fmt.Println("Enemy Manager is nil")
		return
	}
	if em.EnemyManager == nil {
		em.EnemyManager = make([]EnemyManager, 0)
	}

	// 1. check level for enemy
	for _, platform := range level {
		if isEnemy(platform.TileInfo.TileType) {
			// 2. if found the enemy put it into the enemy manager accordingly.
			// This handles finding the correct manager and respecting the limit (step 3)
			em.enemySpawnPos(platform.X, platform.Y)
		}
	}

	// 4. call an ManageEnemiesWorkers and make gorotine for managing different EnemyManager using waitGroup and mutex syncing
	for i := range em.EnemyManager {
		em.wg.Add(1)
		go em.ManageEnemiesWorkers(i)
	}
}

func isEnemy(t core.TileType) bool {
	switch t {
	case core.EnemyBasic:
		return true
	default:
		return false
	}
}

func (em *ParallelEnemyManager) enemySpawnPos(x, y float64) {
	// 3. each emyManager should oversee atmost 10 enemy as of now as given in mxEnInManager
	// Try to add to existing managers first
	for i := range em.EnemyManager {
		if len(em.EnemyManager[i].Enemies) < mxEnInManager {
			em.EnemyManager[i].Enemies = append(em.EnemyManager[i].Enemies, em.EnemyManager[i].InitEnemy(core.Position{X: x, Y: y}))
			fmt.Println("Added enemy to manager: ", em.EnemyManager[i].ID, " with ID: ", em.EnemyManager[i].Enemies[len(em.EnemyManager[i].Enemies)-1].ID)
			return
		}
	}

	// If all managers are full or none exist, create a new one
	var base EnemyManager
	newMgrID := fmt.Sprintf("EM-%d", len(em.EnemyManager))
	newMgr := base.InitEnemyManager(newMgrID)
	newMgr.Enemies = append(newMgr.Enemies, newMgr.InitEnemy(core.Position{X: x, Y: y}))
	em.EnemyManager = append(em.EnemyManager, newMgr)
}

func (em *ParallelEnemyManager) ManageEnemiesWorkers(workerID int) {
	defer em.wg.Done()

	// 4. ... mutex syncing
	// Use mutex to protect access to shared resources if needed, or to synchronize start.
	// Here we lock briefly to demonstrate syncing as requested.
	em.mutex.Lock()
	if workerID < len(em.EnemyManager) {
		// Example: Logging or initialization requiring sync
		// fmt.Printf("Worker %d managing %d enemies\n", workerID, len(em.EnemyManager[workerID].Enemies))
	}
	em.mutex.Unlock()

	// Actual management logic would go here (e.g., update loop)
	// For now, this task only requested the structure.
}
