# Kingsday Game 🧡

A King's Day balloon pop game built with [Ebitengine](https://ebitengine.org/).

Pop balloons before they float away! Build combos for multiplied points. 60 seconds on the clock.

## Quick Start

```bash
go run .
```

## Build

```bash
# Native (macOS/Linux/Windows)
go build -o kingsday_game .

# Windows cross-compile
GOOS=windows go build -o kingsday_game.exe .

# WASM for web
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o game.wasm .
```

Or use [just](https://github.com/casey/just):

```bash
just          # list all recipes
just run      # run native
just build    # build native
just wasm     # build wasm
just serve    # serve wasm locally on :8080
just clean    # clean build artifacts
```

## Web Deploy

Deploying to the web requires 3 files:

| File | Size | Purpose |
|---|---|---|
| `index.html` | ~1KB | Page shell |
| `wasm_exec.js` | ~17KB | Go WASM runtime |
| `game.wasm` | ~12MB raw / ~3MB gzipped | Game code |

Any static host works — GitHub Pages, Netlify, Vercel, Cloudflare Pages.

### Caching

- `game.wasm` should be served with `Cache-Control: public, max-age=31536000, immutable`
- On new releases, change the filename (e.g. `game-<hash>.wasm`) or add a query param (`game.wasm?v=2`) to bust the cache
- Most static hosts auto-gzip, so users only download ~3MB

### Local test

```bash
just serve
# then open http://localhost:8080
```

## Controls

- **Click** balloons to pop them
- **Space** to start / continue
