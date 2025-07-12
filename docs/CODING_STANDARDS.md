# Go Game Project Coding Standards

## Overview

This document outlines the coding standards and best practices for the Go game project, focusing on 2025 Go features, clean architecture principles, and maintainable code.

## Table of Contents

1. [Go 2025 Best Practices](#go-2025-best-practices)
2. [Documentation Standards](#documentation-standards)
3. [Architecture Guidelines](#architecture-guidelines)
4. [Error Handling](#error-handling)
5. [Testing Standards](#testing-standards)
6. [Performance Guidelines](#performance-guidelines)
7. [Security Considerations](#security-considerations)

## Go 2025 Best Practices

### 1.1 Generics Usage

Use generics for:
- Container types (Result, Option, Repository)
- Algorithm implementations
- Type-safe collections
- Generic interfaces

```go
// ✅ Good: Generic Result type
type Result[T any] struct {
    value T
    err   error
}

// ✅ Good: Generic Repository interface
type Repository[T any, ID comparable] interface {
    Get(id ID) (T, error)
    Create(entity T) error
    Update(entity T) error
    Delete(id ID) error
}

// ❌ Bad: Non-generic approach
type StringResult struct {
    value string
    err   error
}
```

### 1.2 Structured Logging

Use `log/slog` for structured logging:

```go
// ✅ Good: Structured logging with slog
logger.Info("Player moved",
    "player_id", playerID,
    "old_position", oldPos,
    "new_position", newPos,
    "delta_time", deltaTime)

// ❌ Bad: String concatenation
logger.Info("Player " + playerID + " moved from " + oldPos.String() + " to " + newPos.String())
```

### 1.3 Context Propagation

Always propagate context through function calls:

```go
// ✅ Good: Context-aware functions
func (s *System) Update(ctx context.Context, deltaTime float64) error {
    if ctx.Err() != nil {
        return errors.FromContext(ctx, errors.SystemUpdateFailed, "context cancelled")
    }
    // ... implementation
}

// ❌ Bad: No context
func (s *System) Update(deltaTime float64) error {
    // ... implementation
}
```

### 1.4 Error Handling

Use wrapped errors and custom error types:

```go
// ✅ Good: Wrapped errors with context
if err != nil {
    return errors.Wrapf(err, errors.AssetLoadFailed, "failed to load sprite %s", spriteName)
}

// ✅ Good: Custom error types
err := errors.NewGameError(errors.EntityNotFound, "player entity not found").
    WithContext("entity_id", entityID).
    WithContext("world_id", worldID)

// ❌ Bad: Basic error handling
if err != nil {
    return err
}
```

## Documentation Standards

### 2.1 Package Documentation

Every package must have a package comment:

```go
// Package movement provides ECS movement systems for game entities.
// It handles player movement, enemy AI movement, and projectile trajectories.
//
// The package follows the Entity-Component-System architecture pattern
// and provides both keyboard and touch input support for movement.
package movement
```

### 2.2 Function Documentation

Document all exported functions with complete examples:

```go
// NewMovementSystem creates a new movement system with the provided dependencies.
// The system handles entity movement based on input and physics calculations.
//
// Example:
//
//	config := &config.GameConfig{ScreenSize: config.Size{Width: 800, Height: 600}}
//	logger := logger.NewWithConfig(nil)
//	inputHandler := input.NewHandler()
//	
//	system := NewMovementSystem(world, config, logger, inputHandler)
//	if err := system.Initialize(context.Background()); err != nil {
//	    return err
//	}
//
// Parameters:
//   - world: The ECS world containing entities
//   - config: Game configuration including screen dimensions
//   - logger: Logger for debug and error messages
//   - inputHandler: Input handler for movement commands
//
// Returns:
//   - *MovementSystem: Configured movement system
func NewMovementSystem(world donburi.World, config *config.GameConfig, logger common.Logger, inputHandler common.GameInputHandler) *MovementSystem {
    // ... implementation
}
```

### 2.3 Type Documentation

Document all exported types:

```go
// MovementSystem updates entity positions based on velocity or input.
// It is responsible for moving the player, starfield, and any other moving entities.
//
// The system processes movement in the following order:
// 1. Player input processing
// 2. Physics calculations
// 3. Position updates
// 4. Boundary checking
//
// Thread Safety: This type is not safe for concurrent use.
type MovementSystem struct {
    world        donburi.World
    config       *config.GameConfig
    logger       common.Logger
    inputHandler common.GameInputHandler
}
```

### 2.4 Interface Documentation

Document interfaces with usage examples:

```go
// MovementInputHandler handles movement-specific input for game entities.
// Implementations should provide smooth, responsive movement input
// that can be easily tested and mocked.
//
// Example implementation:
//
//	type KeyboardMovementHandler struct {
//	    keys map[ebiten.Key]bool
//	}
//	
//	func (h *KeyboardMovementHandler) GetMovementInput() math.Angle {
//	    if h.keys[ebiten.KeyLeft] {
//	        return math.Angle(-5.0)
//	    }
//	    if h.keys[ebiten.KeyRight] {
//	        return math.Angle(5.0)
//	    }
//	    return 0
//	}
type MovementInputHandler interface {
    GetMovementInput() math.Angle
}
```

## Architecture Guidelines

### 3.1 Clean Architecture Principles

Follow the dependency rule: dependencies should only point inward.

```
internal/
├── app/          # Application layer (highest level)
├── game/         # Game logic layer
├── scenes/       # Scene management layer
├── ecs/          # Entity-Component-System layer
├── input/        # Input handling layer
├── ui/           # User interface layer
├── config/       # Configuration layer
├── common/       # Shared interfaces and types
└── errors/       # Error handling layer
```

### 3.2 Interface Segregation Principle

Keep interfaces small and focused:

```go
// ✅ Good: Small, focused interfaces
type MovementInputHandler interface {
    GetMovementInput() math.Angle
}

type ActionInputHandler interface {
    IsQuitPressed() bool
    IsPausePressed() bool
    IsShootPressed() bool
}

// ❌ Bad: Large, monolithic interface
type GameInputHandler interface {
    GetMovementInput() math.Angle
    IsQuitPressed() bool
    IsPausePressed() bool
    IsShootPressed() bool
    GetTouchState() *TouchState
    GetMousePosition() Point
    // ... many more methods
}
```

### 3.3 Dependency Injection

Use dependency injection for all dependencies:

```go
// ✅ Good: Dependency injection
type MovementSystem struct {
    world        donburi.World
    config       *config.GameConfig
    logger       common.Logger
    inputHandler common.GameInputHandler
}

func NewMovementSystem(world donburi.World, config *config.GameConfig, logger common.Logger, inputHandler common.GameInputHandler) *MovementSystem {
    return &MovementSystem{
        world:        world,
        config:       config,
        logger:       logger,
        inputHandler: inputHandler,
    }
}

// ❌ Bad: Creating dependencies inside
func NewMovementSystem(world donburi.World) *MovementSystem {
    return &MovementSystem{
        world:        world,
        config:       config.NewGameConfig(), // ❌ Creating dependency
        logger:       logger.NewLogger(),     // ❌ Creating dependency
        inputHandler: input.NewHandler(),     // ❌ Creating dependency
    }
}
```

## Error Handling

### 4.1 Error Types

Use custom error types for better error classification:

```go
// ✅ Good: Custom error types
type GameError struct {
    Code      ErrorCode
    Message   string
    Cause     error
    Timestamp time.Time
    Context   map[string]interface{}
}

// ✅ Good: Error codes
const (
    AssetNotFound    ErrorCode = "ASSET_NOT_FOUND"
    EntityNotFound   ErrorCode = "ENTITY_NOT_FOUND"
    SystemUpdateFailed ErrorCode = "SYSTEM_UPDATE_FAILED"
)
```

### 4.2 Error Wrapping

Always wrap errors with context:

```go
// ✅ Good: Error wrapping with context
if err != nil {
    return errors.Wrapf(err, errors.AssetLoadFailed, "failed to load sprite %s", spriteName)
}

// ✅ Good: Error builder pattern
err := errors.NewErrorBuilder(errors.ResourceLoadFailed, "failed to load texture").
    WithCause(originalErr).
    WithContext("file_path", filePath).
    WithContext("attempted_size", size).
    Build()
```

### 4.3 Result Types

Use Result types for better error handling:

```go
// ✅ Good: Result type usage
func LoadSprite(name string) common.Result[*ebiten.Image] {
    if sprite, err := loadSpriteFromFile(name); err != nil {
        return common.Err[*ebiten.Image](err)
    }
    return common.Ok(sprite)
}

// Usage
result := LoadSprite("player.png")
if result.IsErr() {
    logger.Error("Failed to load sprite", "error", result.Error())
    return
}
sprite := result.Value()
```

## Testing Standards

### 5.1 Test Structure

Follow the Arrange-Act-Assert pattern:

```go
func TestMovementSystem_PlayerMovement(t *testing.T) {
    // Arrange
    suite := NewTestSuite(t)
    defer suite.Cleanup()
    
    player := suite.CreateTestPlayer(common.Point{X: 100, Y: 100})
    suite.InputHandler().SetMovementInput(math.Angle(5.0))
    
    // Act
    err := suite.MovementSystem().Update(suite.Context(), 1.0/60.0)
    
    // Assert
    require.NoError(t, err)
    suite.AssertPosition(player, common.Point{X: 105, Y: 100})
}
```

### 5.2 Mock Usage

Use mocks for external dependencies:

```go
// ✅ Good: Mock usage
func TestInputHandler_GetMovementInput(t *testing.T) {
    // Arrange
    handler := NewMockInputHandler()
    expectedInput := math.Angle(5.0)
    handler.SetMovementInput(expectedInput)
    
    // Act
    result := handler.GetMovementInput()
    
    // Assert
    assert.Equal(t, expectedInput, result)
}
```

### 5.3 Integration Tests

Write integration tests for system interactions:

```go
func TestIntegration_PlayerMovementAndCollision(t *testing.T) {
    // Arrange
    suite := NewIntegrationTestSuite(t)
    defer suite.Cleanup()
    
    player := suite.CreateTestPlayer(common.Point{X: 100, Y: 100})
    enemy := suite.CreateTestEnemy(common.Point{X: 110, Y: 100})
    
    // Act
    suite.RunGameLoop(10, 1.0/60.0)
    
    // Assert
    suite.AssertEntityExists(player)
    suite.AssertEntityNotExists(enemy) // Should be destroyed by collision
}
```

## Performance Guidelines

### 6.1 Memory Management

- Use object pools for frequently allocated objects
- Pre-allocate slices with known capacity
- Avoid unnecessary allocations in hot paths

```go
// ✅ Good: Pre-allocated slice
func (s *System) Update(ctx context.Context) error {
    entities := make([]donburi.Entity, 0, s.expectedEntityCount)
    // ... populate entities
}

// ❌ Bad: Dynamic allocation
func (s *System) Update(ctx context.Context) error {
    var entities []donburi.Entity // ❌ Will cause allocations
    // ... populate entities
}
```

### 6.2 Benchmarking

Write benchmarks for performance-critical code:

```go
func BenchmarkMovementSystem_Update(b *testing.B) {
    suite := NewBenchmarkSuite()
    defer suite.Cleanup()
    
    entities := suite.CreateBenchmarkEntities(1000)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        for _, entity := range entities {
            // ... update logic
        }
    }
}
```

## Security Considerations

### 7.1 Input Validation

Always validate external input:

```go
// ✅ Good: Input validation
func (s *System) SetPosition(pos common.Point) error {
    if pos.X < 0 || pos.Y < 0 {
        return errors.NewGameError(errors.ValidationFailed, "position cannot be negative")
    }
    if pos.X > s.config.ScreenSize.Width || pos.Y > s.config.ScreenSize.Height {
        return errors.NewGameError(errors.ValidationFailed, "position out of bounds")
    }
    s.position = pos
    return nil
}
```

### 7.2 Resource Limits

Implement resource limits to prevent DoS attacks:

```go
// ✅ Good: Resource limits
const (
    MaxEntities = 10000
    MaxProjectiles = 1000
    MaxEnemies = 500
)

func (s *System) SpawnEntity(entityType EntityType) error {
    if s.entityCount >= MaxEntities {
        return errors.NewGameError(errors.ResourceExhausted, "entity limit reached")
    }
    // ... spawn logic
}
```

## Code Review Checklist

Before submitting code for review, ensure:

- [ ] All exported functions/types are documented
- [ ] Error handling follows the established patterns
- [ ] Tests cover all new functionality
- [ ] No linter errors or warnings
- [ ] Cyclomatic complexity is under 10 for all functions
- [ ] Interfaces follow ISP (max 3-4 methods)
- [ ] Dependencies are injected, not created
- [ ] Context is propagated through function calls
- [ ] Structured logging is used appropriately
- [ ] Performance considerations are addressed
- [ ] Security implications are considered

## Tools and Automation

### Required Tools

- `golangci-lint` for linting
- `gocyclo` for cyclomatic complexity analysis
- `go test` with coverage
- `go vet` for static analysis
- `go fmt` for code formatting

### CI/CD Integration

All code must pass:
- Linting checks
- Unit tests with >80% coverage
- Integration tests
- Performance benchmarks
- Security scans

## Conclusion

These standards ensure code quality, maintainability, and adherence to Go best practices. Regular reviews and updates to these standards help maintain high code quality as the project evolves. 