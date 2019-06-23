package dynamicstruct

import (
	"errors"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type ofType int

const (
	// SampleString is sample init string value
	SampleString = ""
	// SampleInt is sample init int value
	SampleInt = 0
	// SampleFloat is sample init float value
	SampleFloat = 0.0
	// SampleBool is sample init bool value
	SampleBool = false

	tMap ofType = iota
	tFunc
	tChanBoth
	tChanRecv
	tChanSend
	tStruct
	tSlice
	tPrmtv
)

var (
	// ErrSample is sample init error value
	ErrSample = errors.New("SampleError")
)

// Builder is the interface that builds a dynamic and runtime struct.
type Builder interface {
	AddString(name string) Builder
	AddInt(name string) Builder
	AddFloat(name string) Builder
	AddBool(name string) Builder
	AddMap(name string, ke interface{}, ve interface{}) Builder
	AddFunc(name string, eargs []interface{}, erets []interface{}) Builder
	AddChanBoth(name string, e interface{}) Builder
	AddChanRecv(name string, e interface{}) Builder
	AddChanSend(name string, e interface{}) Builder
	AddStruct(name string, i interface{}, isPtr bool) Builder
	AddStructPtr(name string, i interface{}) Builder
	AddSlice(name string, e interface{}) Builder
	Remove(name string) Builder
	Exists(name string) bool
	NumField() int
	Build() DynamicStruct
	BuildNonPtr() DynamicStruct
}

// BuilderImpl is the default Builder implementation.
type BuilderImpl struct {
	fields map[string]reflect.Type
}

// NewBuilder returns a concrete Builder
func NewBuilder() Builder {
	return &BuilderImpl{fields: map[string]reflect.Type{}}
}

type addParam struct {
	name     string
	intfs    []interface{}
	keyIntfs []interface{}
	ot       ofType
	isPtr    bool
}

// AddString returns a Builder that was added a string field named by name parameter.
func (b *BuilderImpl) AddString(name string) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{SampleString},
		ot:    tPrmtv,
		isPtr: false,
	}
	b.add(p)
	return b
}

// AddInt returns a Builder that was added a int field named by name parameter.
func (b *BuilderImpl) AddInt(name string) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{SampleInt},
		ot:    tPrmtv,
		isPtr: false,
	}
	b.add(p)
	return b
}

// AddFloat returns a Builder that was added a float64 field named by name parameter.
func (b *BuilderImpl) AddFloat(name string) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{SampleFloat},
		ot:    tPrmtv,
		isPtr: false,
	}
	b.add(p)
	return b
}

// AddBool returns a Builder that was added a bool field named by name parameter.
func (b *BuilderImpl) AddBool(name string) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{SampleBool},
		ot:    tPrmtv,
		isPtr: false,
	}
	b.add(p)
	return b
}

// AddMap returns a Builder that was added a map field named by name parameter.
// Type of map key is type of ke.
// Type of map value is type of ve.
func (b *BuilderImpl) AddMap(name string, ke interface{}, ve interface{}) Builder {
	p := &addParam{
		name:     name,
		intfs:    []interface{}{ve},
		keyIntfs: []interface{}{ke},
		ot:       tMap,
		isPtr:    false,
	}
	b.add(p)
	return b
}

// AddFunc returns a Builder that was added a func field named by name parameter.
// Types of func args are types of eargs.
// Types of func returns are types of erets.
func (b *BuilderImpl) AddFunc(name string, eargs []interface{}, erets []interface{}) Builder {
	p := &addParam{
		name:     name,
		intfs:    erets,
		keyIntfs: eargs,
		ot:       tFunc,
		isPtr:    false,
	}
	b.add(p)
	return b
}

// AddChanBoth returns a Builder that was added a BothDir chan field named by name parameter.
// Type of chan is type of e.
func (b *BuilderImpl) AddChanBoth(name string, e interface{}) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{e},
		ot:    tChanBoth,
		isPtr: false,
	}
	b.add(p)
	return b
}

// AddChanRecv returns a Builder that was added a RecvDir chan field named by name parameter.
// Type of chan is type of e.
func (b *BuilderImpl) AddChanRecv(name string, e interface{}) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{e},
		ot:    tChanRecv,
		isPtr: false,
	}
	b.add(p)
	return b
}

// AddChanSend returns a Builder that was added a SendDir chan field named by name parameter.
// Type of chan is type of e.
func (b *BuilderImpl) AddChanSend(name string, e interface{}) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{e},
		ot:    tChanSend,
		isPtr: false,
	}
	b.add(p)
	return b
}

// AddStruct returns a Builder that was added a struct field named by name parameter.
// Type of struct is type of i.
func (b *BuilderImpl) AddStruct(name string, i interface{}, isPtr bool) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{i},
		ot:    tStruct,
		isPtr: isPtr,
	}
	b.add(p)
	return b
}

// AddStructPtr returns a Builder that was added a struct pointer field named by name parameter.
// Type of struct is type of i.
func (b *BuilderImpl) AddStructPtr(name string, i interface{}) Builder {
	return b.AddStruct(name, i, true)
}

