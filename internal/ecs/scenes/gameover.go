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
	// Handle game over input
	if s.manager.inputHandler.IsShootPressed() || s.manager.inputHandler.IsPausePressed() {
		// Return to menu
		s.manager.SwitchScene(SceneMenu)
	}
	return nil
}

func (s *GameOverScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	// Draw game over text
	centerX := float64(s.manager.config.ScreenSize.Width) / 2
	centerY := float64(s.manager.config.ScreenSize.Height) / 2
	drawCenteredText(screen, "GAME OVER", centerX, centerY-50, 1.0)

	// Draw instruction text
	drawCenteredText(screen, "Press SPACE or ESC to return to menu", centerX, centerY+50, 0.8)
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
