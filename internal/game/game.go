package game

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/logger"
	"github.com/solarlune/resolv"
)

//go:embed assets/*
var assets embed.FS

type GimlarGame struct {
	config      *GameConfig
	player      *Player
	stars       []Star
	space       *resolv.Space
	prevX       float64
	prevY       float64
	input       InputHandlerInterface
	starImage   *ebiten.Image
	gameStarted bool
}

func NewGimlarGame(config *GameConfig, input InputHandlerInterface) (*GimlarGame, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if input == nil {
		return nil, errors.New("input handler cannot be nil")
	}

	// Create a single star image that will be used for all stars
	starImage := ebiten.NewImage(1, 1)
	starImage.Fill(color.White)

	g := &GimlarGame{
		config:      config,
		input:       input,
		prevX:       0,
		prevY:       0,
		starImage:   starImage,
		gameStarted: false,
	}

	// Initialize stars
	g.stars = initializeStars(config.NumStars, starImage)

	// Load the player sprite
	imageData, err := assets.ReadFile("assets/player.png")
	if err != nil {
		return nil, fmt.Errorf("failed to load player image: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode player image: %w", err)
	}

	g.player, err = NewPlayer(input, config, img)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	g.space = resolv.NewSpace(config.ScreenWidth, config.ScreenHeight, config.PlayerWidth, config.PlayerHeight)
	g.space.Add(g.player.Object)

	return g, nil
}

func (g *GimlarGame) Run() error {
	ebiten.SetWindowSize(g.config.ScreenWidth, g.config.ScreenHeight)
	return ebiten.RunGame(g)
}

func (g *GimlarGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.config.ScreenWidth, g.config.ScreenHeight
}

func (g *GimlarGame) Update() error {
	// Update the stars
	g.updateStars()

	// Update the player's state
	g.player.Update()
	g.player.updatePosition()

	// Log the player's position after updating if it has changed
	pos := g.player.Object.Position()
	if pos.X != g.prevX || pos.Y != g.prevY {
		logger.GlobalLogger.Debug("Player position after update", "X", pos.X, "Y", pos.Y)
		g.prevX = pos.X
		g.prevY = pos.Y
	}

	return nil
}

func (g *GimlarGame) Draw(screen *ebiten.Image) {
	// Draw the stars
	g.drawStars(screen)

	// Draw the player
	g.drawPlayer(screen)

	// Draw debug info if debug is true
	if g.config.Debug {
		g.DrawDebugInfo(screen)
	}
}

func (g *GimlarGame) drawPlayer(screen *ebiten.Image) {
	g.player.Draw(screen)
}

func (g *GimlarGame) GetRadius() float64 {
	return g.config.Radius
}

// GetSpeed returns the game's speed
func (g *GimlarGame) GetSpeed() float64 {
	return g.config.Speed
}

// GetPlayer returns the game's player
func (g *GimlarGame) GetPlayer() *Player {
	return g.player
}

// GetStars returns the game's stars
func (g *GimlarGame) GetStars() []Star {
	return g.stars
}

// GetSpace returns the game's space
func (g *GimlarGame) GetSpace() *resolv.Space {
	return g.space
}
