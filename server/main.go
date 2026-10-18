// Command server serves the compiled WebAssembly game as a static web app.
// It is the container entrypoint used for the Render deployment: it listens on
// $PORT (Render sets this) and serves the ./web bundle produced by the build.
package main

import (
	"log"
	"mime"
	"net/http"
	"os"
	"time"
)

func main() {
	// WebAssembly streaming instantiation requires the correct MIME type.
	if err := mime.AddExtensionType(".wasm", "application/wasm"); err != nil {
		log.Printf("warning: could not register .wasm mime type: %v", err)
	}

	port := envOr("PORT", "8080")
	root := envOr("WEB_ROOT", "./web")

	mux := http.NewServeMux()
	mux.Handle("/", noCacheHTML(http.FileServer(http.Dir(root))))

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           logRequests(mux),
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("serving %q on :%s", root, port)
	log.Fatal(srv.ListenAndServe())
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// noCacheHTML prevents stale pages and binaries after redeploys. The wasm
// binary has a fixed name (not content-hashed), so it must also be revalidated.
func noCacheHTML(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if p == "/" || p == "/index.html" || p == "/game.wasm" {
			w.Header().Set("Cache-Control", "no-cache")
		}
		next.ServeHTTP(w, r)
	})
}
