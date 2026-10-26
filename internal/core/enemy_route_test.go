package core

import "testing"

// A flyer flies around a platform between it and the player, and no enemy
// attacks through a tile (#35).

// ledgeWorld is waveWorld with a one-tile-thick ledge 4 tiles over the floor,
// from x=300 to x=900.
const ledgeWorldY = waveGroundY - 4*LevelTileHeight

func ledgeWorld() *DynamicQuadtree {
	return waveWorld(0, 2000, AABB{X: 300, Y: ledgeWorldY, Width: 10 * LevelTileWidth, Height: LevelTileHeight})
}

// TestFlyerGoesAroundLedgeToReachPlayer: a flyer under a ledge, with the player
// standing on it, flies around the ledge and bites them.
func TestFlyerGoesAroundLedgeToReachPlayer(t *testing.T) {
	for _, startX := range []float64{420, 700} { // nearer the left end, nearer the right end
		qt := ledgeWorld()
		e := newTestFlyer(Position{X: startX, Y: waveGroundY - 150})
		p := playerOn(600, ledgeWorldY)

		for i := 0; i < 10*int(tps()) && p.Combat.Health == p.Combat.MaxHealth; i++ {
			UpdateEnemy(e, p, qt)
			e.UpdateEnemyAnimation()
		}
		if p.Combat.Health == p.Combat.MaxHealth {
			t.Errorf("from x=%v: flyer ended at (%v, %v) without hitting the player on the ledge", startX, e.Pos.X, e.Pos.Y)
		}
	}
}

// TestFlyerDoesNotAttackThroughWall: a flyer whose attack box reaches past a
// wall doesn't strike the player behind it.
func TestFlyerDoesNotAttackThroughWall(t *testing.T) {
	bp := testFlyerBlueprint
	bp.Stats.AttackRange = 150 // reaches past a one-tile wall
	wall := AABB{X: 600, Y: waveGroundY - 4*LevelTileHeight, Width: LevelTileWidth, Height: 4 * LevelTileHeight}
	qt := waveWorld(0, 2000, wall)
	e := newTestFlyer(Position{X: wall.X - bp.BodyWidth, Y: waveGroundY - bp.BodyHeight})
	e.bp = &bp
	p := playerOn(wall.X+wall.Width+20, waveGroundY)
	if !e.bp.Attacker.InRange(e, p.GetBounds()) {
		t.Fatal("setup: the player must be inside the flyer's attack box")
	}

	e.decideBehavior(p, qt)
	if e.State == EnemyAttack {
		t.Error("attacking through the wall")
	}
}

// TestDiverDoesNotDiveOntoLedgeAbovePlayer: a Smasher-like diver over a ledge,
// with the player under it, doesn't smash onto the ledge.
func TestDiverDoesNotDiveOntoLedgeAbovePlayer(t *testing.T) {
	qt := ledgeWorld()
	e := newTestDiver(600, 4*LevelTileHeight+150) // 150 px over the ledge
	p := standingPlayer(600)                      // on the floor, under the ledge
	if !e.bp.Attacker.InRange(e, p.GetBounds()) {
		t.Fatal("setup: the player must be in the diver's column")
	}

	e.decideBehavior(p, qt)
	if e.State == EnemyAttack {
		t.Error("diving onto the ledge over the player")
	}
}
