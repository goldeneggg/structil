package dynamicstruct

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/iancoleman/strcase"
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

// JSONToDynamicStructInterface returns an interface via DynamicStruct.DecodeMap from JSON data.
// jsonData argument must be a byte array data of JSON.
//
// This method supports known format JSON and unknown format JSON.
// But when JSON format is known, this method is not recommended. Because this method is suitable for unknown JSON with heavy and slow reflection functions.
//
// Field names in DynamicStruct are converted to CamelCase automatically
// - e.g. "hoge" JSON field is converted to "Hoge".
// - e.g. "huga_field" JSON field is converted to "HugaField".
func JSONToDynamicStructInterface(jsonData []byte) (interface{}, error) {
	var unmarshalled interface{}
	err := json.Unmarshal(jsonData, &unmarshalled)
	if err != nil {
		return nil, err
	}

	// FIXME: unmarshalled が配列の場合がある
	m, ok := unmarshalled.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can not cast type from unmarshalled unknownJSON[%#v] to map", unmarshalled)
	}

	var camelizedKey, tag string
	camelizedFieldMap := make(map[string]interface{}, len(m))
	b := NewBuilder()

	for k, v := range m {
		camelizedKey = strcase.ToCamel(k)
		camelizedFieldMap[camelizedKey] = v
		tag = fmt.Sprintf(`json:"%s"`, k)

		switch value := v.(type) {
		case bool:
			b = b.AddBoolWithTag(camelizedKey, tag)
		case float64:
			b = b.AddFloat64WithTag(camelizedKey, tag)
		case string:
			b = b.AddStringWithTag(camelizedKey, tag)
		case []interface{}:
			b = b.AddSliceWithTag(camelizedKey, interface{}(value[0]), tag)
		case map[string]interface{}:
			for kk, vv := range value {
				b = b.AddMapWithTag(camelizedKey, kk, interface{}(vv), tag)
				break
			}
		case nil:
			// FIXME:
			// fmt.Printf("@@@ (nil) k: %s, v: %#v, value: %#v\n", k, v, value)
		default:
			return nil, fmt.Errorf("jsonData %#v has invalid typed key", m)
		}
	}

	ds := b.Build()
	intf, err := ds.DecodeMap(camelizedFieldMap)
	if err != nil {
		return nil, err
	}

	return intf, nil
}
