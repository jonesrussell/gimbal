package game

import (
	"context"
	"time"

	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/errors"
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

	// Check if quit has been requested by a scene
	if g.sceneManager.IsQuitRequested() {
		g.logger.Info("Quit requested, stopping game loop")
		return errors.NewGameError(errors.StateTransition, "application shutdown requested")
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

	// Update Gyruss system (handles all enemy spawning, paths, behaviors, attacks, firing)
	gyrussUpdateFunc := func() error {
		return g.gyrussSystem.Update(ctx, deltaTime)
	}
	if err := g.updateSystemWithTiming("gyruss", gyrussUpdateFunc); err != nil {
		return err
	}

	// Update stage state machine (drives wave/boss transitions; reads wave manager counts)
	g.stageStateMachine.Update(ctx, deltaTime)

	// Process queued events so onBossDefeated runs before level completion check (same frame)
	g.eventSystem.ProcessEvents()

	// Update weapon system (player weapons)
	weaponUpdateFunc := func() error {
		g.weaponSystem.Update(deltaTime)
		return nil
	}
	if err := g.updateSystemWithTiming("weapon", weaponUpdateFunc); err != nil {
		return err
	}

	// Check for level completion (boss killed)
	g.checkLevelCompletion()

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
