package ecs

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/internal/common"
)

type OptionsScene struct {
	manager *SceneManager
}

func NewOptionsScene(manager *SceneManager) *OptionsScene {
	return &OptionsScene{manager: manager}
}

func (s *OptionsScene) Update() error {
	if s.manager.inputHandler.GetLastEvent() != common.InputEventNone {
		s.manager.SwitchScene(SceneMenu)
	}
	return nil
}

func (s *OptionsScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	drawCenteredText(
		screen,
		"OPTIONS\nComing Soon!",
		float64(s.manager.config.ScreenSize.Width)/2,
		float64(s.manager.config.ScreenSize.Height)/2,
		1.0,
	)
}

func (s *OptionsScene) Enter()             { s.manager.logger.Debug("Entering options scene") }
func (s *OptionsScene) Exit()              { s.manager.logger.Debug("Exiting options scene") }
func (s *OptionsScene) GetType() SceneType { return SceneOptions }
