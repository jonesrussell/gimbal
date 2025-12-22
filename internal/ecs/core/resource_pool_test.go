package core

import (
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

