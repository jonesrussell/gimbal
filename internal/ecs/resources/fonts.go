package resources

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/jonesrussell/gimbal/assets"
)

// loadDefaultFont loads the default game font
func (rm *ResourceManager) loadDefaultFont() error {
	fontBytes, err := assets.Assets.ReadFile("fonts/PressStart2P.ttf")
	if err != nil {
		rm.logger.Error("failed to read font", "error", err)
		return err
	}
	fontTTF, err := opentype.Parse(fontBytes)
	if err != nil {
		rm.logger.Error("failed to parse font", "error", err)
		return err
	}
	otFace, err := opentype.NewFace(fontTTF, &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		rm.logger.Error("failed to create opentype face", "error", err)
		return err
	}
	rm.defaultFont = text.NewGoXFace(otFace)
	return nil
}

// GetDefaultFont returns the default game font
func (rm *ResourceManager) GetDefaultFont() text.Face {
	return rm.defaultFont
}
