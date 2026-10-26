package core

import (
	"slices"
	"testing"
)

// testSecondLevel registers a second level for the test. Only one level file
// exists, so it reuses level 1's grid with another background and size.
func testSecondLevel(t *testing.T) string {
	t.Helper()
	const name = "testLevel2"
	def := Levels[StartLevel]
	def.Background, def.Width = "../assets/levels/aquaticAbyss/background/aquaticAbyss_day.png", 7000
	Levels[name] = def
	t.Cleanup(func() { delete(Levels, name) })
	return name
}

// TestSwitchLevel: switching makes another level current (platforms in a new
// quadtree, portals, enemies, background, size); coming back restores the
// enemies as they were left; an unknown name changes nothing.
func TestSwitchLevel(t *testing.T) {
	WorldInit()
	second := testSecondLevel(t)

	w, err := NewWorld(StartLevel)
	if err != nil {
		t.Fatal(err)
	}
	if w.Name != StartLevel || len(w.Level.Portals) == 0 || len(w.Enemies) != len(w.Level.Enemies) {
		t.Fatalf("new world: level %q, %d portals, %d of %d enemies; want %q with portals and every enemy",
			w.Name, len(w.Level.Portals), len(w.Enemies), len(w.Level.Enemies), StartLevel)
	}
	for _, c := range w.Level.Portals {
		if c.Level != StartLevel {
			t.Fatalf("portal at %v belongs to level %q, want %q", c.Pos, c.Level, StartLevel)
		}
	}
	firstQT := w.Quadtree
	w.Enemies = w.Enemies[1:] // one enemy killed and culled

	if err := w.SwitchLevel(second); err != nil {
		t.Fatal(err)
	}
	if w.Name != second || w.Def.Background != Levels[second].Background || w.Def.Width != 7000 {
		t.Errorf("after switch: level %q, background %q, width %v; want %q's", w.Name, w.Def.Background, w.Def.Width, second)
	}
	if w.Quadtree == firstQT || w.Quadtree.Root.Bounds.Width != 7000 {
		t.Error("after switch: the quadtree was not rebuilt for the new level")
	}
	if got := w.Quadtree.Retrieve(w.Level.Platforms[0].GetBounds()); !containsCollider(got, &w.Level.Platforms[0]) {
		t.Error("after switch: the new level's platforms are not in the quadtree")
	}
	if len(w.Enemies) != len(w.Level.Enemies) {
		t.Errorf("after switch: %d enemies, want a fresh %d", len(w.Enemies), len(w.Level.Enemies))
	}

	if err := w.SwitchLevel(StartLevel); err != nil {
		t.Fatal(err)
	}
	if w.Name != StartLevel || len(w.Enemies) != len(w.Level.Enemies) {
		t.Errorf("back in %q: %d enemies, want all %d (non-bosses revive on entering, D4)", w.Name, len(w.Enemies), len(w.Level.Enemies))
	}

	if err := w.SwitchLevel("noSuchLevel"); err == nil || w.Name != StartLevel {
		t.Errorf("unknown level: err %v, current %q; want an error and %q kept", err, w.Name, StartLevel)
	}
}

func containsCollider(cs []Collider, c Collider) bool {
	for _, x := range cs {
		if x == c {
			return true
		}
	}
	return false
}

// TestTravelMovesPlayer: travelling to a portal in the current level puts the
// player standing in it, still; one in another level switches to that level
// first. A portal that doesn't exist is an error and moves nothing.
func TestTravelMovesPlayer(t *testing.T) {
	WorldInit()
	second := testSecondLevel(t)
	w, err := NewWorld(StartLevel)
	if err != nil {
		t.Fatal(err)
	}
	p := InitPlayer(nil)
	p.StartAt(w.Level)
	p.Physics.VelX, p.Physics.VelY = 300, -200

	there := w.Level.Portals[len(w.Level.Portals)-1]
	if err := w.Travel(p, there.Ref()); err != nil {
		t.Fatal(err)
	}
	if p.Pos != there.StandPos() || p.Physics.VelX != 0 || p.Physics.VelY != 0 || w.Name != StartLevel {
		t.Errorf("same level: at %v vel (%v, %v) in %q, want still at %v in %q",
			p.Pos, p.Physics.VelX, p.Physics.VelY, w.Name, there.StandPos(), StartLevel)
	}

	far := there.Ref()
	far.Level = second
	if err := w.Travel(p, far); err != nil {
		t.Fatal(err)
	}
	if w.Name != second || p.Pos != there.StandPos() {
		t.Errorf("other level: in %q at %v, want %q at %v", w.Name, p.Pos, second, there.StandPos())
	}

	before := p.Pos
	if err := w.Travel(p, PortalRef{Level: second, Cell: Position{X: -60}}); err == nil || p.Pos != before {
		t.Errorf("missing portal: err %v, at %v; want an error and the player kept at %v", err, p.Pos, before)
	}
}

