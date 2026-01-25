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

#### Components (`ecs/core/components.go`)
- **Core**: Position, Sprite, Movement, Orbital, Health, Size, Speed, Angle, Scale, EnemyTypeID
- **Gyruss**: EntryPath, BehaviorState, ScaleAnimation, AttackPattern, FirePattern, RetreatTimer
- **Tags**: PlayerTag, StarTag, EnemyTag, ProjectileTag, EnemyProjectileTag

#### Systems (`ecs/systems/`)

**Gyruss System** (`gyruss/`) - Main coordinator for Gyruss-style gameplay:
- Manages stage loading, wave spawning, and boss fights
- Coordinates all subsystems (path, behavior, attack, fire, power-up)
- Implements `EnemySystemInterface` for collision system compatibility

**Subsystems**:
- `path/` - Entry path animations (spiral_in, arc_sweep, straight_in, loop_entry)
- `behavior/` - Enemy state machine (entering → orbiting → attacking → retreating)
- `attack/` - Rush attack patterns (single_rush, paired_rush, loopback_rush, suicide_dive)
- `fire/` - Enemy projectile patterns (single_shot, burst, spray)
- `animation/` - Scale animations with easing
- `powerup/` - Power-up spawning and collection (double_shot, extra_life)

**Other Systems**:
- `collision/` - Entity collision detection with timeout protection
- `health/` - Entity health and invincibility
- `movement/` - Entity movement patterns
- `weapon/` - Player weapon firing
- `enemy/` - Wave manager and enemy spawner (used by GyrussSystem)

### Managers (`ecs/managers/`)
- **GameStateManager** - Core game state (pause, game over, victory, timing)
- **ScoreManager** - Score tracking and bonus lives
- **LevelManager** - Level number tracking
- **ResourceManager** - Sprites, audio, and asset caching
- **StageLoader** - JSON stage configuration loading

### Scene System
Scenes implement the `Scene` interface (Update, Draw, Enter, Exit, GetType) and are registered in `scenes/registry.go`.

### Configuration
- Game constants: `config/constants.go` (game-wide), `ecs/constants.go` (ECS-specific), `game/constants.go` (game loop)
- Runtime config: Functional options (`WithDebug()`, `WithSpeed()`)
- Entity configs: JSON files in `assets/entities/` (player.json)
- Stage configs: JSON files in `assets/stages/` (stage_01.json through stage_06.json)

## Gyruss Stage System

Stages are JSON files in `assets/stages/` with this structure:

```json
{
  "stage_number": 1,
  "planet": "Earth",
  "metadata": { "name": "Stage 1", "description": "..." },
  "waves": [
    {
      "wave_id": "wave_1",
      "spawn_sequence": [
        {
          "enemy_type": "basic",
          "count": 8,
          "entry_path": { "type": "spiral_in", "duration": 2.0 },
          "behavior": { "post_entry": "orbit_then_attack" },
          "attack_pattern": { "type": "single_rush" },
          "fire_pattern": { "type": "single_shot" }
        }
      ],
      "on_clear": "next_wave"
    }
  ],
  "boss": { "enabled": true, "boss_type": "earth_boss", "health": 10 },
  "difficulty": { "enemy_speed_multiplier": 1.0 }
}
```

### Enemy Entry Paths
- `spiral_in` - Spiral from center to orbit position
- `arc_sweep` - Arc sweep entry
- `straight_in` - Direct entry
- `loop_entry` - Looping entry pattern

### Enemy Behaviors
- `orbit_only` - Only orbit, no attacks
- `orbit_then_attack` - Orbit for duration, then attack
- `immediate_attack` - Attack immediately after entry
- `hover_center_then_orbit` - Hover at center, then move to orbit

### Attack Patterns
- `single_rush` - Single enemy rushes player
- `paired_rush` - Two enemies rush together
- `loopback_rush` - Rush with return loop
- `suicide_dive` - One-way dive at player

## Coding Conventions

### Function Design
- Keep functions under 30 lines when possible
- Break down long functions into smaller, focused helpers
- Use descriptive function names that indicate single responsibility

### Receiver Naming
Use consistent receiver names: `gs` for GyrussSystem, `ws` for WeaponSystem, `cs` for CollisionSystem, `evt` for EventSystem, `scoreMgr` for ScoreManager, `sceneMgr` for SceneManager.

### Error Handling
```go
// Use custom GameError with error codes
return errors.NewGameError(errors.AssetNotFound, "player sprite not found")
return errors.NewGameErrorWithCause(errors.SystemInitFailed, "failed to init", err)

// Fluent error building with context
errors.NewErrorBuilder(errors.AssetNotFound, "sprite missing").
    WithCause(err).
    WithContext("sprite_name", name).
    Build()

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
Pass context through the call chain for proper lifecycle management:
- Use `context.Background()` for initialization and startup operations
- Use `context.WithTimeout()` for resource loading operations
- Add cancellation checks in loops and long operations:
```go
select {
case <-ctx.Done():
    return ctx.Err()
default:
}
```

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
3. Call Update in `game/game_systems.go`

### New Stage
1. Create JSON file in `assets/stages/stage_XX.json`
2. Follow existing stage structure (see stage_01.json)
3. Update `StageLoader.GetTotalStages()` if adding beyond stage 6

### New Sprite Asset
1. Add to `assets/sprites/`
2. Configure in `ecs/managers/resource/sprite_creation.go`

### New Audio Track
1. Add OGG file to `assets/sounds/`
2. Configure in `ecs/managers/resource/audio.go`
3. Play via `ResourceManager.GetAudioPlayer().PlayMusic()`

## Performance Notes

- Collision detection timeout: `config.CollisionTimeout` (half frame budget)
- Slow system threshold: `config.SlowSystemThreshold` (5ms)
- `ImagePool` reuses ebiten.Image instances - use `GetImage()` and `ReturnImage()`
- Debug logging interval: `config.DebugLogInterval` (60 frames)
