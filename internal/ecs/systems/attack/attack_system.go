package attack

import (
	"context"
	"math"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// AttackSystem orchestrates complex attack patterns
type AttackSystem struct {
	world        donburi.World
	config       *config.GameConfig
	screenCenter common.Point
	executors    map[core.AttackPatternType]AttackExecutor
}

// AttackExecutor defines interface for attack pattern execution
type AttackExecutor interface {
	Execute(entry *donburi.Entry, data *core.AttackPatternData, deltaTime float64)
	IsComplete(entry *donburi.Entry, data *core.AttackPatternData) bool
}

// NewAttackSystem creates a new attack system
func NewAttackSystem(
	world donburi.World,
	cfg *config.GameConfig,
) *AttackSystem {
	as := &AttackSystem{
		world:  world,
		config: cfg,
		screenCenter: common.Point{
			X: float64(cfg.ScreenSize.Width) / 2,
			Y: float64(cfg.ScreenSize.Height) / 2,
		},
		executors: make(map[core.AttackPatternType]AttackExecutor),
	}

	// Register attack executors
	as.registerExecutors()

	return as
}

// registerExecutors registers all attack pattern executors
func (as *AttackSystem) registerExecutors() {
	as.executors[core.AttackSingleRush] = &SingleRushExecutor{system: as}
	as.executors[core.AttackPairedRush] = &PairedRushExecutor{system: as}
	as.executors[core.AttackLoopbackRush] = &LoopbackRushExecutor{system: as}
	as.executors[core.AttackSuicideDive] = &SuicideDiveExecutor{system: as}
}

// Update processes all entities with active attacks
func (as *AttackSystem) Update(ctx context.Context, deltaTime float64) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Query entities with AttackPattern component and active attacks
	query.NewQuery(
		filter.Contains(core.AttackPattern),
	).Each(as.world, func(entry *donburi.Entry) {
		attackData := core.AttackPattern.Get(entry)
		if attackData.IsActive {
			as.executeAttack(entry, attackData, deltaTime)
		}
	})

	return nil
}

// executeAttack executes the appropriate attack pattern
func (as *AttackSystem) executeAttack(entry *donburi.Entry, data *core.AttackPatternData, deltaTime float64) {
	executor, exists := as.executors[data.PatternType]
	if !exists {
		return
	}

	executor.Execute(entry, data, deltaTime)
}

// GetPlayerPosition finds the player position for targeting
func (as *AttackSystem) GetPlayerPosition() common.Point {
	var playerPos common.Point
	query.NewQuery(
		filter.Contains(core.PlayerTag),
	).Each(as.world, func(entry *donburi.Entry) {
		if entry.HasComponent(core.Position) {
			pos := core.Position.Get(entry)
			playerPos = *pos
		}
	})
	return playerPos
}

// GetScreenCenter returns the screen center
func (as *AttackSystem) GetScreenCenter() common.Point {
	return as.screenCenter
}

// SingleRushExecutor executes single rush attacks
type SingleRushExecutor struct {
	system *AttackSystem
}

func (sre *SingleRushExecutor) Execute(entry *donburi.Entry, data *core.AttackPatternData, deltaTime float64) {
	if !entry.HasComponent(core.Position) {
		return
	}

	pos := core.Position.Get(entry)
	target := data.TargetPosition

	// Move toward target
	dx := target.X - pos.X
	dy := target.Y - pos.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance < 1 {
		return
	}

	speed := data.RushSpeed
	if speed == 0 {
		speed = 300.0
	}

	moveDistance := math.Min(speed*deltaTime, distance)
	pos.X += (dx / distance) * moveDistance
	pos.Y += (dy / distance) * moveDistance
}

func (sre *SingleRushExecutor) IsComplete(entry *donburi.Entry, data *core.AttackPatternData) bool {
	if !entry.HasComponent(core.Position) {
		return true
	}

	pos := core.Position.Get(entry)
	dx := data.TargetPosition.X - pos.X
	dy := data.TargetPosition.Y - pos.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	return distance < 50 // Complete when near target
}

// PairedRushExecutor executes coordinated two-enemy rush attacks
type PairedRushExecutor struct {
	system *AttackSystem
}

