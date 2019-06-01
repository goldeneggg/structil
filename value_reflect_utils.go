package structil

import "reflect"

func VInterface(v reflect.Value) interface{} {
	if VCanInterface(v) {
		return v.Interface()
	} else {
		return nil
	}
}

func VCanInterface(v reflect.Value) bool {
	return v.IsValid() && v.CanInterface()
}

func VCompare(v1 reflect.Value, v2 reflect.Value) bool {
	return VInterface(v1) == VInterface(v2)
}
