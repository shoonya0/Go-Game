package core

import (
	"slices"
	"testing"
)

// testSecondLevel registers a second level for the test: level 1's grid (so
// its portal cells are the same) with another background.
func testSecondLevel(t *testing.T) string {
	t.Helper()
	const name = "testLevel2"
	def := Levels[StartLevel]
	def.Background = "../assets/levels/aquaticAbyss/background/aquaticAbyss_day.png"
	Levels[name] = def
	t.Cleanup(func() { delete(Levels, name) })
	return name
}

// TestSwitchLevel: switching makes another level current (platforms in a new
// quadtree, portals, enemies, background); coming back restores the
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
	if w.Name != second || w.Def.Background != Levels[second].Background {
		t.Errorf("after switch: level %q, background %q; want %q's", w.Name, w.Def.Background, second)
	}
	if w.Quadtree == firstQT || w.Quadtree.Root.Bounds.Width != w.Level.Width {
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

// TestLevelPortals: every level has at least two portals, so a player who
// arrives at its first one can still find another and open the next level
// (levels decision D5).
func TestLevelPortals(t *testing.T) {
	WorldInit()
	for name, def := range Levels {
		if n := len(LoadLevel(def.File, def.Key).Portals); n < 2 {
			t.Errorf("level %q has %d portals, want at least 2", name, n)
		}
	}
}

// TestLastPortalOpensNextLevel (levels decision D5): saving the last unused
// portal of a level takes the player to the next level's first portal, which
// becomes used and the checkpoint; saving it again later, or saving a portal
// while others are still unused, does not; the last level opens nothing.
func TestLastPortalOpensNextLevel(t *testing.T) {
	WorldInit()
	w, err := NewWorld(StartLevel)
	if err != nil {
		t.Fatal(err)
	}
	p := InitPlayer(nil)
	p.StartAt(w.Level) // the first portal is used
	last := w.Level.Portals[len(w.Level.Portals)-1]
	save := func(c Portal) {
		p.Pos = c.StandPos()
		p.updatePortals(w.Level.Portals, &InputState{Interact: true})
	}

	save(w.Level.Portals[0])
	if p.LevelDone() {
		t.Fatal("saving an already used portal, with others unused, finished the level")
	}
	save(last)
	if !p.LevelDone() {
		t.Fatalf("saving the last unused portal of %q didn't finish it", StartLevel)
	}
	opened, err := w.OpenNextLevel(p)
	if err != nil || !opened {
		t.Fatalf("OpenNextLevel = %v, %v; want true, nil", opened, err)
	}
	entry := w.Level.Portals[0]
	if w.Name != DeepForestNight || p.Pos != entry.StandPos() || p.Checkpoint != entry.Ref() || !slices.Contains(p.ActivePortals, entry.Ref()) {
		t.Fatalf("after opening: in %q at %v, checkpoint %v; want %q at its first portal %v, used and saved",
			w.Name, p.Pos, p.Checkpoint, DeepForestNight, entry.StandPos())
	}
	p.updatePortals(w.Level.Portals, &InputState{})
	if p.LevelDone() {
		t.Error("arriving by teleport finished the level")
	}

	if err := w.Travel(p, last.Ref()); err != nil { // back to day, save the last portal again
		t.Fatal(err)
	}
	save(last)
	if p.LevelDone() {
		t.Error("saving the last portal again opened the next level a second time")
	}

	w.Name = LevelOrder[len(LevelOrder)-1]
	if opened, err := w.OpenNextLevel(p); opened || err != nil {
		t.Errorf("the last level opened %v (err %v), want nothing", opened, err)
	}
}

// TestOldSaveOpensNextLevel: a save from before levels D5 can hold every
// portal of a sub-level as used without the next sub-level open. Saving any
// of them again opens it, once (code review of PR #28: the player was stuck).
func TestOldSaveOpensNextLevel(t *testing.T) {
	WorldInit()
	w, err := NewWorld(StartLevel)
	if err != nil {
		t.Fatal(err)
	}
	p := InitPlayer(nil)
	p.StartAt(w.Level)
	for _, c := range w.Level.Portals { // as Continue restores them
		p.activatePortal(c)
	}
	save := func(c Portal) {
		p.Pos = c.StandPos()
		p.updatePortals(w.Level.Portals, &InputState{Interact: true})
	}

	first := w.Level.Portals[0]
	save(first)
	if !p.LevelDone() {
		t.Fatal("every portal used, next level not open: saving one didn't open it")
	}
	if opened, err := w.OpenNextLevel(p); err != nil || !opened {
		t.Fatalf("OpenNextLevel = %v, %v; want true, nil", opened, err)
	}
	if err := w.Travel(p, first.Ref()); err != nil {
		t.Fatal(err)
	}
	save(first)
	if p.LevelDone() {
		t.Error("the next level is open: saving again opened it a second time")
	}
}

// TestOpenNextLevelWithoutPortalKeepsWorld: when the next level has no portal
// to arrive at, OpenNextLevel reports an error and leaves the world and the
// player as they were, not in the next level at the old coordinates.
func TestOpenNextLevelWithoutPortalKeepsWorld(t *testing.T) {
	WorldInit()
	const empty = "testNoPortals"
	Levels[empty] = Levels[StartLevel]
	levelPortals[empty] = []Position{} // its grid read as having no portal
	order := LevelOrder
	LevelOrder = []string{StartLevel, empty}
	t.Cleanup(func() { delete(Levels, empty); delete(levelPortals, empty); LevelOrder = order })

	w, err := NewWorld(StartLevel)
	if err != nil {
		t.Fatal(err)
	}
	p := InitPlayer(nil)
	p.StartAt(w.Level)
	pos, qt := p.Pos, w.Quadtree
	if opened, err := w.OpenNextLevel(p); opened || err == nil {
		t.Fatalf("OpenNextLevel = %v, %v; want false and an error", opened, err)
	}
	if w.Name != StartLevel || w.Quadtree != qt || p.Pos != pos {
		t.Errorf("after the failed open: in %q at %v (quadtree changed: %v); want %q at %v, unchanged",
			w.Name, p.Pos, w.Quadtree != qt, StartLevel, pos)
	}
}

// TestLevelSizeFromGrid: a level is as big as its grid, so the camera and the
// quadtree cover every row and column (deepForest's grids are 36 rows, more
// than the 34 the other levels have).
func TestLevelSizeFromGrid(t *testing.T) {
	WorldInit()
	for name, def := range Levels {
		w, err := NewWorld(name)
		if err != nil {
			t.Fatal(err)
		}
		grid := levelGrid(t, def)
		wantW, wantH := float64(len(grid[0])*LevelTileWidth), float64(len(grid)*LevelTileHeight)
		if b := w.Quadtree.Root.Bounds; b.Width != wantW || b.Height != wantH {
			t.Errorf("level %q: quadtree %vx%v, want its grid's %vx%v", name, b.Width, b.Height, wantW, wantH)
		}
	}
}

// TestLevelBorderKeepsPlayerIn: every level is closed by an invisible border
// just outside its edges, so a player pushing left, right or up from inside
// stays in the level (the user: "cover it with rock, not shown").
func TestLevelBorderKeepsPlayerIn(t *testing.T) {
	WorldInit()
	for name := range Levels {
		w, err := NewWorld(name)
		if err != nil {
			t.Fatal(err)
		}
		width := w.Level.Width
		cases := []struct {
			edge  string
			start Position
			in    InputState
			velY  float64
		}{
			{"left", Position{X: 5, Y: 60}, InputState{Direction: Direction{LeftRight: -1}}, 0},
			{"right", Position{X: width - bodyWidth - 5, Y: 60}, InputState{Direction: Direction{LeftRight: 1}}, 0},
			{"top", Position{X: width / 2, Y: 5}, InputState{}, -JumpForce},
		}
		for _, c := range cases {
			p := InitPlayer(nil)
			p.Pos = c.start
			p.Physics.VelY = c.velY
			for range int(tps()) { // one second
				UpdatePlayer(p, &c.in, w.Quadtree, nil)
				if p.Pos.X < 0 || p.Pos.X+bodyWidth > width || p.Pos.Y < 0 {
					t.Fatalf("%s, %s edge: the player left the level, at %v", name, c.edge, p.Pos)
				}
			}
		}
	}
}
