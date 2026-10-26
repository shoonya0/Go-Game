package core

import "testing"

// TestSaveRoundTrip: a snapshot encodes and decodes to the same progress, and
// is its own copy of the active portals.
func TestSaveRoundTrip(t *testing.T) {
	WorldInit()
	w, err := NewWorld(StartLevel)
	if err != nil {
		t.Fatal(err)
	}
	p := InitPlayer(nil)
	p.StartAt(w.Level)
	p.activatePortal(w.Level.Portals[len(w.Level.Portals)-1])

	s := p.Snapshot()
	p.ActivePortals[0].Cell.X = -1
	b, err := EncodeSave(s)
	if err != nil {
		t.Fatal(err)
	}
	got, err := DecodeSave(b)
	if err != nil {
		t.Fatal(err)
	}
	if got.Checkpoint != s.Checkpoint || len(got.ActivePortals) != 2 || got.ActivePortals[0] != w.Level.Portals[0].Ref() {
		t.Errorf("decoded %+v, want checkpoint %v and both portals, start first", got, s.Checkpoint)
	}
}

// TestSaveDecodeRejectsBadData: garbage, another version and a save without a
// checkpoint all count as no save.
func TestSaveDecodeRejectsBadData(t *testing.T) {
	for name, b := range map[string]string{
		"garbage":       "{not json",
		"old version":   `{"version":0,"checkpoint":{"Level":"deepForest1"}}`,
		"no checkpoint": `{"version":1}`,
	} {
		if _, err := DecodeSave([]byte(b)); err == nil {
			t.Errorf("%s: no error, want one", name)
		}
	}
}

// TestSaveContinue: Continue on a fresh game stands the player in the saved
// checkpoint, in its level, with the saved portals active; a portal that no
// longer exists is dropped.
func TestSaveContinue(t *testing.T) {
	WorldInit()
	second := testSecondLevel(t)
	w, err := NewWorld(StartLevel)
	if err != nil {
		t.Fatal(err)
	}
	start, last := w.Level.Portals[0], w.Level.Portals[len(w.Level.Portals)-1]
	far := last.Ref()
	far.Level = second
	gone := PortalRef{Level: StartLevel, Cell: Position{X: -60}}
	s := SaveData{Version: saveVersion, Checkpoint: far, ActivePortals: []PortalRef{start.Ref(), gone, far}}

	p := InitPlayer(nil)
	if err := w.Continue(p, s); err != nil {
		t.Fatal(err)
	}
	if w.Name != second || p.Pos != last.StandPos() || p.Checkpoint != far || p.RespawnPos != last.StandPos() {
		t.Errorf("continued in %q at %v, checkpoint %v; want %q standing in %v", w.Name, p.Pos, p.Checkpoint, second, far)
	}
	if len(p.ActivePortals) != 2 || p.ActivePortals[0] != start.Ref() || p.ActivePortals[1] != far {
		t.Errorf("active portals %v, want [%v %v] (the missing one dropped)", p.ActivePortals, start.Ref(), far)
	}
}

// TestSaveContinueMissingCheckpoint: a checkpoint in an unknown level, or at a
// cell with no portal, is an error and leaves the player where they were.
func TestSaveContinueMissingCheckpoint(t *testing.T) {
	WorldInit()
	w, err := NewWorld(StartLevel)
	if err != nil {
		t.Fatal(err)
	}
	p := InitPlayer(nil)
	p.StartAt(w.Level)
	before := p.Pos
	for _, cp := range []PortalRef{{Level: "noSuchLevel"}, {Level: StartLevel, Cell: Position{X: -60}}} {
		s := SaveData{Version: saveVersion, Checkpoint: cp, ActivePortals: []PortalRef{cp}}
		if err := w.Continue(p, s); err == nil || p.Pos != before || w.Name != StartLevel {
			t.Errorf("checkpoint %v: err %v, at %v in %q; want an error and nothing moved", cp, err, p.Pos, w.Name)
		}
	}
}

// TestSaveAutosaveOnChangeOnly: the start portal is due once; after it is
// written nothing is due until a portal is saved or the checkpoint moves.
func TestSaveAutosaveOnChangeOnly(t *testing.T) {
	WorldInit()
	w, err := NewWorld(StartLevel)
	if err != nil {
		t.Fatal(err)
	}
	p := InitPlayer(nil)
	p.StartAt(w.Level)
	start, last := w.Level.Portals[0], w.Level.Portals[len(w.Level.Portals)-1]

	var a Autosave
	s, due := a.Due(p)
	if !due || s.Checkpoint != start.Ref() {
		t.Fatalf("new game: due %v with checkpoint %v, want the start portal due", due, s.Checkpoint)
	}
	a.Written(s)
	if _, due := a.Due(p); due {
		t.Error("nothing changed: due, want not")
	}
	p.activatePortal(last)
	if s, due = a.Due(p); !due || s.Checkpoint != last.Ref() {
		t.Errorf("portal saved: due %v with checkpoint %v, want %v due", due, s.Checkpoint, last.Ref())
	}
	a.Written(s)
	p.activatePortal(start) // back to an old portal: only the checkpoint moves
	if _, due := a.Due(p); !due {
		t.Error("checkpoint moved: not due, want due")
	}
}
