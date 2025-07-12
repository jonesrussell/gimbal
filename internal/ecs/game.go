package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	"github.com/jonesrussell/gimbal/internal/ecs/resources"
	scenes "github.com/jonesrussell/gimbal/internal/ecs/scenes"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/collision"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/health"
	"github.com/jonesrussell/gimbal/internal/ecs/ui_ebitenui"
)

// ECSGame represents the main game state using ECS
type ECSGame struct {
	world        donburi.World
	config       *common.GameConfig
	inputHandler common.GameInputHandler
	logger       common.Logger

	// Event system
	eventSystem *EventSystem

	// Resource management
	resourceManager *resources.ResourceManager

	// Game state management
	stateManager *GameStateManager
	scoreManager *managers.ScoreManager
	levelManager *LevelManager

	// Scene management
	sceneManager *scenes.SceneManager

	// Combat systems
	enemySystem     *EnemySystem
	weaponSystem    *WeaponSystem
	collisionSystem *collision.CollisionSystem
	healthSystem    *health.HealthSystem

	// 2025: EbitenUI responsive design system
	responsiveUI *ui_ebitenui.ResponsiveUI

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
	return g.updateCurrentScene()
}

// Draw renders the game
func (g *ECSGame) Draw(screen *ebiten.Image) {
	// Use scene manager to draw the current scene
	g.sceneManager.Draw(screen)

	// 2025: Render responsive HUD overlay
	if g.sceneManager.GetCurrentScene().GetType() == scenes.ScenePlaying {
		g.renderResponsiveHUD(screen)
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
func (g *ECSGame) GetLevelManager() *LevelManager {
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

// renderResponsiveHUD renders the 2025 responsive HUD overlay
func (g *ECSGame) renderResponsiveHUD(screen *ebiten.Image) {
	// Update responsive UI layout
	width, height := ebiten.WindowSize()
	g.responsiveUI.UpdateResponsiveLayout(width, height)

	// Draw the EbitenUI system
	g.responsiveUI.Draw(screen)
}