// TestTravelMakesDestinationCheckpoint: the portal travelled to becomes the
// checkpoint (D10): drawn gold, and the player respawns standing in it.
func TestTravelMakesDestinationCheckpoint(t *testing.T) {
	WorldInit()
	w, err := NewWorld(StartLevel)
	if err != nil {
		t.Fatal(err)
	}
	p := InitPlayer(nil)
	p.StartAt(w.Level)
	from, there := w.Level.Portals[0], w.Level.Portals[len(w.Level.Portals)-1]
	p.activatePortal(there)
	p.activatePortal(from) // the hold at the start portal saved it

	if err := w.Travel(p, there.Ref()); err != nil {
		t.Fatal(err)
	}
	if p.Checkpoint != there.Ref() || p.RespawnPos != there.StandPos() {
		t.Errorf("after travel: checkpoint %v, respawn %v; want the destination %v, %v",
			p.Checkpoint, p.RespawnPos, there.Ref(), there.StandPos())
	}
	if p.PortalLook(there) != PortalCheckpoint || p.PortalLook(from) != PortalActive {
		t.Errorf("after travel: destination look %v, start look %v; want checkpoint, active",
			p.PortalLook(there), p.PortalLook(from))
	}
}

// TestTravelRevivesEnemies: after any travel, also inside the same level,
// the level's non-boss enemies are back fresh at their spawns (D11).
func TestTravelRevivesEnemies(t *testing.T) {
	WorldInit()
	w, err := NewWorld(StartLevel)
	if err != nil {
		t.Fatal(err)
	}
	p := InitPlayer(nil)
	p.StartAt(w.Level)
	w.Enemies[1].TakeDamage(1)
	w.Enemies[1].Pos.X += 200
	w.Enemies = w.Enemies[1:] // the first was killed and culled

	if err := w.Travel(p, w.Level.Portals[len(w.Level.Portals)-1].Ref()); err != nil {
		t.Fatal(err)
	}
	if len(w.Enemies) != len(w.Level.Enemies) {
		t.Fatalf("after travel: %d enemies, want all %d revived", len(w.Enemies), len(w.Level.Enemies))
	}
	for i, e := range w.Enemies {
		want := NewEnemy(w.Level.Enemies[i].Kind, w.Level.Enemies[i].Pos)
		if e.bp != want.bp || e.Pos != want.Pos || e.Health != want.Health || e.State != EnemyPatrol {
			t.Errorf("enemy %d: %s at %v, health %d, state %d; want a fresh %s at %v, health %d, patrolling",
				i, e.bp.Name, e.Pos, e.Health, e.State, want.bp.Name, want.Pos, want.Health)
		}
	}
}

// TestRespawnInCheckpointLevel: dying in another level than the checkpoint's
// brings the player back to the checkpoint's level, standing in it, with that
// level's enemies revived.
func TestRespawnInCheckpointLevel(t *testing.T) {
	WorldInit()
	second := testSecondLevel(t)
	w, err := NewWorld(StartLevel)
	if err != nil {
		t.Fatal(err)
	}
	p := InitPlayer(nil)
	p.StartAt(w.Level)
	w.Enemies = w.Enemies[1:] // one killed in the checkpoint's level

	// into the other level without travelling (a travel would move the checkpoint, D10)
	stand := w.Level.Portals[0].StandPos()
	if err := w.SwitchLevel(second); err != nil {
		t.Fatal(err)
	}
	p.placeAt(stand)
	p.respawn() // as the death animation ends
	w.Respawned(p)

	if w.Name != StartLevel || p.Pos != w.Level.Portals[0].StandPos() {
		t.Errorf("respawned in %q at %v, want %q at %v", w.Name, p.Pos, StartLevel, w.Level.Portals[0].StandPos())
	}
	if len(w.Enemies) != len(w.Level.Enemies) {
		t.Errorf("%d enemies after respawn, want all %d revived", len(w.Enemies), len(w.Level.Enemies))
	}
}

