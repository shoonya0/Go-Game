package core

import "testing"

// A ground enemy chases only a player on ground it can walk to (#34).

// playerOn is a player whose feet are at y, centred near x.
func playerOn(x, feetY float64) *PlayerRuntime {
	return &PlayerRuntime{Pos: Position{X: x, Y: feetY - bodyHeight}, Combat: Combat{Health: 100, MaxHealth: 100}}
}

// floorEnemy is the test enemy standing on waveWorld's floor at x.
func floorEnemy(x float64) *EnemyRuntime {
	return newTestEnemy(Position{X: x, Y: waveGroundY - testEnemyBlueprint.BodyHeight})
}

func TestGroundEnemyChaseNeedsWalkableGround(t *testing.T) {
	ledgeY := waveGroundY - 3*LevelTileHeight
	ledge := AABB{X: 600, Y: ledgeY, Width: 5 * LevelTileWidth, Height: LevelTileHeight}

	cases := []struct {
		name  string
		qt    *DynamicQuadtree
		e     *EnemyRuntime
		p     *PlayerRuntime
		chase bool
	}{
		{"same floor", waveWorld(0, 1500), floorEnemy(400), playerOn(700, waveGroundY), true},
		{"player jumping over the same floor", waveWorld(0, 1500), floorEnemy(400), playerOn(700, waveGroundY-2*LevelTileHeight), true},
		{"player on a ledge above", waveWorld(0, 1500, ledge), floorEnemy(400), playerOn(700, ledgeY), false},
		{"player below the enemy's ledge", waveWorld(0, 1500, ledge),
			newTestEnemy(Position{X: 700, Y: ledgeY - testEnemyBlueprint.BodyHeight}), playerOn(400, waveGroundY), false},
		{"gap between them", waveWorld(0, 600, AABB{X: 720, Y: waveGroundY, Width: 800, Height: LevelTileHeight}),
			floorEnemy(400), playerOn(800, waveGroundY), false},
		{"wall between them", waveWorld(0, 1500, AABB{X: 600, Y: waveGroundY - 2*LevelTileHeight, Width: LevelTileWidth, Height: 2 * LevelTileHeight}),
			floorEnemy(400), playerOn(750, waveGroundY), false},
		{"across a pool", poolWorld(), floorEnemy(450), playerOn(840, waveGroundY), true},
		{"player wading in a pool", poolWorld(), floorEnemy(450), playerOn(700, waterDepthFloor), true},
	}
	for _, c := range cases {
		c.e.decideBehavior(c.p, c.qt)
		if chased := c.e.State == EnemyAlert; chased != c.chase {
			t.Errorf("%s: chasing=%v (state %d), want %v", c.name, chased, c.e.State, c.chase)
		}
	}
}

// TestFlyerChasesPlayerOnAnotherPlatform: the ground rule is for walkers only.
func TestFlyerChasesPlayerOnAnotherPlatform(t *testing.T) {
	ledgeY := waveGroundY - 3*LevelTileHeight
	qt := waveWorld(0, 1500, AABB{X: 600, Y: ledgeY, Width: 5 * LevelTileWidth, Height: LevelTileHeight})
	e := newTestFlyer(Position{X: 400, Y: waveGroundY - 300})

	e.decideBehavior(playerOn(700, ledgeY), qt)
	if e.State != EnemyAlert {
		t.Errorf("state %d, want EnemyAlert", e.State)
	}
}
