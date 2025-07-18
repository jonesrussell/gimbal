package scenes

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/common"
)

type StudioIntroScene struct {
	manager   *SceneManager
	font      text.Face
	startTime time.Time
	minTime   float64
	maxTime   float64
	finished  bool
}

func NewStudioIntroScene(manager *SceneManager, font text.Face) *StudioIntroScene {
	return &StudioIntroScene{
		manager:   manager,
		font:      font,
		startTime: time.Now(),
		minTime:   2.0, // Minimum 2 seconds
		maxTime:   4.0, // Maximum 4 seconds
		finished:  false,
	}
}

func (s *StudioIntroScene) Update() error {
	elapsed := time.Since(s.startTime).Seconds()
	if s.finished {
		return nil
	}
	// Allow skip after minTime with any key or mouse
	if elapsed >= s.minTime {
		input := s.manager.GetInputHandler()
		if input != nil && (input.GetLastEvent() != common.InputEventNone) {
			s.finished = true
			s.manager.SwitchScene(SceneTitleScreen)
			return nil
		}
	}
	// Auto-advance after maxTime
	if elapsed >= s.maxTime {
		s.finished = true
		s.manager.SwitchScene(SceneTitleScreen)
	}
	return nil
}

func (s *StudioIntroScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	elapsed := time.Since(s.startTime).Seconds()
	fadeProgress := elapsed / s.maxTime
	if fadeProgress > 1.0 {
		fadeProgress = 1.0
	}
	drawCenteredTextWithOptions(screen, TextDrawOptions{
		Text:  "GIMBAL STUDIOS",
		X:     float64(s.manager.config.ScreenSize.Width) / 2,
		Y:     float64(s.manager.config.ScreenSize.Height) / 2,
		Alpha: fadeProgress,
		Font:  s.font,
	})
	drawCenteredTextWithOptions(screen, TextDrawOptions{
		Text:  "Presents",
		X:     float64(s.manager.config.ScreenSize.Width) / 2,
		Y:     float64(s.manager.config.ScreenSize.Height)/2 + 50,
		Alpha: fadeProgress * 0.8,
		Font:  s.font,
	})
	drawCenteredTextWithOptions(screen, TextDrawOptions{
		Text:  "Press any key...",
		X:     float64(s.manager.config.ScreenSize.Width) / 2,
		Y:     float64(s.manager.config.ScreenSize.Height)/2 + 40,
		Alpha: 1.0,
		Font:  s.font,
	})
}

func (s *StudioIntroScene) Enter() {
	s.manager.logger.Debug("Entering studio intro scene")
	s.startTime = time.Now()
}

func (s *StudioIntroScene) Exit() {
	s.manager.logger.Debug("Exiting studio intro scene")
}

func (s *StudioIntroScene) GetType() SceneType {
	return SceneStudioIntro
}
