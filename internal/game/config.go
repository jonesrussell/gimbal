package game

import "image"

const (
	defaultScreenWidth  = 640
	defaultScreenHeight = 480
	defaultPlayerSize   = 16
	defaultNumStars     = 100
	defaultSpeed        = 0.04
	defaultStarSize     = 5.0
	defaultStarSpeed    = 2.0
	defaultAngleStep    = 0.05
	defaultRadiusRatio  = 0.75
	screenCenterDivisor = 2
)

// GameConfig holds all configurable parameters for the game
type GameConfig struct {
	ScreenWidth  int
	ScreenHeight int
	PlayerWidth  int
	PlayerHeight int
	Radius       float64
	Center       image.Point
	NumStars     int
	Debug        bool
	Speed        float64
	StarSize     float64
	StarSpeed    float64
	AngleStep    float64
}

// GameOption is a function that modifies a GameConfig
type GameOption func(*GameConfig)

// WithScreenSize sets the screen dimensions
func WithScreenSize(width, height int) GameOption {
	return func(c *GameConfig) {
		c.ScreenWidth = width
		c.ScreenHeight = height
		c.Center = image.Point{X: width / screenCenterDivisor, Y: height / screenCenterDivisor}
		c.Radius = float64(height/screenCenterDivisor) * defaultRadiusRatio
	}
}

// WithPlayerSize sets the player dimensions
func WithPlayerSize(width, height int) GameOption {
	return func(c *GameConfig) {
		c.PlayerWidth = width
		c.PlayerHeight = height
	}
}

// WithNumStars sets the number of stars
func WithNumStars(num int) GameOption {
	return func(c *GameConfig) {
		c.NumStars = num
	}
}

// WithDebug enables debug mode
func WithDebug(debug bool) GameOption {
	return func(c *GameConfig) {
		c.Debug = debug
	}
}

// WithSpeed sets the game speed
func WithSpeed(speed float64) GameOption {
	return func(c *GameConfig) {
		c.Speed = speed
	}
}

// WithStarSettings sets star-related parameters
func WithStarSettings(size, speed float64) GameOption {
	return func(c *GameConfig) {
		c.StarSize = size
		c.StarSpeed = speed
	}
}

// WithAngleStep sets the angle step for player rotation
func WithAngleStep(step float64) GameOption {
	return func(c *GameConfig) {
		c.AngleStep = step
	}
}

// DefaultConfig returns a default game configuration
func DefaultConfig() *GameConfig {
	return &GameConfig{
		ScreenWidth:  defaultScreenWidth,
		ScreenHeight: defaultScreenHeight,
		PlayerWidth:  defaultPlayerSize,
		PlayerHeight: defaultPlayerSize,
		Radius:       float64(defaultScreenHeight/screenCenterDivisor) * defaultRadiusRatio,
		Center:       image.Point{X: defaultScreenWidth / screenCenterDivisor, Y: defaultScreenHeight / screenCenterDivisor},
		NumStars:     defaultNumStars,
		Debug:        false,
		Speed:        defaultSpeed,
		StarSize:     defaultStarSize,
		StarSpeed:    defaultStarSpeed,
		AngleStep:    defaultAngleStep,
	}
}

// NewConfig creates a new GameConfig with the given options
func NewConfig(opts ...GameOption) *GameConfig {
	config := DefaultConfig()
	for _, opt := range opts {
		opt(config)
	}
	return config
}
