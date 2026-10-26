package core

import "testing"

// setSecondAttack makes every new swing pick attack_2 (true) or attack_1.
func setSecondAttack(t *testing.T, second bool) {
	old := rollSecondAttack
	rollSecondAttack = func() bool { return second }
	t.Cleanup(func() { rollSecondAttack = old })
}

// swingUntilHit runs a grounded skeleton facing a player in reach, a tick at a
// time as the game does (AI, then animation), and returns the swing frame on
// which the player first lost health, or -1 if no hit landed within the swing.
func swingUntilHit(t *testing.T, kind EnemyKind) (hitFrame int, e *EnemyRuntime) {
	t.Helper()
	qt := waveWorld(0, 1200)
	bp := blueprintFor(kind)
	e = newTestEnemy(Position{X: 300, Y: waveGroundY - bp.BodyHeight})
	e.bp, e.Scale, e.Health = bp, bp.Scale, bp.Stats.MaxHealth
	e.Animations = bp.newAnimations()
	p := &PlayerRuntime{Pos: Position{X: 330, Y: waveGroundY - bodyHeight}, Combat: Combat{Health: 100, MaxHealth: 100}}

	for range 3 * int(tps()) {
		UpdateEnemy(e, p, qt)
		if p.Combat.Health < 100 {
			return e.swingFrame(), e
		}
		if e.State != EnemyAttack && e.PrevState == EnemyAttack {
			return -1, e // the swing ended without a hit
		}
		e.UpdateEnemyAnimation()
	}
	return -1, e
}

func TestMeleeHitLandsOnHitFrame(t *testing.T) {
	for _, tc := range []struct {
		name   string
		kind   EnemyKind
		second bool
		clip   ClipName
	}{
		{"skeleton attack_1", KindSkeletonWhite, false, ClipAttack1},
		{"skeleton attack_2", KindSkeletonWhite, true, ClipAttack2},
		{"mushroom", KindMushroom, false, ClipAttack1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setSecondAttack(t, tc.second)
			got, e := swingUntilHit(t, tc.kind)
			want := e.bp.Clips[tc.clip].HitFrame
			if want == 0 {
				t.Fatalf("%v has no HitFrame set; the test needs one past frame 0", tc.clip)
			}
			if got != want {
				t.Errorf("hit landed on swing frame %d, want %d (the %v HitFrame)", got, want, tc.clip)
			}
			if e.clipFor(EnemyAttack) != tc.clip {
				t.Errorf("swing played %v, want %v", e.clipFor(EnemyAttack), tc.clip)
			}
		})
	}
}

// A species without attack_2 art always swings with attack_1, whatever the roll.
func TestSecondAttackFallsBackWithoutArt(t *testing.T) {
	setSecondAttack(t, true)
	e := newAnimatedTestEnemy(Position{}) // test enemy: no attack_2
	e.bp.Attacker.Start(e)
	if e.swing.second || e.clipFor(EnemyAttack) != ClipAttack1 {
		t.Errorf("test enemy swing second=%v clip=%v, want attack_1", e.swing.second, e.clipFor(EnemyAttack))
	}
}

// One swing hits at most once, even if the cooldown has run out by a later frame.
func TestMeleeSwingHitsOnce(t *testing.T) {
	setSecondAttack(t, false)
	e := newAnimatedTestEnemy(Position{})
	near := &PlayerRuntime{Pos: Position{X: 30}, Combat: Combat{Health: 100, MaxHealth: 100}}

	e.decideBehavior(near, nil) // swing starts; the test enemy hits on frame 0
	e.PrevState = e.State
	e.attackTimer = 0
	for range 5 {
		e.CurrAnimFrame++
		e.decideBehavior(near, nil)
	}
	if want := 100 - e.bp.Stats.Damage; near.Combat.Health != want {
		t.Errorf("player health = %d after one swing, want %d (one hit)", near.Combat.Health, want)
	}
}

// A swing whose hit frame comes during the cooldown lands as soon as the
// cooldown is ready, rather than missing: the cooldown is longer than an
// attack clip, so missing would drop every other swing.
func TestMeleeSwingLandsAfterCooldown(t *testing.T) {
	setSecondAttack(t, false)
	e := newAnimatedTestEnemy(Position{})
	near := &PlayerRuntime{Pos: Position{X: 30}, Combat: Combat{Health: 100, MaxHealth: 100}}

	e.attackTimer = 1 // still cooling down from the last swing
	e.decideBehavior(near, nil)
	e.PrevState = e.State
	if near.Combat.Health != 100 {
		t.Fatalf("hit landed during the cooldown: health = %d", near.Combat.Health)
	}
	e.attackTimer = 0 // cooldown ready a few frames into the swing
	e.CurrAnimFrame += 2
	e.decideBehavior(near, nil)
	if want := 100 - e.bp.Stats.Damage; near.Combat.Health != want {
		t.Errorf("health = %d once the cooldown was ready, want %d (the blocked swing lands)", near.Combat.Health, want)
	}
}

// Swing after swing, every swing lands exactly once: an enemy whose cooldown
// outlasts its attack clip waits before swinging again instead of whiffing.
func TestMeleeEverySwingLands(t *testing.T) {
	for _, tc := range []struct {
		name   string
		bp     *EnemyBlueprint
		second bool
	}{
		{"skeleton attack_1", blueprintFor(KindSkeletonWhite), false},
		{"skeleton attack_2", blueprintFor(KindSkeletonWhite), true},
		{"mushroom", blueprintFor(KindMushroom), false},
		{"white golem", blueprintFor(KindWhiteGolem), false},
		{"flying demon", blueprintFor(KindFlyingDemon), false},
		{"test enemy, hit on frame 0", &testEnemyBlueprint, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setSecondAttack(t, tc.second)
			qt := waveWorld(0, 1200)
			bp := tc.bp
			e := newTestEnemy(Position{X: 300, Y: waveGroundY - bp.BodyHeight})
			e.bp, e.Scale, e.Health = bp, bp.Scale, bp.Stats.MaxHealth
			e.Animations = bp.newAnimations()
			p := &PlayerRuntime{Pos: Position{X: 330, Y: waveGroundY - bodyHeight}, Combat: Combat{Health: 1000, MaxHealth: 1000}}

			swings, hits := 0, 0
			for range 6 * int(tps()) {
				health := p.Combat.Health
				UpdateEnemy(e, p, qt)
				if p.Combat.Health < health {
					hits++
				}
				if e.State == EnemyAttack && e.PrevState != EnemyAttack {
					if swings > hits {
						t.Fatalf("swing %d ended without a hit (%d hits)", swings, hits)
					}
					swings++
				}
				e.UpdateEnemyAnimation()
			}
			if swings < 4 || hits < swings-1 {
				t.Errorf("%d swings, %d hits in 6s; want every swing to land", swings, hits)
			}
		})
	}
}
