package config

const (
	DefaultScreenWidth  = 1280
	DefaultScreenHeight = 720
	DefaultPlayerSize   = 48 // Reduced for better gameplay balance
	DefaultNumStars     = 25
	DefaultSpeed        = 0.04
	DefaultStarSize     = 5.0
	DefaultStarSpeed    = 40.0
	DefaultAngleStep    = 0.05
	DefaultRadiusRatio  = 0.8
	// CenterDivisor is used to calculate the center point by dividing dimensions
	CenterDivisor = 2

	// Star field defaults
	DefaultStarSpawnRadiusMin = 30.0
	DefaultStarSpawnRadiusMax = 80.0
	DefaultStarMinScale       = 0.3
	DefaultStarMaxScale       = 1.0
	DefaultStarScaleDistance  = 200.0
	DefaultStarResetMargin    = 50.0
)

// GameConfig holds all configurable parameters for the game
type GameConfig struct {
	ScreenSize Size
	PlayerSize Size
	Radius     float64
	NumStars   int
	Debug      bool
	Invincible bool // Player invincibility (only works when Debug is true)
	Speed      float64
	StarSize   float64
	StarSpeed  float64
	AngleStep  float64

	// Star field configuration
	StarSpawnRadiusMin float64
	StarSpawnRadiusMax float64
	StarMinScale       float64
	StarMaxScale       float64
	StarScaleDistance  float64
	StarResetMargin    float64
}

// StarFieldSettings groups star field configuration parameters
type StarFieldSettings struct {
	SpawnRadiusMin float64
	SpawnRadiusMax float64
	MinScale       float64
	MaxScale       float64
	ScaleDistance  float64
	ResetMargin    float64
}

// GameOption is a function that modifies a GameConfig
type GameOption func(*GameConfig)

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

// WithInvincible enables player invincibility (only works when Debug is enabled)
// Note: The constraint that invincible only works when Debug is enabled is enforced
// at the application level, not in this option function.
func WithInvincible(invincible bool) GameOption {
	return func(c *GameConfig) {
		c.Invincible = invincible
	}
}

// DefaultConfig returns a default game configuration
func DefaultConfig() *GameConfig {
	return &GameConfig{
		ScreenSize: Size{
			Width:  DefaultScreenWidth,
			Height: DefaultScreenHeight,
		},
		PlayerSize: Size{
			Width:  DefaultPlayerSize,
			Height: DefaultPlayerSize,
		},
		Radius:     float64(DefaultScreenHeight/CenterDivisor) * DefaultRadiusRatio, // Use height since it's smaller
		NumStars:   DefaultNumStars,
		Debug:      false,
		Invincible: false,
		Speed:      DefaultSpeed,
		StarSize:   DefaultStarSize,
		StarSpeed:  DefaultStarSpeed,
		AngleStep:  DefaultAngleStep,

		// Star field defaults
		StarSpawnRadiusMin: DefaultStarSpawnRadiusMin,
		StarSpawnRadiusMax: DefaultStarSpawnRadiusMax,
		StarMinScale:       DefaultStarMinScale,
		StarMaxScale:       DefaultStarMaxScale,
		StarScaleDistance:  DefaultStarScaleDistance,
		StarResetMargin:    DefaultStarResetMargin,
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

// Size represents dimensions
type Size struct {
	Width, Height int
}
