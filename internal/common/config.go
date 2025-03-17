package common

const (
	DefaultScreenWidth  = 640
	DefaultScreenHeight = 480
	DefaultPlayerSize   = 16
	DefaultNumStars     = 100
	DefaultSpeed        = 0.04
	DefaultStarSize     = 5.0
	DefaultStarSpeed    = 2.0
	DefaultAngleStep    = 0.05
	DefaultRadiusRatio  = 0.75
)

// GameConfig holds all configurable parameters for the game
type GameConfig struct {
	ScreenSize Size
	PlayerSize Size
	Radius     float64
	NumStars   int
	Debug      bool
	Speed      float64
	StarSize   float64
	StarSpeed  float64
	AngleStep  float64
}

// GameOption is a function that modifies a GameConfig
type GameOption func(*GameConfig)

// WithScreenSize sets the screen dimensions
func WithScreenSize(width, height int) GameOption {
	return func(c *GameConfig) {
		c.ScreenSize = Size{Width: width, Height: height}
		c.Radius = float64(height/2) * DefaultRadiusRatio
	}
}

// WithPlayerSize sets the player dimensions
func WithPlayerSize(width, height int) GameOption {
	return func(c *GameConfig) {
		c.PlayerSize = Size{Width: width, Height: height}
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
		ScreenSize: Size{Width: DefaultScreenWidth, Height: DefaultScreenHeight},
		PlayerSize: Size{Width: DefaultPlayerSize, Height: DefaultPlayerSize},
		Radius:     float64(DefaultScreenHeight/2) * DefaultRadiusRatio,
		NumStars:   DefaultNumStars,
		Debug:      false,
		Speed:      DefaultSpeed,
		StarSize:   DefaultStarSize,
		StarSpeed:  DefaultStarSpeed,
		AngleStep:  DefaultAngleStep,
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
