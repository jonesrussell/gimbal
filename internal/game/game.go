package game

import (
	"context"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/events"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/collision"
	enemysys "github.com/jonesrussell/gimbal/internal/ecs/systems/enemy"
	healthsys "github.com/jonesrussell/gimbal/internal/ecs/systems/health"
	weaponsys "github.com/jonesrussell/gimbal/internal/ecs/systems/weapon"
	"github.com/jonesrussell/gimbal/internal/scenes"
	"github.com/jonesrussell/gimbal/internal/ui"
)

// ECSGame represents the main game state using ECS
type ECSGame struct {
	world        donburi.World
	config       *config.GameConfig
	inputHandler common.GameInputHandler
	logger       common.Logger

	// Event system
	eventSystem *events.EventSystem

	// Resource management
	resourceManager *resources.ResourceManager

	// Game state management
	stateManager *GameStateManager
	scoreManager *managers.ScoreManager
	levelManager *managers.LevelManager

	// Scene management
	sceneManager *scenes.SceneManager

	// Combat systems
	enemySystem     *enemysys.EnemySystem
	weaponSystem    *weaponsys.WeaponSystem
	collisionSystem *collision.CollisionSystem
	healthSystem    *healthsys.HealthSystem

	// 2025: EbitenUI responsive design system
	ui common.GameUI

	// Entity references
	playerEntity donburi.Entity
	starEntities []donburi.Entity
	frameCount   int // For debug logging
}

// Update updates the game state
func (g *ECSGame) Update() error {
	g.frameCount++
	if g.frameCount%60 == 0 {
		g.logger.Debug("Game loop running",
			"frame", g.frameCount,
			"scene", g.sceneManager.GetCurrentScene(),
			"entities", g.world.Len(),
			"fps", ebiten.ActualFPS(),
			"player_valid", g.playerEntity != 0)
	}
	// Handle input
	g.inputHandler.HandleInput()

	if err := g.sceneManager.Update(); err != nil {
		g.logger.Error("Scene manager update failed", "error", err)
	}

	if err := g.ui.Update(); err != nil {
		g.logger.Error("UI update failed", "error", err)
	}

	// NEW: Add ECS system updates during gameplay
	currentScene := g.sceneManager.GetCurrentScene()
	if currentScene != nil && currentScene.GetType() == scenes.ScenePlaying {
		deltaTime := 1.0 / 60.0 // Or use actual delta time
		ctx := context.Background()

		if g.healthSystem != nil {
			start := time.Now()
			if err := g.healthSystem.Update(ctx); err != nil {
				g.logger.Error("Health system update failed", "error", err)
			}
			dur := time.Since(start)
			if dur > 5*time.Millisecond {
				g.logger.Warn("Slow system update", "system", "health", "duration", dur)
			}
		}
		if g.enemySystem != nil {
			start := time.Now()
			g.enemySystem.Update(deltaTime)
			dur := time.Since(start)
			if dur > 5*time.Millisecond {
				g.logger.Warn("Slow system update", "system", "enemy", "duration", dur)
			}
		}
		if g.weaponSystem != nil {
			start := time.Now()
			g.weaponSystem.Update(deltaTime)
			dur := time.Since(start)
			if dur > 5*time.Millisecond {
				g.logger.Warn("Slow system update", "system", "weapon", "duration", dur)
			}
		}
		if g.collisionSystem != nil {
			start := time.Now()
			if err := g.collisionSystem.Update(ctx); err != nil {
				g.logger.Error("Collision system update failed", "error", err)
			}
			dur := time.Since(start)
			if dur > 5*time.Millisecond {
				g.logger.Warn("Slow system update", "system", "collision", "duration", dur)
			}
		}
		// Add other ECS systems here if needed
		g.logger.Debug("ECS systems updated", "delta", deltaTime)
	}

	current, maximum := g.healthSystem.GetPlayerHealth()
	healthPercent := 1.0
	if maximum > 0 {
		healthPercent = float64(current) / float64(maximum)
	}
	uiData := ui.HUDData{
		Score:  g.scoreManager.GetScore(),
		Lives:  current,
		Level:  g.levelManager.GetLevel(),
		Health: healthPercent,
	}
	if hudUI, ok := g.ui.(interface{ UpdateHUD(ui.HUDData) }); ok {
		hudUI.UpdateHUD(uiData)
	}

	return nil
}

// Draw renders the game
func (g *ECSGame) Draw(screen *ebiten.Image) {
	// Use scene manager to draw the current scene
	g.sceneManager.Draw(screen)

	// 2025: Render responsive HUD overlay
	if g.sceneManager.GetCurrentScene().GetType() == scenes.ScenePlaying {
		g.ui.Draw(screen)
	}
}

// Layout returns the game's logical screen size
func (g *ECSGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.config.ScreenSize.Width, g.config.ScreenSize.Height
}

// Cleanup cleans up resources
func (g *ECSGame) Cleanup(ctx context.Context) {
	g.logger.Debug("Cleaning up ECS game")

	// Clean up resources
	if g.resourceManager != nil {
		if err := g.resourceManager.Cleanup(ctx); err != nil {
			g.logger.Error("Failed to cleanup resource manager", "error", err)
		}
	}

	// Donburi handles entity cleanup automatically
}

// IsPaused returns the pause state
func (g *ECSGame) IsPaused() bool {
	return g.stateManager.IsPaused()
}

// GetScoreManager returns the score manager
func (g *ECSGame) GetScoreManager() *managers.ScoreManager {
	return g.scoreManager
}

// GetLevelManager returns the level manager
func (g *ECSGame) GetLevelManager() *managers.LevelManager {
	return g.levelManager
}

// SetInputHandler sets the input handler (for testing)
func (g *ECSGame) SetInputHandler(handler common.GameInputHandler) {
	g.inputHandler = handler
}

// GetInputHandler returns the current input handler
func (g *ECSGame) GetInputHandler() common.GameInputHandler {
	return g.inputHandler
}
