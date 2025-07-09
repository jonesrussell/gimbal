# Refactoring TODO - SOLID Principles & Best Practices

## 🔴 **Priority 1: Critical Issues (Fix First)**

### 1.1 Complete Constants Extraction
**File**: `internal/common/constants.go`
**Issue**: Missing constants from refactoring plan
**Tasks**:
- [ ] Add `MovementSpeedDegreesPerFrame = 5`
- [ ] Add `MinTouchDuration = 10`
- [ ] Add `TouchThreshold = 50`
- [ ] Review codebase for any remaining magic numbers

### 1.2 Create System Manager
**File**: `internal/ecs/system_manager.go` (new file)
**Issue**: Hardcoded system execution order, tight coupling
**Tasks**:
- [ ] Create `System` interface with `Update()` and `Name()` methods
- [ ] Create `SystemManager` struct with systems slice
- [ ] Implement `AddSystem()` method
- [ ] Implement `UpdateAll()` method with error handling
- [ ] Update main game loop to use SystemManager
- [ ] Refactor existing systems to implement System interface

## 🟡 **Priority 2: Moderate Issues**

### 2.1 Extract Game State Management
**File**: `internal/ecs/game_state.go` (new file)
**Issue**: Mixed responsibilities in ECSGame
**Tasks**:
- [ ] Create `GameState` struct with isPaused, score, level fields
- [ ] Create `GameStateManager` struct with state and events
- [ ] Implement `TogglePause()` method with event emission
- [ ] Add methods for score and level management
- [ ] Integrate with existing ECSGame struct

### 2.2 Create Error Strategy
**File**: `internal/common/errors.go` (new file)
**Issue**: Inconsistent error handling
**Tasks**:
- [ ] Create `GameError` struct with Code, Message, Cause fields
- [ ] Implement `Error()` method for GameError
- [ ] Define standard error constants (ErrAssetNotFound, ErrEntityInvalid, etc.)
- [ ] Update existing error handling to use GameError
- [ ] Add error codes for different failure scenarios

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

## 📊 **Implementation Checklist**

### Phase 1: Core Infrastructure
- [ ] Complete constants extraction
- [ ] Create System Manager
- [ ] Update main game loop to use new infrastructure

### Phase 2: State & Error Management
- [ ] Implement Game State Management
- [ ] Create Error Strategy
- [ ] Update existing code to use new error handling

### Phase 3: Advanced Features
- [ ] Implement Component Registry
- [ ] Add System Dependencies
- [ ] Create Configuration Validation

## 🧪 **Testing Requirements**

### Unit Tests
- [ ] Test System Manager with mock systems
- [ ] Test Game State Manager state transitions
- [ ] Test Error Strategy with various error scenarios
- [ ] Test Component Registry registration and retrieval
- [ ] Test System Dependencies topological sort
- [ ] Test Configuration Validation with valid/invalid configs

### Integration Tests
- [ ] Test System Manager with real systems
- [ ] Test Game State integration with existing systems
- [ ] Test error propagation through system chain
- [ ] Test component lifecycle with registry

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