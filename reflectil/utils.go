package reflectil

import (
	"reflect"
)

// ToI returns a converted interface from rv.
func ToI(rv reflect.Value) interface{} {
	if rv.IsValid() && rv.CanInterface() {
		return rv.Interface()
	}
	return nil
}

// ElemTypeOf returns a element Type from i.
// If the i's type's Kind is Array, Chan, Map, Ptr, or Slice then this returns the element Type.
// Otherwise this returns i's original Type.
func ElemTypeOf(i interface{}) reflect.Type {
	if i == nil {
		return nil
	}

	t := reflect.TypeOf(i)
	k := t.Kind()

	switch k {
	case reflect.Array, reflect.Chan, reflect.Ptr, reflect.Slice:
		return t.Elem()
	default:
		return t
	}
}

/*
// Note: Publicize candidate
func elemOf(i interface{}) reflect.Value {
	v := reflect.Indirect(reflect.ValueOf(i))
	k := v.Kind()

	if k == reflect.Invalid {
		return v
	}

	if k == reflect.Interface {
		return v.Elem()
	}

	return v
}

// Note: Publicize candidate
func settableOf(i interface{}) reflect.Value {
	// i's Kind must be Interface or Ptr(if else, occur panic)
	return reflect.ValueOf(i).Elem()
}

// Note: Publicize candidate
func clone(i interface{}) interface{} {
	return reflect.Indirect(reflect.ValueOf(i)).Interface()
}

// Note: Publicize candidate
func newSettable(typ reflect.Type) reflect.Value {
	return reflect.New(typ).Elem()
}

// Note: Publicize candidate
func isImplements(i interface{}, t interface{}) bool {
	typ := reflect.TypeOf(t).Elem()
	return isImplementsType(i, typ)
}

func isImplementsType(i interface{}, typ reflect.Type) bool {
	v := reflect.ValueOf(i)
	return typ.Implements(v.Type())
}

func genericsTypeOf() reflect.Type {
	return reflect.TypeOf((*interface{})(nil)).Elem()
}

func newGenericsSettable() reflect.Value {
	return newSettable(genericsTypeOf())
}

func privateFieldValueOf(i interface{}, name string) reflect.Value {
	sv := settableOf(i)
	f := sv.FieldByName(name)
	return reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
}
*/
