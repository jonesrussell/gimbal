package health

import (
	"context"
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/contracts"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// HealthSystem manages player health, invincibility, and respawning with context support
type HealthSystem struct {
	world      donburi.World
	config     *config.GameConfig
	registry   contracts.SystemRegistry
	logger     common.Logger
	lastUpdate time.Time
}

// NewHealthSystem creates a new health system with proper dependency injection
func NewHealthSystem(
	world donburi.World,
	config *config.GameConfig,
	registry contracts.SystemRegistry,
	logger common.Logger,
) *HealthSystem {
	return &HealthSystem{
		world:      world,
		config:     config,
		registry:   registry,
		logger:     logger,
		lastUpdate: time.Now(),
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
	if err := hs.checkGameOverCondition(ctx); err != nil {
		return err
	}

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
	if err := hs.registry.Events().EmitPlayerDamaged(ctx, entity, damage, health.Current); err != nil {
		return err
	}

	hs.logger.Debug("Entity damaged", "damage", damage, "remaining_health", health.Current)

	// Check if entity should respawn or game over
	if health.Current > 0 {
		if err := hs.respawnEntity(ctx, entity); err != nil {
			return err
		}
	} else {
		if err := hs.registry.State().SetGameOver(ctx, true); err != nil {
			return err
		}
		if err := hs.registry.Events().EmitGameOver(ctx); err != nil {
			return err
		}
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
func (hs *HealthSystem) GetHealth(ctx context.Context, entity donburi.Entity) (current, max int, ok bool) {
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

// AddLife adds a life to an entity
func (hs *HealthSystem) AddLife(ctx context.Context, entity donburi.Entity) error {
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
	health.Current++
	if health.Current > health.Maximum {
		health.Current = health.Maximum
	}

	core.Health.SetValue(entry, *health)

	hs.logger.Debug("Life added to entity", "new_lives", health.Current)
	return nil
}

// checkGameOverCondition checks if the game should end
func (hs *HealthSystem) checkGameOverCondition(ctx context.Context) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Get player health
	players := make([]donburi.Entity, 0)
	query.NewQuery(
		filter.And(
			filter.Contains(core.PlayerTag),
			filter.Contains(core.Health),
		),
	).Each(hs.world, func(entry *donburi.Entry) {
		players = append(players, entry.Entity())
	})

	if len(players) == 0 {
		return nil
	}

	playerEntity := players[0]
	playerEntry := hs.world.Entry(playerEntity)
	if !playerEntry.Valid() {
		return nil
	}

	health := core.Health.Get(playerEntry)
	if health.Current <= 0 {
		if err := hs.registry.State().SetGameOver(ctx, true); err != nil {
			return err
		}
		if err := hs.registry.Events().EmitGameOver(ctx); err != nil {
			return err
		}
		hs.logger.Debug("Game over condition met")
	}

	return nil
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
