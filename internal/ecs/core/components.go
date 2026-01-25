package core

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/math"
)

// MovementPattern represents movement behavior patterns
// This is a type alias to avoid import cycles - it should match enemy.MovementPattern
type MovementPattern int

// Re-export timing constants for use by other packages
var (
	DefaultInvincibilityDuration = config.DefaultInvincibilityDuration
)

// Component tags for different entity types
var (
	// PlayerTag marks an entity as a player
	PlayerTag = donburi.NewTag()
	// StarTag marks an entity as a star
	StarTag = donburi.NewTag()
	// EnemyTag marks an entity as an enemy
	EnemyTag = donburi.NewTag()
	// ProjectileTag marks an entity as a player projectile
	ProjectileTag = donburi.NewTag()
	// EnemyProjectileTag marks an entity as an enemy projectile
	EnemyProjectileTag = donburi.NewTag()
	// PowerUpTag marks an entity as a collectible power-up
	PowerUpTag = donburi.NewTag()
)

// Components
var (
	// Position component stores entity position
	Position = donburi.NewComponentType[common.Point]()
	// Sprite component stores the entity's sprite
	Sprite = donburi.NewComponentType[*ebiten.Image]()
	// Movement component stores movement data
	Movement = donburi.NewComponentType[MovementData]()
	// Orbital component stores orbital movement data
	Orbital = donburi.NewComponentType[OrbitalData]()
	// Size component stores entity dimensions
	Size = donburi.NewComponentType[config.Size]()
	// Speed component stores movement speed
	Speed = donburi.NewComponentType[float64]()
	// Angle component stores rotation angle
	Angle = donburi.NewComponentType[math.Angle]()
	// Scale component stores scaling factor
	Scale = donburi.NewComponentType[float64]()
	// Health component stores entity health data
	Health = donburi.NewComponentType[HealthData]()
	// EnemyTypeID component stores the enemy type identifier
	EnemyTypeID = donburi.NewComponentType[int]()

	// Gyruss-style components for enemy entry and behavior

	// EntryPath component stores parametric entry path data for enemy warp-in
	EntryPath = donburi.NewComponentType[EntryPathData]()
	// BehaviorState component stores enemy behavior state machine data
	BehaviorState = donburi.NewComponentType[BehaviorStateData]()
	// ScaleAnimation component stores visual scaling animation data
	ScaleAnimation = donburi.NewComponentType[ScaleAnimationData]()
	// AttackPattern component stores attack behavior configuration
	AttackPattern = donburi.NewComponentType[AttackPatternData]()
	// FirePattern component stores firing behavior configuration
	FirePattern = donburi.NewComponentType[FirePatternData]()
	// RetreatTimer component stores retreat timeout data
	RetreatTimer = donburi.NewComponentType[RetreatTimerData]()
	// PowerUpData component identifies power-up type and behavior
	PowerUpData = donburi.NewComponentType[PowerUpTypeData]()
)

// MovementData represents movement information
type MovementData struct {
	Velocity    common.Point
	MaxSpeed    float64
	Pattern     MovementPattern // Movement pattern type (should match enemy.MovementPattern)
	PatternTime time.Duration   // Time accumulator for pattern-based movement
	BaseAngle   float64         // Base angle for pattern calculations
	BaseSpeed   float64         // Base speed for pattern calculations
}

// OrbitalData represents orbital movement information
type OrbitalData struct {
	Center       common.Point
	Radius       float64
	OrbitalAngle math.Angle
	FacingAngle  math.Angle
}

// HealthData represents health and invincibility information
type HealthData struct {
	Current               int           // Current health/lives
	Maximum               int           // Maximum health/lives
	InvincibilityTime     time.Duration // Time remaining for invincibility
	IsInvincible          bool          // Whether entity is currently invincible
	InvincibilityDuration time.Duration // Duration of invincibility when hit
}

// NewHealthData creates a new health data with default invincibility duration
func NewHealthData(current, maximum int) HealthData {
	return HealthData{
		Current:               current,
		Maximum:               maximum,
		InvincibilityTime:     0,
		IsInvincible:          false,
		InvincibilityDuration: DefaultInvincibilityDuration,
	}
}

