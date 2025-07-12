package collision

import (
	"context"

	"github.com/yohamta/donburi"
)

type HealthSystemInterface interface {
	DamagePlayer(entity donburi.Entity, damage int)
	IsPlayerInvincible() bool
}

type EventSystemInterface interface {
	EmitPlayerDamaged(entity donburi.Entity, damage, remaining int)
	EmitGameOver()
	EmitEnemyDestroyed(ctx context.Context, entity donburi.Entity, points int) error
}

type EnemySystemInterface interface {
	DestroyEnemy(entity donburi.Entity) int // returns points
}
