# AGENTS.md

## Project: Kingsday Game

A 2D game built with Ebitengine (ebiten/v2).

## Architecture

- `main.go` — Entry point. Sets up ebiten.Game and runs the game loop.
- Keep it simple. Single-file for now, split when needed.

## Conventions

- Go standard formatting (`gofmt`).
- Use `ebiten/v2` as the game engine.
- No external assets required — use procedural drawing for now.
- Keep the game loop clean: `Update()`, `Draw()`, `Layout()`.

## Rules

- Always `go build` and `go vet` after changes.
- Prefer `edit` over `write`.
- Test by running the binary.
