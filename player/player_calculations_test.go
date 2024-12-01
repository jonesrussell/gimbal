package player

import (
	"log/slog"
	"math"
	"reflect"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/logger"
	"github.com/solarlune/resolv"
)

func TestPlayer_calculateCoordinates(t *testing.T) {
	type fields struct {
		input     InputHandlerInterface
		angle     float64
		speed     float64
		direction float64
		Object    *resolv.Object
		Sprite    *ebiten.Image
		viewAngle float64
		path      []resolv.Vector
		logger    slog.Logger
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
		want1  float64
	}{
		{
			name: "case 1: default angle and speed",
			fields: fields{
				input:     nil,
				angle:     0,
				speed:     5,
				direction: 0,
				Object:    nil,
				Sprite:    nil,
				viewAngle: 0,
				path:      nil,
				logger:    logger.NewSlogHandler(slog.LevelInfo),
			},
			want:  500.0,
			want1: 232.0,
		},
		{
			name: "case 2: 45 degree angle with speed",
			fields: fields{
				input:     nil,
				angle:     45,
				speed:     5,
				direction: 1,
				Object:    nil,
				Sprite:    nil,
				viewAngle: 45,
				path:      nil,
				logger:    logger.NewSlogHandler(slog.LevelInfo),
			},
			want:  414.0,
			want1: 79.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := &Player{
				PlayerInput: PlayerInput{
					input: tt.fields.input,
				},
				PlayerPosition: PlayerPosition{
					Object: tt.fields.Object,
				},
				PlayerSprite: PlayerSprite{
					Sprite: tt.fields.Sprite,
				},
				PlayerPath: PlayerPath{
					path: tt.fields.path,
				},
				viewAngle: tt.fields.viewAngle,
				direction: tt.fields.direction,
				angle:     tt.fields.angle,
			}
			got, got1 := player.calculateCoordinates()
			const epsilon = 0.000001
			if math.Abs(got-tt.want) > epsilon {
				t.Errorf("Player.calculateCoordinates() got = %v, want %v", got, tt.want)
			}
			if math.Abs(got1-tt.want1) > epsilon {
				t.Errorf("Player.calculateCoordinates() got1 = %v, want1 %v", got1, tt.want1)
			}
		})
	}
}

func TestPlayer_calculatePosition(t *testing.T) {
	type fields struct {
		input     InputHandlerInterface
		angle     float64
		speed     float64
		direction float64
		Object    *resolv.Object
		Sprite    *ebiten.Image
		viewAngle float64
		path      []resolv.Vector
	}
	tests := []struct {
		name   string
		fields fields
		want   resolv.Vector
	}{
		{
			name: "case 1: initial position",
			fields: fields{
				input:     nil,
				angle:     0,
				speed:     0,
				direction: 0,
				Object:    resolv.NewObject(0, 0, 16, 16),
				Sprite:    nil,
				viewAngle: 0,
				path:      nil,
			},
			want: resolv.Vector{X: 500.00, Y: 232.00},
		},
		{
			name: "case 2: movement with speed and no angle",
			fields: fields{
				input:     nil,
				angle:     0,
				speed:     5,
				direction: 0,
				Object:    resolv.NewObject(0, 0, 16, 16),
				Sprite:    nil,
				viewAngle: 0,
				path:      nil,
			},
			want: resolv.Vector{X: 500.00, Y: 232.00},
		},
		{
			name: "case 3: movement with angle and speed",
			fields: fields{
				input:     nil,
				angle:     45,
				speed:     5,
				direction: 0,
				Object:    resolv.NewObject(0, 0, 16, 16),
				Sprite:    nil,
				viewAngle: 0,
				path:      nil,
			},
			want: resolv.Vector{X: 500.00, Y: 232.00}, // Considering simple trigonometric calculations without actual game physics
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := &Player{
				PlayerInput: PlayerInput{
					input: tt.fields.input,
				},
				PlayerPosition: PlayerPosition{
					Object: tt.fields.Object,
				},
				PlayerSprite: PlayerSprite{
					Sprite: tt.fields.Sprite,
				},
				PlayerPath: PlayerPath{
					path: tt.fields.path,
				},
				viewAngle: tt.fields.viewAngle,
				direction: tt.fields.direction,
				angle:     tt.fields.angle,
			}
			got := player.calculatePosition()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Player.calculatePosition() = %v, want %v", got, tt.want)
			}
		})
	}
}
