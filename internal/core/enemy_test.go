package core

import "testing"

// newTestEnemy builds an enemy without loading a sprite sheet, so the logic can be
// tested without a graphics context. It wires only the fields the AI/combat code
// touches.
func newTestEnemy(pos Position) *EnemyRuntime {
	bp := blueprintFor(KindSuccubus)
	return &EnemyRuntime{
		bp:        bp,
		State:     EnemyPatrol,
		PrevState: EnemyPatrol,
		Scale:     bp.Scale,
		Pos:       pos,
		Health:    bp.Stats.MaxHealth,
		spawnX:    pos.X,
		patrolDir: 1,
	}
}

func TestBlueprintRegistered(t *testing.T) {
	bp := blueprintFor(KindSuccubus)
	if bp == nil || bp.Name != "Succubus" {
		t.Fatalf("blueprintFor(KindSuccubus) = %+v, want a Succubus blueprint", bp)
	}
	e := newAnimatedTestEnemy(Position{})
	for _, st := range []EnemyState{EnemyIdle, EnemyPatrol, EnemyAlert, EnemyAttack, EnemyHurt, EnemyDead} {
		if e.animFor(st) == nil {
			t.Errorf("succubus blueprint missing animation for state %d", st)
		}
	}
}

func TestEnemyTakeDamage(t *testing.T) {
	e := newTestEnemy(Position{})

	e.TakeDamage(10)
	if e.State != EnemyHurt {
		t.Errorf("after non-lethal hit state = %d, want EnemyHurt", e.State)
	}
	if e.Health != e.bp.Stats.MaxHealth-10 {
		t.Errorf("health = %d, want %d", e.Health, e.bp.Stats.MaxHealth-10)
	}

	e.TakeDamage(1000) // lethal
	if e.State != EnemyDead || e.Health != 0 {
		t.Errorf("after lethal hit state=%d health=%d, want EnemyDead/0", e.State, e.Health)
	}

	e.TakeDamage(5) // no-op once dead
	if e.State != EnemyDead || e.Health != 0 {
		t.Errorf("damage after death changed state/health: state=%d health=%d", e.State, e.Health)
	}
}

func TestEnemyPatrolTurnsAround(t *testing.T) {
	e := newTestEnemy(Position{X: 0})
	rangePx := e.bp.Stats.PatrolTiles * LevelTileWidth

	e.Pos.X = rangePx + 10 // past the right end of the patrol path
	e.patrol()
	if e.patrolDir != -1 || e.VelX >= 0 {
		t.Errorf("past right end: patrolDir=%v velX=%v, want -1 and negative", e.patrolDir, e.VelX)
	}

	e.Pos.X = -rangePx - 10 // past the left end
	e.patrol()
	if e.patrolDir != 1 || e.VelX <= 0 {
		t.Errorf("past left end: patrolDir=%v velX=%v, want 1 and positive", e.patrolDir, e.VelX)
	}
}

func TestEnemyDetectsChasesAndAttacks(t *testing.T) {
	e := newTestEnemy(Position{X: 0})
	st := e.bp.Stats

	// Player far away: patrol.
	far := &PlayerRuntime{Pos: Position{X: st.DetectRange + 500}, Combat: Combat{Health: 100, MaxHealth: 100}}
	e.decideBehavior(far)
	if e.State != EnemyPatrol {
		t.Errorf("far player: state=%d, want EnemyPatrol", e.State)
	}

	// Player within detection but out of reach: chase toward them.
	mid := &PlayerRuntime{Pos: Position{X: st.DetectRange - 10}, Combat: Combat{Health: 100, MaxHealth: 100}}
	e.decideBehavior(mid)
	if e.State != EnemyAlert || e.VelX <= 0 {
		t.Errorf("mid player: state=%d velX=%v, want EnemyAlert chasing right", e.State, e.VelX)
	}

	// Player in reach: attack and deal one hit.
	near := &PlayerRuntime{Pos: Position{X: 30}, Combat: Combat{Health: 100, MaxHealth: 100}}
	e.attackTimer = 0
	e.decideBehavior(near)
	if e.State != EnemyAttack {
		t.Errorf("near player: state=%d, want EnemyAttack", e.State)
	}
	if near.Combat.Health != 100-st.Damage {
		t.Errorf("player health = %d, want %d after one hit", near.Combat.Health, 100-st.Damage)
	}

	// A second immediate tick must not deal damage (cooldown active).
	e.decideBehavior(near)
	if near.Combat.Health != 100-st.Damage {
		t.Errorf("player took damage during cooldown: health = %d", near.Combat.Health)
	}
}

// TestBlueprintsValid checks every registered species is complete, so a new
// enemy fails here instead of panicking or vanishing mid-game. Markers, sheet
// bounds and health bars have their own tests.
func TestBlueprintsValid(t *testing.T) {
	names := map[string]EnemyKind{}
	for kind, bp := range enemyBlueprints {
		if bp.Name == "" {
			t.Errorf("kind %d: empty Name", kind)
		}
		if other, dup := names[bp.Name]; dup {
			t.Errorf("name %q used by kinds %d and %d", bp.Name, other, kind)
		}
		names[bp.Name] = kind

		for _, name := range requiredClips {
			if !bp.Has(name) {
				t.Errorf("%s: no %v clip, which every enemy needs", bp.Name, name)
			}
		}
		for st := range stateClips {
			if st == EnemyGuard {
				continue // optional trait, checked below
			}
			if _, ok := bp.resolveClip(stateClips[st]); !ok {
				t.Errorf("%s: nothing to play in state %d", bp.Name, st)
			}
		}
		if bp.Mover == nil || bp.Attacker == nil {
			t.Errorf("%s: Mover=%v Attacker=%v, want both set", bp.Name, bp.Mover, bp.Attacker)
		}
		if bp.BodyWidth <= 0 || bp.BodyHeight <= 0 || bp.Scale <= 0 {
			t.Errorf("%s: body %vx%v scale %v, want all positive", bp.Name, bp.BodyWidth, bp.BodyHeight, bp.Scale)
		}
		st := bp.Stats
		for stat, v := range map[string]float64{
			"Damage": float64(st.Damage), "WalkSpeed": st.WalkSpeed, "ChaseSpeed": st.ChaseSpeed,
			"DetectRange": st.DetectRange, "AttackRange": st.AttackRange,
			"AttackCooldown": st.AttackCooldown, "PatrolTiles": st.PatrolTiles,
		} {
			if v <= 0 {
				t.Errorf("%s: %s = %v, want > 0", bp.Name, stat, v)
			}
		}
		if st.GuardChance < 0 || st.GuardChance > 1 || (st.GuardChance > 0 && !bp.Has(ClipGuard)) {
			t.Errorf("%s: GuardChance %v (guard clip: %v), want 0-1 and a guard clip when > 0", bp.Name, st.GuardChance, bp.Has(ClipGuard))
		}
		if st.AttackRange >= st.DetectRange {
			t.Errorf("%s: AttackRange %v >= DetectRange %v, so it would never chase", bp.Name, st.AttackRange, st.DetectRange)
		}
	}
}
