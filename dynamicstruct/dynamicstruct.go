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
	// TODO:
	// Clone() DynamicStruct
}

// impl is the default DynamicStruct implementation.
type impl struct {
	name       string
	structType reflect.Type
	isPtr      bool
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

	return strbuilder.String()
}

/*
// JSONToDynamicStructInterface returns an interface via DynamicStruct.DecodeMap from JSON data.
// jsonData argument must be a byte array data of JSON.
//
// This method supports known format JSON and unknown format JSON.
// But when JSON format is known, this method is not recommended. Because this method is suitable for unknown JSON with heavy and slow reflection functions.
//
// Field names in DynamicStruct are converted to CamelCase automatically
// - e.g. "hoge" JSON field is converted to "Hoge".
// - e.g. "huga_field" JSON field is converted to "HugaField".
//
// Deprecated: This function is very experimental and it will be removed soon
func JSONToDynamicStructInterface(jsonData []byte) (interface{}, error) {
	// FIXME:
	// want to add json validation. but is json.Valid(data) too slow?
	// See: https://stackoverflow.com/questions/22128282/how-to-check-string-is-in-json-format

	var unmarshalledJSON interface{}
	err := json.Unmarshal(jsonData, &unmarshalledJSON)
	if err != nil {
		return nil, err
	}

	return parseUnmarshalledJSON(unmarshalledJSON)
}

func parseUnmarshalledJSON(unmarshalledJSON interface{}) (interface{}, error) {
	switch t := unmarshalledJSON.(type) {
	case map[string]interface{}:
		return mapToDynamicStructInterface(t)
	case []interface{}:
		var i interface{}
		var err error
		iArr := make([]interface{}, len(t))
		for idx, elemJSON := range t {
			// call this function recursively
			i, err = parseUnmarshalledJSON(elemJSON)
			if err != nil {
				return nil, err
			}

			iArr[idx] = i
		}

		return iArr, nil
	}

	return nil, fmt.Errorf("unexpected return. unmarshalledJSON %+v is not map or array", unmarshalledJSON)
}

func mapToDynamicStructInterface(m map[string]interface{}) (interface{}, error) {
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
			// Note: Is this ok?
			b = b.AddInterfaceWithTag(camelizedKey, false, tag)
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
*/
