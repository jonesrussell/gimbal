package ecs

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Renderer handles all drawing operations
type Renderer struct {
	colors MenuColors
	layout MenuLayout
	calc   *AnimationCalculator
}

// NewRenderer creates a new renderer
func NewRenderer(colors MenuColors, layout MenuLayout, calc *AnimationCalculator) *Renderer {
	return &Renderer{
		colors: colors,
		layout: layout,
		calc:   calc,
	}
}

// DrawPauseMenu renders the entire pause menu
func (r *Renderer) DrawPauseMenu(screen *ebiten.Image, width, height, selection int, items []MenuItem) {
	r.drawOverlay(screen, width, height)
	r.drawTitle(screen, width)
	r.drawMenuItems(screen, width, selection, items)
	r.drawHintText(screen, width, height)
}

func (r *Renderer) drawOverlay(screen *ebiten.Image, width, height int) {
	alpha := uint8(128 * r.calc.GetFadeInAlpha())
	overlay := ebiten.NewImage(width, height)
	overlay.Fill(color.RGBA{r.colors.Overlay.R, r.colors.Overlay.G, r.colors.Overlay.B, alpha})
	screen.DrawImage(overlay, &ebiten.DrawImageOptions{})
}

func (r *Renderer) drawTitle(screen *ebiten.Image, width int) {
	alpha := r.calc.GetFadeInAlpha() * r.calc.GetTitlePulse()

	op := &text.DrawOptions{}
	op.GeoM.Scale(r.layout.TitleScale, r.layout.TitleScale)
	op.GeoM.Translate(float64(width)/2-75, r.layout.TitleY)
	op.ColorScale.SetR(float32(r.colors.Title.R) / 255)
	op.ColorScale.SetG(float32(r.colors.Title.G) / 255)
	op.ColorScale.SetB(float32(r.colors.Title.B) / 255)
	op.ColorScale.SetA(float32(alpha))

	text.Draw(screen, "PAUSED", defaultFontFace, op)
}

func (r *Renderer) drawMenuItems(screen *ebiten.Image, width, selection int, items []MenuItem) {
	for i, item := range items {
		itemY := r.layout.MenuStartY + float64(i)*r.layout.MenuItemSpacing
		isSelected := i == selection

		if isSelected {
			r.drawSelectedItem(screen, width, item.Text, itemY)
		} else {
			r.drawUnselectedItem(screen, width, item.Text, itemY)
		}
	}
}

func (r *Renderer) drawSelectedItem(screen *ebiten.Image, width int, text string, y float64) {
	// Draw chevron
	chevronX := r.calc.GetChevronPosition(float64(width)/2 + r.layout.ChevronOffsetX)
	chevronAlpha := r.calc.GetPulseValue(4.0) * r.calc.GetFadeInAlpha()
	r.drawChevron(screen, chevronX, y, chevronAlpha)

	// Draw background
	bgAlpha := 0.6 * r.calc.GetFadeInAlpha()
	r.drawSelectionBackground(screen, width, text, y, bgAlpha)

	// Draw text with scaling
	scale := r.calc.GetScaleValue()
	alpha := r.calc.GetFadeInAlpha() * r.calc.GetPulseValue(r.calc.config.PulseSpeed)
	r.drawText(screen, width, text, y, alpha, scale, r.colors.SelectedText)
}

func (r *Renderer) drawUnselectedItem(screen *ebiten.Image, width int, text string, y float64) {
	alpha := r.calc.GetFadeInAlpha() * 0.7
	r.drawText(screen, width, text, y, alpha, 1.0, r.colors.Text)
}

func (r *Renderer) drawChevron(screen *ebiten.Image, x, y, alpha float64) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y-5)
	op.ColorScale.SetR(float32(r.colors.Chevron.R) / 255)
	op.ColorScale.SetG(float32(r.colors.Chevron.G) / 255)
	op.ColorScale.SetB(float32(r.colors.Chevron.B) / 255)
	op.ColorScale.SetA(float32(alpha))
	text.Draw(screen, ">", defaultFontFace, op)
}

func (r *Renderer) drawSelectionBackground(screen *ebiten.Image, width int, itemText string, y, alpha float64) {
	textWidth, textHeight := text.Measure(itemText, defaultFontFace, 0)
	w := int(textWidth)
	h := int(textHeight)

	// Draw border
	r.drawBackgroundLayer(screen, width, w, h, y, alpha, r.colors.Border, r.layout.BorderWidth)

	// Draw main background
	r.drawBackgroundLayer(screen, width, w, h, y, alpha, r.colors.Background, 0)
}

func (r *Renderer) drawBackgroundLayer(screen *ebiten.Image, width, w, h int, y, alpha float64, col color.RGBA, extraSize int) {
	rectW := w + r.layout.PaddingX*2 + extraSize*2
	rectH := h + r.layout.PaddingY*2 + extraSize*2

	rect := ebiten.NewImage(rectW, rectH)
	layerColor := color.RGBA{col.R, col.G, col.B, uint8(float64(col.A) * alpha)}
	rect.Fill(layerColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		float64(width)/2-float64(rectW)/2,
		y-float64(rectH)/2,
	)
	screen.DrawImage(rect, op)
}

func (r *Renderer) drawText(screen *ebiten.Image, width int, textStr string, y, alpha, scale float64, col color.RGBA) {
	textWidth, textHeight := text.Measure(textStr, defaultFontFace, 0)

	op := &text.DrawOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(
		float64(width)/2-textWidth*scale/2,
		y-textHeight*scale/2,
	)
	op.ColorScale.SetR(float32(col.R) / 255)
	op.ColorScale.SetG(float32(col.G) / 255)
	op.ColorScale.SetB(float32(col.B) / 255)
	op.ColorScale.SetA(float32(alpha))

	text.Draw(screen, textStr, defaultFontFace, op)
}

func (r *Renderer) drawHintText(screen *ebiten.Image, width, height int) {
	hintText := "Press ESC to resume quickly"
	alpha := 0.6 * r.calc.GetFadeInAlpha() * r.calc.GetHintPulse()

	textWidth, _ := text.Measure(hintText, defaultFontFace, 0)

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(width)/2-textWidth/2, float64(height)-60)
	op.ColorScale.SetR(float32(r.colors.HintText.R) / 255)
	op.ColorScale.SetG(float32(r.colors.HintText.G) / 255)
	op.ColorScale.SetB(float32(r.colors.HintText.B) / 255)
	op.ColorScale.SetA(float32(alpha))

	text.Draw(screen, hintText, defaultFontFace, op)
}
