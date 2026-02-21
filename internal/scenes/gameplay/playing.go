// Package gameplay provides the main gameplay scene for the game.
// It handles the active game state, entity updates, rendering, and game logic during play.
package gameplay

import (
	"context"
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	v2text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/enemy"
	"github.com/jonesrussell/gimbal/internal/scenes"
	"github.com/jonesrussell/gimbal/internal/scenes/bossintro"
)

// Music track name constants
const (
	musicTrackMain   = "game_music_main"
	musicTrackLevel1 = "game_music_level_1"
	musicTrackLevel2 = "game_music_level_2"
	musicTrackBoss   = "game_music_boss"
)

type PlayingScene struct {
	manager      *scenes.SceneManager
	screenShake  float64 // Screen shake intensity (0 = no shake)
	font         v2text.Face
	scoreManager *managers.ScoreManager
	resourceMgr  *resources.ResourceManager

	// Level title display
	levelTitleStartTime time.Time
	showLevelTitle      bool
	currentLevelNumber  int
	levelTitleDuration  float64 // Duration to show title in seconds

	// Music state tracking
	currentMusicTrack string // Track which music is currently playing
	bossWasActive     bool   // Track if boss was active last frame

	// Boss intro overlay
	bossIntroOverlay *bossintro.BossIntroOverlay
}

func NewPlayingScene(
	manager *scenes.SceneManager,
	font v2text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) *PlayingScene {
	scene := &PlayingScene{
		manager:            manager,
		font:               font,
		scoreManager:       scoreManager,
		resourceMgr:        resourceMgr,
		levelTitleDuration: 3.0, // Show title for 3 seconds
	}

	// Create boss intro overlay
	scene.bossIntroOverlay = bossintro.NewBossIntroOverlay(manager, font, scoreManager, resourceMgr)

	// UI is now handled by the main game's EbitenUI system

	return scene
}

func (s *PlayingScene) Update() error {
	// Update screen shake
	if s.screenShake > 0 {
		s.screenShake -= 0.1 // Reduce shake over time
		if s.screenShake < 0 {
			s.screenShake = 0
		}
	}

	// Update level title display
	if s.showLevelTitle {
		elapsed := time.Since(s.levelTitleStartTime).Seconds()
		if elapsed >= s.levelTitleDuration {
			s.showLevelTitle = false
		}
	}

	// Update boss intro overlay
	if s.bossIntroOverlay != nil {
		deltaTime := 1.0 / 60.0 // Assume 60 FPS
		s.bossIntroOverlay.Update(deltaTime)
	}

	// Check for boss and switch music accordingly
	s.updateBossMusic()

	return nil
}

// updateBossMusic checks if boss is active and switches music accordingly
func (s *PlayingScene) updateBossMusic() {
	// Check if boss is active by querying the world
	bossActive := s.isBossActive()

	// If boss state changed, switch music
	if bossActive != s.bossWasActive {
		if bossActive {
			s.triggerBossIntro()
			s.switchToBossMusic()
		} else if s.bossWasActive {
			s.switchToLevelMusic()
		}
		s.bossWasActive = bossActive
	}
}

// triggerBossIntro triggers the boss intro overlay
func (s *PlayingScene) triggerBossIntro() {
	if s.bossIntroOverlay == nil {
		return
	}

	// Get current stage number
	stageNumber := 1
	if levelMgr := s.manager.GetLevelManager(); levelMgr != nil {
		stageNumber = levelMgr.GetLevel()
	}

	// Get boss type from stage config (simplified - use stage number to determine boss type)
	bossTypes := []string{"earth", "mars", "jupiter", "saturn", "uranus", "neptune"}
	bossType := "earth"
	if stageNumber > 0 && stageNumber <= len(bossTypes) {
		bossType = bossTypes[stageNumber-1]
	}

	s.bossIntroOverlay.Trigger(stageNumber, bossType)
}

// isBossActive checks if there's an active boss in the world
func (s *PlayingScene) isBossActive() bool {
	bossCount := 0
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
			filter.Contains(core.EnemyTypeID),
		),
	).Each(s.manager.GetWorld(), func(entry *donburi.Entry) {
		typeID := core.EnemyTypeID.Get(entry)
		if typeID != nil && *typeID == int(enemy.EnemyTypeBoss) {
			bossCount++
		}
	})
	return bossCount > 0
}

