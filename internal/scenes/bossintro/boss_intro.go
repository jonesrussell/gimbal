package bossintro

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
	"github.com/jonesrussell/gimbal/internal/scenes/transitions"
)

const (
	bossIntroDuration = 1.5 // seconds
)

// BossIntroOverlay represents a boss intro overlay that appears on top of gameplay
type BossIntroOverlay struct {
	manager       *scenes.SceneManager
	font          text.Face
	resourceMgr   *resources.ResourceManager
	scoreManager  *managers.ScoreManager
	startTime     time.Time
	active        bool
	bossSprite    *ebiten.Image
	warningOverlay *ebiten.Image
	flashTransition *transitions.FlashTransition
	stageNumber   int
	bossType      string
	soundPlayed   bool
}

// NewBossIntroOverlay creates a new boss intro overlay
func NewBossIntroOverlay(
	manager *scenes.SceneManager,
	font text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) *BossIntroOverlay {
	config := manager.GetConfig()
	flashTransition := transitions.NewFlashTransition(0.3, 0.8, config.ScreenSize.Width, config.ScreenSize.Height)

	return &BossIntroOverlay{
		manager:        manager,
		font:           font,
		resourceMgr:    resourceMgr,
		scoreManager:   scoreManager,
		flashTransition: flashTransition,
		active:        false,
	}
}

// Trigger triggers the boss intro overlay
func (b *BossIntroOverlay) Trigger(stageNumber int, bossType string) {
	b.stageNumber = stageNumber
	b.bossType = bossType
	b.startTime = time.Now()
	b.active = true
	b.soundPlayed = false
	b.flashTransition.Reset()

	// Load boss sprite
	if b.resourceMgr != nil {
		// Try to load boss portrait
		if bossPortrait, ok := b.resourceMgr.GetSprite(context.Background(), fmt.Sprintf("boss_portrait_%s", bossType)); ok {
			b.bossSprite = bossPortrait
		} else if bossSprite, ok := b.resourceMgr.GetSprite(context.Background(), "enemy_boss"); ok {
			b.bossSprite = bossSprite
		}

		// Load warning overlay
		if warningOverlay, ok := b.resourceMgr.GetSprite(context.Background(), "warning_overlay"); ok {
			b.warningOverlay = warningOverlay
		}
	}
}

// Update updates the boss intro overlay
func (b *BossIntroOverlay) Update(deltaTime float64) {
	if !b.active {
		return
	}

	elapsed := time.Since(b.startTime).Seconds()

	// Update flash transition
	b.flashTransition.Update(deltaTime)

	// Auto-dismiss after duration
	if elapsed >= bossIntroDuration {
		b.active = false
		return
	}

	// Play warning sound at start
	if !b.soundPlayed && elapsed > 0.05 {
		b.playWarningSound()
		b.soundPlayed = true
	}
}

// Draw draws the boss intro overlay on top of the screen
func (b *BossIntroOverlay) Draw(screen *ebiten.Image) {
	if !b.active {
		return
	}

	config := b.manager.GetConfig()
	centerX := float64(config.ScreenSize.Width) / 2
	centerY := float64(config.ScreenSize.Height) / 2
	elapsed := time.Since(b.startTime).Seconds()

	// Draw red tint overlay
	if b.warningOverlay != nil {
		op := &ebiten.DrawImageOptions{}
		overlayAlpha := 0.6 * (1.0 - math.Min(1.0, elapsed/bossIntroDuration)) // Fade out
		op.ColorScale.SetA(float32(overlayAlpha))
		screen.DrawImage(b.warningOverlay, op)
	} else {
		// Fallback: draw red tint
		overlay := ebiten.NewImage(config.ScreenSize.Width, config.ScreenSize.Height)
		overlayAlpha := uint8(0.6 * 255 * (1.0 - math.Min(1.0, elapsed/bossIntroDuration)))
		overlay.Fill(color.RGBA{255, 0, 0, overlayAlpha})
		screen.DrawImage(overlay, nil)
	}

	// Draw flash effect
	if elapsed < 0.3 {
		b.flashTransition.Draw(screen, nil, nil)
	}

	// Draw boss sprite with zoom-in animation
	if b.bossSprite != nil {
		op := &ebiten.DrawImageOptions{}
		zoomProgress := math.Min(1.0, elapsed/0.8)
		// Ease-out zoom: start fast, slow down
		easedProgress := 1.0 - math.Pow(1.0-zoomProgress, 3)
		scale := 0.3 + easedProgress*1.2 // Scale from 0.3 to 1.5
		
		spriteWidth := float64(b.bossSprite.Bounds().Dx()) * scale
		spriteHeight := float64(b.bossSprite.Bounds().Dy()) * scale
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(centerX-spriteWidth/2, centerY-spriteHeight/2-40)
		
		// Fade in
		alpha := math.Min(1.0, elapsed/0.3)
		op.ColorScale.SetA(float32(alpha))
		screen.DrawImage(b.bossSprite, op)
	}

	// Draw warning text
	warningText := "WARNING"
	fadeAlpha := math.Min(1.0, elapsed/0.2)
	scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
		Text:  warningText,
		X:     centerX,
		Y:     centerY - 80,
		Alpha: fadeAlpha,
		Font:  b.font,
	})

	// Draw "BOSS APPROACHING" text
	bossText := "BOSS APPROACHING"
	scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
		Text:  bossText,
		X:     centerX,
		Y:     centerY + 100,
		Alpha: fadeAlpha,
		Font:  b.font,
	})
}

// IsActive returns whether the overlay is currently active
func (b *BossIntroOverlay) IsActive() bool {
	return b.active
}

// Dismiss manually dismisses the overlay
func (b *BossIntroOverlay) Dismiss() {
	b.active = false
}

func (b *BossIntroOverlay) playWarningSound() {
	if b.resourceMgr == nil {
		return
	}

	audioPlayer := b.resourceMgr.GetAudioPlayer()
	if audioPlayer == nil {
		return
	}

	// Play warning sound effect
	b.manager.GetLogger().Debug("Boss intro warning sound should play")
}
