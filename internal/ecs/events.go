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
func (evt *EventSystem) EmitPlayerMoved(pos common.Point, angle common.Angle) {
	PlayerMovedEventType.Publish(evt.world, PlayerMovedEvent{
		Position: pos,
		Angle:    angle,
	})
}

// EmitStarCollected emits a star collected event
func (evt *EventSystem) EmitStarCollected(pos common.Point, score int) {
	StarCollectedEventType.Publish(evt.world, StarCollectedEvent{
		Position: pos,
		Score:    score,
	})
}

// EmitScoreChanged emits a score changed event
func (evt *EventSystem) EmitScoreChanged(oldScore, newScore int) {
	ScoreChangedEventType.Publish(evt.world, ScoreChangedEvent{
		OldScore: oldScore,
		NewScore: newScore,
		Delta:    newScore - oldScore,
	})
}

// EmitGamePaused emits a game paused event
func (evt *EventSystem) EmitGamePaused() {
	GameStateEventType.Publish(evt.world, GameStateEvent{IsPaused: true})
}

// EmitGameResumed emits a game resumed event
func (evt *EventSystem) EmitGameResumed() {
	GameStateEventType.Publish(evt.world, GameStateEvent{IsPaused: false})
}

// SubscribeToPlayerMoved subscribes to player moved events
func (evt *EventSystem) SubscribeToPlayerMoved(callback events.Subscriber[PlayerMovedEvent]) {
	PlayerMovedEventType.Subscribe(evt.world, callback)
}

// SubscribeToStarCollected subscribes to star collected events
func (evt *EventSystem) SubscribeToStarCollected(callback events.Subscriber[StarCollectedEvent]) {
	StarCollectedEventType.Subscribe(evt.world, callback)
}

// SubscribeToScoreChanged subscribes to score changed events
func (evt *EventSystem) SubscribeToScoreChanged(callback events.Subscriber[ScoreChangedEvent]) {
	ScoreChangedEventType.Subscribe(evt.world, callback)
}

// SubscribeToGameState subscribes to game state events
func (evt *EventSystem) SubscribeToGameState(callback events.Subscriber[GameStateEvent]) {
	GameStateEventType.Subscribe(evt.world, callback)
}

// ProcessEvents processes all pending events
func (evt *EventSystem) ProcessEvents() {
	PlayerMovedEventType.ProcessEvents(evt.world)
	StarCollectedEventType.ProcessEvents(evt.world)
	ScoreChangedEventType.ProcessEvents(evt.world)
	GameStateEventType.ProcessEvents(evt.world)
}
