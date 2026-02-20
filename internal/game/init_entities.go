package game

import (
	"context"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/errors"
)

// createGameEntities creates all game entities
func (g *ECSGame) createGameEntities(ctx context.Context) error {
	return g.createEntities(ctx)
}

// createEntities creates all game entities
func (g *ECSGame) createEntities(ctx context.Context) error {
	// Get sprites from resource manager
	playerSprite, ok := g.resourceManager.GetSprite(ctx, resources.SpritePlayer)
	if !ok {
		return errors.NewGameError(errors.AssetNotFound, "player sprite not found")
	}

	starSprite, ok := g.resourceManager.GetSprite(ctx, resources.SpriteStar)
	if !ok {
		return errors.NewGameError(errors.AssetNotFound, "star sprite not found")
	}

	// Create player with config from JSON
	if g.playerConfig == nil {
		return errors.NewGameError(errors.ConfigMissing, "player config not loaded")
	}
	g.playerEntity = core.CreatePlayer(g.world, playerSprite, g.config, g.playerConfig)

	// Create star field
	g.starEntities = core.CreateStarField(g.world, starSprite, g.config)

	return nil
}
