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
	"github.com/jonesrussell/gimbal/internal/ecs/viewport"
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

	// 2025: Advanced responsive design system
	viewport           *viewport.AdvancedViewportManager
	fluidGrid          *viewport.FluidGrid
	responsiveHUD      *viewport.GameHUD
	responsiveRenderer *viewport.ResponsiveRenderer
	accessibility      *viewport.AccessibilityConfig

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
	// 2025: Update advanced viewport with responsive techniques
	g.viewport.UpdateAdvanced(outsideWidth, outsideHeight)

	// 2025: Adaptive logical screen size based on device class
	switch g.viewport.GetDeviceClass() {
	case string(viewport.DeviceClassMobile):
		if g.viewport.GetOrientation() == string(viewport.OrientationPortrait) {
			return 1080, 1920 // Mobile portrait
		}
		return 1920, 1080 // Mobile landscape
	case string(viewport.DeviceClassTablet):
		return 1440, 1080 // Tablet optimized
	case string(viewport.DeviceClassUltrawide):
		return 2560, 1080 // Ultrawide support
	default:
		return 1920, 1080 // Standard desktop
	}
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
	// Update fluid grid container dimensions
	width, height := g.viewport.GetCurrentDimensions()
	g.fluidGrid.UpdateContainer(float64(width), float64(height))

	// Render the responsive HUD using the performance-optimized renderer
	g.responsiveRenderer.RenderFrame(screen, g.viewport, g.responsiveHUD)
}
