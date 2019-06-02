package structil

import (
	"fmt"
	"reflect"
)

func IElemOf(i interface{}) reflect.Value {
	v := reflect.Indirect(reflect.ValueOf(i))
	k := v.Kind()

	if k == reflect.Invalid {
		return v
	}

	if k == reflect.Interface {
		return v.Elem()
	} else {
		return v
	}
}

func ISettableOf(ptr interface{}) reflect.Value {
	v := reflect.ValueOf(ptr)
	k := v.Kind()
	if k != reflect.Ptr {
		panic(fmt.Sprintf("ptr should be Ptr but it is %+v", k))
	}

	return v.Elem()
}
