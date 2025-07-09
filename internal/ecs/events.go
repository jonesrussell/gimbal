package ecs

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"

	"github.com/jonesrussell/gimbal/internal/common"
)

// Event data structures
type PlayerMovedEvent struct {
	Position common.Point
	Angle    common.Angle
}

type StarCollectedEvent struct {
	Position common.Point
	Score    int
}

type ScoreChangedEvent struct {
	OldScore int
	NewScore int
	Delta    int
}

type GameStateEvent struct {
	IsPaused bool
}

// Event types for the game
var (
	PlayerMovedEventType   = events.NewEventType[PlayerMovedEvent]()
	StarCollectedEventType = events.NewEventType[StarCollectedEvent]()
	ScoreChangedEventType  = events.NewEventType[ScoreChangedEvent]()
	GameStateEventType     = events.NewEventType[GameStateEvent]()
)

// EventSystem manages game events
type EventSystem struct {
	world donburi.World
}

// NewEventSystem creates a new event system
func NewEventSystem(world donburi.World) *EventSystem {
	return &EventSystem{
		world: world,
	}
}

// EmitPlayerMoved emits a player moved event
func (es *EventSystem) EmitPlayerMoved(pos common.Point, angle common.Angle) {
	PlayerMovedEventType.Publish(es.world, PlayerMovedEvent{
		Position: pos,
		Angle:    angle,
	})
}

// EmitStarCollected emits a star collected event
func (es *EventSystem) EmitStarCollected(pos common.Point, score int) {
	StarCollectedEventType.Publish(es.world, StarCollectedEvent{
		Position: pos,
		Score:    score,
	})
}

// EmitScoreChanged emits a score changed event
func (es *EventSystem) EmitScoreChanged(oldScore, newScore int) {
	ScoreChangedEventType.Publish(es.world, ScoreChangedEvent{
		OldScore: oldScore,
		NewScore: newScore,
		Delta:    newScore - oldScore,
	})
}

// EmitGamePaused emits a game paused event
func (es *EventSystem) EmitGamePaused() {
	GameStateEventType.Publish(es.world, GameStateEvent{IsPaused: true})
}

// EmitGameResumed emits a game resumed event
func (es *EventSystem) EmitGameResumed() {
	GameStateEventType.Publish(es.world, GameStateEvent{IsPaused: false})
}

// SubscribeToPlayerMoved subscribes to player moved events
func (es *EventSystem) SubscribeToPlayerMoved(callback events.Subscriber[PlayerMovedEvent]) {
	PlayerMovedEventType.Subscribe(es.world, callback)
}

// SubscribeToStarCollected subscribes to star collected events
func (es *EventSystem) SubscribeToStarCollected(callback events.Subscriber[StarCollectedEvent]) {
	StarCollectedEventType.Subscribe(es.world, callback)
}

// SubscribeToScoreChanged subscribes to score changed events
func (es *EventSystem) SubscribeToScoreChanged(callback events.Subscriber[ScoreChangedEvent]) {
	ScoreChangedEventType.Subscribe(es.world, callback)
}

// SubscribeToGameState subscribes to game state events
func (es *EventSystem) SubscribeToGameState(callback events.Subscriber[GameStateEvent]) {
	GameStateEventType.Subscribe(es.world, callback)
}

// ProcessEvents processes all pending events
func (es *EventSystem) ProcessEvents() {
	PlayerMovedEventType.ProcessEvents(es.world)
	StarCollectedEventType.ProcessEvents(es.world)
	ScoreChangedEventType.ProcessEvents(es.world)
	GameStateEventType.ProcessEvents(es.world)
}