func (s *PlayingScene) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.Black)

	// Apply screen shake if active
	if s.screenShake > 0 {
		// Get image from pool instead of creating new one
		shakenImage := s.manager.GetImagePool().GetImage(screen.Bounds().Dx(), screen.Bounds().Dy())
		defer s.manager.GetImagePool().ReturnImage(shakenImage)

		s.drawGameContent(shakenImage)
		op := &ebiten.DrawImageOptions{}
		shakeOffset := s.screenShake * 5
		op.GeoM.Translate(shakeOffset, shakeOffset)
		screen.DrawImage(shakenImage, op)
	} else {
		s.drawGameContent(screen)
	}

	// Draw level title overlay if showing
	if s.showLevelTitle {
		s.drawLevelTitle(screen)
	}

	// Draw boss intro overlay if active (draws on top of everything)
	if s.bossIntroOverlay != nil && s.bossIntroOverlay.IsActive() {
		s.bossIntroOverlay.Draw(screen)
	}
}

func (s *PlayingScene) drawGameContent(screen *ebiten.Image) {
	// Use optimized render system if available
	if renderOptimizer := s.manager.GetRenderOptimizer(); renderOptimizer != nil {
		renderOptimizer.OptimizedRenderSystem(s.manager.GetWorld(), screen)
	} else {
		// Fallback to original render system
		renderWrapper := core.NewRenderSystemWrapper(screen)
		if err := renderWrapper.Update(s.manager.GetWorld()); err != nil {
			log.Printf("[ERROR] Render system failed: %v", err)
		}
	}

	if s.manager.GetConfig().Debug {
		s.drawDebugInfo(screen)
	}
}

// UI elements are now handled by the main game's EbitenUI system

// TriggerScreenShake triggers a screen shake effect
func (s *PlayingScene) TriggerScreenShake() {
	s.screenShake = 1.0 // Set shake intensity
}

// drawDebugInfo renders debug information when Config.Debug is enabled.
// Debug overlay (FPS, entity count, etc.) is drawn by the ECS debug renderer; this is a no-op placeholder.
func (s *PlayingScene) drawDebugInfo(screen *ebiten.Image) {}

func (s *PlayingScene) Enter() {
	// Show level title when entering playing scene
	if levelManager := s.manager.GetLevelManager(); levelManager != nil {
		s.ShowLevelTitle(levelManager.GetLevel())
	}

	// Start background music
	s.startBackgroundMusic()
	s.bossWasActive = false // Reset boss tracking
}

func (s *PlayingScene) Exit() {
	s.showLevelTitle = false

	// Stop background music
	s.stopBackgroundMusic()
}

// startBackgroundMusic starts playing the background music based on current level
func (s *PlayingScene) startBackgroundMusic() {
	audioPlayer := s.resourceMgr.GetAudioPlayer()
	if audioPlayer == nil {
		return
	}

	musicName := s.getLevelMusicName()
	musicRes, ok := s.resourceMgr.GetAudio(context.Background(), musicName)
	if !ok {
		return
	}

	if err := audioPlayer.PlayMusic(musicName, musicRes, 0.7); err != nil {
		log.Printf("[WARN] Failed to play background music: music=%s error=%v", musicName, err)
	}
}

// stopBackgroundMusic stops the background music
func (s *PlayingScene) stopBackgroundMusic() {
	audioPlayer := s.resourceMgr.GetAudioPlayer()
	if audioPlayer == nil {
		return
	}

	// Stop all gameplay music tracks
	audioPlayer.StopMusic(musicTrackLevel1)
	audioPlayer.StopMusic(musicTrackBoss)
	audioPlayer.StopMusic(musicTrackMain)
	s.currentMusicTrack = ""
}

// switchToBossMusic switches from level music to boss music
func (s *PlayingScene) switchToBossMusic() {
	audioPlayer := s.resourceMgr.GetAudioPlayer()
	if audioPlayer == nil {
		return
	}

	// Stop current level music
	audioPlayer.StopMusic(musicTrackLevel1)
	audioPlayer.StopMusic(musicTrackMain)

	musicRes, ok := s.resourceMgr.GetAudio(context.Background(), musicTrackBoss)
	if !ok {
		return
	}

	if err := audioPlayer.PlayMusic(musicTrackBoss, musicRes, 0.7); err != nil {
		log.Printf("[WARN] Failed to play boss music: error=%v", err)
	} else {
		s.currentMusicTrack = musicTrackBoss
	}
}

