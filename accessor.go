package structil

import (
	"fmt"
	"reflect"
)

const (
	initStr     = ""
	initInt     = 0
	initFloat64 = 0.0
	initBool    = false
)

type Accessor interface {
	GetRV(name string) reflect.Value
	Get(name string) interface{}
	GetString(name string) string
	GetInt64(name string) int64
	GetFloat64(name string) float64
	GetBool(name string) bool
	IsStruct(name string) bool
	IsSlice(name string) bool
	IsInterface(name string) bool
	MapStructs(name string, f func(int, Accessor) interface{}) ([]interface{}, error)
}

type aImpl struct {
	rv       reflect.Value
	cachedRV map[string]reflect.Value
	cachedI  map[string]interface{}
}

func NewAccessor(st interface{}) (Accessor, error) {
	if st == nil {
		return nil, fmt.Errorf("value of passed argument %+v is nil", st)
	}

	rv := reflect.ValueOf(st)
	kind := rv.Kind()

	if kind != reflect.Ptr && kind != reflect.Struct {
		return nil, fmt.Errorf("%v is not supported kind", kind)
	}

	if kind == reflect.Ptr {
		if rv.IsNil() {
			return nil, fmt.Errorf("value of passed argument %+v is nil", rv)
		}
		// TODO: maybe require syncrhonization control when st is pointer?
		rv = reflect.Indirect(rv)
	}

	return &aImpl{
		rv:       rv,
		cachedRV: map[string]reflect.Value{},
		cachedI:  map[string]interface{}{},
	}, nil
}

func (a *aImpl) GetRV(name string) reflect.Value {
	return a.getRV(name, true)
}

func (a *aImpl) getRV(name string, isIndirect bool) reflect.Value {
	frv, ok := a.cachedRV[name]
	if !ok {
		frv = a.recacheRV(name, isIndirect)
	}

	return frv
}

func (a *aImpl) recacheRV(name string, isIndirect bool) reflect.Value {
	frv := a.rv.FieldByName(name)
	if isIndirect {
		frv = reflect.Indirect(frv)
	}
	a.cachedRV[name] = frv

	return a.cachedRV[name]
}

func (a *aImpl) Get(name string) interface{} {
	intf, ok := a.cachedI[name]
	if !ok {
		intf = a.recacheI(name)
	}

	return intf
}

func (a *aImpl) recacheI(name string) interface{} {
	frv := a.GetRV(name)
	if frv.IsValid() && frv.CanInterface() {
		a.cachedI[name] = frv.Interface()
	} else {
		a.cachedI[name] = nil
	}

	return a.cachedI[name]
}

func (a *aImpl) GetString(name string) string {
	return a.GetRV(name).String()
}

func (a *aImpl) GetInt64(name string) int64 {
	return a.GetRV(name).Int()
}

func (a *aImpl) GetFloat64(name string) float64 {
	return a.GetRV(name).Float()
}

func (a *aImpl) GetBool(name string) bool {
	return a.GetRV(name).Bool()
}

func (a *aImpl) IsStruct(name string) bool {
	return a.is(name, reflect.Struct)
}

func (a *aImpl) IsSlice(name string) bool {
	return a.is(name, reflect.Slice)
}

func (a *aImpl) IsInterface(name string) bool {
	return a.is(name, reflect.Interface)
}

func (a *aImpl) is(name string, exp reflect.Kind) bool {
	frv := a.GetRV(name)
	return frv.Kind() == exp
}

func (a *aImpl) MapStructs(name string, f func(int, Accessor) interface{}) ([]interface{}, error) {
	if !a.IsSlice(name) {
		return nil, fmt.Errorf("field %s is not slice", name)
	}

	var vi reflect.Value
	var ac Accessor
	var err error
	var res []interface{}
	srv := a.GetRV(name)

	for i := 0; i < srv.Len(); i++ {
		vi = srv.Index(i)
		ac, err = NewAccessor(vi.Interface())
		if err != nil {
			res = append(res, nil)
			continue
		}

		res = append(res, f(i, ac))
	}

	return res, nil
}
