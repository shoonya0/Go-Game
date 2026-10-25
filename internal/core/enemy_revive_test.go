package core

import "testing"

// testBossKind registers a boss kind (a mushroom flagged Boss) for the test.
func testBossKind(t *testing.T) EnemyKind {
	t.Helper()
	const kind = EnemyKind(1000)
	bp := *enemyBlueprints[KindMushroom]
	bp.Name, bp.Marker, bp.Boss = "Test Boss", "", true
	enemyBlueprints[kind] = &bp
	t.Cleanup(func() { delete(enemyBlueprints, kind) })
	return kind
}

// TestRespawnRevivesEnemies: after the player respawns, every non-boss enemy
// is back at its spawn, alive, at full health and patrolling, including those
// killed and already culled.
func TestRespawnRevivesEnemies(t *testing.T) {
	spawns := []EnemySpawn{
		{Kind: KindMushroom, Pos: Position{X: 120, Y: 60}},
		{Kind: KindSkeletonWhite, Pos: Position{X: 600, Y: 60}},
	}
	enemies := SpawnEnemies(spawns)
	enemies[1].TakeDamage(1)
	enemies[1].Pos.X += 200
	enemies = enemies[1:] // the mushroom was killed and culled

	got := RespawnEnemies(spawns, enemies)
	if len(got) != len(spawns) {
		t.Fatalf("%d enemies after respawn, want %d", len(got), len(spawns))
	}
	for i, e := range got {
		want := NewEnemy(spawns[i].Kind, spawns[i].Pos)
		if e.bp != want.bp || e.Pos != want.Pos || e.Health != want.Health || e.State != EnemyPatrol {
			t.Errorf("enemy %d: %s at %v, health %d, state %d; want a fresh %s at %v, health %d, patrolling",
				i, e.bp.Name, e.Pos, e.Health, e.State, want.bp.Name, want.Pos, want.Health)
		}
	}
}

// TestKilledBossStaysDead: a killed boss is not revived; a living boss is
// left exactly as it is, hurt and where it stands.
func TestKilledBossStaysDead(t *testing.T) {
	boss := testBossKind(t)
	spawns := []EnemySpawn{
		{Kind: boss, Pos: Position{X: 100}},
		{Kind: KindMushroom, Pos: Position{X: 300}},
	}

	enemies := SpawnEnemies(spawns)
	got := RespawnEnemies(spawns, enemies[1:]) // the boss was killed and culled
	if len(got) != 1 || got[0].bp.Boss {
		t.Fatalf("after killing the boss: %d enemies (boss among them: %v), want only the mushroom", len(got), len(got) > 0 && got[0].bp.Boss)
	}

	enemies = SpawnEnemies(spawns)
	living := enemies[0]
	living.TakeDamage(1)
	living.Pos.X += 50
	got = RespawnEnemies(spawns, enemies)
	kept := 0
	for _, e := range got {
		if e.bp.Boss {
			kept++
			if e != living {
				t.Error("the living boss was replaced; want it kept as it is")
			}
		}
	}
	if kept != 1 || len(got) != 2 {
		t.Errorf("living boss: %d bosses among %d enemies, want 1 among 2", kept, len(got))
	}
}
