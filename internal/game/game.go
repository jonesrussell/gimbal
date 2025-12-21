package game

import (
	"context"
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	v2text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

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

// convertWaveConfigs converts managers.WaveConfig to enemy.WaveConfig
func convertWaveConfigs(managerWaves []managers.WaveConfig) []enemysys.WaveConfig {
	enemyWaves := make([]enemysys.WaveConfig, len(managerWaves))
	for i, mw := range managerWaves {
		enemyTypes := make([]enemysys.EnemyType, len(mw.EnemyTypes))
		for j, et := range mw.EnemyTypes {
			enemyTypes[j] = enemysys.EnemyType(et)
		}
		enemyWaves[i] = enemysys.WaveConfig{
			FormationType:   enemysys.FormationType(mw.FormationType),
			EnemyCount:      mw.EnemyCount,
			EnemyTypes:      enemyTypes,
			SpawnDelay:      mw.SpawnDelay,
			Timeout:         mw.Timeout,
			InterWaveDelay:  mw.InterWaveDelay,
			MovementPattern: enemysys.MovementPattern(mw.MovementPattern),
		}
	}
	return enemyWaves
}

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
	// Get current level config to check completion conditions
	levelConfig := g.levelManager.GetCurrentLevelConfig()
	if levelConfig == nil {
		return
	}

	// Check completion conditions
	conditions := levelConfig.CompletionConditions
	canComplete := true

	// Check if boss kill is required
	if conditions.RequireBossKill {
		if !g.enemySystem.WasBossSpawned() || g.enemySystem.IsBossActive() {
			canComplete = false
		}
	}

	// Check if all waves are required
	if conditions.RequireAllWaves {
		if g.enemySystem.GetWaveManager().HasMoreWaves() {
			canComplete = false
		}
	}

	// Check if all enemies must be killed
	if conditions.RequireAllEnemiesKilled {
		// This would require checking active enemy count
		// For now, we'll assume boss kill + all waves = all enemies killed
		// TODO: Implement active enemy count check
		_ = conditions // Avoid unused variable warning
	}

	if canComplete {
		// Level complete!
		currentLevel := g.levelManager.GetLevel()
		g.logger.Debug("Level complete", "level", currentLevel)
		g.levelManager.IncrementLevel()

		// Load next level's configuration
		nextLevelConfig := g.levelManager.GetCurrentLevelConfig()
		if nextLevelConfig != nil {
			enemyWaves := convertWaveConfigs(nextLevelConfig.Waves)
			g.enemySystem.LoadLevelConfig(enemyWaves, &nextLevelConfig.Boss)
			g.logger.Debug("Next level config loaded",
				"level", nextLevelConfig.LevelNumber,
				"waves", len(nextLevelConfig.Waves))

			// Show level title for new level
			if currentScene := g.sceneManager.GetCurrentScene(); currentScene != nil {
				if playingScene, ok := currentScene.(interface{ ShowLevelTitle(int) }); ok {
					playingScene.ShowLevelTitle(nextLevelConfig.LevelNumber)
				}
			}
		} else {
			// No more levels, just reset
			g.enemySystem.Reset()
		}

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

// drawWaveDebugInfo draws wave information at the bottom of the screen
func (g *ECSGame) drawWaveDebugInfo(screen *ebiten.Image) {
	if g.enemySystem == nil {
		return
	}

	waveManager := g.enemySystem.GetWaveManager()
	if waveManager == nil {
		return
	}

	screenHeight := float64(g.config.ScreenSize.Height)
	lineHeight := 20.0
	x := 10.0

	currentWave := waveManager.GetCurrentWave()
	if currentWave == nil {
		// No active wave - check if boss is active or spawning
		if !waveManager.HasMoreWaves() {
			// All waves complete - show boss info if boss is active or spawning
			if g.enemySystem.IsBossActive() || g.enemySystem.WasBossSpawned() {
				g.drawBossDebugInfo(screen, x, screenHeight, lineHeight)
				return
			}
			// Boss not spawned yet - show spawn timer
			if g.enemySystem.WasBossSpawned() {
				// Boss was spawned but is dead
				g.drawDebugText(screen, "Boss: Defeated", x, screenHeight-lineHeight)
				return
			}
			// Boss spawning soon
			g.drawDebugText(screen, "Boss: Spawning soon...", x, screenHeight-lineHeight)
			return
		}

		// Still have waves - show waiting status
		var statusText string
		if waveManager.IsWaiting() {
			statusText = "Wave: Waiting for next wave..."
		} else {
			statusText = "Wave: Starting..."
		}
		g.drawDebugText(screen, statusText, x, screenHeight-lineHeight)
		return
	}

	// Format formation type
	formationName := g.formatFormationType(currentWave.Config.FormationType)

	// Format enemy types
	enemyTypesStr := g.formatEnemyTypes(currentWave.Config.EnemyTypes)

	// Calculate progress
	progress := float64(currentWave.EnemiesKilled) / float64(currentWave.Config.EnemyCount) * 100
	if currentWave.Config.EnemyCount == 0 {
		progress = 0
	}

	// Calculate number of lines to determine starting Y position
	// Wave, Formation, Enemies, Spawned, Types, Pattern, Status, Timer
	numLines := 8
	startY := screenHeight - float64(numLines)*lineHeight - 20 // Increased margin to prevent cutoff

	// Draw wave information from bottom up
	y := startY
	g.drawDebugText(screen, fmt.Sprintf("Wave %d/%d", currentWave.WaveIndex+1, waveManager.GetWaveCount()), x, y)
	y += lineHeight
	g.drawDebugText(screen, fmt.Sprintf("Formation: %s", formationName), x, y)
	y += lineHeight
	enemyText := fmt.Sprintf("Enemies: %d/%d (%.0f%%)",
		currentWave.EnemiesKilled, currentWave.Config.EnemyCount, progress)
	g.drawDebugText(screen, enemyText, x, y)
	y += lineHeight
	g.drawDebugText(screen, fmt.Sprintf("Spawned: %d", currentWave.EnemiesSpawned), x, y)
	y += lineHeight
	g.drawDebugText(screen, fmt.Sprintf("Types: %s", enemyTypesStr), x, y)
	y += lineHeight
	patternText := fmt.Sprintf("Pattern: %s",
		g.formatMovementPattern(currentWave.Config.MovementPattern))
	g.drawDebugText(screen, patternText, x, y)
	y += lineHeight
	if currentWave.IsSpawning {
		g.drawDebugText(screen, "Status: Spawning", x, y)
	} else if currentWave.IsComplete {
		g.drawDebugText(screen, "Status: Complete", x, y)
	} else {
		g.drawDebugText(screen, "Status: Active", x, y)
	}
	y += lineHeight
	g.drawDebugText(screen, fmt.Sprintf("Timer: %.1fs", currentWave.WaveTimer), x, y)
}

// drawDebugText draws text with a semi-transparent background
func (g *ECSGame) drawDebugText(screen *ebiten.Image, text string, x, y float64) {
	ctx := g.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	fontFace, err := g.resourceManager.GetDefaultFont(ctx)
	if err != nil {
		return
	}

	// Measure text size
	width, height := v2text.Measure(text, fontFace, 0)

	// Draw semi-transparent black rectangle behind text
	padding := float32(4.0)
	vector.DrawFilledRect(screen,
		float32(x)-padding,
		float32(y)-float32(height)-padding,
		float32(width)+padding*2,
		float32(height)+padding*2,
		color.RGBA{0, 0, 0, 150}, false)

	// Draw text on top
	op := &v2text.DrawOptions{}
	op.GeoM.Translate(x, y)
	v2text.Draw(screen, text, fontFace, op)
}

// formatFormationType formats a formation type as a string
func (g *ECSGame) formatFormationType(ft enemysys.FormationType) string {
	switch ft {
	case enemysys.FormationLine:
		return "Line"
	case enemysys.FormationCircle:
		return "Circle"
	case enemysys.FormationV:
		return "V"
	case enemysys.FormationDiamond:
		return "Diamond"
	case enemysys.FormationDiagonal:
		return "Diagonal"
	case enemysys.FormationSpiral:
		return "Spiral"
	case enemysys.FormationRandom:
		return "Random"
	default:
		return "Unknown"
	}
}

// formatEnemyTypes formats enemy types as a string
func (g *ECSGame) formatEnemyTypes(types []enemysys.EnemyType) string {
	if len(types) == 0 {
		return "None"
	}

	typeCounts := make(map[string]int)
	for _, t := range types {
		typeCounts[t.String()]++
	}

	result := ""
	first := true
	for name, count := range typeCounts {
		if !first {
			result += ", "
		}
		if count > 1 {
			result += fmt.Sprintf("%s x%d", name, count)
		} else {
			result += name
		}
		first = false
	}
	return result
}

// formatMovementPattern formats a movement pattern as a string
func (g *ECSGame) formatMovementPattern(mp enemysys.MovementPattern) string {
	switch mp {
	case enemysys.MovementPatternNormal:
		return "Normal"
	case enemysys.MovementPatternZigzag:
		return "Zigzag"
	case enemysys.MovementPatternAccelerating:
		return "Accelerating"
	case enemysys.MovementPatternPulsing:
		return "Pulsing"
	default:
		return "Unknown"
	}
}

// drawBossDebugInfo draws boss debug information
func (g *ECSGame) drawBossDebugInfo(screen *ebiten.Image, x, screenHeight, lineHeight float64) {
	bossEntry := g.findBossEntity()
	if bossEntry == nil {
		g.drawBossStatus(screen, x, screenHeight, lineHeight)
		return
	}
	g.drawBossDetails(screen, bossEntry, x, screenHeight, lineHeight)
}

// findBossEntity finds the boss entity in the world
func (g *ECSGame) findBossEntity() *donburi.Entry {
	var bossEntry *donburi.Entry
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
			filter.Contains(core.EnemyTypeID),
		),
	).Each(g.world, func(entry *donburi.Entry) {
		typeID := core.EnemyTypeID.Get(entry)
		if enemysys.EnemyType(*typeID) == enemysys.EnemyTypeBoss {
			bossEntry = entry
		}
	})
	return bossEntry
}