// TestRespawnWithoutCheckpointStays: in a level with no portals there is no
// checkpoint; a respawn stays in the level and still revives its enemies.
func TestRespawnWithoutCheckpointStays(t *testing.T) {
	WorldInit()
	w, err := NewWorld(StartLevel)
	if err != nil {
		t.Fatal(err)
	}
	p := InitPlayer(nil) // never started at a portal
	w.Enemies = w.Enemies[1:]
	w.Respawned(p)
	if w.Name != StartLevel || len(w.Enemies) != len(w.Level.Enemies) {
		t.Errorf("in %q with %d of %d enemies, want %q with all revived", w.Name, len(w.Enemies), len(w.Level.Enemies), StartLevel)
	}
}

// TestEnteringLevelRevivesEnemies: coming back into a level brings back its
// non-boss enemies fresh at their spawns, killed ones included (D4).
func TestEnteringLevelRevivesEnemies(t *testing.T) {
	WorldInit()
	second := testSecondLevel(t)
	w, err := NewWorld(StartLevel)
	if err != nil {
		t.Fatal(err)
	}
	if len(w.Enemies) < 2 {
		t.Fatalf("level %q has %d enemies, want at least 2", StartLevel, len(w.Enemies))
	}
	w.Enemies[1].TakeDamage(1)
	w.Enemies[1].Pos.X += 200
	w.Enemies = w.Enemies[1:] // the first was killed and culled

	if err := w.SwitchLevel(second); err != nil {
		t.Fatal(err)
	}
	if err := w.SwitchLevel(StartLevel); err != nil {
		t.Fatal(err)
	}
	if len(w.Enemies) != len(w.Level.Enemies) {
		t.Fatalf("back in %q: %d enemies, want all %d revived", StartLevel, len(w.Enemies), len(w.Level.Enemies))
	}
	for i, e := range w.Enemies {
		want := NewEnemy(w.Level.Enemies[i].Kind, w.Level.Enemies[i].Pos)
		if e.bp != want.bp || e.Pos != want.Pos || e.Health != want.Health || e.State != EnemyPatrol {
			t.Errorf("enemy %d: %s at %v, health %d, state %d; want a fresh %s at %v, health %d, patrolling",
				i, e.bp.Name, e.Pos, e.Health, e.State, want.bp.Name, want.Pos, want.Health)
		}
	}
}

// TestLevelTiles: every level's grid uses only its own tile types, and the
// pairs it keeps apart never touch (levels spec: e.g. larva never near grass
// in boulderBastion).
func TestLevelTiles(t *testing.T) {
	WorldInit()
	for name, def := range Levels {
		level := LoadLevel(def.File, def.Key)
		if len(def.Tiles) == 0 {
			t.Errorf("level %q lists no tiles", name)
		}
		bad := map[TileType]int{}
		for _, p := range level.Platforms {
			if !slices.Contains(def.Tiles, p.TileInfo.TileType) {
				bad[p.TileInfo.TileType]++
			}
		}
		for tt, n := range bad {
			t.Errorf("level %q: %d %s cells, but its tiles are %v", name, n, tt, def.Tiles)
		}
		for _, pair := range def.KeepApart {
			if x, y, ok := tilesTouch(level, pair[0], pair[1]); ok {
				t.Errorf("level %q: %s touches %s at cell (%d, %d)", name, pair[0], pair[1], x, y)
			}
		}
	}
}

// TestTilesTouch: the KeepApart check finds side and diagonal neighbours and
// nothing farther.
func TestTilesTouch(t *testing.T) {
	g, l := slices.Index(tileRowOrder, Grass)+1, slices.Index(tileRowOrder, Larva)+1 // grid codes
	cases := []struct {
		name  string
		grid  [][]int
		touch bool
	}{
		{"side", [][]int{{g, l}}, true},
		{"diagonal", [][]int{{g, 0}, {0, l}}, true},
		{"one cell apart", [][]int{{g, 0, l}}, false},
		{"no larva", [][]int{{g, g}}, false},
	}
	for _, c := range cases {
		if _, _, got := tilesTouch(BuildLevel(c.grid), Grass, Larva); got != c.touch {
			t.Errorf("%s: touch = %v, want %v", c.name, got, c.touch)
		}
	}
}

// tilesTouch finds a cell of type a with a cell of type b among its eight
// neighbours, and returns a's cell.
func tilesTouch(level Level, a, b TileType) (x, y int, ok bool) {
	cells := map[[2]int]TileType{}
	for _, p := range level.Platforms {
		cells[[2]int{int(p.X) / LevelTileWidth, int(p.Y) / LevelTileHeight}] = p.TileInfo.TileType
	}
	for c, tt := range cells {
		if tt != a {
			continue
		}
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if cells[[2]int{c[0] + dx, c[1] + dy}] == b {
					return c[0], c[1], true
				}
			}
		}
	}
	return 0, 0, false
}
