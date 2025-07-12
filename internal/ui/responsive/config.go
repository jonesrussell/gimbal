package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/ui/core"
)

// Config holds all UI configuration
type Config struct {
	Font        text.Face
	HeartSprite *ebiten.Image
	AmmoSprite  *ebiten.Image
}

// Validate ensures the configuration is valid
func (c *Config) Validate() error {
	if c.Font == nil {
		return core.ErrInvalidFont
	}
	if c.HeartSprite == nil {
		return core.ErrInvalidHeartSprite
	}
	return nil
}
