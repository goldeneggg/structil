package dynamicstruct

import (
	"errors"
	"reflect"
	"sync"

	"github.com/goldeneggg/structil/util"
)

const (
	capBuilderField   = 10
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
)

var (
	// ErrSample is sample init error value
	ErrSample = errors.New("SampleError")
)

// Builder is the interface that builds a dynamic and runtime struct.
type Builder struct {
	name  string
	bfMap builderFieldMap
	mu    sync.RWMutex
	err   error
}

// NewBuilder returns a concrete Builder
func NewBuilder() *Builder {
	return &Builder{
		name:  defaultStructName,
		bfMap: make(builderFieldMap),
	}
}

type builderField struct {
	name string
	typ  reflect.Type
	tag  reflect.StructTag
}

type builderFieldMap map[string]*builderField

func (b *Builder) getFieldMap(key string) *builderField {
	b.mu.RLock()
	defer b.mu.RUnlock()

	r := b.bfMap[key]
	return r
}

func (b *Builder) putFieldMap(key string, bf *builderField) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.bfMap[key] = bf
}

func (b *Builder) deleteFieldMap(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.bfMap, key)
}

func (b *Builder) hasFieldMap(key string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	_, ok := b.bfMap[key]
	return ok
}

func (b *Builder) lenFieldMap() int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return len(b.bfMap)
}

func (b *Builder) addFieldFunc(name string, isPtr bool, tag string, f func() reflect.Type) *Builder {
	defer func() {
		err := util.RecoverToError(recover())

		// keep 1st recoverd error
		if err != nil && b.err == nil {
			b.err = err
		}
	}()

	typ := f()

	if isPtr {
		typ = reflect.PtrTo(typ)
	}

	b.putFieldMap(name, &builderField{
		name: name,
		typ:  typ,
		tag:  reflect.StructTag(tag),
	})

	return b
}

// AddString returns a Builder that was added a string field named by name parameter.
func (b *Builder) AddString(name string) *Builder {
	b.AddStringWithTag(name, "")
	return b
}

// AddStringWithTag returns a Builder that was added a string field with tag named by name parameter.
func (b *Builder) AddStringWithTag(name string, tag string) *Builder {
	f := func() reflect.Type {
		return reflect.TypeOf(SampleString)
	}
	b.addFieldFunc(name, false, tag, f)

	return b
}

// AddInt returns a Builder that was added a int field named by name parameter.
func (b *Builder) AddInt(name string) *Builder {
	b.AddIntWithTag(name, "")
	return b
}

// AddIntWithTag returns a Builder that was added a int field with tag named by name parameter.
func (b *Builder) AddIntWithTag(name string, tag string) *Builder {
	f := func() reflect.Type {
		return reflect.TypeOf(SampleInt)
	}
	b.addFieldFunc(name, false, tag, f)

	return b
}

// AddByte returns a Builder that was added a byte field named by name parameter.
func (b *Builder) AddByte(name string) *Builder {
	b.AddByteWithTag(name, "")
	return b
}

// AddByteWithTag returns a Builder that was added a byte field with tag named by name parameter.
func (b *Builder) AddByteWithTag(name string, tag string) *Builder {
	f := func() reflect.Type {
		return reflect.TypeOf(SampleByte)
	}
	b.addFieldFunc(name, false, tag, f)

	return b
}

// AddFloat32 returns a Builder that was added a float32 field named by name parameter.
func (b *Builder) AddFloat32(name string) *Builder {
	b.AddFloat32WithTag(name, "")
	return b
}

// AddFloat32WithTag returns a Builder that was added a float32 field with tag named by name parameter.
func (b *Builder) AddFloat32WithTag(name string, tag string) *Builder {
	f := func() reflect.Type {
		return reflect.TypeOf(SampleFloat32)
	}
	b.addFieldFunc(name, false, tag, f)

	return b
}

// AddFloat64 returns a Builder that was added a float64 field named by name parameter.
func (b *Builder) AddFloat64(name string) *Builder {
	b.AddFloat64WithTag(name, "")
	return b
}

// AddFloat64WithTag returns a Builder that was added a float64 field with tag named by name parameter.
func (b *Builder) AddFloat64WithTag(name string, tag string) *Builder {
	f := func() reflect.Type {
		return reflect.TypeOf(SampleFloat64)
	}
	b.addFieldFunc(name, false, tag, f)

	return b
}

// AddBool returns a Builder that was added a bool field named by name parameter.
func (b *Builder) AddBool(name string) *Builder {
	b.AddBoolWithTag(name, "")
	return b
}

