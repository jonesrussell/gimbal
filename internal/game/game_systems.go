package game

import (
	"context"
	"time"

	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/scenes"
)

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

