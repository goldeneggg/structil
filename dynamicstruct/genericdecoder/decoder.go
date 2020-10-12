package genericdecoder

import (
	"fmt"

	"github.com/iancoleman/strcase"

	"github.com/goldeneggg/structil/dynamicstruct"
)

// GenericDecoder is the interface for decoding generic data(JSON, YAML, and more).
type GenericDecoder interface {
	Decode(data []byte) (*DecodedResult, error)
}

// DecodedResult is the result of GenericDecoder.Decode.
type DecodedResult struct {
	dynamicstruct.DynamicStruct
	DecodedInterface interface{}
}

// ui must be a unmarshalled interface from JSON, and others
func decode(ui interface{}, ds dynamicstruct.DynamicStruct) (*DecodedResult, error) {
	var err error

	switch t := ui.(type) {
	case map[string]interface{}:
		return decodeMap(t, ds)
	case []interface{}:
		// TODO: should check length and if length == 1, then call decodeMap directly and once instead of current implementation.
		var drElem *DecodedResult
		var dsOnce dynamicstruct.DynamicStruct
		iArr := make([]interface{}, len(t))
		for idx, elemIntf := range t {
			// call this function recursively
			// we want to build DynamicStruct only once
			// FIXME: current code can not support "omitempty" field for JSON array
			// TODO: DynamicStruct bulding not only once but only once *with omitempty support*
			drElem, err = decode(elemIntf, dsOnce)
			if err != nil {
				return nil, err
			}
			dsOnce = drElem.DynamicStruct

			iArr[idx] = drElem.DecodedInterface
		}

		return &DecodedResult{
			DynamicStruct:    dsOnce,
			DecodedInterface: iArr,
		}, nil
	}

	return nil, fmt.Errorf("unexpected return. unmarshalledJSON %+v is not map or array", ui)
}

func decodeMap(m map[string]interface{}, ds dynamicstruct.DynamicStruct) (*DecodedResult, error) {
	dr := &DecodedResult{
		DynamicStruct: ds,
	}
	var err error

	// camelizedKeys for building DynamicStruct with exported fields
	// e.g. if json item name is "hoge_huga", same field name in DynamicStruct is "HogeHuga"
	// FIXME: support case that input json field names are not snake_case but camelCase
	camelizedKeys, camelizedMap := camelizeMap(m)

	if dr.DynamicStruct == nil {
		dr.DynamicStruct, err = buildDynamicStruct(m, camelizedKeys)
		if err != nil {
			return nil, err
		}
	}

	dr.DecodedInterface, err = dr.DynamicStruct.DecodeMap(camelizedMap)
	if err != nil {
		return nil, err
	}

	return dr, nil
}

func camelizeMap(m map[string]interface{}) (map[string]string, map[string]interface{}) {
	camelizedKeys := make(map[string]string, len(m))
	camelizedMap := make(map[string]interface{}, len(m))

	for k, v := range m {
		camelizedKeys[k] = strcase.ToCamel(k)
		camelizedMap[camelizedKeys[k]] = v
	}

	return camelizedKeys, camelizedMap
}

func buildDynamicStruct(m map[string]interface{}, camelizedKeys map[string]string) (dynamicstruct.DynamicStruct, error) {
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
		name = camelizedKeys[k]

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
		case nil:
			// Note: Is this ok?
			b = b.AddInterfaceWithTag(name, false, tag)
		default:
			return nil, fmt.Errorf("jsonData %#v has invalid typed key", m)
		}
	}

	return b.Build(), nil
}
