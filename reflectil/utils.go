package reflectil

import (
	"reflect"
	"unsafe"
)

// ToI returns a converted interface from rv.
func ToI(rv reflect.Value) interface{} {
	if rv.IsValid() && rv.CanInterface() {
		return rv.Interface()
	}
	return nil
}

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
