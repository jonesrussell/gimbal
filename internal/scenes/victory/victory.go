package victory

import (
	"context"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/scenes"
	"github.com/jonesrussell/gimbal/internal/scenes/effects"
)

const (
	victoryDuration = 8.0 // seconds
	epilogueStart   = 4.0 // seconds - when epilogue text appears
)

type VictoryScene struct {
	manager       *scenes.SceneManager
	font          text.Face
	resourceMgr   *resources.ResourceManager
	scoreManager  *managers.ScoreManager
	startTime     time.Time
	starfield     *effects.Starfield
	missionBanner *ebiten.Image
	starfieldBg   *ebiten.Image
	scrollOffset  float64
	soundPlayed   bool
}

func NewVictoryScene(
	manager *scenes.SceneManager,
	font text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) *VictoryScene {
	config := manager.GetConfig()
	starfield := effects.NewStarfield(
		config.ScreenSize.Width,
		config.ScreenSize.Height,
		150, // More stars for ending
		0.5, // Slower speed
	)

	return &VictoryScene{
		manager:      manager,
		font:         font,
		resourceMgr:  resourceMgr,
		scoreManager: scoreManager,
		startTime:    time.Now(),
		starfield:    starfield,
		scrollOffset: 0,
	}
}

func (s *VictoryScene) Update() error {
	elapsed := time.Since(s.startTime).Seconds()

	// Update starfield
	deltaTime := 1.0 / 60.0
	s.starfield.Update(deltaTime)

	// Update scroll offset for planetary fly-through effect
	s.scrollOffset += 0.5 * deltaTime * 60.0

	// Play victory fanfare at start
	if !s.soundPlayed && elapsed > 0.1 {
		s.playVictoryFanfare()
		s.soundPlayed = true
	}

	// Auto-advance to credits after duration
	if elapsed >= victoryDuration {
		s.manager.SwitchScene(scenes.SceneCredits)
		return nil
	}

	return nil
}

func (s *VictoryScene) Draw(screen *ebiten.Image) {
	config := s.manager.GetConfig()
	centerX := float64(config.ScreenSize.Width) / 2
	centerY := float64(config.ScreenSize.Height) / 2
	elapsed := time.Since(s.startTime).Seconds()

	// Draw starfield background
	screen.Fill(color.Black)
	if s.starfieldBg != nil {
		// Draw tiled starfield background with scroll
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, -math.Mod(s.scrollOffset, float64(s.starfieldBg.Bounds().Dy())))
		screen.DrawImage(s.starfieldBg, op)
		// Draw second tile for seamless scrolling
		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Translate(0, -math.Mod(s.scrollOffset, float64(s.starfieldBg.Bounds().Dy()))+float64(s.starfieldBg.Bounds().Dy()))
		screen.DrawImage(s.starfieldBg, op2)
	} else {
		// Fallback to animated starfield
		s.starfield.Draw(screen)
	}

	// Fade in mission complete banner
	bannerAlpha := math.Min(1.0, elapsed/1.0)
	if s.missionBanner != nil {
		op := &ebiten.DrawImageOptions{}
		bannerWidth := float64(s.missionBanner.Bounds().Dx())
		bannerHeight := float64(s.missionBanner.Bounds().Dy())
		op.GeoM.Translate(centerX-bannerWidth/2, centerY-bannerHeight/2-60)
		op.ColorScale.SetA(float32(bannerAlpha))
		screen.DrawImage(s.missionBanner, op)
	} else {
		// Fallback text
		scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
			Text:  "MISSION COMPLETE",
			X:     centerX,
			Y:     centerY - 60,
			Alpha: bannerAlpha,
			Font:  s.font,
		})
	}

	// Draw epilogue text after delay
	if elapsed >= epilogueStart {
		epilogueAlpha := math.Min(1.0, (elapsed-epilogueStart)/1.0)
		epilogueText := "The galaxy is safe once more."
		scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
			Text:  epilogueText,
			X:     centerX,
			Y:     centerY + 80,
			Alpha: epilogueAlpha * 0.8,
			Font:  s.font,
		})
	}
}

func (s *VictoryScene) Enter() {
	s.manager.GetLogger().Debug("Entering victory scene")
	s.startTime = time.Now()
	s.scrollOffset = 0
	s.starfield.Reset()
	s.soundPlayed = false

	// Load assets
	if s.resourceMgr != nil {
		if banner, ok := s.resourceMgr.GetSprite(context.Background(), "mission_complete"); ok {
			s.missionBanner = banner
		}
		if bg, ok := s.resourceMgr.GetSprite(context.Background(), "starfield_bg"); ok {
			s.starfieldBg = bg
		}
	}
}

func (s *VictoryScene) Exit() {
	s.manager.GetLogger().Debug("Exiting victory scene")
}

func (s *VictoryScene) GetType() scenes.SceneType {
	return scenes.SceneVictory
}

func (s *VictoryScene) playVictoryFanfare() {
	if s.resourceMgr == nil {
		return
	}

	audioPlayer := s.resourceMgr.GetAudioPlayer()
	if audioPlayer == nil {
		return
	}

	// Play victory fanfare
	s.manager.GetLogger().Debug("Victory fanfare should play")
}
