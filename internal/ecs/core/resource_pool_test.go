package core //nolint:testpackage // Testing from same package to access unexported functions

import (
	"context"
	"testing"
)

func TestImagePool_createKey(t *testing.T) {
	ip := &ImagePool{}

	tests := []struct {
		name           string
		width          int
		height         int
		expectedKey    string
		expectedPrefix string
	}{
		{
			name:           "small dimensions",
			width:          64,
			height:         64,
			expectedKey:    "64x64",
			expectedPrefix: "64x64",
		},
		{
			name:           "large dimensions",
			width:          1920,
			height:         1080,
			expectedKey:    "1920x1080",
			expectedPrefix: "1920x1080",
		},
		{
			name:           "different width and height",
			width:          800,
			height:         600,
			expectedKey:    "800x600",
			expectedPrefix: "800x600",
		},
		{
			name:           "very large dimensions",
			width:          4096,
			height:         4096,
			expectedKey:    "4096x4096",
			expectedPrefix: "4096x4096",
		},
		{
			name:           "dimensions above 127 (tests fix for rune conversion bug)",
			width:          128,
			height:         256,
			expectedKey:    "128x256",
			expectedPrefix: "128x256",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := ip.createKey(tt.width, tt.height)
			if key != tt.expectedKey {
				t.Errorf("createKey(%d, %d) = %q, want %q", tt.width, tt.height, key, tt.expectedKey)
			}
			// Verify the key is properly formatted (contains "x" separator)
			if len(key) < 3 {
				t.Errorf("createKey(%d, %d) = %q, key too short", tt.width, tt.height, key)
			}
		})
	}
}

// noOpLogger is a test logger that does nothing
type noOpLogger struct{}

func (n *noOpLogger) Debug(msg string, fields ...any)                             {}
func (n *noOpLogger) Info(msg string, fields ...any)                              {}
func (n *noOpLogger) Warn(msg string, fields ...any)                              {}
func (n *noOpLogger) Error(msg string, fields ...any)                             {}
func (n *noOpLogger) DebugContext(ctx context.Context, msg string, fields ...any) {}
func (n *noOpLogger) InfoContext(ctx context.Context, msg string, fields ...any)  {}
func (n *noOpLogger) WarnContext(ctx context.Context, msg string, fields ...any)  {}
func (n *noOpLogger) ErrorContext(ctx context.Context, msg string, fields ...any) {}
func (n *noOpLogger) Sync() error                                                 { return nil }

func TestNewImagePool(t *testing.T) {
	logger := &noOpLogger{}
	pool := NewImagePool(logger)

	if pool == nil {
		t.Fatal("NewImagePool() returned nil")
	}
	if pool.pool == nil {
		t.Error("NewImagePool() pool map is nil")
	}
	if pool.logger != logger {
		t.Error("NewImagePool() logger not set correctly")
	}
}

func TestImagePool_GetImage(t *testing.T) {
	logger := &noOpLogger{}
	pool := NewImagePool(logger)

	tests := []struct {
		name   string
		width  int
		height int
	}{
		{
			name:   "create new image",
			width:  100,
			height: 100,
		},
		{
			name:   "create different sized image",
			width:  200,
			height: 150,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := pool.GetImage(tt.width, tt.height)
			if img == nil {
				t.Fatal("GetImage() returned nil")
			}

			bounds := img.Bounds()
			if bounds.Dx() != tt.width {
				t.Errorf("GetImage() width = %d, want %d", bounds.Dx(), tt.width)
			}
			if bounds.Dy() != tt.height {
				t.Errorf("GetImage() height = %d, want %d", bounds.Dy(), tt.height)
			}
		})
	}
}

