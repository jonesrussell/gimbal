package game

import (
	"context"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/debug"
	"github.com/jonesrussell/gimbal/internal/ecs/events"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/collision"
	enemysys "github.com/jonesrussell/gimbal/internal/ecs/systems/enemy"
	healthsys "github.com/jonesrussell/gimbal/internal/ecs/systems/health"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/movement"
	weaponsys "github.com/jonesrussell/gimbal/internal/ecs/systems/weapon"
	"github.com/jonesrussell/gimbal/internal/scenes"
)

// ECSGame represents the main game state using ECS
type ECSGame struct {
	world        donburi.World
	config       *config.GameConfig
	inputHandler common.GameInputHandler
	logger       common.Logger

	// Context for game lifecycle
	ctx    context.Context
	cancel context.CancelFunc

	// Event system
	eventSystem *events.EventSystem

	// Resource management
	resourceManager *resources.ResourceManager

	// Game state management
	stateManager *GameStateManager
	scoreManager *managers.ScoreManager
	levelManager *managers.LevelManager

	// Entity configurations (loaded from JSON)
	playerConfig *managers.PlayerConfig

	// Scene management
	sceneManager *scenes.SceneManager

	// Combat systems
	enemySystem       *enemysys.EnemySystem
	enemyWeaponSystem *enemysys.EnemyWeaponSystem
	weaponSystem      *weaponsys.WeaponSystem
	collisionSystem   *collision.CollisionSystem
	healthSystem      *healthsys.HealthSystem

	// Movement system
	movementSystem *movement.MovementSystem

	// 2025: EbitenUI responsive design system
	ui common.GameUI

	// Entity references
	playerEntity donburi.Entity
	starEntities []donburi.Entity
	frameCount   int // For debug logging

	// Performance optimization
	renderOptimizer *core.RenderOptimizer
	imagePool       *core.ImagePool
	perfMonitor     *debug.PerformanceMonitor

	// Debug system
	renderDebugger  *debug.RenderingDebugger
	showDebugInfo   bool
	debugKeyPressed bool
}

// Update updates the game state
func (g *ECSGame) Update() error {
	g.updatePerformanceMonitoring()
	g.updateDebugLogging()
	g.updateDebugInput()

	if err := g.updateCoreSystems(); err != nil {
		return err
	}

	// Use the game's context for proper lifecycle management
	ctx := g.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	if err := g.updateGameplaySystems(ctx); err != nil {
		return err
	}

	// Process queued events (score updates, damage events, etc.)
	g.eventSystem.ProcessEvents()

	g.updateHUD()
	g.endPerformanceMonitoring()

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

	// Render debug overlay if enabled
	if g.showDebugInfo && g.renderDebugger != nil {
		g.renderDebugger.StartFrame()
		g.renderDebugger.RenderDebugInfo(screen, g.world)
	}

	// Render wave debug info at top of screen when DEBUG is enabled
	if g.config.Debug && g.sceneManager.GetCurrentScene().GetType() == scenes.ScenePlaying {
		g.drawWaveDebugInfo(screen)
	}
}

// Layout returns the game's logical screen size
func (g *ECSGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.config.ScreenSize.Width, g.config.ScreenSize.Height
}

// Cleanup cleans up resources
func (g *ECSGame) Cleanup(ctx context.Context) {
	g.logger.Debug("Cleaning up ECS game")

	// Cancel the game context to signal shutdown to all systems
	if g.cancel != nil {
		g.cancel()
	}

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
