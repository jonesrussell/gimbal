package ecs

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type PausedScene struct {
	manager *SceneManager
}

func NewPausedScene(manager *SceneManager) *PausedScene {
	return &PausedScene{manager: manager}
}

func (s *PausedScene) Update() error {
	// Handle pause menu input (to be implemented)
	return nil
}

func (s *PausedScene) Draw(screen *ebiten.Image) {
	// Draw pause overlay
	screen.Fill(color.Black)
	// Draw pause text (simplified)
}

func (s *PausedScene) Enter() {
	s.manager.logger.Debug("Entering paused scene")
}

func (s *PausedScene) Exit() {
	s.manager.logger.Debug("Exiting paused scene")
}

func (s *PausedScene) GetType() SceneType {
	return ScenePaused
}
