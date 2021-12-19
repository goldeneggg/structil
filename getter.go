package structil

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"

	"github.com/goldeneggg/structil/util"
)

// Getter is the struct that wraps the basic Getter method.
type Getter struct {
	rv     reflect.Value // Value of input interface (this is struct)
	numf   int           // Field nums
	names  []string      // Field names
	fields map[string]*getterField
	mu     sync.RWMutex
}

// NewGetter returns a concrete Getter that uses and obtains from i.
// i must be a struct or struct pointer.
func NewGetter(i interface{}) (*Getter, error) {
	stVal, err := toStructValue(i)
	if err != nil {
		return nil, err
	}

	g := &Getter{
		rv:   stVal,
		numf: stVal.NumField(),
	}
	g.names = make([]string, 0, g.numf)
	g.fields = make(map[string]*getterField, g.numf)

	for idx := 0; idx < g.numf; idx++ {
		gf := g.newGetterField(idx)
		g.names = append(g.names, gf.name)
		g.fields[gf.name] = gf
	}

	return g, nil
}

// toStructValue returns a reflect.Value that can generate to Getter.
// i must be a struct or struct pointer.
func toStructValue(i interface{}) (reflect.Value, error) {
	rv := reflect.ValueOf(i)
	kind := rv.Kind()
	if kind != reflect.Ptr && kind != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("kind [%v] is not either struct or pointer. i = [%+v]", kind, i)
	}
	if kind == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}
	if !rv.IsValid() {
		return reflect.Value{}, fmt.Errorf("reflect.Value is invalid. i = [%+v]", i)
	}

	return rv, nil
}

type getterField struct {
	name     string
	sFld     reflect.StructField
	typ      reflect.Type
	indirect reflect.Value // is Value via reflect.Indirect(v)
	intf     interface{}
}

func (gf *getterField) isKind(kind reflect.Kind) bool {
	return gf.indirect.Kind() == kind
}

func (g *Getter) newGetterField(idx int) *getterField {
	sFld := g.rv.Type().Field(idx)
	v := g.rv.Field(idx)
	indirect := reflect.Indirect(v)

	return &getterField{
		name:     sFld.Name,
		sFld:     sFld,
		typ:      sFld.Type,
		indirect: indirect,
		intf:     util.ToI(indirect),
	}
}

// NumField returns num of struct field.
func (g *Getter) NumField() int {
	return g.numf
}

// Names returns names of struct field.
func (g *Getter) Names() []string {
	return g.names
}

// goroutine-safely access to a getterField by name
func (g *Getter) getSafely(name string) (*getterField, bool) {
	// FIXME: sync.Onceを使ってリファクタ検討
	g.mu.RLock()
	defer g.mu.RUnlock()

	gf, ok := g.fields[name]
	return gf, ok
}

// goroutine-safely and kind-safely access to a getterField by name
func (g *Getter) getSafelyKindly(name string, kind reflect.Kind) (*getterField, bool) {
	gf, ok := g.getSafely(name)
	return gf, ok && gf.isKind(kind)
}

// Has tests whether the original struct has a field named "name".
func (g *Getter) Has(name string) bool {
	_, ok := g.getSafely(name)
	return ok
}

// GetType returns the reflect.Type object of the original struct field named "name".
// 2nd return value will be false if the original struct does not have a "name" field.
func (g *Getter) GetType(name string) (reflect.Type, bool) {
	gf, ok := g.getSafely(name)
	if ok {
		return gf.typ, true
	}

	return nil, false
}

// GetValue returns the reflect.Value object of the original struct field named "name".
// 2nd return value will be false if the original struct does not have a "name" field.
func (g *Getter) GetValue(name string) (reflect.Value, bool) {
	gf, ok := g.getSafely(name)
	if ok {
		return gf.indirect, true
	}

	return reflect.Value{}, false
}

// Get returns the interface of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
func (g *Getter) Get(name string) (interface{}, bool) {
	gf, ok := g.getSafely(name)
	if ok {
		return gf.intf, true
	}

	return nil, false
}

// ToMap returns a map converted from this Getter.
func (g *Getter) ToMap() map[string]interface{} {
	g.mu.RLock()
	defer g.mu.RUnlock()

	m := make(map[string]interface{})
	for name, gf := range g.fields {
		m[name] = gf.intf
	}

	return m
}

// IsSlice reports whether type of the original struct field named name is slice.
func (g *Getter) IsSlice(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Slice)
	return ok
}

// Slice returns the slice of interface of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not slice of interface.
func (g *Getter) Slice(name string) ([]interface{}, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Slice)
	if !ok {
		return nil, false
	}

	len := gf.indirect.Len()

	// See: https://golang.org/doc/faq#convert_slice_of_interface
	iSlice := make([]interface{}, len)
	for i := 0; i < len; i++ {
		iSlice[i] = gf.indirect.Index(i).Interface()
	}
	return iSlice, true
}

