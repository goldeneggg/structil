package dynamicstruct

import (
	"errors"
	"reflect"

	"github.com/goldeneggg/structil/util"
)

type pattern int

const (
	defaultStructName = "DynamicStruct"

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

	patternMap pattern = iota
	patternFunc
	patternChanBoth
	patternChanRecv
	patternChanSend
	patternStruct
	patternSlice
	patternPrmtv
	patternInterface
	patternDynamicStruct
)

var (
	// ErrSample is sample init error value
	ErrSample = errors.New("SampleError")
)

// Builder is the interface that builds a dynamic and runtime struct.
type Builder struct {
	name   string
	fields map[string]reflect.Type
	tags   map[string]reflect.StructTag // optional
	err    error
}

// NewBuilder returns a concrete Builder
func NewBuilder() *Builder {
	return &Builder{
		name:   defaultStructName,
		fields: map[string]reflect.Type{},
		tags:   map[string]reflect.StructTag{},
	}
}

type addParam struct {
	name     string
	intfs    []interface{}
	keyIntfs []interface{}
	pattern  pattern
	isPtr    bool
	tag      string
}

// AddString returns a Builder that was added a string field named by name parameter.
func (b *Builder) AddString(name string) *Builder {
	b.AddStringWithTag(name, "")
	return b
}

// AddStringWithTag returns a Builder that was added a string field with tag named by name parameter.
func (b *Builder) AddStringWithTag(name string, tag string) *Builder {
	p := &addParam{
		name:    name,
		intfs:   []interface{}{SampleString},
		pattern: patternPrmtv,
		isPtr:   false,
		tag:     tag,
	}
	b.add(p)
	return b
}

// AddInt returns a Builder that was added a int field named by name parameter.
func (b *Builder) AddInt(name string) *Builder {
	b.AddIntWithTag(name, "")
	return b
}

// AddIntWithTag returns a Builder that was added a int field with tag named by name parameter.
func (b *Builder) AddIntWithTag(name string, tag string) *Builder {
	p := &addParam{
		name:    name,
		intfs:   []interface{}{SampleInt},
		pattern: patternPrmtv,
		isPtr:   false,
		tag:     tag,
	}
	b.add(p)
	return b
}

// AddByte returns a Builder that was added a byte field named by name parameter.
func (b *Builder) AddByte(name string) *Builder {
	b.AddByteWithTag(name, "")
	return b
}

// AddByteWithTag returns a Builder that was added a byte field with tag named by name parameter.
func (b *Builder) AddByteWithTag(name string, tag string) *Builder {
	p := &addParam{
		name:    name,
		intfs:   []interface{}{SampleByte},
		pattern: patternPrmtv,
		isPtr:   false,
		tag:     tag,
	}
	b.add(p)
	return b
}

// AddFloat32 returns a Builder that was added a float32 field named by name parameter.
func (b *Builder) AddFloat32(name string) *Builder {
	b.AddFloat32WithTag(name, "")
	return b
}

// AddFloat32WithTag returns a Builder that was added a float32 field with tag named by name parameter.
func (b *Builder) AddFloat32WithTag(name string, tag string) *Builder {
	p := &addParam{
		name:    name,
		intfs:   []interface{}{SampleFloat32},
		pattern: patternPrmtv,
		isPtr:   false,
		tag:     tag,
	}
	b.add(p)
	return b
}

// AddFloat64 returns a Builder that was added a float64 field named by name parameter.
func (b *Builder) AddFloat64(name string) *Builder {
	b.AddFloat64WithTag(name, "")
	return b
}

// AddFloat64WithTag returns a Builder that was added a float64 field with tag named by name parameter.
func (b *Builder) AddFloat64WithTag(name string, tag string) *Builder {
	p := &addParam{
		name:    name,
		intfs:   []interface{}{SampleFloat64},
		pattern: patternPrmtv,
		isPtr:   false,
		tag:     tag,
	}
	b.add(p)
	return b
}

// AddBool returns a Builder that was added a bool field named by name parameter.
func (b *Builder) AddBool(name string) *Builder {
	b.AddBoolWithTag(name, "")
	return b
}

