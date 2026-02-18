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

	// parallel processing
	wg          sync.WaitGroup
	mutex       sync.Mutex
	WorkerCount int

	// worker pool channels
	workSignal []chan struct{} // per-worker channel to signal "start updating"
	done       chan struct{}   // shared channel workers use to signal "done"
	quit       chan struct{}   // close this to shut down all workers
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

	return ParallelEnemyManager{
		EnemyManager:  managers,
		enemyBasicImg: core.LoadImage(enemySpriteSheetPath),
		Animations:    InitEnemyAnimations(),
		PartyManager:  InitPartyManager(),
		WorkerCount:   workerCount,
	}
}

// Shutdown signals all persistent workers to exit and waits for them to finish.
func (em *ParallelEnemyManager) Shutdown() {
	if em.quit != nil {
		close(em.quit)
		em.wg.Wait()
	}
	em.EnemyManager = nil
}

func (em *ParallelEnemyManager) AddEnemyToLevel(level []core.Platform) {
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
			em.spawnEnemy(platform.X, platform.Y)
		}
	}

	// 3. Start persistent worker goroutines (once, at level load)
	em.startWorkers()
}

// startWorkers spawns one long-lived goroutine per EnemyManager.
// Each worker blocks on its own channel waiting for a signal to update.
func (em *ParallelEnemyManager) startWorkers() {
	workerCount := len(em.EnemyManager)
	em.workSignal = make([]chan struct{}, workerCount)
	em.done = make(chan struct{}, workerCount)
	em.quit = make(chan struct{})

	for i := 0; i < workerCount; i++ {
		em.workSignal[i] = make(chan struct{})
		em.wg.Add(1)
		go em.worker(i)
	}
	fmt.Printf("Started %d persistent enemy workers\n", workerCount)
}

// worker is a long-lived goroutine that waits for a signal each frame.
func (em *ParallelEnemyManager) worker(id int) {
	defer em.wg.Done()

	for {
		select {
		case <-em.quit:
			return
		case <-em.workSignal[id]:
			em.mutex.Lock()
			if id < len(em.EnemyManager) {
				em.EnemyManager[id].Update()
			}
			em.mutex.Unlock()

			em.done <- struct{}{}
		}
	}
}

// Update is called every frame. It signals all workers to process their
// enemies, then waits for all of them to finish before returning.
func (em *ParallelEnemyManager) Update() {
	if em.workSignal == nil {
		return
	}

	// Signal all workers to start
	for i := range em.workSignal {
		em.workSignal[i] <- struct{}{}
	}

	// Wait for all workers to finish this frame
	for range em.workSignal {
		<-em.done
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

func (em *ParallelEnemyManager) spawnEnemy(x, y float64) {
	for i := range em.EnemyManager {
		if len(em.EnemyManager[i].Enemies) < mxEnInManager {
			em.EnemyManager[i].Enemies = append(em.EnemyManager[i].Enemies, em.EnemyManager[i].InitEnemy(core.Position{X: x, Y: y}))
			return
		}
	}

	var base EnemyManager
	newMgrID := fmt.Sprintf("EM-%d", len(em.EnemyManager))
	newMgr := base.InitEnemyManager(newMgrID)
	newMgr.Enemies = append(newMgr.Enemies, newMgr.InitEnemy(core.Position{X: x, Y: y}))
	em.EnemyManager = append(em.EnemyManager, newMgr)
}
