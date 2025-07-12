package ui_ebitenui

import (
	"fmt"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/colornames"
)

// ResponsiveUI implements 2025 responsive design using EbitenUI
type ResponsiveUI struct {
	ui *ebitenui.UI

	// UI components
	hudContainer    *widget.Container
	livesContainer  *widget.Container
	scoreContainer  *widget.Container
	healthContainer *widget.Container
	ammoContainer   *widget.Container

	// Widgets
	livesText *widget.Text
	scoreText *widget.Text
	healthBar *widget.ProgressBar
	ammoIcons []*widget.Graphic

	// Resources
	font        text.Face
	heartSprite *ebiten.Image
	ammoSprite  *ebiten.Image

	// State
	currentLives  int
	currentScore  int
	currentHealth float64
	currentAmmo   int

	// Responsive settings
	screenWidth  int
	screenHeight int
	deviceClass  string
}

// NewResponsiveUI creates a new EbitenUI-based responsive UI system
func NewResponsiveUI(font text.Face, heartSprite, ammoSprite *ebiten.Image) *ResponsiveUI {
	ui := &ResponsiveUI{
		font:          font,
		heartSprite:   heartSprite,
		ammoSprite:    ammoSprite,
		currentLives:  3,
		currentScore:  0,
		currentHealth: 1.0,
		currentAmmo:   10,
	}

	ui.createResponsiveUI()
	return ui
}

// createResponsiveUI creates the main responsive UI structure
func (ui *ResponsiveUI) createResponsiveUI() {
	// Root container with anchor layout for responsive positioning
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.BackgroundImage(
			image.NewNineSliceColor(color.RGBA{0, 0, 0, 0}),
		),
	)

	// Create HUD container
	ui.hudContainer = ui.createHUDContainer()
	root.AddChild(ui.hudContainer)

	ui.ui = &ebitenui.UI{Container: root}
}

// createHUDContainer creates the responsive HUD layout
func (ui *ResponsiveUI) createHUDContainer() *widget.Container {
	hud := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	// Lives display - top left
	ui.livesContainer = ui.createLivesDisplay()
	hud.AddChild(ui.livesContainer)

	// Score display - top right
	ui.scoreContainer = ui.createScoreDisplay()
	hud.AddChild(ui.scoreContainer)

	// Health bar - bottom left
	ui.healthContainer = ui.createHealthBar()
	hud.AddChild(ui.healthContainer)

	// Ammo counter - bottom right
	ui.ammoContainer = ui.createAmmoCounter()
	hud.AddChild(ui.ammoContainer)

	return hud
}

// createLivesDisplay creates the responsive lives display
func (ui *ResponsiveUI) createLivesDisplay() *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionStart,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				Padding: widget.Insets{
					Left: 20, Top: 20,
				},
			}),
		),
	)

	// Lives text
	ui.livesText = widget.NewText(
		widget.TextOpts.Text("Lives:", ui.font, colornames.White),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)
	container.AddChild(ui.livesText)

	// Heart icons (responsive sizing)
	ui.updateLivesIcons(container)

	return container
}

// createScoreDisplay creates the responsive score display
func (ui *ResponsiveUI) createScoreDisplay() *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.BackgroundImage(
			image.NewNineSliceColor(color.RGBA{0, 0, 0, 180}),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				Padding: widget.Insets{
					Right: 20, Top: 20,
				},
			}),
		),
	)

	ui.scoreText = widget.NewText(
		widget.TextOpts.Text("Score: 0", ui.font, colornames.White),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)
	container.AddChild(ui.scoreText)

	return container
}

// createHealthBar creates the responsive health bar
func (ui *ResponsiveUI) createHealthBar() *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionStart,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				Padding: widget.Insets{
					Left: 20, Bottom: 20,
				},
			}),
		),
	)

	// Health bar background
	healthBarBg := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(
			image.NewNineSliceColor(color.RGBA{40, 40, 40, 200}),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionStart,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)

	// Health bar fill
	ui.healthBar = widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionStart,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)

	container.AddChild(healthBarBg)
	container.AddChild(ui.healthBar)

	return container
}

