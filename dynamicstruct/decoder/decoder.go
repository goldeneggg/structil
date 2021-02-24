package decoder

import (
	"fmt"

	"github.com/iancoleman/strcase"

	"github.com/goldeneggg/structil/dynamicstruct"
)

// DecodedResult is the result of Decoder.Decode.
// deprecated
type DecodedResult struct {
	*dynamicstruct.DynamicStruct
	Interface interface{}
}

type Decoder struct {
	data []byte
	dt   DataType
	unm  interface{}
	ds   dynamicstruct.DynamicStruct
}

func NewJSON(data []byte) (*Decoder, error) {
	return New(data, TypeJSON)
}

func NewYAML(data []byte) (*Decoder, error) {
	return New(data, TypeYAML)
}

func New(data []byte, dt DataType) (d *Decoder, err error) {
	var intf interface{}
	err = dt.Unmarshal(data, &intf)

	d = &Decoder{
		data: data,
		dt:   dt,
		unm:  intf,
	}

	return
}

func (d *Decoder) DynamicStruct(nest bool, useTag bool) (*dynamicstruct.DynamicStruct, error) {
	return d.toDs(d.unm, nest, useTag)
}

func (d *Decoder) toDs(i interface{}, nest bool, useTag bool) (*dynamicstruct.DynamicStruct, error) {
	switch t := i.(type) {
	case map[string]interface{}:
		return d.toDsFromStringMap(t, nest, useTag)
	// FIXME: map[interface{}]interface{} support (for YAML)
	// case map[interface{}]interface{}:
	// 	m := make(map[string]interface{})
	// 	for k, v := range t {
	// 		m[fmt.Sprintf("%v", k)] = v
	// 	}
	// 	return decodeMap(m, dt, ds)
	case []interface{}:
		if len(t) > 0 {
			if len(t) == 1 {
				return d.toDs(t[0], nest, useTag)
			}

			// TODO: seek an element that have max size of t. And call d.toDs with this element
			// 配列内の構造が可変なケースを考慮して、最も大きい構造の要素を取り出してその要素に対してtoDsを呼ぶようにする
			return d.toDs(t[0], nest, useTag)
		}
	}

	return nil, fmt.Errorf("unexpected interface: %#v", i)
}

func (d *Decoder) toDsFromStringMap(m map[string]interface{}, nest bool, useTag bool) (*dynamicstruct.DynamicStruct, error) {
	var tag, name string
	b := dynamicstruct.NewBuilder()

	for k, v := range m {
		// TODO: "json" changes dynamic. e.g. "yaml", "xml" and others
		// TODO: apply initialisms theories. See: https://github.com/golang/go/wiki/CodeReviewComments#initialisms
		//   (and more golint theories validations)

		// TODO: add "omitempty"? (e.g. when key is missing, type should be a pointer and have "omitempty")
		// TODO: add ",string", ",boolean" extra options?
		// See: https://golang.org/pkg/encoding/json/#Marshal
		// See: https://m-zajac.github.io/json2go/
		if useTag {
			tag = fmt.Sprintf(`%s:"%s"`, d.dt.String(), k)
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
			// FIXME: fix nest mode support or not
			b = b.AddSliceWithTag(name, interface{}(value[0]), tag)
		case map[string]interface{}:
			// TODO: nest mode support
			if nest {
				nds, err := d.toDsFromStringMap(value, nest, useTag)
				if err != nil {
					return nil, err
				}
				b = b.AddDynamicStructPtr(name, nds)
			} else {
				for kk := range value {
					b = b.AddMapWithTag(name, kk, nil, tag)
					// only one addition
					break
				}
			}

		// case map[interface{}]interface{}:
		// 	m := make(map[string]interface{})
		// 	for k, v := range value {
		// 		m[fmt.Sprintf("%v", k)] = v
		// 	}

		// 	for kk, vv := range m {
		// 		b = b.AddMapWithTag(name, kk, interface{}(vv), tag)
		// 		break
		// 	}
		case nil:
			// Note: Is this ok?
			b = b.AddInterfaceWithTag(name, false, tag)
		default:
			return nil, fmt.Errorf("value %#v has invalid type. m is %#v", value, m)
		}
	}

	return b.Build()
}

