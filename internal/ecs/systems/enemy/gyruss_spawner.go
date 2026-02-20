package enemy

import (
	"context"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/dbg"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
)

// GyrussSpawner spawns enemies with Gyruss-style components
type GyrussSpawner struct {
	world        donburi.World
	gameConfig   *config.GameConfig
	resourceMgr  *resources.ResourceManager
	screenCenter common.Point

	// Sprite cache
	enemySprites map[string]*ebiten.Image
}

// NewGyrussSpawner creates a new Gyruss-style enemy spawner
func NewGyrussSpawner(
	world donburi.World,
	gameConfig *config.GameConfig,
	resourceMgr *resources.ResourceManager,
) *GyrussSpawner {
	return &GyrussSpawner{
		world:       world,
		gameConfig:  gameConfig,
		resourceMgr: resourceMgr,
		screenCenter: common.Point{
			X: float64(gameConfig.ScreenSize.Width) / 2,
			Y: float64(gameConfig.ScreenSize.Height) / 2,
		},
		enemySprites: make(map[string]*ebiten.Image),
	}
}

// SpawnEnemy spawns an enemy with full Gyruss components from group config
//
//nolint:funlen // Entity setup requires setting many related components sequentially
func (gs *GyrussSpawner) SpawnEnemy(ctx context.Context, groupConfig *managers.EnemyGroupConfig, spawnIndex int) donburi.Entity {
	// Create entity with all Gyruss components
	entity := gs.world.Create(
		core.EnemyTag,
		core.Position,
		core.Sprite,
		core.Size,
		core.Health,
		core.Movement,
		core.EnemyTypeID,
		core.EntryPath,
		core.BehaviorState,
		core.ScaleAnimation,
		core.AttackPattern,
		core.FirePattern,
		core.RetreatTimer,
	)
	entry := gs.world.Entry(entity)

	// Get enemy type data
	enemyType := gs.getEnemyType(groupConfig.EnemyType)
	sprite := gs.getSprite(ctx, groupConfig.EnemyType)

	// Calculate spawn angle based on index (distribute around center)
	baseAngle := float64(spawnIndex) * (2 * math.Pi / float64(groupConfig.Count))

	// Set initial position at screen center (enemies warp in from center)
	core.Position.SetValue(entry, gs.screenCenter)

	// Set sprite
	core.Sprite.SetValue(entry, sprite)

	// Set size (will be scaled by scale animation)
	size := 32 // Default size
	core.Size.SetValue(entry, config.Size{Width: size, Height: size})

	// Set health
	health := gs.getHealthForType(groupConfig.EnemyType)
	core.Health.SetValue(entry, core.NewHealthData(health, health))

	// Set enemy type ID
	core.EnemyTypeID.SetValue(entry, int(enemyType))

	// Set entry path
	pathType := managers.PathTypeFromString(groupConfig.EntryPath.Type)
	direction := managers.ConvertDirection(groupConfig.EntryPath.Parameters.RotationDirection)
	orbitRadius := gs.getOrbitRadius()

	// Calculate end position on orbit ring
	endX := gs.screenCenter.X + orbitRadius*math.Cos(baseAngle)
	endY := gs.screenCenter.Y + orbitRadius*math.Sin(baseAngle)

	core.EntryPath.SetValue(entry, core.EntryPathData{
		PathType:      core.PathType(pathType),
		Progress:      0.0,
		Duration:      groupConfig.EntryPath.Duration,
		ElapsedTime:   0.0,
		StartPosition: gs.screenCenter,
		EndPosition:   common.Point{X: endX, Y: endY},
		Parameters: core.PathParams{
			SpiralTurns:       groupConfig.EntryPath.Parameters.SpiralTurns,
			ArcAngle:          groupConfig.EntryPath.Parameters.ArcAngle,
			RotationDirection: direction,
			StartRadius:       groupConfig.EntryPath.Parameters.StartRadius,
		},
		IsComplete: false,
	})

	// Set behavior state
	behaviorType := managers.BehaviorFromString(groupConfig.Behavior.PostEntry)
	orbitDir := managers.ConvertDirection(groupConfig.Behavior.OrbitDirection)
	orbitDuration := time.Duration(groupConfig.Behavior.OrbitDuration * float64(time.Second))
	core.BehaviorState.SetValue(entry, core.BehaviorStateData{
		CurrentState:      core.StateEntering,
		PreviousState:     core.StateEntering,
		StateTime:         0,
		PostEntryBehavior: core.PostEntryBehavior(behaviorType),
		OrbitDuration:     orbitDuration,
		AttackCooldown:    5 * time.Second,
		AttackCount:       0,
		MaxAttacks:        groupConfig.Behavior.MaxAttacks,
		OrbitDirection:    orbitDir,
		OrbitSpeed:        groupConfig.Behavior.OrbitSpeed,
		TargetOrbitAngle:  baseAngle * 180 / math.Pi,
	})

	// Set scale animation
	easing := managers.EasingFromString(groupConfig.ScaleAnimation.Easing)
	core.ScaleAnimation.SetValue(entry, core.ScaleAnimationData{
		StartScale:  groupConfig.ScaleAnimation.StartScale,
		TargetScale: groupConfig.ScaleAnimation.EndScale,
		Progress:    0.0,
		Duration:    groupConfig.EntryPath.Duration, // Match entry path duration
		ElapsedTime: 0.0,
		Easing:      core.EasingType(easing),
		IsComplete:  false,
	})

	// Set attack pattern
	attackType := managers.AttackPatternFromString(groupConfig.AttackPattern.Type)
	attackCooldown := time.Duration(groupConfig.AttackPattern.Cooldown * float64(time.Second))
	core.AttackPattern.SetValue(entry, core.AttackPatternData{
		PatternType:    core.AttackPatternType(attackType),
		RushSpeed:      groupConfig.AttackPattern.RushSpeed,
		ReturnSpeed:    groupConfig.AttackPattern.ReturnSpeed,
		AttackDuration: attackCooldown,
		AttackTimer:    0,
		TargetPosition: common.Point{},
		ReturnPosition: common.Point{},
		IsActive:       false,
		PairEntityID:   0,
	})

	// Set fire pattern
	fireType := managers.FirePatternFromString(groupConfig.FirePattern.Type)
	core.FirePattern.SetValue(entry, core.FirePatternData{
		PatternType:        core.FirePatternType(fireType),
		FireRate:           groupConfig.FirePattern.FireRate,
		BurstCount:         groupConfig.FirePattern.BurstCount,
		SprayAngle:         groupConfig.FirePattern.SprayAngle,
		ProjectileCount:    groupConfig.FirePattern.ProjectileCount,
		CanFireWhileOrbit:  groupConfig.FirePattern.FireWhileOrbit,
		CanFireWhileAttack: groupConfig.FirePattern.FireWhileAttack,
		LastFireTime:       0,
	})

	// Set retreat timer
	core.RetreatTimer.SetValue(entry, core.RetreatTimerData{
		TimeoutDuration: time.Duration(groupConfig.Retreat.Timeout * float64(time.Second)),
		ElapsedTime:     0,
		IsRetreating:    false,
		RetreatSpeed:    groupConfig.Retreat.Speed,
		RetreatAngle:    0,
	})

	// Set initial movement (will be updated by path system)
	core.Movement.SetValue(entry, core.MovementData{
		Velocity: common.Point{X: 0, Y: 0},
		MaxSpeed: groupConfig.AttackPattern.RushSpeed,
	})

	dbg.Log(dbg.Spawn, "Gyruss enemy spawned (type=%s)", groupConfig.EnemyType)

	return entity
}