// AddSlice returns a Builder that was added a slice field named by name parameter.
// Type of slice is type of e.
func (b *BuilderImpl) AddSlice(name string, e interface{}) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{e},
		ot:    tSlice,
		isPtr: false,
	}
	b.add(p)
	return b
}

func (b *BuilderImpl) add(p *addParam) {
	var typeOf reflect.Type

	switch p.ot {
	case tMap:
		typeOf = reflect.MapOf(reflect.TypeOf(p.keyIntfs[0]), reflect.TypeOf(p.intfs[0]))
	case tFunc:
		aTypes := make([]reflect.Type, len(p.keyIntfs))
		for i := 0; i < len(p.keyIntfs); i++ {
			aTypes[i] = reflect.TypeOf(p.keyIntfs[i])
		}

		vTypes := make([]reflect.Type, len(p.intfs))
		for i := 0; i < len(p.intfs); i++ {
			vTypes[i] = reflect.TypeOf(p.intfs[i])
		}
		// TODO: variadic support
		typeOf = reflect.FuncOf(aTypes, vTypes, false)
	case tChanBoth:
		typeOf = reflect.ChanOf(reflect.BothDir, reflect.TypeOf(p.intfs[0]))
	case tChanRecv:
		typeOf = reflect.ChanOf(reflect.RecvDir, reflect.TypeOf(p.intfs[0]))
	case tChanSend:
		typeOf = reflect.ChanOf(reflect.SendDir, reflect.TypeOf(p.intfs[0]))
	case tStruct:
		it := reflect.TypeOf(p.intfs[0])
		if it.Kind() == reflect.Ptr {
			it = it.Elem()
		}
		fs := make([]reflect.StructField, it.NumField())
		for i := 0; i < it.NumField(); i++ {
			fs[i] = it.Field(i)
		}
		typeOf = reflect.StructOf(fs)
	case tSlice:
		typeOf = reflect.SliceOf(reflect.TypeOf(p.intfs[0]))
	default:
		typeOf = reflect.TypeOf(p.intfs[0])
	}

	if p.isPtr {
		typeOf = reflect.PtrTo(typeOf)
	}

	b.fields[p.name] = typeOf
}

// Remove returns a Builder that was removed a field named by name parameter.
func (b *BuilderImpl) Remove(name string) Builder {
	delete(b.fields, name)
	return b
}

// Exists returns true if the specified name field exists
func (b *BuilderImpl) Exists(name string) bool {
	_, ok := b.fields[name]
	return ok
}

// NumField returns the number of built struct fields.
func (b *BuilderImpl) NumField() int {
	return len(b.fields)
}

// Build returns a concrete struct pointer built by Builder.
func (b *BuilderImpl) Build() DynamicStruct {
	return b.build(true)
}

// BuildNonPtr returns a concrete struct built by Builder.
func (b *BuilderImpl) BuildNonPtr() DynamicStruct {
	return b.build(false)
}

func (b *BuilderImpl) build(isPtr bool) DynamicStruct {
	var i int
	fs := make([]reflect.StructField, len(b.fields))
	for name, typ := range b.fields {
		fs[i] = reflect.StructField{Name: name, Type: typ}
		i++
	}

	return newDs(fs, isPtr)
}

// DynamicStruct is the interface that built dynamic struct by Builder.Build().
type DynamicStruct interface {
	NumField() int
	Field(i int) reflect.StructField
	FieldByName(name string) (reflect.StructField, bool)
	IsPtr() bool
	Interface() interface{}
	DecodeMap(m map[string]interface{}) (interface{}, error)
}

// Impl is the default DynamicStruct implementation.
type Impl struct {
	structType reflect.Type
	isPtr      bool
	intf       interface{}
}

func newDs(fs []reflect.StructField, isPtr bool) DynamicStruct {
	ds := &Impl{structType: reflect.StructOf(fs), isPtr: isPtr}

	n := reflect.New(ds.structType)
	if isPtr {
		ds.intf = n.Interface()
	} else {
		ds.intf = reflect.Indirect(n).Interface()
	}

	return ds
}

// NumField returns the number of built struct fields.
func (ds *Impl) NumField() int {
	return ds.structType.NumField()
}

// Field returns the i'th field of the built struct.
func (ds *Impl) Field(i int) reflect.StructField {
	return ds.structType.Field(i)
}

// FieldByName returns the struct field with the given name
// and a boolean indicating if the field was found.
func (ds *Impl) FieldByName(name string) (reflect.StructField, bool) {
	return ds.structType.FieldByName(name)
}

// IsPtr reports whether the built struct type is pointer.
func (ds *Impl) IsPtr() bool {
	return ds.isPtr
}

// Interface returns the interface of built struct.
func (ds *Impl) Interface() interface{} {
	return ds.intf
}

// DecodeMap returns the interface that was decoded from input map.
func (ds *Impl) DecodeMap(m map[string]interface{}) (interface{}, error) {
	if !ds.IsPtr() {
		return nil, errors.New("DecodeMap can execute only if dynamic struct is pointer. But this is false")
	}

	err := mapstructure.Decode(m, &ds.intf)
	return ds.intf, err
}
