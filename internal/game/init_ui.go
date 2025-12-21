package game

import (
	"context"
	"fmt"

	"github.com/jonesrussell/gimbal/internal/errors"
	uicore "github.com/jonesrussell/gimbal/internal/ui/core"
	"github.com/jonesrussell/gimbal/internal/ui/responsive"
)

// initializeUI sets up the game UI
func (g *ECSGame) initializeUI(ctx context.Context) error {
	font, err := g.resourceManager.GetDefaultFont(ctx)
	if err != nil {
		return errors.NewGameErrorWithCause(errors.AssetLoadFailed, "failed to get default font", err)
	}

	heartSprite, err := g.resourceManager.GetUISprite(ctx, "heart", uicore.HeartIconSize)
	if err != nil {
		return fmt.Errorf("failed to load heart sprite: %w", err)
	}

	ammoSprite, err := g.resourceManager.GetUISprite(ctx, "ammo", uicore.AmmoIconSize)
	if err != nil {
		g.logger.Warn("Failed to load ammo sprite, using fallback", "error", err)
		ammoSprite = nil // Will use fallback in UI
	}

	uiConfig := &responsive.Config{
		Font:        font,
		HeartSprite: heartSprite,
		AmmoSprite:  ammoSprite,
	}

	gameUI, err := responsive.NewResponsiveUI(uiConfig)
	if err != nil {
		return fmt.Errorf("failed to create game UI: %w", err)
	}
	g.ui = gameUI
	return nil
}

