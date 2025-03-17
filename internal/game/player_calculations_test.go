package game

import (
	"log/slog"
	"reflect"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/logger"
	"github.com/solarlune/resolv"
)

func TestPlayer_calculateCoordinates(t *testing.T) {
	type fields struct {
		input     InputHandlerInterface
		angle     float64
		speed     float64
		direction float64
		Object    *resolv.ConvexPolygon
		Sprite    *ebiten.Image
		viewAngle float64
		path      []resolv.Vector
		logger    slog.Logger
		config    *GameConfig
	}
	type args struct {
		angle float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
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
				config:    DefaultConfig(),
			},
			args: args{
				angle: 0,
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
				config:    DefaultConfig(),
			},
			args: args{
				angle: 45,
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
				config:    tt.fields.config,
			}
			got, got1 := player.calculateCoordinates(tt.args.angle)
			if got != tt.want {
				t.Errorf("Player.calculateCoordinates() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Player.calculateCoordinates() got1 = %v, want %v", got1, tt.want1)
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
		Object    *resolv.ConvexPolygon
		Sprite    *ebiten.Image
		viewAngle float64
		path      []resolv.Vector
		config    *GameConfig
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
				Object:    resolv.NewRectangle(0, 0, 16, 16),
				Sprite:    nil,
				viewAngle: 0,
				path:      nil,
				config:    DefaultConfig(),
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
				Object:    resolv.NewRectangle(0, 0, 16, 16),
				Sprite:    nil,
				viewAngle: 0,
				path:      nil,
				config:    DefaultConfig(),
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
				Object:    resolv.NewRectangle(0, 0, 16, 16),
				Sprite:    nil,
				viewAngle: 0,
				path:      nil,
				config:    DefaultConfig(),
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
				config:    tt.fields.config,
			}
			got := player.calculatePosition()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Player.calculatePosition() = %v, want %v", got, tt.want)
			}
		})
	}
}
