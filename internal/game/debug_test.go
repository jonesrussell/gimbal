package game

import (
	"reflect"
	"testing"
)

func TestNewDebugger(t *testing.T) {
	tests := []struct {
		name string
		want *Debugger
	}{
		{
			name: "Test Debugger Creation",
			want: &Debugger{
				debug: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDebugger(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDebugger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDebugger_IsDebugMode(t *testing.T) {
	type fields struct {
		debug bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Test Debug Mode Enabled",
			fields: fields{
				debug: true,
			},
			want: true,
		},
		{
			name: "Test Debug Mode Disabled",
			fields: fields{
				debug: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Debugger{
				debug: tt.fields.debug,
			}
			if got := d.IsDebugMode(); got != tt.want {
				t.Errorf("Debugger.IsDebugMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDebugger_DebugMode(t *testing.T) {
	type fields struct {
		debug bool
	}
	type args struct {
		debug bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test Debug Mode Enabled",
			fields: fields{
				debug: false,
			},
			args: args{
				debug: true,
			},
		},
		{
			name: "Test Debug Mode Disabled",
			fields: fields{
				debug: true,
			},
			args: args{
				debug: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Debugger{
				debug: tt.fields.debug,
			}
			d.DebugMode(tt.args.debug)
			if d.debug != tt.args.debug {
				t.Errorf("Debugger.DebugMode() = %v, want %v", d.debug, tt.args.debug)
			}
		})
	}
}
