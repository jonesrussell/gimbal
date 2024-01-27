package game

import (
	"testing"

	input "github.com/quasilyte/ebitengine-input"
)

func TestNewGimlarGame(t *testing.T) {
	type args struct {
		speed float64
	}
	tests := []struct {
		name string
		args args
		want *GimlarGame
	}{
		{
			name: "Test with speed 1.0",
			args: args{speed: 1.0},
			want: &GimlarGame{
				speed: 1.0,
				// Add other fields as needed
			},
		},
		{
			name: "Test with speed 2.0",
			args: args{speed: 2.0},
			want: &GimlarGame{
				speed: 2.0,
				// Add other fields as needed
			},
		},
		// Add more test cases as needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := NewGimlarGame(tt.args.speed)
			if got.speed != tt.want.speed {
				t.Errorf("NewGimlarGame().speed = %v, want %v", got.speed, tt.want.speed)
			}
			// Compare other fields as needed
		})
	}
}

func TestGimlarGame_Run(t *testing.T) {
	type fields struct {
		p           *Player
		inputSystem input.System
		speed       float64
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test with valid game",
			fields: fields{
				p:           &Player{},      // Replace with a valid Player
				inputSystem: input.System{}, // Replace with a valid input.System
				speed:       1.0,
			},
			wantErr: false,
		},
		// Add more test cases as needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHandler := &MockHandler{}
			player, err := NewPlayer(mockHandler, tt.fields.speed)
			if err != nil {
				t.Fatalf("Failed to create new player: %v", err)
			}
			g := &GimlarGame{
				p:           player,
				inputSystem: tt.fields.inputSystem,
				speed:       tt.fields.speed,
			}
			if err := g.Run(); (err != nil) != tt.wantErr {
				t.Errorf("GimlarGame.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