// SpawnBoss spawns a boss with Gyruss components
//
//nolint:funlen // Entity setup requires setting many related components sequentially
func (gs *GyrussSpawner) SpawnBoss(ctx context.Context, bossConfig *managers.StageBossConfig) donburi.Entity {
	entity := gs.world.Create(
		core.EnemyTag,
		core.Position,
		core.Sprite,
		core.Size,
		core.Health,
		core.Movement,
		core.EnemyTypeID,
		core.EntryPath,
		core.BehaviorState,
		core.ScaleAnimation,
		core.AttackPattern,
		core.FirePattern,
	)
	entry := gs.world.Entry(entity)

	// Get boss sprite
	sprite := gs.getSprite(ctx, bossConfig.BossType)

	// Set initial position at center
	core.Position.SetValue(entry, gs.screenCenter)

	// Set sprite
	core.Sprite.SetValue(entry, sprite)

	// Set size
	core.Size.SetValue(entry, config.Size{Width: bossConfig.Size, Height: bossConfig.Size})

	// Set health
	core.Health.SetValue(entry, core.NewHealthData(bossConfig.Health, bossConfig.Health))

	// Set enemy type as boss
	core.EnemyTypeID.SetValue(entry, int(EnemyTypeBoss))

	// Set entry path
	pathType := managers.PathTypeFromString(bossConfig.EntryPath.Type)
	bossOrbitRadius := gs.getOrbitRadius() * 0.6 // Boss orbits closer

	core.EntryPath.SetValue(entry, core.EntryPathData{
		PathType:      core.PathType(pathType),
		Progress:      0.0,
		Duration:      bossConfig.EntryPath.Duration,
		ElapsedTime:   0.0,
		StartPosition: gs.screenCenter,
		EndPosition:   common.Point{X: gs.screenCenter.X + bossOrbitRadius, Y: gs.screenCenter.Y},
		Parameters: core.PathParams{
			RotationDirection: 1,
		},
		IsComplete: false,
	})

	// Set behavior state
	behaviorType := managers.BehaviorFromString(bossConfig.Behavior.PostEntry)
	orbitDir := managers.ConvertDirection(bossConfig.Behavior.OrbitDirection)
	bossOrbitDuration := time.Duration(bossConfig.Behavior.OrbitDuration * float64(time.Second))
	core.BehaviorState.SetValue(entry, core.BehaviorStateData{
		CurrentState:      core.StateEntering,
		PreviousState:     core.StateEntering,
		StateTime:         0,
		PostEntryBehavior: core.PostEntryBehavior(behaviorType),
		OrbitDuration:     bossOrbitDuration,
		AttackCooldown:    time.Duration(bossConfig.AttackPattern.Cooldown * float64(time.Second)),
		AttackCount:       0,
		MaxAttacks:        bossConfig.Behavior.MaxAttacks,
		OrbitDirection:    orbitDir,
		OrbitSpeed:        bossConfig.Behavior.OrbitSpeed,
		TargetOrbitAngle:  0,
	})

	// Set scale animation (boss scales up from center)
	core.ScaleAnimation.SetValue(entry, core.ScaleAnimationData{
		StartScale:  0.2,
		TargetScale: 1.0,
		Progress:    0.0,
		Duration:    bossConfig.EntryPath.Duration,
		ElapsedTime: 0.0,
		Easing:      core.EasingEaseOut,
		IsComplete:  false,
	})

	// Set attack pattern
	attackType := managers.AttackPatternFromString(bossConfig.AttackPattern.Type)
	bossAttackCooldown := time.Duration(bossConfig.AttackPattern.Cooldown * float64(time.Second))
	core.AttackPattern.SetValue(entry, core.AttackPatternData{
		PatternType:    core.AttackPatternType(attackType),
		RushSpeed:      bossConfig.AttackPattern.RushSpeed,
		ReturnSpeed:    bossConfig.AttackPattern.ReturnSpeed,
		AttackDuration: bossAttackCooldown,
		AttackTimer:    0,
		TargetPosition: common.Point{},
		ReturnPosition: common.Point{},
		IsActive:       false,
		PairEntityID:   0,
	})

	// Set fire pattern
	fireType := managers.FirePatternFromString(bossConfig.FirePattern.Type)
	core.FirePattern.SetValue(entry, core.FirePatternData{
		PatternType:        core.FirePatternType(fireType),
		FireRate:           bossConfig.FirePattern.FireRate,
		BurstCount:         bossConfig.FirePattern.BurstCount,
		SprayAngle:         bossConfig.FirePattern.SprayAngle,
		ProjectileCount:    bossConfig.FirePattern.ProjectileCount,
		CanFireWhileOrbit:  bossConfig.FirePattern.FireWhileOrbit,
		CanFireWhileAttack: bossConfig.FirePattern.FireWhileAttack,
		LastFireTime:       0,
	})

	// Set movement
	core.Movement.SetValue(entry, core.MovementData{
		Velocity: common.Point{X: 0, Y: 0},
		MaxSpeed: bossConfig.AttackPattern.RushSpeed,
	})

	dbg.Log(dbg.Spawn, "Gyruss boss entity created (type=%s)", bossConfig.BossType)

	return entity
}