// AddBoolWithTag returns a Builder that was added a bool field with tag named by name parameter.
func (b *Builder) AddBoolWithTag(name string, tag string) *Builder {
	p := &addParam{
		name:    name,
		intfs:   []interface{}{SampleBool},
		pattern: patternPrmtv,
		isPtr:   false,
		tag:     tag,
	}
	b.add(p)
	return b
}

// AddMap returns a Builder that was added a map field named by name parameter.
// Type of map key is type of ki.
// Type of map value is type of vi.
func (b *Builder) AddMap(name string, ki interface{}, vi interface{}) *Builder {
	b.AddMapWithTag(name, ki, vi, "")
	return b
}

// AddMapWithTag returns a Builder that was added a map field with tag named by name parameter.
// Type of map key is type of ki.
// Type of map value is type of vi.
func (b *Builder) AddMapWithTag(name string, ki interface{}, vi interface{}, tag string) *Builder {
	p := &addParam{
		name:     name,
		intfs:    []interface{}{vi},
		keyIntfs: []interface{}{ki},
		pattern:  patternMap,
		isPtr:    false,
		tag:      tag,
	}
	b.add(p)
	return b
}

// AddFunc returns a Builder that was added a func field named by name parameter.
// Types of func args are types of in.
// Types of func returns are types of out.
func (b *Builder) AddFunc(name string, in []interface{}, out []interface{}) *Builder {
	b.AddFuncWithTag(name, in, out, "")
	return b
}

// AddFuncWithTag returns a Builder that was added a func field with tag named by name parameter.
// Types of func args are types of in.
// Types of func returns are types of out.
func (b *Builder) AddFuncWithTag(name string, in []interface{}, out []interface{}, tag string) *Builder {
	p := &addParam{
		name:     name,
		intfs:    out,
		keyIntfs: in,
		pattern:  patternFunc,
		isPtr:    false,
		tag:      tag,
	}
	b.add(p)
	return b
}

// AddChanBoth returns a Builder that was added a BothDir chan field named by name parameter.
// Type of chan is type of i.
func (b *Builder) AddChanBoth(name string, i interface{}) *Builder {
	b.AddChanBothWithTag(name, i, "")
	return b
}

// AddChanBothWithTag returns a Builder that was added a BothDir chan field with tag named by name parameter.
// Type of chan is type of i.
func (b *Builder) AddChanBothWithTag(name string, i interface{}, tag string) *Builder {
	p := &addParam{
		name:    name,
		intfs:   []interface{}{i},
		pattern: patternChanBoth,
		isPtr:   false,
		tag:     tag,
	}
	b.add(p)
	return b
}

// AddChanRecv returns a Builder that was added a RecvDir chan field named by name parameter.
// Type of chan is type of i.
func (b *Builder) AddChanRecv(name string, i interface{}) *Builder {
	b.AddChanRecvWithTag(name, i, "")
	return b
}

// AddChanRecvWithTag returns a Builder that was added a RecvDir chan field with tag named by name parameter.
// Type of chan is type of i.
func (b *Builder) AddChanRecvWithTag(name string, i interface{}, tag string) *Builder {
	p := &addParam{
		name:    name,
		intfs:   []interface{}{i},
		pattern: patternChanRecv,
		isPtr:   false,
		tag:     tag,
	}
	b.add(p)
	return b
}

// AddChanSend returns a Builder that was added a SendDir chan field named by name parameter.
// Type of chan is type of i.
func (b *Builder) AddChanSend(name string, i interface{}) *Builder {
	b.AddChanSendWithTag(name, i, "")
	return b
}

// AddChanSendWithTag returns a Builder that was added a SendDir chan field with tag named by name parameter.
// Type of chan is type of i.
func (b *Builder) AddChanSendWithTag(name string, i interface{}, tag string) *Builder {
	p := &addParam{
		name:    name,
		intfs:   []interface{}{i},
		pattern: patternChanSend,
		isPtr:   false,
		tag:     tag,
	}
	b.add(p)
	return b
}

