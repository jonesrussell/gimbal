package ecs

// ECS-wide constants for configuration and magic numbers

// Enemy system
const (
	DefaultEnemySpawnIntervalFrames = 60 // 1 second at 60fps
	DefaultEnemyWaveCount           = 5
	EnemySpawnRadiusMin             = 50.0
	EnemySpawnRadiusMax             = 150.0
	EnemyWaveMargin                 = 50.0
	EnemyBossRadius                 = 100.0
	EnemySwarmDroneSize             = 16
	EnemyHeavyCruiserSize           = 32
	EnemyBossSize                   = 64
	EnemyAsteroidSize               = 24
)

// Weapon/projectile system
const (
	DefaultWeaponFireIntervalFrames = 10 // 6 shots/sec at 60fps
	DefaultProjectileSpeed          = 5.0
	DefaultProjectileSize           = 4
	ProjectileOffset                = 20.0
	ProjectileMargin                = 50.0
)
