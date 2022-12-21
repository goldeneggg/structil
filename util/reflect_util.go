package util

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// ToI returns a converted interface from rv.
func ToI(rv reflect.Value) any {
	if rv.IsValid() && rv.CanInterface() {
		return rv.Interface()
	}
	return nil
}

// ElemTypeOf returns a element Type from i.
// If the i's type's Kind is Array, Chan, Map, Ptr, or Slice then this returns the element Type.
// Otherwise this returns i's original Type.
func ElemTypeOf(i any) reflect.Type {
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

// RecoverToError returns an error converted from recoverd panic information.
func RecoverToError(r any) (err error) {
	if r != nil {
		msg := fmt.Sprintf("\n%v\n", r) + stackTrace()
		err = fmt.Errorf("unexpected panic occurred: %s", msg)
	}

	return
}

func stackTrace() string {
	msgs := make([]string, 0, 10)

	for d := 0; ; d++ {
		pc, file, line, ok := runtime.Caller(d)
		if !ok {
			break
		}
		msgs = append(msgs, fmt.Sprintf(" -> %d: %s: %s:%d", d, runtime.FuncForPC(pc).Name(), file, line))
	}

	return strings.Join(msgs, "\n")
}

/*
// Note: Publicize candidate
func elemOf(i any) reflect.Value {
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
func settableOf(i any) reflect.Value {
	// i's Kind must be Interface or Ptr(if else, occur panic)
	return reflect.ValueOf(i).Elem()
}

// Note: Publicize candidate
func clone(i any) any {
	return reflect.Indirect(reflect.ValueOf(i)).Interface()
}

// Note: Publicize candidate
func newSettable(typ reflect.Type) reflect.Value {
	return reflect.New(typ).Elem()
}

// Note: Publicize candidate
func isImplements(i any, t any) bool {
	typ := reflect.TypeOf(t).Elem()
	return isImplementsType(i, typ)
}

func isImplementsType(i any, typ reflect.Type) bool {
	v := reflect.ValueOf(i)
	return typ.Implements(v.Type())
}

func genericsTypeOf() reflect.Type {
	return reflect.TypeOf((*any)(nil)).Elem()
}

func newGenericsSettable() reflect.Value {
	return newSettable(genericsTypeOf())
}

func privateFieldValueOf(i any, name string) reflect.Value {
	sv := settableOf(i)
	f := sv.FieldByName(name)
	return reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
}
*/