// IsBool reports whether type of the original struct field named name is bool.
func (g *Getter) IsBool(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Bool)
	return ok
}

// Bool returns the byte of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not bool.
func (g *Getter) Bool(name string) (bool, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Bool)
	if !ok {
		return false, false
	}

	res, ok := gf.intf.(bool)
	return res, ok
}

// IsByte reports whether type of the original struct field named name is byte.
func (g *Getter) IsByte(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Uint8)
	return ok
}

// Byte returns the byte of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not byte.
func (g *Getter) Byte(name string) (byte, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Uint8)
	if !ok {
		return 0, false
	}

	res, ok := gf.intf.(byte)
	return res, ok
}

// IsBytes reports whether type of the original struct field named name is []byte.
func (g *Getter) IsBytes(name string) bool {
	gf, ok := g.getSafelyKindly(name, reflect.Slice)
	if !ok {
		return false
	}

	return gf.typ.Elem().Kind() == reflect.Uint8
}

// Bytes returns the []byte of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not []byte.
func (g *Getter) Bytes(name string) ([]byte, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Slice)
	if !ok {
		return nil, false
	}

	res, ok := gf.intf.([]byte)
	return res, ok
}

// IsString reports whether type of the original struct field named name is string.
func (g *Getter) IsString(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.String)
	return ok
}

// String returns the string of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not string.
func (g *Getter) String(name string) (string, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.String)
	if !ok {
		return "", false
	}

	res, ok := gf.intf.(string)
	return res, ok
}

// IsInt reports whether type of the original struct field named name is int.
func (g *Getter) IsInt(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Int)
	return ok
}

// Int returns the int of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not int.
func (g *Getter) Int(name string) (int, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Int)
	if !ok {
		return 0, false
	}

	res, ok := gf.intf.(int)
	return res, ok
}

// IsInt8 reports whether type of the original struct field named name is int8.
func (g *Getter) IsInt8(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Int8)
	return ok
}

// Int8 returns the int8 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not int8.
func (g *Getter) Int8(name string) (int8, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Int8)
	if !ok {
		return 0, false
	}

	res, ok := gf.intf.(int8)
	return res, ok
}

// IsInt16 reports whether type of the original struct field named name is int16.
func (g *Getter) IsInt16(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Int16)
	return ok
}

// Int16 returns the int16 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not int16.
func (g *Getter) Int16(name string) (int16, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Int16)
	if !ok {
		return 0, false
	}

	res, ok := gf.intf.(int16)
	return res, ok
}

// IsInt32 reports whether type of the original struct field named name is int32.
func (g *Getter) IsInt32(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Int32)
	return ok
}

// Int32 returns the int32 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not int32.
func (g *Getter) Int32(name string) (int32, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Int32)
	if !ok {
		return 0, false
	}

	res, ok := gf.intf.(int32)
	return res, ok
}

// IsInt64 reports whether type of the original struct field named name is int64.
func (g *Getter) IsInt64(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Int64)
	return ok
}

// Int64 returns the int64 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not int64.
func (g *Getter) Int64(name string) (int64, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Int64)
	if !ok {
		return 0, false
	}

	res, ok := gf.intf.(int64)
	return res, ok
}

// IsUint reports whether type of the original struct field named name is uint.
func (g *Getter) IsUint(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Uint)
	return ok
}

// Uint returns the uint of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not uint.
func (g *Getter) Uint(name string) (uint, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Uint)
	if !ok {
		return 0, false
	}

	res, ok := gf.intf.(uint)
	return res, ok
}

// IsUint8 reports whether type of the original struct field named name is uint8.
func (g *Getter) IsUint8(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Uint8)
	return ok
}

// Uint8 returns the uint8 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not uint8.
func (g *Getter) Uint8(name string) (uint8, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Uint8)
	if !ok {
		return 0, false
	}

	res, ok := gf.intf.(uint8)
	return res, ok
}

// IsUint16 reports whether type of the original struct field named name is uint16.
func (g *Getter) IsUint16(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Uint16)
	return ok
}

// Uint16 returns the uint16 of the original struct field named name.Getter
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not uint16.
func (g *Getter) Uint16(name string) (uint16, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Uint16)
	if !ok {
		return 0, false
	}

	res, ok := gf.intf.(uint16)
	return res, ok
}

// IsUint32 reports whether type of the original struct field named name is uint32.
func (g *Getter) IsUint32(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Uint32)
	return ok
}

// Uint32 returns the uint32 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not uint32.
func (g *Getter) Uint32(name string) (uint32, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Uint32)
	if !ok {
		return 0, false
	}

	res, ok := gf.intf.(uint32)
	return res, ok
}

