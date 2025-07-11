package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameOverScene struct {
	manager *SceneManager
}

func NewGameOverScene(manager *SceneManager) *GameOverScene {
	return &GameOverScene{manager: manager}
}

func (s *GameOverScene) Update() error {
	// Handle game over input (to be implemented)
	return nil
}

func (s *GameOverScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	// Draw game over screen (to be implemented)
}

func (s *GameOverScene) Enter() {
	s.manager.logger.Debug("Entering game over scene")
}

func (s *GameOverScene) Exit() {
	s.manager.logger.Debug("Exiting game over scene")
}

func (s *GameOverScene) GetType() SceneType {
	return SceneGameOver
}
