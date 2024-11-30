### Phase 1: Code Modernization & Architecture

1. **Update Dependencies**
- Upgrade to Go 1.23 for performance improvements and new features
- Update Ebitengine to v2.6.6 (latest stable)
- Add `go.uber.org/dig` for dependency injection
- Add `go.uber.org/zap` for better logging
- Add `github.com/stretchr/testify` for testing
- Add `github.com/vektra/mockery/v2` for mocks

2. **Restructure Project Layout**
```
gimbal/
├── cmd/
│   └── gimbal/
├── internal/
│   ├── engine/      # Game engine components
│   ├── entities/    # Game entities (player, enemies)
│   ├── systems/     # Game systems (collision, scoring)
│   ├── assets/      # Game assets
│   └── ui/          # UI components
├── pkg/             # Reusable packages
└── web/            # Web/WASM specific code
```

3. **Implement Entity Component System (ECS)**
- Use `github.com/bytearena/ecs` for ECS implementation
- Components:
  - Position
  - Sprite
  - Collision
  - Movement
  - Health
  - Score

### Phase 2: Core Features for Alpha

1. **Player Mechanics**
- Smooth circular movement
- Shooting mechanics
- Basic collision detection
- Physics system integration

2. **Enemy System**
- Basic enemy spawning
- Simple movement patterns
- Collision with player/bullets
- Basic AI behavior system

3. **Scoring System**
- Basic point system
- Score display
- High score persistence
- Score multipliers

4. **Game States**
- Title screen
- Game over screen
- Pause functionality
- State persistence

### Phase 3: Testing & CI/CD

1. **Testing Infrastructure**
- Unit tests using testify
- Integration tests
- Performance benchmarks
- Mock generation with mockery
- Test coverage reporting

2. **CI/CD Pipeline**
```yaml
name: Gimbal CI/CD

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - name: Install dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y libgl1-mesa-dev xorg-dev
      - name: Test
        run: go test -v ./... -coverprofile=coverage.txt
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.txt
      - name: Build
        run: make build

  release:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - name: Build releases
        run: |
          make build/linux
          make build/win32
          make build/web
```

### Phase 4: Alpha Release Checklist

1. **Documentation**
- Update README with alpha status
- Add CONTRIBUTING.md
- Add CHANGELOG.md
- Generate and host GoDoc
- Add architecture diagrams

2. **Distribution**
- Create GitHub release
- Set up itch.io page
- Enable GitHub Pages for web version
- Create installation instructions

3. **Monitoring**
- Add Sentry.io for error tracking
- Set up Google Analytics for web version
- Add telemetry for gameplay metrics
- Implement crash reporting

### Implementation Plan

1. First, let's update the main game loop in `core/game.go`:

```go
package engine

type GameState int

const (
    StateTitle GameState = iota
    StatePlaying
    StatePaused
    StateGameOver
)

type Game struct {
    state     GameState
    world     *ecs.World
    systems   []System
    renderer  *Renderer
    input     *InputSystem
    assets    *AssetManager
    logger    *zap.Logger
}

func NewGame(logger *zap.Logger) (*Game, error) {
    g := &Game{
        state:  StateTitle,
        world:  ecs.NewWorld(),
        logger: logger,
    }
    
    // Initialize systems
    g.systems = []System{
        NewMovementSystem(g.world),
        NewCollisionSystem(g.world),
        NewRenderSystem(g.world),
    }
    
    return g, nil
}
```

2. Then update the player implementation to use ECS:

```go
package entities

type Player struct {
    *ecs.BasicEntity
    *components.Position
    *components.Sprite
    *components.Movement
    *components.Health
    *components.Weapon
}

func NewPlayer(world *ecs.World) *Player {
    player := &Player{
        BasicEntity: ecs.NewBasic(),
        Position:    &components.Position{},
        Sprite:      &components.Sprite{},
        Movement:    &components.Movement{},
        Health:      &components.Health{Lives: 3},
        Weapon:      &components.Weapon{},
    }
    
    world.AddEntity(player)
    return player
}
```

3. Next Steps:
- Set up dependency injection with dig
- Implement basic systems (Movement, Collision, Render)
- Add state management
- Create asset loading system
- Implement input handling
