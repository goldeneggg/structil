package dynamicstruct

import (
	"errors"
	"reflect"
)

type ofType int

const (
	// SampleString is sample string value
	SampleString = ""
	// SampleInt is sample int value
	SampleInt = int(1)
	// SampleByte is sample byte value
	SampleByte = byte(1)
	// SampleFloat32 is sample float32 value
	SampleFloat32 = float32(1.1)
	// SampleFloat64 is sample float64 value
	SampleFloat64 = float64(1.1)
	// SampleBool is sample bool value
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
	AddStringWithTag(name string, tag string) Builder
	AddInt(name string) Builder
	AddIntWithTag(name string, tag string) Builder
	AddByte(name string) Builder
	AddByteWithTag(name string, tag string) Builder
	AddFloat32(name string) Builder
	AddFloat32WithTag(name string, tag string) Builder
	AddFloat64(name string) Builder
	AddFloat64WithTag(name string, tag string) Builder
	AddBool(name string) Builder
	AddBoolWithTag(name string, tag string) Builder
	AddMap(name string, ke interface{}, ve interface{}) Builder
	AddMapWithTag(name string, ke interface{}, ve interface{}, tag string) Builder
	AddFunc(name string, eargs []interface{}, erets []interface{}) Builder
	AddFuncWithTag(name string, eargs []interface{}, erets []interface{}, tag string) Builder
	AddChanBoth(name string, e interface{}) Builder
	AddChanBothWithTag(name string, e interface{}, tag string) Builder
	AddChanRecv(name string, e interface{}) Builder
	AddChanRecvWithTag(name string, e interface{}, tag string) Builder
	AddChanSend(name string, e interface{}) Builder
	AddChanSendWithTag(name string, e interface{}, tag string) Builder
	AddStruct(name string, i interface{}, isPtr bool) Builder
	AddStructWithTag(name string, i interface{}, isPtr bool, tag string) Builder
	AddStructPtr(name string, i interface{}) Builder
	AddStructPtrWithTag(name string, i interface{}, tag string) Builder
	AddSlice(name string, e interface{}) Builder
	AddSliceWithTag(name string, e interface{}, tag string) Builder
	Remove(name string) Builder
	Exists(name string) bool
	NumField() int
	Build() DynamicStruct
	BuildNonPtr() DynamicStruct
}

// BuilderImpl is the default Builder implementation.
type BuilderImpl struct {
	fields map[string]reflect.Type
	tags   map[string]reflect.StructTag
}

// NewBuilder returns a concrete Builder
func NewBuilder() Builder {
	return &BuilderImpl{fields: map[string]reflect.Type{}, tags: map[string]reflect.StructTag{}}
}

type addParam struct {
	name  string
	intfs []interface{}
	//intfTypesIsInterface bool
	keyIntfs []interface{}
	//keyTypesIsInterface  bool
	ot    ofType
	isPtr bool
	tag   string
}

// AddString returns a Builder that was added a string field named by name parameter.
func (b *BuilderImpl) AddString(name string) Builder {
	b.AddStringWithTag(name, "")
	return b
}

// AddStringWithTag returns a Builder that was added a string field with tag named by name parameter.
func (b *BuilderImpl) AddStringWithTag(name string, tag string) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{SampleString},
		ot:    tPrmtv,
		isPtr: false,
		tag:   tag,
	}
	b.add(p)
	return b
}

// AddInt returns a Builder that was added a int field named by name parameter.
func (b *BuilderImpl) AddInt(name string) Builder {
	b.AddIntWithTag(name, "")
	return b
}

// AddIntWithTag returns a Builder that was added a int field with tag named by name parameter.
func (b *BuilderImpl) AddIntWithTag(name string, tag string) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{SampleInt},
		ot:    tPrmtv,
		isPtr: false,
		tag:   tag,
	}
	b.add(p)
	return b
}

// AddByte returns a Builder that was added a byte field named by name parameter.
func (b *BuilderImpl) AddByte(name string) Builder {
	b.AddByteWithTag(name, "")
	return b
}

// AddByteWithTag returns a Builder that was added a byte field with tag named by name parameter.
func (b *BuilderImpl) AddByteWithTag(name string, tag string) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{SampleByte},
		ot:    tPrmtv,
		isPtr: false,
		tag:   tag,
	}
	b.add(p)
	return b
}

// AddFloat32 returns a Builder that was added a float32 field named by name parameter.
func (b *BuilderImpl) AddFloat32(name string) Builder {
	b.AddFloat32WithTag(name, "")
	return b
}

// AddFloat32WithTag returns a Builder that was added a float32 field with tag named by name parameter.
func (b *BuilderImpl) AddFloat32WithTag(name string, tag string) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{SampleFloat32},
		ot:    tPrmtv,
		isPtr: false,
		tag:   tag,
	}
	b.add(p)
	return b
}

// AddFloat64 returns a Builder that was added a float64 field named by name parameter.
func (b *BuilderImpl) AddFloat64(name string) Builder {
	b.AddFloat64WithTag(name, "")
	return b
}

// AddFloat64WithTag returns a Builder that was added a float64 field with tag named by name parameter.
func (b *BuilderImpl) AddFloat64WithTag(name string, tag string) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{SampleFloat64},
		ot:    tPrmtv,
		isPtr: false,
		tag:   tag,
	}
	b.add(p)
	return b
}

// AddBool returns a Builder that was added a bool field named by name parameter.
func (b *BuilderImpl) AddBool(name string) Builder {
	b.AddBoolWithTag(name, "")
	return b
}

