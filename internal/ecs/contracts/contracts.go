package contracts

import (
	"context"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
)

// System represents a generic ECS system
type System[T any] interface {
	Update(ctx context.Context, data T) error
	Initialize(ctx context.Context) error
	Cleanup(ctx context.Context) error
}

// HealthSystem manages entity health, damage, and invincibility
type HealthSystem interface {
	System[HealthUpdateData]
	DamageEntity(ctx context.Context, entity donburi.Entity, damage int) common.Result[HealthResult]
	HealEntity(ctx context.Context, entity donburi.Entity, amount int) common.Result[HealthResult]
	IsInvincible(ctx context.Context, entity donburi.Entity) bool
	GetHealth(ctx context.Context, entity donburi.Entity) common.Result[HealthResult]
	AddLife(ctx context.Context, entity donburi.Entity) common.Result[HealthResult]
}

// HealthUpdateData contains data needed for health system updates
type HealthUpdateData struct {
	DeltaTime float64
}

// HealthResult contains health-related operation results
type HealthResult struct {
	CurrentHealth int
	MaxHealth     int
	IsInvincible  bool
}

// EventSystem manages game events and event handling
type EventSystem interface {
	System[EventUpdateData]
	Publish(ctx context.Context, event Event) error
	Subscribe(ctx context.Context, eventType EventType, handler EventHandler) func() // Returns unsubscribe function
	EmitPlayerDamaged(ctx context.Context, entity donburi.Entity, damage, remaining int) error
	EmitGameOver(ctx context.Context) error
	EmitEnemyDestroyed(ctx context.Context, entity donburi.Entity, points int) error
}

// EventUpdateData contains data needed for event system updates
type EventUpdateData struct {
	DeltaTime float64
}

// EnemySystem manages enemy spawning, movement, and destruction
type EnemySystem interface {
	System[EnemyUpdateData]
	DestroyEnemy(ctx context.Context, entity donburi.Entity) common.Result[EnemyResult]
	SpawnEnemy(ctx context.Context, position common.Point) common.Result[donburi.Entity]
	GetActiveCount(ctx context.Context) int
}

// EnemyUpdateData contains data needed for enemy system updates
type EnemyUpdateData struct {
	DeltaTime float64
	PlayerPos common.Point
}

// EnemyResult contains enemy-related operation results
type EnemyResult struct {
	Points int
	Entity donburi.Entity
}

// WeaponSystem manages weapon firing and projectile management
type WeaponSystem interface {
	System[WeaponUpdateData]
	FireWeapon(ctx context.Context, entity donburi.Entity, direction common.Point) common.Result[WeaponResult]
	GetProjectileCount(ctx context.Context) int
}

// WeaponUpdateData contains data needed for weapon system updates
type WeaponUpdateData struct {
	DeltaTime float64
	PlayerPos common.Point
}

// WeaponResult contains weapon-related operation results
type WeaponResult struct {
	ProjectileEntity donburi.Entity
	Success          bool
}

// ResourceSystem manages game assets and resources
type ResourceSystem interface {
	System[ResourceUpdateData]
	LoadSprite(ctx context.Context, name, path string) common.Result[*ebiten.Image]
	LoadSound(ctx context.Context, name, path string) common.Result[Sound]
	GetSprite(ctx context.Context, name string) common.Result[*ebiten.Image]
	GetDefaultFont(ctx context.Context) common.Result[Font]
}

// ResourceUpdateData contains data needed for resource system updates
type ResourceUpdateData struct {
	DeltaTime float64
}

// ScoreSystem manages scoring and high scores
type ScoreSystem interface {
	System[ScoreUpdateData]
	AddScore(ctx context.Context, points int) common.Result[ScoreResult]
	GetScore(ctx context.Context) int
	GetHighScore(ctx context.Context) int
	ResetScore(ctx context.Context) error
}

// ScoreUpdateData contains data needed for score system updates
type ScoreUpdateData struct {
	DeltaTime float64
}

// ScoreResult contains score-related operation results
type ScoreResult struct {
	CurrentScore int
	HighScore    int
	PointsAdded  int
}

// StateSystem manages game state transitions
type StateSystem interface {
	System[StateUpdateData]
	SetPaused(ctx context.Context, paused bool) error
	IsPaused(ctx context.Context) bool
	SetGameOver(ctx context.Context, gameOver bool) error
	IsGameOver(ctx context.Context) bool
	GetGameState(ctx context.Context) common.Result[GameState]
}

// StateUpdateData contains data needed for state system updates
type StateUpdateData struct {
	DeltaTime float64
}

// Event represents a game event
type Event interface {
	Type() EventType
	Timestamp() int64
	Data() map[string]interface{}
}

// EventType represents the type of event
type EventType int

const (
	EventTypePlayerDamaged EventType = iota
	EventTypeGameOver
	EventTypeEnemyDestroyed
	EventTypeScoreChanged
	EventTypeLifeAdded
)

// EventHandler processes game events
type EventHandler func(ctx context.Context, event Event) error

// Sound represents an audio resource
type Sound interface {
	Play(ctx context.Context) error
	Stop(ctx context.Context) error
	SetVolume(ctx context.Context, volume float64) error
}

// Font represents a text font resource
type Font interface {
	MeasureText(text string) (width, height int)
	DrawText(ctx context.Context, screen *ebiten.Image, text string, x, y int, color Color) error
}

// GameState represents the current game state
type GameState struct {
	Paused   bool
	GameOver bool
	Score    int
	Level    int
}

// Color represents a color value
type Color struct {
	R, G, B, A uint8
}