func (pre *PairedRushExecutor) Execute(entry *donburi.Entry, data *core.AttackPatternData, deltaTime float64) {
	// Similar to single rush but could coordinate with pair
	if !entry.HasComponent(core.Position) {
		return
	}

	pos := core.Position.Get(entry)
	target := data.TargetPosition

	dx := target.X - pos.X
	dy := target.Y - pos.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance < 1 {
		return
	}

	speed := data.RushSpeed
	if speed == 0 {
		speed = 280.0 // Slightly slower for paired
	}

	moveDistance := math.Min(speed*deltaTime, distance)
	pos.X += (dx / distance) * moveDistance
	pos.Y += (dy / distance) * moveDistance
}

func (pre *PairedRushExecutor) IsComplete(entry *donburi.Entry, data *core.AttackPatternData) bool {
	if !entry.HasComponent(core.Position) {
		return true
	}

	pos := core.Position.Get(entry)
	dx := data.TargetPosition.X - pos.X
	dy := data.TargetPosition.Y - pos.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	return distance < 50
}

// LoopbackRushExecutor executes loopback rush attacks
type LoopbackRushExecutor struct {
	system *AttackSystem
}

func (lre *LoopbackRushExecutor) Execute(entry *donburi.Entry, data *core.AttackPatternData, deltaTime float64) {
	if !entry.HasComponent(core.Position) {
		return
	}

	pos := core.Position.Get(entry)
	center := lre.system.GetScreenCenter()

	// Calculate distance to center
	dx := center.X - pos.X
	dy := center.Y - pos.Y
	distToCenter := math.Sqrt(dx*dx + dy*dy)

	speed := data.RushSpeed
	if speed == 0 {
		speed = 350.0
	}

	if distToCenter > 30 {
		// Rush toward center
		moveDistance := math.Min(speed*deltaTime, distToCenter)
		pos.X += (dx / distToCenter) * moveDistance
		pos.Y += (dy / distToCenter) * moveDistance
	} else {
		// Passed through center, continue outward and loop back
		// This creates the loopback effect
		angle := math.Atan2(pos.Y-center.Y, pos.X-center.X)
		moveDistance := speed * deltaTime
		pos.X += math.Cos(angle) * moveDistance
		pos.Y += math.Sin(angle) * moveDistance
	}
}

func (lre *LoopbackRushExecutor) IsComplete(entry *donburi.Entry, data *core.AttackPatternData) bool {
	if !entry.HasComponent(core.Position) {
		return true
	}

	pos := core.Position.Get(entry)
	center := lre.system.GetScreenCenter()

	// Complete when back at orbit radius
	dx := pos.X - center.X
	dy := pos.Y - center.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	// Complete when passed through center and reached orbit radius
	orbitRadius := float64(lre.system.config.ScreenSize.Height) / 2 * 0.8
	return distance > orbitRadius
}

// SuicideDiveExecutor executes suicide dive attacks
type SuicideDiveExecutor struct {
	system *AttackSystem
}

func (sde *SuicideDiveExecutor) Execute(entry *donburi.Entry, data *core.AttackPatternData, deltaTime float64) {
	if !entry.HasComponent(core.Position) {
		return
	}

	pos := core.Position.Get(entry)
	// Dive at player position
	playerPos := sde.system.GetPlayerPosition()

	dx := playerPos.X - pos.X
	dy := playerPos.Y - pos.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance < 1 {
		return
	}

	speed := data.RushSpeed
	if speed == 0 {
		speed = 400.0 // Fast suicide dive
	}

	moveDistance := math.Min(speed*deltaTime, distance)
	pos.X += (dx / distance) * moveDistance
	pos.Y += (dy / distance) * moveDistance
}

func (sde *SuicideDiveExecutor) IsComplete(entry *donburi.Entry, data *core.AttackPatternData) bool {
	if !entry.HasComponent(core.Position) {
		return true
	}

	pos := core.Position.Get(entry)
	playerPos := sde.system.GetPlayerPosition()

	dx := playerPos.X - pos.X
	dy := playerPos.Y - pos.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	// Complete when very close to player (collision should happen)
	return distance < 20
}
