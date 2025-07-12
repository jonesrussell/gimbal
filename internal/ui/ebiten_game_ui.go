package ui

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/internal/common"
)

type EbitenGameUI struct {
	responsiveUI *ResponsiveUI
	paused       bool
	deviceClass  string
}

func NewEbitenGameUI(responsiveUI *ResponsiveUI) *EbitenGameUI {
	return &EbitenGameUI{
		responsiveUI: responsiveUI,
	}
}

func (e *EbitenGameUI) Update() error {
	e.responsiveUI.Update()
	return nil
}

func (e *EbitenGameUI) Draw(screen *ebiten.Image) {
	e.responsiveUI.Draw(screen)
}

func (e *EbitenGameUI) UpdateScore(score int) {
	e.responsiveUI.UpdateScore(score)
}

func (e *EbitenGameUI) UpdateLives(lives int) {
	e.responsiveUI.UpdateLives(lives)
}

func (e *EbitenGameUI) ShowPauseMenu(visible bool) {
	e.paused = visible
	// TODO: Implement pause menu overlay in ResponsiveUI
}

func (e *EbitenGameUI) SetDeviceClass(deviceClass string) {
	e.deviceClass = deviceClass
	e.responsiveUI.UpdateResponsiveLayout(e.responsiveUI.screenWidth, e.responsiveUI.screenHeight)
}

// Data-driven update for HUD
func (e *EbitenGameUI) UpdateHUD(data common.HUDData) {
	e.UpdateScore(data.Score)
	e.UpdateLives(data.Lives)
	// TODO: Add health, level, etc.
}
