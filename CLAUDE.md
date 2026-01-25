# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Gimbal is a Gyruss-style arcade game written in Go using:
- **Ebiten** (v2.8.8) - 2D game engine
- **Donburi** (v1.15.7) - Entity Component System (ECS) framework
- **EbitenUI** - UI library for game interfaces
- **Zap** - Structured logging

## Build & Run Commands

```bash
# Development (debug mode with hot reload)
task dev:run                    # Run with debug features
task dev:hot                    # Hot reload with Air
task dev:serve                  # WebAssembly at localhost:4242

# Building
task builds:current             # Build for current platform
task builds:linux               # Build for Linux
task builds:windows             # Build for Windows
task builds:web                 # Build for WebAssembly
task builds:all                 # Build all platforms

# Testing
task tests:all                  # All tests with race detection
task tests:short                # Fast tests only
task tests:coverage             # Generate HTML coverage report
go test ./internal/input/...    # Run specific package tests

# Code Quality
task lint                       # Format, vet, and lint
task lint:fix                   # Auto-fix lint issues
task deadcode:check             # Check for dead code

# Dependencies
task install:tools              # Install Air, wasmserve, mockgen
task deps:tidy                  # Tidy and download modules
```

## Architecture

### Dependency Injection
The `app/Container` manages all dependencies with ordered initialization:
1. Logger → 2. Config → 3. Input Handler → 4. Game Instance

### ECS Architecture
- **Components** (`ecs/core/components.go`): Position, Sprite, Movement, Orbital, Health, Size, Speed, Angle, Scale, EnemyTypeID
- **Tags**: PlayerTag, StarTag, EnemyTag, ProjectileTag, EnemyProjectileTag
- **Systems** (`ecs/systems/`):
  - `collision/` - Entity collision detection with timeout protection
  - `enemy/` - Enemy spawning, movement, shooting, wave management
  - `health/` - Entity health and invincibility
  - `movement/` - Entity movement patterns
  - `weapon/` - Player weapon firing

### Managers (`ecs/managers/`)
- **GameStateManager** - Core game state (pause, game over, victory, timing)
- **ScoreManager** - Score tracking and bonus lives
- **LevelManager** - Level progression and configuration
- **ResourceManager** - Sprites, audio, and asset caching
- **WaveManager** - Wave spawning and completion

### Scene System
Scenes implement the `Scene` interface (Update, Draw, Enter, Exit, GetType) and are registered in `scenes/registry.go`.

### Configuration
- Game constants: `config/constants.go`
- Runtime config: Functional options (`WithDebug()`, `WithSpeed()`)
- Entity configs: JSON files in `assets/entities/` (player.json, enemies.json)
- Level configs: JSON files in `assets/levels/`

## Coding Conventions

### Receiver Naming
Use consistent receiver names: `es` for EnemySystem, `ws` for WeaponSystem, `cs` for CollisionSystem, `evt` for EventSystem, `scoreMgr` for ScoreManager, `sceneMgr` for SceneManager.

### Error Handling
```go
// Use custom GameError with error codes
return errors.NewGameError(errors.AssetNotFound, "player sprite not found")
return errors.NewGameErrorWithCause(errors.SystemInitFailed, "failed to init", err)

// Use errors.As() for unwrapping (not type assertions)
var gameErr *GameError
if errors.As(err, &gameErr) {
    // Handle GameError
}
```

### Logging
```go
// Use structured logging with key-value pairs
g.logger.Debug("Player created", "entity_id", entity, "position", pos)
g.logger.Error("System failed", "system", name, "error", err)
```

### Context
Pass context through the call chain for proper lifecycle management.

### Graceful Shutdown
Use `SceneManager.RequestQuit()` instead of `os.Exit()` for proper cleanup.

## Adding New Features

### New ECS Component
1. Define in `ecs/core/components.go`
2. Register: `MyComponent = donburi.NewComponentType[MyData]()`

### New Scene
1. Create in `internal/scenes/<name>/`
2. Implement Scene interface
3. Register in `scenes/registry.go` with factory function

### New System
1. Create in `internal/ecs/systems/<name>/`
2. Initialize in `game/init_systems.go`
3. Call Update in `game/game.go`

### New Sprite Asset
1. Add to `assets/sprites/`
2. Configure in `ecs/managers/resource/sprite_creation.go`

### New Audio Track
1. Add OGG file to `assets/sounds/`
2. Configure in `ecs/managers/resource/audio.go`
3. Play via `ResourceManager.GetAudioPlayer().PlayMusic()`

## Level System

Levels are JSON files in `assets/levels/` containing:
- Waves with formations (Line, Circle, V, Diamond, Diagonal, Spiral, Random)
- Enemy types (Basic=0, Heavy=1, Boss=2)
- Movement patterns (Normal, Zigzag, Accelerating, Pulsing)
- Boss configuration and difficulty multipliers
- Completion conditions

## Performance Notes

- Collision detection timeout: `config.CollisionTimeout` (half frame budget)
- Slow system threshold: `config.SlowSystemThreshold` (5ms)
- `ImagePool` reuses ebiten.Image instances - use `GetImage()` and `ReturnImage()`
- Debug logging interval: `config.DebugLogInterval` (60 frames)
