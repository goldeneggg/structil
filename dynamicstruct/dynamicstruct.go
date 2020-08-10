package dynamicstruct

import (
	"errors"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// DynamicStruct is the interface that built dynamic struct by Builder.Build().
type DynamicStruct interface {
	NumField() int
	Field(i int) reflect.StructField
	FieldByName(name string) (reflect.StructField, bool)
	IsPtr() bool
	Interface() interface{}
	DecodeMap(m map[string]interface{}) (interface{}, error)
}

// impl is the default DynamicStruct implementation.
type impl struct {
	structType reflect.Type
	isPtr      bool
	intf       interface{}
}

// newDynamicStruct returns a concrete DynamicStruct
// Note: Create DynamicStruct via Builder.Build(), instead of calling this method directly.
//
// func New(fs []reflect.StructField, isPtr bool) DynamicStruct {
func newDynamicStruct(fs []reflect.StructField, isPtr bool) DynamicStruct {
	ds := &impl{
		structType: reflect.StructOf(fs),
		isPtr:      isPtr,
	}

	n := reflect.New(ds.structType)
	if isPtr {
		ds.intf = n.Interface()
	} else {
		ds.intf = reflect.Indirect(n).Interface()
	}

	return ds
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

// Interface returns the interface of built struct.
func (ds *impl) Interface() interface{} {
	return ds.intf
}

// DecodeMap returns the interface that was decoded from input map.
func (ds *impl) DecodeMap(m map[string]interface{}) (interface{}, error) {
	if !ds.IsPtr() {
		return nil, errors.New("DecodeMap can execute only if dynamic struct is pointer. But this is false")
	}

	i := ds.Interface()
	err := mapstructure.Decode(m, &i)
	return i, err
}
