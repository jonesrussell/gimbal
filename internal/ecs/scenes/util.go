package ecs

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/jonesrussell/gimbal/assets"
)

var defaultFontFace text.Face

func init() {
	fontBytes, err := assets.Assets.ReadFile("fonts/PressStart2P.ttf")
	if err != nil {
		log.Fatalf("failed to read font: %v", err)
	}
	fontTTF, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatalf("failed to parse font: %v", err)
	}

	// Create font face for opentype
	opentypeFace, err := opentype.NewFace(fontTTF, &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("failed to create opentype face: %v", err)
	}

	// Create text/v2 face from opentype face
	defaultFontFace = text.NewGoXFace(opentypeFace)
}

// drawCenteredText draws text centered on screen (helper method for scenes)
func drawCenteredText(screen *ebiten.Image, textStr string, x, y, alpha float64) {
	width, height := text.Measure(textStr, defaultFontFace, 0)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(int(x)-int(width)/2), float64(int(y)+int(height)/2))
	op.ColorScale.SetR(1)
	op.ColorScale.SetG(1)
	op.ColorScale.SetB(1)
	op.ColorScale.SetA(float32(alpha))
	text.Draw(screen, textStr, defaultFontFace, op)
}
