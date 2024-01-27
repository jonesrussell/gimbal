package game

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

func TestNewPlayer(t *testing.T) {
	type args struct {
		input HandlerInterface
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
			_, err := NewPlayer(tt.args.input, tt.args.speed)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPlayer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Rest of the file...

func TestPlayer_Update(t *testing.T) {
	type fields struct {
		input     HandlerInterface
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
				input:     &HandlerWrapper{&MockHandler{}}, // Wrap MockHandler with HandlerWrapper
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
				input:     &HandlerWrapper{&MockHandler{}}, // Wrap MockHandler with HandlerWrapper
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
			p := &Player{
				input:     tt.fields.input,
				speed:     tt.fields.speed,
				angle:     tt.fields.angle,
				direction: tt.fields.direction,
				Object:    tt.fields.Object,
			}
			switch tt.name {
			case "Test with MoveLeft action":
				p.input.ActionIsPressed(ActionMoveLeft)
			case "Test with MoveRight action":
				p.input.ActionIsPressed(ActionMoveRight)
			}
			p.Update()
			if p.direction != tt.want {
				t.Errorf("Player.Update() direction = %v, want %v", p.direction, tt.want)
			}
		})
	}
}

func TestPlayer_Draw(t *testing.T) {
	type fields struct {
		input     HandlerInterface
		speed     float64
		angle     float64
		direction float64
		Object    *resolv.Object
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test with valid Player",
			fields: fields{
				input:     &MockHandler{}, // Use MockHandler
				speed:     1.0,
				angle:     0.0,
				direction: 0.0,
				Object:    resolv.NewObject(0, 0, 20, 20),
			},
		},
		// Add more test cases as needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{
				input:     tt.fields.input,
				speed:     tt.fields.speed,
				angle:     tt.fields.angle,
				direction: tt.fields.direction,
				Object:    tt.fields.Object,
			}
			img := ebiten.NewImage(20, 20) // Create a new image
			p.Draw(img)
			// Add assertions to verify the drawing operation
		})
	}
}
