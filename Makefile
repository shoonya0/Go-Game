# Build and run targets for the game (native + WebAssembly) and the web server.
WASM_EXEC := $(shell go env GOROOT)/lib/wasm/wasm_exec.js
LDFLAGS   := -s -w

.PHONY: run build build-wasm serve docker clean

## run the game natively in a desktop window
run:
	go run ./cmd

## build a native desktop binary into ./bin
build:
	go build -ldflags="$(LDFLAGS)" -o bin/game ./cmd

## compile the game to WebAssembly into ./web
build-wasm:
	@mkdir -p web
	GOOS=js GOARCH=wasm go build -ldflags="$(LDFLAGS)" -o web/game.wasm ./cmd
	cp "$(WASM_EXEC)" web/wasm_exec.js
	@echo "built web/game.wasm ($$(du -h web/game.wasm | cut -f1))"

## build the wasm bundle, then serve it locally on http://localhost:8080
serve: build-wasm
	go run ./server

## build the production container image
docker:
	docker build -t go-game .

## remove build artifacts
clean:
	rm -rf bin web/game.wasm web/wasm_exec.js
