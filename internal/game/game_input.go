package game

import (
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	weaponsys "github.com/jonesrussell/gimbal/internal/ecs/systems/weapon"
	"github.com/jonesrussell/gimbal/internal/scenes"
)

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
		// Update game state manager BEFORE switching scene
		g.stateManager.SetPaused(true)
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
