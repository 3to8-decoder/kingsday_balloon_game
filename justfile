# Kingsday Game — build recipes
# https://github.com/casey/just

# Default: list available recipes
default:
    @just --list

# Run the game natively
run:
    go run .

# Build native binary
build:
    go build -o kingsday_game .

# Build WASM for web (stripped, small)
wasm:
    GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o game.wasm .

# Serve WASM build locally for testing
serve: wasm
    python3 -m http.server 8080

# Clean build artifacts
clean:
    rm -f kingsday_game game.wasm

# Run go vet
vet:
    go vet ./...

# Check everything builds (native + wasm)
check: vet build wasm
    @echo "✅ All builds pass"
