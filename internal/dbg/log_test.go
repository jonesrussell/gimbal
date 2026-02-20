package dbg

import (
	"bytes"
	"os"
	"testing"
)

func TestLog_outputFormat(t *testing.T) {
	Enable()
	defer Enable() // leave enabled for other tests
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = old }()

	Log(Event, "test %s", "message")
	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	out := buf.String()
	if out != "[EVENT] test message\n" {
		t.Errorf("Log output = %q, want [EVENT] test message\n", out)
	}
}

func TestLog_respectsDisable(t *testing.T) {
	Disable()
	defer Enable()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = old }()

	Log(Event, "should not appear")
	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	if buf.Len() != 0 {
		t.Errorf("Log should produce no output when disabled, got %q", buf.String())
	}
}
