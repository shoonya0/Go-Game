package core

import "testing"

const waveGroundY = 500.0

// waveWorld builds a quadtree with a grass floor from x0 to x1 at waveGroundY,
// plus any extra solid boxes (walls).
func waveWorld(x0, x1 float64, walls ...AABB) *DynamicQuadtree {
	qt := NewDynamicQuadtree(AABB{X: 0, Y: 0, Width: 4000, Height: 2000})
	grass := Tile{TileType: Grass}
	for x := x0; x < x1; x += LevelTileWidth {
		qt.Insert(&Platform{X: x, Y: waveGroundY, Width: LevelTileWidth, Height: LevelTileHeight, TileInfo: grass})
	}
	for _, w := range walls {
		qt.Insert(&Platform{X: w.X, Y: w.Y, Width: w.Width, Height: w.Height, TileInfo: grass})
	}
	return qt
}

// uPlayer returns a player standing on the floor at x, mid-way through the U
// swing on the given frame.
func uPlayer(x float64, frame int) *PlayerRuntime {
	p := attackingPlayer(Position{X: x, Y: waveGroundY - bodyHeight}, PlayerStateSpecialAttack3, frame)
	p.Animations = InitPlayerAnimations()
	return p
}

// standOnWaveFloor puts e on waveWorld's floor, grounded as a tick of
// UpdateEnemy would leave it.
func standOnWaveFloor(e *EnemyRuntime) {
	e.Pos.Y = waveGroundY - e.bp.BodyHeight
	e.onGround = true
}

// tickWaves runs the wave update for n ticks.
func tickWaves(p *PlayerRuntime, qt *DynamicQuadtree, enemies []*EnemyRuntime, n int) {
	for i := 0; i < n; i++ {
		p.UpdateGroundWaves(qt, enemies)
	}
}

func TestGroundWaveCastOncePerSwing(t *testing.T) {
	qt := waveWorld(0, 2000)
	p := uPlayer(500, groundWaveSpawnFrame-1)

	tickWaves(p, qt, nil, 1)
	if len(p.waves) != 0 {
		t.Fatal("wave cast before the sword reaches the ground")
	}
	p.CurrAnimFrame = groundWaveSpawnFrame
	tickWaves(p, qt, nil, 1)
	p.CurrAnimFrame++
	tickWaves(p, qt, nil, 1)
	if len(p.waves) != 1 {
		t.Fatalf("waves = %d, want 1 per swing", len(p.waves))
	}

	w := p.waves[0]
	u := p.Animations[SpecialAttack3]
	if w.clip.AnimStartFrame != u.AnimStartFrame+u.TotalFrames || w.clip.SpriteSheetYPosition != u.SpriteSheetYPosition {
		t.Errorf("wave clip = %+v, want the frames after U's on row %d", w.clip, u.SpriteSheetYPosition)
	}
	if w.Bottom != waveGroundY {
		t.Errorf("wave runs at y %v, want the ground at %v", w.Bottom, waveGroundY)
	}
}

func TestGroundWaveRunsTheWayThePlayerFaces(t *testing.T) {
	qt := waveWorld(0, 2000)
	for _, left := range []bool{false, true} {
		p := uPlayer(1000, groundWaveSpawnFrame)
		p.FlipX = left
		tickWaves(p, qt, nil, 1)
		start := p.waves[0].X
		tickWaves(p, qt, nil, int(tps()/2))

		moved := p.waves[0].X - start
		if left && moved >= 0 || !left && moved <= 0 {
			t.Errorf("facing left=%v: wave moved %v px", left, moved)
		}
		if left == (p.waves[0].X > centerX(p.GetBounds())) {
			t.Errorf("facing left=%v: wave is on the wrong side of the player", left)
		}
	}
}

func TestGroundWaveFadesAfterItsFrames(t *testing.T) {
	qt := waveWorld(0, 3000)
	p := uPlayer(200, groundWaveSpawnFrame)
	tickWaves(p, qt, nil, 1)

	life := int(tps() * groundWaveFrames / groundWaveFPS)
	tickWaves(p, qt, nil, life-2)
	if len(p.waves) != 1 {
		t.Fatal("wave ended before its frames played")
	}
	tickWaves(p, qt, nil, 3)
	if len(p.waves) != 0 {
		t.Fatal("wave still alive after its last frame")
	}
}

