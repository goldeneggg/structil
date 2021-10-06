package decoder

import (
	"fmt"

	"github.com/iancoleman/strcase"

	"github.com/goldeneggg/structil"
	"github.com/goldeneggg/structil/dynamicstruct"
)

// Decoder is the struct that decodes some marshaled data like JSON and YAML.
type Decoder struct {
	data     []byte // original data
	dt       dataType
	unm      interface{}            // unmarshaled result from data to JSON/YAML/etc
	unmMapsi map[string]interface{} // convert string map from unmarshaled result
	ds       *dynamicstruct.DynamicStruct
	dsi      interface{} // unmarshaled result from data to DynamicStruct
}

func newDecoder(data []byte, dt dataType) (d *Decoder, err error) {
	unm, err := dt.unmarshal(data)

	d = &Decoder{
		data:     data,
		dt:       dt,
		unm:      unm,
		unmMapsi: make(map[string]interface{}),
	}

	switch t := d.unm.(type) {
	case map[string]interface{}:
		// JSON
		d.unmMapsi = t
	case map[interface{}]interface{}:
		// YAML
		d.unmMapsi = mapiiToMapsi(t)
		// FIXME:
		// トップレベルが配列の場合、タイプに関わらず直接のunmarshalはエラーになってしまうので、
		// （暫定的に）0番目の要素を取り出してそれを処理するようにしている
		// という対応を取ろうとした際の名残。不要とわかったら消す
		// case []interface{}:
		// 	if len(t) > 0 {
		// 		switch tt := t[0].(type) {
		// 		case map[string]interface{}:
		// 			d.unmMapsi = tt
		// 		case map[interface{}]interface{}:
		// 			d.unmMapsi = mapiiToMapsi(tt)
		// 		}
		// 	}
	}

	return
}

// FromJSON returns a concrete Decoder for JSON.
func FromJSON(data []byte) (*Decoder, error) {
	return newDecoder(data, typeJSON)
}

// FromYAML returns a concrete Decoder for YAML.
func FromYAML(data []byte) (*Decoder, error) {
	return newDecoder(data, typeYAML)
}

// JSONToGetter returns a structil.Getter with a decoded JSON via DynamicStruct.
func JSONToGetter(data []byte, nest bool) (*structil.Getter, error) {
	d, err := FromJSON(data)
	if err != nil {
		return nil, err
	}

	// FIXME: when nest = true, failed to unmershal array_struct_field
	// "json: cannot unmarshal array into Go struct field .array_struct_field of type struct { Vvvv string "json:\"vvvv\""; Kkk string "json:\"kkk\"" }"
	return d.dsToGetter(nest)
}

// YAMLToGetter returns a structil.Getter with a decoded YAML via DynamicStruct.
// FIXME: this function has a problem caused by map[interface{}]interface{}.
func YAMLToGetter(data []byte, nest bool) (*structil.Getter, error) {
	d, err := FromYAML(data)
	if err != nil {
		return nil, err
	}

	// FIXME: when nest = true, failed to unmershal array_struct_field
	return d.dsToGetter(nest)
}

func (d *Decoder) dsToGetter(nest bool) (*structil.Getter, error) {
	ds, err := d.DynamicStruct(nest, true)
	if err != nil {
		return nil, err
	}

	dsi, err := d.decodeToDynamicStruct(ds)
	if err != nil {
		return nil, err
	}

	return structil.NewGetter(dsi)
}

func (d *Decoder) decodeToDynamicStruct(ds *dynamicstruct.DynamicStruct) (interface{}, error) {
	d.dsi = ds.NewInterface()
	if err := d.dt.unmarshalWithIPtr(d.data, &d.dsi); err != nil {
		return nil, err
	}

	return d.dsi, nil
}

// Data returns an original data as []byte.
func (d *Decoder) Data() []byte {
	return d.data
}

/*
// Interface returns a unmarshaled interface from original data.
func (d *Decoder) Interface() interface{} {
	return d.unm
}
*/

/*
// Map returns a map of unmarshaled interface from original data.
func (d *Decoder) Map() (map[string]interface{}, error) {
	return d.dt.intfToStringMap(d.Interface())
}
*/

// DynamicStruct returns a decoded DynamicStruct with unmarshaling data to DynamicStruct interface.
func (d *Decoder) DynamicStruct(nest bool, useTag bool) (*dynamicstruct.DynamicStruct, error) {
	var err error

	d.ds, err = d.toDs(d.unm, nest, useTag)
	if err != nil {
		return nil, err
	}

	return d.ds, err
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
		return d.toDsFromStringMap(mapiiToMapsi(t), nest, useTag)
	}

	return nil, fmt.Errorf("unexpected interface: %#v", i)
}

func (d *Decoder) toDsFromStringMap(m map[string]interface{}, nest bool, useTag bool) (*dynamicstruct.DynamicStruct, error) {
	var tag, name string
	var err error
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
					b, err = d.addForStringMap(b, vv, true, tag, name, nest, useTag)
					if err != nil {
						return nil, err
					}

					if nest {
						nds, err := d.toDsFromStringMap(vv, nest, useTag)
						if err != nil {
							return nil, err
						}
						b = b.AddDynamicStructSliceWithTag(name, nds, false, tag)
					} else {
						b = b.AddSliceWithTag(name, interface{}(vv), tag)
					}
				case map[interface{}]interface{}:
					m := mapiiToMapsi(vv)
					b, err = d.addForStringMap(b, m, true, tag, name, nest, useTag)
					if err != nil {
						return nil, err
					}
				default:
					b = b.AddSliceWithTag(name, interface{}(vv), tag)
				}
			}
		case map[string]interface{}:
			b, err = d.addForStringMap(b, value, false, tag, name, nest, useTag)
			if err != nil {
				return nil, err
			}

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
			m := mapiiToMapsi(value)
			b, err = d.addForStringMap(b, m, false, tag, name, nest, useTag)
			if err != nil {
				return nil, err
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

func (d *Decoder) addForStringMap(
	b *dynamicstruct.Builder,
	m map[string]interface{},
	forSlice bool,
	tag string,
	name string,
	nest bool,
	useTag bool) (*dynamicstruct.Builder, error) {

	// Note: forSlice judgement is top priority.
	// FIXME: support nested dynamicstruct for slice element
	// (currently, this causes a panic when field type is "array_struct")
	if nest {
		nds, err := d.toDsFromStringMap(m, nest, useTag)
		if err != nil {
			return b, err
		}
		b = b.AddDynamicStructWithTag(name, nds, false, tag)
	} else if forSlice {
		b = b.AddSliceWithTag(name, interface{}(m), tag)
	} else {
		for kk := range m {
			b = b.AddMapWithTag(name, kk, nil, tag)
			// only one addition
			break
		}
	}

	return b, nil
}

// convert map[interface{}]interface{} to map[string]interface{}
func mapiiToMapsi(mapii map[interface{}]interface{}) map[string]interface{} {
	mapsi := make(map[string]interface{})
	for k, v := range mapii {
		switch t := v.(type) {
		case map[interface{}]interface{}:
			mapsi[fmt.Sprintf("%v", k)] = mapiiToMapsi(t)
		default:
			mapsi[fmt.Sprintf("%v", k)] = v
		}
	}

	return mapsi
}