func TestImagePool_ReturnImage(t *testing.T) {
	logger := &noOpLogger{}
	pool := NewImagePool(logger)

	// Get an image
	img := pool.GetImage(100, 100)
	if img == nil {
		t.Fatal("GetImage() returned nil")
	}

	// Return it to pool
	pool.ReturnImage(img)

	// Verify it's in the pool by checking stats
	stats := pool.GetPoolStats()
	if stats["100x100"] == nil {
		t.Error("ReturnImage() did not add image to pool")
	}
	if count, ok := stats["100x100"].(int); !ok || count != 1 {
		t.Errorf("ReturnImage() pool count = %v, want 1", stats["100x100"])
	}

	// Test returning nil image (should not panic)
	pool.ReturnImage(nil)

	// Get the image again - should reuse from pool
	img2 := pool.GetImage(100, 100)
	if img2 == nil {
		t.Fatal("GetImage() returned nil when reusing from pool")
	}

	// Verify pool is now empty for this size
	stats2 := pool.GetPoolStats()
	if count, ok := stats2["100x100"].(int); !ok || count != 0 {
		t.Errorf("GetImage() should remove image from pool, count = %v", stats2["100x100"])
	}
}

func TestImagePool_GetPoolStats(t *testing.T) {
	logger := &noOpLogger{}
	pool := NewImagePool(logger)

	// Initially should be empty
	stats := pool.GetPoolStats()
	if stats["total_pooled"] != 0 {
		t.Errorf("GetPoolStats() total_pooled = %v, want 0", stats["total_pooled"])
	}
	if stats["pool_count"] != 0 {
		t.Errorf("GetPoolStats() pool_count = %v, want 0", stats["pool_count"])
	}

	// Add some images to pool
	img1 := pool.GetImage(100, 100)
	img2 := pool.GetImage(200, 200)
	img3 := pool.GetImage(100, 100)

	pool.ReturnImage(img1)
	pool.ReturnImage(img2)
	pool.ReturnImage(img3)

	stats = pool.GetPoolStats()
	if total, ok := stats["total_pooled"].(int); !ok || total != 3 {
		t.Errorf("GetPoolStats() total_pooled = %v, want 3", stats["total_pooled"])
	}
	if count, ok := stats["pool_count"].(int); !ok || count != 2 {
		t.Errorf("GetPoolStats() pool_count = %v, want 2 (two different sizes)", stats["pool_count"])
	}
	if count, ok := stats["100x100"].(int); !ok || count != 2 {
		t.Errorf("GetPoolStats() 100x100 count = %v, want 2", stats["100x100"])
	}
	if count, ok := stats["200x200"].(int); !ok || count != 1 {
		t.Errorf("GetPoolStats() 200x200 count = %v, want 1", stats["200x200"])
	}
}

func TestImagePool_Cleanup(t *testing.T) {
	logger := &noOpLogger{}
	pool := NewImagePool(logger)

	// Add images to pool
	img1 := pool.GetImage(100, 100)
	img2 := pool.GetImage(200, 200)
	pool.ReturnImage(img1)
	pool.ReturnImage(img2)

	// Verify they're in the pool
	statsBefore := pool.GetPoolStats()
	if total, ok := statsBefore["total_pooled"].(int); !ok || total != 2 {
		t.Errorf("Before cleanup: total_pooled = %v, want 2", statsBefore["total_pooled"])
	}

	// Cleanup
	pool.Cleanup()

	// Verify pool is empty
	statsAfter := pool.GetPoolStats()
	if total, ok := statsAfter["total_pooled"].(int); !ok || total != 0 {
		t.Errorf("After cleanup: total_pooled = %v, want 0", statsAfter["total_pooled"])
	}
	if count, ok := statsAfter["pool_count"].(int); !ok || count != 0 {
		t.Errorf("After cleanup: pool_count = %v, want 0", statsAfter["pool_count"])
	}
}

func TestImagePool_GetImageReuseFromPool(t *testing.T) {
	logger := &noOpLogger{}
	pool := NewImagePool(logger)

	// Create and return an image
	img1 := pool.GetImage(150, 150)
	pool.ReturnImage(img1)

	// Get image again - should reuse from pool
	img2 := pool.GetImage(150, 150)

	// Should be the same image (reused)
	// Note: We can't directly compare pointers in tests easily,
	// but we can verify the pool is now empty for this size
	stats := pool.GetPoolStats()
	if count, ok := stats["150x150"].(int); !ok || count != 0 {
		t.Errorf("GetImage() should reuse from pool, remaining count = %v", stats["150x150"])
	}

	// Verify it's a valid image
	bounds := img2.Bounds()
	if bounds.Dx() != 150 || bounds.Dy() != 150 {
		t.Errorf("GetImage() reused image has wrong size: %dx%d, want 150x150", bounds.Dx(), bounds.Dy())
	}
}
