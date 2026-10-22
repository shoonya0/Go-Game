package core

import "testing"

func TestEnemiesHaveAtLeastTwoWholeBars(t *testing.T) {
	for kind, bp := range enemyBlueprints {
		hp := bp.Stats.MaxHealth
		if hp < MinEnemyHealthBars*EnemyHealthPerBar || hp%EnemyHealthPerBar != 0 {
			t.Errorf("%s (kind %d): %d HP, want whole bars of %d and at least %d bars", bp.Name, kind, hp, EnemyHealthPerBar, MinEnemyHealthBars)
		}
	}
	if got := newTestEnemy(Position{}).healthBars(); got != 2 {
		t.Errorf("succubus shows %d bars, want 2", got)
	}
}

func TestEnemyHealthRegenerates(t *testing.T) {
	e := newTestEnemy(Position{})
	e.TakeDamage(EnemyHealthPerBar) // one bar down
	for i := 0; i < int(tps()); i++ {
		e.regenerateHealth()
	}
	want := e.bp.Stats.MaxHealth - EnemyHealthPerBar + int(EnemyRegenBarsPerSecond*EnemyHealthPerBar)
	if e.Health != want {
		t.Errorf("after 1s: health=%d, want %d", e.Health, want)
	}

	for i := 0; i < 60*int(tps()); i++ {
		e.regenerateHealth()
	}
	if e.Health != e.bp.Stats.MaxHealth {
		t.Errorf("after a minute: health=%d, want capped at %d", e.Health, e.bp.Stats.MaxHealth)
	}
}

func TestDeadEnemyDoesNotRegenerate(t *testing.T) {
	e := newTestEnemy(Position{X: 100, Y: 100})
	e.TakeDamage(e.Health)
	for i := 0; i < 2*int(tps()); i++ {
		UpdateEnemy(e, staminaPlayer(), nil)
	}
	if e.Health != 0 {
		t.Errorf("dead enemy regenerated to %d", e.Health)
	}
}
