package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// TestBackgroundsOnly: the levels route serves level backgrounds and nothing
// else (no level grids, no folder listings, no path tricks).
func TestBackgroundsOnly(t *testing.T) {
	root := t.TempDir()
	for _, f := range []string{"forest/background/forest_day.png", "forest/level_1.json"} {
		p := filepath.Join(root, filepath.FromSlash(f))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	h := http.StripPrefix("/levels/", backgroundsOnly(http.FileServer(http.Dir(root))))

	cases := map[string]int{
		"/levels/forest/background/forest_day.png":  http.StatusOK,
		"/levels/forest/background/":                http.StatusNotFound,
		"/levels/forest/level_1.json":               http.StatusNotFound,
		"/levels/forest/":                           http.StatusNotFound,
		"/levels/forest/background/../level_1.json": http.StatusNotFound,
	}
	for url, want := range cases {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, url, nil))
		if rec.Code != want {
			t.Errorf("GET %s = %d, want %d", url, rec.Code, want)
		}
	}
}
