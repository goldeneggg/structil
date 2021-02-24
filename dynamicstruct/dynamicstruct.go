package dynamicstruct

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// DynamicStruct is the struct that built dynamic struct by Builder.Build().
type DynamicStruct struct {
	name  string
	rt    reflect.Type
	isPtr bool
	// sortedFields  string  // TODO: for performance tuning
	def string
}

// TODO: add "sortedFields" slice string argument
func newDynamicStruct(fields []reflect.StructField, isPtr bool) *DynamicStruct {
	return newDynamicStructWithName(fields, isPtr, defaultStructName)
}

// newDynamicStructWithName returns a concrete DynamicStruct
// Note: Create DynamicStruct via Builder.Build(), instead of calling this method directly.
func newDynamicStructWithName(fields []reflect.StructField, isPtr bool, name string) *DynamicStruct {
	return &DynamicStruct{
		name:  name,
		rt:    reflect.StructOf(fields),
		isPtr: isPtr,
	}
}

// Name returns the name of this.
func (ds *DynamicStruct) Name() string {
	return ds.name
}

// Type returns the reflect.Type for struct type of this.
func (ds *DynamicStruct) Type() reflect.Type {
	return ds.rt
}

// NumField returns the number of built struct fields.
func (ds *DynamicStruct) NumField() int {
	return ds.rt.NumField()
}

// Field returns the i'th field of the built struct.
func (ds *DynamicStruct) Field(i int) reflect.StructField {
	return ds.rt.Field(i)
}

// FieldByName returns the struct field with the given name
// and a boolean indicating if the field was found.
func (ds *DynamicStruct) FieldByName(name string) (reflect.StructField, bool) {
	return ds.rt.FieldByName(name)
}

// IsPtr reports whether the built struct type is pointer.
func (ds *DynamicStruct) IsPtr() bool {
	return ds.isPtr
}

// NewInterface returns the new interface value of built struct.
func (ds *DynamicStruct) NewInterface() interface{} {
	rv := reflect.New(ds.rt)
	if ds.isPtr {
		return rv.Interface()
	}

	return reflect.Indirect(rv).Interface()
}

// DecodeMap returns the interface that was decoded from input map.
func (ds *DynamicStruct) DecodeMap(m map[string]interface{}) (interface{}, error) {
	if !ds.IsPtr() {
		return nil, errors.New("DecodeMap can execute only if dynamic struct is pointer. But this is false")
	}

	i := ds.NewInterface()
	err := mapstructure.Decode(m, &i)
	return i, err
}

// Definition returns the struct definition string with field indention by TAB.
// Fields are sorted by field name.
func (ds *DynamicStruct) Definition() string {
	// build definition only once
	if ds.def != "" {
		return ds.def
	}

	// TODO: performance optimization
	// TODO: nested DynamicStruct fields are also sorted as same as high-level DynamicStruct
	fields := make([]reflect.StructField, ds.NumField())
	for i := 0; i < ds.NumField(); i++ {
		fields[i] = ds.Field(i)
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})

	var stb strings.Builder
	indent := "\t"

	stb.WriteString("type " + ds.Name() + " struct {\n")
	for _, field := range fields {
		stb.WriteString(indent)
		stb.WriteString(field.Name)
		stb.WriteString(" ")
		stb.WriteString(field.Type.String())
		if field.Tag != "" {
			stb.WriteString(" ")
			stb.WriteString(fmt.Sprintf("`%s`", field.Tag))
		}
		stb.WriteString("\n")
	}
	stb.WriteString("}")

	ds.def = stb.String()
	return ds.def
}
