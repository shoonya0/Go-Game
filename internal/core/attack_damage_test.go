package core

import "testing"

func TestAttackDamageInBars(t *testing.T) {
	bars := func(hp int) float64 { return float64(hp) / EnemyHealthPerBar }
	total := func(st PlayerStateType) (sum int) {
		for _, h := range playerStrikes[st].Hits {
			sum += h.Damage
		}
		return sum
	}
	for _, c := range []struct {
		name string
		st   PlayerStateType
		want float64
	}{
		{"J", PlayerStateWeakAttack, 0.5},
		{"I", PlayerStateStrongAttack, 0.5},
		{"jump + J", PlayerStateWeakAttackInAir, 0.5},
		{"jump + I", PlayerStateStrongAttackInAir, 1.5},
		{"K", PlayerStateSpecialAttack1, 1.5},
		{"L", PlayerStateSpecialAttack2, 3},
	} {
		if got := bars(total(c.st)); got != c.want {
			t.Errorf("%s deals %v bars, want %v", c.name, got, c.want)
		}
	}
}

// swingThrough plays every frame of the current attack, hitting enemies on each.
func swingThrough(p *PlayerRuntime, enemies []*EnemyRuntime) {
	n := p.Animations[p.State.GetPlayerState()].TotalFrames
	for f := 0; f < n; f++ {
		p.CurrAnimFrame = f
		p.HitEnemies(enemies)
	}
}

func TestComboAttacksHitTwicePerSwing(t *testing.T) {
	for _, c := range []struct {
		name          string
		st            PlayerStateType
		first, second int
	}{
		{"K", PlayerStateSpecialAttack1, SpecialAttack1Damage1, SpecialAttack1Damage2},
		{"L", PlayerStateSpecialAttack2, SpecialAttack2Damage1, SpecialAttack2Damage2},
	} {
		p := comboPlayer(c.st, 0)
		e := newTestEnemy(Position{X: centerX(p.GetBounds()) + 10, Y: p.Pos.Y})
		e.bp = &EnemyBlueprint{BodyWidth: 46, BodyHeight: 100, Stats: EnemyStats{MaxHealth: 200}}
		e.Health = 200

		hit := playerStrikes[c.st].Hits
		for f := hit[0].FirstFrame; f <= hit[0].LastFrame; f++ {
			p.CurrAnimFrame = f
			p.HitEnemies([]*EnemyRuntime{e})
		}
		if e.Health != 200-c.first {
			t.Errorf("%s first hit: health=%d, want %d", c.name, e.Health, 200-c.first)
		}

		swingThrough(p, []*EnemyRuntime{e})
		if e.Health != 200-c.first-c.second {
			t.Errorf("%s whole swing: health=%d, want %d (each hit once)", c.name, e.Health, 200-c.first-c.second)
		}
	}
}

func TestLIsSlowerThanK(t *testing.T) {
	a := InitPlayerAnimations()
	if a[SpecialAttack2].AnimationSpeed >= a[SpecialAttack1].AnimationSpeed {
		t.Errorf("L plays at %v fps, K at %v: want L slower", a[SpecialAttack2].AnimationSpeed, a[SpecialAttack1].AnimationSpeed)
	}
}
