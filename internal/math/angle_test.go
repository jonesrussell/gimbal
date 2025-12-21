package math

import (
	"math"
	"testing"
)

func TestAngle_ToRadians(t *testing.T) {
	tests := []struct {
		name  string
		input Angle
		want  float64
	}{
		{"zero degrees", Angle(0), 0},
		{"90 degrees", Angle(90), math.Pi / 2},
		{"180 degrees", Angle(180), math.Pi},
		{"270 degrees", Angle(270), 3 * math.Pi / 2},
		{"360 degrees", Angle(360), 2 * math.Pi},
		{"45 degrees", Angle(45), math.Pi / 4},
		{"negative angle", Angle(-90), -math.Pi / 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.ToRadians()
			if math.Abs(got-tt.want) > 1e-10 {
				t.Errorf("ToRadians() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromRadians(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  Angle
	}{
		{"zero radians", 0, Angle(0)},
		{"pi/2 radians", math.Pi / 2, Angle(90)},
		{"pi radians", math.Pi, Angle(180)},
		{"3pi/2 radians", 3 * math.Pi / 2, Angle(270)},
		{"2pi radians", 2 * math.Pi, Angle(360)},
		{"pi/4 radians", math.Pi / 4, Angle(45)},
		{"negative radians", -math.Pi / 2, Angle(-90)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromRadians(tt.input)
			if math.Abs(float64(got-tt.want)) > 1e-10 {
				t.Errorf("FromRadians() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAngle_Add(t *testing.T) {
	tests := []struct {
		name   string
		angle1 Angle
		angle2 Angle
		want   Angle
	}{
		{"add zero", Angle(45), Angle(0), Angle(45)},
		{"add positive", Angle(30), Angle(60), Angle(90)},
		{"add negative", Angle(90), Angle(-30), Angle(60)},
		{"add negative result", Angle(30), Angle(-60), Angle(-30)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.angle1.Add(tt.angle2)
			if got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAngle_Sub(t *testing.T) {
	tests := []struct {
		name   string
		angle1 Angle
		angle2 Angle
		want   Angle
	}{
		{"subtract zero", Angle(45), Angle(0), Angle(45)},
		{"subtract positive", Angle(90), Angle(30), Angle(60)},
		{"subtract negative", Angle(60), Angle(-30), Angle(90)},
		{"subtract negative result", Angle(30), Angle(90), Angle(-60)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.angle1.Sub(tt.angle2)
			if got != tt.want {
				t.Errorf("Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAngle_Mul(t *testing.T) {
	tests := []struct {
		name   string
		angle  Angle
		scalar float64
		want   Angle
	}{
		{"multiply by zero", Angle(45), 0, Angle(0)},
		{"multiply by one", Angle(45), 1, Angle(45)},
		{"multiply by two", Angle(45), 2, Angle(90)},
		{"multiply by half", Angle(90), 0.5, Angle(45)},
		{"multiply negative angle", Angle(-45), 2, Angle(-90)},
		{"multiply by negative scalar", Angle(45), -1, Angle(-45)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.angle.Mul(tt.scalar)
			if math.Abs(float64(got-tt.want)) > 1e-10 {
				t.Errorf("Mul() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAngle_Div(t *testing.T) {
	tests := []struct {
		name   string
		angle  Angle
		scalar float64
		want   Angle
	}{
		{"divide by one", Angle(90), 1, Angle(90)},
		{"divide by two", Angle(90), 2, Angle(45)},
		{"divide by half", Angle(45), 0.5, Angle(90)},
		{"divide negative angle", Angle(-90), 2, Angle(-45)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.angle.Div(tt.scalar)
			if math.Abs(float64(got-tt.want)) > 1e-10 {
				t.Errorf("Div() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAngle_Normalize(t *testing.T) {
	tests := []struct {
		name  string
		input Angle
		want  Angle
	}{
		{"positive angle within range", Angle(45), Angle(45)},
		{"zero angle", Angle(0), Angle(0)},
		{"360 degrees", Angle(360), Angle(0)},
		{"negative angle", Angle(-90), Angle(270)},
		{"negative angle multiple", Angle(-450), Angle(270)},
		{">360 angle", Angle(450), Angle(90)},
		{"720 degrees", Angle(720), Angle(0)},
		{"1080 degrees", Angle(1080), Angle(0)},
		{"large negative", Angle(-360), Angle(0)},
		{"boundary case 359", Angle(359), Angle(359)},
		{"boundary case -1", Angle(-1), Angle(359)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Normalize()
			if got != tt.want {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAngle_Validate(t *testing.T) {
	tests := []struct {
		name  string
		input Angle
		want  Angle
	}{
		{"valid angle", Angle(90), Angle(90)},
		{"zero angle", Angle(0), Angle(0)},
		{"180 degrees", Angle(180), Angle(180)},
		{"270 degrees (max)", Angle(270), Angle(270)},
		{"angle below min", Angle(-190), FromRadians(MinAngle)},
		{"angle above max", Angle(280), FromRadians(MaxAngle)},
		{"very negative angle", Angle(-1000), FromRadians(MinAngle)},
		{"very large angle", Angle(1000), FromRadians(MaxAngle)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Validate()
			// For clamped values, check they're within valid range
			rad := got.ToRadians()
			if rad < MinAngle || rad > MaxAngle {
				t.Errorf("Validate() returned angle %v (rad: %v) outside valid range [%v, %v]",
					got, rad, MinAngle, MaxAngle)
			}
			// For exact matches, verify exact value
			if tt.input.ToRadians() >= MinAngle && tt.input.ToRadians() <= MaxAngle {
				if got != tt.want {
					t.Errorf("Validate() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestAngleConstants(t *testing.T) {
	if AngleUp != 0 {
		t.Errorf("AngleUp = %v, want 0", AngleUp)
	}
	if AngleRight != 90 {
		t.Errorf("AngleRight = %v, want 90", AngleRight)
	}
	if AngleDown != 180 {
		t.Errorf("AngleDown = %v, want 180", AngleDown)
	}
	if AngleLeft != 270 {
		t.Errorf("AngleLeft = %v, want 270", AngleLeft)
	}
	if DegreesInCircle != 360.0 {
		t.Errorf("DegreesInCircle = %v, want 360.0", DegreesInCircle)
	}
}

func TestAngleRoundTrip(t *testing.T) {
	// Test that converting to radians and back preserves the value
	tests := []struct {
		name  string
		input Angle
	}{
		{"zero", Angle(0)},
		{"45 degrees", Angle(45)},
		{"90 degrees", Angle(90)},
		{"180 degrees", Angle(180)},
		{"270 degrees", Angle(270)},
		{"360 degrees", Angle(360)},
		{"negative", Angle(-45)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rad := tt.input.ToRadians()
			got := FromRadians(rad)
			if math.Abs(float64(got-tt.input)) > 1e-10 {
				t.Errorf("Round trip failed: input %v -> radians %v -> angle %v", tt.input, rad, got)
			}
		})
	}
}
