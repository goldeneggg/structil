package structil

import (
	"fmt"
	"reflect"

	"github.com/goldeneggg/structil/reflectil"
)

// Getter is the interface that wraps the basic Getter method.
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

// NewGetter returns a concrete Getter that uses and obtains from i.
// i must be a struct or struct pointer.
func NewGetter(i interface{}) (Getter, error) {
	if i == nil {
		return nil, fmt.Errorf("value of passed argument %+v is nil.", i)
	}

	rv := reflect.ValueOf(i)
	kind := rv.Kind()

	if kind != reflect.Ptr && kind != reflect.Struct {
		return nil, fmt.Errorf("%v is not supported kind.", kind)
	}

	if kind == reflect.Ptr {
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

// Has tests whether the original struct has a field named "name" arg.
func (g *gImpl) Has(name string) bool {
	_, ok := g.hases[name]
	if !ok {
		g.cache(name)
	}

	return g.hases[name]
}

func (g *gImpl) cache(name string) {
	frv := g.rv.FieldByName(name)
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

// GetType returns the reflect.Type object of the original struct field named "name".
// It panics if the original struct does not have a field named "name".
func (g *gImpl) GetType(name string) reflect.Type {
	g.panicIfNotHave(name)

	_, ok := g.types[name]
	if !ok {
		g.cache(name)
	}

	return g.types[name]
}

func (g *gImpl) panicIfNotHave(name string) {
	if !g.Has(name) {
		panic(fmt.Sprintf("field name %s does not exist in the original struct.", name))
	}
}

// GetValue returns the reflect.Value object of the original struct field named "name".
// It panics if the original struct does not have a field named "name".
func (g *gImpl) GetValue(name string) reflect.Value {
	g.panicIfNotHave(name)

	_, ok := g.values[name]
	if !ok {
		g.cache(name)
	}

	return g.values[name]
}

// Get returns the interface of the original struct field named "name".
// It panics if the original struct does not have a field named "name".
func (g *gImpl) Get(name string) interface{} {
	g.panicIfNotHave(name)

	_, ok := g.intfs[name]
	if !ok {
		g.cache(name)
	}

	return g.intfs[name]
}

// Bytes returns the []byte of the original struct field named "name".
// It panics if the original struct does not have a field named "name".
// It panics if type of the original struct field named "name" is not []byte.
func (g *gImpl) Bytes(name string) []byte {
	if v, ok := g.Get(name).([]byte); ok {
		return v
	} else {
		panic(fmt.Sprintf("field name %s is not []byte type. value: %v", name, g.GetValue(name)))
	}
}

// String returns the string of the original struct field named "name".
// It panics if the original struct does not have a field named "name".
// It panics if type of the original struct field named "name" is not string.
func (g *gImpl) String(name string) string {
	if v, ok := g.Get(name).(string); ok {
		return v
	} else {
		panic(fmt.Sprintf("field name %s is not string type. value: %v", name, g.GetValue(name)))
	}
}

// Int64 returns the int64 of the original struct field named "name".
// It panics if the original struct does not have a field named "name".
// It panics if type of the original struct field named "name" is not int64.
func (g *gImpl) Int64(name string) int64 {
	if v, ok := g.Get(name).(int64); ok {
		return v
	} else {
		panic(fmt.Sprintf("field name %s is not int64 type. value: %v", name, g.GetValue(name)))
	}
}

// Uint64 returns the uint64 of the original struct field named "name".
// It panics if the original struct does not have a field named "name".
// It panics if type of the original struct field named "name" is not uint64.
func (g *gImpl) Uint64(name string) uint64 {
	if v, ok := g.Get(name).(uint64); ok {
		return v
	} else {
		panic(fmt.Sprintf("field name %s is not uint64 type. value: %v", name, g.GetValue(name)))
	}
}

// Float64 returns the float64 of the original struct field named "name".
// It panics if the original struct does not have a field named "name".
// It panics if type of the original struct field named "name" is not float64.
func (g *gImpl) Float64(name string) float64 {
	if v, ok := g.Get(name).(float64); ok {
		return v
	} else {
		panic(fmt.Sprintf("field name %s is not float64 type. value: %v", name, g.GetValue(name)))
	}
}

// Bool returns the bool of the original struct field named "name".
// It panics if the original struct does not have a field named "name".
// It panics if type of the original struct field named "name" is not bool.
func (g *gImpl) Bool(name string) bool {
	if v, ok := g.Get(name).(bool); ok {
		return v
	} else {
		panic(fmt.Sprintf("field name %s is not bool type. value: %v", name, g.GetValue(name)))
	}
}

// IsBytes reports whether type of the original struct field named "name" is []byte.
func (g *gImpl) IsBytes(name string) bool {
	return g.IsSlice(name) && g.GetType(name).Elem().Kind() == reflect.Uint8
}

// IsString reports whether type of the original struct field named "name" is string.
func (g *gImpl) IsString(name string) bool {
	return g.is(name, reflect.String)
}

// IsInt64 reports whether type of the original struct field named "name" is int64.
func (g *gImpl) IsInt64(name string) bool {
	return g.is(name, reflect.Int64)
}

// IsUint64 reports whether type of the original struct field named "name" is uint64.
func (g *gImpl) IsUint64(name string) bool {
	return g.is(name, reflect.Uint64)
}

// IsFloat64 reports whether type of the original struct field named "name" is float64.
func (g *gImpl) IsFloat64(name string) bool {
	return g.is(name, reflect.Float64)
}

// IsBool reports whether type of the original struct field named "name" is bool.
func (g *gImpl) IsBool(name string) bool {
	return g.is(name, reflect.Bool)
}

// IsMap reports whether type of the original struct field named "name" is map.
func (g *gImpl) IsMap(name string) bool {
	return g.is(name, reflect.Map)
}

// IsFunc reports whether type of the original struct field named "name" is func.
func (g *gImpl) IsFunc(name string) bool {
	return g.is(name, reflect.Func)
}

// IsChan reports whether type of the original struct field named "name" is chan.
func (g *gImpl) IsChan(name string) bool {
	return g.is(name, reflect.Chan)
}

// IsStruct reports whether type of the original struct field named "name" is struct.
func (g *gImpl) IsStruct(name string) bool {
	return g.is(name, reflect.Struct)
}

// IsSlice reports whether type of the original struct field named "name" is slice.
func (g *gImpl) IsSlice(name string) bool {
	return g.is(name, reflect.Slice)
}

func (g *gImpl) is(name string, exp reflect.Kind) bool {
	if !g.Has(name) {
		return false
	}

	frv := g.GetValue(name)
	return frv.Kind() == exp
}

// MapGet returns the interface slice of mapped values of the original struct field named "name".
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