// IsUint64 reports whether type of the original struct field named name is uint64.
func (g *Getter) IsUint64(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Uint64)
	return ok
}

// Uint64 returns the uint64 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not uint64.
func (g *Getter) Uint64(name string) (uint64, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Uint64)
	if !ok {
		return 0, false
	}

	res, ok := gf.intf.(uint64)
	return res, ok
}

// IsUintptr reports whether type of the original struct field named name is uintptr.
func (g *Getter) IsUintptr(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Uintptr)
	return ok
}

// Uintptr returns the uintptr of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not uintptr.
func (g *Getter) Uintptr(name string) (uintptr, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Uintptr)
	if !ok {
		return 0, false
	}

	res, ok := gf.intf.(uintptr)
	return res, ok
}

// IsFloat32 reports whether type of the original struct field named name is float32.
func (g *Getter) IsFloat32(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Float32)
	return ok
}

// Float32 returns the float32 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not float32.
func (g *Getter) Float32(name string) (float32, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Float32)
	if !ok {
		return 0, false
	}

	res, ok := gf.intf.(float32)
	return res, ok
}

// IsFloat64 reports whether type of the original struct field named name is float64.
func (g *Getter) IsFloat64(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Float64)
	return ok
}

// Float64 returns the float64 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not float64.
func (g *Getter) Float64(name string) (float64, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Float64)
	if !ok {
		return 0, false
	}

	res, ok := gf.intf.(float64)
	return res, ok
}

// IsComplex64 reports whether type of the original struct field named name is []byte.
func (g *Getter) IsComplex64(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Complex64)
	return ok
}

// Complex64 returns the complex64 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not complex64.
func (g *Getter) Complex64(name string) (complex64, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Complex64)
	if !ok {
		return 0, false
	}

	res, ok := gf.intf.(complex64)
	return res, ok
}

// IsComplex128 reports whether type of the original struct field named name is []byte.
func (g *Getter) IsComplex128(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Complex128)
	return ok
}

// Complex128 returns the complex128 of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not complex128.
func (g *Getter) Complex128(name string) (complex128, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Complex128)
	if !ok {
		return 0, false
	}

	res, ok := gf.intf.(complex128)
	return res, ok
}

// IsUnsafePointer reports whether type of the original struct field named name is []byte.
func (g *Getter) IsUnsafePointer(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.UnsafePointer)
	return ok
}

// UnsafePointer returns the unsafe.Pointer of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not unsafe.Pointer.
func (g *Getter) UnsafePointer(name string) (unsafe.Pointer, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.UnsafePointer)
	if !ok {
		return nil, false
	}

	res, ok := gf.intf.(unsafe.Pointer)
	return res, ok
}

// IsMap reports whether type of the original struct field named name is map.
func (g *Getter) IsMap(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Map)
	return ok
}

// IsFunc reports whether type of the original struct field named name is func.
func (g *Getter) IsFunc(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Func)
	return ok
}

// IsChan reports whether type of the original struct field named name is chan.
func (g *Getter) IsChan(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Chan)
	return ok
}

// IsStruct reports whether type of the original struct field named name is struct.
func (g *Getter) IsStruct(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Struct)
	return ok
}

// IsArray reports whether type of the original struct field named name is slice.
func (g *Getter) IsArray(name string) bool {
	_, ok := g.getSafelyKindly(name, reflect.Array)
	return ok
}

// GetGetter returns the Getter of interface of the original struct field named name.
// 2nd return value will be false if the original struct does not have a "name" field.
// 2nd return value will be false if type of the original struct "name" field is not struct or struct pointer.
func (g *Getter) GetGetter(name string) (*Getter, bool) {
	gf, ok := g.getSafelyKindly(name, reflect.Struct)
	if !ok {
		return nil, false
	}

	g, err := NewGetter(gf.intf)
	return g, err == nil
}

// MapGet returns the interface slice of mapped values of the original struct field named name.
func (g *Getter) MapGet(name string, f func(int, *Getter) (interface{}, error)) ([]interface{}, error) {
	gf, ok := g.getSafelyKindly(name, reflect.Slice)
	if !ok {
		return nil, fmt.Errorf("field %s does not exist or is not slice type", name)
	}

	var vi reflect.Value
	var eg *Getter
	var err error
	var r interface{}

	res := make([]interface{}, gf.indirect.Len())

	for i := 0; i < gf.indirect.Len(); i++ {
		vi = gf.indirect.Index(i)
		eg, err = NewGetter(util.ToI(vi))
		if err != nil {
			return nil, fmt.Errorf("fail NewGetter: %w", err)
		}

		r, err = f(i, eg)
		if err != nil {
			return nil, fmt.Errorf("fail MapGet func: %w", err)
		}

		res[i] = r
	}

	return res, nil
}
