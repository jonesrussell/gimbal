package game

import (
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
)

// checkLevelCompletion checks if the stage is complete and advances to the next
func (g *ECSGame) checkLevelCompletion() {
	// Check if current stage is complete (boss defeated)
	if !g.gyrussSystem.IsStageComplete() {
		return
	}

	g.handleLevelComplete()
}

// checkCompletionConditions checks all completion conditions
func (g *ECSGame) checkCompletionConditions(conditions managers.CompletionConditions) bool {
	// Check if boss kill is required
	if conditions.RequireBossKill {
		if !g.gyrussSystem.WasBossSpawned() || g.gyrussSystem.IsBossActive() {
			return false
		}
	}

	// Check if all waves are required
	if conditions.RequireAllWaves {
		if g.gyrussSystem.GetWaveManager().HasMoreWaves() {
			return false
		}
	}

	return true
}

// handleLevelComplete handles level completion actions
func (g *ECSGame) handleLevelComplete() {
	currentStage := g.gyrussSystem.GetCurrentStage()
	g.logger.Debug("Stage complete", "stage", currentStage)

	// Update level manager to track progression
	g.levelManager.IncrementLevel()

	// Load next stage
	if err := g.gyrussSystem.LoadNextStage(); err != nil {
		g.logger.Warn("Failed to load next stage, resetting", "error", err)
		g.gyrussSystem.Reset()
		return
	}

	nextStage := g.gyrussSystem.GetCurrentStage()
	g.logger.Debug("Next stage loaded", "stage", nextStage)

	// Show level title for new stage
	if currentScene := g.sceneManager.GetCurrentScene(); currentScene != nil {
		if playingScene, ok := currentScene.(interface{ ShowLevelTitle(int) }); ok {
			playingScene.ShowLevelTitle(nextStage)
		}
	}
}
