package logger

import "reflect"

// equalValues compares two values for equality
func equalValues(a, b any) bool {
	if a == nil || b == nil {
		return a == b
	}

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	if va.Type() != vb.Type() {
		return false
	}

	return compareByKind(va, vb)
}

// compareByKind delegates comparison based on reflect.Kind
func compareByKind(va, vb reflect.Value) bool {
	switch va.Kind() {
	case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128, reflect.String, reflect.UnsafePointer:
		return comparePrimitive(va, vb)
	case reflect.Slice, reflect.Array:
		return equalSlicesOrArrays(va, vb)
	case reflect.Map:
		return equalMaps(va, vb)
	case reflect.Struct:
		return equalStructs(va, vb)
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Pointer:
		return comparePrimitive(va, vb)
	default:
		return false
	}
}

// comparePrimitive compares primitive types using direct interface comparison
func comparePrimitive(va, vb reflect.Value) bool {
	return va.Interface() == vb.Interface()
}

func equalSlicesOrArrays(va, vb reflect.Value) bool {
	if va.Len() != vb.Len() {
		return false
	}
	for i := 0; i < va.Len(); i++ {
		if !equalValues(va.Index(i).Interface(), vb.Index(i).Interface()) {
			return false
		}
	}
	return true
}

func equalMaps(va, vb reflect.Value) bool {
	if va.Len() != vb.Len() {
		return false
	}
	for _, k := range va.MapKeys() {
		if !equalValues(va.MapIndex(k).Interface(), vb.MapIndex(k).Interface()) {
			return false
		}
	}
	return true
}

func equalStructs(va, vb reflect.Value) bool {
	for i := 0; i < va.NumField(); i++ {
		if !equalValues(va.Field(i).Interface(), vb.Field(i).Interface()) {
			return false
		}
	}
	return true
}

