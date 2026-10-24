package core

import "testing"

// newGuardEnemy is a test enemy facing right whose species can guard (its hurt
// art stands in for a guard clip), with the guard roll fixed to pass or fail.
func newGuardEnemy(t *testing.T, rollPasses bool) *EnemyRuntime {
	t.Helper()
	old := rollGuard
	t.Cleanup(func() { rollGuard = old })
	rollGuard = func() float64 {
		if rollPasses {
			return 0
		}
		return 1
	}

	bp := *blueprintFor(KindSuccubus) // copy; its Clips map is replaced below, not shared
	bp.Clips = map[ClipName]EnemyClip{}
	for name, c := range blueprintFor(KindSuccubus).Clips {
		bp.Clips[name] = c
	}
	bp.Clips[ClipGuard] = bp.Clips[ClipHurt]
	bp.Stats.GuardChance = 0.5

	e := newAnimatedTestEnemy(Position{X: 500})
	e.bp = &bp
	e.Animations = bp.newAnimations()
	return e
}

func TestGuardBlocksFrontHit(t *testing.T) {
	e := newGuardEnemy(t, true)
	front := centerX(e.GetBounds()) + 40 // facing right, hit from the right

	e.TakeHitFrom(10, front)
	if e.Health != e.bp.Stats.MaxHealth || e.State != EnemyGuard {
		t.Fatalf("front hit: health=%d state=%d, want no damage and EnemyGuard", e.Health, e.State)
	}

	e.VelX = 100
	e.decideBehavior(farPlayer(e))
	if e.State != EnemyGuard || e.VelX != 0 {
		t.Errorf("while guarding: state=%d velX=%v, want still guarding and standing", e.State, e.VelX)
	}

	e.PrevState = e.State
	tickEnemyAnim(e, clipTicks(e, EnemyGuard))
	if e.State != EnemyPatrol {
		t.Errorf("after the guard clip: state=%d, want back to EnemyPatrol", e.State)
	}
}

func TestGuardIgnoresBackHit(t *testing.T) {
	e := newGuardEnemy(t, true)
	back := centerX(e.GetBounds()) - 40 // facing right, hit from the left

	e.TakeHitFrom(10, back)
	if e.Health != e.bp.Stats.MaxHealth-10 || e.State != EnemyHurt {
		t.Errorf("back hit: health=%d state=%d, want 10 damage and EnemyHurt", e.Health, e.State)
	}

	e = newGuardEnemy(t, true)
	e.FlipX = true // now facing left, so the same side is its front
	e.TakeHitFrom(10, centerX(e.GetBounds())-40)
	if e.State != EnemyGuard {
		t.Errorf("hit from the left while facing left: state=%d, want EnemyGuard", e.State)
	}

	e = newGuardEnemy(t, true)
	e.bp.FacesLeft = true // art drawn facing left (the mushroom): facing is still FlipX
	e.TakeHitFrom(10, back)
	if e.State != EnemyHurt {
		t.Errorf("left-facing art, facing right, hit from the left: state=%d, want EnemyHurt", e.State)
	}
}

// TestGuardNeedsClipChanceAndRoll checks the hit lands whenever one of the
// conditions for a guard is missing.
func TestGuardNeedsClipChanceAndRoll(t *testing.T) {
	for name, setup := range map[string]func(t *testing.T) *EnemyRuntime{
		"no guard clip": func(t *testing.T) *EnemyRuntime {
			e := newGuardEnemy(t, true)
			delete(e.bp.Clips, ClipGuard)
			return e
		},
		"zero chance": func(t *testing.T) *EnemyRuntime {
			e := newGuardEnemy(t, true)
			e.bp.Stats.GuardChance = 0
			return e
		},
		"roll fails": func(t *testing.T) *EnemyRuntime { return newGuardEnemy(t, false) },
		"mid-swing": func(t *testing.T) *EnemyRuntime {
			e := newGuardEnemy(t, true)
			e.State = EnemyAttack
			return e
		},
	} {
		e := setup(t)
		e.TakeHitFrom(10, centerX(e.GetBounds())+40)
		if e.Health != e.bp.Stats.MaxHealth-10 || e.State != EnemyHurt {
			t.Errorf("%s: health=%d state=%d, want the hit to land", name, e.Health, e.State)
		}
	}
}

// TestPlayerHitsGuardedFromTheFront checks both player hit paths, a swing and
// the U ground wave, report the side they come from: an enemy turned toward
// the player blocks, one with its back to the player is hit.
func TestPlayerHitsGuardedFromTheFront(t *testing.T) {
	for _, facesPlayer := range []bool{true, false} {
		p := comboPlayer(PlayerStateWeakAttack, 0) // facing right, enemy on its right
		e := newGuardEnemy(t, true)
		e.Pos = Position{X: centerX(p.GetBounds()) + 10, Y: p.Pos.Y}
		e.FlipX = facesPlayer // facing left = toward the player
		swingThrough(p, []*EnemyRuntime{e})
		if blocked := e.Health == e.bp.Stats.MaxHealth; blocked != facesPlayer {
			t.Errorf("swing, enemy facing the player=%v: health=%d, want blocked=%v", facesPlayer, e.Health, facesPlayer)
		}

		qt := waveWorld(0, 2000)
		e = newGuardEnemy(t, true)
		e.Pos = Position{X: 700}
		b := e.GetBounds()
		e.Pos.Y += waveGroundY - (b.Y + b.Height) // stand it on the floor
		e.FlipX = facesPlayer
		tickWaves(uPlayer(500, groundWaveSpawnFrame), qt, []*EnemyRuntime{e}, int(tps()))
		if blocked := e.Health == e.bp.Stats.MaxHealth; blocked != facesPlayer {
			t.Errorf("ground wave, enemy facing the player=%v: health=%d, want blocked=%v", facesPlayer, e.Health, facesPlayer)
		}
	}
}
