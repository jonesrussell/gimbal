# Variables
USERNAME := jonesrussell
PROJECTNAME := gimbal
VERSION := v0.1.0
GO = go
GO_LDFLAGS = -ldflags "-s -w"
BINARY_DIR = bin
BINARY_NAME = gimbal

# Check if WASM_EXEC_PATH is set, if not, set it
ifndef WASM_EXEC_PATH
WASM_EXEC_PATH="$(shell go env GOROOT)/misc/wasm/wasm_exec.js"
endif

# Check if ITCH_USERNAME is set, if not, set it
ifndef ITCH_USERNAME
ITCH_USERNAME=jonesrussell
endif

ITCH_PATH=gimbal

# Help target
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## Serve target: Starts a local server for development
.PHONY: serve
serve:
	@if ! command -v wasmserve >/dev/null 2>&1; then \
		echo "wasmserve could not be found. You can install it with: go install github.com/hajimehoshi/wasmserve@latest" >&2; \
		exit 1; \
	fi
	@echo "Hosting game on http://localhost:4242"
	(cd cmd/gimbal/; wasmserve -http=":4242" -allow-origin='*' -tags .)

# Build Linux target
.PHONY: build/linux
build/linux:
	@echo "Making build Linux build..."
	rm -rf build/linux
	mkdir -p build/linux/$(PROJECTNAME)
	go build -tags build -o build/linux/$(PROJECTNAME)/$(PROJECTNAME) ./cmd/gimbal

# Build Win32 target
.PHONY: build/win32
build/win32:
	@echo "Making build Win32 build..."
	rm -rf build/win32
	mkdir -p build/win32/$(PROJECTNAME)
	GOOS=windows go build -tags build -o build/win32/$(PROJECTNAME)/$(PROJECTNAME).exe ./cmd/gimbal

# Build WebAssembly target
.PHONY: build/web
build/web:
	@echo "Making build wasm build..."
	rm -rf build/web
	mkdir -p build/web
	GOOS=js GOARCH=wasm go build -tags "build,js" -o build/web/game.wasm ./cmd/gimbal
	cp -r html/* build/web
	cp $(WASM_EXEC_PATH) build/web

# Build all target
.PHONY: build
build: build/linux build/win32 build/web

# Upload to itch.io target
.PHONY: itch
itch: clean/build build
	butler push --if-changed build/win32 $(ITCH_USERNAME)/$(ITCH_PATH):windows
	butler push --if-changed build/linux $(ITCH_USERNAME)/$(ITCH_PATH):linux-amd64
	butler push build/web $(ITCH_USERNAME)/$(ITCH_PATH):web
	@echo "Project is live on http://$(ITCH_USERNAME).itch.io/$(ITCH_PATH)"

# Run target
.PHONY: run
run:
	go run ./cmd/gimbal

# Clean build target
.PHONY: clean/build
clean/build:
	rm -rf build

# Clean all target
.PHONY: clean
clean: clean_build
	rm -rf game.wasm site.zip

# Test target
.PHONY: test
test:
	go test -v ./...

# All target
.PHONY: all
all: fmt lint test build

# Format target
.PHONY: fmt
fmt:
	$(GO) fmt ./...

# Lint target
.PHONY: lint
lint:
	golangci-lint run

# Vet target
.PHONY: vet
vet:
	$(GO) vet ./...

## Mod target: Tidies and downloads Go modules
.PHONY: mod
mod:
	go mod tidy
	go clean -modcache
	go mod download
