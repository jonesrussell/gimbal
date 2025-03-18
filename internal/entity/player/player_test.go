package player_test

import (
	"os"
	"testing"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/entity/player"
	"github.com/jonesrussell/gimbal/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Set up test environment
	os.Setenv("EBITEN_TEST", "1")

	// Run tests
	code := m.Run()

	// Clean up
	os.Unsetenv("EBITEN_TEST")
	os.Exit(code)
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockLogger := logger.NewMock()

	tests := []struct {
		name    string
		config  *common.EntityConfig
		sprite  player.Drawable
		wantErr bool
	}{
		{
			name: "valid config and sprite",
			config: &common.EntityConfig{
				Position: common.Point{X: 100, Y: 100},
				Size:     common.Size{Width: 32, Height: 32},
				Speed:    1.0,
				Radius:   100,
			},
			sprite:  player.NewTestSprite(32, 32),
			wantErr: false,
		},
		{
			name:    "nil config",
			config:  nil,
			sprite:  player.NewTestSprite(32, 32),
			wantErr: true,
		},
		{
			name: "nil sprite",
			config: &common.EntityConfig{
				Position: common.Point{X: 100, Y: 100},
				Size:     common.Size{Width: 32, Height: 32},
				Speed:    1.0,
				Radius:   100,
			},
			sprite:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := player.New(tt.config, tt.sprite, mockLogger)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.config, got.Config())
			assert.Equal(t, tt.sprite, got.Sprite())
			assert.InEpsilon(t, tt.config.Speed, got.GetSpeed(), 1e-6)
		})
	}
}

func TestPlayer_Angle(t *testing.T) {
	t.Parallel()

	mockLogger := logger.NewMock()

	// Create test config
	config := &common.EntityConfig{
		Position: common.Point{X: 100, Y: 100},
		Size:     common.Size{Width: 32, Height: 32},
		Speed:    1.0,
		Radius:   50.0,
	}

	// Create test sprite
	sprite := player.NewTestSprite(32, 32)

	// Create player
	p, err := player.New(config, sprite, mockLogger)
	require.NoError(t, err)

	// Test angle setting and getting
	testAngle := common.Angle(45)
	p.SetAngle(testAngle)
	assert.InEpsilon(t, float64(testAngle), float64(p.GetAngle()), 1e-6)

	// Test facing angle
	testFacingAngle := common.Angle(90)
	p.SetFacingAngle(testFacingAngle)
	assert.InEpsilon(t, float64(testFacingAngle), float64(p.GetFacingAngle()), 1e-6)
}

func TestPlayer_Position(t *testing.T) {
	t.Parallel()

	mockLogger := logger.NewMock()

	// Create test config
	config := &common.EntityConfig{
		Position: common.Point{X: 100, Y: 100},
		Size:     common.Size{Width: 32, Height: 32},
		Speed:    1.0,
		Radius:   50.0,
	}

	// Create test sprite
	sprite := player.NewTestSprite(32, 32)

	// Create player
	p, err := player.New(config, sprite, mockLogger)
	require.NoError(t, err)
	require.NotNil(t, p)

	// Test initial position (should be at top of circle, 0 degrees)
	pos := p.GetPosition()
	assert.InEpsilon(t, 100.0, pos.X, 1e-6) // center.X + radius * sin(0) = 100 + 50 * 0
	assert.InEpsilon(t, 50.0, pos.Y, 1e-6)  // center.Y - radius * cos(0) = 100 - 50 * 1

	// Test position after setting angle to 90 degrees (right)
	p.SetAngle(90)
	pos = p.GetPosition()
	assert.InEpsilon(t, 150.0, pos.X, 1e-6) // center.X + radius * sin(90) = 100 + 50 * 1
	assert.InEpsilon(t, 100.0, pos.Y, 1e-6) // center.Y - radius * cos(90) = 100 - 50 * 0
}

func TestPlayer_Update(t *testing.T) {
	t.Parallel()

	mockLogger := logger.NewMock()

	tests := []struct {
		name     string
		angle    common.Angle
		expected common.Point
	}{
		{
			name:     "0 degrees",
			angle:    0,
			expected: common.Point{X: 100, Y: 50}, // 100 + 50 (radius)
		},
		{
			name:     "90 degrees",
			angle:    90,
			expected: common.Point{X: 150, Y: 100}, // 100 + 50 (radius)
		},
		{
			name:     "180 degrees",
			angle:    180,
			expected: common.Point{X: 100, Y: 150}, // 100 + 50 (radius)
		},
		{
			name:     "270 degrees",
			angle:    270,
			expected: common.Point{X: 50, Y: 100}, // 100 - 50 (radius)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create test config
			config := &common.EntityConfig{
				Position: common.Point{X: 100, Y: 100},
				Size:     common.Size{Width: 32, Height: 32},
				Speed:    1.0,
				Radius:   50.0,
			}

			// Create test sprite
			sprite := player.NewTestSprite(32, 32)

			// Create player
			p, err := player.New(config, sprite, mockLogger)
			require.NoError(t, err)
			require.NotNil(t, p)

			p.SetAngle(tt.angle)
			p.Update()

			pos := p.GetPosition()
			assert.InEpsilon(t, tt.expected.X, pos.X, 0.0001)
			assert.InEpsilon(t, tt.expected.Y, pos.Y, 0.0001)
		})
	}
}
