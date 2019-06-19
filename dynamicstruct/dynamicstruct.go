package dynamicstruct

import (
	"reflect"
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
	smpMap  map[interface{}]interface{}
	smpFunc func([]interface{}) []interface{}
)

// DynamicStruct is the interface that builds a dynamic and runtime struct.
type DynamicStruct interface {
	AddString(name string) DynamicStruct
	AddInt(name string) DynamicStruct
	AddFloat(name string) DynamicStruct
	AddBool(name string) DynamicStruct
	AddMap(name string, ke interface{}, ve interface{}) DynamicStruct
	AddFunc(name string, eargs []interface{}, erets []interface{}) DynamicStruct
	AddChanBoth(name string, e interface{}) DynamicStruct
	AddChanRecv(name string, e interface{}) DynamicStruct
	AddChanSend(name string, e interface{}) DynamicStruct
	AddStruct(name string, i interface{}, isPtr bool) DynamicStruct
	AddStructPtr(name string, i interface{}) DynamicStruct
	AddSlice(name string, e interface{}) DynamicStruct
	Build() interface{}
	BuildNonPtr() interface{}
	NumBuiltField() int
	BuiltField(i int) reflect.StructField
	Remove(name string) DynamicStruct
	Exists(name string) bool
}

// DynamicStructImpl is the default DynamicStruct implementation.
type DynamicStructImpl struct {
	fields     map[string]reflect.Type
	structType reflect.Type
}

// New returns a concrete DynamicStruct
func New() DynamicStruct {
	return &DynamicStructImpl{fields: map[string]reflect.Type{}}
}

type addParam struct {
	name     string
	intfs    []interface{}
	keyIntfs []interface{}
	ot       ofType
	isPtr    bool
}

// AddString returns a DynamicStruct that was added a string field named by "name" parameter.
func (ds *DynamicStructImpl) AddString(name string) DynamicStruct {
	p := &addParam{
		name:  name,
		intfs: []interface{}{SampleString},
		ot:    tPrmtv,
		isPtr: false,
	}
	ds.add(p)
	return ds
}

// AddInt returns a DynamicStruct that was added a int field named by "name" parameter.
func (ds *DynamicStructImpl) AddInt(name string) DynamicStruct {
	p := &addParam{
		name:  name,
		intfs: []interface{}{SampleInt},
		ot:    tPrmtv,
		isPtr: false,
	}
	ds.add(p)
	return ds
}

// AddFloat returns a DynamicStruct that was added a float64 field named by "name" parameter.
func (ds *DynamicStructImpl) AddFloat(name string) DynamicStruct {
	p := &addParam{
		name:  name,
		intfs: []interface{}{SampleFloat},
		ot:    tPrmtv,
		isPtr: false,
	}
	ds.add(p)
	return ds
}

// AddBool returns a DynamicStruct that was added a bool field named by "name" parameter.
func (ds *DynamicStructImpl) AddBool(name string) DynamicStruct {
	p := &addParam{
		name:  name,
		intfs: []interface{}{SampleBool},
		ot:    tPrmtv,
		isPtr: false,
	}
	ds.add(p)
	return ds
}

// AddMap returns a DynamicStruct that was added a map field named by "name" parameter.
// Type of map key is type of "ke" and type of map value is type of "ve".
func (ds *DynamicStructImpl) AddMap(name string, ke interface{}, ve interface{}) DynamicStruct {
	p := &addParam{
		name:     name,
		intfs:    []interface{}{ve},
		keyIntfs: []interface{}{ke},
		ot:       tMap,
		isPtr:    false,
	}
	ds.add(p)
	return ds
}

// AddFunc returns a DynamicStruct that was added a func field named by "name" parameter.
// Types of func args are types of "eargs" and types of func returns are types of "erets".
func (ds *DynamicStructImpl) AddFunc(name string, eargs []interface{}, erets []interface{}) DynamicStruct {
	p := &addParam{
		name:     name,
		intfs:    erets,
		keyIntfs: eargs,
		ot:       tFunc,
		isPtr:    false,
	}
	ds.add(p)
	return ds
}

// AddChanBoth returns a DynamicStruct that was added a BothDir chan field named by "name" parameter.
// Type of chan is type of "e".
func (ds *DynamicStructImpl) AddChanBoth(name string, e interface{}) DynamicStruct {
	p := &addParam{
		name:  name,
		intfs: []interface{}{e},
		ot:    tChanBoth,
		isPtr: false,
	}
	ds.add(p)
	return ds
}

