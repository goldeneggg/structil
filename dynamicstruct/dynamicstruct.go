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
	name   string
	fields []reflect.StructField
	rt     reflect.Type
	isPtr  bool
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
		name:   name,
		fields: fields,
		rt:     reflect.StructOf(fields),
		isPtr:  isPtr,
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

// Fields returns the all fields of the built struct.
func (ds *DynamicStruct) Fields() []reflect.StructField {
	return ds.fields
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

	var stb strings.Builder
	ds.def = definition(&stb, ds.fields, ds.name, 1)
	return ds.def
}

func definition(stbp *strings.Builder, flds []reflect.StructField, name string, indentLevel int) string {
	sortedFlds := sortFields(flds)

	if indentLevel == 1 {
		stbp.WriteString("type ")
	}
	if name != "" {
		stbp.WriteString(name + " ")
	}
	stbp.WriteString("struct {\n")

	indent := strings.Repeat("\t", indentLevel)
	//indent := "\t"
	for _, sf := range sortedFlds {
		stbp.WriteString(indent)
		stbp.WriteString(sf.Name)
		stbp.WriteString(" ")

		nt := sf.Type
		if nt.Kind() == reflect.Ptr {
			nt = nt.Elem()
		}
		if nt.Kind() == reflect.Struct {
			// recursively call if type is struct
			nflds := make([]reflect.StructField, nt.NumField())
			for i := 0; i < nt.NumField(); i++ {
				nflds[i] = nt.Field(i)
			}
			var nstb strings.Builder
			stbp.WriteString(definition(&nstb, nflds, "", indentLevel+1))
		} else {
			stbp.WriteString(sf.Type.String())
			if sf.Tag != "" {
				stbp.WriteString(" ")
				stbp.WriteString(fmt.Sprintf("`%s`", sf.Tag))
			}
		}

		stbp.WriteString("\n")
	}

	if indentLevel > 1 {
		stbp.WriteString(strings.Repeat("\t", indentLevel-1))
	}
	stbp.WriteString("}")

	return stbp.String()
}

func sortFields(fields []reflect.StructField) []reflect.StructField {
	// TODO: performance optimization
	sfs := make([]reflect.StructField, len(fields))
	for i := 0; i < len(fields); i++ {
		sfs[i] = fields[i]
	}
	sort.Slice(sfs, func(i, j int) bool {
		return sfs[i].Name < sfs[j].Name
	})

	return sfs
}
