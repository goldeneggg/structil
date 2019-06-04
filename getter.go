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
	Get(name string) interface{}
	GetString(name string) string
	GetInt64(name string) int64
	GetUint64(name string) uint64
	GetFloat64(name string) float64
	GetBool(name string) bool
	IsString(name string) bool
	IsInt64(name string) bool
	IsUint64(name string) bool
	IsFloat64(name string) bool
	IsBool(name string) bool
	IsStruct(name string) bool
	IsSlice(name string) bool
	IsMap(name string) bool
	IsFunc(name string) bool
	IsChan(name string) bool
	IsInterface(name string) bool
	MapGet(name string, f func(int, Getter) interface{}) ([]interface{}, error)
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

func (g *gImpl) Get(name string) interface{} {
	_, ok := g.cachedI[name]
	if !ok {
		g.cache(name)
	}

	return g.cachedI[name]
}

func (g *gImpl) getRV(name string) reflect.Value {
	_, ok := g.cachedRV[name]
	if !ok {
		// TODO: non-indirectを取り扱うか決める
		g.cache(name)
	}

	return g.cachedRV[name]
}

func (g *gImpl) cache(name string) {
	frv := g.rv.FieldByName(name)
	frv = reflect.Indirect(frv)
	g.cachedRV[name] = frv

	g.cachedI[name] = toI(frv)
}

func (g *gImpl) GetString(name string) string {
	// TODO: reflect.Value has String() method because it implements the Stringer interface.
	// So this method does not occur panic.
	return g.getRV(name).String()
}

func (g *gImpl) GetInt64(name string) int64 {
	return g.getRV(name).Int()
}

func (g *gImpl) GetUint64(name string) uint64 {
	return g.getRV(name).Uint()
}

func (g *gImpl) GetFloat64(name string) float64 {
	return g.getRV(name).Float()
}

func (g *gImpl) GetBool(name string) bool {
	return g.getRV(name).Bool()
}

func (g *gImpl) IsString(name string) bool {
	return g.is(name, reflect.String)
}

func (g *gImpl) IsInt64(name string) bool {
	return g.is(name, reflect.Int)
}

func (g *gImpl) IsUint64(name string) bool {
	return g.is(name, reflect.Uint)
}

func (g *gImpl) IsFloat64(name string) bool {
	return g.is(name, reflect.Float64)
}

func (g *gImpl) IsBool(name string) bool {
	return g.is(name, reflect.Bool)
}

func (g *gImpl) IsStruct(name string) bool {
	return g.is(name, reflect.Struct)
}

func (g *gImpl) IsSlice(name string) bool {
	return g.is(name, reflect.Slice)
}

func (g *gImpl) IsMap(name string) bool {
	return g.is(name, reflect.Map)
}

func (g *gImpl) IsFunc(name string) bool {
	return g.is(name, reflect.Func)
}

func (g *gImpl) IsChan(name string) bool {
	return g.is(name, reflect.Chan)
}

func (g *gImpl) IsInterface(name string) bool {
	return g.is(name, reflect.Interface)
}

func (g *gImpl) is(name string, exp reflect.Kind) bool {
	frv := g.getRV(name)
	return frv.Kind() == exp
}

func (g *gImpl) MapGet(name string, f func(int, Getter) interface{}) ([]interface{}, error) {
	if !g.IsSlice(name) {
		return nil, fmt.Errorf("field %s is not slice", name)
	}

	var vi reflect.Value
	var ac Getter
	var err error
	var res []interface{}
	srv := g.getRV(name)

	for i := 0; i < srv.Len(); i++ {
		vi = srv.Index(i)
		ac, err = NewGetter(toI(vi))
		if err != nil {
			res = append(res, nil)
			continue
		}

		res = append(res, f(i, ac))
	}

	return res, nil
}

// TODO: candidates of moving to utils
func toI(rv reflect.Value) interface{} {
	if rv.IsValid() && rv.CanInterface() {
		return rv.Interface()
	} else {
		return nil
	}
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
func compareStructure(i1 interface{}, i2 interface{}) bool {
	// TODO: 2つのstructの構造比較
	return false
}
