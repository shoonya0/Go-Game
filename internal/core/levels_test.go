package core

import "testing"

// testSecondLevel registers a second level for the test. Only one level file
// exists, so it reuses level 1's grid with another background and size.
func testSecondLevel(t *testing.T) string {
	t.Helper()
	const name = "testLevel2"
	def := Levels[StartLevel]
	def.Background, def.Width = "../assets/levels/deepForest/background/LBackground_1.png", 7000
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
	left := len(w.Enemies)

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
	if w.Name != StartLevel || len(w.Enemies) != left {
		t.Errorf("back in %q: %d enemies, want the %d left there", w.Name, len(w.Enemies), left)
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
