# syntax=docker/dockerfile:1

# ---- build stage --------------------------------------------------------
FROM golang:1.24 AS build
WORKDIR /src

# Cache module downloads across builds.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compile the game to WebAssembly. Assets are embedded into the binary
# (see assets/embed.go) except the level backgrounds, which the game downloads
# per level (assets/backgrounds_js.go); ship those next to the server.
RUN GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o /out/web/game.wasm ./cmd \
 && cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" /out/web/wasm_exec.js \
 && cp web/index.html /out/web/index.html \
 && mkdir -p /out/levels \
 && cd assets/levels && find . -path './*/background/*' -type f -exec cp --parents {} /out/levels/ \;

# Compile the static file server as a small, static (CGO-free) binary.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server ./server

# ---- run stage ----------------------------------------------------------
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /out/ /app/
ENV WEB_ROOT=/app/web
ENV LEVELS_ROOT=/app/levels
EXPOSE 8080
USER nonroot:nonroot
CMD ["/app/server"]
