package resolv

import (
	"github.com/jonesrussell/gimbal/internal/components"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

func Add(space *donburi.Entry, objects ...*donburi.Entry) {
	for _, obj := range objects {
		components.Space.Get(space).Add(GetObject(obj))
	}
}

func SetObject(entry *donburi.Entry, obj *resolv.Object) {
	components.Object.Set(entry, obj)
}

func GetObject(entry *donburi.Entry) *resolv.Object {
	return components.Object.Get(entry)
}
