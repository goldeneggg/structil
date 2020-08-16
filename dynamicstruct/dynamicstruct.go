package dynamicstruct

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// DynamicStruct is the interface that built dynamic struct by Builder.Build().
type DynamicStruct interface {
	Name() string
	NumField() int
	Field(i int) reflect.StructField
	FieldByName(name string) (reflect.StructField, bool)
	IsPtr() bool
	NewInterface() interface{}
	DecodeMap(m map[string]interface{}) (interface{}, error)
	Definition() string
}

// impl is the default DynamicStruct implementation.
type impl struct {
	name       string
	structType reflect.Type
	isPtr      bool
	definition string
}

func newDynamicStruct(fields []reflect.StructField, isPtr bool) DynamicStruct {
	return newDynamicStructWithName(fields, isPtr, defaultStructName)
}

// newDynamicStructWithName returns a concrete DynamicStruct
// Note: Create DynamicStruct via Builder.Build(), instead of calling this method directly.
func newDynamicStructWithName(fields []reflect.StructField, isPtr bool, name string) DynamicStruct {
	ds := &impl{
		name:       name,
		structType: reflect.StructOf(fields),
		isPtr:      isPtr,
	}

	return ds
}

// Name returns the name of this.
func (ds *impl) Name() string {
	return ds.name
}

// NumField returns the number of built struct fields.
func (ds *impl) NumField() int {
	return ds.structType.NumField()
}

// Field returns the i'th field of the built struct.
func (ds *impl) Field(i int) reflect.StructField {
	return ds.structType.Field(i)
}

// FieldByName returns the struct field with the given name
// and a boolean indicating if the field was found.
func (ds *impl) FieldByName(name string) (reflect.StructField, bool) {
	return ds.structType.FieldByName(name)
}

// IsPtr reports whether the built struct type is pointer.
func (ds *impl) IsPtr() bool {
	return ds.isPtr
}

// NewInterface returns the new interface value of built struct.
func (ds *impl) NewInterface() interface{} {
	rv := reflect.New(ds.structType)
	if ds.isPtr {
		return rv.Interface()
	}

	return reflect.Indirect(rv).Interface()
}

// DecodeMap returns the interface that was decoded from input map.
func (ds *impl) DecodeMap(m map[string]interface{}) (interface{}, error) {
	if !ds.IsPtr() {
		return nil, errors.New("DecodeMap can execute only if dynamic struct is pointer. But this is false")
	}

	i := ds.NewInterface()
	err := mapstructure.Decode(m, &i)
	return i, err
}

// Definition returns the struct definition string with field indention by TAB.
// Fields are sorted by field name.
func (ds *impl) Definition() string {
	// TODO: build definition only once
	if ds.definition != "" {
		return ds.definition
	}

	sortedFields := make([]reflect.StructField, ds.NumField())
	for i := 0; i < ds.NumField(); i++ {
		sortedFields[i] = ds.Field(i)
	}
	sort.Slice(sortedFields, func(i, j int) bool {
		return sortedFields[i].Name < sortedFields[j].Name
	})

	var strbuilder strings.Builder
	indent := "\t"

	strbuilder.WriteString("type " + ds.Name() + " struct {\n")
	for _, field := range sortedFields {
		strbuilder.WriteString(indent)
		strbuilder.WriteString(field.Name)
		strbuilder.WriteString(" ")
		strbuilder.WriteString(field.Type.String())
		if field.Tag != "" {
			strbuilder.WriteString(" ")
			strbuilder.WriteString(fmt.Sprintf("`%s`", field.Tag))
		}
		strbuilder.WriteString("\n")
	}
	strbuilder.WriteString("}")

	ds.definition = strbuilder.String()
	return ds.definition
}
