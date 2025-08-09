package core

import (
	"fmt"
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/colornames"
)

// HUDBuilder builds HUD components
type HUDBuilder struct {
	font          text.Face
	spriteManager *SpriteManager
}

// NewHUDBuilder creates a new HUD builder
func NewHUDBuilder(font text.Face, spriteManager *SpriteManager) *HUDBuilder {
	return &HUDBuilder{
		font:          font,
		spriteManager: spriteManager,
	}
}

// BuildHUD creates the main HUD container
func (hb *HUDBuilder) BuildHUD() *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(
			image.NewNineSliceColor(color.NRGBA{0, 0, 0, SemiTransparentBlack}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(HUDSpacing),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    HUDPaddingV,
				Left:   HUDPadding,
				Right:  HUDPadding,
				Bottom: HUDPaddingV,
			}),
		)),
	)
}

// BuildLivesContainer creates the lives display container
func (hb *HUDBuilder) BuildLivesContainer() *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(IconSpacing),
		)),
	)
}

// BuildLivesIcon creates a heart icon for lives display
func (hb *HUDBuilder) BuildLivesIcon() *widget.Graphic {
	return widget.NewGraphic(
		widget.GraphicOpts.Image(hb.spriteManager.GetHeartSprite()),
		widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.MinSize(HeartIconSize, HeartIconSize)),
	)
}

// BuildLivesText creates the lives text widget
func (hb *HUDBuilder) BuildLivesText(lives int) *widget.Text {
	return widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("x%d", lives), hb.font, colornames.White),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.MinSize(32, 24)),
	)
}

// BuildScoreText creates the score text widget
func (hb *HUDBuilder) BuildScoreText(score int) *widget.Text {
	return widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Score: %d", score), hb.font, colornames.White),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.MinSize(120, 24)),
	)
}

// BuildAmmoContainer creates the ammo display container
func (hb *HUDBuilder) BuildAmmoContainer() *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(IconSpacing),
		)),
	)
}

// BuildAmmoIcon creates an ammo icon
func (hb *HUDBuilder) BuildAmmoIcon() *widget.Graphic {
	return widget.NewGraphic(
		widget.GraphicOpts.Image(hb.spriteManager.GetAmmoSprite()),
		widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.MinSize(AmmoIconSize, AmmoIconSize)),
	)
}

// BuildAmmoText creates fallback ammo text
func (hb *HUDBuilder) BuildAmmoText() *widget.Text {
	return widget.NewText(
		widget.TextOpts.Text("Ammo: ∞", hb.font, colornames.White),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.MinSize(60, 24)),
	)
}
