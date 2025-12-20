# CLAUDE.md

This file provides guidance for Claude Code when working with the Gimbal codebase.

## Project Overview

Gimbal is a Gyruss-style arcade game written in Go using:
- **Ebiten** (v2.8.8) - 2D game engine
- **Donburi** (v1.15.7) - Entity Component System (ECS) framework
- **EbitenUI** - UI library for game interfaces
- **Zap** - Structured logging

## Build & Run

```bash
# Run the game
go run .

# Build
go build -o gimbal .

# Run tests
go test ./...

# Run with hot reload (development)
air
```

## Project Structure

```
gimbal/
├── main.go                    # Application entry point
├── internal/
│   ├── app/                   # Dependency injection container
│   ├── common/                # Shared interfaces and types
│   ├── config/                # Configuration and constants
│   ├── ecs/
│   │   ├── core/              # ECS components, tags, factories
│   │   ├── debug/             # Debug rendering and performance monitoring
│   │   ├── events/            # Event system
│   │   ├── managers/          # Score, level, resource managers
│   │   └── systems/           # ECS systems (collision, enemy, health, movement, weapon)
│   ├── errors/                # Custom error types with codes
│   ├── game/                  # Main game loop and initialization
│   ├── input/                 # Input handling (keyboard, touch, mouse)
│   ├── logger/                # Zap-based logging
│   ├── math/                  # Angle utilities
│   ├── scenes/                # Scene management (intro, menu, gameplay, pause, etc.)
│   └── ui/                    # Responsive UI components
└── assets/                    # Game assets (sprites, fonts, audio)
```

## Key Patterns

### Dependency Injection
The `app/Container` manages all dependencies with ordered initialization:
1. Logger → 2. Config → 3. Input Handler → 4. Game Instance

### ECS Architecture
- **Components**: Position, Sprite, Movement, Orbital, Health, Size, Speed, Angle, Scale
- **Tags**: PlayerTag, StarTag, EnemyTag, ProjectileTag
- **Systems**: HealthSystem, MovementSystem, CollisionSystem, EnemySystem, WeaponSystem

### Configuration
- Game constants are in `config/constants.go`
- Runtime config uses functional options pattern (`WithDebug()`, `WithSpeed()`, etc.)
- Environment variables loaded via `godotenv` and `envconfig`

## Coding Conventions

### Error Handling
Use the custom `errors.GameError` type with error codes:
```go
return errors.NewGameError(errors.AssetNotFound, "player sprite not found")
return errors.NewGameErrorWithCause(errors.SystemInitFailed, "failed to init", err)
```

### Logging
Use structured logging with key-value pairs:
```go
g.logger.Debug("Player created", "entity_id", entity, "position", pos)
g.logger.Error("System failed", "system", name, "error", err)
```

### Interfaces
- `common.Logger` - Logging interface
- `common.GameInputHandler` - Composite input interface
- `common.HealthProvider` - Health system access
- `common.GameUI` - UI interface

### Context Usage
Pass context through the call chain for proper lifecycle management:
```go
func (g *ECSGame) Update() error {
    ctx := g.ctx // Use game's context
    if err := g.updateGameplaySystems(ctx); err != nil {
        return err
    }
}
```

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/input/...
```

Use the `TestableInputHandler` interface for simulating input in tests.

## Common Tasks

### Adding a New ECS Component
1. Define component type in `ecs/core/components.go`
2. Register with Donburi: `MyComponent = donburi.NewComponentType[MyData]()`

### Adding a New Scene
1. Create scene file in `internal/scenes/`
2. Implement `Scene` interface (Update, Draw, Enter, Exit, GetType)
3. Register in `scenes/registry.go`

### Adding a New System
1. Create system in `internal/ecs/systems/<name>/`
2. Initialize in `game/game_init.go` `createGameplaySystems()`
3. Call Update in `game/game.go` `updateGameplaySystems()`

## Performance Notes

- Collision detection has a timeout of `config.CollisionTimeout` (half frame budget)
- Systems taking longer than `config.SlowSystemThreshold` (5ms) are logged as warnings
- Use `config.DebugLogInterval` (60 frames) for periodic debug logging
- The `RenderOptimizer` and `ImagePool` are available for rendering optimization