// AddStruct returns a Builder that was added a struct field named by name parameter.
// Type of struct is type of i.
func (b *Builder) AddStruct(name string, i interface{}, isPtr bool) *Builder {
	b.AddStructWithTag(name, i, isPtr, "")
	return b
}

// AddStructWithTag returns a Builder that was added a struct field with tag named by name parameter.
// Type of struct is type of i.
func (b *Builder) AddStructWithTag(name string, i interface{}, isPtr bool, tag string) *Builder {
	p := &addParam{
		name:    name,
		intfs:   []interface{}{i},
		pattern: patternStruct,
		isPtr:   isPtr,
		tag:     tag,
	}
	b.add(p)
	return b
}

// AddStructPtr returns a Builder that was added a struct pointer field named by name parameter.
// Type of struct is type of i.
func (b *Builder) AddStructPtr(name string, i interface{}) *Builder {
	return b.AddStruct(name, i, true)
}

// AddStructPtrWithTag returns a Builder that was added a struct pointer field with tag named by name parameter.
// Type of struct is type of i.
func (b *Builder) AddStructPtrWithTag(name string, i interface{}, tag string) *Builder {
	return b.AddStructWithTag(name, i, true, tag)
}

// AddSlice returns a Builder that was added a slice field named by name parameter.
// Type of slice is type of i.
func (b *Builder) AddSlice(name string, i interface{}) *Builder {
	b.AddSliceWithTag(name, i, "")
	return b
}

// AddSliceWithTag returns a Builder that was added a slice field with tag named by name parameter.
// Type of slice is type of i.
func (b *Builder) AddSliceWithTag(name string, i interface{}, tag string) *Builder {
	p := &addParam{
		name:    name,
		intfs:   []interface{}{i},
		pattern: patternSlice,
		isPtr:   false,
		tag:     tag,
	}
	b.add(p)
	return b
}

// AddInterface returns a Builder that was added a interface{} field named by name parameter.
func (b *Builder) AddInterface(name string, isPtr bool) *Builder {
	b.AddInterfaceWithTag(name, isPtr, "")
	return b
}

// AddInterfaceWithTag returns a Builder that was added a interface{} field with tag named by name parameter.
func (b *Builder) AddInterfaceWithTag(name string, isPtr bool, tag string) *Builder {
	p := &addParam{
		name:    name,
		intfs:   []interface{}{(*interface{})(nil)},
		pattern: patternInterface,
		isPtr:   isPtr,
		tag:     tag,
	}
	b.add(p)
	return b
}

// AddDynamicStruct returns a Builder that was added a DynamicStruct field named by name parameter.
func (b *Builder) AddDynamicStruct(name string, ds *DynamicStruct, isPtr bool) *Builder {
	b.AddDynamicStructWithTag(name, ds, isPtr, "")
	return b
}

// AddDynamicStructWithTag returns a Builder that was added a DynamicStruct field with tag named by name parameter.
func (b *Builder) AddDynamicStructWithTag(name string, ds *DynamicStruct, isPtr bool, tag string) *Builder {
	p := &addParam{
		name:    name,
		intfs:   []interface{}{ds.NewInterface()}, // use ds.NewInterface() for building concrete fields of ds
		pattern: patternStruct,                    // use patternStruct for building concrete fields of ds
		isPtr:   isPtr,
		tag:     tag,
	}
	b.add(p)
	return b
}

// AddDynamicStructPtr returns a Builder that was added a DynamicStruct pointer field named by name parameter.
func (b *Builder) AddDynamicStructPtr(name string, ds *DynamicStruct) *Builder {
	return b.AddDynamicStruct(name, ds, true)
}

// AddDynamicStructPtrWithTag returns a Builder that was added a DynamicStruct pointer field with tag named by name parameter.
func (b *Builder) AddDynamicStructPtrWithTag(name string, ds *DynamicStruct, tag string) *Builder {
	return b.AddDynamicStructWithTag(name, ds, true, tag)
}

