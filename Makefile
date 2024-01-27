USERNAME := jonesrussell
PROJECTNAME := gimbal
VERSION := v1.0.0
GO = go
GO_LDFLAGS = -ldflags "-s -w"
BINARY_DIR = bin
BINARY_NAME = game

ifndef WASM_EXEC_PATH
WASM_EXEC_PATH="$(shell go env GOROOT)/misc/wasm/wasm_exec.js"
endif

ifndef ITCH_USERNAME
ITCH_USERNAME=jonesrussell
endif

ITCH_PATH=gimbal

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## serve: serve WebAssembly version locally
.PHONY: serve
serve:
	@if ! command -v wasmserve >/dev/null 2>&1; then \
		echo "wasmserve could not be found"; \
		echo "You can install it with: go install github.com/hajimehoshi/wasmserve@latest"; \
		exit 1; \
	fi
	@echo "Hosting game on http://localhost:4242"
	(cd cmd/gimbal/; wasmserve -http=":4242" -allow-origin='*' -tags .)
	
## build/linux: build build Linux version
.PHONY: build/linux
build/linux:
	@echo "Making build Linux build..."
	rm -rf build/linux
	mkdir -p build/linux/$(PROJECTNAME)
	go build -tags build -o build/linux/$(PROJECTNAME)/$(PROJECTNAME) ./cmd/gimbal

## build/win32: build build Win32 version
.PHONY: build/win32
build/win32:
	@echo "Making build Win32 build..."
	rm -rf build/win32
	mkdir -p build/win32/$(PROJECTNAME)
	GOOS=windows go build -tags build -o build/win32/$(PROJECTNAME)/$(PROJECTNAME).exe ./cmd/gimbal

## build/web: build build WebAssembly version
.PHONY: build/web
build/web:
	@echo "Making build wasm build..."
	rm -rf build/web
	mkdir -p build/web
	GOOS=js GOARCH=wasm go build -tags "build,js" -o build/web/${PROJECTNAME}.wasm ./cmd/gimbal
	cp -r html/* build/web
	cp $(WASM_EXEC_PATH) build/web

## build: build all build versions
.PHONY: build
build: build/linux build/win32 build/web

## itch: upload all build versions on itch.io
.PHONY: itch
itch: clean/build build
	butler push --if-changed build/win32 $(ITCH_USERNAME)/$(ITCH_PATH):windows
	butler push --if-changed build/linux $(ITCH_USERNAME)/$(ITCH_PATH):linux-amd64
	butler push build/web $(ITCH_USERNAME)/$(ITCH_PATH):web
	@echo "Project is live on http://$(ITCH_USERNAME).itch.io/$(ITCH_PATH)"

## run: run game (dev)
.PHONY: run
run:
	go run ./cmd/gimbal

## clean/build: remove all previosly built build versions
.PHONY: clean/build
clean/build:
	rm -rf build

## clean: clean all build artifacts
.PHONY: clean
clean: clean_build
	rm -rf game.wasm site.zip

## test: run tests
.PHONY: test
test:
	go test -v ./...

all: fmt lint test build

fmt:
	$(GO) fmt ./...

lint:
	golangci-lint run

vet:
	$(GO) vet ./...
