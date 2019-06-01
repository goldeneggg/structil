package structil

import "reflect"

func IElemOf(i interface{}) reflect.Value {
	typ := reflect.TypeOf(i)
	k := typ.Kind()
	if k == reflect.Ptr {
		return reflect.ValueOf(i).Elem()
	} else {
		return reflect.ValueOf(&i).Elem()
	}
}
