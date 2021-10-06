package structil

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/goldeneggg/structil/util"
)

// Getter is the struct that wraps the basic Getter method.
type Getter struct {
	rv     reflect.Value            // Value of input interface
	numf   int                      // Field nums
	names  []string                 // Field names
	hases  map[string]bool          // Field existing condition map of struct fields
	types  map[string]reflect.Type  // Type map of struct fields
	values map[string]reflect.Value // Value map of indirected struct fields
	intfs  map[string]interface{}   // interface map of struct fields
	cached map[string]bool
}

// NewGetter returns a concrete Getter that uses and obtains from i.
// i must be a struct or struct pointer.
func NewGetter(i interface{}) (*Getter, error) {
	rv := reflect.ValueOf(i)
	kind := rv.Kind()

	if kind != reflect.Ptr && kind != reflect.Struct {
		return nil, fmt.Errorf("%+v is not supported kind: %v. value: %+v", i, kind, rv)
	}

	if kind == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}

	if !rv.IsValid() {
		return nil, fmt.Errorf("%+v is invalid argument. value: %+v", i, rv)
	}

	numf := rv.NumField()

	return &Getter{
		rv:     rv,
		numf:   numf,
		names:  make([]string, 0, numf),
		hases:  map[string]bool{},
		values: map[string]reflect.Value{},
		types:  map[string]reflect.Type{},
		intfs:  map[string]interface{}{},
		cached: map[string]bool{},
	}, nil
}

// NumField returns num of struct field.
func (g *Getter) NumField() int {
	return g.numf
}

// Has tests whether the original struct has a field named name arg.
func (g *Getter) Has(name string) bool {
	if _, ok := g.cached[name]; !ok {
		g.cache(name)
	}

	return g.hases[name]
}

func (g *Getter) cache(name string) {
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
	g.intfs[name] = util.ToI(frv)
	g.cached[name] = true
}

// Names returns names of struct field.
func (g *Getter) Names() []string {
	// to setup g.names is run only once
	if g.numf > 0 && len(g.names) == 0 {
		var sf reflect.StructField
		for i := 0; i < g.numf; i++ {
			sf = g.rv.Type().Field(i)
			g.names = append(g.names, sf.Name)
		}
	}

	return g.names
}

// GetType returns the reflect.Type object of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
func (g *Getter) GetType(name string) (reflect.Type, bool) {
	if _, ok := g.cached[name]; !ok {
		g.cache(name)
	}

	return g.types[name], g.hases[name]
}

// GetValue returns the reflect.Value object of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
func (g *Getter) GetValue(name string) (reflect.Value, bool) {
	if _, ok := g.cached[name]; !ok {
		g.cache(name)
	}

	return g.values[name], g.hases[name]
}

// Get returns the interface of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
func (g *Getter) Get(name string) (interface{}, bool) {
	if _, ok := g.cached[name]; !ok {
		g.cache(name)
	}

	return g.intfs[name], g.hases[name]
}

// Bool returns the byte of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not bool.
func (g *Getter) Bool(name string) (bool, bool) {
	v, has := g.Get(name)
	if !has {
		return false, has
	}

	res, ok := v.(bool)
	return res, ok
}

// Byte returns the byte of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not byte.
func (g *Getter) Byte(name string) (byte, bool) {
	v, has := g.Get(name)
	if !has {
		return 0, has
	}

	res, ok := v.(byte)
	return res, ok
}

// Bytes returns the []byte of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not []byte.
func (g *Getter) Bytes(name string) ([]byte, bool) {
	v, has := g.Get(name)
	if !has {
		return nil, has
	}

	res, ok := v.([]byte)
	return res, ok
}

// String returns the string of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not string.
func (g *Getter) String(name string) (string, bool) {
	v, has := g.Get(name)
	if !has {
		return "", has
	}

	res, ok := v.(string)
	return res, ok
}

// Int returns the int of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not int.
func (g *Getter) Int(name string) (int, bool) {
	v, has := g.Get(name)
	if !has {
		return 0, has
	}

	res, ok := v.(int)
	return res, ok
}

// Int8 returns the int8 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not int8.
func (g *Getter) Int8(name string) (int8, bool) {
	v, has := g.Get(name)
	if !has {
		return 0, has
	}

	res, ok := v.(int8)
	return res, ok
}

// Int16 returns the int16 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not int16.
func (g *Getter) Int16(name string) (int16, bool) {
	v, has := g.Get(name)
	if !has {
		return 0, has
	}

	res, ok := v.(int16)
	return res, ok
}

