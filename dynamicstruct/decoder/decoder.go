package decoder

import (
	"encoding/json"
	"fmt"

	"github.com/iancoleman/strcase"

	"github.com/goldeneggg/structil"
	"github.com/goldeneggg/structil/dynamicstruct"
)

// Decoder is the struct that decodes some marshaled data like JSON and YAML.
type Decoder struct {
	data []byte
	dt   dataType
	unm  interface{}
	ds   *dynamicstruct.DynamicStruct
}

// FromJSON returns a concrete Decoder for JSON.
func FromJSON(data []byte) (*Decoder, error) {
	return newDecoder(data, typeJSON)
}

// JSONToI returns a decoded interface from JSON via DynamicStruct.
func JSONToI(data []byte) (interface{}, error) {
	d, err := FromJSON(data)
	if err != nil {
		return nil, err
	}

	return d.toDynamicStructI()
}

// JSONToGetter returns a structil.Getter with a decoded JSON via DynamicStruct.
func JSONToGetter(data []byte) (*structil.Getter, error) {
	intf, err := JSONToI(data)
	if err != nil {
		return nil, err
	}

	return structil.NewGetter(intf)
}

// FromYAML returns a concrete Decoder for YAML.
func FromYAML(data []byte) (*Decoder, error) {
	return newDecoder(data, typeYAML)
}

// YAMLToI returns a decoded interface from YAML via DynamicStruct.
// Note: The gopkg.in/yaml.v2 package returns an unmarshaled interface as "map[interface{}]interface{}" type.
func YAMLToI(data []byte) (interface{}, error) {
	d, err := FromYAML(data)
	if err != nil {
		return nil, err
	}

	return d.toDynamicStructI()
}

// YAMLToGetter returns a structil.Getter with a decoded YAML via DynamicStruct.
func YAMLToGetter(data []byte) (*structil.Getter, error) {
	intf, err := YAMLToI(data)
	if err != nil {
		return nil, err
	}

	// Note:
	// YAMLToI returns a map[interface{}]interface{} so via json.Marshal for generating []byte data.
	// FIXME: json.Marshal does not support "map[interface{}]interface{}", and is expected "map[string]interface{}"
	jsonD, err := json.Marshal(intf)
	if err != nil {
		return nil, err
	}

	return JSONToGetter(jsonD)
}

func newDecoder(data []byte, dt dataType) (d *Decoder, err error) {
	var intf interface{}
	err = dt.unmarshal(data, &intf)

	d = &Decoder{
		data: data,
		dt:   dt,
		unm:  intf,
	}

	return
}

func (d *Decoder) toDynamicStructI() (interface{}, error) {
	if d.ds == nil {
		_, err := d.DynamicStruct(true, true)
		if err != nil {
			return nil, err
		}
	}

	intf := d.ds.NewInterface()
	if err := d.dt.unmarshal(d.data, &intf); err != nil {
		return nil, err
	}

	return intf, nil
}

// RawData returns an original data as []byte.
func (d *Decoder) RawData() []byte {
	return d.data
}

// Interface returns a unmarshaled interface from original data.
func (d *Decoder) Interface() interface{} {
	return d.unm
}

// Map returns a map of unmarshaled interface from original data.
func (d *Decoder) Map() (map[string]interface{}, error) {
	return d.dt.intfToStringMap(d.Interface())
}

// DynamicStruct returns a decoded DynamicStruct.
func (d *Decoder) DynamicStruct(nest bool, useTag bool) (*dynamicstruct.DynamicStruct, error) {
	ds, err := d.toDs(d.unm, nest, useTag)
	d.ds = ds
	return ds, err
}

func (d *Decoder) toDs(i interface{}, nest bool, useTag bool) (*dynamicstruct.DynamicStruct, error) {
	switch t := i.(type) {
	case map[string]interface{}:
		return d.toDsFromStringMap(t, nest, useTag)
	case []interface{}:
		if len(t) > 0 {
			if len(t) == 1 {
				return d.toDs(t[0], nest, useTag)
			}

			// TODO: seek an element that have max size of t. And call d.toDs with this element
			// 配列内の構造が可変なケースを考慮して、最も大きい構造の要素を取り出してその要素に対してtoDsを呼ぶようにする
			// See: https://stackoverflow.com/questions/44257522/how-to-get-memory-size-of-variable-in-go
			//   should use "unsafe.Sizeof(var)"?
			return d.toDs(t[0], nest, useTag)
		}
	// YAML support
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for k, v := range t {
			m[fmt.Sprintf("%v", k)] = v
		}
		return d.toDsFromStringMap(m, nest, useTag)
	}

	return nil, fmt.Errorf("unexpected interface: %#v", i)
}

func (d *Decoder) toDsFromStringMap(m map[string]interface{}, nest bool, useTag bool) (*dynamicstruct.DynamicStruct, error) {
	var tag, name string
	b := dynamicstruct.NewBuilder()

	for k, v := range m {
		// TODO: apply initialisms theories. See: https://github.com/golang/go/wiki/CodeReviewComments#initialisms
		//   (and more golint theories validations)

		// TODO: add "omitempty"? (e.g. when key is missing, type should be a pointer and have "omitempty")
		// TODO: add ",string", ",boolean" extra options?
		// See: https://golang.org/pkg/encoding/json/#Marshal
		// See: https://m-zajac.github.io/json2go/
		if useTag {
			tag = fmt.Sprintf(`%s:"%s"`, d.dt.string(), k)
		}

		name = strcase.ToCamel(k)

		// See: https://golang.org/pkg/encoding/json/#Unmarshal
		switch value := v.(type) {
		case bool:
			b = b.AddBoolWithTag(name, tag)
		case float64:
			b = b.AddFloat64WithTag(name, tag)
		case string:
			b = b.AddStringWithTag(name, tag)
		case []interface{}:
			if len(value) > 0 {
				// FIXME: fix nest mode support or not
				switch vv := value[0].(type) {
				case map[string]interface{}:
					if nest {
						nds, err := d.toDsFromStringMap(vv, nest, useTag)
						if err != nil {
							return nil, err
						}
						b = b.AddDynamicStructSliceWithTag(name, nds, false, tag)
					} else {
						b = b.AddSliceWithTag(name, interface{}(vv), tag)
					}
				default:
					b = b.AddSliceWithTag(name, interface{}(vv), tag)
				}
			}
		case map[string]interface{}:
			if nest {
				nds, err := d.toDsFromStringMap(value, nest, useTag)
				if err != nil {
					return nil, err
				}
				b = b.AddDynamicStructWithTag(name, nds, false, tag)
			} else {
				for kk := range value {
					b = b.AddMapWithTag(name, kk, nil, tag)
					// only one addition
					break
				}
			}
		// YAML support
		case int:
			b = b.AddIntWithTag(name, tag)
		// YAML support
		case map[interface{}]interface{}:
			m := make(map[string]interface{})
			for k, v := range value {
				m[fmt.Sprintf("%v", k)] = v
			}

			if nest {
				nds, err := d.toDsFromStringMap(m, nest, useTag)
				if err != nil {
					return nil, err
				}
				b = b.AddDynamicStruct(name, nds, false)
			} else {
				for kk := range m {
					b = b.AddMapWithTag(name, kk, nil, tag)
					// only one addition
					break
				}
			}
		case nil:
			b = b.AddInterfaceWithTag(name, false, tag)
		default:
			return nil, fmt.Errorf("value %#v has invalid type. m is %#v", value, m)
		}
	}

	return b.Build()
}
