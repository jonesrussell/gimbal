package ecs

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/internal/common"
)

type CreditsScene struct {
	manager *SceneManager
}

func NewCreditsScene(manager *SceneManager) *CreditsScene {
	return &CreditsScene{manager: manager}
}

func (s *CreditsScene) Update() error {
	if s.manager.inputHandler.GetLastEvent() != common.InputEventNone {
		s.manager.SwitchScene(SceneMenu)
	}
	return nil
}

func (s *CreditsScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	drawCenteredText(
		screen,
		"CREDITS\nGimbal Studios\n2025",
		float64(s.manager.config.ScreenSize.Width)/2,
		float64(s.manager.config.ScreenSize.Height)/2,
		1.0,
	)
}

func (s *CreditsScene) Enter()             { s.manager.logger.Debug("Entering credits scene") }
func (s *CreditsScene) Exit()              { s.manager.logger.Debug("Exiting credits scene") }
func (s *CreditsScene) GetType() SceneType { return SceneCredits }
