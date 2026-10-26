//go:build !js

package system

import (
	"os"
	"testing"
)

// useTempSaveDir points the records at a fresh folder for the test.
func useTempSaveDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir() + "/PirateAdventure" // not created yet: StoreRecord makes it
	old := saveDir
	saveDir = func() (string, error) { return dir, nil }
	t.Cleanup(func() { saveDir = old })
	return dir
}

// TestSaveFileRoundTrip: a stored record reads back, a second store replaces
// it, and after a delete there is none; deleting again is no error.
func TestSaveFileRoundTrip(t *testing.T) {
	dir := useTempSaveDir(t)
	if _, ok := LoadRecord("save"); ok {
		t.Fatal("fresh folder: found a record, want none")
	}
	for _, want := range []string{`{"a":1}`, `{"b":2}`} {
		if err := StoreRecord("save", []byte(want)); err != nil {
			t.Fatal(err)
		}
		if got, ok := LoadRecord("save"); !ok || string(got) != want {
			t.Errorf("read %q (ok %v), want %q", got, ok, want)
		}
	}
	if _, err := os.Stat(dir + "/save.json.tmp"); !os.IsNotExist(err) {
		t.Errorf("temporary file left behind (stat err %v)", err)
	}
	for range 2 {
		if err := DeleteRecord("save"); err != nil {
			t.Fatal(err)
		}
	}
	if _, ok := LoadRecord("save"); ok {
		t.Error("after delete: found a record, want none")
	}
}

// TestSaveFileKeysApart: records under different keys don't touch each
// other, so wiping the save keeps the settings.
func TestSaveFileKeysApart(t *testing.T) {
	useTempSaveDir(t)
	if err := StoreRecord("save", []byte("progress")); err != nil {
		t.Fatal(err)
	}
	if err := StoreRecord("settings", []byte("volume")); err != nil {
		t.Fatal(err)
	}
	if err := DeleteRecord("save"); err != nil {
		t.Fatal(err)
	}
	if got, ok := LoadRecord("settings"); !ok || string(got) != "volume" {
		t.Errorf("settings after wiping the save: %q (ok %v), want %q", got, ok, "volume")
	}
}