// ============================================================================
// Gyruss-Style Component Data Structures
// ============================================================================

// PathType represents the type of parametric entry path
type PathType int

const (
	PathTypeSpiralIn      PathType = iota // Spiral inward from edge
	PathTypeArcSweep                      // Arc sweep from edge
	PathTypeStraightIn                    // Straight line from edge
	PathTypeLoopEntry                     // Loop path entry
	PathTypeRandomOutward                 // Random outward movement
)

// BehaviorStateType represents enemy behavior states
type BehaviorStateType int

const (
	StateEntering   BehaviorStateType = iota // Following entry path
	StateOrbiting                            // Orbiting at ring
	StateAttacking                           // Executing attack pattern
	StateRetreating                          // Returning to orbit or exiting
	StateHovering                            // Hovering at center
)

// PostEntryBehavior defines what enemies do after entering
type PostEntryBehavior int

const (
	BehaviorOrbitOnly            PostEntryBehavior = iota // Just orbit
	BehaviorOrbitThenAttack                               // Orbit then attack
	BehaviorImmediateAttack                               // Attack immediately
	BehaviorHoverCenterThenOrbit                          // Hover at center, then orbit
)

// AttackPatternType represents attack movement patterns
type AttackPatternType int

const (
	AttackNone         AttackPatternType = iota // No attack
	AttackSingleRush                            // Rush toward player, return to orbit
	AttackPairedRush                            // Coordinated two-enemy rush
	AttackLoopbackRush                          // Loop through center and return
	AttackSuicideDive                           // Direct dive at player
)

// FirePatternType represents firing patterns
type FirePatternType int

const (
	FireNone       FirePatternType = iota // No firing
	FireSingleShot                        // Aimed single shot at player
	FireBurst                             // 3-shot burst
	FireSpray                             // Multi-directional spread
)

// PowerUpType represents power-up types
type PowerUpType int

const (
	PowerUpDoubleShot PowerUpType = iota // Double shot weapon
	PowerUpExtraLife                     // Extra life
)

// EasingType represents easing functions for animations
type EasingType int

const (
	EasingLinear    EasingType = iota // Linear interpolation
	EasingEaseIn                      // Slow start, fast end
	EasingEaseOut                     // Fast start, slow end
	EasingEaseInOut                   // Slow start and end
)

// EntryPathData defines parametric entry path for enemy warp-in
type EntryPathData struct {
	PathType      PathType     // Type of entry path
	Progress      float64      // 0.0 to 1.0
	Duration      float64      // Total duration in seconds
	ElapsedTime   float64      // Time elapsed since entry started
	StartPosition common.Point // Starting position (near center)
	EndPosition   common.Point // Target position (on orbit ring)
	Parameters    PathParams   // Path-specific parameters
	IsComplete    bool         // Whether entry is complete
}

// PathParams stores path-specific configuration
type PathParams struct {
	SpiralTurns       float64 // Number of spiral turns (for spiral paths)
	ArcAngle          float64 // Arc angle in degrees (for arc paths)
	RotationDirection int     // 1 for clockwise, -1 for counter-clockwise
	CurveIntensity    float64 // General curvature intensity
	StartRadius       float64 // Starting radius from center
}

// BehaviorStateData stores enemy behavior state machine
type BehaviorStateData struct {
	CurrentState      BehaviorStateType // Current behavior state
	PreviousState     BehaviorStateType // Previous state (for transitions)
	StateTime         time.Duration     // Time in current state
	PostEntryBehavior PostEntryBehavior // Behavior after entry completes
	OrbitDuration     time.Duration     // Time to orbit before attacking
	AttackCooldown    time.Duration     // Cooldown between attacks
	AttackCount       int               // Number of attacks performed
	MaxAttacks        int               // Max attacks before retreating (0 = unlimited)
	OrbitDirection    int               // 1 for clockwise, -1 for counter-clockwise
	OrbitSpeed        float64           // Orbital speed in degrees per second
	TargetOrbitAngle  float64           // Target angle on orbit ring
}

