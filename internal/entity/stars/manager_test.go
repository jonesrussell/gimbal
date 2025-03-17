package stars_test

import (
	"os"
	"testing"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/entity/stars"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testStarSize  = 2.0
	testStarSpeed = 1.0
	testNumStars  = 10
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

func TestNewManager(t *testing.T) {
	t.Parallel()

	bounds := common.Size{Width: 640, Height: 480}
	m := stars.NewManager(bounds, testNumStars, testStarSize, testStarSpeed)
	require.NotNil(t, m)
	assert.Len(t, m.GetStars(), testNumStars)
}

func TestManager_Update(t *testing.T) {
	t.Parallel()

	bounds := common.Size{Width: 640, Height: 480}
	m := stars.NewManager(bounds, testNumStars, testStarSize, testStarSpeed)
	require.NotNil(t, m)

	// Test star movement
	m.Update()
	assert.NotEmpty(t, m.GetStars())

	// Test star position
	star := m.GetStars()[0]
	assert.InEpsilon(t, testStarSize, star.GetSize(), 0.0001)
	assert.InEpsilon(t, testStarSpeed, star.GetSpeed(), 0.0001)
}

func TestManager_Draw(t *testing.T) {
	t.Parallel()

	bounds := common.Size{Width: 800, Height: 600}
	manager := stars.NewManager(bounds, 1, testStarSize, testStarSpeed)
	require.NotNil(t, manager)

	// Draw should not panic
	manager.Draw(nil)
}
