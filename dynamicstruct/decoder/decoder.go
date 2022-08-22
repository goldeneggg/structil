package decoder

import (
	"fmt"
	"math/big"

	"github.com/iancoleman/strcase"

	"github.com/goldeneggg/structil"
	"github.com/goldeneggg/structil/dynamicstruct"
)

// Decoder is the struct that decodes some marshaled data like JSON and YAML.
type Decoder struct {
	typ  dataType
	data []byte
	intf interface{}
	m    map[string]interface{}
	ds   *dynamicstruct.DynamicStruct
}

// FromJSON returns a concrete Decoder for JSON.
func FromJSON(data []byte) (*Decoder, error) {
	return newDecoder(data, typeJSON)
}

// FromYAML returns a concrete Decoder for YAML.
func FromYAML(data []byte) (*Decoder, error) {
	return newDecoder(data, typeYAML)
}

// FromHCL returns a concrete Decoder for HCL.
func FromHCL(data []byte) (*Decoder, error) {
	return newDecoder(data, typeHCL)
}

// FromXML returns a concrete Decoder for XML.
// FIXME: This function is still a future candidate (returned error now)
func FromXML(data []byte) (*Decoder, error) {
	return newDecoder(data, end) // FIXME: "end" is provisional type
}

func newDecoder(data []byte, typ dataType) (*Decoder, error) {
	intf, err := typ.unmarshal(data)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	switch t := intf.(type) {
	case map[string]interface{}:
		// JSON
		m = t
	case []interface{}:
		if len(t) > 0 {
			// The items in the array must be same for all elements.
			// So the first element is used to process
			switch tt := t[0].(type) {
			case map[string]interface{}:
				m = tt
			default:
				return nil, fmt.Errorf("unexpected type of t[0] [%v]", tt)
			}
		}
	default:
		return nil, fmt.Errorf("unexpected type of unmashalled interface [%v]", t)
	}

	return &Decoder{
		typ:  typ,
		data: data,
		intf: intf,
		m:    m,
	}, nil
}

// JSONToGetter returns a structil.Getter with a decoded JSON via DynamicStruct.
func JSONToGetter(data []byte, nest bool) (*structil.Getter, error) {
	d, err := FromJSON(data)
	if err != nil {
		return nil, fmt.Errorf("fail to JSONToGetter: %w", err)
	}

	return d.dsToGetter(nest)
}

// YAMLToGetter returns a structil.Getter with a decoded YAML via DynamicStruct.
func YAMLToGetter(data []byte, nest bool) (*structil.Getter, error) {
	d, err := FromYAML(data)
	if err != nil {
		return nil, fmt.Errorf("fail to YAMLToGetter: %w", err)
	}

	return d.dsToGetter(nest)
}

// FIXME: HCL decode method supports only struct pointer or map pointer
// func HCLToGetter(data []byte, nest bool) (*structil.Getter, error) {
// 	d, err := FromHCL(data)
// 	if err != nil {
// 		return nil, fmt.Errorf("fail to HCLToGetter: %w", err)
// 	}

// 	return d.dsToGetter(nest)
// }

// DynamicStruct returns a decoded DynamicStruct with unmarshaling data to DynamicStruct interface.
func (d *Decoder) DynamicStruct(nest bool, useTag bool) (*dynamicstruct.DynamicStruct, error) {
	var err error

	d.ds, err = d.toDsFromStringMap(d.m, nest, useTag)
	if err != nil {
		return nil, err
	}

	return d.ds, err
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
	data, err := d.typ.marshal(d.m)
	if err != nil {
		return nil, fmt.Errorf("fail to d.typ.marshal: %w", err)
	}

	// ds.NewInterface() returns a struct *pointer*
	dsi := ds.NewInterface()
	// must use "dsi" (not "&dsi"). because "dsi" is pointer
	err = d.typ.unmarshalWithPtr(data, dsi)
	if err != nil {
		return nil, fmt.Errorf("fail to d.typ.unmarshalWithPtr: %w", err)
	}

	return dsi, nil
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
			tag = fmt.Sprintf(`%s:"%s"`, d.typ.string(), k)
		}

		// FIXME: the first character of k should be only alpha-numeric
		// "@" や "/" は置換対応が必要かも
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
						b = b.AddDynamicStructSliceWithTag(name, nds, tag)
					} else {
						b = b.AddSliceWithTag(name, interface{}(vv), tag)
					}
				default:
					// FIXME: 配列要素を全て "interface{}" にキャストしているが、型を明示したい
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
		case nil:
			b = b.AddInterfaceWithTag(name, false, tag)
		// HCL support
		case *big.Float:
			b = b.AddFloat64WithTag(name, tag)
		default:
			return nil, fmt.Errorf("unsupported type of map-value. key = [%s] value = %#v", k, value)
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
	// FIXME: スライス内の要素のDynamicStruct変換サポート
	// (currently, this causes a panic when field type is "array_struct")
	if nest {
		nds, err := d.toDsFromStringMap(m, nest, useTag)
		if err != nil {
			return b, fmt.Errorf("addForStringMap nest = [%t], useTag = [%t]: %w", nest, useTag, err)
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