// getEnemyType converts string type to EnemyType enum
func (gs *GyrussSpawner) getEnemyType(typeStr string) EnemyType {
	switch typeStr {
	case EnemyTypeStrBasic:
		return EnemyTypeBasic
	case EnemyTypeStrHeavy:
		return EnemyTypeHeavy
	case EnemyTypeStrBoss:
		return EnemyTypeBoss
	case EnemyTypeStrSatellite:
		return EnemyTypeBasic // Map to basic for now
	default:
		return EnemyTypeBasic
	}
}

// getHealthForType returns health for an enemy type
func (gs *GyrussSpawner) getHealthForType(typeStr string) int {
	switch typeStr {
	case EnemyTypeStrBasic:
		return 1
	case EnemyTypeStrHeavy:
		return 3
	case EnemyTypeStrSatellite:
		return 1
	default:
		return 1
	}
}

// getOrbitRadius returns the orbit radius for enemies
func (gs *GyrussSpawner) getOrbitRadius() float64 {
	// Use the smaller screen dimension to calculate orbit radius
	minDim := float64(gs.gameConfig.ScreenSize.Width)
	if float64(gs.gameConfig.ScreenSize.Height) < minDim {
		minDim = float64(gs.gameConfig.ScreenSize.Height)
	}
	return minDim * 0.35 // 35% of screen for orbit
}