// AddBoolWithTag returns a Builder that was added a bool field with tag named by name parameter.
func (b *BuilderImpl) AddBoolWithTag(name string, tag string) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{SampleBool},
		ot:    tPrmtv,
		isPtr: false,
		tag:   tag,
	}
	b.add(p)
	return b
}

// AddMap returns a Builder that was added a map field named by name parameter.
// Type of map key is type of ke.
// Type of map value is type of ve.
func (b *BuilderImpl) AddMap(name string, ke interface{}, ve interface{}) Builder {
	b.AddMapWithTag(name, ke, ve, "")
	return b
}

// AddMapWithTag returns a Builder that was added a map field with tag named by name parameter.
// Type of map key is type of ke.
// Type of map value is type of ve.
func (b *BuilderImpl) AddMapWithTag(name string, ke interface{}, ve interface{}, tag string) Builder {
	p := &addParam{
		name:     name,
		intfs:    []interface{}{ve},
		keyIntfs: []interface{}{ke},
		ot:       tMap,
		isPtr:    false,
		tag:      tag,
	}
	b.add(p)
	return b
}

// AddFunc returns a Builder that was added a func field named by name parameter.
// Types of func args are types of eargs.
// Types of func returns are types of erets.
func (b *BuilderImpl) AddFunc(name string, eargs []interface{}, erets []interface{}) Builder {
	b.AddFuncWithTag(name, eargs, erets, "")
	return b
}

// AddFuncWithTag returns a Builder that was added a func field with tag named by name parameter.
// Types of func args are types of eargs.
// Types of func returns are types of erets.
func (b *BuilderImpl) AddFuncWithTag(name string, eargs []interface{}, erets []interface{}, tag string) Builder {
	p := &addParam{
		name:     name,
		intfs:    erets,
		keyIntfs: eargs,
		ot:       tFunc,
		isPtr:    false,
		tag:      tag,
	}
	b.add(p)
	return b
}

// AddChanBoth returns a Builder that was added a BothDir chan field named by name parameter.
// Type of chan is type of e.
func (b *BuilderImpl) AddChanBoth(name string, e interface{}) Builder {
	b.AddChanBothWithTag(name, e, "")
	return b
}

// AddChanBothWithTag returns a Builder that was added a BothDir chan field with tag named by name parameter.
// Type of chan is type of e.
func (b *BuilderImpl) AddChanBothWithTag(name string, e interface{}, tag string) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{e},
		ot:    tChanBoth,
		isPtr: false,
		tag:   tag,
	}
	b.add(p)
	return b
}

// AddChanRecv returns a Builder that was added a RecvDir chan field named by name parameter.
// Type of chan is type of e.
func (b *BuilderImpl) AddChanRecv(name string, e interface{}) Builder {
	b.AddChanRecvWithTag(name, e, "")
	return b
}

// AddChanRecvWithTag returns a Builder that was added a RecvDir chan field with tag named by name parameter.
// Type of chan is type of e.
func (b *BuilderImpl) AddChanRecvWithTag(name string, e interface{}, tag string) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{e},
		ot:    tChanRecv,
		isPtr: false,
		tag:   tag,
	}
	b.add(p)
	return b
}

// AddChanSend returns a Builder that was added a SendDir chan field named by name parameter.
// Type of chan is type of e.
func (b *BuilderImpl) AddChanSend(name string, e interface{}) Builder {
	b.AddChanSendWithTag(name, e, "")
	return b
}

// AddChanSendWithTag returns a Builder that was added a SendDir chan field with tag named by name parameter.
// Type of chan is type of e.
func (b *BuilderImpl) AddChanSendWithTag(name string, e interface{}, tag string) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{e},
		ot:    tChanSend,
		isPtr: false,
		tag:   tag,
	}
	b.add(p)
	return b
}

// AddStruct returns a Builder that was added a struct field named by name parameter.
// Type of struct is type of i.
func (b *BuilderImpl) AddStruct(name string, i interface{}, isPtr bool) Builder {
	b.AddStructWithTag(name, i, isPtr, "")
	return b
}

// AddStructWithTag returns a Builder that was added a struct field with tag named by name parameter.
// Type of struct is type of i.
func (b *BuilderImpl) AddStructWithTag(name string, i interface{}, isPtr bool, tag string) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{i},
		ot:    tStruct,
		isPtr: isPtr,
		tag:   tag,
	}
	b.add(p)
	return b
}

// AddStructPtr returns a Builder that was added a struct pointer field named by name parameter.
// Type of struct is type of i.
func (b *BuilderImpl) AddStructPtr(name string, i interface{}) Builder {
	return b.AddStruct(name, i, true)
}

// AddStructPtrWithTag returns a Builder that was added a struct pointer field with tag named by name parameter.
// Type of struct is type of i.
func (b *BuilderImpl) AddStructPtrWithTag(name string, i interface{}, tag string) Builder {
	return b.AddStructWithTag(name, i, true, tag)
}

// AddSlice returns a Builder that was added a slice field named by name parameter.
// Type of slice is type of e.
func (b *BuilderImpl) AddSlice(name string, e interface{}) Builder {
	b.AddSliceWithTag(name, e, "")
	return b
}

// AddSliceWithTag returns a Builder that was added a slice field with tag named by name parameter.
// Type of slice is type of e.
func (b *BuilderImpl) AddSliceWithTag(name string, e interface{}, tag string) Builder {
	p := &addParam{
		name:  name,
		intfs: []interface{}{e},
		ot:    tSlice,
		isPtr: false,
		tag:   tag,
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
	b.tags[p.name] = reflect.StructTag(p.tag)
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
		fs[i] = reflect.StructField{
			Name: name,
			Type: typ,
			Tag:  b.tags[name],
		}
		i++
	}

	return newDynamicStruct(fs, isPtr)
}
