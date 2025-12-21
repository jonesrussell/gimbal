package config

import (
	"math"
	"testing"
	"time"
)

func TestFrameDuration(t *testing.T) {
	// FrameDuration should be 1/60 seconds
	expected := time.Second / 60
	if FrameDuration != expected {
		t.Errorf("FrameDuration = %v, want %v", FrameDuration, expected)
	}
}

func TestDeltaTime(t *testing.T) {
	// DeltaTime should be 1/60 (approximately 0.016666...)
	expected := 1.0 / 60.0
	tolerance := 0.0001
	if math.Abs(DeltaTime-expected) > tolerance {
		t.Errorf("DeltaTime = %v, want approximately %v", DeltaTime, expected)
	}
}

func TestTargetFPS(t *testing.T) {
	if TargetFPS != 60 {
		t.Errorf("TargetFPS = %d, want 60", TargetFPS)
	}
}

func TestCollisionTimeout(t *testing.T) {
	// CollisionTimeout should be half of FrameDuration
	expected := FrameDuration / 2
	if CollisionTimeout != expected {
		t.Errorf("CollisionTimeout = %v, want %v", CollisionTimeout, expected)
	}
}

func TestSlowSystemThreshold(t *testing.T) {
	expected := 5 * time.Millisecond
	if SlowSystemThreshold != expected {
		t.Errorf("SlowSystemThreshold = %v, want %v", SlowSystemThreshold, expected)
	}
}

func TestDefaultInvincibilityDuration(t *testing.T) {
	expected := 2 * time.Second
	if DefaultInvincibilityDuration != expected {
		t.Errorf("DefaultInvincibilityDuration = %v, want %v", DefaultInvincibilityDuration, expected)
	}
}

func TestDefaultPlayerHealth(t *testing.T) {
	if DefaultPlayerHealth != 3 {
		t.Errorf("DefaultPlayerHealth = %d, want 3", DefaultPlayerHealth)
	}
}

func TestDefaultPlayerMaxHealth(t *testing.T) {
	if DefaultPlayerMaxHealth != 3 {
		t.Errorf("DefaultPlayerMaxHealth = %d, want 3", DefaultPlayerMaxHealth)
	}
}

func TestDebugLogInterval(t *testing.T) {
	if DebugLogInterval != 60 {
		t.Errorf("DebugLogInterval = %d, want 60", DebugLogInterval)
	}
}

func TestDefaultScreenDimensions(t *testing.T) {
	if DefaultScreenWidth != 1280 {
		t.Errorf("DefaultScreenWidth = %d, want 1280", DefaultScreenWidth)
	}
	if DefaultScreenHeight != 720 {
		t.Errorf("DefaultScreenHeight = %d, want 720", DefaultScreenHeight)
	}
}

func TestDefaultPlayerSize(t *testing.T) {
	if DefaultPlayerSize != 48 {
		t.Errorf("DefaultPlayerSize = %d, want 48", DefaultPlayerSize)
	}
}

func TestDefaultNumStars(t *testing.T) {
	if DefaultNumStars != 100 {
		t.Errorf("DefaultNumStars = %d, want 100", DefaultNumStars)
	}
}

func TestDefaultSpeed(t *testing.T) {
	if DefaultSpeed != 0.04 {
		t.Errorf("DefaultSpeed = %v, want 0.04", DefaultSpeed)
	}
}

func TestDefaultStarSize(t *testing.T) {
	if DefaultStarSize != 5.0 {
		t.Errorf("DefaultStarSize = %v, want 5.0", DefaultStarSize)
	}
}

func TestDefaultStarSpeed(t *testing.T) {
	if DefaultStarSpeed != 40.0 {
		t.Errorf("DefaultStarSpeed = %v, want 40.0", DefaultStarSpeed)
	}
}

func TestDefaultAngleStep(t *testing.T) {
	if DefaultAngleStep != 0.05 {
		t.Errorf("DefaultAngleStep = %v, want 0.05", DefaultAngleStep)
	}
}