func TestGroundWaveNeedsGround(t *testing.T) {
	const floorEnd = 12 * LevelTileWidth // 720, on a tile edge
	qt := waveWorld(0, floorEnd)

	p := uPlayer(500, groundWaveSpawnFrame)
	tickWaves(p, qt, nil, 1)
	for i := 0; i < int(tps()) && len(p.waves) > 0; i++ {
		if lead := p.waves[0].X + groundWaveWidth/2; lead > floorEnd+1 {
			t.Fatalf("wave front at %v ran past the floor's end %v", lead, floorEnd)
		}
		tickWaves(p, qt, nil, 1)
	}
	if len(p.waves) != 0 {
		t.Fatal("wave kept going past the end of the floor")
	}

	edge := uPlayer(floorEnd-bodyWidth, groundWaveSpawnFrame) // feet at the edge, facing off it
	tickWaves(edge, qt, nil, 1)
	if len(edge.waves) != 0 {
		t.Error("wave cast with no ground in front of the player")
	}
}

func TestGroundWaveStopsAtWall(t *testing.T) {
	wall := AABB{X: 800, Y: waveGroundY - 3*LevelTileHeight, Width: LevelTileWidth, Height: 3 * LevelTileHeight}
	qt := waveWorld(0, 2000, wall)
	p := uPlayer(500, groundWaveSpawnFrame)
	tickWaves(p, qt, nil, 1)

	for i := 0; i < int(tps()*2) && len(p.waves) > 0; i++ {
		if p.waves[0].Box().Intersects(wall) {
			t.Fatal("wave entered the wall")
		}
		tickWaves(p, qt, nil, 1)
	}
	if len(p.waves) != 0 {
		t.Fatal("wave went through the wall")
	}
}

func TestGroundWaveHitsEachEnemyOnce(t *testing.T) {
	qt := waveWorld(0, 2000)
	e := newTestEnemy(Position{X: 700})
	standOnWaveFloor(e)
	enemies := []*EnemyRuntime{e}

	p := uPlayer(500, groundWaveSpawnFrame)
	tickWaves(p, qt, enemies, int(tps())) // the wave runs through the enemy

	if want := e.bp.Stats.MaxHealth - GroundWaveDamage; e.Health != want {
		t.Errorf("enemy health = %d, want %d (one wave hit)", e.Health, want)
	}
}

// TestGroundWaveMissesFlyers: the wave runs along the floor, so it passes
// under an enemy that isn't on the ground, even one low enough to touch it.
func TestGroundWaveMissesFlyers(t *testing.T) {
	qt := waveWorld(0, 2000)
	e := newTestFlyer(Position{X: 700, Y: waveGroundY - testFlyerBlueprint.BodyHeight - 1})
	p := uPlayer(500, groundWaveSpawnFrame)

	touched := false
	for range int(tps()) {
		tickWaves(p, qt, []*EnemyRuntime{e}, 1)
		for _, w := range p.waves {
			touched = touched || w.Box().Intersects(e.GetBounds())
		}
	}
	if !touched {
		t.Fatal("the wave never reached the flyer's body; the test proves nothing")
	}
	if e.Health != e.bp.Stats.MaxHealth {
		t.Errorf("flyer health = %d, want %d (the wave passes under it)", e.Health, e.bp.Stats.MaxHealth)
	}
}

func TestPlayerHitboxesIncludeGroundWaves(t *testing.T) {
	qt := waveWorld(0, 2000)
	p := uPlayer(500, groundWaveSpawnFrame)
	tickWaves(p, qt, nil, 1)

	boxes := p.Hitboxes()
	last := boxes[len(boxes)-1]
	if last.Kind != HitboxAttack || !last.Active || last.Box != p.waves[0].Box() {
		t.Errorf("last hitbox = %+v, want the live wave's box", last)
	}
}
