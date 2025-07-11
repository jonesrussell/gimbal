package ecs

import (
	"github.com/yohamta/donburi/ecs"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	scenes "github.com/jonesrussell/gimbal/internal/ecs/scenes"
	systems "github.com/jonesrussell/gimbal/internal/ecs/systems"
)

// updateCurrentScene handles scene-specific updates
func (g *ECSGame) updateCurrentScene() error {
	currentScene := g.sceneManager.GetCurrentScene()

	switch currentScene.GetType() {
	case scenes.SceneStudioIntro, scenes.SceneTitleScreen, scenes.SceneMenu:
		g.handleMenuInput()
	case scenes.ScenePlaying:
		return g.updatePlayingScene()
	case scenes.ScenePaused:
		g.handlePauseInput()
	}

	return nil
}

// updatePlayingScene handles gameplay updates
func (g *ECSGame) updatePlayingScene() error {
	// Check if paused
	if g.stateManager.IsPaused() {
		return nil
	}

	// Check for pause input
	if g.inputHandler.IsPausePressed() {
		g.stateManager.TogglePause()
		g.sceneManager.SwitchScene(scenes.ScenePaused)
		return nil
	}

	// Update player input and movement
	inputAngle := g.updatePlayerInput()

	// Update combat systems
	g.updateCombatSystems()

	// Update ECS systems
	g.updateECSSystems()

	// Handle player movement events
	g.handlePlayerMovementEvents(inputAngle)

	// Process all events
	g.eventSystem.ProcessEvents()

	return nil
}

// updatePlayerInput handles player input and movement
func (g *ECSGame) updatePlayerInput() common.Angle {
	inputAngle := g.inputHandler.GetMovementInput()
	core.PlayerInputSystem(g.world, inputAngle)
	return inputAngle
}

// updateCombatSystems updates all combat-related systems
func (g *ECSGame) updateCombatSystems() {
	g.handleWeaponFiring()
	g.enemySystem.Update(1.0) // Assuming 60fps, so deltaTime = 1.0
	g.weaponSystem.Update(1.0)
	g.collisionSystem.Update()
	g.healthSystem.Update() // Update health system (invincibility, game over, etc.)
}

// updateECSSystems runs all ECS systems
func (g *ECSGame) updateECSSystems() {
	// Get font from resource manager
	font := g.resourceManager.GetDefaultFont()
	if font == nil {
		g.logger.Error("Failed to get default font for score display")
		return
	}

	// Run core systems
	core.MovementSystem(g.world)
	core.OrbitalMovementSystem(g.world)
	core.StarMovementSystem(&ecs.ECS{World: g.world}, g.config)

	// Run score display system
	systems.ScoreDisplaySystem(g.world, font, g.scoreManager)()

// handlePlayerMovementEvents emits events when player moves
func (g *ECSGame) handlePlayerMovementEvents(inputAngle common.Angle) {
	if inputAngle != 0 {
		playerEntry := g.world.Entry(g.playerEntity)
		if playerEntry.Valid() {
			pos := core.Position.Get(playerEntry)
			orb := core.Orbital.Get(playerEntry)
			g.eventSystem.EmitPlayerMoved(*pos, orb.OrbitalAngle)
		}
	}
}

// handleWeaponFiring handles weapon firing based on input
func (g *ECSGame) handleWeaponFiring() {
	// Get player position and angle
	playerEntry := g.world.Entry(g.playerEntity)
	if !playerEntry.Valid() {
		return
	}

	pos := core.Position.Get(playerEntry)
	orb := core.Orbital.Get(playerEntry)

	// Check for shoot input (Space key)
	if g.inputHandler.IsShootPressed() {
		g.weaponSystem.FireWeapon(WeaponTypePrimary, *pos, orb.FacingAngle)
	}
}

// handleMenuInput handles input for menu scenes
func (g *ECSGame) handleMenuInput() {
	currentScene := g.sceneManager.GetCurrentScene()

	switch currentScene.GetType() {
	case scenes.SceneTitleScreen:
		// Any key to continue to main menu
		if g.inputHandler.GetLastEvent() != common.InputEventNone {
			g.sceneManager.SwitchScene(scenes.SceneMenu)
		}
	case scenes.SceneMenu:
		// Handle menu navigation
		if currentScene.GetType() == scenes.SceneMenu {
			if menuScene, ok := currentScene.(*scenes.MenuScene); ok {
				// Menu input is handled within the scene itself
				_ = menuScene // Use the variable to avoid unused variable warning
			}
		}
	}
}

// handlePauseInput handles input for pause menu
func (g *ECSGame) handlePauseInput() {
	// Pause scene handles its own input (ESC debounce logic)
	// No additional input handling needed here
}