// Int32 returns the int32 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not int32.
func (g *Getter) Int32(name string) (int32, bool) {
	v, has := g.Get(name)
	if !has {
		return 0, has
	}

	res, ok := v.(int32)
	return res, ok
}

// Int64 returns the int64 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not int64.
func (g *Getter) Int64(name string) (int64, bool) {
	v, has := g.Get(name)
	if !has {
		return 0, has
	}

	res, ok := v.(int64)
	return res, ok
}

// Uint returns the uint of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not uint.
func (g *Getter) Uint(name string) (uint, bool) {
	v, has := g.Get(name)
	if !has {
		return 0, has
	}

	res, ok := v.(uint)
	return res, ok
}

// Uint8 returns the uint8 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not uint8.
func (g *Getter) Uint8(name string) (uint8, bool) {
	v, has := g.Get(name)
	if !has {
		return 0, has
	}

	res, ok := v.(uint8)
	return res, ok
}

// Uint16 returns the uint16 of the original struct field named name.Getter
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not uint16.
func (g *Getter) Uint16(name string) (uint16, bool) {
	v, has := g.Get(name)
	if !has {
		return 0, has
	}

	res, ok := v.(uint16)
	return res, ok
}

// Uint32 returns the uint32 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not uint32.
func (g *Getter) Uint32(name string) (uint32, bool) {
	v, has := g.Get(name)
	if !has {
		return 0, has
	}

	res, ok := v.(uint32)
	return res, ok
}

// Uint64 returns the uint64 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not uint64.
func (g *Getter) Uint64(name string) (uint64, bool) {
	v, has := g.Get(name)
	if !has {
		return 0, has
	}

	res, ok := v.(uint64)
	return res, ok
}

// Uintptr returns the uintptr of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not uintptr.
func (g *Getter) Uintptr(name string) (uintptr, bool) {
	v, has := g.Get(name)
	if !has {
		return 0, has
	}

	res, ok := v.(uintptr)
	return res, ok
}

// Float32 returns the float32 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not float32.
func (g *Getter) Float32(name string) (float32, bool) {
	v, has := g.Get(name)
	if !has {
		return 0, has
	}

	res, ok := v.(float32)
	return res, ok
}

// Float64 returns the float64 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not float64.
func (g *Getter) Float64(name string) (float64, bool) {
	v, has := g.Get(name)
	if !has {
		return 0, has
	}

	res, ok := v.(float64)
	return res, ok
}

// Complex64 returns the complex64 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not complex64.
func (g *Getter) Complex64(name string) (complex64, bool) {
	v, has := g.Get(name)
	if !has {
		return 0, has
	}

	res, ok := v.(complex64)
	return res, ok
}

// Complex128 returns the complex128 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not complex128.
func (g *Getter) Complex128(name string) (complex128, bool) {
	v, has := g.Get(name)
	if !has {
		return 0, has
	}

	res, ok := v.(complex128)
	return res, ok
}

// UnsafePointer returns the unsafe.Pointer of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not unsafe.Pointer.
func (g *Getter) UnsafePointer(name string) (unsafe.Pointer, bool) {
	v, has := g.Get(name)
	if !has {
		return nil, has
	}

	res, ok := v.(unsafe.Pointer)
	return res, ok
}

// Slice returns the slice of interface of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not slice of interface.
func (g *Getter) Slice(name string) ([]interface{}, bool) {
	if !g.IsSlice(name) {
		return nil, false
	}

	// See: https://golang.org/doc/faq#convert_slice_of_interface
	v, _ := g.GetValue(name)
	iSlice := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		iSlice[i] = v.Index(i).Interface()
	}
	return iSlice, true
}

// IsByte reports whether type of the original struct field named name is byte.
func (g *Getter) IsByte(name string) bool {
	return g.is(name, reflect.Uint8)
}

// IsBytes reports whether type of the original struct field named name is []byte.
func (g *Getter) IsBytes(name string) bool {
	if g.IsSlice(name) {
		t, ok := g.GetType(name)
		if ok && t.Elem().Kind() == reflect.Uint8 {
			return true
		}
	}
	return false
}

// IsString reports whether type of the original struct field named name is string.
func (g *Getter) IsString(name string) bool {
	return g.is(name, reflect.String)
}

// IsInt reports whether type of the original struct field named name is int.
func (g *Getter) IsInt(name string) bool {
	return g.is(name, reflect.Int)
}