// AddChanRecv returns a DynamicStruct that was added a RecvDir chan field named by "name" parameter.
// Type of chan is type of "e".
func (ds *DynamicStructImpl) AddChanRecv(name string, e interface{}) DynamicStruct {
	p := &addParam{
		name:  name,
		intfs: []interface{}{e},
		ot:    tChanRecv,
		isPtr: false,
	}
	ds.add(p)
	return ds
}

// AddChanSend returns a DynamicStruct that was added a SendDir chan field named by "name" parameter.
// Type of chan is type of "e".
func (ds *DynamicStructImpl) AddChanSend(name string, e interface{}) DynamicStruct {
	p := &addParam{
		name:  name,
		intfs: []interface{}{e},
		ot:    tChanSend,
		isPtr: false,
	}
	ds.add(p)
	return ds
}

// AddStruct returns a DynamicStruct that was added a struct field named by "name" parameter.
// Type of struct is type of "i".
func (ds *DynamicStructImpl) AddStruct(name string, i interface{}, isPtr bool) DynamicStruct {
	p := &addParam{
		name:  name,
		intfs: []interface{}{i},
		ot:    tStruct,
		isPtr: isPtr,
	}
	ds.add(p)
	return ds
}

// AddStruct returns a DynamicStruct that was added a struct pointer field named by "name" parameter.
// Type of struct is type of "i".
func (ds *DynamicStructImpl) AddStructPtr(name string, i interface{}) DynamicStruct {
	return ds.AddStruct(name, i, true)
}

// AddSlice returns a DynamicStruct that was added a slice field named by "name" parameter.
// Type of slice is type of "e".
func (ds *DynamicStructImpl) AddSlice(name string, e interface{}) DynamicStruct {
	p := &addParam{
		name:  name,
		intfs: []interface{}{e},
		ot:    tSlice,
		isPtr: false,
	}
	ds.add(p)
	return ds
}

func (ds *DynamicStructImpl) add(p *addParam) {
	it := reflect.TypeOf(p.intfs[0])
	var typeOf reflect.Type

	switch p.ot {
	case tMap:
		kt := reflect.TypeOf(p.keyIntfs[0])
		typeOf = reflect.MapOf(kt, it)
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
		typeOf = reflect.ChanOf(reflect.BothDir, it)
	case tChanRecv:
		typeOf = reflect.ChanOf(reflect.RecvDir, it)
	case tChanSend:
		typeOf = reflect.ChanOf(reflect.SendDir, it)
	case tStruct:
		if it.Kind() == reflect.Ptr {
			it = it.Elem()
		}
		fs := make([]reflect.StructField, it.NumField())
		for i := 0; i < it.NumField(); i++ {
			fs[i] = it.Field(i)
		}
		typeOf = reflect.StructOf(fs)
	case tSlice:
		typeOf = reflect.SliceOf(it)
	default:
		typeOf = it
	}

	if p.isPtr {
		typeOf = reflect.PtrTo(typeOf)
	}

	ds.fields[p.name] = typeOf
}

// Build returns a concrete struct pointer built by DynamicStruct.
func (ds *DynamicStructImpl) Build() interface{} {
	return ds.build(true)
}

// BuildNonPtr returns a concrete struct built by DynamicStruct.
func (ds *DynamicStructImpl) BuildNonPtr() interface{} {
	return ds.build(false)
}

func (ds *DynamicStructImpl) build(isPtr bool) interface{} {
	var i int
	fs := make([]reflect.StructField, len(ds.fields))

	for name, typ := range ds.fields {
		fs[i] = reflect.StructField{Name: name, Type: typ}
		i++
	}
	ds.structType = reflect.StructOf(fs)
	n := reflect.New(ds.structType)

	if isPtr {
		return n.Interface()
	} else {
		return reflect.Indirect(n).Interface()
	}
}

// NumBuiltField returns the number of built struct fields.
func (ds *DynamicStructImpl) NumBuiltField() int {
	return ds.structType.NumField()
}

// BuiltField returns the i'th field of the built struct.
func (ds *DynamicStructImpl) BuiltField(i int) reflect.StructField {
	return ds.structType.Field(i)
}

// Remove returns a DynamicStruct that was removed a field named by "name" parameter.
func (ds *DynamicStructImpl) Remove(name string) DynamicStruct {
	delete(ds.fields, name)
	return ds
}

// Exists returns true if the specified name field exists
func (ds *DynamicStructImpl) Exists(name string) bool {
	_, ok := ds.fields[name]
	return ok
}
