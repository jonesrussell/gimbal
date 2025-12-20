package events

import (
	"context"
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/math"
)

// Event data structures
type PlayerMovedEvent struct {
	Position common.Point
	Angle    math.Angle
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

type PlayerDamagedEvent struct {
	PlayerEntity   donburi.Entity
	Damage         int
	RemainingLives int
}

type GameOverEvent struct {
	Reason string
}

type LifeAddedEvent struct {
	PlayerEntity donburi.Entity
	NewLives     int
}

type EnemyDestroyedEvent struct {
	Entity donburi.Entity
	Points int
	Time   int64
}

// Event types for the game
var (
	PlayerMovedEventType    = events.NewEventType[PlayerMovedEvent]()
	StarCollectedEventType  = events.NewEventType[StarCollectedEvent]()
	ScoreChangedEventType   = events.NewEventType[ScoreChangedEvent]()
	GameStateEventType      = events.NewEventType[GameStateEvent]()
	PlayerDamagedEventType  = events.NewEventType[PlayerDamagedEvent]()
	GameOverEventType       = events.NewEventType[GameOverEvent]()
	LifeAddedEventType      = events.NewEventType[LifeAddedEvent]()
	EnemyDestroyedEventType = events.NewEventType[EnemyDestroyedEvent]()
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
func (evt *EventSystem) EmitPlayerMoved(pos common.Point, angle math.Angle) {
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

// EmitPlayerDamaged emits a player damaged event
func (evt *EventSystem) EmitPlayerDamaged(playerEntity donburi.Entity, damage, remainingLives int) {
	PlayerDamagedEventType.Publish(evt.world, PlayerDamagedEvent{
		PlayerEntity:   playerEntity,
		Damage:         damage,
		RemainingLives: remainingLives,
	})
}

// EmitGameOver emits a game over event
func (evt *EventSystem) EmitGameOver() {
	GameOverEventType.Publish(evt.world, GameOverEvent{Reason: "No lives remaining"})
}

// EmitLifeAdded emits a life added event
func (evt *EventSystem) EmitLifeAdded(playerEntity donburi.Entity, newLives int) {
	LifeAddedEventType.Publish(evt.world, LifeAddedEvent{
		PlayerEntity: playerEntity,
		NewLives:     newLives,
	})
}

// EmitEnemyDestroyed emits an event when an enemy is destroyed with points awarded
func (evt *EventSystem) EmitEnemyDestroyed(ctx context.Context, entity donburi.Entity, points int) error {
	EnemyDestroyedEventType.Publish(evt.world, EnemyDestroyedEvent{
		Entity: entity,
		Points: points,
		Time:   time.Now().Unix(),
	})
	return nil
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

// SubscribeToPlayerDamaged subscribes to player damaged events
func (evt *EventSystem) SubscribeToPlayerDamaged(callback events.Subscriber[PlayerDamagedEvent]) {
	PlayerDamagedEventType.Subscribe(evt.world, callback)
}

// SubscribeToGameOver subscribes to game over events
func (evt *EventSystem) SubscribeToGameOver(callback events.Subscriber[GameOverEvent]) {
	GameOverEventType.Subscribe(evt.world, callback)
}

// SubscribeToLifeAdded subscribes to life added events
func (evt *EventSystem) SubscribeToLifeAdded(callback events.Subscriber[LifeAddedEvent]) {
	LifeAddedEventType.Subscribe(evt.world, callback)
}

// SubscribeToEnemyDestroyed subscribes to enemy destroyed events
func (evt *EventSystem) SubscribeToEnemyDestroyed(callback events.Subscriber[EnemyDestroyedEvent]) {
	EnemyDestroyedEventType.Subscribe(evt.world, callback)
}

// ProcessEvents processes all pending events
func (evt *EventSystem) ProcessEvents() {
	PlayerMovedEventType.ProcessEvents(evt.world)
	StarCollectedEventType.ProcessEvents(evt.world)
	ScoreChangedEventType.ProcessEvents(evt.world)
	GameStateEventType.ProcessEvents(evt.world)
	PlayerDamagedEventType.ProcessEvents(evt.world)
	GameOverEventType.ProcessEvents(evt.world)
	LifeAddedEventType.ProcessEvents(evt.world)
	EnemyDestroyedEventType.ProcessEvents(evt.world)
}
