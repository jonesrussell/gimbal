package game

import (
	"context"
	"time"

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
	"github.com/jonesrussell/gimbal/internal/ui/state"
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

// updatePerformanceMonitoring handles performance monitoring for the frame
func (g *ECSGame) updatePerformanceMonitoring() {
	if g.perfMonitor != nil {
		g.perfMonitor.StartFrame()
	}
}

// updateDebugLogging handles periodic debug logging
func (g *ECSGame) updateDebugLogging() {
	g.frameCount++
	if g.frameCount%config.DebugLogInterval == 0 {
		g.logger.Debug("Game loop running",
			"frame", g.frameCount,
			"scene", g.sceneManager.GetCurrentScene(),
			"entities", g.world.Len(),
			"fps", ebiten.ActualFPS(),
			"player_valid", g.playerEntity != 0)
	}
}

// updateDebugInput handles debug key input
func (g *ECSGame) updateDebugInput() {
	if ebiten.IsKeyPressed(ebiten.KeyF3) && !g.debugKeyPressed {
		g.showDebugInfo = !g.showDebugInfo
		if g.renderDebugger != nil {
			g.renderDebugger.Toggle()
		}
		g.debugKeyPressed = true
		g.logger.Debug("Debug overlay toggled", "enabled", g.showDebugInfo)
	} else if !ebiten.IsKeyPressed(ebiten.KeyF3) {
		g.debugKeyPressed = false
	}
}

// updateCoreSystems updates scene manager and UI
func (g *ECSGame) updateCoreSystems() error {
	g.inputHandler.HandleInput()

	// Handle pause input
	g.handlePauseInput()

	if err := g.sceneManager.Update(); err != nil {
		g.logger.Error("Scene manager update failed", "error", err)
		return err
	}

	if err := g.ui.Update(); err != nil {
		g.logger.Error("UI update failed", "error", err)
		return err
	}

	return nil
}

// updateGameplaySystems updates ECS systems during gameplay
func (g *ECSGame) updateGameplaySystems(ctx context.Context) error {
	currentScene := g.sceneManager.GetCurrentScene()
	isPlayingScene := currentScene != nil && currentScene.GetType() == scenes.ScenePlaying
	if !isPlayingScene {
		return nil
	}

	deltaTime := config.DeltaTime

	// Handle shooting input
	g.handleShootingInput()

	systems := []struct {
		name     string
		updateFn func() error
	}{
		{"health", func() error { return g.healthSystem.Update(ctx) }},
		{"movement", func() error { return g.movementSystem.Update(ctx, deltaTime) }},
		{"collision", func() error { return g.collisionSystem.Update(ctx) }},
	}

	for _, system := range systems {
		if err := g.updateSystemWithTiming(system.name, system.updateFn); err != nil {
			return err
		}
	}

	// Update systems without error returns
	enemyUpdateFunc := func() error {
		return g.enemySystem.Update(ctx, deltaTime)
	}
	if err := g.updateSystemWithTiming("enemy", enemyUpdateFunc); err != nil {
		return err
	}
	enemyWeaponUpdateFunc := func() error {
		g.enemyWeaponSystem.Update(deltaTime)
		return nil
	}
	if err := g.updateSystemWithTiming("enemy_weapon", enemyWeaponUpdateFunc); err != nil {
		return err
	}
	weaponUpdateFunc := func() error {
		g.weaponSystem.Update(deltaTime)
		return nil
	}
	if err := g.updateSystemWithTiming("weapon", weaponUpdateFunc); err != nil {
		return err
	}

	// Check for level completion (boss killed)
	g.checkLevelCompletion()

	g.logger.Debug("ECS systems updated", "delta", deltaTime)
	return nil
}

// checkLevelCompletion checks if the boss is killed and advances the level
func (g *ECSGame) checkLevelCompletion() {
	// Check if boss was spawned but is now killed
	if g.enemySystem.WasBossSpawned() && !g.enemySystem.IsBossActive() {
		// Boss was killed, level complete!
		g.logger.Debug("Level complete - boss defeated")
		g.levelManager.IncrementLevel()

		// Reset enemy system for next level
		g.enemySystem.Reset()

		// TODO: Add level complete event/UI notification
	}
}

// handlePauseInput processes pause input and switches to pause scene
func (g *ECSGame) handlePauseInput() {
	currentScene := g.sceneManager.GetCurrentScene()

	// Only handle pause in playing scene
	if currentScene == nil || currentScene.GetType() != scenes.ScenePlaying {
		return
	}

	// Check if pause key is pressed
	if g.inputHandler.IsPausePressed() {
		g.logger.Debug("Pause key pressed, switching to pause scene")
		g.sceneManager.SwitchScene(scenes.ScenePaused)
	}
}

// handleShootingInput processes shooting input and fires weapons
func (g *ECSGame) handleShootingInput() {
	// Only handle shooting if we have a valid player entity
	if g.playerEntity == 0 {
		g.logger.Debug("No player entity found, skipping shooting input")
		return
	}

	// Check if shoot key is pressed
	if g.inputHandler.IsShootPressed() {
		// Get player position and angle
		playerEntry := g.world.Entry(g.playerEntity)
		if !playerEntry.Valid() {
			g.logger.Debug("Player entity invalid, skipping shooting input")
			return
		}

		pos := core.Position.Get(playerEntry)
		orbital := core.Orbital.Get(playerEntry)

		if pos == nil || orbital == nil {
			g.logger.Debug("Player position or orbital data missing, skipping shooting input")
			return
		}

		// Fire weapon with player position and facing angle
		if g.weaponSystem.FireWeapon(weaponsys.WeaponTypePrimary, *pos, orbital.FacingAngle) {
			g.logger.Debug("Weapon fired", "position", pos, "angle", orbital.FacingAngle)
		} else {
			g.logger.Debug("Weapon fire blocked by timing",
				"fire_timer", g.weaponSystem.GetFireTimer(),
				"fire_interval", g.weaponSystem.GetFireInterval())
		}
	}
}

// updateSystemWithTiming updates a system with performance timing
func (g *ECSGame) updateSystemWithTiming(systemName string, updateFn func() error) error {
	start := time.Now()
	err := updateFn()
	dur := time.Since(start)

	if err != nil {
		g.logger.Error("System update failed", "system", systemName, "error", err)
		return err
	}

	if dur > config.SlowSystemThreshold {
		g.logger.Warn("Slow system update", "system", systemName, "duration", dur)
	}

	return nil
}

// updateHUD updates the heads-up display
func (g *ECSGame) updateHUD() {
	current, maximum := g.healthSystem.GetPlayerHealth()
	healthPercent := 1.0
	if maximum > 0 {
		healthPercent = float64(current) / float64(maximum)
	}

	uiData := state.HUDData{
		Score:  g.scoreManager.GetScore(),
		Lives:  current,
		Level:  g.levelManager.GetLevel(),
		Health: healthPercent,
	}

	if hudUI, ok := g.ui.(interface{ UpdateHUD(state.HUDData) }); ok {
		hudUI.UpdateHUD(uiData)
	}
}

// endPerformanceMonitoring ends performance monitoring for the frame
func (g *ECSGame) endPerformanceMonitoring() {
	if g.perfMonitor != nil {
		g.perfMonitor.EndFrame()
	}
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
