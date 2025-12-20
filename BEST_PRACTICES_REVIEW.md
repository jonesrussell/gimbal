# Gimbal Best Practices Review

**Date:** December 20, 2025
**Go Version:** Currently 1.24, recommended upgrade to 1.25
**Codebase:** ~9,097 lines across 80 Go files

---

## Executive Summary

Gimbal is a well-structured Gyruss-style arcade game using Ebiten and Donburi ECS. The codebase demonstrates solid architectural patterns including dependency injection, interface-based design, and ECS architecture. However, there are opportunities for improvement in DRY compliance, SRP adherence, and Go 1.25 adoption.

---

## 1. Go 1.25 Upgrade Opportunities

### Current State
The project uses **Go 1.24**. Go 1.25 was released August 12, 2025 with several relevant features.

### Recommended Upgrades

#### 1.1 Use `testing/synctest` for Concurrent Tests
Go 1.25 introduces `testing/synctest` for testing time-dependent concurrent code:

```go
// Before: Flaky time-based tests
func TestWeaponCooldown(t *testing.T) {
    time.Sleep(100 * time.Millisecond)
    // ...
}

// After: Deterministic testing with synctest
func TestWeaponCooldown(t *testing.T) {
    synctest.Run(func() {
        // Time advances deterministically
    })
}
```

**Applicable files:**
- `internal/ecs/systems/weapon/`
- `internal/ecs/systems/collision/`
- `internal/input/handler_test.go`

#### 1.2 Experimental `encoding/json/v2`
For any JSON configuration loading, consider the experimental json/v2 package for better performance and API.

#### 1.3 Improved Stack Allocation
Go 1.25 allocates more slice backing storage on the stack. Review allocations in hot paths:
- `internal/ecs/systems/collision/system.go:94` - `enemies := make([]donburi.Entity, 0)`
- Consider pre-allocating with capacity estimates

#### 1.4 Update `go.mod`
```diff
-go 1.24
+go 1.25
```

---

## 2. DRY (Don't Repeat Yourself) Violations

### 2.1 Duplicate Error Code Definitions
**File:** `internal/errors/errors.go:92-149` vs `errors.go:154-211`

**Issue:** Error codes are defined twice - once as `const string` and again as `const ErrorCode`:

```go
// Lines 92-149: String constants
const (
    ErrAssetNotFound    = "ASSET_NOT_FOUND"
    ErrAssetLoadFailed  = "ASSET_LOAD_FAILED"
    // ...
)

// Lines 154-211: ErrorCode constants (same values!)
const (
    AssetNotFound    ErrorCode = "ASSET_NOT_FOUND"
    AssetLoadFailed  ErrorCode = "ASSET_LOAD_FAILED"
    // ...
)
```

**Fix:** Remove the string constants (lines 92-149) and use only `ErrorCode` type:
```go
// Keep only typed constants
const (
    AssetNotFound    ErrorCode = "ASSET_NOT_FOUND"
    // ...
)
```

### 2.2 Repeated Context Map Initialization
**Files:** `internal/errors/errors.go` (multiple locations)

**Issue:** `make(map[string]interface{})` pattern repeated in 5+ places:
- Line 34, 44, 58, 68, 238, 248, 263

**Fix:** Extract to helper:
```go
func newContextMap() map[string]interface{} {
    return make(map[string]interface{})
}
```

### 2.3 Dead Code Comments Pattern
**Files:** Multiple files with "removed - dead code" comments

**Locations:**
- `internal/app/container.go:139` - `// GetInputHandler removed - dead code`
- `internal/app/container.go:148` - `// IsInitialized removed - dead code`
- `internal/app/container.go:181` - `// SetInputHandler removed - dead code`
- `internal/config/config.go:59,61,63,87`

**Fix:** Delete these comments entirely. Version control tracks removed code.

### 2.4 Duplicate AABB Collision Logic Potential
**File:** `internal/ecs/systems/collision/detection.go`

The bounding box calculation pattern could be extracted:
```go
// Current: Inline calculation
left1 := pos1.X
right1 := pos1.X + float64(size1.Width)
// ...

// Suggested: Extract to method
type BoundingBox struct {
    Left, Right, Top, Bottom float64
}

func NewBoundingBox(pos common.Point, size config.Size) BoundingBox {
    return BoundingBox{
        Left:   pos.X,
        Right:  pos.X + float64(size.Width),
        Top:    pos.Y,
        Bottom: pos.Y + float64(size.Height),
    }
}
```

