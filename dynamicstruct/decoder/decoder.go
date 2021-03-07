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
				b = b.AddDynamicStruct(name, nds, false)
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