// IsInt8 reports whether type of the original struct field named name is int8.
func (g *Getter) IsInt8(name string) bool {
	return g.is(name, reflect.Int8)
}

// IsInt16 reports whether type of the original struct field named name is int16.
func (g *Getter) IsInt16(name string) bool {
	return g.is(name, reflect.Int16)
}

// IsInt32 reports whether type of the original struct field named name is int32.
func (g *Getter) IsInt32(name string) bool {
	return g.is(name, reflect.Int32)
}

// IsInt64 reports whether type of the original struct field named name is int64.
func (g *Getter) IsInt64(name string) bool {
	return g.is(name, reflect.Int64)
}

// IsUint reports whether type of the original struct field named name is uint.
func (g *Getter) IsUint(name string) bool {
	return g.is(name, reflect.Uint)
}

// IsUint8 reports whether type of the original struct field named name is uint8.
func (g *Getter) IsUint8(name string) bool {
	return g.is(name, reflect.Uint8)
}

// IsUint16 reports whether type of the original struct field named name is uint16.
func (g *Getter) IsUint16(name string) bool {
	return g.is(name, reflect.Uint16)
}

// IsUint32 reports whether type of the original struct field named name is uint32.
func (g *Getter) IsUint32(name string) bool {
	return g.is(name, reflect.Uint32)
}

// IsUint64 reports whether type of the original struct field named name is uint64.
func (g *Getter) IsUint64(name string) bool {
	return g.is(name, reflect.Uint64)
}

// IsUintptr reports whether type of the original struct field named name is uintptr.
func (g *Getter) IsUintptr(name string) bool {
	return g.is(name, reflect.Uintptr)
}

// IsFloat32 reports whether type of the original struct field named name is float32.
func (g *Getter) IsFloat32(name string) bool {
	return g.is(name, reflect.Float32)
}

// IsFloat64 reports whether type of the original struct field named name is float64.
func (g *Getter) IsFloat64(name string) bool {
	return g.is(name, reflect.Float64)
}

// IsBool reports whether type of the original struct field named name is bool.
func (g *Getter) IsBool(name string) bool {
	return g.is(name, reflect.Bool)
}

// IsComplex64 reports whether type of the original struct field named name is []byte.
func (g *Getter) IsComplex64(name string) bool {
	return g.is(name, reflect.Complex64)
}

// IsComplex128 reports whether type of the original struct field named name is []byte.
func (g *Getter) IsComplex128(name string) bool {
	return g.is(name, reflect.Complex128)
}

// IsUnsafePointer reports whether type of the original struct field named name is []byte.
func (g *Getter) IsUnsafePointer(name string) bool {
	return g.is(name, reflect.UnsafePointer)
}

// IsMap reports whether type of the original struct field named name is map.
func (g *Getter) IsMap(name string) bool {
	return g.is(name, reflect.Map)
}

// IsFunc reports whether type of the original struct field named name is func.
func (g *Getter) IsFunc(name string) bool {
	return g.is(name, reflect.Func)
}

// IsChan reports whether type of the original struct field named name is chan.
func (g *Getter) IsChan(name string) bool {
	return g.is(name, reflect.Chan)
}

// IsStruct reports whether type of the original struct field named name is struct.
func (g *Getter) IsStruct(name string) bool {
	return g.is(name, reflect.Struct)
}

// IsSlice reports whether type of the original struct field named name is slice.
func (g *Getter) IsSlice(name string) bool {
	return g.is(name, reflect.Slice)
}

// IsArray reports whether type of the original struct field named name is slice.
func (g *Getter) IsArray(name string) bool {
	return g.is(name, reflect.Array)
}

func (g *Getter) is(name string, exp reflect.Kind) bool {
	if !g.Has(name) {
		return false
	}

	return g.values[name].Kind() == exp
}

// MapGet returns the interface slice of mapped values of the original struct field named name.
func (g *Getter) MapGet(name string, f func(int, *Getter) (interface{}, error)) ([]interface{}, error) {
	if !g.IsSlice(name) {
		return nil, fmt.Errorf("field %s does not exist or is not slice type", name)
	}

	var vi reflect.Value
	var eg *Getter
	var err error
	var r interface{}

	srv, _ := g.GetValue(name)
	res := make([]interface{}, srv.Len())

	for i := 0; i < srv.Len(); i++ {
		vi = srv.Index(i)
		eg, err = NewGetter(util.ToI(vi))
		if err != nil {
			return nil, err
		}

		r, err = f(i, eg)
		if err != nil {
			return nil, err
		}

		res[i] = r
	}

	return res, nil
}
