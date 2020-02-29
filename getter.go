package structil

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/goldeneggg/structil/reflectil"
)

// Getter is the interface that wraps the basic Getter method.
type Getter interface {
	NumField() int
	Has(name string) bool
	GetType(name string) reflect.Type
	GetValue(name string) reflect.Value
	Get(name string) interface{}
	EGet(name string) (interface{}, error)
	Byte(name string) byte
	Bytes(name string) []byte
	String(name string) string
	Int(name string) int
	Int8(name string) int8
	Int16(name string) int16
	Int32(name string) int32
	Int64(name string) int64
	Uint(name string) uint
	Uint8(name string) uint8
	Uint16(name string) uint16
	Uint32(name string) uint32
	Uint64(name string) uint64
	Uintptr(name string) uintptr
	Float64(name string) float64
	Bool(name string) bool
	Complex64(name string) complex64
	Complex128(name string) complex128
	UnsafePointer(name string) unsafe.Pointer
	IsByte(name string) bool
	IsBytes(name string) bool
	IsString(name string) bool
	IsInt(name string) bool
	IsInt8(name string) bool
	IsInt16(name string) bool
	IsInt32(name string) bool
	IsInt64(name string) bool
	IsUint(name string) bool
	IsUint8(name string) bool
	IsUint16(name string) bool
	IsUint32(name string) bool
	IsUint64(name string) bool
	IsUintptr(name string) bool
	IsFloat64(name string) bool
	IsBool(name string) bool
	IsComplex64(name string) bool
	IsComplex128(name string) bool
	IsUnsafePointer(name string) bool
	IsMap(name string) bool
	IsFunc(name string) bool
	IsChan(name string) bool
	IsStruct(name string) bool
	IsSlice(name string) bool
	IsArray(name string) bool
	MapGet(name string, f func(int, Getter) (interface{}, error)) ([]interface{}, error)
}

// GetterImpl is the default Getter implementation.
type GetterImpl struct {
	rv     reflect.Value // Value of input interface
	numf   int
	hases  map[string]bool
	types  map[string]reflect.Type  // Type map of struct fields
	values map[string]reflect.Value // Value map of indirected struct fields
	intfs  map[string]interface{}   // interface map of struct fields
}

// NewGetter returns a concrete Getter that uses and obtains from i.
// i must be a struct or struct pointer.
func NewGetter(i interface{}) (Getter, error) {
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

	return &GetterImpl{
		rv:     rv,
		numf:   rv.NumField(),
		hases:  map[string]bool{},
		values: map[string]reflect.Value{},
		types:  map[string]reflect.Type{},
		intfs:  map[string]interface{}{},
	}, nil
}

// NumField returns num of struct field.
func (g *GetterImpl) NumField() int {
	return g.numf
}

// Has tests whether the original struct has a field named name arg.
func (g *GetterImpl) Has(name string) bool {
	_, ok := g.hases[name]
	if !ok {
		g.cache(name)
	}

	return g.hases[name]
}

