package decoder

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

// DataType is implemented by any value that has String and Unmershal method,
/*
type DataType interface {
	String() string
	Unmarshal([]byte, interface{}) error
	Marshal(interface{}) ([]byte, error)
	IntfToStringMap(interface{}) (map[string]interface{}, error)
}
*/

// DefaultDataType is the type of original data format
type dataType int

const (
	// TypeJSON is the type sign of JSON
	typeJSON dataType = iota

	// TypeYAML is the type sign of YAML
	typeYAML

	// FIXME: futures as follows

	// TypeXML is the type sign of XML
	// TypeXML

	// TypeTOML is the type sign of TOML
	// TypeTOML

	// TypeCSV is the type sign of CSV
	// TypeCSV
)

var formats = [...]string{
	typeJSON: "json",
	typeYAML: "yaml",
}

func (dt dataType) string() string {
	if dt >= 0 && int(dt) < len(formats) {
		return formats[dt]
	}
	return ""
}

func (dt dataType) unmarshal(data []byte, ptr interface{}) (err error) {
	switch dt {
	case typeJSON:
		err = json.Unmarshal(data, ptr)
	case typeYAML:
		// Note: The gopkg.in/yaml.v2 package returns an unmarshaled interface as "map[interface{}]interface{}" type.
		err = yaml.Unmarshal(data, ptr)
	default:
		err = fmt.Errorf("invalid datatype for Unmarshal: %v", dt)
	}

	return
}

// FIXME: add tests and examples
func (dt dataType) marshal(v interface{}) (data []byte, err error) {
	switch dt {
	case typeJSON:
		data, err = json.Marshal(v)
	case typeYAML:
		data, err = yaml.Marshal(v)
	default:
		err = fmt.Errorf("invalid datatype for Marshal: %v", dt)
	}

	return
}

func (dt dataType) intfToStringMap(v interface{}) (map[string]interface{}, error) {
	ms := make(map[string]interface{})

	if v == nil {
		return ms, nil
	}

	var ok bool

	switch dt {
	case typeJSON:
		ms, ok = v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expect to be returned map[string]interface{} but not. v = %#v", v)
		}
		return ms, nil
	case typeYAML:
		mi, ok := v.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("expect to be returned map[interface{}]interface{} but not. v = %#v", v)
		}

		for kmi, vmi := range mi {
			kms, ok := kmi.(string)
			if !ok {
				return nil, fmt.Errorf("expect to be returned string, but %#v", kms)
			}
			ms[kms] = vmi
		}
		return ms, nil
	default:
		return nil, fmt.Errorf("invalid datatype for Marshal: %v", dt)
	}
}
