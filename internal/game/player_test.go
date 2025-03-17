package game_test

import (
	_ "image/png"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/game"
	"github.com/solarlune/resolv"
)

func TestNewPlayer(t *testing.T) {
	type args struct {
		input  game.InputHandlerInterface
		config *game.GameConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test with valid input and config",
			args: args{
				input:  &game.MockHandler{},
				config: game.DefaultConfig(),
			},
			wantErr: false,
		},
		{
			name: "Test with nil input",
			args: args{
				input:  nil,
				config: game.DefaultConfig(),
			},
			wantErr: true,
		},
		{
			name: "Test with nil config",
			args: args{
				input:  &game.MockHandler{},
				config: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image := ebiten.NewImage(600, 480)
			_, err := game.NewPlayer(tt.args.input, tt.args.config, image)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPlayer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlayer_Update(t *testing.T) {
	type fields struct {
		input     game.InputHandlerInterface
		config    *game.GameConfig
		angle     float64
		direction float64
		Object    *resolv.ConvexPolygon
	}

	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "Test with MoveLeft action",
			fields: fields{
				input:     game.NewMockHandler(),
				config:    game.DefaultConfig(),
				angle:     0.0,
				direction: 0.0,
				Object:    resolv.NewRectangle(0, 0, 20, 20),
			},
			want: -1.0,
		},
		{
			name: "Test with MoveRight action",
			fields: fields{
				input:     game.NewMockHandler(),
				config:    game.DefaultConfig(),
				angle:     0.0,
				direction: 0.0,
				Object:    resolv.NewRectangle(0, 0, 20, 20),
			},
			want: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image := ebiten.NewImage(600, 480)
			p, err := game.NewPlayer(tt.fields.input, tt.fields.config, image)
			if err != nil {
				t.Fatalf("Failed to create new player: %v", err)
			}
			p.SetAngle(tt.fields.angle)
			p.SetDirection(tt.fields.direction)
			p.Object = tt.fields.Object
			switch tt.name {
			case "Test with MoveLeft action":
				tt.fields.input.(*game.MockHandler).PressKey(ebiten.KeyLeft)
			case "Test with MoveRight action":
				tt.fields.input.(*game.MockHandler).PressKey(ebiten.KeyRight)
			}
			p.Update()
			if p.GetDirection() != tt.want {
				t.Errorf("Player.Update() direction = %v, want %v", p.GetDirection(), tt.want)
			}
		})
	}
}
