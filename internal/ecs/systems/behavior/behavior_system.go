package behavior

import (
	"context"
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/dbg"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// BehaviorSystem manages enemy behavior state machines
type BehaviorSystem struct {
	world         donburi.World
	config        *config.GameConfig
	logger        common.Logger
	stateRegistry *StateRegistry
	screenCenter  common.Point
}

// NewBehaviorSystem creates a new behavior system
func NewBehaviorSystem(
	world donburi.World,
	cfg *config.GameConfig,
	logger common.Logger,
) *BehaviorSystem {
	bs := &BehaviorSystem{
		world:         world,
		config:        cfg,
		logger:        logger,
		stateRegistry: NewStateRegistry(),
		screenCenter: common.Point{
			X: float64(cfg.ScreenSize.Width) / 2,
			Y: float64(cfg.ScreenSize.Height) / 2,
		},
	}

	// Register all state handlers
	bs.registerStates()

	return bs
}

// registerStates registers all state handlers
func (bs *BehaviorSystem) registerStates() {
	bs.stateRegistry.Register(NewEnteringState(bs.logger))
	bs.stateRegistry.Register(NewOrbitingState(bs.config, bs.logger))
	bs.stateRegistry.Register(NewAttackingState(bs.config, bs.logger))
	bs.stateRegistry.Register(NewRetreatingState(bs.config, bs.logger))
	bs.stateRegistry.Register(NewHoveringState(bs.config, bs.logger))
}

// Update processes all entities with behavior states
func (bs *BehaviorSystem) Update(ctx context.Context, deltaTime float64) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Query all entities with BehaviorState component
	query.NewQuery(
		filter.Contains(core.BehaviorState),
	).Each(bs.world, func(entry *donburi.Entry) {
		bs.updateEntity(entry, deltaTime)
	})

	return nil
}

// updateEntity updates a single entity's behavior state
func (bs *BehaviorSystem) updateEntity(entry *donburi.Entry, deltaTime float64) {
	behaviorData := core.BehaviorState.Get(entry)

	// Get current state handler
	handler := bs.stateRegistry.Get(behaviorData.CurrentState)
	if handler == nil {
		bs.logger.Warn("No handler for state",
			"state", behaviorData.CurrentState,
			"entity", entry.Entity())
		return
	}

	// Update state time
	behaviorData.StateTime += time.Duration(deltaTime * float64(time.Second))

	// Update current state
	handler.Update(entry, behaviorData, deltaTime)

	// Check for state transition
	nextState := handler.NextState(entry, behaviorData)
	if nextState != behaviorData.CurrentState {
		bs.transitionState(entry, behaviorData, handler, nextState)
	}

	// Save updated behavior data
	core.BehaviorState.SetValue(entry, *behaviorData)
}

// transitionState handles state transitions
func (bs *BehaviorSystem) transitionState(
	entry *donburi.Entry,
	data *core.BehaviorStateData,
	currentHandler StateHandler,
	nextState core.BehaviorStateType,
) {
	// Exit current state
	currentHandler.Exit(entry, data)

	// Update state
	data.PreviousState = data.CurrentState
	data.CurrentState = nextState
	data.StateTime = 0

	// Get next state handler
	nextHandler := bs.stateRegistry.Get(nextState)
	if nextHandler == nil {
		bs.logger.Warn("No handler for next state",
			"state", nextState,
			"entity", entry.Entity())
		return
	}

	// Enter next state
	nextHandler.Enter(entry, data)

	dbg.Log(dbg.State, "Behavior state transition %v → %v", data.PreviousState, data.CurrentState)
}

// GetScreenCenter returns the screen center
func (bs *BehaviorSystem) GetScreenCenter() common.Point {
	return bs.screenCenter
}

// ForceTransition forces an entity to a specific state
func (bs *BehaviorSystem) ForceTransition(entry *donburi.Entry, targetState core.BehaviorStateType) {
	if !entry.HasComponent(core.BehaviorState) {
		return
	}

	behaviorData := core.BehaviorState.Get(entry)
	currentHandler := bs.stateRegistry.Get(behaviorData.CurrentState)

	if currentHandler != nil {
		bs.transitionState(entry, behaviorData, currentHandler, targetState)
		core.BehaviorState.SetValue(entry, *behaviorData)
	}
}
