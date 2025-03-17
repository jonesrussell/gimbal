package game_test

import (
	"log/slog"
	"reflect"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/game"
	"github.com/jonesrussell/gimbal/internal/logger"
	"github.com/solarlune/resolv"
)

func TestPlayer_CalculateCoordinates(t *testing.T) {
	type fields struct {
		Input     game.InputHandlerInterface
		Angle     float64
		Speed     float64
		Direction float64
		Object    *resolv.ConvexPolygon
		Sprite    *ebiten.Image
		ViewAngle float64
		Path      []resolv.Vector
		Logger    slog.Logger
		Config    *game.GameConfig
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
				Input:     nil,
				Angle:     0,
				Speed:     5,
				Direction: 0,
				Object:    nil,
				Sprite:    nil,
				ViewAngle: 0,
				Path:      nil,
				Logger:    logger.NewSlogHandler(slog.LevelInfo),
				Config:    game.DefaultConfig(),
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
				Input:     nil,
				Angle:     45,
				Speed:     5,
				Direction: 1,
				Object:    nil,
				Sprite:    nil,
				ViewAngle: 45,
				Path:      nil,
				Logger:    logger.NewSlogHandler(slog.LevelInfo),
				Config:    game.DefaultConfig(),
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
			player := &game.Player{
				PlayerInput: game.PlayerInput{
					Input: tt.fields.Input,
				},
				PlayerPosition: game.PlayerPosition{
					Object: tt.fields.Object,
				},
				PlayerSprite: game.PlayerSprite{
					Sprite: tt.fields.Sprite,
				},
				PlayerPath: game.PlayerPath{
					Path: tt.fields.Path,
				},
				ViewAngle: tt.fields.ViewAngle,
				Direction: tt.fields.Direction,
				Angle:     tt.fields.Angle,
				Config:    tt.fields.Config,
			}
			got, got1 := player.CalculateCoordinates(tt.args.angle)
			if got != tt.want {
				t.Errorf("Player.CalculateCoordinates() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Player.CalculateCoordinates() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPlayer_CalculatePosition(t *testing.T) {
	type fields struct {
		Input     game.InputHandlerInterface
		Angle     float64
		Speed     float64
		Direction float64
		Object    *resolv.ConvexPolygon
		Sprite    *ebiten.Image
		ViewAngle float64
		Path      []resolv.Vector
		Config    *game.GameConfig
	}
	tests := []struct {
		name   string
		fields fields
		want   resolv.Vector
	}{
		{
			name: "case 1: initial position",
			fields: fields{
				Input:     nil,
				Angle:     0,
				Speed:     0,
				Direction: 0,
				Object:    resolv.NewRectangle(0, 0, 16, 16),
				Sprite:    nil,
				ViewAngle: 0,
				Path:      nil,
				Config:    game.DefaultConfig(),
			},
			want: resolv.Vector{X: 500.00, Y: 232.00},
		},
		{
			name: "case 2: movement with speed and no angle",
			fields: fields{
				Input:     nil,
				Angle:     0,
				Speed:     5,
				Direction: 0,
				Object:    resolv.NewRectangle(0, 0, 16, 16),
				Sprite:    nil,
				ViewAngle: 0,
				Path:      nil,
				Config:    game.DefaultConfig(),
			},
			want: resolv.Vector{X: 500.00, Y: 232.00},
		},
		{
			name: "case 3: movement with angle and speed",
			fields: fields{
				Input:     nil,
				Angle:     45,
				Speed:     5,
				Direction: 0,
				Object:    resolv.NewRectangle(0, 0, 16, 16),
				Sprite:    nil,
				ViewAngle: 0,
				Path:      nil,
				Config:    game.DefaultConfig(),
			},
			want: resolv.Vector{X: 500.00, Y: 232.00}, // Simple trigonometric calculations
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := &game.Player{
				PlayerInput: game.PlayerInput{
					Input: tt.fields.Input,
				},
				PlayerPosition: game.PlayerPosition{
					Object: tt.fields.Object,
				},
				PlayerSprite: game.PlayerSprite{
					Sprite: tt.fields.Sprite,
				},
				PlayerPath: game.PlayerPath{
					Path: tt.fields.Path,
				},
				ViewAngle: tt.fields.ViewAngle,
				Direction: tt.fields.Direction,
				Angle:     tt.fields.Angle,
				Config:    tt.fields.Config,
			}
			got := player.CalculatePosition()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Player.CalculatePosition() = %v, want %v", got, tt.want)
			}
		})
	}
}
