# Refactoring TODO - SOLID Principles & Best Practices

## 🔴 **Priority 1: Critical Issues (Fix First)**

### 1.1 Complete Constants Extraction
**File**: `internal/common/constants.go`
**Issue**: Missing constants from refactoring plan
**Tasks**:
- [x] Add `MovementSpeedDegreesPerFrame = 5`
- [x] Add `MinTouchDuration = 10`
- [x] Add `TouchThreshold = 50`
- [x] Review codebase for any remaining magic numbers

### 1.2 Create System Manager
**File**: `internal/ecs/system_manager.go` (new file)
**Issue**: Hardcoded system execution order, tight coupling
**Tasks**:
- [x] Create `System` interface with `Update()` and `Name()` methods
- [x] Create `SystemManager` struct with systems slice
- [x] Implement `AddSystem()` method
- [x] Implement `UpdateAll()` method with error handling
- [x] Update main game loop to use SystemManager
- [x] Refactor existing systems to implement System interface

### 1.3 Implement Dependency Injection for Input
**File**: `internal/ecs/game.go` and `internal/common/interfaces.go`
**Issue**: ECS layer directly depends on input layer (architectural violation)
**Tasks**:
- [ ] Move input interfaces to `internal/common/interfaces.go`
- [ ] Update ECS game constructor to accept input handler via DI
- [ ] Create input factory/provider in main.go
- [ ] Remove direct input package imports from ECS
- [ ] Update game initialization to use DI pattern

## 🟡 **Priority 2: Moderate Issues**

### 2.1 Extract Game State Management
**File**: `internal/ecs/game_state.go` (new file)
**Issue**: Mixed responsibilities in ECSGame
**Tasks**:
- [x] Create `GameState` struct with isPaused, score, level fields
- [x] Create `GameStateManager` struct with state and events
- [x] Implement `TogglePause()` method with event emission
- [x] Add methods for score and level management
- [x] Integrate with existing ECSGame struct

### 2.2 Create Error Strategy
**File**: `internal/common/errors.go` (new file)
**Issue**: Inconsistent error handling
**Tasks**:
- [x] Create `GameError` struct with Code, Message, Cause fields
- [x] Implement `Error()` method for GameError
- [x] Define standard error constants (ErrAssetNotFound, ErrEntityInvalid, etc.)
- [x] Update existing error handling to use GameError
- [x] Add error codes for different failure scenarios

### 2.3 Create Application Container
**File**: `internal/app/container.go` (new file)
**Issue**: Manual dependency wiring in main.go
**Tasks**:
- [ ] Create `Container` struct for dependency management
- [ ] Implement constructor methods for each dependency
- [ ] Add proper initialization order
- [ ] Implement graceful shutdown handling
- [ ] Move dependency wiring from main.go to container

## 🟢 **Priority 3: Enhancements**

### 3.1 Create Component Registry
**File**: `internal/ecs/component_registry.go` (new file)
**Issue**: Components scattered across files
**Tasks**:
- [ ] Create `ComponentRegistry` struct with components map
- [ ] Implement `Register()` method
- [ ] Implement `Get()` method
- [ ] Add component registration for all existing components
- [ ] Update component creation to use registry

### 3.2 Add System Dependencies
**File**: `internal/ecs/system_dependencies.go` (new file)
**Issue**: No dependency management between systems
**Tasks**:
- [ ] Create `SystemDependency` struct with Name and Dependencies
- [ ] Create `SystemGraph` struct with dependencies map
- [ ] Implement `AddDependency()` method
- [ ] Implement `GetExecutionOrder()` with topological sort
- [ ] Define system dependencies (e.g., MovementSystem before RenderSystem)
- [ ] Integrate with SystemManager

### 3.3 Create Configuration Validation
**File**: `internal/common/config_validator.go` (new file)
**Issue**: No configuration validation
**Tasks**:
- [ ] Create `ConfigValidator` struct
- [ ] Implement `Validate()` method for GameConfig
- [ ] Add validation for screen size (must be positive)
- [ ] Add validation for radius (must be positive)
- [ ] Add validation for other config parameters
- [ ] Integrate validation into game startup

