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
		var drElem *DecodedResult
		var dsOnce dynamicstruct.DynamicStruct
		iArr := make([]interface{}, len(t))
		for idx, elemIntf := range t {
			// call this function recursively
			// we want to build DynamicStruct only once
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

	// camelize key for building DynamicStruct with exported fields
	// e.g. if json item name is "hoge_huga", same field name in DynamicStruct is "HogeHuga"
	cm := camelizeMapKey(m)

	if dr.DynamicStruct == nil {
		dr.DynamicStruct, err = buildDynamicStruct(cm)
		if err != nil {
			return nil, err
		}
	}

	dr.DecodedInterface, err = dr.DynamicStruct.DecodeMap(cm)
	if err != nil {
		return nil, err
	}

	return dr, nil
}

func camelizeMapKey(m map[string]interface{}) map[string]interface{} {
	camelizedMap := make(map[string]interface{}, len(m))
	for k, v := range m {
		camelizedMap[strcase.ToCamel(k)] = v
	}

	return camelizedMap
}

func buildDynamicStruct(m map[string]interface{}) (dynamicstruct.DynamicStruct, error) {
	var tag string
	b := dynamicstruct.NewBuilder()

	for k, v := range m {
		tag = fmt.Sprintf(`json:"%s"`, k)

		// See: https://golang.org/pkg/encoding/json/#Unmarshal
		switch value := v.(type) {
		case bool:
			b = b.AddBoolWithTag(k, tag)
		case float64:
			b = b.AddFloat64WithTag(k, tag)
		case string:
			b = b.AddStringWithTag(k, tag)
		case []interface{}:
			b = b.AddSliceWithTag(k, interface{}(value[0]), tag)
		case map[string]interface{}:
			for kk, vv := range value {
				b = b.AddMapWithTag(k, kk, interface{}(vv), tag)
				break
			}
		case nil:
			// Note: Is this ok?
			b = b.AddInterfaceWithTag(k, false, tag)
		default:
			return nil, fmt.Errorf("jsonData %#v has invalid typed key", m)
		}
	}

	return b.Build(), nil
}
