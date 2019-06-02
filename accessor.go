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
	GetRV(name string) (reflect.Value, error)
	Get(name string) (interface{}, error)
	GetString(name string) (string, error)
	GetInt(name string) (int, error)
	GetFloat64(name string) (float64, error)
	GetBool(name string) (bool, error)
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
	// TODO st == nilならすぐエラーでいい

	rv := reflect.ValueOf(st)
	kind := rv.Kind()

	if kind != reflect.Ptr && kind != reflect.Struct {
		return nil, fmt.Errorf("%v is not supported kind", kind)
	}

	if kind == reflect.Ptr {
		if rv.IsNil() {
			return nil, fmt.Errorf("value of passed argument %+v is nil", rv)
		}
		// TODO: ポインタの場合、入力元内容変更の影響を受けるので対応検討
		// （並行処理で使った場合の考慮etc）
		rv = reflect.Indirect(rv)
	}

	return &aImpl{
		rv:       rv,
		cachedRV: map[string]reflect.Value{},
		cachedI:  map[string]interface{}{},
	}, nil
}

func (a *aImpl) GetRV(name string) (reflect.Value, error) {
	return a.getRV(name, true)
}

func (a *aImpl) getRV(name string, isIndirect bool) (reflect.Value, error) {
	frv, ok := a.cachedRV[name]
	if !ok {
		var err error
		frv, err = a.recacheRV(name, isIndirect)
		if err != nil {
			return reflect.ValueOf(nil), err
		}
	}

	return frv, nil
}

func (a *aImpl) recacheRV(name string, isIndirect bool) (reflect.Value, error) {
	frv := a.rv.FieldByName(name)
	kind := frv.Kind()
	if kind == reflect.Invalid {
		return reflect.ValueOf(nil), fmt.Errorf("name %s is invalid. frv: %+v", name, frv)
	}

	if isIndirect {
		frv = reflect.Indirect(frv)
	}
	a.cachedRV[name] = frv

	return frv, nil
}

func (a *aImpl) Get(name string) (interface{}, error) {
	intf, ok := a.cachedI[name]
	if !ok {
		var err error
		intf, err = a.recacheI(name)
		if err != nil {
			return nil, err
		}
	}

	return intf, nil
}

func (a *aImpl) recacheI(name string) (interface{}, error) {
	frv, err := a.GetRV(name)
	if err != nil {
		return nil, err
	}

	var intf interface{}
	if frv.IsValid() && frv.CanInterface() {
		intf = frv.Interface()
	}
	a.cachedI[name] = intf

	return intf, nil
}

func (a *aImpl) GetString(name string) (string, error) {
	intf, err := a.Get(name)
	if err != nil {
		return initStr, err
	}

	res, ok := intf.(string)
	if !ok {
		return initStr, fmt.Errorf("field %s is not string %+v", name, intf)
	}

	return res, nil
}

func (a *aImpl) GetInt(name string) (int, error) {
	intf, err := a.Get(name)
	if err != nil {
		return initInt, err
	}

	res, ok := intf.(int)
	if !ok {
		return initInt, fmt.Errorf("field %s is not int %+v", name, intf)
	}

	return res, nil
}

func (a *aImpl) GetFloat64(name string) (float64, error) {
	intf, err := a.Get(name)
	if err != nil {
		return initInt, err
	}

	res, ok := intf.(float64)
	if !ok {
		return initFloat64, fmt.Errorf("field %s is not float64 %+v", name, intf)
	}

	return res, nil
}

func (a *aImpl) GetBool(name string) (bool, error) {
	intf, err := a.Get(name)
	if err != nil {
		return initBool, err
	}

	res, ok := intf.(bool)
	if !ok {
		return initBool, fmt.Errorf("field %s is not bool %+v", name, intf)
	}

	return res, nil
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
	frv, err := a.GetRV(name)
	if err != nil {
		return false
	}

	return frv.Kind() == exp
}

func (a *aImpl) MapStructs(name string, f func(int, Accessor) interface{}) ([]interface{}, error) {
	if !a.IsSlice(name) {
		return nil, fmt.Errorf("field %s is not slice", name)
	}

	srv, err := a.GetRV(name)
	if err != nil {
		return nil, err
	}

	var vi reflect.Value
	var ac Accessor
	var res []interface{}

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