/*
// Decode decodes original data to interface via DynamicStruct.DecodeMap.
// data argument must be a byte array data of valid format(JSON, YAML, TOML).
func Decode(data []byte, dt DataType) (*DecodedResult, error) {
	var intf interface{}
	var err error

	switch dt {
	case TypeJSON:
		err = unmarshalJSON(data, &intf)
	case TypeYAML:
		err = unmarshalYAML(data, &intf)
	default:
		err = fmt.Errorf("invalid datatype: %v", dt)
	}
	if err != nil {
		return nil, err
	}

	dr, err := decode(intf, dt, nil)
	if err != nil {
		return nil, err
	}

	return dr, nil
}

func unmarshalJSON(data []byte, ptr interface{}) error {
	// FIXME:
	// want to add json validation. but is json.Valid(data) too slow?
	// See: https://stackoverflow.com/questions/22128282/how-to-check-string-is-in-json-format
	return json.Unmarshal(data, ptr)
}

func unmarshalYAML(data []byte, ptr interface{}) error {
	return yaml.Unmarshal(data, ptr)
}

// ui must be a unmarshalled interface from JSON, and others
func decode(ui interface{}, dt DataType, ds *dynamicstruct.DynamicStruct) (*DecodedResult, error) {
	var err error

	switch t := ui.(type) {
	case map[string]interface{}:
		return decodeMap(t, dt, ds)
	// case map[interface{}]interface{}:
	// 	m := make(map[string]interface{})
	// 	for k, v := range t {
	// 		m[fmt.Sprintf("%v", k)] = v
	// 	}
	// 	return decodeMap(m, dt, ds)
	case []interface{}:
		// TODO: should check length and if length == 1, then call decodeMap directly and once instead of current implementation.
		var dr *DecodedResult
		var ds *dynamicstruct.DynamicStruct
		iArr := make([]interface{}, len(t))
		for idx, elemIntf := range t {
			// call this function recursively
			// we want to build DynamicStruct only once
			// FIXME: current code can not support "omitempty" field for JSON array
			// TODO: DynamicStruct bulding not only once but only once *with omitempty support*
			dr, err = decode(elemIntf, dt, ds)
			if err != nil {
				return nil, err
			}
			ds = dr.DynamicStruct

			iArr[idx] = dr.Interface
		}

		return &DecodedResult{
			DynamicStruct: ds,
			Interface:     iArr,
		}, nil
	}

	return nil, fmt.Errorf("unexpected unmarshalled interface: %#v", ui)
}

func decodeMap(m map[string]interface{}, dt DataType, ds *dynamicstruct.DynamicStruct) (*DecodedResult, error) {
	dr := &DecodedResult{
		DynamicStruct: ds,
	}
	var err error

	// ck for building DynamicStruct with exported fields
	// e.g. if json item name is "hoge_huga", same field name in DynamicStruct is "HogeHuga"
	// FIXME: support case that input json field names are not snake_case but camelCase
	cKeys, cMap := camelizeMap(m)

	if dr.DynamicStruct == nil {
		dr.DynamicStruct, err = buildDynamicStruct(m, cKeys, dt)
		if err != nil {
			return nil, err
		}
	}

	dr.Interface, err = dr.DynamicStruct.DecodeMap(cMap)
	if err != nil {
		return nil, err
	}

	return dr, nil
}

// TODO: implement
func fromStringMapToDynamicStruct(m map[string]interface{}, recur bool, format string) (*DynamicStruct, error) {
	var tag, name string
	b := dynamicstruct.NewBuilder()

	for k, v := range m {
		// TODO: "json" changes dynamic. e.g. "yaml", "xml" and others
		// TODO: apply initialisms theories. See: https://github.com/golang/go/wiki/CodeReviewComments#initialisms
		//   (and more golint theories validations)

		// TODO: add "omitempty"? (e.g. when key is missing, type should be a pointer and have "omitempty")
		// TODO: add ",string", ",boolean" extra options?
		// See: https://golang.org/pkg/encoding/json/#Marshal
		// See: https://m-zajac.github.io/json2go/
		if format != "" {
			tag = fmt.Sprintf(`%s:"%s"`, format, k)
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
			b = b.AddSliceWithTag(name, interface{}(value[0]), tag)
		case map[string]interface{}:
			for kk, vv := range value {
				b = b.AddMapWithTag(name, kk, interface{}(vv), tag)
				break
			}
		// case map[interface{}]interface{}:
		// 	m := make(map[string]interface{})
		// 	for k, v := range value {
		// 		m[fmt.Sprintf("%v", k)] = v
		// 	}

		// 	for kk, vv := range m {
		// 		b = b.AddMapWithTag(name, kk, interface{}(vv), tag)
		// 		break
		// 	}
		case nil:
			// Note: Is this ok?
			b = b.AddInterfaceWithTag(name, false, tag)
		default:
			return nil, fmt.Errorf("value %#v has invalid type. m is %#v", value, m)
		}
	}

	return b.Build()
}

// deprecated
func camelizeMap(m map[string]interface{}) (map[string]string, map[string]interface{}) {
	cKeys := make(map[string]string, len(m))
	cMap := make(map[string]interface{}, len(m))

	for k, v := range m {
		cKeys[k] = strcase.ToCamel(k)
		cMap[cKeys[k]] = v
	}

	return cKeys, cMap
}

// deprecated
func buildDynamicStruct(m map[string]interface{}, cKeys map[string]string, dt DataType) (*dynamicstruct.DynamicStruct, error) {
	var tag, name string
	b := dynamicstruct.NewBuilder()

	for k, v := range m {
		// TODO: "json" changes dynamic. e.g. "yaml", "xml" and others
		// TODO: apply initialisms theories. See: https://github.com/golang/go/wiki/CodeReviewComments#initialisms
		//   (and more golint theories validations)
		// TODO: add "omitempty"? (e.g. when key is missing, type should be a pointer and have "omitempty")
		// TODO: add ",string", ",boolean" extra options?
		// See: https://golang.org/pkg/encoding/json/#Marshal
		// See: https://m-zajac.github.io/json2go/
		tag = fmt.Sprintf(`json:"%s"`, k)
		name = cKeys[k]

		// See: https://golang.org/pkg/encoding/json/#Unmarshal
		switch value := v.(type) {
		case bool:
			b = b.AddBoolWithTag(name, tag)
		case float64:
			b = b.AddFloat64WithTag(name, tag)
		case string:
			b = b.AddStringWithTag(name, tag)
		case []interface{}:
			b = b.AddSliceWithTag(name, interface{}(value[0]), tag)
		case map[string]interface{}:
			for kk, vv := range value {
				b = b.AddMapWithTag(name, kk, interface{}(vv), tag)
				break
			}
		// case map[interface{}]interface{}:
		// 	m := make(map[string]interface{})
		// 	for k, v := range value {
		// 		m[fmt.Sprintf("%v", k)] = v
		// 	}

		// 	for kk, vv := range m {
		// 		b = b.AddMapWithTag(name, kk, interface{}(vv), tag)
		// 		break
		// 	}
		case nil:
			// Note: Is this ok?
			b = b.AddInterfaceWithTag(name, false, tag)
		default:
			return nil, fmt.Errorf("value %#v has invalid type. m is %#v", value, m)
		}
	}

	return b.Build()
}
*/