// switchToLevelMusic switches from boss music back to level music
func (s *PlayingScene) switchToLevelMusic() {
	audioPlayer := s.resourceMgr.GetAudioPlayer()
	if audioPlayer == nil {
		return
	}

	// Stop boss music
	audioPlayer.StopMusic(musicTrackBoss)

	// Determine which level music to play
	musicName := s.getLevelMusicName()

	musicRes, ok := s.resourceMgr.GetAudio(context.Background(), musicName)
	if !ok {
		return
	}

	if err := audioPlayer.PlayMusic(musicName, musicRes, 0.7); err != nil {
		log.Printf("[WARN] Failed to play level music: music=%s error=%v", musicName, err)
	} else {
		s.currentMusicTrack = musicName
	}
}

// getLevelMusicName determines which level music to play based on current level
func (s *PlayingScene) getLevelMusicName() string {
	levelManager := s.manager.GetLevelManager()
	if levelManager == nil {
		return musicTrackLevel1
	}

	// Use level-specific music (e.g., game_music_level_1 for level 1)
	if levelManager.GetLevel() == 1 {
		return musicTrackLevel1
	}
	if levelManager.GetLevel() == 2 {
		return musicTrackLevel2
	}

	// Fallback to main music for other levels (can be extended later)
	return musicTrackMain
}

// ShowLevelTitle displays the level title overlay
func (s *PlayingScene) ShowLevelTitle(levelNumber int) {
	s.currentLevelNumber = levelNumber
	s.levelTitleStartTime = time.Now()
	s.showLevelTitle = true
}

// drawLevelTitle draws the level title overlay
func (s *PlayingScene) drawLevelTitle(screen *ebiten.Image) {
	if s.font == nil {
		return
	}

	elapsed := time.Since(s.levelTitleStartTime).Seconds()
	alpha := s.calculateTitleAlpha(elapsed)
	if alpha <= 0 {
		return
	}

	s.drawTitleOverlay(screen, alpha)
	titleText, descText := s.getTitleText()
	s.drawTitleText(screen, titleText, descText, alpha)
}

// calculateTitleAlpha calculates fade alpha (fade in for first 0.5s, fade out for last 0.5s)
func (s *PlayingScene) calculateTitleAlpha(elapsed float64) float64 {
	alpha := 1.0
	if elapsed < 0.5 {
		alpha = elapsed / 0.5 // Fade in
	} else if elapsed > s.levelTitleDuration-0.5 {
		alpha = (s.levelTitleDuration - elapsed) / 0.5 // Fade out
	}
	return alpha
}

// drawTitleOverlay draws the semi-transparent background overlay
func (s *PlayingScene) drawTitleOverlay(screen *ebiten.Image, alpha float64) {
	bgColor := color.RGBA{0, 0, 0, uint8(200 * alpha)}
	overlay := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	overlay.Fill(bgColor)

	op := &ebiten.DrawImageOptions{}
	op.ColorScale.SetA(float32(alpha))
	screen.DrawImage(overlay, op)
}

// getTitleText gets title and description text for the current level
func (s *PlayingScene) getTitleText() (titleText, descText string) {
	titleText = fmt.Sprintf("STAGE %d", s.currentLevelNumber)
	return titleText, descText
}

// drawTitleText draws the title and description text
func (s *PlayingScene) drawTitleText(screen *ebiten.Image, titleText, descText string, alpha float64) {
	titleWidth, titleHeight := v2text.Measure(titleText, s.font, 0)
	screenWidth := float64(s.manager.GetConfig().ScreenSize.Width)
	screenHeight := float64(s.manager.GetConfig().ScreenSize.Height)

	titleX := (screenWidth - float64(titleWidth)) / 2
	titleY := screenHeight/2 - float64(titleHeight) - 30

	// Draw title
	titleOp := &v2text.DrawOptions{}
	titleOp.GeoM.Translate(titleX, titleY)
	titleOp.ColorScale.SetA(float32(alpha))
	v2text.Draw(screen, titleText, s.font, titleOp)

	// Draw description if available
	if descText != "" {
		descWidth, _ := v2text.Measure(descText, s.font, 0)
		descX := (screenWidth - float64(descWidth)) / 2
		descY := titleY + float64(titleHeight) + 20

		descOp := &v2text.DrawOptions{}
		descOp.GeoM.Translate(descX, descY)
		descOp.ColorScale.SetA(float32(alpha * 0.8)) // Slightly more transparent
		v2text.Draw(screen, descText, s.font, descOp)
	}
}

func (s *PlayingScene) GetType() scenes.SceneType {
	return scenes.ScenePlaying
}
