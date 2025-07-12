package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	textv2 "github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/common"
)

type UIConfig struct {
	Font  interface{} // Should be textv2.Face
	Theme interface{} // Should be *ebiten.Image
}

type UIFactory interface {
	CreateGameUI(config UIConfig) common.GameUI
}

type EbitenUIFactory struct{}

func (f *EbitenUIFactory) CreateGameUI(config UIConfig) common.GameUI {
	font, ok := config.Font.(textv2.Face)
	if !ok {
		panic("UIConfig.Font must be textv2.Face")
	}
	heartSprite, ok := config.Theme.(*ebiten.Image)
	if !ok {
		panic("UIConfig.Theme must be *ebiten.Image")
	}
	responsiveUI := NewResponsiveUI(font, heartSprite, nil)
	return NewEbitenGameUI(responsiveUI)
}
