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
	inputHandler := s.manager.GetInputHandler()
	if inputHandler.IsShootPressed() || inputHandler.IsPausePressed() {
		// Return to menu
		s.manager.SwitchScene(SceneMenu)
	}
	return nil
}

func (s *GameOverScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	// Draw game over text
	config := s.manager.GetConfig()
	centerX := float64(config.ScreenSize.Width) / 2
	centerY := float64(config.ScreenSize.Height) / 2
	DrawCenteredTextWithOptions(screen, TextDrawOptions{
		Text:  "GAME OVER",
		X:     centerX,
		Y:     centerY - 50,
		Alpha: 1.0,
		Font:  s.font,
	})

	// Draw instruction text
	DrawCenteredTextWithOptions(screen, TextDrawOptions{
		Text:  "Press SPACE or ESC to return to menu",
		X:     centerX,
		Y:     centerY + 50,
		Alpha: 0.8,
		Font:  s.font,
	})
}

func (s *GameOverScene) Enter() {
	s.manager.GetLogger().Debug("Entering game over scene")
}

func (s *GameOverScene) Exit() {
	s.manager.GetLogger().Debug("Exiting game over scene")
}

func (s *GameOverScene) GetType() SceneType {
	return SceneGameOver
}