### 3.4 Consider Event-Driven Input Architecture
**File**: `internal/events/input_events.go` (new file)
**Issue**: Tight coupling between input and game logic
**Tasks**:
- [ ] Create `InputEvent` struct with Type and Data fields
- [ ] Create `EventBus` interface with Publish/Subscribe methods
- [ ] Implement simple event bus for input events
- [ ] Update input handler to emit events instead of direct calls
- [ ] Update game systems to subscribe to relevant input events
- [ ] Add event-driven pause/quit handling

## 📊 **Implementation Checklist**

### Phase 1: Core Infrastructure & Architecture
- [x] Complete constants extraction
- [x] Create System Manager
- [x] Update main game loop to use new infrastructure
- [ ] **Fix architectural violation with dependency injection**

### Phase 2: State & Error Management
- [x] Implement Game State Management
- [x] Create Error Strategy
- [x] Update existing code to use new error handling
- [ ] Create application container for DI

### Phase 3: Advanced Features
- [ ] Implement Component Registry
- [ ] Add System Dependencies
- [ ] Create Configuration Validation
- [ ] Consider event-driven input architecture

## 🧪 **Testing Requirements**

### Unit Tests
- [ ] Test System Manager with mock systems
- [ ] Test Game State Manager state transitions
- [ ] Test Error Strategy with various error scenarios
- [ ] **Test dependency injection with mock input handlers**
- [ ] Test Component Registry registration and retrieval
- [ ] Test System Dependencies topological sort
- [ ] Test Configuration Validation with valid/invalid configs
- [ ] Test event bus with multiple subscribers

### Integration Tests
- [ ] Test System Manager with real systems
- [ ] Test Game State integration with existing systems
- [ ] Test error propagation through system chain
- [ ] **Test DI container initialization and shutdown**
- [ ] Test component lifecycle with registry
- [ ] Test event-driven input handling

### Regression Tests
- [ ] Verify existing gameplay functionality works
- [ ] Verify rendering still works correctly
- [ ] Verify input handling still works
- [ ] Verify star field movement still works

## 📈 **Success Criteria**

- [ ] All magic numbers replaced with named constants
- [ ] System execution order is configurable and explicit
- [ ] Game state is centralized and manageable
- [ ] Error handling is consistent and informative
- [ ] **ECS layer no longer directly depends on input layer**
- [ ] **Input interfaces are properly abstracted**
- [ ] **Dependency injection is implemented and tested**
- [ ] Components are centrally registered and managed
- [ ] System dependencies are explicit and validated
- [ ] Configuration is validated at startup
- [ ] All existing functionality preserved
- [ ] Code is more maintainable and testable
- [ ] Performance is not degraded

## 🔍 **Code Review Checklist**

Before marking any task complete:
- [ ] Code follows Go best practices
- [ ] Error handling is appropriate
- [ ] Tests are written and passing
- [ ] Documentation is updated
- [ ] No new magic numbers introduced
- [ ] No tight coupling created
- [ ] Single responsibility principle followed
- [ ] Open/closed principle maintained
- [ ] **Dependency inversion principle followed**
- [ ] **Interfaces are defined where they're used, not where they're implemented**

## 🏗️ **Architecture Notes**

### Dependency Injection Implementation
```go
// internal/common/interfaces.go
type InputHandler interface {
    HandleInput()
    IsKeyPressed(key ebiten.Key) bool
    GetMovementInput() Angle
    IsQuitPressed() bool
    IsPausePressed() bool
}

// internal/ecs/game.go
func NewECSGame(config *common.GameConfig, logger common.Logger, inputHandler common.InputHandler) (*ECSGame, error) {
    // Constructor now accepts dependencies
}

// main.go or internal/app/container.go
func initializeApp() {
    logger := logger.New()
    inputHandler := input.New(logger)
    game, err := ecs.NewECSGame(config, logger, inputHandler)
}
```

### Future Considerations
- **Event-driven architecture** for complex input scenarios
- **Command pattern** for input actions
- **Hexagonal architecture** for multiple input/output adapters
- **Plugin system** for extensible game systems