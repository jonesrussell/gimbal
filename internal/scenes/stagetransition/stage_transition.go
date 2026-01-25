package stagetransition

import (
	"context"
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/scenes"
)

const (
	stageTransitionDuration = 1.5 // seconds
)

type StageTransitionScene struct {
	manager      *scenes.SceneManager
	font         text.Face
	resourceMgr  *resources.ResourceManager
	scoreManager *managers.ScoreManager
	startTime    time.Time
	nextPlanet   string
	bossPortrait *ebiten.Image
	warpFrames   []*ebiten.Image
	currentFrame int
	soundPlayed  bool
}

func NewStageTransitionScene(
	manager *scenes.SceneManager,
	font text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) *StageTransitionScene {
	return &StageTransitionScene{
		manager:      manager,
		font:         font,
		resourceMgr:  resourceMgr,
		scoreManager: scoreManager,
		startTime:    time.Now(),
		nextPlanet:   "Mars",
		warpFrames:   make([]*ebiten.Image, 0, 8),
	}
}

// SetNextPlanet sets the next planet name for the transition
func (s *StageTransitionScene) SetNextPlanet(planetName string) {
	s.nextPlanet = planetName
}

func (s *StageTransitionScene) Update() error {
	elapsed := time.Since(s.startTime).Seconds()

	// Auto-advance after duration
	if elapsed >= stageTransitionDuration {
		// Load next stage and set stage intro info before switching
		s.loadNextStage()
		s.manager.SwitchScene(scenes.SceneStageIntro)
		return nil
	}

	// Update warp tunnel frame
	if len(s.warpFrames) > 0 {
		frameIndex := int((elapsed / stageTransitionDuration) * float64(len(s.warpFrames)))
		if frameIndex >= len(s.warpFrames) {
			frameIndex = len(s.warpFrames) - 1
		}
		s.currentFrame = frameIndex
	}

	// Play warp sound at start
	if !s.soundPlayed && elapsed > 0.1 {
		s.playWarpSound()
		s.soundPlayed = true
	}

	return nil
}

func (s *StageTransitionScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	config := s.manager.GetConfig()
	centerX := float64(config.ScreenSize.Width) / 2
	centerY := float64(config.ScreenSize.Height) / 2
	elapsed := time.Since(s.startTime).Seconds()

	// Draw warp tunnel frame
	if len(s.warpFrames) > 0 && s.currentFrame < len(s.warpFrames) && s.warpFrames[s.currentFrame] != nil {
		op := &ebiten.DrawImageOptions{}
		// Fade effect
		alpha := 1.0 - (elapsed/stageTransitionDuration)*0.5 // Fade out slightly
		op.ColorScale.SetA(float32(alpha))
		screen.DrawImage(s.warpFrames[s.currentFrame], op)
	}

	// Draw "TRAVELING TO" text
	travelText := fmt.Sprintf("TRAVELING TO %s", s.nextPlanet)
	fadeAlpha := math.Min(1.0, elapsed/0.3)
	scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
		Text:  travelText,
		X:     centerX,
		Y:     centerY - 40,
		Alpha: fadeAlpha,
		Font:  s.font,
	})

	// Draw boss portrait if available
	if s.bossPortrait != nil {
		op := &ebiten.DrawImageOptions{}
		portraitScale := 0.3 + 0.7*math.Min(1.0, elapsed/0.5) // Zoom in
		portraitWidth := float64(s.bossPortrait.Bounds().Dx()) * portraitScale
		op.GeoM.Scale(portraitScale, portraitScale)
		op.GeoM.Translate(centerX-portraitWidth/2, centerY+20)
		op.ColorScale.SetA(float32(fadeAlpha))
		screen.DrawImage(s.bossPortrait, op)
	}
}

func (s *StageTransitionScene) Enter() {
	s.manager.GetLogger().Debug("Entering stage transition scene",
		"next_planet", s.nextPlanet)
	s.startTime = time.Now()
	s.soundPlayed = false
	s.currentFrame = 0

	// Load warp tunnel frames
	if s.resourceMgr != nil {
		s.warpFrames = make([]*ebiten.Image, 0, 8)
		for i := 1; i <= 8; i++ {
			frameName := fmt.Sprintf("warp_tunnel_%02d", i)
			if frame, ok := s.resourceMgr.GetSprite(context.Background(), frameName); ok {
				s.warpFrames = append(s.warpFrames, frame)
			}
		}

		// Load boss portrait (try to match planet name)
		bossName := s.nextPlanet
		if bossPortrait, ok := s.resourceMgr.GetSprite(context.Background(), fmt.Sprintf("boss_portrait_%s", bossName)); ok {
			s.bossPortrait = bossPortrait
		}
	}
}

func (s *StageTransitionScene) Exit() {
	s.manager.GetLogger().Debug("Exiting stage transition scene")
}

func (s *StageTransitionScene) GetType() scenes.SceneType {
	return scenes.SceneStageTransition
}

func (s *StageTransitionScene) playWarpSound() {
	if s.resourceMgr == nil {
		return
	}

	audioPlayer := s.resourceMgr.GetAudioPlayer()
	if audioPlayer == nil {
		return
	}

	// Play warp sound effect
	s.manager.GetLogger().Debug("Stage transition warp sound should play")
}

// loadNextStage loads the next stage into the game system
func (s *StageTransitionScene) loadNextStage() {
	// The stage loading is handled by the game's level completion system
	// This is called to ensure the stage is ready before showing the intro
	s.manager.GetLogger().Debug("Stage transition complete, next stage should be loaded")
}
