package ui

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
	livesText  *widget.Text
	scoreText  *widget.Text
	healthBar  *widget.ProgressBar
	ammoIcons  []*widget.Graphic
	heartIcons []*widget.Graphic // Track heart icons separately

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
	// Ensure heartSprite is 32x32 for UI
	if heartSprite != nil {
		if heartSprite.Bounds().Dx() != HeartIconSize || heartSprite.Bounds().Dy() != HeartIconSize {
			scaled := ebiten.NewImage(HeartIconSize, HeartIconSize)
			op := &ebiten.DrawImageOptions{}
			scaleX := float64(HeartIconSize) / float64(heartSprite.Bounds().Dx())
			scaleY := float64(HeartIconSize) / float64(heartSprite.Bounds().Dy())
			op.GeoM.Scale(scaleX, scaleY)
			scaled.DrawImage(heartSprite, op)
			heartSprite = scaled
		}
	}
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
	// Main HUD row at top left
	hudRow := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(
			image.NewNineSliceColor(color.NRGBA{0, 0, 0, 180}), // semi-transparent black
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(24),
			widget.RowLayoutOpts.Padding(widget.Insets{Top: 12, Left: 16, Right: 16, Bottom: 12}),
		)),
	)

	// Lives: icon + value
	livesContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(4),
		)),
	)
	livesIcon := widget.NewGraphic(
		widget.GraphicOpts.Image(ui.heartSprite),
		widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.MinSize(24, 24)),
	)
	ui.livesText = widget.NewText(
		widget.TextOpts.Text("x3", ui.font, colornames.White),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.MinSize(32, 24)),
	)
	livesContainer.AddChild(livesIcon)
	livesContainer.AddChild(ui.livesText)

	// Score
	ui.scoreText = widget.NewText(
		widget.TextOpts.Text("Score: 0", ui.font, colornames.White),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.MinSize(120, 24)),
	)

	// Ammo: icon + value (or icons)
	ammoContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(4),
		)),
	)
	// Use up to 10 ammo icons, or a single icon + value for infinite
	ui.ammoIcons = ui.ammoIcons[:0]
	if ui.ammoSprite != nil {
		for i := 0; i < ui.currentAmmo && i < 10; i++ {
			ammoIcon := widget.NewGraphic(
				widget.GraphicOpts.Image(ui.ammoSprite),
				widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.MinSize(16, 16)),
			)
			ammoContainer.AddChild(ammoIcon)
			ui.ammoIcons = append(ui.ammoIcons, ammoIcon)
		}
	} else {
		// fallback: text
		ammoText := widget.NewText(
			widget.TextOpts.Text("Ammo: ∞", ui.font, colornames.White),
			widget.TextOpts.WidgetOpts(widget.WidgetOpts.MinSize(60, 24)),
		)
		ammoContainer.AddChild(ammoText)
	}

	// Add all sections to the HUD row
	hudRow.AddChild(livesContainer)
	hudRow.AddChild(ui.scoreText)
	hudRow.AddChild(ammoContainer)

	// Root container with anchor layout, HUD row anchored top left
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	root.AddChild(hudRow)

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
			widget.RowLayoutOpts.Spacing(5), // Reduced spacing
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

	// Don't call updateLivesIcons here - let it be called by UpdateLives
	// ui.updateLivesIcons(container)

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

	// Health bar fill with required TrackImage.Idle
	ui.healthBar = widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionStart,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ProgressBarOpts.Images(
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}), // Track (background)
				Hover: image.NewNineSliceColor(color.NRGBA{R: 120, G: 120, B: 120, A: 255}),
			},
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{R: 0, G: 255, B: 0, A: 255}), // Fill (foreground)
				Hover: image.NewNineSliceColor(color.NRGBA{R: 0, G: 200, B: 0, A: 255}),
			},
		),
		widget.ProgressBarOpts.Values(1, 3, 3),
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
	// Remove existing heart icons
	for _, heartIcon := range ui.heartIcons {
		container.RemoveChild(heartIcon)
	}
	ui.heartIcons = ui.heartIcons[:0] // Clear the slice

	// Add new heart icons
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
		ui.heartIcons = append(ui.heartIcons, heartIcon)
	}
}

// updateAmmoIcons updates the ammo icons based on current ammo
func (ui *ResponsiveUI) updateAmmoIcons(container *widget.Container) {
	// Remove existing ammo icons
	for _, ammoIcon := range ui.ammoIcons {
		container.RemoveChild(ammoIcon)
	}
	ui.ammoIcons = ui.ammoIcons[:0] // Clear the slice

	// Use ammo sprite if available, otherwise use heart sprite or create a simple colored image
	var ammoIconImage *ebiten.Image
	if ui.ammoSprite != nil {
		ammoIconImage = ui.ammoSprite
	} else if ui.heartSprite != nil {
		ammoIconImage = ui.heartSprite
	} else {
		// Create a simple yellow square for ammo icons
		ammoIconImage = ebiten.NewImage(16, 16)
		ammoIconImage.Fill(color.NRGBA{R: 255, G: 255, B: 0, A: 255}) // Yellow square
	}

	for i := 0; i < ui.currentAmmo && i < 10; i++ { // Limit to 10 visible icons
		ammoIcon := widget.NewGraphic(
			widget.GraphicOpts.Image(ammoIconImage),
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
