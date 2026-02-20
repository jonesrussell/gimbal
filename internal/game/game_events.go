package game

import (
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/ecs/events"
	"github.com/jonesrussell/gimbal/internal/scenes"
	"github.com/jonesrussell/gimbal/internal/ui/presenter"
)

// setupEventSubscriptions sets up event handlers
func (g *ECSGame) setupEventSubscriptions() {
	// Initialize HUD presenter with current game state
	g.initializeHUDPresenter()
	// Subscribe to player movement events (no per-move logging to avoid 60/s spam)
	g.eventSystem.SubscribeToPlayerMoved(func(w donburi.World, event events.PlayerMovedEvent) {})

	// Subscribe to game state events
	g.eventSystem.SubscribeToGameState(func(w donburi.World, event events.GameStateEvent) {
		g.logger.Debug("Game state changed", "is_paused", event.IsPaused)
	})

	// Subscribe to score changes
	g.eventSystem.SubscribeToScoreChanged(func(w donburi.World, event events.ScoreChangedEvent) {
		g.logger.Debug("Score changed",
			"old_score", event.OldScore,
			"new_score", event.NewScore,
			"delta", event.Delta)
	})

	// Subscribe to game over events
	g.eventSystem.SubscribeToGameOver(func(w donburi.World, event events.GameOverEvent) {
		g.logger.Debug("Game over triggered", "reason", event.Reason)
		g.sceneManager.SwitchScene(scenes.SceneGameOver)
	})

	// Subscribe to player damage events for screen shake
	g.eventSystem.SubscribeToPlayerDamaged(func(w donburi.World, event events.PlayerDamagedEvent) {
		g.logger.Debug("Player damaged", "damage", event.Damage, "remaining_lives", event.RemainingLives)
		// TODO: Re-enable after PlayingScene is moved to gameplay package
		// Trigger screen shake if we're in the playing scene
		// if g.sceneManager.GetCurrentScene().GetType() == scenes.ScenePlaying {
		// 	if playingScene, ok := g.sceneManager.GetCurrentScene().(*scenes.PlayingScene); ok {
		// 		playingScene.TriggerScreenShake()
		// 	}
		// }
	})

	// Subscribe to enemy destroyed events for scoring
	g.eventSystem.SubscribeToEnemyDestroyed(func(w donburi.World, event events.EnemyDestroyedEvent) {
		g.scoreManager.AddScore(event.Points)
		g.logger.Debug(
			"Score added from enemy destruction",
			"points", event.Points,
			"total_score", g.scoreManager.GetScore(),
		)
	})
}

// initializeHUDPresenter creates and subscribes the HUD presenter to events
func (g *ECSGame) initializeHUDPresenter() {
	// Get initial values from game state
	current, maximum := g.healthSystem.GetPlayerHealth()
	initialHealth := 1.0
	if maximum > 0 {
		initialHealth = float64(current) / float64(maximum)
	}

	g.hudPresenter = presenter.NewHUDPresenter(
		g.scoreManager.GetScore(),
		current,
		g.levelManager.GetLevel(),
		initialHealth,
	)

	// Subscribe to events
	g.hudPresenter.Subscribe(g.eventSystem)
	g.logger.Debug("HUD presenter initialized and subscribed to events")
}
