# Gimbal

A Gyruss-style arcade shooter built with Go, featuring an Entity Component System (ECS) architecture and modern game development practices.

![Go Version](https://img.shields.io/badge/Go-1.25-blue)
![License](https://img.shields.io/badge/License-MIT-green)

## Overview

Gimbal is a space shooter where you control a ship orbiting around the center of the screen, fighting waves of enemies as you progress through levels. Built with clean architecture, comprehensive testing, and modern Go development practices.

## Features

- 🎮 **Arcade-style gameplay** - Gyruss-inspired orbital combat
- 🏗️ **ECS Architecture** - Entity Component System for flexible, maintainable code
- 🎯 **Multiple Game Systems** - Collision detection, enemy AI, weapons, health, movement
- 📊 **Level System** - Configurable levels with waves, formations, and bosses
- 🎨 **Modern UI** - Responsive UI with EbitenUI
- 🧪 **Well Tested** - Comprehensive unit test coverage
- 🚀 **Cross-platform** - Builds for Linux, Windows, and WebAssembly

## Tech Stack

- **Go 1.25** - Programming language
- **Ebiten v2.8.8** - 2D game engine
- **Donburi v1.15.7** - Entity Component System framework
- **EbitenUI v0.6.2** - UI library
- **Zap v1.27.0** - Structured logging
- **Task** - Build automation (replaces Makefile)

## Prerequisites

- Go 1.25 or later
- [Task](https://taskfile.dev/) for build automation (optional, but recommended)
- For Linux builds: X11 development libraries (installed automatically in CI)

## Installation

### Clone the repository

```bash
git clone https://github.com/jonesrussell/gimbal.git
cd gimbal
```

### Install dependencies

```bash
go mod download
```

### Install development tools (optional)

```bash
task install:tools
```

This installs:
- Air (hot reloading)
- wasmserve (WebAssembly server)
- mockgen (mock generation)

## Running

### Development mode

```bash
# Run with debug features enabled
task dev:run

# Or use hot reloading
task dev:hot

# Or use standard Go command
go run -tags dev .
```

### Production mode

```bash
go run .
```

### WebAssembly (local)

```bash
task dev:serve
# Then open http://localhost:4242 in your browser
```

## Building

### Using Task (Recommended)

```bash
# Build for current platform
task builds:current

# Build for Linux
task builds:linux

# Build for Windows
task builds:windows

# Build for WebAssembly
task builds:web

# Build all platforms
task builds:all
```

### Using Go directly

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -tags build -ldflags "-s -w" -o gimbal .

# Windows
GOOS=windows GOARCH=amd64 go build -tags build -ldflags "-s -w" -o gimbal.exe .

# WebAssembly
GOOS=js GOARCH=wasm go build -tags "build,js" -ldflags "-s -w" -o game.wasm .
```

## Controls

- **Arrow Keys / WASD** - Move player around the orbital path
- **Space** - Fire weapons
- **P** - Pause game
- **ESC** - Return to menu (when paused)

## Project Structure

```
gimbal/
├── assets/              # Game assets (sprites, fonts, entities, levels)
├── internal/
│   ├── app/            # Dependency injection container
│   ├── common/         # Shared interfaces and utilities
│   ├── config/         # Configuration management
│   ├── ecs/            # ECS components, systems, and managers
│   │   ├── core/       # Core components and queries
│   │   ├── systems/    # Game systems (collision, enemy, health, movement, weapon)
│   │   └── managers/   # Resource, score, level managers
│   ├── errors/         # Custom error types
│   ├── game/           # Main game loop and initialization
│   ├── input/          # Input handling
│   ├── logger/         # Structured logging
│   ├── math/           # Math utilities (angles, etc.)
│   └── scenes/         # Scene management (intro, menu, gameplay, pause, gameover)
├── .github/workflows/  # CI/CD workflows
├── main.go             # Application entry point
└── Taskfile.yml        # Build automation
```

## Architecture

### Entity Component System (ECS)

The game uses the Donburi ECS framework for a flexible, data-driven architecture:

- **Components** - Data containers (Position, Sprite, Health, etc.)
- **Systems** - Logic processors (Movement, Collision, Enemy AI, etc.)
- **Entities** - Unique identifiers for game objects

### Systems

- **Collision System** - Handles entity collisions with timeout protection
- **Enemy System** - Manages enemy spawning, movement patterns, and waves
- **Health System** - Handles damage, invincibility, and death
- **Movement System** - Processes entity movement and orbital mechanics
- **Weapon System** - Manages player and enemy projectiles
- **Score Manager** - Tracks score, high score, and bonus lives
- **Level Manager** - Handles level progression and configuration

### Configuration

- **JSON-based** - Entity and level configurations loaded from JSON files
- **Environment variables** - Runtime configuration via `.env` file
- **Functional options** - Clean configuration API

## Testing

### Run tests

```bash
# Run all tests
task tests:all

# Run fast tests only
task tests:short

# Run with coverage
go test ./... -cover

# Generate HTML coverage report
task tests:coverage
```

### Test Coverage

Current test coverage includes:
- Math utilities (100%)
- Error handling (93%)
- Score management
- Configuration validation (63%)
- Context utilities
- JSON utilities

## Development

### Code Quality

```bash
# Lint code
task lint

# Auto-fix lint issues
task lint:fix

# Check for dead code
task deadcode:check
```

### Adding New Features

1. **New Component** - Add to `internal/ecs/core/components.go`
2. **New System** - Create in `internal/ecs/systems/<name>/`
3. **New Scene** - Create in `internal/scenes/<name>/` and register in `scenes/registry.go`
4. **New Configuration** - Add to `internal/config/` and load via JSON or env vars

See [CLAUDE.md](CLAUDE.md) for detailed development guidelines.

## Releases

Releases are automatically built via GitHub Actions when you push a tag:

```bash
git tag v1.0.0
git push origin v1.0.0
```

This will:
- Build Linux and Windows binaries
- Create a GitHub release
- Attach binaries to the release

See [`.github/workflows/release.yml`](.github/workflows/release.yml) for details.

## Configuration

Copy `.env.example` to `.env` and customize:

```bash
cp env.example .env
```

Key configuration options:
- `LOG_LEVEL` - Logging verbosity (DEBUG, INFO, WARN, ERROR)
- `DEBUG` - Enable debug mode
- `GAME_WINDOW_WIDTH` / `GAME_WINDOW_HEIGHT` - Window dimensions

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Write tests for new features
- Follow Go coding standards
- Use structured logging
- Document public APIs
- Keep commits atomic and well-described

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by the classic arcade game Gyruss
- Built with the excellent [Ebiten](https://ebitengine.org/) game engine
- Uses [Donburi](https://github.com/yohamta/donburi) for ECS architecture
