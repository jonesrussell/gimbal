package game

import (
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
}

// Update updates the game state
func (g *ECSGame) Update() error {
	// Handle input
	g.inputHandler.HandleInput()

	// Update scene manager
	if err := g.sceneManager.Update(); err != nil {
		g.logger.Error("Scene update failed", "error", err)
		return err
	}

	// Update based on current scene
	if err := g.ui.Update(); err != nil {
		return err
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

// Layout implements ebiten.Game interface
func (g *ECSGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// 2025: Responsive layout based on screen size
	aspectRatio := float64(outsideWidth) / float64(outsideHeight)

	// Mobile portrait
	if outsideWidth < 768 && aspectRatio < 1.0 {
		return 1080, 1920
	}

	// Mobile landscape / tablet
	if outsideWidth < 1024 {
		return 1440, 1080
	}

	// Desktop standard
	if outsideWidth < 1920 {
		return 1920, 1080
	}

	// Ultrawide support
	return outsideWidth, 1080
}

// Cleanup cleans up resources
func (g *ECSGame) Cleanup() {
	g.logger.Debug("Cleaning up ECS game")

	// Clean up resources
	if g.resourceManager != nil {
		g.resourceManager.Cleanup()
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