func TestDefaultRadiusRatio(t *testing.T) {
	if DefaultRadiusRatio != 0.8 {
		t.Errorf("DefaultRadiusRatio = %v, want 0.8", DefaultRadiusRatio)
	}
}

func TestCenterDivisor(t *testing.T) {
	if CenterDivisor != 2 {
		t.Errorf("CenterDivisor = %d, want 2", CenterDivisor)
	}
}

func TestDefaultStarSpawnRadiusMin(t *testing.T) {
	if DefaultStarSpawnRadiusMin != 30.0 {
		t.Errorf("DefaultStarSpawnRadiusMin = %v, want 30.0", DefaultStarSpawnRadiusMin)
	}
}

func TestDefaultStarSpawnRadiusMax(t *testing.T) {
	if DefaultStarSpawnRadiusMax != 80.0 {
		t.Errorf("DefaultStarSpawnRadiusMax = %v, want 80.0", DefaultStarSpawnRadiusMax)
	}
}

func TestDefaultStarMinScale(t *testing.T) {
	if DefaultStarMinScale != 0.3 {
		t.Errorf("DefaultStarMinScale = %v, want 0.3", DefaultStarMinScale)
	}
}

func TestDefaultStarMaxScale(t *testing.T) {
	if DefaultStarMaxScale != 1.0 {
		t.Errorf("DefaultStarMaxScale = %v, want 1.0", DefaultStarMaxScale)
	}
}

func TestDefaultStarScaleDistance(t *testing.T) {
	if DefaultStarScaleDistance != 200.0 {
		t.Errorf("DefaultStarScaleDistance = %v, want 200.0", DefaultStarScaleDistance)
	}
}

func TestDefaultStarResetMargin(t *testing.T) {
	if DefaultStarResetMargin != 50.0 {
		t.Errorf("DefaultStarResetMargin = %v, want 50.0", DefaultStarResetMargin)
	}
}

func TestTimingConstants_Relationship(t *testing.T) {
	// Verify the relationship between timing constants
	calculatedFrameDuration := time.Second / time.Duration(TargetFPS)
	if FrameDuration != calculatedFrameDuration {
		t.Errorf("FrameDuration should equal time.Second/TargetFPS. Got %v, want %v",
			FrameDuration, calculatedFrameDuration)
	}

	calculatedDeltaTime := 1.0 / float64(TargetFPS)
	tolerance := 0.0001
	if math.Abs(DeltaTime-calculatedDeltaTime) > tolerance {
		t.Errorf("DeltaTime should equal 1/TargetFPS. Got %v, want approximately %v", DeltaTime, calculatedDeltaTime)
	}

	if CollisionTimeout != FrameDuration/2 {
		t.Errorf("CollisionTimeout should equal FrameDuration/2. Got %v, want %v", CollisionTimeout, FrameDuration/2)
	}
}

func TestStarConstants_Relationship(t *testing.T) {
	// Verify star spawn radius relationship
	if DefaultStarSpawnRadiusMin >= DefaultStarSpawnRadiusMax {
		t.Errorf("DefaultStarSpawnRadiusMin (%v) should be less than DefaultStarSpawnRadiusMax (%v)",
			DefaultStarSpawnRadiusMin, DefaultStarSpawnRadiusMax)
	}

	// Verify star scale relationship
	if DefaultStarMinScale >= DefaultStarMaxScale {
		t.Errorf("DefaultStarMinScale (%v) should be less than DefaultStarMaxScale (%v)",
			DefaultStarMinScale, DefaultStarMaxScale)
	}
}

func TestHealthConstants_Relationship(t *testing.T) {
	// Health and max health should match
	if DefaultPlayerHealth != DefaultPlayerMaxHealth {
		t.Errorf("DefaultPlayerHealth (%d) should equal DefaultPlayerMaxHealth (%d)",
			DefaultPlayerHealth, DefaultPlayerMaxHealth)
	}
}
