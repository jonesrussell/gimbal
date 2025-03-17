package player_test

import (
	"image"
	"testing"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/entity/player"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockSprite implements the SpriteImage interface for testing
type MockSprite struct {
	bounds image.Rectangle
}

func (m *MockSprite) Bounds() image.Rectangle {
	return m.bounds
}

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  *common.EntityConfig
		sprite  player.SpriteImage
		wantErr bool
	}{
		{
			name: "valid config and sprite",
			config: &common.EntityConfig{
				Position: common.Point{X: 100, Y: 100},
				Size:     common.Size{Width: 32, Height: 32},
				Radius:   100,
				Speed:    2.0,
			},
			sprite: &MockSprite{
				bounds: image.Rect(0, 0, 32, 32),
			},
			wantErr: false,
		},
		{
			name:    "nil config",
			config:  nil,
			sprite:  &MockSprite{},
			wantErr: true,
		},
		{
			name: "nil sprite",
			config: &common.EntityConfig{
				Position: common.Point{X: 100, Y: 100},
				Size:     common.Size{Width: 32, Height: 32},
				Radius:   100,
				Speed:    2.0,
			},
			sprite:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := player.New(tt.config, tt.sprite)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.config, got.Config())
			assert.Equal(t, tt.sprite, got.Sprite())
		})
	}
}

func TestPlayer_Angle(t *testing.T) {
	t.Parallel()

	// Create test config
	config := &common.EntityConfig{
		Position: common.Point{X: 100, Y: 100},
		Size:     common.Size{Width: 32, Height: 32},
		Speed:    1.0,
		Radius:   50.0,
	}

	// Create mock sprite
	sprite := player.NewMockImage(32, 32)

	// Create player
	p, err := player.New(config, sprite)
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

	// Create test config
	config := &common.EntityConfig{
		Position: common.Point{X: 100, Y: 100},
		Size:     common.Size{Width: 32, Height: 32},
		Speed:    1.0,
		Radius:   50.0,
	}

	// Create mock sprite
	sprite := player.NewMockImage(32, 32)

	// Create player
	p, err := player.New(config, sprite)
	require.NoError(t, err)

	// Test position setting and getting
	testPos := common.Point{X: 200, Y: 200}
	p.SetPosition(testPos)
	assert.Equal(t, testPos, p.GetPosition())
}
