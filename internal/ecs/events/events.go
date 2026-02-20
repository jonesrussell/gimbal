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

type LevelChangedEvent struct {
	OldLevel int
	NewLevel int
}

// Stage progression events
type WaveStartedEvent struct {
	WaveIndex int
}

type WaveCompletedEvent struct {
	WaveIndex int
}

type BossSpawnRequestedEvent struct{}

type BossSpawnedEvent struct{}

type BossDefeatedEvent struct{}

type StageCompletedEvent struct {
	StageNumber int
}

// Event types for the game
var (
	PlayerMovedEventType        = events.NewEventType[PlayerMovedEvent]()
	StarCollectedEventType      = events.NewEventType[StarCollectedEvent]()
	ScoreChangedEventType       = events.NewEventType[ScoreChangedEvent]()
	GameStateEventType          = events.NewEventType[GameStateEvent]()
	PlayerDamagedEventType      = events.NewEventType[PlayerDamagedEvent]()
	GameOverEventType           = events.NewEventType[GameOverEvent]()
	LifeAddedEventType          = events.NewEventType[LifeAddedEvent]()
	EnemyDestroyedEventType     = events.NewEventType[EnemyDestroyedEvent]()
	LevelChangedEventType       = events.NewEventType[LevelChangedEvent]()
	WaveStartedEventType        = events.NewEventType[WaveStartedEvent]()
	WaveCompletedEventType      = events.NewEventType[WaveCompletedEvent]()
	BossSpawnRequestedEventType = events.NewEventType[BossSpawnRequestedEvent]()
	BossSpawnedEventType        = events.NewEventType[BossSpawnedEvent]()
	BossDefeatedEventType       = events.NewEventType[BossDefeatedEvent]()
	StageCompletedEventType     = events.NewEventType[StageCompletedEvent]()
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

// EmitLevelChanged emits a level changed event
func (evt *EventSystem) EmitLevelChanged(oldLevel, newLevel int) {
	LevelChangedEventType.Publish(evt.world, LevelChangedEvent{
		OldLevel: oldLevel,
		NewLevel: newLevel,
	})
}

// SubscribeToLevelChanged subscribes to level changed events
func (evt *EventSystem) SubscribeToLevelChanged(callback events.Subscriber[LevelChangedEvent]) {
	LevelChangedEventType.Subscribe(evt.world, callback)
}

// EmitWaveStarted emits a wave started event
func (evt *EventSystem) EmitWaveStarted(waveIndex int) {
	WaveStartedEventType.Publish(evt.world, WaveStartedEvent{WaveIndex: waveIndex})
}

// SubscribeToWaveStarted subscribes to wave started events
func (evt *EventSystem) SubscribeToWaveStarted(callback events.Subscriber[WaveStartedEvent]) {
	WaveStartedEventType.Subscribe(evt.world, callback)
}

// EmitWaveCompleted emits a wave completed event
func (evt *EventSystem) EmitWaveCompleted(waveIndex int) {
	WaveCompletedEventType.Publish(evt.world, WaveCompletedEvent{WaveIndex: waveIndex})
}

// SubscribeToWaveCompleted subscribes to wave completed events
func (evt *EventSystem) SubscribeToWaveCompleted(callback events.Subscriber[WaveCompletedEvent]) {
	WaveCompletedEventType.Subscribe(evt.world, callback)
}

// EmitBossSpawnRequested emits a boss spawn requested event
func (evt *EventSystem) EmitBossSpawnRequested() {
	BossSpawnRequestedEventType.Publish(evt.world, BossSpawnRequestedEvent{})
}

// SubscribeToBossSpawnRequested subscribes to boss spawn requested events
func (evt *EventSystem) SubscribeToBossSpawnRequested(callback events.Subscriber[BossSpawnRequestedEvent]) {
	BossSpawnRequestedEventType.Subscribe(evt.world, callback)
}

// EmitBossSpawned emits a boss spawned event
func (evt *EventSystem) EmitBossSpawned() {
	BossSpawnedEventType.Publish(evt.world, BossSpawnedEvent{})
}

// SubscribeToBossSpawned subscribes to boss spawned events
func (evt *EventSystem) SubscribeToBossSpawned(callback events.Subscriber[BossSpawnedEvent]) {
	BossSpawnedEventType.Subscribe(evt.world, callback)
}

// EmitBossDefeated emits a boss defeated event
func (evt *EventSystem) EmitBossDefeated() {
	BossDefeatedEventType.Publish(evt.world, BossDefeatedEvent{})
}

// SubscribeToBossDefeated subscribes to boss defeated events
func (evt *EventSystem) SubscribeToBossDefeated(callback events.Subscriber[BossDefeatedEvent]) {
	BossDefeatedEventType.Subscribe(evt.world, callback)
}

// EmitStageCompleted emits a stage completed event
func (evt *EventSystem) EmitStageCompleted(stageNumber int) {
	StageCompletedEventType.Publish(evt.world, StageCompletedEvent{StageNumber: stageNumber})
}

// SubscribeToStageCompleted subscribes to stage completed events
func (evt *EventSystem) SubscribeToStageCompleted(callback events.Subscriber[StageCompletedEvent]) {
	StageCompletedEventType.Subscribe(evt.world, callback)
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
	LevelChangedEventType.ProcessEvents(evt.world)
	WaveStartedEventType.ProcessEvents(evt.world)
	WaveCompletedEventType.ProcessEvents(evt.world)
	BossSpawnRequestedEventType.ProcessEvents(evt.world)
	BossSpawnedEventType.ProcessEvents(evt.world)
	BossDefeatedEventType.ProcessEvents(evt.world)
	StageCompletedEventType.ProcessEvents(evt.world)
}