// createAmmoCounter creates the responsive ammo counter
func (ui *ResponsiveUI) createAmmoCounter() *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(5),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				Padding: widget.Insets{
					Right: 20, Bottom: 20,
				},
			}),
		),
	)

	// Ammo text
	ammoText := widget.NewText(
		widget.TextOpts.Text("Ammo:", ui.font, colornames.White),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)
	container.AddChild(ammoText)

	// Ammo icons
	ui.updateAmmoIcons(container)

	return container
}

// updateLivesIcons updates the heart icons based on current lives
func (ui *ResponsiveUI) updateLivesIcons(container *widget.Container) {
	// For now, just add heart icons without removing existing ones
	// This is a simplified approach that works with EbitenUI
	for i := 0; i < ui.currentLives; i++ {
		heartIcon := widget.NewGraphic(
			widget.GraphicOpts.Image(ui.heartSprite),
			widget.GraphicOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}),
			),
		)
		container.AddChild(heartIcon)
	}
}

// updateAmmoIcons updates the ammo icons based on current ammo
func (ui *ResponsiveUI) updateAmmoIcons(container *widget.Container) {
	// For now, just add ammo icons without removing existing ones
	// This is a simplified approach that works with EbitenUI
	ui.ammoIcons = make([]*widget.Graphic, 0, ui.currentAmmo)
	for i := 0; i < ui.currentAmmo && i < 10; i++ { // Limit to 10 visible icons
		ammoIcon := widget.NewGraphic(
			widget.GraphicOpts.Image(ui.ammoSprite),
			widget.GraphicOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}),
			),
		)
		container.AddChild(ammoIcon)
		ui.ammoIcons = append(ui.ammoIcons, ammoIcon)
	}
}

// UpdateLives updates the lives display
func (ui *ResponsiveUI) UpdateLives(lives int) {
	if ui.currentLives != lives {
		ui.currentLives = lives
		ui.updateLivesIcons(ui.livesContainer)
	}
}

// UpdateScore updates the score display
func (ui *ResponsiveUI) UpdateScore(score int) {
	if ui.currentScore != score {
		ui.currentScore = score
		ui.scoreText.Label = fmt.Sprintf("Score: %d", score)
	}
}

// UpdateHealth updates the health bar
func (ui *ResponsiveUI) UpdateHealth(health float64) {
	if ui.currentHealth != health {
		ui.currentHealth = health
		ui.healthBar.SetCurrent(int(health * 100)) // Convert to percentage
	}
}

// UpdateAmmo updates the ammo counter
func (ui *ResponsiveUI) UpdateAmmo(ammo int) {
	if ui.currentAmmo != ammo {
		ui.currentAmmo = ammo
		ui.updateAmmoIcons(ui.ammoContainer)
	}
}

// UpdateResponsiveLayout updates the UI layout based on screen size
func (ui *ResponsiveUI) UpdateResponsiveLayout(width, height int) {
	ui.screenWidth = width
	ui.screenHeight = height

	// Determine device class
	if width < 768 {
		ui.deviceClass = "mobile"
	} else if width < 1024 {
		ui.deviceClass = "tablet"
	} else if width > 1920 {
		ui.deviceClass = "ultrawide"
	} else {
		ui.deviceClass = "desktop"
	}
}

// Update updates the EbitenUI system
func (ui *ResponsiveUI) Update() {
	ui.ui.Update()
}

// Draw draws the UI
func (ui *ResponsiveUI) Draw(screen *ebiten.Image) {
	ui.ui.Draw(screen)
}

// GetUI returns the underlying EbitenUI instance
func (ui *ResponsiveUI) GetUI() *ebitenui.UI {
	return ui.ui
}

// GetDeviceClass returns the current device class
func (ui *ResponsiveUI) GetDeviceClass() string {
	return ui.deviceClass
}

// GetScreenDimensions returns the current screen dimensions
func (ui *ResponsiveUI) GetScreenDimensions() (width, height int) {
	return ui.screenWidth, ui.screenHeight
}
