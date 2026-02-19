package stageintro

import (
	"context"
	"fmt"
	"image/color"
	"math"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/scenes"
)

const (
	stageIntroDuration = 2.5 // seconds
)

type StageIntroScene struct {
	manager      *scenes.SceneManager
	font         text.Face
	resourceMgr  *resources.ResourceManager
	scoreManager *managers.ScoreManager
	startTime    time.Time
	stageNumber  int
	fromPlanet   string
	toPlanet     string
	message      string
	planetSprite *ebiten.Image
	gridOverlay  *ebiten.Image
	soundPlayed  bool
}

func NewStageIntroScene(
	manager *scenes.SceneManager,
	font text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) *StageIntroScene {
	return &StageIntroScene{
		manager:      manager,
		font:         font,
		resourceMgr:  resourceMgr,
		scoreManager: scoreManager,
		startTime:    time.Now(),
		stageNumber:  1,
		fromPlanet:   "Earth",
		toPlanet:     "Mars",
		message:      "ENEMY ACTIVITY DETECTED",
	}
}

// SetStageInfo sets the stage information for the intro
func (s *StageIntroScene) SetStageInfo(stageNumber int, fromPlanet, toPlanet, message string) {
	s.stageNumber = stageNumber
	s.fromPlanet = fromPlanet
	s.toPlanet = toPlanet
	if message == "" {
		s.message = "ENEMY ACTIVITY DETECTED"
	} else {
		s.message = message
	}
}

func (s *StageIntroScene) Update() error {
	elapsed := time.Since(s.startTime).Seconds()

	// Auto-advance after duration
	if elapsed >= stageIntroDuration {
		// Load the stage before switching to playing
		s.loadStage()
		s.manager.SwitchScene(scenes.ScenePlaying)
		return nil
	}

	// Play whoosh sound at start
	if !s.soundPlayed && elapsed > 0.1 {
		s.playWhooshSound()
		s.soundPlayed = true
	}

	return nil
}

func (s *StageIntroScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	config := s.manager.GetConfig()
	centerX := float64(config.ScreenSize.Width) / 2
	centerY := float64(config.ScreenSize.Height) / 2
	elapsed := time.Since(s.startTime).Seconds()

	// Fade in effect
	fadeAlpha := math.Min(1.0, elapsed/0.5)

	// Draw scanning grid overlay
	if s.gridOverlay != nil {
		op := &ebiten.DrawImageOptions{}
		gridAlpha := 0.3 + 0.2*math.Sin(elapsed*4.0) // Pulsing effect
		op.ColorScale.SetA(float32(gridAlpha * fadeAlpha))
		// Center the grid on the screen
		gridWidth := float64(s.gridOverlay.Bounds().Dx())
		gridHeight := float64(s.gridOverlay.Bounds().Dy())
		op.GeoM.Translate(centerX-gridWidth/2, centerY-gridHeight/2)
		screen.DrawImage(s.gridOverlay, op)
	}

	// Draw stage number
	stageText := fmt.Sprintf("STAGE %d", s.stageNumber)
	scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
		Text:  stageText,
		X:     centerX,
		Y:     centerY - 120,
		Alpha: fadeAlpha,
		Font:  s.font,
	})

	// Draw planet sprite
	if s.planetSprite != nil {
		op := &ebiten.DrawImageOptions{}
		planetScale := 0.5 + 0.5*math.Min(1.0, elapsed/0.8) // Zoom in effect
		planetWidth := float64(s.planetSprite.Bounds().Dx()) * planetScale
		planetHeight := float64(s.planetSprite.Bounds().Dy()) * planetScale
		op.GeoM.Scale(planetScale, planetScale)
		op.GeoM.Translate(centerX-planetWidth/2, centerY-planetHeight/2-20)
		op.ColorScale.SetA(float32(fadeAlpha))
		screen.DrawImage(s.planetSprite, op)
	}

	// Draw planet route
	routeText := fmt.Sprintf("%s → %s", s.fromPlanet, s.toPlanet)
	scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
		Text:  routeText,
		X:     centerX,
		Y:     centerY + 60,
		Alpha: fadeAlpha,
		Font:  s.font,
	})

	// Draw message
	scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
		Text:  s.message,
		X:     centerX,
		Y:     centerY + 120,
		Alpha: fadeAlpha * 0.8,
		Font:  s.font,
	})
}

func (s *StageIntroScene) Enter() {
	s.manager.GetLogger().Debug("Entering stage intro scene",
		"stage", s.stageNumber,
		"from", s.fromPlanet,
		"to", s.toPlanet)
	s.startTime = time.Now()
	s.soundPlayed = false

	// If stage info not set, try to get it from level manager
	if s.stageNumber == 0 {
		if levelMgr := s.manager.GetLevelManager(); levelMgr != nil {
			s.stageNumber = levelMgr.GetLevel()
			// Map stage numbers to planets (stage 1 = Earth→Mars, stage 2 = Mars→Jupiter, …)
			planets := []string{"Earth", "Mars", "Jupiter", "Saturn", "Uranus", "Neptune"}
			if s.stageNumber == 1 {
				s.fromPlanet = "Earth"
				s.toPlanet = "Mars"
			} else if s.stageNumber >= 2 && s.stageNumber < len(planets) {
				s.fromPlanet = planets[s.stageNumber-1]
				s.toPlanet = planets[s.stageNumber]
			} else if s.stageNumber == len(planets) {
				s.fromPlanet = planets[s.stageNumber-2]
				s.toPlanet = planets[s.stageNumber-1]
			}
		}
	}

	// Load planet sprite
	if s.resourceMgr != nil {
		planetName := s.toPlanet
		// Convert to lowercase for asset lookup
		planetNameLower := strings.ToLower(planetName)
		if planetSprite, ok := s.resourceMgr.GetSprite(context.Background(), fmt.Sprintf("planet_%s", planetNameLower)); ok {
			s.planetSprite = planetSprite
		}

		// Load scanning grid
		if gridSprite, ok := s.resourceMgr.GetSprite(context.Background(), "scanning_grid"); ok {
			s.gridOverlay = gridSprite
		}
	}
}

func (s *StageIntroScene) Exit() {
	s.manager.GetLogger().Debug("Exiting stage intro scene")
}

func (s *StageIntroScene) GetType() scenes.SceneType {
	return scenes.SceneStageIntro
}

func (s *StageIntroScene) playWhooshSound() {
	if s.resourceMgr == nil {
		return
	}

	audioPlayer := s.resourceMgr.GetAudioPlayer()
	if audioPlayer == nil {
		return
	}

	// Play whoosh sound effect (one-shot, not music)
	// Note: For sound effects, we'd need a separate sound effect player
	// For now, we'll just log that it should play
	s.manager.GetLogger().Debug("Stage intro whoosh sound should play")
}

// loadStage loads the stage into the game system
func (s *StageIntroScene) loadStage() {
	// The stage should already be loaded by the stage transition scene
	// or by the game system when starting. This is a placeholder for
	// any additional stage setup needed.
	s.manager.GetLogger().Debug("Stage intro complete, loading stage", "stage", s.stageNumber)
}