func (g *GetterImpl) cache(name string) {
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

// GetType returns the reflect.Type object of the original struct field named name.
// It panics if the original struct does not have a field named name.
func (g *GetterImpl) GetType(name string) reflect.Type {
	g.panicIfNotHave(name)

	_, ok := g.types[name]
	if !ok {
		g.cache(name)
	}

	return g.types[name]
}

func (g *GetterImpl) panicIfNotHave(name string) {
	if !g.Has(name) {
		panic(fmt.Sprintf("field name %s does not exist in the original struct.", name))
	}
}

// GetValue returns the reflect.Value object of the original struct field named name.
// It panics if the original struct does not have a field named name.
func (g *GetterImpl) GetValue(name string) reflect.Value {
	g.panicIfNotHave(name)

	_, ok := g.values[name]
	if !ok {
		g.cache(name)
	}

	return g.values[name]
}

// Get returns the interface of the original struct field named name.
// It panics if the original struct does not have a field named name.
func (g *GetterImpl) Get(name string) interface{} {
	g.panicIfNotHave(name)

	_, ok := g.intfs[name]
	if !ok {
		g.cache(name)
	}

	return g.intfs[name]
}

// EGet returns the interface of the original struct field named name.
// It returns an error if the original struct does not have a field named name.
func (g *GetterImpl) EGet(name string) (intf interface{}, err error) {
	defer func() {
		err = reflectil.RecoverToError(recover())
	}()

	intf = g.Get(name)

	return
}

// Byte returns the byte of the original struct field named name.
// It panics if the original struct does not have a field named name.
// It panics if type of the original struct field named name is not byte.
func (g *GetterImpl) Byte(name string) byte {
	if v, ok := g.Get(name).(byte); ok {
		return v
	}
	panic(fmt.Sprintf("field name %s is not byte type. value kind: %v", name, g.GetValue(name).Kind()))
}

// Bytes returns the []byte of the original struct field named name.
// It panics if the original struct does not have a field named name.
// It panics if type of the original struct field named name is not []byte.
func (g *GetterImpl) Bytes(name string) []byte {
	if v, ok := g.Get(name).([]byte); ok {
		return v
	}
	panic(fmt.Sprintf("field name %s is not []byte type. value kind: %v", name, g.GetValue(name).Kind()))
}

// String returns the string of the original struct field named name.
// It panics if the original struct does not have a field named name.
// It panics if type of the original struct field named name is not string.
func (g *GetterImpl) String(name string) string {
	if v, ok := g.Get(name).(string); ok {
		return v
	}
	panic(fmt.Sprintf("field name %s is not string type. value kind: %v", name, g.GetValue(name).Kind()))
}

// Int returns the int of the original struct field named name.
// It panics if the original struct does not have a field named name.
// It panics if type of the original struct field named name is not int.
func (g *GetterImpl) Int(name string) int {
	if v, ok := g.Get(name).(int); ok {
		return v
	}
	panic(fmt.Sprintf("field name %s is not int type. value kind: %v", name, g.GetValue(name).Kind()))
}

// Int8 returns the int8 of the original struct field named name.
// It panics if the original struct does not have a field named name.
// It panics if type of the original struct field named name is not int8.
func (g *GetterImpl) Int8(name string) int8 {
	if v, ok := g.Get(name).(int8); ok {
		return v
	}
	panic(fmt.Sprintf("field name %s is not int8 type. value kind: %v", name, g.GetValue(name).Kind()))
}

// Int16 returns the int16 of the original struct field named name.
// It panics if the original struct does not have a field named name.
// It panics if type of the original struct field named name is not int16.
func (g *GetterImpl) Int16(name string) int16 {
	if v, ok := g.Get(name).(int16); ok {
		return v
	}
	panic(fmt.Sprintf("field name %s is not int16 type. value kind: %v", name, g.GetValue(name).Kind()))
}

// Int32 returns the int32 of the original struct field named name.
// It panics if the original struct does not have a field named name.
// It panics if type of the original struct field named name is not int32.
func (g *GetterImpl) Int32(name string) int32 {
	if v, ok := g.Get(name).(int32); ok {
		return v
	}
	panic(fmt.Sprintf("field name %s is not int32 type. value kind: %v", name, g.GetValue(name).Kind()))
}

// Int64 returns the int64 of the original struct field named name.
// It panics if the original struct does not have a field named name.
// It panics if type of the original struct field named name is not int64.
func (g *GetterImpl) Int64(name string) int64 {
	if v, ok := g.Get(name).(int64); ok {
		return v
	}
	panic(fmt.Sprintf("field name %s is not int64 type. value kind: %v", name, g.GetValue(name).Kind()))
}

// Uint returns the uint of the original struct field named name.
// It panics if the original struct does not have a field named name.
// It panics if type of the original struct field named name is not uint.
func (g *GetterImpl) Uint(name string) uint {
	if v, ok := g.Get(name).(uint); ok {
		return v
	}
	panic(fmt.Sprintf("field name %s is not uint type. value kind: %v", name, g.GetValue(name).Kind()))
}

// Uint8 returns the uint8 of the original struct field named name.
// It panics if the original struct does not have a field named name.
// It panics if type of the original struct field named name is not uint8.
func (g *GetterImpl) Uint8(name string) uint8 {
	if v, ok := g.Get(name).(uint8); ok {
		return v
	}
	panic(fmt.Sprintf("field name %s is not uint8 type. value kind: %v", name, g.GetValue(name).Kind()))
}

// Uint16 returns the uint16 of the original struct field named name.
// It panics if the original struct does not have a field named name.
// It panics if type of the original struct field named name is not uint16.
func (g *GetterImpl) Uint16(name string) uint16 {
	if v, ok := g.Get(name).(uint16); ok {
		return v
	}
	panic(fmt.Sprintf("field name %s is not uint16 type. value kind: %v", name, g.GetValue(name).Kind()))
}

// Uint32 returns the uint32 of the original struct field named name.
// It panics if the original struct does not have a field named name.
// It panics if type of the original struct field named name is not uint32.
func (g *GetterImpl) Uint32(name string) uint32 {
	if v, ok := g.Get(name).(uint32); ok {
		return v
	}
	panic(fmt.Sprintf("field name %s is not uint32 type. value kind: %v", name, g.GetValue(name).Kind()))
}

// Uint64 returns the uint64 of the original struct field named name.
// It panics if the original struct does not have a field named name.
// It panics if type of the original struct field named name is not uint64.
func (g *GetterImpl) Uint64(name string) uint64 {
	if v, ok := g.Get(name).(uint64); ok {
		return v
	}
	panic(fmt.Sprintf("field name %s is not uint64 type. value kind: %v", name, g.GetValue(name).Kind()))
}

// Uintptr returns the uintptr of the original struct field named name.
// It panics if the original struct does not have a field named name.
// It panics if type of the original struct field named name is not uintptr.
func (g *GetterImpl) Uintptr(name string) uintptr {
	if v, ok := g.Get(name).(uintptr); ok {
		return v
	}
	panic(fmt.Sprintf("field name %s is not uintptr type. value kind: %v", name, g.GetValue(name).Kind()))
}

// Float64 returns the float64 of the original struct field named name.
// It panics if the original struct does not have a field named name.
// It panics if type of the original struct field named name is not float64.
func (g *GetterImpl) Float64(name string) float64 {
	if v, ok := g.Get(name).(float64); ok {
		return v
	}
	panic(fmt.Sprintf("field name %s is not float64 type. value kind: %v", name, g.GetValue(name).Kind()))
}

// Bool returns the bool of the original struct field named name.
// It panics if the original struct does not have a field named name.
// It panics if type of the original struct field named name is not bool.
func (g *GetterImpl) Bool(name string) bool {
	if v, ok := g.Get(name).(bool); ok {
		return v
	}
	panic(fmt.Sprintf("field name %s is not bool type. value kind: %v", name, g.GetValue(name).Kind()))
}

// Complex64 returns the complex64 of the original struct field named name.
// It panics if the original struct does not have a field named name.
// It panics if type of the original struct field named name is not complex64.
func (g *GetterImpl) Complex64(name string) complex64 {
	if v, ok := g.Get(name).(complex64); ok {
		return v
	}
	panic(fmt.Sprintf("field name %s is not complex64 type. value kind: %v", name, g.GetValue(name).Kind()))
}

// Complex128 returns the complex128 of the original struct field named name.
// It panics if the original struct does not have a field named name.
// It panics if type of the original struct field named name is not complex128.
func (g *GetterImpl) Complex128(name string) complex128 {
	if v, ok := g.Get(name).(complex128); ok {
		return v
	}
	panic(fmt.Sprintf("field name %s is not complex128 type. value kind: %v", name, g.GetValue(name).Kind()))
}

// UnsafePointer returns the unsafe.Pointer of the original struct field named name.
// It panics if the original struct does not have a field named name.
// It panics if type of the original struct field named name is not unsafe.Pointer.
func (g *GetterImpl) UnsafePointer(name string) unsafe.Pointer {
	if v, ok := g.Get(name).(unsafe.Pointer); ok {
		return v
	}
	panic(fmt.Sprintf("field name %s is not complex128 type. value kind: %v", name, g.GetValue(name).Kind()))
}

// IsByte reports whether type of the original struct field named name is byte.
func (g *GetterImpl) IsByte(name string) bool {
	return g.is(name, reflect.Uint8)
}

// IsBytes reports whether type of the original struct field named name is []byte.
func (g *GetterImpl) IsBytes(name string) bool {
	return g.IsSlice(name) && g.GetType(name).Elem().Kind() == reflect.Uint8
}

// IsString reports whether type of the original struct field named name is string.
func (g *GetterImpl) IsString(name string) bool {
	return g.is(name, reflect.String)
}

// IsInt reports whether type of the original struct field named name is int.
func (g *GetterImpl) IsInt(name string) bool {
	return g.is(name, reflect.Int)
}

// IsInt8 reports whether type of the original struct field named name is int8.
func (g *GetterImpl) IsInt8(name string) bool {
	return g.is(name, reflect.Int8)
}

// IsInt16 reports whether type of the original struct field named name is int16.
func (g *GetterImpl) IsInt16(name string) bool {
	return g.is(name, reflect.Int16)
}

// IsInt32 reports whether type of the original struct field named name is int32.
func (g *GetterImpl) IsInt32(name string) bool {
	return g.is(name, reflect.Int32)
}

// IsInt64 reports whether type of the original struct field named name is int64.
func (g *GetterImpl) IsInt64(name string) bool {
	return g.is(name, reflect.Int64)
}

// IsUint reports whether type of the original struct field named name is uint.
func (g *GetterImpl) IsUint(name string) bool {
	return g.is(name, reflect.Uint)
}

// IsUint8 reports whether type of the original struct field named name is uint8.
func (g *GetterImpl) IsUint8(name string) bool {
	return g.is(name, reflect.Uint8)
}

// IsUint16 reports whether type of the original struct field named name is uint16.
func (g *GetterImpl) IsUint16(name string) bool {
	return g.is(name, reflect.Uint16)
}

// IsUint32 reports whether type of the original struct field named name is uint32.
func (g *GetterImpl) IsUint32(name string) bool {
	return g.is(name, reflect.Uint32)
}

// IsUint64 reports whether type of the original struct field named name is uint64.
func (g *GetterImpl) IsUint64(name string) bool {
	return g.is(name, reflect.Uint64)
}

// IsUintptr reports whether type of the original struct field named name is uintptr.
func (g *GetterImpl) IsUintptr(name string) bool {
	return g.is(name, reflect.Uintptr)
}

// IsFloat64 reports whether type of the original struct field named name is float64.
func (g *GetterImpl) IsFloat64(name string) bool {
	return g.is(name, reflect.Float64)
}

// IsBool reports whether type of the original struct field named name is bool.
func (g *GetterImpl) IsBool(name string) bool {
	return g.is(name, reflect.Bool)
}

// IsComplex64 reports whether type of the original struct field named name is []byte.
func (g *GetterImpl) IsComplex64(name string) bool {
	return g.is(name, reflect.Complex64)
}

// IsComplex128 reports whether type of the original struct field named name is []byte.
func (g *GetterImpl) IsComplex128(name string) bool {
	return g.is(name, reflect.Complex128)
}

// IsUnsafePointer reports whether type of the original struct field named name is []byte.
func (g *GetterImpl) IsUnsafePointer(name string) bool {
	return g.is(name, reflect.UnsafePointer)
}

// IsMap reports whether type of the original struct field named name is map.
func (g *GetterImpl) IsMap(name string) bool {
	return g.is(name, reflect.Map)
}

// IsFunc reports whether type of the original struct field named name is func.
func (g *GetterImpl) IsFunc(name string) bool {
	return g.is(name, reflect.Func)
}

// IsChan reports whether type of the original struct field named name is chan.
func (g *GetterImpl) IsChan(name string) bool {
	return g.is(name, reflect.Chan)
}

// IsStruct reports whether type of the original struct field named name is struct.
func (g *GetterImpl) IsStruct(name string) bool {
	return g.is(name, reflect.Struct)
}

// IsSlice reports whether type of the original struct field named name is slice.
func (g *GetterImpl) IsSlice(name string) bool {
	return g.is(name, reflect.Slice)
}

// IsArray reports whether type of the original struct field named name is slice.
func (g *GetterImpl) IsArray(name string) bool {
	return g.is(name, reflect.Array)
}

func (g *GetterImpl) is(name string, exp reflect.Kind) bool {
	if !g.Has(name) {
		return false
	}

	frv := g.GetValue(name)
	return frv.Kind() == exp
}

// MapGet returns the interface slice of mapped values of the original struct field named name.
func (g *GetterImpl) MapGet(name string, f func(int, Getter) (interface{}, error)) ([]interface{}, error) {
	if !g.IsSlice(name) {
		return nil, fmt.Errorf("field %s is not slice", name)
	}

	var vi reflect.Value
	var eg Getter
	var err error
	var r interface{}

	srv := g.GetValue(name)
	res := make([]interface{}, srv.Len())

	for i := 0; i < srv.Len(); i++ {
		vi = srv.Index(i)
		eg, err = NewGetter(reflectil.ToI(vi))
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
