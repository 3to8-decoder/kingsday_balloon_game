# AGENTS.md

## Project: Kingsday Game

A 2D balloon pop game built with Ebitengine (ebiten/v2).

## Architecture

- `main.go` — Entry point. Game struct, update/draw loop, all game logic.
- `metal_silence_darwin.go` — macOS: suppresses Metal startup warning on stderr.
- `metal_silence_other.go` — Non-macOS: no-op stub for the above.
- `index.html` — Web shell for WASM deployment.
- `wasm_exec.js` — Go WASM runtime (vendored from Go toolchain).

## Build Conventions

- Use `just` (justfile) for all build commands.
- **Native:** `go build -o kingsday_game .`
- **WASM:** `GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o game.wasm .`
  - `-s -w` strips debug info — critical for WASM size (12MB → still 12MB raw, but gzips to ~3MB).
- **game.wasm is gitignored** — always rebuild from source.
- After any code change: `go vet` and `go build` must pass.

## Conventions

- Go standard formatting (`gofmt`).
- Use `ebiten/v2` as the game engine.
- No external assets required — procedural drawing for now.
- Keep the game loop clean: `Update()`, `Draw()`, `Layout()`.
- Prefer `edit` over `write`.
- Test by running the binary.

## Web Deploy

- 3 files: `index.html`, `wasm_exec.js`, `game.wasm`
- Cache `game.wasm` aggressively with immutable headers
- Cache-bust by renaming or versioning the wasm filename