---

## 3. SRP (Single Responsibility Principle) Violations

### 3.1 GameInputHandler Interface Explosion
**File:** `internal/common/interfaces.go:112-120`

**Issue:** `GameInputHandler` composes 7 interfaces, making it unwieldy:

```go
type GameInputHandler interface {
    InputHandler          // Core input handling
    MovementInputHandler  // Movement
    ActionInputHandler    // Actions (pause, shoot, quit)
    TouchInputHandler     // Touch input
    MouseInputHandler     // Mouse input
    EventInputHandler     // Event tracking
    TestableInputHandler  // Testing utilities
}
```

**Fix:** Consider interface segregation:
```go
// For game systems that only need movement
type MovementController interface {
    GetMovementInput() math.Angle
}

// For systems that need actions
type ActionController interface {
    IsPausePressed() bool
    IsShootPressed() bool
}

// Full interface only where truly needed
type FullInputHandler interface {
    MovementController
    ActionController
    TouchInputHandler
    // ...
}
```

### 3.2 SceneManager Has Too Many Responsibilities
**File:** `internal/scenes/manager.go`

**Current Responsibilities:**
1. Scene lifecycle management (Enter/Exit)
2. Scene switching
3. Debug rendering coordination
4. Resource management (optimizer, image pool)
5. Health system reference holding
6. Input handler reference
7. Resume callback management

**Fix:** Extract concerns:
```go
// SceneManager - only scene lifecycle
type SceneManager struct {
    scenes       map[SceneType]Scene
    currentScene Scene
}

// DebugOverlay - separate debug concerns
type DebugOverlay struct {
    renderer *debug.DebugRenderer
}

// ResourceContext - resource sharing
type ResourceContext struct {
    optimizer *core.RenderOptimizer
    imagePool *core.ImagePool
}
```

### 3.3 ECSGame Struct Has 15+ Fields
**File:** `internal/game/game.go:27-73`

**Issue:** `ECSGame` manages too many concerns directly:
- World management
- Event system
- Resource management
- State management (stateManager, scoreManager, levelManager)
- Scene management
- 5 gameplay systems
- UI
- Debug tools
- Performance monitoring

**Fix:** Group related fields into sub-structs:
```go
type ECSGame struct {
    world  donburi.World
    config *config.GameConfig

    // Group related concerns
    managers   *GameManagers   // score, level, state, resource
    systems    *GameSystems    // health, movement, collision, enemy, weapon
    rendering  *RenderContext  // optimizer, pool, debugger

    inputHandler common.GameInputHandler
    logger       common.Logger
}

type GameSystems struct {
    Health    *healthsys.HealthSystem
    Movement  *movement.MovementSystem
    Collision *collision.CollisionSystem
    Enemy     *enemysys.EnemySystem
    Weapon    *weaponsys.WeaponSystem
}
```

### 3.4 healthSystem as interface{}
**File:** `internal/scenes/manager.go:59,222-228`

**Issue:** Using `interface{}` loses type safety:
```go
healthSystem interface{} // Health system interface for scenes to access

func (sceneMgr *SceneManager) SetHealthSystem(healthSystem interface{}) {
    sceneMgr.healthSystem = healthSystem
}
```

**Fix:** Define a proper interface:
```go
type HealthProvider interface {
    GetPlayerHealth() (current, maximum int)
    ApplyDamage(entity donburi.Entity, damage int)
}
```

---

## 4. Additional Best Practice Improvements

### 4.1 Magic Numbers Should Be Constants
**File:** `internal/input/handler.go:12-14`

```go
const (
    PlayerMovementSpeed = 5  // Good!
    MinTouchDuration    = 10 // Good!
    TouchThreshold      = 50 // Good!
)
```

**But missing in:**
- `internal/ecs/core/components.go:76` - `InvincibilityDuration: 2.0`
- `internal/game/game.go:137` - `deltaTime := 1.0 / 60.0`
- `internal/game/game.go:237` - `5*time.Millisecond`
- `internal/ecs/systems/collision/system.go:56` - `16*time.Millisecond`

