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
│   │   ├── managers/          # Score, level, resource, entity config managers
│   │   │   ├── resource/      # Resource manager for sprites, audio, and assets
│   │   │   ├── level_config.go    # Level configuration structures
│   │   │   ├── level_loader.go    # Level loading from JSON files
│   │   │   ├── level_definitions.go # Default level definitions
│   │   │   ├── level.go           # LevelManager for level progression
│   │   │   ├── score.go           # ScoreManager for scoring
│   │   │   ├── entity_config.go   # Entity configuration structures
│   │   │   └── entity_loader.go    # Entity config loading from JSON
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
- **Components**: Position, Sprite, Movement, Orbital, Health, Size, Speed, Angle, Scale, EnemyTypeID
- **Tags**: PlayerTag, StarTag, EnemyTag, ProjectileTag, EnemyProjectileTag
- **Systems**:
  - `collision/` - CollisionSystem for entity collision detection
  - `enemy/` - EnemySystem for enemy spawning/movement, EnemyWeaponSystem for enemy shooting (uses sprite assets for projectiles), WaveManager for wave management
  - `health/` - HealthSystem for entity health and invincibility
  - `movement/` - MovementSystem for entity movement patterns
  - `weapon/` - WeaponSystem for player weapon firing

### Configuration
- Game constants are in `config/constants.go`
- Runtime config uses functional options pattern (`WithDebug()`, `WithSpeed()`, etc.)
- Environment variables loaded via `godotenv` and `envconfig`
- Entity configurations loaded from JSON files in `assets/entities/`:
  - `player.json` - Player configuration (health, size, sprite, invincibility)
  - `enemies.json` - Enemy type configurations (health, speed, size, movement patterns)
- Level configurations loaded from JSON files in `assets/levels/` or use defaults

### Sprite Assets
- Sprites are loaded via the ResourceManager from `assets/sprites/`
- Game sprites: `player.png`, `heart.png`, `enemy.png`, `enemy_heavy.png`, `enemy_boss.png`, `enemy_ammo.png`, `enemy_heavy_ammo.png`
- Enemy projectile sprites: `enemy_ammo.png` (Basic enemies), `enemy_heavy_ammo.png` (Heavy enemies, also used for Boss as fallback)
- All sprites use fallback colored placeholders if asset loading fails

### Audio Assets
- Audio files are loaded via the ResourceManager from `assets/sounds/`
- Format: OGG/Vorbis (`.ogg` files) - recommended for background music
- Music tracks:
  - `game_music_main.ogg` - Menu music (plays in main menu)
  - `game_music_level_1.ogg` - Level 1 gameplay music
  - `game_music_boss.ogg` - Boss battle music (plays when boss appears)
- Audio is optional - game continues without audio if initialization fails (e.g., no audio device in containers)
- Audio automatically loops during playback
- Music switches dynamically: level music → boss music when boss spawns, boss music → level music when boss is defeated

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
1. Create scene directory in `internal/scenes/` (e.g., `scenes/myscene/`)
2. Implement `Scene` interface (Update, Draw, Enter, Exit, GetType)
3. Register in `scenes/registry.go` using `RegisterScene()` with a factory function
4. Available scenes: intro, mainmenu, gameplay/playing, pause, gameover, menu

### Adding a New System
1. Create system directory in `internal/ecs/systems/<name>/`
2. Initialize in `game/init_systems.go` `createGameplaySystems()`
3. Call Update in `game/game.go` `updateGameplaySystems()`
4. Systems are organized by domain (collision, enemy, health, movement, weapon)

### Adding Enemy Projectile Sprites
1. Add sprite asset to `assets/sprites/` (e.g., `enemy_ammo.png`)
2. Add sprite configuration to `internal/ecs/managers/resource/sprite_creation.go` in `loadGameSprites()`
3. Update `EnemyWeaponSystem.initializeProjectileSprites()` to load the sprite for the appropriate enemy type
4. The system automatically uses the sprite based on the enemy type when firing projectiles

### Adding Audio/Music
1. Add OGG/Vorbis audio file to `assets/sounds/` (e.g., `game_music_level_2.ogg`)
2. Add audio configuration to `internal/ecs/managers/resource/audio.go` in `LoadAllAudio()`
3. Use `ResourceManager.GetAudioPlayer()` to access the audio player
4. Call `PlayMusic(name, audioResource, volume)` to play music (volume 0.0-1.0)
5. Music automatically loops - use `StopMusic(name)` to stop playback
6. For scene-specific music, add playback in scene's `Enter()` method and stop in `Exit()` method
7. For dynamic music switching (e.g., boss battles), check game state in scene's `Update()` method

## Performance Notes

- Collision detection has a timeout of `config.CollisionTimeout` (half frame budget)
- Systems taking longer than `config.SlowSystemThreshold` (5ms) are logged as warnings
- Use `config.DebugLogInterval` (60 frames) for periodic debug logging
- The `RenderOptimizer` and `ImagePool` are available for rendering optimization

## Level System

- Levels are defined in JSON files in `assets/levels/` or use default definitions
- Each level contains:
  - `LevelNumber` - Level identifier
  - `Metadata` - Name, description, music track, background theme
  - `Waves` - Array of wave configurations with formations, enemy types, spawn delays
  - `Boss` - Boss configuration (enabled, health, movement, shooting)
  - `Difficulty` - Multipliers for enemy speed, health, spawn rate, etc.
  - `CompletionConditions` - Requirements to complete level (boss kill, all waves, etc.)
- `LevelManager` manages level progression and provides current level config
- `WaveManager` (in enemy system) handles wave spawning and completion
- Enemy types: Basic (0), Heavy (1), Boss (2)
- Formation types: Line, Circle, V, Diamond, Diagonal, Spiral, Random
- Movement patterns: Normal, Zigzag, Accelerating, Pulsing

## Audio System

- Audio uses Ebiten's `audio/vorbis` package for OGG/Vorbis playback
- `ResourceManager` loads audio files from embedded assets (`assets/sounds/*`)
- `AudioPlayer` manages background music playback with looping support
- Audio initialization is optional - game runs without audio if device unavailable
- Music tracks:
  - Menu music (`game_music_main`) - Plays in main menu scene
  - Level music (`game_music_level_1`) - Plays during gameplay (level-specific)
  - Boss music (`game_music_boss`) - Plays when boss appears, switches back to level music when defeated
- Music volume defaults to 70% (0.7) for background music
- Audio context uses 44100 Hz sample rate
- Audio files are decoded and cached for efficient playback
