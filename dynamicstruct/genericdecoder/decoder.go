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

func decode(intf interface{}, ds dynamicstruct.DynamicStruct) (*DecodedResult, error) {
	var err error

	switch t := intf.(type) {
	case map[string]interface{}:
		return newDecodedResult(t, ds)
	case []interface{}:
		var drElem *DecodedResult
		dsOnce := ds
		iArr := make([]interface{}, len(t))
		for idx, elemIntf := range t {
			// call this function recursively
			// we want to build DynamicStruct once, so "ds" argument is assigned
			drElem, err = decode(elemIntf, dsOnce)
			if err != nil {
				return nil, err
			}

			dsOnce = drElem.DynamicStruct
			iArr[idx] = drElem.DecodedInterface
		}

		return &DecodedResult{
			DynamicStruct:    ds,
			DecodedInterface: iArr,
		}, nil
	}

	return nil, fmt.Errorf("unexpected return. unmarshalledJSON %+v is not map or array", intf)
}

func newDecodedResult(m map[string]interface{}, ds dynamicstruct.DynamicStruct) (*DecodedResult, error) {
	dr := &DecodedResult{}
	var err error

	cm := camelizeMapKey(m)

	if ds != nil {
		dr.DynamicStruct = ds
	} else {
		dr.DynamicStruct, err = dynamicStructFromMap(cm)
		if err != nil {
			return nil, err
		}
	}

	dr.DecodedInterface, err = decodeMap(dr.DynamicStruct, cm)
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

func dynamicStructFromMap(m map[string]interface{}) (dynamicstruct.DynamicStruct, error) {
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

	ds := b.Build()

	return ds, nil
}

func decodeMap(ds dynamicstruct.DynamicStruct, m map[string]interface{}) (interface{}, error) {
	intf, err := ds.DecodeMap(m)
	if err != nil {
		return nil, err
	}

	return intf, nil
}
