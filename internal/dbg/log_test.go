package dbg

import (
	"bytes"
	"os"
	"testing"
)

func TestLog_outputFormat(t *testing.T) {
	Enable()
	defer Disable()
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
	defer Disable()
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

func TestTrace_lifecycle(t *testing.T) {
	ClearTrace()
	defer ClearTrace()

	if TraceRequested() {
		t.Fatal("TraceRequested should be false initially")
	}

	Trace()
	if !TraceRequested() {
		t.Fatal("TraceRequested should be true after Trace()")
	}

	ClearTrace()
	if TraceRequested() {
		t.Fatal("TraceRequested should be false after ClearTrace()")
	}

	// Repeated ClearTrace calls must be idempotent.
	ClearTrace()
	if TraceRequested() {
		t.Fatal("TraceRequested should still be false after second ClearTrace()")
	}
}

func TestTrace_oneFrameSemantics(t *testing.T) {
	Disable()
	defer Disable()
	ClearTrace()
	defer ClearTrace()

	Trace()

	// Simulate frame start: enable logging for this frame if trace was requested.
	wasEnabled := IsEnabled()
	if TraceRequested() {
		Enable()
	}

	if !IsEnabled() {
		t.Fatal("logging should be enabled during the trace frame")
	}

	// Simulate frame end: restore pre-trace state and clear the flag.
	if TraceRequested() {
		if !wasEnabled {
			Disable()
		}
		ClearTrace()
	}

	if IsEnabled() {
		t.Fatal("logging should be disabled after trace frame when it was disabled before")
	}
	if TraceRequested() {
		t.Fatal("trace flag should be cleared after frame end")
	}
}
