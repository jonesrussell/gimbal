package game

import (
	"fmt"
	_ "image/png"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

func TestNewPlayer(t *testing.T) {
	type args struct {
		input InputHandlerInterface
		speed float64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test with valid input and speed",
			args: args{
				input: &MockHandler{}, // Use MockHandler
				speed: 1.0,
			},
			wantErr: false,
		},
		{
			name: "Test with nil input",
			args: args{
				input: nil,
				speed: 1.0,
			},
			wantErr: true,
		},
		// Add more test cases as needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image := ebiten.NewImage(600, 480)
			_, err := NewPlayer(tt.args.input, tt.args.speed, &Debugger{}, image)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPlayer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlayer_Update(t *testing.T) {
	type fields struct {
		input     InputHandlerInterface
		speed     float64
		angle     float64
		direction float64
		Object    *resolv.Object
	}

	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "Test with MoveLeft action",
			fields: fields{
				input:     NewMockHandler(), // Use MockHandler
				speed:     1.0,
				angle:     0.0,
				direction: 0.0,
				Object:    resolv.NewObject(0, 0, 20, 20),
			},
			want: -1.0,
		},
		{
			name: "Test with MoveRight action",
			fields: fields{
				input:     NewMockHandler(), // Use MockHandler
				speed:     1.0,
				angle:     0.0,
				direction: 0.0,
				Object:    resolv.NewObject(0, 0, 20, 20),
			},
			want: 1.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image := ebiten.NewImage(600, 480)
			p, err := NewPlayer(tt.fields.input, tt.fields.speed, &Debugger{}, image)
			if err != nil {
				t.Fatalf("Failed to create new player: %v", err)
			}
			p.angle = tt.fields.angle
			p.direction = tt.fields.direction
			p.Object = tt.fields.Object
			switch tt.name {
			case "Test with MoveLeft action":
				tt.fields.input.(*MockHandler).PressKey(ebiten.KeyLeft)
			case "Test with MoveRight action":
				tt.fields.input.(*MockHandler).PressKey(ebiten.KeyRight)
			}
			p.Update()
			if p.direction != tt.want {
				fmt.Println(tt.name)
				t.Errorf("Player.Update() direction = %v, want %v", p.direction, tt.want)
			}
		})
	}
}
