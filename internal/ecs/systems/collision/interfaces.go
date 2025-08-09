package collision

import (
	"context"

	"github.com/yohamta/donburi"
)

// HealthSystemInterface defines the interface for health system interactions during collision
type HealthSystemInterface interface {
	DamagePlayer(entity donburi.Entity, damage int)
	IsPlayerInvincible() bool
}

// EventSystemInterface defines the interface for event system interactions during collision
type EventSystemInterface interface {
	EmitPlayerDamaged(entity donburi.Entity, damage, remaining int)
	EmitGameOver()
	EmitEnemyDestroyed(ctx context.Context, entity donburi.Entity, points int) error
}

// EnemySystemInterface defines the interface for enemy system interactions during collision
type EnemySystemInterface interface {
	DestroyEnemy(entity donburi.Entity) int // returns points
}
