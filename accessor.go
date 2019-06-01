package structil

import (
	"fmt"
	"reflect"
)

const (
	errStr  = ""
	errInt  = -1
	errBool = false
)

type Accessor interface {
	GetRV(name string) (reflect.Value, error)
	Get(name string) (interface{}, error)
	GetString(name string) (string, error)
	GetInt(name string) (int, error)
	GetBool(name string) (bool, error)
	IsStruct(name string) bool
	IsSlice(name string) bool
	MapStructs(name string, f func(int, Accessor) interface{}) ([]interface{}, error)
}

type accessorImpl struct {
	rv      reflect.Value
	cachedV map[string]reflect.Value
	cachedI map[string]interface{}
}

func NewAccessor(st interface{}) (Accessor, error) {
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

	return &accessorImpl{
		rv:      rv,
		cachedV: map[string]reflect.Value{},
		cachedI: map[string]interface{}{},
	}, nil
}

func (a *accessorImpl) GetRV(name string) (reflect.Value, error) {
	return a.getRV(name, true)
}

func (a *accessorImpl) getRV(name string, isIndirect bool) (reflect.Value, error) {
	frv, ok := a.cachedV[name]
	if !ok {
		frv = a.rv.FieldByName(name)
		kind := frv.Kind()
		if kind == reflect.Invalid {
			return reflect.ValueOf(nil), fmt.Errorf("name %s is invalid. frv: %+v", name, frv)
		}

		if isIndirect {
			frv = reflect.Indirect(frv)
		}
		a.cachedV[name] = frv
	}

	return frv, nil
}

func (a *accessorImpl) Get(name string) (interface{}, error) {
	intf, ok := a.cachedI[name]
	if !ok {
		frv, err := a.GetRV(name)
		if err != nil {
			return nil, err
		}

		intf = frv.Interface()
		a.cachedI[name] = intf
	}

	return intf, nil
}

func (a *accessorImpl) GetString(name string) (string, error) {
	intf, err := a.Get(name)
	if err != nil {
		return errStr, err
	}

	res, ok := intf.(string)
	if !ok {
		return errStr, fmt.Errorf("field %s is not string %+v", name, intf)
	}

	return res, nil
}

func (a *accessorImpl) GetInt(name string) (int, error) {
	intf, err := a.Get(name)
	if err != nil {
		return errInt, err
	}

	res, ok := intf.(int)
	if !ok {
		return errInt, fmt.Errorf("field %s is not int %+v", name, intf)
	}

	return res, nil
}

func (a *accessorImpl) GetBool(name string) (bool, error) {
	intf, err := a.Get(name)
	if err != nil {
		return errBool, err
	}

	res, ok := intf.(bool)
	if !ok {
		return errBool, fmt.Errorf("field %s is not bool %+v", name, intf)
	}

	return res, nil
}

func (a *accessorImpl) IsStruct(name string) bool {
	return a.is(name, reflect.Struct)
}

func (a *accessorImpl) IsSlice(name string) bool {
	return a.is(name, reflect.Slice)
}

func (a *accessorImpl) is(name string, exp reflect.Kind) bool {
	frv, err := a.GetRV(name)
	if err != nil {
		return false
	}

	kind := frv.Kind()
	return kind == exp
}

func (a *accessorImpl) MapStructs(name string, f func(int, Accessor) interface{}) ([]interface{}, error) {
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
