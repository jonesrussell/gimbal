package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type GameOverScene struct {
	manager *SceneManager
	font    text.Face
}

func NewGameOverScene(manager *SceneManager, font text.Face) *GameOverScene {
	return &GameOverScene{
		manager: manager,
		font:    font,
	}
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
	drawCenteredTextWithOptions(screen, TextDrawOptions{
		Text:  "GAME OVER",
		X:     centerX,
		Y:     centerY - 50,
		Alpha: 1.0,
		Font:  s.font,
	})

	// Draw instruction text
	drawCenteredTextWithOptions(screen, TextDrawOptions{
		Text:  "Press SPACE or ESC to return to menu",
		X:     centerX,
		Y:     centerY + 50,
		Alpha: 0.8,
		Font:  s.font,
	})
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
