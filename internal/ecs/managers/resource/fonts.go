package resources

import (
	"context"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/jonesrussell/gimbal/assets"
	"github.com/jonesrussell/gimbal/internal/errors"
)

// loadDefaultFont loads the default game font
func (rm *ResourceManager) loadDefaultFont(ctx context.Context) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

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
func (rm *ResourceManager) GetDefaultFont(ctx context.Context) (text.Face, error) {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if rm.defaultFont == nil {
		// Try to load the font if not already loaded
		if err := rm.loadDefaultFont(ctx); err != nil {
			return nil, err
		}
	}
	if rm.defaultFont == nil {
		return nil, errors.NewGameError(errors.AssetLoadFailed, "default font not loaded")
	}
	return rm.defaultFont, nil
}