// getSprite gets or creates a sprite for an enemy type
func (gs *GyrussSpawner) getSprite(ctx context.Context, typeStr string) *ebiten.Image {
	// Check cache
	if sprite, ok := gs.enemySprites[typeStr]; ok {
		return sprite
	}

	// Try to load from resource manager
	spriteName := typeStr + "_enemy"
	sprite, exists := gs.resourceMgr.GetSprite(ctx, spriteName)
	if exists {
		gs.enemySprites[typeStr] = sprite
		return sprite
	}

	// Create placeholder
	size := 32
	sprite = ebiten.NewImage(size, size)

	// Color based on type
	var clr color.RGBA
	switch typeStr {
	case EnemyTypeStrBasic:
		clr = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Red
	case EnemyTypeStrHeavy:
		clr = color.RGBA{R: 255, G: 165, B: 0, A: 255} // Orange
	case EnemyTypeStrSatellite:
		clr = color.RGBA{R: 100, G: 100, B: 255, A: 255} // Light blue
	case "earth_boss":
		clr = color.RGBA{R: 128, G: 0, B: 128, A: 255} // Purple
	default:
		clr = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Red
	}
	sprite.Fill(clr)

	gs.enemySprites[typeStr] = sprite
	dbg.Log(dbg.System, "Created placeholder sprite (type=%s)", typeStr)

	return sprite
}
