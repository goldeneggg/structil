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

// TODO: prettize error logging if error
type Getter interface {
	GetRV(name string) reflect.Value
	Get(name string) interface{}
	GetString(name string) string
	GetInt64(name string) int64
	GetFloat64(name string) float64
	GetBool(name string) bool
	IsStruct(name string) bool
	IsSlice(name string) bool
	IsInterface(name string) bool
	MapStructs(name string, f func(int, Getter) interface{}) ([]interface{}, error)
}

// TODO: implement common panic handler
// if non-exist name assigned, suggest nealy name and pretty error print
type gImpl struct {
	rv       reflect.Value
	cachedRV map[string]reflect.Value
	cachedI  map[string]interface{}
}

func NewGetter(i interface{}) (Getter, error) {
	if i == nil {
		return nil, fmt.Errorf("value of passed argument %+v is nil", i)
	}

	rv := reflect.ValueOf(i)
	kind := rv.Kind()

	// Invalid kind is handled here too.
	if kind != reflect.Ptr && kind != reflect.Struct {
		return nil, fmt.Errorf("%v is not supported kind", kind)
	}

	if kind == reflect.Ptr {
		if rv.IsNil() {
			return nil, fmt.Errorf("value of passed argument %+v is nil", rv)
		}
		// TODO: maybe require syncrhonization control when i is pointer?
		rv = reflect.Indirect(rv)
	}

	return &gImpl{
		rv:       rv,
		cachedRV: map[string]reflect.Value{},
		cachedI:  map[string]interface{}{},
	}, nil
}

// TODO: map => struct => Getter
func NewGetterFromMap(m map[string]interface{}) (Getter, error) {
	return nil, nil
}

// TODO: non-indirectを取り扱うか決める
func (g *gImpl) GetRV(name string) reflect.Value {
	return g.getRV(name, true)
}

func (g *gImpl) getRV(name string, isIndirect bool) reflect.Value {
	frv, ok := g.cachedRV[name]
	if !ok {
		frv = g.recacheRV(name, isIndirect)
	}

	return frv
}

func (g *gImpl) recacheRV(name string, isIndirect bool) reflect.Value {
	frv := g.rv.FieldByName(name)
	if isIndirect {
		frv = reflect.Indirect(frv)
	}
	g.cachedRV[name] = frv

	return g.cachedRV[name]
}

func (g *gImpl) Get(name string) interface{} {
	intf, ok := g.cachedI[name]
	if !ok {
		intf = g.recacheI(name)
	}

	return intf
}

func (g *gImpl) recacheI(name string) interface{} {
	frv := g.GetRV(name)
	if frv.IsValid() && frv.CanInterface() {
		g.cachedI[name] = frv.Interface()
	} else {
		g.cachedI[name] = nil
	}

	return g.cachedI[name]
}

func (g *gImpl) GetString(name string) string {
	return g.GetRV(name).String()
}

func (g *gImpl) GetInt64(name string) int64 {
	return g.GetRV(name).Int()
}

func (g *gImpl) GetFloat64(name string) float64 {
	return g.GetRV(name).Float()
}

func (g *gImpl) GetBool(name string) bool {
	return g.GetRV(name).Bool()
}

func (g *gImpl) IsStruct(name string) bool {
	return g.is(name, reflect.Struct)
}

func (g *gImpl) IsSlice(name string) bool {
	return g.is(name, reflect.Slice)
}

func (g *gImpl) IsInterface(name string) bool {
	return g.is(name, reflect.Interface)
}

func (g *gImpl) is(name string, exp reflect.Kind) bool {
	frv := g.GetRV(name)
	return frv.Kind() == exp
}

func (g *gImpl) MapStructs(name string, f func(int, Getter) interface{}) ([]interface{}, error) {
	if !g.IsSlice(name) {
		return nil, fmt.Errorf("field %s is not slice", name)
	}

	var vi reflect.Value
	var ac Getter
	var err error
	var res []interface{}
	srv := g.GetRV(name)

	for i := 0; i < srv.Len(); i++ {
		vi = srv.Index(i)
		ac, err = NewGetter(vi.Interface())
		if err != nil {
			res = append(res, nil)
			continue
		}

		res = append(res, f(i, ac))
	}

	return res, nil
}

// TODO: candidates of moving to utils
func newSettable(typ reflect.Type) reflect.Value {
	return reflect.New(typ).Elem()
}

// TODO: candidates of moving to utils
func settableOf(i interface{}) reflect.Value {
	// i's Kind must be Interface or Ptr(if else, occur panic)
	return reflect.ValueOf(i).Elem()
}

// TODO: candidates of moving to utils
// 2つのstructの構造比較
func compareStructure(i1 interface{}, i2 interface{}) bool {
	return false
}
