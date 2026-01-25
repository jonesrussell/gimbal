package behavior

import (
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// StateHandler defines interface for behavior states
type StateHandler interface {
	// Enter is called when transitioning into this state
	Enter(entry *donburi.Entry, data *core.BehaviorStateData)

	// Update is called every frame while in this state
	Update(entry *donburi.Entry, data *core.BehaviorStateData, deltaTime float64)

	// Exit is called when transitioning out of this state
	Exit(entry *donburi.Entry, data *core.BehaviorStateData)

	// NextState determines the next state based on conditions
	NextState(entry *donburi.Entry, data *core.BehaviorStateData) core.BehaviorStateType

	// StateType returns the state this handler manages
	StateType() core.BehaviorStateType
}

// StateRegistry holds all state handlers
type StateRegistry struct {
	handlers map[core.BehaviorStateType]StateHandler
}

// NewStateRegistry creates a new state registry
func NewStateRegistry() *StateRegistry {
	return &StateRegistry{
		handlers: make(map[core.BehaviorStateType]StateHandler),
	}
}

// Register adds a state handler to the registry
func (sr *StateRegistry) Register(handler StateHandler) {
	sr.handlers[handler.StateType()] = handler
}

// Get retrieves a state handler by type
func (sr *StateRegistry) Get(stateType core.BehaviorStateType) StateHandler {
	if handler, exists := sr.handlers[stateType]; exists {
		return handler
	}
	return nil
}

// GetAll returns all registered handlers
func (sr *StateRegistry) GetAll() map[core.BehaviorStateType]StateHandler {
	return sr.handlers
}
