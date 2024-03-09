package game

import (
	"log/slog"
	"reflect"
	"testing"

	"github.com/jonesrussell/gimbal/internal/logger"
)

func TestNewGimlarGame(t *testing.T) {
	tests := []struct {
		name    string
		speed   float64
		want    *GimlarGame
		wantErr bool
	}{
		{
			name:  "valid speed",
			speed: 100.0,
			want: &GimlarGame{
				speed:  100.0,
				logger: logger.NewSlogHandler(slog.LevelDebug), // Assuming you want to test with debug level
			},
			wantErr: false,
		},
		{
			name:    "invalid speed",
			speed:   -100.0,
			want:    &GimlarGame{},
			wantErr: true,
		},
		// New test cases
		{
			name:    "zero speed",
			speed:   0.0,
			want:    &GimlarGame{},
			wantErr: false,
		},
		{
			name:    "small positive speed",
			speed:   0.1,
			want:    &GimlarGame{},
			wantErr: false,
		},
		{
			name:    "large positive speed",
			speed:   1000.0,
			want:    &GimlarGame{},
			wantErr: false,
		},
		{
			name:    "small negative speed",
			speed:   -0.1,
			want:    &GimlarGame{},
			wantErr: false,
		},
		{
			name:    "large negative speed",
			speed:   -1000.0,
			want:    &GimlarGame{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGimlarGame(tt.speed)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGimlarGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if g.speed != tt.want.speed || !reflect.DeepEqual(g.logger, tt.want.logger) {
				t.Errorf("NewGimlarGame() = %v, want %v", g, tt.want)
			}
		})
	}
}
