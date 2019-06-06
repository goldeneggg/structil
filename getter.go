package structil

import (
	"fmt"
	"reflect"

	"github.com/goldeneggg/structil/reflectil"
)

type Getter interface {
	Has(name string) bool
	GetType(name string) reflect.Type
	GetValue(name string) reflect.Value
	Get(name string) interface{}
	Bytes(name string) []byte
	String(name string) string
	Int64(name string) int64
	Uint64(name string) uint64
	Float64(name string) float64
	Bool(name string) bool
	IsBytes(name string) bool
	IsString(name string) bool
	IsInt64(name string) bool
	IsUint64(name string) bool
	IsFloat64(name string) bool
	IsBool(name string) bool
	IsMap(name string) bool
	IsFunc(name string) bool
	IsChan(name string) bool
	IsStruct(name string) bool
	IsSlice(name string) bool
	MapGet(name string, f func(int, Getter) interface{}) ([]interface{}, error)
}

type gImpl struct {
	rv     reflect.Value // Value of input interface
	hases  map[string]bool
	types  map[string]reflect.Type  // Type map of struct fields
	values map[string]reflect.Value // Value map of indirected struct fields
	intfs  map[string]interface{}   // interface map of struct fields
}

func NewGetter(i interface{}) (Getter, error) {
	if i == nil {
		return nil, fmt.Errorf("value of passed argument %+v is nil", i)
	}

	rv := reflect.ValueOf(i)
	kind := rv.Kind()

	if kind != reflect.Ptr && kind != reflect.Struct {
		return nil, fmt.Errorf("%v is not supported kind", kind)
	}

	if kind == reflect.Ptr {
		if rv.IsNil() {
			return nil, fmt.Errorf("value of passed argument %+v is nil", rv)
		}

		// indirect is required when kind is Ptr
		rv = reflect.Indirect(rv)
	}

	return &gImpl{
		rv:     rv,
		hases:  map[string]bool{},
		values: map[string]reflect.Value{},
		types:  map[string]reflect.Type{},
		intfs:  map[string]interface{}{},
	}, nil
}

func (g *gImpl) Has(name string) bool {
	_, ok := g.hases[name]
	if !ok {
		g.cache(name)
	}

	return g.hases[name]
}

func (g *gImpl) GetType(name string) reflect.Type {
	_, ok := g.types[name]
	if !ok {
		g.cache(name)
	}

	return g.types[name]
}

func (g *gImpl) cache(name string) {
	frv := g.rv.FieldByName(name) // XXX: This code is slow
	if frv.IsValid() {
		g.types[name] = frv.Type()
		g.hases[name] = true
	} else {
		g.types[name] = nil
		g.hases[name] = false
	}

	frv = reflect.Indirect(frv)
	g.values[name] = frv

	g.intfs[name] = reflectil.ToI(frv)
}

func (g *gImpl) GetValue(name string) reflect.Value {
	_, ok := g.values[name]
	if !ok {
		g.cache(name)
	}

	return g.values[name]
}

func (g *gImpl) Get(name string) interface{} {
	_, ok := g.intfs[name]
	if !ok {
		g.cache(name)
	}

	return g.intfs[name]
}

func (g *gImpl) Bytes(name string) []byte {
	return g.GetValue(name).Bytes()
}

func (g *gImpl) String(name string) string {
	// Note:
	// reflect.Value has String() method because it implements the Stringer interface.
	// So this method does not occur panic.
	return g.GetValue(name).String()
}

func (g *gImpl) Int64(name string) int64 {
	return g.GetValue(name).Int()
}

func (g *gImpl) Uint64(name string) uint64 {
	return g.GetValue(name).Uint()
}

func (g *gImpl) Float64(name string) float64 {
	return g.GetValue(name).Float()
}

func (g *gImpl) Bool(name string) bool {
	return g.GetValue(name).Bool()
}

func (g *gImpl) IsBytes(name string) bool {
	return g.IsSlice(name) && g.GetType(name).Elem().Kind() == reflect.Uint8
}

func (g *gImpl) IsString(name string) bool {
	return g.is(name, reflect.String)
}

func (g *gImpl) IsInt64(name string) bool {
	return g.is(name, reflect.Int64)
}

func (g *gImpl) IsUint64(name string) bool {
	return g.is(name, reflect.Uint64)
}

func (g *gImpl) IsFloat64(name string) bool {
	return g.is(name, reflect.Float64)
}

func (g *gImpl) IsBool(name string) bool {
	return g.is(name, reflect.Bool)
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

func (g *gImpl) IsStruct(name string) bool {
	return g.is(name, reflect.Struct)
}

func (g *gImpl) IsSlice(name string) bool {
	return g.is(name, reflect.Slice)
}

func (g *gImpl) is(name string, exp reflect.Kind) bool {
	frv := g.GetValue(name)
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
	srv := g.GetValue(name)

	for i := 0; i < srv.Len(); i++ {
		vi = srv.Index(i)
		ac, err = NewGetter(reflectil.ToI(vi))
		if err != nil {
			res = append(res, nil)
			continue
		}

		res = append(res, f(i, ac))
	}

	return res, nil
}
