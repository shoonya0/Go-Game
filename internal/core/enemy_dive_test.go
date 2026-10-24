package core

import (
	"math"
	"slices"
	"testing"
)

// testDiverBlueprint is the test flyer with a dive: it hovers 150 px over the
// player and smashes down. Wind-up 0.5 s, impact clip 0.5 s.
var testDiverBlueprint = func() EnemyBlueprint {
	bp := testFlyerBlueprint
	bp.Name = "Test diver"
	bp.Marker = "~^"
	bp.Mover = flyingMover{LeashTiles: 5, HoverHeight: 150}
	bp.Attacker = diveAttack{FallSpeed: 600, MaxDrop: 300}
	bp.Stats.AttackRange = 40
	bp.Stats.AttackCooldown = 1.5
	bp.Clips = map[ClipName]EnemyClip{}
	for name, c := range testFlyerBlueprint.Clips {
		bp.Clips[name] = c
	}
	bp.Clips[ClipAttack1] = testEnemyClip(3, 0, 6, 12, true) // wind-up: the dive ends it, not the clip
	bp.Clips[ClipAttackLoop] = testEnemyClip(3, 0, 2, 12, true)
	bp.Clips[ClipAttackEnd] = testEnemyClip(3, 2, 4, 8, false)
	return bp
}()

// newTestDiver is a test diver whose center is straight above x, its feet
// height px above the floor.
func newTestDiver(x, height float64) *EnemyRuntime {
	bp := &testDiverBlueprint
	e := newTestEnemy(Position{X: x - bp.BodyWidth/2, Y: waveGroundY - height - bp.BodyHeight})
	e.bp = bp
	e.Animations = bp.newAnimations()
	return e
}

// standingPlayer is a player standing on the wave floor, centered on x.
func standingPlayer(x float64) *PlayerRuntime {
	return &PlayerRuntime{Pos: Position{X: x - bodyWidth/2, Y: waveGroundY - bodyHeight}, Combat: Combat{Health: 1000, MaxHealth: 1000}}
}

// diveTick runs one game tick for e, as the game loop does.
func diveTick(e *EnemyRuntime, p *PlayerRuntime, qt *DynamicQuadtree) {
	UpdateEnemy(e, p, qt)
	e.UpdateEnemyAnimation()
}

func TestDiveInRangeOnlyBelow(t *testing.T) {
	e := newTestDiver(400, 150)
	for _, tc := range []struct {
		name string
		x, y float64 // player center x, top y
		want bool
	}{
		{"straight below", 400, waveGroundY - bodyHeight, true},
		{"below, inside the column", 430, waveGroundY - bodyHeight, true},
		{"below, outside the column", 480, waveGroundY - bodyHeight, false},
		{"level with it", 400, e.Pos.Y, false},
	} {
		target := AABB{X: tc.x - bodyWidth/2, Y: tc.y, Width: bodyWidth, Height: bodyHeight}
		if got := e.bp.Attacker.InRange(e, target); got != tc.want {
			t.Errorf("%s: InRange = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestDiveAttackPhases: over a standing player, the diver winds up in place,
// falls, lands with one hit, then climbs back over the player and dives again
// once its cooldown is over.
func TestDiveAttackPhases(t *testing.T) {
	qt := waveWorld(0, 1200)
	e := newTestDiver(400, 150)
	p := standingPlayer(400)
	startY := e.Pos.Y
	windup := e.bp.Clips[ClipAttack1].duration()

	var phases []ClipName
	hits, impacts := 0, 0
	for i := range 6 * int(tps()) {
		health := p.Combat.Health
		diveTick(e, p, qt)
		if p.Combat.Health < health {
			hits++
		}
		if e.State != EnemyAttack {
			continue
		}
		clip := e.clipFor(EnemyAttack)
		if len(phases) == 0 || phases[len(phases)-1] != clip {
			phases = append(phases, clip)
			if clip == ClipAttackEnd {
				impacts++
				if !e.bp.Mover.Grounded(e) {
					t.Errorf("tick %d: impact while not on the ground (y=%v)", i, e.Pos.Y)
				}
			}
		}
		if impacts == 0 && float64(i) < windup*tps()-1 && e.Pos.Y != startY {
			t.Fatalf("tick %d: moved to y=%v during the wind-up, want to hold at %v", i, e.Pos.Y, startY)
		}
	}

	want := []ClipName{ClipAttack1, ClipAttackLoop, ClipAttackEnd, ClipAttack1, ClipAttackLoop, ClipAttackEnd}
	if len(phases) < len(want) || !slices.Equal(phases[:len(want)], want) {
		t.Errorf("phases %v, want at least %v", phases, want)
	}
	if hits != impacts || hits < 2 {
		t.Errorf("%d hits for %d impacts in 6 s, want one hit per impact and at least 2", hits, impacts)
	}
}

// TestDiveCommitsWhenPlayerLeaves: once the wind-up starts the dive goes
// through, even if the player walks off; it lands without hitting anyone.
func TestDiveCommitsWhenPlayerLeaves(t *testing.T) {
	qt := waveWorld(0, 1200)
	e := newTestDiver(400, 150)
	p := standingPlayer(400)

	diveTick(e, p, qt)
	if e.State != EnemyAttack {
		t.Fatalf("state %d over the player, want EnemyAttack", e.State)
	}
	p.Pos.X += 1000 // gone, out of sight
	landed := false
	for range 2 * int(tps()) {
		diveTick(e, p, qt)
		landed = landed || (e.State == EnemyAttack && e.clipFor(EnemyAttack) == ClipAttackEnd)
	}
	if !landed || p.Combat.Health != p.Combat.MaxHealth {
		t.Errorf("landed=%v player health=%d, want the dive to land and miss", landed, p.Combat.Health)
	}
}

// TestDiveEndsOverPit: with no ground below, the fall gives up and the diver
// flies home instead of falling forever.
func TestDiveEndsOverPit(t *testing.T) {
	qt := waveWorld(0, 0) // no floor
	e := newTestDiver(400, 150)
	p := standingPlayer(400)
	spawnY := e.Pos.Y

	diveTick(e, p, qt)
	p.Pos.X += 1000 // out of sight, so it goes home after the dive
	for range 8 * int(tps()) {
		diveTick(e, p, qt)
	}
	if e.State == EnemyAttack || math.Abs(e.Pos.Y-spawnY) > 1 {
		t.Errorf("after 8 s: state=%d y=%v, want patrolling back at its spawn height %v", e.State, e.Pos.Y, spawnY)
	}
}
