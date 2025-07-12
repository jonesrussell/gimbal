package contracts

import (
	"context"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
)

// HealthSystem manages entity health, damage, and invincibility
type HealthSystem interface {
	DamageEntity(ctx context.Context, entity donburi.Entity, damage int) error
	HealEntity(ctx context.Context, entity donburi.Entity, amount int) error
	IsInvincible(ctx context.Context, entity donburi.Entity) bool
	GetHealth(ctx context.Context, entity donburi.Entity) (current, maxHealth int, ok bool)
	AddLife(ctx context.Context, entity donburi.Entity) error
	Update(ctx context.Context) error
}

// EventSystem manages game events and event handling
type EventSystem interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(ctx context.Context, eventType EventType, handler EventHandler) error
	EmitPlayerDamaged(ctx context.Context, entity donburi.Entity, damage, remaining int) error
	EmitGameOver(ctx context.Context) error
	EmitEnemyDestroyed(ctx context.Context, entity donburi.Entity, points int) error
	Update(ctx context.Context) error
}

// EnemySystem manages enemy spawning, movement, and destruction
type EnemySystem interface {
	DestroyEnemy(ctx context.Context, entity donburi.Entity) (points int, err error)
	SpawnEnemy(ctx context.Context, position common.Point) (donburi.Entity, error)
	GetActiveCount(ctx context.Context) int
	Update(ctx context.Context, deltaTime float64) error
}

// WeaponSystem manages weapon firing and projectile management
type WeaponSystem interface {
	FireWeapon(ctx context.Context, entity donburi.Entity, direction common.Point) error
	GetProjectileCount(ctx context.Context) int
	Update(ctx context.Context, deltaTime float64) error
}

// ResourceSystem manages game assets and resources
type ResourceSystem interface {
	LoadSprite(ctx context.Context, name, path string) (*ebiten.Image, error)
	LoadSound(ctx context.Context, name, path string) (Sound, error)
	GetSprite(ctx context.Context, name string) (*ebiten.Image, bool)
	GetDefaultFont(ctx context.Context) (Font, error)
	Cleanup(ctx context.Context) error
}

// ScoreSystem manages scoring and high scores
type ScoreSystem interface {
	AddScore(ctx context.Context, points int) error
	GetScore(ctx context.Context) int
	GetHighScore(ctx context.Context) int
	ResetScore(ctx context.Context) error
}

// StateSystem manages game state transitions
type StateSystem interface {
	SetPaused(ctx context.Context, paused bool) error
	IsPaused(ctx context.Context) bool
	SetGameOver(ctx context.Context, gameOver bool) error
	IsGameOver(ctx context.Context) bool
	GetGameState(ctx context.Context) GameState
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
