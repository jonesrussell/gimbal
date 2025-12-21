package config

import "time"

// Game timing constants
const (
	// TargetFPS is the target frames per second for the game
	TargetFPS = 60

	// FrameDuration is the duration of a single frame at target FPS
	FrameDuration = time.Second / TargetFPS

	// DeltaTime is the fixed delta time per frame (in seconds)
	DeltaTime = 1.0 / float64(TargetFPS)

	// CollisionTimeout is the maximum time budget for collision detection per frame
	// Using half the frame budget to leave room for other systems
	CollisionTimeout = FrameDuration / 2

	// SlowSystemThreshold is the duration above which a system update is considered slow
	SlowSystemThreshold = 5 * time.Millisecond
)

// Health and combat constants
const (
	// DefaultInvincibilityDuration is the default duration of invincibility after being hit
	DefaultInvincibilityDuration = 2 * time.Second

	// DefaultPlayerHealth is the default starting health for the player
	DefaultPlayerHealth = 3

	// DefaultPlayerMaxHealth is the default maximum health for the player
	DefaultPlayerMaxHealth = 3
)

// Debug constants
const (
	// DebugLogInterval is how often (in frames) to log debug information
	DebugLogInterval = 60
)