// AddBoolWithTag returns a Builder that was added a bool field with tag named by name parameter.
func (b *Builder) AddBoolWithTag(name string, tag string) *Builder {
	f := func() reflect.Type {
		return reflect.TypeOf(SampleBool)
	}
	b.addFieldFunc(name, false, tag, f)

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
	f := func() reflect.Type {
		var vt reflect.Type
		if vi == nil {
			vt = reflect.TypeOf((*interface{})(nil)).Elem()
		} else {
			vt = reflect.TypeOf(vi)
		}
		return reflect.MapOf(reflect.TypeOf(ki), vt)
	}
	b.addFieldFunc(name, false, tag, f)

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
	f := func() reflect.Type {
		it := make([]reflect.Type, len(in))
		for i := 0; i < len(in); i++ {
			it[i] = reflect.TypeOf(in[i])
		}
		ot := make([]reflect.Type, len(out))
		for i := 0; i < len(out); i++ {
			ot[i] = reflect.TypeOf(out[i])
		}
		return reflect.FuncOf(it, ot, false)
	}
	b.addFieldFunc(name, false, tag, f)

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
	f := func() reflect.Type {
		return reflect.ChanOf(reflect.BothDir, reflect.TypeOf(i))
	}
	b.addFieldFunc(name, false, tag, f)

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
	f := func() reflect.Type {
		return reflect.ChanOf(reflect.RecvDir, reflect.TypeOf(i))
	}
	b.addFieldFunc(name, false, tag, f)

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
	f := func() reflect.Type {
		return reflect.ChanOf(reflect.SendDir, reflect.TypeOf(i))
	}
	b.addFieldFunc(name, false, tag, f)

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
	f := func() reflect.Type {
		iType := reflect.TypeOf(i)
		if iType.Kind() == reflect.Ptr {
			iType = iType.Elem()
		}
		fields := make([]reflect.StructField, iType.NumField())
		for i := 0; i < iType.NumField(); i++ {
			fields[i] = iType.Field(i)
		}
		return reflect.StructOf(fields)
	}
	b.addFieldFunc(name, isPtr, tag, f)

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
	f := func() reflect.Type {
		return reflect.SliceOf(reflect.TypeOf(i))
	}
	b.addFieldFunc(name, false, tag, f)

	return b
}

// AddSlicePtr returns a Builder that was added a slice pointer field named by name parameter.
// Type of slice is elem type of i.
func (b *Builder) AddSlicePtr(name string, i interface{}) *Builder {
	b.AddSliceWithTag(name, i, "")
	return b
}

// AddSlicePtrWithTag returns a Builder that was added a slice pointer field with tag named by name parameter.
// Type of slice is elem type of i.
func (b *Builder) AddSlicePtrWithTag(name string, i interface{}, tag string) *Builder {
	f := func() reflect.Type {
		return reflect.SliceOf(reflect.TypeOf(i))
	}
	b.addFieldFunc(name, true, tag, f)

	return b
}

// AddInterface returns a Builder that was added a interface{} field named by name parameter.
func (b *Builder) AddInterface(name string, isPtr bool) *Builder {
	b.AddInterfaceWithTag(name, isPtr, "")
	return b
}

// AddInterfaceWithTag returns a Builder that was added a interface{} field with tag named by name parameter.
func (b *Builder) AddInterfaceWithTag(name string, isPtr bool, tag string) *Builder {
	f := func() reflect.Type {
		return reflect.TypeOf((*interface{})(nil)).Elem()
	}
	b.addFieldFunc(name, isPtr, tag, f)

	return b
}

// AddDynamicStruct returns a Builder that was added a DynamicStruct field named by name parameter.
func (b *Builder) AddDynamicStruct(name string, ds *DynamicStruct, isPtr bool) *Builder {
	b.AddDynamicStructWithTag(name, ds, isPtr, "")
	return b
}

// AddDynamicStructWithTag returns a Builder that was added a DynamicStruct field with tag named by name parameter.
func (b *Builder) AddDynamicStructWithTag(name string, ds *DynamicStruct, isPtr bool, tag string) *Builder {
	b.AddStructWithTag(name, ds.NewInterface(), isPtr, tag)

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

// AddDynamicStructSlice returns a Builder that was added a DynamicStruct slice field named by name parameter.
func (b *Builder) AddDynamicStructSlice(name string, ds *DynamicStruct) *Builder {
	return b.AddDynamicStructSliceWithTag(name, ds, "")
}

// AddDynamicStructSliceWithTag returns a Builder that was added a DynamicStruct slice field with tag named by name parameter.
func (b *Builder) AddDynamicStructSliceWithTag(name string, ds *DynamicStruct, tag string) *Builder {
	b.AddSliceWithTag(name, ds.NewInterface(), tag)

	return b
}

// NumField returns the number of built struct fields.
func (b *Builder) NumField() int {
	return b.lenFieldMap()
}

// Exists returns true if the specified name field exists
func (b *Builder) Exists(name string) bool {
	return b.hasFieldMap(name)
}

// GetStructName returns the name of this DynamicStruct.
func (b *Builder) GetStructName() string {
	return b.name
}

// SetStructName returns a Builder that was set the name of DynamicStruct.
// Default name is "DynamicStruct"
func (b *Builder) SetStructName(name string) *Builder {
	b.name = name
	return b
}

// GetTag returns the tag of the specified name field.
func (b *Builder) GetTag(name string) string {
	if b.hasFieldMap(name) {
		return string(b.getFieldMap(name).tag)
	}
	return ""
}

// SetTag returns a Builder that was set the tag for the specific field.
// Expected tag string is 'TYPE1:"FIELDNAME1" TYPEn:"FIELDNAMEn"' format (e.g. json:"id" etc)
func (b *Builder) SetTag(name string, tag string) *Builder {
	if b.hasFieldMap(name) {
		b.getFieldMap(name).tag = reflect.StructTag(tag)
	}
	return b
}

// Remove returns a Builder that was removed a field named by name parameter.
func (b *Builder) Remove(name string) *Builder {
	b.deleteFieldMap(name)
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

	var i int
	fields := make([]reflect.StructField, b.lenFieldMap())

	for key, bf := range b.bfMap {
		fields[i] = reflect.StructField{
			Name: key,
			Type: bf.typ,
			Tag:  bf.tag,
		}
		i++
	}

	ds, err = newDynamicStruct(fields, isPtr, b.GetStructName())

	return
}
