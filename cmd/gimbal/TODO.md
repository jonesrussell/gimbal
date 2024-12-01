## Current Project Structure
```
gimbal/
├── cmd/
│   └── gimbal/
│       ├── main.go          # Main application entry point
│       └── TODO.md          # This file
├── internal/
│   ├── config/             # Configuration management
│   │   ├── config.go       # Configuration loading and management
│   │   └── config.development.json # Development configuration
│   ├── engine/            # Game engine components
│   │   ├── game.go        # Game engine implementation
│   │   ├── types.go       # Engine interfaces and types
│   │   ├── constants.go   # Game states and interfaces
│   │   └── stars.go       # Star system implementation
│   ├── entities/          # [TODO] Game entities
│   ├── systems/           # [TODO] Game systems
│   ├── assets/            # [TODO] Game assets
│   └── ui/                # [TODO] UI components
└── go.mod                 # Module definition

[TODO] Additional directories to be created:
├── pkg/                   # Reusable packages
└── web/                  # Web/WASM specific code
```

### Phase 1: Code Modernization & Architecture

1. **Update Dependencies**
- [ ] Upgrade to Go 1.23 for performance improvements and new features
- [ ] Update Ebitengine to v2.6.6 (latest stable)
- [x] Add `go.uber.org/dig` for dependency injection
- [x] Add `go.uber.org/zap` for better logging
- [ ] Add `github.com/stretchr/testify` for testing
- [ ] Add `github.com/vektra/mockery/v2` for mocks

2. **Restructure Project Layout**
- [x] Create basic directory structure
- [x] Move packages to internal/
- [x] Create engine package with basic interfaces
- [x] Implement game engine structure
- [x] Move core/constants.go to engine/
- [x] Move core/stars.go to engine/
- [x] Remove deprecated core directory
- [ ] Fix import paths in main.go
- [ ] Set up DI container
- [ ] Ensure all packages are properly exposed and importable

3. **Configuration Management**
- [x] Create `config` package for centralized configuration
- [x] Move config to internal/
- [x] Implement configuration injection using `dig`
- [x] Move screen-dependent values from constants.go to config
- [x] Add configuration for number of stars
- [ ] Support different config profiles
- [ ] Add configuration validation
- [ ] Implement hot-reloading for development
- [ ] Add environment variable support
- [ ] Create configuration documentation

4. **Code Quality**
- [ ] Fix linting errors:
  - [x] Fix undefined config.Screen in stars.go
  - [ ] Fix import cycle in game package
  - [ ] Fix unused variables in game_test.go
  - [ ] Fix player package redeclarations
  - [ ] Fix player calculation tests
- [ ] Add golangci-lint to CI pipeline
- [ ] Set up pre-commit hooks for linting

5. **Immediate Next Steps**
1. Fix remaining lint errors:
   - [ ] Clean up player package constants
   - [ ] Update player calculation tests
   - [ ] Remove unused test variables
   - [ ] Fix type mismatches in tests

2. Create assets package for resource management:
   - [ ] Implement AssetManager interface
   - [ ] Add basic image loading functionality
   - [ ] Add error handling for missing assets
   - [ ] Integrate with DI container

3. Implement proper star initialization:
   - [ ] Move star image loading to asset manager
   - [ ] Add error handling for failed initialization
   - [ ] Add star configuration validation
   - [ ] Implement star pool for better performance

4. Configuration Fixes:
   - [ ] Implement proper config.New function in config package
   - [ ] Fix config initialization in engine/stars.go
   - [ ] Add proper config injection in game package
   - [ ] Document config initialization pattern

5. Type Definition Fixes:
   - [ ] Fix Game type definition in game package
   - [ ] Ensure proper type exports
   - [ ] Fix circular dependencies if any
   - [ ] Add missing interface implementations

6. Test Cleanup:
   - [ ] Remove unused speed variables from game tests
   - [ ] Properly use player variable in player tests
   - [ ] Add proper test assertions
   - [ ] Implement test helpers for common setup

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
