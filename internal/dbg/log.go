package dbg

import (
	"fmt"
	"sync/atomic"
)

type Category string

const (
	Event  Category = "EVENT"
	State  Category = "STATE"
	System Category = "SYSTEM"
	World  Category = "WORLD"
	Spawn  Category = "SPAWN"
)

var enabled uint32 = 1
var traceNextFrame uint32

func Enable()  { atomic.StoreUint32(&enabled, 1) }
func Disable() { atomic.StoreUint32(&enabled, 0) }

// Trace enables full dbg output for the next frame only (call Trace(), then one Update() will log; then disabled again).
func Trace() { atomic.StoreUint32(&traceNextFrame, 1) }

// TraceRequested returns whether trace-next-frame was requested.
func TraceRequested() bool { return atomic.LoadUint32(&traceNextFrame) == 1 }

// ClearTrace clears the trace-next-frame flag (called at end of Update after that frame).
func ClearTrace() { atomic.StoreUint32(&traceNextFrame, 0) }

func Log(cat Category, msg string, args ...any) {
	if atomic.LoadUint32(&enabled) == 0 {
		return
	}
	fmt.Printf("[%s] %s\n", cat, fmt.Sprintf(msg, args...))
}
