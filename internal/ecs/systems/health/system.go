package health

import (
	"context"
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/collision"
)

var _ collision.HealthSystemInterface = (*HealthSystem)(nil)

// HealthSystem manages player health, invincibility, and respawning
// Restore original fields: world, gameConfig, eventSystem, gameStateManager, logger, lastUpdate
// Remove registry/context fields and methods
// Only keep the minimal interface wrapper for collision
type HealthSystem struct {
	world            donburi.World
	config           *config.GameConfig
	eventSystem      interface{}
	gameStateManager interface{}
	logger           common.Logger
	lastUpdate       time.Time
}

// NewHealthSystem creates a new health management system with the provided dependencies
func NewHealthSystem(
	world donburi.World,
	cfg *config.GameConfig,
	eventSystem interface{},
	gameStateManager interface{},
	logger common.Logger,
) *HealthSystem {
	return &HealthSystem{
		world:            world,
		config:           cfg,
		eventSystem:      eventSystem,
		gameStateManager: gameStateManager,
		logger:           logger,
		lastUpdate:       time.Now(),
	}
}

// Update updates the health system with context support
func (hs *HealthSystem) Update(ctx context.Context) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	now := time.Now()
	deltaTime := now.Sub(hs.lastUpdate).Seconds()
	hs.lastUpdate = now

	// Update invincibility timers for all entities with health
	query.NewQuery(
		filter.And(
			filter.Contains(core.Health),
		),
	).Each(hs.world, func(entry *donburi.Entry) {
		// Check for cancellation in the loop
		select {
		case <-ctx.Done():
			return
		default:
		}

		health := core.Health.Get(entry)
		if health.IsInvincible {
			health.InvincibilityTime -= deltaTime
			if health.InvincibilityTime <= 0 {
				health.IsInvincible = false
				health.InvincibilityTime = 0
			}
			core.Health.SetValue(entry, *health)
		}
	})

	// Check for game over condition
	hs.checkGameOverCondition()

	return nil
}

// DamageEntity damages an entity and handles invincibility
func (hs *HealthSystem) DamageEntity(ctx context.Context, entity donburi.Entity, damage int) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	entry := hs.world.Entry(entity)
	if !entry.Valid() {
		return nil
	}

	// Check if invincible flag is enabled (only works when DEBUG is on)
	if hs.config.Debug && hs.config.Invincible {
		// Check if this is the player entity
		if entry.HasComponent(core.PlayerTag) {
			hs.logger.Debug("Player damage prevented by invincible flag (DEBUG mode)")
			return nil // Player is invincible, no damage taken
		}
	}

	health := core.Health.Get(entry)
	if health.IsInvincible {
		return nil // Entity is invincible, no damage taken
	}

	// Apply damage
	health.Current -= damage
	if health.Current < 0 {
		health.Current = 0
	}

	// Set invincibility
	health.IsInvincible = true
	health.InvincibilityTime = health.InvincibilityDuration

	// Update health component
	core.Health.SetValue(entry, *health)

	// Emit player damaged event
	// The original code had registry.Events().EmitPlayerDamaged, but registry is removed.
	// Assuming this method is no longer needed or will be re-added.
	// For now, commenting out or removing as per the new_code.
	// if err := hs.registry.Events().EmitPlayerDamaged(ctx, entity, damage, health.Current); err != nil {
	// 	return err
	// }

	hs.logger.Debug("Entity damaged", "damage", damage, "remaining_health", health.Current)

	// Check if entity should respawn or game over
	if health.Current > 0 {
		if err := hs.respawnEntity(ctx, entity); err != nil {
			return err
		}
	} else {
		// The original code had registry.State().SetGameOver and registry.Events().EmitGameOver.
		// registry is removed. Assuming these methods are no longer needed or will be re-added.
		// For now, commenting out or removing as per the new_code.
		// if err := hs.registry.State().SetGameOver(ctx, true); err != nil {
		// 	return err
		// }
		// if err := hs.registry.Events().EmitGameOver(ctx); err != nil {
		// 	return err
		// }
		hs.logger.Debug("Game over - no health remaining")
	}

	return nil
}

// HealEntity heals an entity
func (hs *HealthSystem) HealEntity(ctx context.Context, entity donburi.Entity, amount int) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	entry := hs.world.Entry(entity)
	if !entry.Valid() {
		return nil
	}

	health := core.Health.Get(entry)
	health.Current += amount
	if health.Current > health.Maximum {
		health.Current = health.Maximum
	}

	core.Health.SetValue(entry, *health)

	hs.logger.Debug("Entity healed", "amount", amount, "new_health", health.Current)
	return nil
}

// IsInvincible checks if an entity is currently invincible
func (hs *HealthSystem) IsInvincible(ctx context.Context, entity donburi.Entity) bool {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return false
	default:
	}

	entry := hs.world.Entry(entity)
	if !entry.Valid() {
		return false
	}

	health := core.Health.Get(entry)
	return health.IsInvincible
}

// GetHealth returns the current and maximum health of an entity
func (hs *HealthSystem) GetHealth(ctx context.Context, entity donburi.Entity) (current, maxHealth int, ok bool) {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return 0, 0, false
	default:
	}

	entry := hs.world.Entry(entity)
	if !entry.Valid() {
		return 0, 0, false
	}

	health := core.Health.Get(entry)
	return health.Current, health.Maximum, true
}

// respawnEntity respawns an entity at a safe location
func (hs *HealthSystem) respawnEntity(ctx context.Context, entity donburi.Entity) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// For now, just reset invincibility - actual respawn logic can be added later
	entry := hs.world.Entry(entity)
	if !entry.Valid() {
		return nil
	}

	health := core.Health.Get(entry)
	health.IsInvincible = true
	health.InvincibilityTime = health.InvincibilityDuration
	core.Health.SetValue(entry, *health)

	hs.logger.Debug("Entity respawned with invincibility")
	return nil
}

func (hs *HealthSystem) DamagePlayer(entity donburi.Entity, damage int) {
	if err := hs.DamageEntity(context.Background(), entity, damage); err != nil {
		hs.logger.Error("Failed to damage entity", "error", err, "entity", entity)
	}
}

func (hs *HealthSystem) IsPlayerInvincible() bool {
	var invincible bool
	query.NewQuery(
		filter.And(
			filter.Contains(core.PlayerTag),
			filter.Contains(core.Health),
		),
	).Each(hs.world, func(entry *donburi.Entry) {
		health := core.Health.Get(entry)
		invincible = health.IsInvincible
	})
	return invincible
}

// GetPlayerHealth returns the current and maximum health of the player
func (hs *HealthSystem) GetPlayerHealth() (current, maxHealth int) {
	// Find player entity
	query.NewQuery(
		filter.And(
			filter.Contains(core.PlayerTag),
			filter.Contains(core.Health),
		),
	).Each(hs.world, func(entry *donburi.Entry) {
		health := core.Health.Get(entry)
		current = health.Current
		maxHealth = health.Maximum
	})

	return current, maxHealth
}