// drawBossStatus draws boss spawn/defeat status
func (g *ECSGame) drawBossStatus(screen *ebiten.Image, x, screenHeight, lineHeight float64) {
	if g.enemySystem.WasBossSpawned() {
		g.drawDebugText(screen, "Boss: Defeated", x, screenHeight-lineHeight)
	} else {
		g.drawDebugText(screen, "Boss: Spawning soon...", x, screenHeight-lineHeight)
	}
}

// drawBossDetails draws detailed boss information
func (g *ECSGame) drawBossDetails(screen *ebiten.Image, bossEntry *donburi.Entry, x, screenHeight, lineHeight float64) {
	pos := core.Position.Get(bossEntry)
	health := core.Health.Get(bossEntry)
	orbital := core.Orbital.Get(bossEntry)
	size := core.Size.Get(bossEntry)

	// Calculate number of lines for boss info
	numLines := 6 // Boss, Health, Position, Orbital Angle, Size, Status
	startY := screenHeight - float64(numLines)*lineHeight - 20

	// Draw boss information from bottom up
	y := startY
	g.drawDebugText(screen, "BOSS", x, y)
	y += lineHeight

	if health != nil {
		healthPercent := float64(health.Current) / float64(health.Maximum) * 100
		healthText := fmt.Sprintf("Health: %d/%d (%.0f%%)",
			health.Current, health.Maximum, healthPercent)
		g.drawDebugText(screen, healthText, x, y)
	} else {
		g.drawDebugText(screen, "Health: Unknown", x, y)
	}
	y += lineHeight

	if pos != nil {
		g.drawDebugText(screen, fmt.Sprintf("Position: (%.0f, %.0f)", pos.X, pos.Y), x, y)
	} else {
		g.drawDebugText(screen, "Position: Unknown", x, y)
	}
	y += lineHeight

	if orbital != nil {
		g.drawDebugText(screen, fmt.Sprintf("Orbital Angle: %.1f°", float64(orbital.OrbitalAngle)), x, y)
	} else {
		g.drawDebugText(screen, "Orbital: Unknown", x, y)
	}
	y += lineHeight

	if size != nil {
		g.drawDebugText(screen, fmt.Sprintf("Size: %dx%d", size.Width, size.Height), x, y)
	} else {
		g.drawDebugText(screen, "Size: Unknown", x, y)
	}
	y += lineHeight

	status := "Active"
	if health != nil && health.Current <= 0 {
		status = "Defeated"
	}
	g.drawDebugText(screen, fmt.Sprintf("Status: %s", status), x, y)
}
