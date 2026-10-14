---
name: go-ebiten-wasm-deploy
description: Ship a Go/Ebitengine game (or any Go binary) to the web as WebAssembly and deploy it — embed assets with go:embed, build with GOOS=js GOARCH=wasm, serve with a small Go static server that sets the application/wasm MIME type, containerize with a multi-stage distroless Dockerfile, and deploy to Render via a blueprint plus GitHub Actions CI. Use when compiling Go to WASM, fixing "asset not found" or a blank screen in the browser, a wrong or missing wasm MIME type, Dockerizing a Go app, or wiring up Render/CI deployment.
---

# Deploy a Go/Ebitengine game to the web (WASM)

The same Go source runs natively and in the browser via WebAssembly. The two things that
break a first attempt are **filesystem access** and the **wasm MIME type**. This skill
handles both and takes it to a running deploy.

## 1. Embed assets (no filesystem in the browser)

`ebitenutil.NewImageFromFile("../assets/x.png")` reads the OS filesystem, which does not
exist in WASM. Embed assets into the binary instead — this also makes native builds
self-contained.

```go
// assets/embed.go
package assets

import "embed"

//go:embed *.png
var FS embed.FS
```

```go
// load from the embedded FS by base name, works on native AND wasm
func LoadImage(p string) *ebiten.Image {
    data, _ := assets.FS.ReadFile(path.Base(p))
    img, _, _ := image.Decode(bytes.NewReader(data))
    return ebiten.NewImageFromImage(img)
}
```

`go:embed` only reaches files at or below the `.go` file's own directory, so put `embed.go`
in the same folder as the assets.

## 2. Build to WebAssembly

```bash
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o web/game.wasm ./cmd
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" web/wasm_exec.js   # Go 1.24+ path
```

(Older Go: `misc/wasm/wasm_exec.js`.) Add an `index.html` that streams and runs it:

```html
<script src="wasm_exec.js"></script>
<script>
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("game.wasm"), go.importObject)
    .then(r => go.run(r.instance));
</script>
```

## 3. Serve with the correct MIME type

`WebAssembly.instantiateStreaming` **requires** `Content-Type: application/wasm`, or the
browser refuses to compile it. A tiny Go static server that listens on `$PORT`:

```go
mime.AddExtensionType(".wasm", "application/wasm")
port := os.Getenv("PORT"); if port == "" { port = "8080" }
http.Handle("/", http.FileServer(http.Dir("./web")))
log.Fatal(http.ListenAndServe(":"+port, nil))
```

## 4. Containerize (multi-stage, distroless)

```dockerfile
FROM golang:1.24 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o /out/web/game.wasm ./cmd \
 && cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" /out/web/wasm_exec.js \
 && cp web/index.html /out/web/index.html
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server ./server

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /out/ /app/
ENV WEB_ROOT=/app/web
USER nonroot:nonroot
CMD ["/app/server"]
```

The server binary is `CGO_ENABLED=0` so it runs on the static distroless base.

## 5. Deploy to Render + CI

```yaml
# render.yaml — New > Blueprint in Render points at the repo
services:
  - type: web
    name: my-game
    runtime: docker
    dockerfilePath: ./Dockerfile
    plan: free
    healthCheckPath: /
    autoDeploy: true
```

CI that builds both targets on every push (`.github/workflows/ci.yml`); the **native**
Ebitengine build needs GL/X11/ALSA headers, the wasm build does not:

```yaml
- run: sudo apt-get update && sudo apt-get install -y libgl1-mesa-dev xorg-dev libasound2-dev
- run: go vet ./...
- run: go build ./...
- run: GOOS=js GOARCH=wasm go build -o /tmp/game.wasm ./cmd
```

## Checklist

- [ ] All runtime file reads go through `go:embed`, not the OS filesystem.
- [ ] `wasm_exec.js` matches the Go toolchain version used to build.
- [ ] Server sends `application/wasm` for `.wasm`.
- [ ] `index.html` is served `no-cache`; the content-hashed wasm may cache.
- [ ] Server binary is `CGO_ENABLED=0` for a static/distroless image.
- [ ] Container listens on `$PORT` and answers the health check path.

## Pitfalls

- Blank canvas + console error about MIME → the server is not sending `application/wasm`.
- "asset not found" in the browser → still reading the filesystem; embed it.
- Huge `.wasm` (tens of MB) is normal (Go runtime + engine + assets). Trim with
  `-ldflags="-s -w"`, enable gzip/Brotli at the edge; TinyGo shrinks further but drops some
  Ebitengine features.

## References

- [Ebitengine — WebAssembly](https://ebitengine.org/en/documents/webassembly.html)
- [Auto-deploying Ebitengine games](https://apocalypsetheory.com/posts/gamedev/ebitengine-just-publish-itchio/)
