# Refactoring Plan: SOLID Principles & Best Practices

## 🔴 **Priority 1: Critical Issues (Fix First)**

### 1.1 Extract Rotation Logic (DRY Violation)
**File**: `internal/ecs/systems.go`
**Issue**: Duplicate rotation logic in RenderSystem
**Solution**: Create `applyRotation()` helper function

```go
func applyRotation(op *ebiten.DrawImageOptions, entry *donburi.Entry, angle float64) {
    var centerX, centerY float64
    if entry.HasComponent(Size) {
        size := Size.Get(entry)
        centerX = float64(size.Width) / 2
        centerY = float64(size.Height) / 2
    } else {
        bounds := (*Sprite.Get(entry)).Bounds()
        centerX = float64(bounds.Dx()) / 2
        centerY = float64(bounds.Dy()) / 2
    }
    
    op.GeoM.Translate(-centerX, -centerY)
    op.GeoM.Rotate(angle * math.Pi / 180)
    op.GeoM.Translate(centerX, centerY)
}
```

### 1.2 Extract Constants (Magic Numbers)
**File**: `internal/common/constants.go` (new file)
**Issue**: Magic numbers scattered throughout code
**Solution**: Centralize constants

```go
package common

const (
    // Angle constants
    FullCircleDegrees = 360
    HalfCircleDegrees = 180
    DegreesToRadians  = math.Pi / 180
    
    // Movement constants
    MovementSpeedDegreesPerFrame = 5
    MinTouchDuration            = 10
    TouchThreshold              = 50
)
```

### 1.3 Create System Manager (Tight Coupling)
**File**: `internal/ecs/system_manager.go` (new file)
**Issue**: Hardcoded system execution order
**Solution**: Centralized system management

```go
type SystemManager struct {
    systems []System
}

type System interface {
    Update(world donburi.World, args ...interface{}) error
    Name() string
}

func (sm *SystemManager) AddSystem(system System) {
    sm.systems = append(sm.systems, system)
}

func (sm *SystemManager) UpdateAll(world donburi.World, args ...interface{}) error {
    for _, system := range sm.systems {
        if err := system.Update(world, args...); err != nil {
            return fmt.Errorf("system %s failed: %w", system.Name(), err)
        }
    }
    return nil
}
```

## 🟡 **Priority 2: Moderate Issues**

### 2.1 Split RenderSystem (SRP Violation)
**Files**: `internal/ecs/systems.go`, `internal/ecs/render_systems.go` (new)
**Issue**: RenderSystem doing too much
**Solution**: Split into specialized systems

```go
// TransformSystem - handles scaling and positioning
func TransformSystem(w donburi.World, screen *ebiten.Image)

// RotationSystem - handles rotation only
func RotationSystem(w donburi.World, screen *ebiten.Image)

// SpriteRenderSystem - handles sprite drawing
func SpriteRenderSystem(w donburi.World, screen *ebiten.Image)
```

### 2.2 Extract Game State Management
**File**: `internal/ecs/game_state.go` (new file)
**Issue**: Mixed responsibilities in ECSGame
**Solution**: Separate state management

```go
type GameState struct {
    isPaused bool
    score    int
    level    int
}

type GameStateManager struct {
    state  *GameState
    events *EventSystem
}

func (gsm *GameStateManager) TogglePause() {
    gsm.state.isPaused = !gsm.state.isPaused
    if gsm.state.isPaused {
        gsm.events.EmitGamePaused()
    } else {
        gsm.events.EmitGameResumed()
    }
}
```

### 2.3 Create Error Strategy
**File**: `internal/common/errors.go` (new file)
**Issue**: Inconsistent error handling
**Solution**: Standardized error types

```go
type GameError struct {
    Code    string
    Message string
    Cause   error
}

func (e *GameError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

var (
    ErrAssetNotFound = &GameError{Code: "ASSET_NOT_FOUND", Message: "asset not found"}
    ErrEntityInvalid = &GameError{Code: "ENTITY_INVALID", Message: "entity is invalid"}
)
```

## 🟢 **Priority 3: Enhancements**

### 3.1 Create Component Registry
**File**: `internal/ecs/component_registry.go` (new file)
**Issue**: Components scattered across files
**Solution**: Centralized component management

```go
type ComponentRegistry struct {
    components map[string]interface{}
}

func (cr *ComponentRegistry) Register(name string, component interface{}) {
    cr.components[name] = component
}

func (cr *ComponentRegistry) Get(name string) (interface{}, bool) {
    comp, exists := cr.components[name]
    return comp, exists
}
```

### 3.2 Add System Dependencies
**File**: `internal/ecs/system_dependencies.go` (new file)
**Issue**: No dependency management between systems
**Solution**: Dependency injection for systems

```go
type SystemDependency struct {
    Name         string
    Dependencies []string
}

type SystemGraph struct {
    dependencies map[string]*SystemDependency
}

func (sg *SystemGraph) AddDependency(name string, deps ...string) {
    sg.dependencies[name] = &SystemDependency{
        Name:         name,
        Dependencies: deps,
    }
}

func (sg *SystemGraph) GetExecutionOrder() ([]string, error) {
    // Topological sort implementation
}
```

### 3.3 Create Configuration Validation
**File**: `internal/common/config_validator.go` (new file)
**Issue**: No configuration validation
**Solution**: Validate configuration at startup

```go
type ConfigValidator struct{}

func (cv *ConfigValidator) Validate(config *GameConfig) error {
    if config.ScreenSize.Width <= 0 || config.ScreenSize.Height <= 0 {
        return &GameError{Code: "INVALID_SCREEN_SIZE", Message: "screen size must be positive"}
    }
    if config.Radius <= 0 {
        return &GameError{Code: "INVALID_RADIUS", Message: "radius must be positive"}
    }
    return nil
}
```

## 📊 **Implementation Priority**

1. **Week 1**: Priority 1 issues (Critical)
   - Extract rotation logic
   - Create constants file
   - Implement system manager

2. **Week 2**: Priority 2 issues (Moderate)
   - Split RenderSystem
   - Extract game state management
   - Implement error strategy

3. **Week 3**: Priority 3 issues (Enhancements)
   - Component registry
   - System dependencies
   - Configuration validation

## 🧪 **Testing Strategy**

1. **Unit Tests**: Each extracted function/component
2. **Integration Tests**: System interactions
3. **Regression Tests**: Ensure existing functionality works
4. **Performance Tests**: Measure impact of refactoring

## 📈 **Expected Benefits**

1. **Maintainability**: Easier to modify and extend
2. **Testability**: Smaller, focused components
3. **Reusability**: Extracted logic can be reused
4. **Readability**: Clearer separation of concerns
5. **Performance**: Better system organization 