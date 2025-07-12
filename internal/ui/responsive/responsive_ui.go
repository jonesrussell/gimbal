package ui

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/internal/ui/core"
	"github.com/jonesrussell/gimbal/internal/ui/rendering"
	"github.com/jonesrussell/gimbal/internal/ui/state"
)

// ResponsiveUI implements a clean, responsive UI system
type ResponsiveUI struct {
	ui            *ebitenui.UI
	state         *state.State
	layout        *rendering.Layout
	hudBuilder    *core.HUDBuilder
	spriteManager *core.SpriteManager

	// UI components
	hudContainer   *widget.Container
	livesContainer *widget.Container
	ammoContainer  *widget.Container
	livesText      *widget.Text
	scoreText      *widget.Text

	// Dynamic widgets
	heartIcons []*widget.Graphic
	ammoIcons  []*widget.Graphic
}

// NewResponsiveUI creates a new responsive UI system with the provided configuration
func NewResponsiveUI(config *Config) (*ResponsiveUI, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	spriteManager := core.NewSpriteManager(config.HeartSprite, config.AmmoSprite)
	hudBuilder := core.NewHUDBuilder(config.Font, spriteManager)

	ui := &ResponsiveUI{
		state:         state.NewState(),
		layout:        rendering.NewLayout(),
		hudBuilder:    hudBuilder,
		spriteManager: spriteManager,
		heartIcons:    make([]*widget.Graphic, 0),
		ammoIcons:     make([]*widget.Graphic, 0),
	}

	if err := ui.build(); err != nil {
		return nil, err
	}

	return ui, nil
}

// build constructs the UI components
func (ui *ResponsiveUI) build() error {
	ui.hudContainer = ui.hudBuilder.BuildHUD()
	ui.livesContainer = ui.hudBuilder.BuildLivesContainer()
	ui.ammoContainer = ui.hudBuilder.BuildAmmoContainer()

	// Build initial components
	ui.rebuildLivesDisplay()
	ui.rebuildScoreDisplay()
	ui.rebuildAmmoDisplay()

	// Assemble HUD
	ui.hudContainer.AddChild(ui.livesContainer)
	ui.hudContainer.AddChild(ui.scoreText)
	ui.hudContainer.AddChild(ui.ammoContainer)

	// Create root container
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	root.AddChild(ui.hudContainer)

	ui.ui = &ebitenui.UI{Container: root}
	return nil
}

// rebuildLivesDisplay rebuilds the lives display
func (ui *ResponsiveUI) rebuildLivesDisplay() {
	ui.clearContainer(ui.livesContainer, ui.heartIcons)
	ui.heartIcons = nil

	// Add heart icon
	heartIcon := ui.hudBuilder.BuildLivesIcon()
	ui.livesContainer.AddChild(heartIcon)
	ui.heartIcons = append(ui.heartIcons, heartIcon)

	// Add lives text
	ui.livesText = ui.hudBuilder.BuildLivesText(ui.state.Lives)
	ui.livesContainer.AddChild(ui.livesText)
}

// rebuildScoreDisplay rebuilds the score display
func (ui *ResponsiveUI) rebuildScoreDisplay() {
	ui.scoreText = ui.hudBuilder.BuildScoreText(ui.state.Score)
}

// rebuildAmmoDisplay rebuilds the ammo display
func (ui *ResponsiveUI) rebuildAmmoDisplay() {
	ui.clearContainer(ui.ammoContainer, ui.ammoIcons)
	ui.ammoIcons = nil

	if ui.state.Ammo <= core.MaxAmmoIcons {
		// Show individual ammo icons
		for i := 0; i < ui.state.Ammo; i++ {
			ammoIcon := ui.hudBuilder.BuildAmmoIcon()
			ui.ammoContainer.AddChild(ammoIcon)
			ui.ammoIcons = append(ui.ammoIcons, ammoIcon)
		}
	} else {
		// Show fallback text for high ammo counts
		ammoText := ui.hudBuilder.BuildAmmoText()
		ui.ammoContainer.AddChild(ammoText)
	}
}

// clearContainer removes widgets from container and clears the slice
func (ui *ResponsiveUI) clearContainer(container *widget.Container, widgets []*widget.Graphic) {
	for _, widget := range widgets {
		container.RemoveChild(widget)
	}
}

// UpdateLives updates the lives display
func (ui *ResponsiveUI) UpdateLives(lives int) {
	newState := *ui.state
	newState.Lives = lives
	if err := newState.Validate(); err != nil {
		return // Optionally log
	}

	if ui.state.Lives != lives {
		ui.state.Lives = lives
		ui.livesText.Label = ui.hudBuilder.BuildLivesText(lives).Label
	}
}

// UpdateScore updates the score display
func (ui *ResponsiveUI) UpdateScore(score int) {
	newState := *ui.state
	newState.Score = score
	if err := newState.Validate(); err != nil {
		return // Optionally log
	}

	if ui.state.Score != score {
		ui.state.Score = score
		ui.scoreText.Label = ui.hudBuilder.BuildScoreText(score).Label
	}
}

// UpdateAmmo updates the ammo display
func (ui *ResponsiveUI) UpdateAmmo(ammo int) {
	newState := *ui.state
	newState.Ammo = ammo
	if err := newState.Validate(); err != nil {
		return // Optionally log
	}

	if ui.state.Ammo != ammo {
		ui.state.Ammo = ammo
		ui.rebuildAmmoDisplay()
	}
}

// UpdateResponsiveLayout updates the layout based on screen dimensions
func (ui *ResponsiveUI) UpdateResponsiveLayout(width, height int) {
	ui.layout.Update(width, height)
}

// Update updates the UI system
func (ui *ResponsiveUI) Update() error {
	ui.ui.Update()
	return nil
}

// Draw renders the UI
func (ui *ResponsiveUI) Draw(screen *ebiten.Image) {
	ui.ui.Draw(screen)
}

// GetDeviceClass returns the current device class
func (ui *ResponsiveUI) GetDeviceClass() string {
	return ui.layout.GetDeviceClass()
}

// GetScreenDimensions returns the current screen dimensions
func (ui *ResponsiveUI) GetScreenDimensions() (width, height int) {
	return ui.layout.GetDimensions()
}

func (ui *ResponsiveUI) SetDeviceClass(deviceClass string) {
	// Map deviceClass string to a width for layout.Update
	switch deviceClass {
	case core.DeviceMobile:
		ui.layout.Update(core.MobileBreakpoint-1, 600)
	case core.DeviceTablet:
		ui.layout.Update(core.TabletBreakpoint-1, 800)
	case core.DeviceUltrawide:
		ui.layout.Update(core.UltrawideBreakpoint+1, 1080)
	default:
		ui.layout.Update(1280, 720)
	}
}

func (ui *ResponsiveUI) ShowPauseMenu(visible bool) {
	// TODO: Implement pause menu display logic if needed
}
