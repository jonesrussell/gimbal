package game

import (
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	enemysys "github.com/jonesrussell/gimbal/internal/ecs/systems/enemy"
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

// checkLevelCompletion checks if the boss is killed and advances the level
func (g *ECSGame) checkLevelCompletion() {
	// Get current level config to check completion conditions
	levelConfig := g.levelManager.GetCurrentLevelConfig()
	if levelConfig == nil {
		return
	}

	// Check completion conditions
	canComplete := g.checkCompletionConditions(levelConfig.CompletionConditions)

	if canComplete {
		g.handleLevelComplete()
	}
}

// checkCompletionConditions checks all completion conditions
func (g *ECSGame) checkCompletionConditions(conditions managers.CompletionConditions) bool {
	// Check if boss kill is required
	if conditions.RequireBossKill {
		if !g.enemySystem.WasBossSpawned() || g.enemySystem.IsBossActive() {
			return false
		}
	}

	// Check if all waves are required
	if conditions.RequireAllWaves {
		if g.enemySystem.GetWaveManager().HasMoreWaves() {
			return false
		}
	}

	// Check if all enemies must be killed
	// This would require checking active enemy count
	// For now, we'll assume boss kill + all waves = all enemies killed
	// TODO: Implement active enemy count check
	_ = conditions.RequireAllEnemiesKilled

	return true
}

// handleLevelComplete handles level completion actions
func (g *ECSGame) handleLevelComplete() {
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
