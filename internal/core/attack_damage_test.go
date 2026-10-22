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
		{"O", PlayerStateSpecialAttack4, 3},
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
		{"O", PlayerStateSpecialAttack4, SpecialAttack4Damage1, SpecialAttack4Damage2},
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

// O's burst (part 1) is shorter than its sword thrust (part 2), so an enemy
// between the two reaches is only hit by the thrust.
func TestOHitsUseTheirOwnReach(t *testing.T) {
	st := playerStrikes[PlayerStateSpecialAttack4]
	burst, thrust := st.reachOn(st.Hits[0].FirstFrame), st.reachOn(st.Hits[1].FirstFrame)
	if burst == 0 || thrust <= burst {
		t.Fatalf("burst reach %v, thrust reach %v: want the thrust to reach further", burst, thrust)
	}

	p := comboPlayer(PlayerStateSpecialAttack4, 0)
	gap := (burst + thrust) / 2 // enemy's near edge sits between the two reaches
	e := newTestEnemy(Position{X: centerX(p.GetBounds()) + gap, Y: p.Pos.Y})
	e.bp = &EnemyBlueprint{BodyWidth: 46, BodyHeight: 100, Stats: EnemyStats{MaxHealth: 200}}
	e.Health = 200

	swingThrough(p, []*EnemyRuntime{e})
	if e.Health != 200-SpecialAttack4Damage2 {
		t.Errorf("health=%d, want only the thrust (%d) to land", e.Health, SpecialAttack4Damage2)
	}
}

// Jump + I spins a full circle, so its box also covers the side the player
// faces away from; other attacks only hit in front.
func TestAirHeavyAttackHitsBothSides(t *testing.T) {
	for _, flip := range []bool{false, true} {
		for _, c := range []struct {
			st   PlayerStateType
			want bool
		}{
			{PlayerStateStrongAttackInAir, true},
			{PlayerStateWeakAttack, false},
		} {
			p := attackingPlayer(Position{X: 300, Y: 100}, c.st, playerStrikes[c.st].Hits[0].FirstFrame)
			p.FlipX = flip
			cx := centerX(p.GetBounds())
			behind := Position{X: cx + 10, Y: p.Pos.Y} // enemy's near edge 10 px behind the centre
			if !flip {
				behind.X = cx - 10 - 46
			}
			e := newTestEnemy(behind)
			e.bp = &EnemyBlueprint{BodyWidth: 46, BodyHeight: 100, Stats: EnemyStats{MaxHealth: 200}}
			e.Health = 200
			p.HitEnemies([]*EnemyRuntime{e})
			if hit := e.Health < 200; hit != c.want {
				t.Errorf("state %d facing left=%v: enemy behind hit=%v, want %v", c.st, flip, hit, c.want)
			}
		}
	}
}
