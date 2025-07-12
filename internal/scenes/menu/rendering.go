package menu

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Draw renders the menu options
func (m *MenuSystem) Draw(screen *ebiten.Image, fadeAlpha float64) {
	for i, option := range m.options {
		y := m.config.MenuY + float64(i*m.config.ItemSpacing)
		isSelected := i == m.selection
		m.drawMenuOption(screen, option.Text, y, isSelected, fadeAlpha)
	}
}

// drawMenuOption renders a single menu option
func (m *MenuSystem) drawMenuOption(
	screen *ebiten.Image,
	option string,
	y float64,
	isSelected bool,
	fadeAlpha float64,
) {
	alpha := fadeAlpha
	scale := 1.0

	if isSelected {
		pulse := m.config.PulseBase + m.config.PulseAmplitude*math.Sin(m.animationTime*m.config.PulseSpeed)
		alpha *= pulse
		scale = 1.0 + 0.05*math.Sin(m.animationTime*m.config.PulseSpeed)

		m.drawSelectionBackground(screen, option, y, m.config.BackgroundAlpha*fadeAlpha)
		m.drawChevron(screen, y, fadeAlpha)
	} else {
		alpha *= m.config.DimmedAlpha
	}

	m.drawOptionText(screen, option, y, alpha, scale)
}

// drawSelectionBackground renders the background highlight for selected items
func (m *MenuSystem) drawSelectionBackground(screen *ebiten.Image, option string, y, alpha float64) {
	width, height := text.Measure(option, m.font, 0)
	w := int(width)
	h := int(height)

	bgColor := m.config.BackgroundColor
	bgColor.A = uint8(float64(bgColor.A) * alpha)

	rect := ebiten.NewImage(w+m.config.PaddingX*2, h+m.config.PaddingY*2)
	rect.Fill(bgColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		float64(m.screenWidth)/2-float64(w+m.config.PaddingX*2)/2,
		y-float64(h+m.config.PaddingY*2)/2+2,
	)
	screen.DrawImage(rect, op)
}

// drawChevron renders the animated selection chevron
func (m *MenuSystem) drawChevron(screen *ebiten.Image, y, fadeAlpha float64) {
	pulse := m.config.PulseBase + m.config.PulseAmplitude*math.Sin(m.animationTime*4.0)
	chevronX := float64(m.screenWidth)/2 - m.config.ChevronOffsetX
	chevronY := y + m.config.ChevronOffsetY

	op := &text.DrawOptions{}
	op.GeoM.Translate(chevronX, chevronY)
	op.ColorScale.SetR(float32(m.config.ChevronColor.R) / 255.0)
	op.ColorScale.SetG(float32(m.config.ChevronColor.G) / 255.0)
	op.ColorScale.SetB(float32(m.config.ChevronColor.B) / 255.0)
	op.ColorScale.SetA(float32(pulse * fadeAlpha))

	text.Draw(screen, ">", m.font, op)
}

// drawOptionText renders the text for a menu option
func (m *MenuSystem) drawOptionText(screen *ebiten.Image, option string, y, alpha, scale float64) {
	op := &text.DrawOptions{}

	width, height := text.Measure(option, m.font, 0)
	centerX := float64(m.screenWidth) / 2
	centerY := y

	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(
		centerX-width*scale/2,
		centerY-height*scale/2,
	)

	textColor := m.config.TextColor
	op.ColorScale.SetR(float32(textColor.R) / 255.0)
	op.ColorScale.SetG(float32(textColor.G) / 255.0)
	op.ColorScale.SetB(float32(textColor.B) / 255.0)
	op.ColorScale.SetA(float32(alpha))

	text.Draw(screen, option, m.font, op)
}
