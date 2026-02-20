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

func Enable()  { atomic.StoreUint32(&enabled, 1) }
func Disable() { atomic.StoreUint32(&enabled, 0) }

func Log(cat Category, msg string, args ...any) {
	if atomic.LoadUint32(&enabled) == 0 {
		return
	}
	fmt.Printf("[%s] %s\n", cat, fmt.Sprintf(msg, args...))
}
