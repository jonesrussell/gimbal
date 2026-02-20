package intro

import (
	"context"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/internal/common"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	scenes "github.com/jonesrussell/gimbal/internal/scenes"
)

type StudioIntroScene struct {
	manager      *scenes.SceneManager
	resourceMgr  *resources.ResourceManager
	startTime    time.Time
	minTime      float64
	maxTime      float64
	finished     bool
	studioScreen *ebiten.Image
}

func NewStudioIntroScene(manager *scenes.SceneManager, resourceMgr *resources.ResourceManager) *StudioIntroScene {
	// Try to load studio screen image
	var studioScreen *ebiten.Image
	if resourceMgr != nil {
		if screen, ok := resourceMgr.GetSprite(context.Background(), "studio_screen"); ok {
			studioScreen = screen
		}
	}

	return &StudioIntroScene{
		manager:      manager,
		resourceMgr:  resourceMgr,
		startTime:    time.Now(),
		minTime:      2.0, // Minimum 2 seconds
		maxTime:      4.0, // Maximum 4 seconds
		finished:     false,
		studioScreen: studioScreen,
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
			s.manager.SwitchScene(scenes.SceneTitleScreen)
			return nil
		}
	}
	// Auto-advance after maxTime
	if elapsed >= s.maxTime {
		s.finished = true
		s.manager.SwitchScene(scenes.SceneTitleScreen)
	}
	return nil
}

func (s *StudioIntroScene) Draw(screen *ebiten.Image) {
	// Draw studio screen image if available
	if s.studioScreen != nil {
		config := s.manager.GetConfig()
		screenWidth := float64(config.ScreenSize.Width)
		screenHeight := float64(config.ScreenSize.Height)
		imgWidth := float64(s.studioScreen.Bounds().Dx())
		imgHeight := float64(s.studioScreen.Bounds().Dy())

		op := &ebiten.DrawImageOptions{}
		// Center the image on screen
		op.GeoM.Translate((screenWidth-imgWidth)/2, (screenHeight-imgHeight)/2)
		screen.DrawImage(s.studioScreen, op)
	}
}

func (s *StudioIntroScene) Enter() {
	s.startTime = time.Now()
}

func (s *StudioIntroScene) Exit() {}

func (s *StudioIntroScene) GetType() scenes.SceneType {
	return scenes.SceneStudioIntro
}
