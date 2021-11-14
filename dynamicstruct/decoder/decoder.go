package decoder

import (
	"fmt"

	"github.com/iancoleman/strcase"

	"github.com/goldeneggg/structil"
	"github.com/goldeneggg/structil/dynamicstruct"
)

// Decoder is the struct that decodes some marshaled data like JSON and YAML.
type Decoder struct {
	dt          dataType
	data        []byte // original data
	unm         interface{}
	dsDecodeMap map[string]interface{} // map for decoding to DymanicStruct
	ds          *dynamicstruct.DynamicStruct
	dsi         interface{} // unmarshaled result from data to DynamicStruct
}

func newDecoder(data []byte, dt dataType) (d *Decoder, err error) {
	unm, err := dt.unmarshal(data)

	d = &Decoder{
		dt:          dt,
		data:        data,
		unm:         unm,
		dsDecodeMap: make(map[string]interface{}),
	}

	switch t := d.unm.(type) {
	case map[string]interface{}:
		// JSON
		d.dsDecodeMap = t
	case map[interface{}]interface{}:
		// YAML
		d.dsDecodeMap = toStringKeyMap(t)
	// FIXME:
	// for top-level is array
	// （暫定的に）0番目の要素を取り出してそれを処理するようにしているが、どうすべきか？
	case []interface{}:
		if len(t) > 0 {
			switch tt := t[0].(type) {
			case map[string]interface{}:
				d.dsDecodeMap = tt
			case map[interface{}]interface{}:
				d.dsDecodeMap = toStringKeyMap(tt)
			}
		}
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
	// must use "d.dsDecodeMap" (not "d.unm"). because key of "d.unm" is not string but interface{}
	data, err := d.dt.marshal(d.dsDecodeMap)
	if err != nil {
		return nil, fmt.Errorf("fail to d.dt.marshal: %w", err)
	}

	// toDsFromStringMap() method uses "Build()" method and this means that ds is build by pointer-mode
	// So ds.NewInterface() returns a struct pointer
	d.dsi = ds.NewInterface()

	// must use "d.dsi" (not "&d.dsi"). because "d.dsi" is pointer
	// if use "&.d.dsi", unmarshal result is not struct but map[interface{}]interface when dt is YAML
	if err := d.dt.unmarshalWithIPtr(data, d.dsi); err != nil {
		return nil, fmt.Errorf("fail to d.dt.unmarshalWithIPtr: %w", err)
	}

	return d.dsi, nil
}

// Data returns an original data as []byte.
func (d *Decoder) Data() []byte {
	return d.data
}

// DynamicStruct returns a decoded DynamicStruct with unmarshaling data to DynamicStruct interface.
func (d *Decoder) DynamicStruct(nest bool, useTag bool) (*dynamicstruct.DynamicStruct, error) {
	var err error

	// d.ds, err = d.toDs(d.unm, nest, useTag)
	d.ds, err = d.toDs(d.dsDecodeMap, nest, useTag)
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
		return d.toDsFromStringMap(toStringKeyMap(t), nest, useTag)
	}

	return nil, fmt.Errorf("unsupported type [%T] for toDs", i)
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

		// FIXME: the first character of k should be only alpha-numeric (e.g. "@" is invalid as first character)
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
						b = b.AddDynamicStructSliceWithTag(name, nds, false, tag)
					} else {
						b = b.AddSliceWithTag(name, interface{}(vv), tag)
					}
				case map[interface{}]interface{}:
					m := toStringKeyMap(vv)
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
			m := toStringKeyMap(value)
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

// convert map[interface{}]interface{} to map[string]interface{}
func toStringKeyMap(mapii map[interface{}]interface{}) map[string]interface{} {
	mapsi := make(map[string]interface{})
	for k, v := range mapii {
		switch vt := v.(type) {
		case []interface{}:
			// for nest array
			mapsi[fmt.Sprintf("%v", k)] = fromArrayToMapValue(vt)
		case map[interface{}]interface{}:
			// for nest object
			mapsi[fmt.Sprintf("%v", k)] = toStringKeyMap(vt)
		default:
			mapsi[fmt.Sprintf("%v", k)] = v
		}
	}

	return mapsi
}

func fromArrayToMapValue(ia []interface{}) interface{} {
	resIa := make([]interface{}, 0, len(ia))
	for _, iv := range ia {
		switch ivt := iv.(type) {
		case []interface{}:
			// for nest array
			resIa = append(resIa, fromArrayToMapValue(ivt))
		case map[interface{}]interface{}:
			// for nest object
			// !!! this is important process for map[interface{}]interface{} to map[string]interface{} for JSON unmarshaling
			resIa = append(resIa, toStringKeyMap(ivt))
		default:
			resIa = append(resIa, ivt)
		}
	}

	return resIa
}