func (b *Builder) add(p *addParam) {
	defer func() {
		err := util.RecoverToError(recover())

		// keep 1st recoverd error
		if err != nil && b.err == nil {
			b.err = err
		}
	}()

	var typeOf reflect.Type

	switch p.pattern {
	case patternMap:
		var vt reflect.Type
		if p.intfs[0] == nil {
			vt = reflect.TypeOf((*interface{})(nil)).Elem()
		} else {
			vt = reflect.TypeOf(p.intfs[0])
		}
		typeOf = reflect.MapOf(reflect.TypeOf(p.keyIntfs[0]), vt)
	case patternFunc:
		inTypes := make([]reflect.Type, len(p.keyIntfs))
		for i := 0; i < len(p.keyIntfs); i++ {
			inTypes[i] = reflect.TypeOf(p.keyIntfs[i])
		}

		outTypes := make([]reflect.Type, len(p.intfs))
		for i := 0; i < len(p.intfs); i++ {
			outTypes[i] = reflect.TypeOf(p.intfs[i])
		}
		// TODO: variadic support
		typeOf = reflect.FuncOf(inTypes, outTypes, false)
	case patternChanBoth:
		typeOf = reflect.ChanOf(reflect.BothDir, reflect.TypeOf(p.intfs[0]))
	case patternChanRecv:
		typeOf = reflect.ChanOf(reflect.RecvDir, reflect.TypeOf(p.intfs[0]))
	case patternChanSend:
		typeOf = reflect.ChanOf(reflect.SendDir, reflect.TypeOf(p.intfs[0]))
	case patternStruct:
		iType := reflect.TypeOf(p.intfs[0])
		if iType.Kind() == reflect.Ptr {
			iType = iType.Elem()
		}

		fields := make([]reflect.StructField, iType.NumField())
		for i := 0; i < iType.NumField(); i++ {
			fields[i] = iType.Field(i)
		}
		typeOf = reflect.StructOf(fields)
	case patternSlice:
		typeOf = reflect.SliceOf(reflect.TypeOf(p.intfs[0]))
	case patternInterface:
		typeOf = reflect.TypeOf(p.intfs[0]).Elem()
	default:
		typeOf = reflect.TypeOf(p.intfs[0])
	}

	if p.isPtr {
		typeOf = reflect.PtrTo(typeOf)
	}

	b.fields[p.name] = typeOf
	b.SetTag(p.name, p.tag)
}

// SetStructName returns a Builder that was set the name of DynamicStruct.
// Default name is "DynamicStruct"
func (b *Builder) SetStructName(name string) *Builder {
	b.name = name
	return b
}

// GetStructName returns the name of this DynamicStruct.
func (b *Builder) GetStructName() string {
	return b.name
}

// SetTag returns a Builder that was set the tag for the specific field.
// Expected tag string is 'TYPE1:"FIELDNAME1" TYPEn:"FIELDNAMEn"' format (e.g. json:"id" etc)
func (b *Builder) SetTag(name string, tag string) *Builder {
	b.tags[name] = reflect.StructTag(tag)
	return b
}

// NumField returns the number of built struct fields.
func (b *Builder) NumField() int {
	return len(b.fields)
}

// Exists returns true if the specified name field exists
func (b *Builder) Exists(name string) bool {
	_, ok := b.fields[name]
	return ok
}

// Remove returns a Builder that was removed a field named by name parameter.
func (b *Builder) Remove(name string) *Builder {
	delete(b.fields, name)
	return b
}

// Build returns a concrete struct pointer built by Builder.
func (b *Builder) Build() (*DynamicStruct, error) {
	return b.build(true)
}

// BuildNonPtr returns a concrete struct built by Builder.
func (b *Builder) BuildNonPtr() (*DynamicStruct, error) {
	return b.build(false)
}

func (b *Builder) build(isPtr bool) (ds *DynamicStruct, err error) {
	if b.err != nil {
		err = b.err
		return
	}

	defer func() {
		err = util.RecoverToError(recover())
		return
	}()

	var i int
	fields := make([]reflect.StructField, len(b.fields))
	for name, typ := range b.fields {
		fields[i] = reflect.StructField{
			Name: name,
			Type: typ,
			Tag:  b.tags[name],
		}
		i++
	}

	ds = newDynamicStructWithName(fields, isPtr, b.GetStructName())

	return
}
