package dbg_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/jonesrussell/gimbal/internal/dbg"
)

func TestLog_outputFormat(t *testing.T) {
	dbg.Enable()
	defer dbg.Disable()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = old }()

	dbg.Log(dbg.Event, "test %s", "message")
	w.Close()
	var buf bytes.Buffer
	if _, readErr := buf.ReadFrom(r); readErr != nil {
		t.Fatalf("ReadFrom: %v", readErr)
	}
	out := buf.String()
	if out != "[EVENT] test message\n" {
		t.Errorf("Log output = %q, want [EVENT] test message\n", out)
	}
}

func TestLog_respectsDisable(t *testing.T) {
	dbg.Disable()
	defer dbg.Disable()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = old }()

	dbg.Log(dbg.Event, "should not appear")
	w.Close()
	var buf bytes.Buffer
	if _, readErr := buf.ReadFrom(r); readErr != nil {
		t.Fatalf("ReadFrom: %v", readErr)
	}
	if buf.Len() != 0 {
		t.Errorf("Log should produce no output when disabled, got %q", buf.String())
	}
}

func TestTrace_lifecycle(t *testing.T) {
	dbg.ClearTrace()
	defer dbg.ClearTrace()

	if dbg.TraceRequested() {
		t.Fatal("TraceRequested should be false initially")
	}

	dbg.Trace()
	if !dbg.TraceRequested() {
		t.Fatal("TraceRequested should be true after Trace()")
	}

	dbg.ClearTrace()
	if dbg.TraceRequested() {
		t.Fatal("TraceRequested should be false after ClearTrace()")
	}

	// Repeated ClearTrace calls must be idempotent.
	dbg.ClearTrace()
	if dbg.TraceRequested() {
		t.Fatal("TraceRequested should still be false after second ClearTrace()")
	}
}

func TestTrace_oneFrameSemantics(t *testing.T) {
	dbg.Disable()
	defer dbg.Disable()
	dbg.ClearTrace()
	defer dbg.ClearTrace()

	dbg.Trace()

	// Simulate frame start: enable logging for this frame if trace was requested.
	wasEnabled := dbg.IsEnabled()
	if dbg.TraceRequested() {
		dbg.Enable()
	}

	if !dbg.IsEnabled() {
		t.Fatal("logging should be enabled during the trace frame")
	}

	// Simulate frame end: restore pre-trace state and clear the flag.
	if dbg.TraceRequested() {
		if !wasEnabled {
			dbg.Disable()
		}
		dbg.ClearTrace()
	}

	if dbg.IsEnabled() {
		t.Fatal("logging should be disabled after trace frame when it was disabled before")
	}
	if dbg.TraceRequested() {
		t.Fatal("trace flag should be cleared after frame end")
	}
}
