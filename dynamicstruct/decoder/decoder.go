package decoder

import (
	"fmt"

	"github.com/iancoleman/strcase"

	"github.com/goldeneggg/structil"
	"github.com/goldeneggg/structil/dynamicstruct"
)

// Decoder is the struct that decodes some marshaled data like JSON and YAML.
type Decoder struct {
	dt        dataType
	orgData   []byte         // original data
	orgIntf   any            // unmarshaled interface from original data
	strKeyMap map[string]any // string key map for decoding to DymanicStruct
	ds        *dynamicstruct.DynamicStruct
	dsi       any // unmarshaled result from data to DynamicStruct
}

func newDecoder(data []byte, dt dataType) (*Decoder, error) {
	unm, err := dt.unmarshal(data)
	if err != nil {
		return nil, err
	}

	dec := &Decoder{
		dt:        dt,
		orgData:   data,
		orgIntf:   unm,
		strKeyMap: make(map[string]any),
	}

	switch t := dec.orgIntf.(type) {
	case map[string]any:
		// JSON
		dec.strKeyMap = t
	// Note: this is dead case with gopkg.in/yaml.v3 (but alive with v2)
	// case map[any]any:
	// 	// YAML
	// 	dec.strKeyMap = toStringKeyMap(t)
	case []any:
		if len(t) > 0 {
			// The items in the array must be same for all elements.
			// So the first element is used to process
			switch tt := t[0].(type) {
			case map[string]any:
				dec.strKeyMap = tt
			// Note: this is dead case with gopkg.in/yaml.v3 (but alive with v2)
			// case map[any]any:
			// 	dec.strKeyMap = toStringKeyMap(tt)
			default:
				return nil, fmt.Errorf("unexpected type of t[0] [%v]", tt)
			}
		}
	default:
		return nil, fmt.Errorf("unexpected type of dec.orgIntf [%v]", t)
	}

	return dec, nil
}

// FromJSON returns a concrete Decoder for JSON.
func FromJSON(data []byte) (*Decoder, error) {
	return newDecoder(data, typeJSON)
}

// FromYAML returns a concrete Decoder for YAML.
func FromYAML(data []byte) (*Decoder, error) {
	return newDecoder(data, typeYAML)
}

// FromXML returns a concrete Decoder for XML.
// FIXME: This function is still a future candidate (returned error now)
func FromXML(data []byte) (*Decoder, error) {
	return newDecoder(data, end) // FIXME: "end" is provisional type
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

func (d *Decoder) decodeToDynamicStruct(ds *dynamicstruct.DynamicStruct) (any, error) {
	// toDsFromStringMap() method uses "Build()" method and this means that ds is build by pointer-mode
	// So ds.NewInterface() returns a struct *pointer*
	d.dsi = ds.NewInterface()

	// must use "d.strKeyMap" (not "d.orgIntf"). because key of "d.orgIntf" is not string but any
	data, err := d.dt.marshal(d.strKeyMap)
	if err != nil {
		return nil, fmt.Errorf("fail to d.dt.marshal: %w", err)
	}

	// must use "d.dsi" (not "&d.dsi"). because "d.dsi" is pointer
	// if use "&.d.dsi", unmarshal result is not struct but map[any]interface when dt is YAML
	if err := d.dt.unmarshalWithIPtr(data, d.dsi); err != nil {
		return nil, fmt.Errorf("fail to d.dt.unmarshalWithIPtr: %w", err)
	}

	return d.dsi, nil
}

// OrgData returns an original data as []byte.
func (d *Decoder) OrgData() []byte {
	return d.orgData
}

// DynamicStruct returns a decoded DynamicStruct with unmarshaling data to DynamicStruct interface.
func (d *Decoder) DynamicStruct(nest bool, useTag bool) (*dynamicstruct.DynamicStruct, error) {
	var err error

	// d.ds, err = d.toDs(d.orgIntf, nest, useTag)
	d.ds, err = d.toDs(d.strKeyMap, nest, useTag)
	if err != nil {
		return nil, err
	}

	return d.ds, err
}

func (d *Decoder) toDs(i any, nest bool, useTag bool) (*dynamicstruct.DynamicStruct, error) {
	switch t := i.(type) {
	case map[string]any:
		return d.toDsFromStringMap(t, nest, useTag)
	}

	return nil, fmt.Errorf("unsupported type [%T] for toDs", i)
}

func (d *Decoder) toDsFromStringMap(m map[string]any, nest bool, useTag bool) (*dynamicstruct.DynamicStruct, error) {
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
		case []any:
			if len(value) > 0 {
				switch vv := value[0].(type) {
				case map[string]any:
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
						b = b.AddSliceWithTag(name, any(vv), tag)
					}
				default:
					// FIXME: 配列要素を全て "any" にキャストしているが、型を明示したい
					b = b.AddSliceWithTag(name, any(vv), tag)
				}
			}
		case map[string]any:
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
		default:
			return nil, fmt.Errorf("unsupported type of map-value. key = [%s] value = %#v", k, value)
		}
	}

	return b.Build()
}

func (d *Decoder) addForStringMap(
	b *dynamicstruct.Builder,
	m map[string]any,
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
		b = b.AddSliceWithTag(name, any(m), tag)
	} else {
		for kk := range m {
			b = b.AddMapWithTag(name, kk, nil, tag)
			// only one addition
			break
		}
	}

	return b, nil
}

// Note: this is dead case with gopkg.in/yaml.v3 (but alive with v2)
// convert map[any]any to map[string]any
// func toStringKeyMap(mapii map[any]any) map[string]any {
// 	mapsi := make(map[string]any)
// 	for k, v := range mapii {
// 		switch vt := v.(type) {
// 		case []any:
// 			// for nest array
// 			mapsi[fmt.Sprintf("%v", k)] = fromArrayToMapValue(vt)
// 		case map[any]any:
// 			// for nest object
// 			mapsi[fmt.Sprintf("%v", k)] = toStringKeyMap(vt)
// 		default:
// 			mapsi[fmt.Sprintf("%v", k)] = v
// 		}
// 	}

// 	return mapsi
// }

// Note: this is dead case with gopkg.in/yaml.v3 (but alive with v2)
// func fromArrayToMapValue(ia []any) any {
// 	resIa := make([]any, 0, len(ia))
// 	for _, iv := range ia {
// 		switch ivt := iv.(type) {
// 		case []any:
// 			// for nest array
// 			resIa = append(resIa, fromArrayToMapValue(ivt))
// 		case map[any]any:
// 			// for nest object
// 			// !!! this is important process for map[any]any to map[string]any for JSON unmarshaling
// 			resIa = append(resIa, toStringKeyMap(ivt))
// 		default:
// 			resIa = append(resIa, ivt)
// 		}
// 	}

// 	return resIa
// }
