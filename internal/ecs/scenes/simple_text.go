package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/common"
)

type SimpleTextScene struct {
	manager   *SceneManager
	text      string
	sceneType SceneType
	font      text.Face
}

func NewSimpleTextScene(manager *SceneManager, textStr string, sceneType SceneType, font text.Face) *SimpleTextScene {
	return &SimpleTextScene{
		manager:   manager,
		text:      textStr,
		sceneType: sceneType,
		font:      font,
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
	drawCenteredTextWithOptions(
		screen,
		TextDrawOptions{
			Text:  s.text,
			X:     float64(s.manager.config.ScreenSize.Width) / 2,
			Y:     float64(s.manager.config.ScreenSize.Height) / 2,
			Alpha: 1.0,
			Font:  s.font,
		},
	)
}

func (s *SimpleTextScene) Enter()             { s.manager.logger.Debug("Entering scene", "scene", s.sceneType) }
func (s *SimpleTextScene) Exit()              { s.manager.logger.Debug("Exiting scene", "scene", s.sceneType) }
func (s *SimpleTextScene) GetType() SceneType { return s.sceneType }
