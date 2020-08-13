package dynamicstruct

import (
	"encoding/json"
	"fmt"

	"github.com/iancoleman/strcase"
)

// GenericDecoder is the decoded object from unknown format JSON/YAML/and others.
type GenericDecoder struct {
	JSONData []byte
	DynamicStruct
	DecodedInterface interface{}
}

// NewGenericDecoder returns a concrete GenericDecoder
// jsonData argument must be a byte array data of valid JSON.
func NewGenericDecoder(jsonData []byte) (*GenericDecoder, error) {
	// FIXME:
	// want to add json validation. but is json.Valid(data) too slow?
	// See: https://stackoverflow.com/questions/22128282/how-to-check-string-is-in-json-format

	return &GenericDecoder{
		JSONData: jsonData,
	}, nil
}

// Decode decodes JSON data to interface via DynamicStruct.DecodeMap.
func (gd *GenericDecoder) Decode() error {
	var unmarshalledJSON interface{}
	err := json.Unmarshal(gd.JSONData, &unmarshalledJSON)
	if err != nil {
		return err
	}

	i, err := gd.parseUnmarshalled(unmarshalledJSON)
	if err != nil {
		return err
	}

	gd.DecodedInterface = i

	return nil
}

func (gd *GenericDecoder) parseUnmarshalled(unmarshalled interface{}) (interface{}, error) {
	switch t := unmarshalled.(type) {
	case map[string]interface{}:
		return gd.buildAndDecode(t)
	case []interface{}:
		var i interface{}
		var err error
		iArr := make([]interface{}, len(t))
		for idx, elemJSON := range t {
			// call this function recursively
			i, err = gd.parseUnmarshalled(elemJSON)
			if err != nil {
				return nil, err
			}

			iArr[idx] = i
		}

		return iArr, nil
	}

	return nil, fmt.Errorf("unexpected return. unmarshalledJSON %+v is not map or array", unmarshalled)
}

func (gd *GenericDecoder) buildAndDecode(m map[string]interface{}) (interface{}, error) {
	var camelizedKey, tag string
	camelizedFieldMap := make(map[string]interface{}, len(m))
	b := NewBuilder()

	for k, v := range m {
		camelizedKey = strcase.ToCamel(k)
		camelizedFieldMap[camelizedKey] = v
		tag = fmt.Sprintf(`json:"%s"`, k)

		switch value := v.(type) {
		case bool:
			b = b.AddBoolWithTag(camelizedKey, tag)
		case float64:
			b = b.AddFloat64WithTag(camelizedKey, tag)
		case string:
			b = b.AddStringWithTag(camelizedKey, tag)
		case []interface{}:
			b = b.AddSliceWithTag(camelizedKey, interface{}(value[0]), tag)
		case map[string]interface{}:
			for kk, vv := range value {
				b = b.AddMapWithTag(camelizedKey, kk, interface{}(vv), tag)
				break
			}
		case nil:
			// Note: Is this ok?
			b = b.AddInterfaceWithTag(camelizedKey, false, tag)
		default:
			return nil, fmt.Errorf("jsonData %#v has invalid typed key", m)
		}
	}

	ds := b.Build()
	intf, err := ds.DecodeMap(camelizedFieldMap)
	if err != nil {
		return nil, err
	}

	return intf, nil
}

// JSONToDynamicStructInterface returns an interface via DynamicStruct.DecodeMap from JSON data.
// jsonData argument must be a byte array data of JSON.
//
// This method supports known format JSON and unknown format JSON.
// But when JSON format is known, this method is not recommended. Because this method is suitable for unknown JSON with heavy and slow reflection functions.
//
// Field names in DynamicStruct are converted to CamelCase automatically
// - e.g. "hoge" JSON field is converted to "Hoge".
// - e.g. "huga_field" JSON field is converted to "HugaField".
func JSONToDynamicStructInterface(jsonData []byte) (interface{}, error) {
	// FIXME:
	// want to add json validation. but is json.Valid(data) too slow?
	// See: https://stackoverflow.com/questions/22128282/how-to-check-string-is-in-json-format

	var unmarshalledJSON interface{}
	err := json.Unmarshal(jsonData, &unmarshalledJSON)
	if err != nil {
		return nil, err
	}

	return parseUnmarshalledJSON(unmarshalledJSON)
}

func parseUnmarshalledJSON(unmarshalledJSON interface{}) (interface{}, error) {
	switch t := unmarshalledJSON.(type) {
	case map[string]interface{}:
		return mapToDynamicStructInterface(t)
	case []interface{}:
		var i interface{}
		var err error
		iArr := make([]interface{}, len(t))
		for idx, elemJSON := range t {
			// call this function recursively
			i, err = parseUnmarshalledJSON(elemJSON)
			if err != nil {
				return nil, err
			}

			iArr[idx] = i
		}

		return iArr, nil
	}

	return nil, fmt.Errorf("unexpected return. unmarshalledJSON %+v is not map or array", unmarshalledJSON)
}

func mapToDynamicStructInterface(m map[string]interface{}) (interface{}, error) {
	var camelizedKey, tag string
	camelizedFieldMap := make(map[string]interface{}, len(m))
	b := NewBuilder()

	for k, v := range m {
		camelizedKey = strcase.ToCamel(k)
		camelizedFieldMap[camelizedKey] = v
		tag = fmt.Sprintf(`json:"%s"`, k)

		switch value := v.(type) {
		case bool:
			b = b.AddBoolWithTag(camelizedKey, tag)
		case float64:
			b = b.AddFloat64WithTag(camelizedKey, tag)
		case string:
			b = b.AddStringWithTag(camelizedKey, tag)
		case []interface{}:
			b = b.AddSliceWithTag(camelizedKey, interface{}(value[0]), tag)
		case map[string]interface{}:
			for kk, vv := range value {
				b = b.AddMapWithTag(camelizedKey, kk, interface{}(vv), tag)
				break
			}
		case nil:
			// Note: Is this ok?
			b = b.AddInterfaceWithTag(camelizedKey, false, tag)
		default:
			return nil, fmt.Errorf("jsonData %#v has invalid typed key", m)
		}
	}

	ds := b.Build()
	intf, err := ds.DecodeMap(camelizedFieldMap)
	if err != nil {
		return nil, err
	}

	return intf, nil
}
