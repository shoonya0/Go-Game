# Deploying the game

The game is written in Go with [Ebitengine](https://ebitengine.org/). All image
assets are embedded into the binary (`assets/embed.go`), so the same code runs as
a native desktop app **and** as WebAssembly in the browser — there is no external
filesystem dependency at runtime.

## Targets

| Command          | Result                                             |
| ---------------- | -------------------------------------------------- |
| `make run`       | Native desktop window (dev loop)                   |
| `make build`     | Native binary → `bin/game`                         |
| `make build-wasm`| WebAssembly bundle → `web/game.wasm` + `wasm_exec.js` |
| `make serve`     | Build wasm, then serve `./web` on `:8080`          |
| `make docker`    | Production container image `go-game`               |

### Windows (PowerShell)

`make` may not be installed. The equivalent commands:

```powershell
# native run
go run ./cmd

# build to WebAssembly
$env:GOOS = "js"; $env:GOARCH = "wasm"
go build -ldflags="-s -w" -o web/game.wasm ./cmd
Remove-Item Env:GOOS, Env:GOARCH
Copy-Item "$(go env GOROOT)/lib/wasm/wasm_exec.js" web/wasm_exec.js

# serve locally
go run ./server
```

Then open http://localhost:8080.

## How the web build works

1. `GOOS=js GOARCH=wasm go build` compiles the game to `game.wasm`.
2. `wasm_exec.js` (shipped with the Go toolchain) is the JS glue that boots the module.
3. `web/index.html` streams and runs the wasm; Ebitengine draws to a `<canvas>`.
4. `server/` is a tiny Go static file server that serves `./web` on `$PORT` and
   sets the `application/wasm` MIME type required for streaming instantiation.

## Deploying to Render

The repo ships a [`render.yaml`](./render.yaml) blueprint and a multi-stage
[`Dockerfile`](./Dockerfile) that builds the wasm bundle and the server, then runs
the server from a distroless image.

1. Push this repo to GitHub.
2. In Render: **New → Blueprint**, point it at the repo. It reads `render.yaml`.
3. Render builds the `Dockerfile` and runs the web service; `autoDeploy` ships every
   push to `main`.

The container listens on Render's injected `$PORT` and passes the `/` health check.

## Continuous integration

[`.github/workflows/ci.yml`](./.github/workflows/ci.yml) runs `go vet`, a native
build, and the WebAssembly build on every push and pull request.

## Notes

- The wasm binary is large (~15–27 MB) because it bundles the Go runtime, Ebitengine,
  and the game's images. `-ldflags="-s -w"` trims it; enabling gzip/Brotli at the
  edge cuts transfer size further. [TinyGo](https://tinygo.org/) can shrink it more
  but does not support every Ebitengine feature.
