# Focused Refactoring: Simplify Player Movement

## 🎯 **Goal**: Make player movement implementation simple and maintainable

## 🔧 **Targeted Changes (Minimal Impact)**

### 1. Extract Movement Constants (5 minutes)
**File**: `internal/common/constants.go` (new)
```go
package common

import "math"

const (
    // Player movement
    PlayerMovementSpeed = 5 // degrees per frame
    
    // Angle calculations
    FullCircleDegrees = 360
    HalfCircleDegrees = 180
    DegreesToRadians  = math.Pi / 180
    
    // Facing angle offset (for Gyruss-style gameplay)
    FacingAngleOffset = 180
)
```

**Benefits**: 
- ✅ No more magic numbers
- ✅ Easy to adjust movement speed
- ✅ Clear intent

### 2. Create Player Movement Helper (10 minutes)
**File**: `internal/ecs/player_movement.go` (new)
```go
package ecs

import "github.com/jonesrussell/gimbal/internal/common"

// PlayerMovement handles all player movement logic
type PlayerMovement struct{}

// UpdateOrbitalAngle updates player's orbital angle based on input
func (pm *PlayerMovement) UpdateOrbitalAngle(orb *OrbitalData, inputAngle common.Angle) {
    if inputAngle != 0 {
        orb.OrbitalAngle += inputAngle
        pm.normalizeAngle(&orb.OrbitalAngle)
    }
}

// UpdateFacingAngle calculates facing angle for Gyruss-style gameplay
func (pm *PlayerMovement) UpdateFacingAngle(orb *OrbitalData) {
    orb.FacingAngle = orb.OrbitalAngle + common.FacingAngleOffset
    pm.normalizeAngle(&orb.FacingAngle)
}

// normalizeAngle keeps angle in [0, 360) range
func (pm *PlayerMovement) normalizeAngle(angle *common.Angle) {
    if *angle < 0 {
        *angle += common.FullCircleDegrees
    } else if *angle >= common.FullCircleDegrees {
        *angle -= common.FullCircleDegrees
    }
}
```

**Benefits**:
- ✅ Single place for all player movement logic
- ✅ Easy to test
- ✅ Easy to modify movement behavior

### 3. Simplify PlayerInputSystem (5 minutes)
**File**: `internal/ecs/systems.go` (modify existing)
```go
// PlayerInputSystem handles player input
func PlayerInputSystem(w donburi.World, inputAngle common.Angle) {
    movement := &PlayerMovement{}
    
    query.NewQuery(
        filter.And(
            filter.Contains(PlayerTag),
            filter.Contains(Orbital),
        ),
    ).Each(w, func(entry *donburi.Entry) {
        orb := Orbital.Get(entry)
        
        // Update orbital angle
        movement.UpdateOrbitalAngle(orb, inputAngle)
        
        // Update facing angle
        movement.UpdateFacingAngle(orb)
    })
}
```

**Benefits**:
- ✅ Cleaner, more readable
- ✅ Logic separated from system
- ✅ Easy to test individual functions

### 4. Create Player Factory Helper (5 minutes)
**File**: `internal/ecs/player_factory.go` (new)
```go
package ecs

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/jonesrussell/gimbal/internal/common"
)

// PlayerFactory creates player entities with consistent configuration
type PlayerFactory struct{}

// CreatePlayer creates a player entity with orbital movement
func (pf *PlayerFactory) CreatePlayer(w donburi.World, sprite *ebiten.Image, config *common.GameConfig) donburi.Entity {
    entity := w.Create(PlayerTag, Position, Sprite, Orbital, Size)
    entry := w.Entry(entity)

    // Set initial position at center
    center := common.Point{
        X: float64(config.ScreenSize.Width) / 2,
        Y: float64(config.ScreenSize.Height) / 2,
    }

    // Set up orbital movement - start at bottom (180 degrees)
    orbitalData := OrbitalData{
        Center:       center,
        Radius:       config.Radius,
        OrbitalAngle: common.HalfCircleDegrees, // 180 degrees
        FacingAngle:  0, // Will be calculated by FacingAngleSystem
    }

    // Set components
    Position.SetValue(entry, center)
    Sprite.SetValue(entry, sprite)
    Orbital.SetValue(entry, orbitalData)
    Size.SetValue(entry, config.PlayerSize)

    return entity
}
```

**Benefits**:
- ✅ Consistent player creation
- ✅ Uses constants instead of magic numbers
- ✅ Easy to modify player setup

## 📊 **Implementation Timeline: 25 minutes total**

1. **Constants file** (5 min) - Extract magic numbers
2. **PlayerMovement helper** (10 min) - Centralize movement logic
3. **Simplify PlayerInputSystem** (5 min) - Use helper functions
4. **PlayerFactory helper** (5 min) - Consistent player creation

## 🎯 **Result: Simple Player Movement**

### Before (Complex):
```go
// Multiple systems, hardcoded values, scattered logic
PlayerInputSystem(g.world, inputAngle)
OrbitalMovementSystem(g.world)
FacingAngleSystem(g.world)

// Magic numbers everywhere
orb.OrbitalAngle += inputAngle
if orb.OrbitalAngle < 0 {
    orb.OrbitalAngle += 360
}
orb.FacingAngle = orb.OrbitalAngle + 180
```

### After (Simple):
```go
// Single system with clear logic
PlayerInputSystem(g.world, inputAngle)

// Clean, readable code
movement.UpdateOrbitalAngle(orb, inputAngle)
movement.UpdateFacingAngle(orb)
```

## 🚀 **Benefits for Future Development**

1. **Easy to Add New Movement Types**: Just add methods to `PlayerMovement`
2. **Easy to Test**: Each function can be unit tested
3. **Easy to Debug**: Clear separation of concerns
4. **Easy to Modify**: Change constants or logic in one place
5. **Easy to Extend**: Add new player behaviors without touching systems

## 🤔 **Should We Do This?**

**Yes, if you want**:
- ✅ Simpler player movement code
- ✅ Easier to add new movement features
- ✅ Better testability
- ✅ Clearer intent

**No, if you prefer**:
- ❌ Keep current complexity
- ❌ Don't want to add more files
- ❌ Current code works fine

**My Recommendation**: This focused refactoring will make player movement much simpler and more maintainable with minimal risk and effort. 