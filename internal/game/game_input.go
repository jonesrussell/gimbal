package game

import (
	"github.com/jonesrussell/gimbal/internal/dbg"
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
		dbg.Log(dbg.System, "Pause key pressed, switching to pause scene")
		// Update game state manager BEFORE switching scene
		g.stateManager.SetPaused(true)
		g.sceneManager.SwitchScene(scenes.ScenePaused)
	}
}

// handleShootingInput processes shooting input and fires weapons
func (g *ECSGame) handleShootingInput() {
	// Only handle shooting if we have a valid player entity
	if g.playerEntity == 0 {
		return
	}

	// Check if shoot key is pressed
	if g.inputHandler.IsShootPressed() {
		// Get player position and angle
		playerEntry := g.world.Entry(g.playerEntity)
		if !playerEntry.Valid() {
			return
		}

		pos := core.Position.Get(playerEntry)
		orbital := core.Orbital.Get(playerEntry)

		if pos == nil || orbital == nil {
			return
		}

		// Fire weapon with player position and facing angle
		if g.weaponSystem.FireWeapon(weaponsys.WeaponTypePrimary, *pos, orbital.FacingAngle) {
			dbg.Log(dbg.System, "Weapon fired")
		}
		// else: fire blocked by timing (no log to avoid spam)
	}
}
