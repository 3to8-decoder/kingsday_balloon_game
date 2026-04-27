# Kingsday Game 🧡

A King's Day balloon pop game built with [Ebitengine](https://ebitengine.org/).

Pop balloons before they float away! Build combos for multiplied points. 60 seconds on the clock.

**▶️ Play now:** [https://3to8-decoder.github.io/kingsday_balloon_game/](https://3to8-decoder.github.io/kingsday_balloon_game/)

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

Any static host works — GitHub Pages, itch.io, Netlify, Vercel, Cloudflare Pages.

### CI/CD (GitHub Pages)

Push to `main` → GitHub Actions builds WASM → deploys to Pages automatically.

Cache busting: the build injects the git short SHA as a query param on `game.wasm` and `wasm_exec.js`, so every deploy busts the browser cache.

### Local test

```bash
just serve
# then open http://localhost:8080
```

## Controls

- **Click / Tap** balloons to pop them
- **Space** to start / continue (desktop)
