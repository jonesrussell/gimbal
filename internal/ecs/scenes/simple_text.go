package ecs

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/internal/common"
)

type SimpleTextScene struct {
	manager   *SceneManager
	text      string
	sceneType SceneType
}

func NewSimpleTextScene(manager *SceneManager, text string, sceneType SceneType) *SimpleTextScene {
	return &SimpleTextScene{
		manager:   manager,
		text:      text,
		sceneType: sceneType,
	}
}

func (s *SimpleTextScene) Update() error {
	if s.manager.inputHandler.GetLastEvent() != common.InputEventNone {
		s.manager.SwitchScene(SceneMenu)
	}
	return nil
}

func (s *SimpleTextScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	drawCenteredText(
		screen,
		s.text,
		float64(s.manager.config.ScreenSize.Width)/2,
		float64(s.manager.config.ScreenSize.Height)/2,
		1.0,
	)
}

func (s *SimpleTextScene) Enter()             { s.manager.logger.Debug("Entering scene", "scene", s.sceneType) }
func (s *SimpleTextScene) Exit()              { s.manager.logger.Debug("Exiting scene", "scene", s.sceneType) }
func (s *SimpleTextScene) GetType() SceneType { return s.sceneType }