**Fix:**
```go
// internal/config/constants.go
const (
    DefaultInvincibilityDuration = 2.0 * time.Second
    TargetFPS                    = 60
    FrameBudget                  = time.Second / TargetFPS
    SlowSystemThreshold          = 5 * time.Millisecond
)
```

### 4.2 Context.Background() in Game Loop
**File:** `internal/game/game.go:281`

**Issue:**
```go
ctx := context.Background()
if err := g.updateGameplaySystems(ctx); err != nil {
```

**Fix:** Pass a cancellation context from the game lifecycle:
```go
type ECSGame struct {
    ctx    context.Context
    cancel context.CancelFunc
    // ...
}

func (g *ECSGame) Update() error {
    if err := g.updateGameplaySystems(g.ctx); err != nil {
```

### 4.3 Error Wrapping Consistency
Use `fmt.Errorf` with `%w` consistently. The codebase mostly does this well, but verify all error returns follow the pattern:

```go
// Good pattern (already used)
return fmt.Errorf("failed to initialize: %w", err)
```

### 4.4 Logging Consistency
The codebase uses structured logging well, but some debug logs run every frame:

**File:** `internal/game/game.go:172`
```go
g.logger.Debug("ECS systems updated", "delta", deltaTime)
```

This logs 60 times per second. Consider:
```go
if g.frameCount%60 == 0 { // Log once per second
    g.logger.Debug("ECS systems updated", "delta", deltaTime)
}
```

### 4.5 Collision Detection Timeout
**File:** `internal/ecs/systems/collision/system.go:56`

```go
ctx, cancel := context.WithTimeout(ctx, 16*time.Millisecond)
```

16ms is the entire frame budget at 60 FPS. If collision takes the full budget, the game will stutter. Consider:
- Using a smaller timeout (8ms)
- Or implementing spatial partitioning to ensure fast collision detection

---

## 5. Testing Improvements

### 5.1 TestableInputHandler Pattern
**File:** `internal/common/interfaces.go:103-107`

Good pattern for testability! Extend this to other systems:
```go
type TestableEnemySystem interface {
    SpawnEnemyAt(pos common.Point) donburi.Entity
    ClearAllEnemies()
}
```

### 5.2 Consider Table-Driven Tests
For collision detection, input handling, etc.:
```go
func TestCollisionDetection(t *testing.T) {
    tests := []struct {
        name     string
        pos1     common.Point
        size1    config.Size
        pos2     common.Point
        size2    config.Size
        expected bool
    }{
        {"overlapping", Point{0, 0}, Size{10, 10}, Point{5, 5}, Size{10, 10}, true},
        {"separate", Point{0, 0}, Size{10, 10}, Point{20, 20}, Size{10, 10}, false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ...
        })
    }
}
```

---

## 6. Architecture Strengths (Keep These!)

1. **Dependency Injection Container** - `internal/app/container.go` is well-designed
2. **Options Pattern** - `internal/config/config.go` uses functional options
3. **ECS Architecture** - Clean separation with Donburi
4. **Scene Management** - Good state machine pattern
5. **Structured Logging** - Consistent use of zap
6. **Custom Error Types** - `internal/errors/errors.go` provides rich context
7. **Interface Segregation** - Many small interfaces in `internal/common/interfaces.go`
8. **Config Validation** - Validation at startup prevents runtime issues

---

## 7. Priority Action Items

| Priority | Item | Effort | Impact |
|----------|------|--------|--------|
| High | Upgrade to Go 1.25 | Low | Medium |
| High | Remove duplicate error code definitions | Low | Medium |
| High | Delete "removed - dead code" comments | Low | Low |
| Medium | Extract magic numbers to constants | Low | Medium |
| Medium | Add proper interface for healthSystem | Low | High |
| Medium | Reduce ECSGame responsibilities | Medium | High |
| Low | Split GameInputHandler interface | Medium | Medium |
| Low | Reduce SceneManager responsibilities | High | Medium |
| Low | Implement spatial partitioning for collision | High | High |

---

## Sources

- [Go 1.25 Release Notes](https://tip.golang.org/doc/go1.25)
- [Go 1.25 is Released - Official Blog](https://go.dev/blog/go1.25)
- [Go 1.25 Interactive Tour](https://antonz.org/go-1-25/)
- [Go 1.25 Upgrade Guide - Leapcell](https://leapcell.io/blog/go-1-25-upgrade-guide)