// ScaleAnimationData stores visual scaling animation
type ScaleAnimationData struct {
	StartScale  float64    // Scale at start of animation
	TargetScale float64    // Scale at end of animation
	Progress    float64    // Animation progress 0.0 to 1.0
	Duration    float64    // Animation duration in seconds
	ElapsedTime float64    // Time elapsed
	Easing      EasingType // Easing function type
	IsComplete  bool       // Whether animation is complete
}

// AttackPatternData stores attack behavior configuration
type AttackPatternData struct {
	PatternType    AttackPatternType // Type of attack pattern
	RushSpeed      float64           // Speed during rush attacks
	ReturnSpeed    float64           // Speed when returning to orbit
	AttackDuration time.Duration     // Duration of attack phase
	AttackTimer    time.Duration     // Timer for current attack
	TargetPosition common.Point      // Target position for attack
	ReturnPosition common.Point      // Position to return to
	IsActive       bool              // Whether attack is currently active
	PairEntityID   int               // Partner entity ID for paired attacks
}

// FirePatternData stores firing behavior configuration
type FirePatternData struct {
	PatternType        FirePatternType // Type of fire pattern
	FireRate           float64         // Shots per second
	BurstCount         int             // Number of shots in a burst
	BurstFired         int             // Shots fired in current burst
	SprayAngle         float64         // Angle spread for spray pattern
	ProjectileCount    int             // Number of projectiles for spray
	LastFireTime       time.Duration   // Time since last shot
	FireCooldown       time.Duration   // Time between shots
	CanFireWhileOrbit  bool            // Fire while orbiting
	CanFireWhileAttack bool            // Fire while attacking
}

// RetreatTimerData stores retreat behavior
type RetreatTimerData struct {
	TimeoutDuration time.Duration // Maximum time before forced retreat
	ElapsedTime     time.Duration // Time elapsed
	IsRetreating    bool          // Whether currently retreating
	RetreatSpeed    float64       // Speed during retreat
	RetreatAngle    float64       // Direction of retreat
}

// PowerUpTypeData identifies power-up type and behavior
type PowerUpTypeData struct {
	Type         PowerUpType   // Type of power-up
	Duration     time.Duration // Duration for temporary power-ups (0 = permanent)
	OrbitalAngle float64       // Angle for orbital movement
	OrbitalSpeed float64       // Speed of orbital movement
	LifeTime     time.Duration // Time alive before despawning
	MaxLifeTime  time.Duration // Maximum time before despawn
}

// NewEntryPathData creates a new entry path data
func NewEntryPathData(pathType PathType, duration float64, start, end common.Point) EntryPathData {
	return EntryPathData{
		PathType:      pathType,
		Progress:      0.0,
		Duration:      duration,
		ElapsedTime:   0.0,
		StartPosition: start,
		EndPosition:   end,
		Parameters:    PathParams{SpiralTurns: 1.5, RotationDirection: 1},
		IsComplete:    false,
	}
}

// NewBehaviorStateData creates a new behavior state data
func NewBehaviorStateData(postEntry PostEntryBehavior, orbitDuration time.Duration) BehaviorStateData {
	return BehaviorStateData{
		CurrentState:      StateEntering,
		PreviousState:     StateEntering,
		StateTime:         0,
		PostEntryBehavior: postEntry,
		OrbitDuration:     orbitDuration,
		AttackCooldown:    5 * time.Second,
		AttackCount:       0,
		MaxAttacks:        3,
		OrbitDirection:    1,
		OrbitSpeed:        45.0, // degrees per second
	}
}

// NewScaleAnimationData creates a new scale animation data
func NewScaleAnimationData(startScale, targetScale, duration float64) ScaleAnimationData {
	return ScaleAnimationData{
		StartScale:  startScale,
		TargetScale: targetScale,
		Progress:    0.0,
		Duration:    duration,
		ElapsedTime: 0.0,
		Easing:      EasingEaseOut,
		IsComplete:  false,
	}
}
