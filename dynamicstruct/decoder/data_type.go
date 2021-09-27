package decoder

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

// DataType is implemented by any value that has String and Unmershal method,
type DataType interface {
	String() string
	Unmarshal([]byte, interface{}) error
	Marshal(interface{}) ([]byte, error)
	IntfToStringMap(interface{}) (map[string]interface{}, error)
}

// DefaultDataType is the type of original data format
type defaultDataType int

const (
	// TypeJSON is the type sign of JSON
	TypeJSON defaultDataType = iota

	// TypeYAML is the type sign of YAML
	TypeYAML

	// FIXME: futures as follows

	// TypeXML is the type sign of XML
	// TypeXML

	// TypeTOML is the type sign of TOML
	// TypeTOML

	// TypeCSV is the type sign of CSV
	// TypeCSV
)

var formats = [...]string{
	TypeJSON: "json",
	TypeYAML: "yaml",
}

func (ddt defaultDataType) String() string {
	if ddt >= 0 && int(ddt) < len(formats) {
		return formats[ddt]
	}
	return ""
}

func (ddt defaultDataType) Unmarshal(data []byte, ptr interface{}) (err error) {
	switch ddt {
	case TypeJSON:
		err = json.Unmarshal(data, ptr)
	case TypeYAML:
		// Note: The gopkg.in/yaml.v2 package returns an unmarshaled interface as "map[interface{}]interface{}" type.
		err = yaml.Unmarshal(data, ptr)
	default:
		err = fmt.Errorf("invalid datatype for Unmarshal: %v", ddt)
	}

	return
}

// FIXME: add tests and examples
func (ddt defaultDataType) Marshal(v interface{}) (data []byte, err error) {
	switch ddt {
	case TypeJSON:
		data, err = json.Marshal(v)
	case TypeYAML:
		data, err = yaml.Marshal(v)
	default:
		err = fmt.Errorf("invalid datatype for Marshal: %v", ddt)
	}

	return
}

func (ddt defaultDataType) IntfToStringMap(v interface{}) (map[string]interface{}, error) {
	ms := make(map[string]interface{})

	if v == nil {
		return ms, nil
	}

	var ok bool

	switch ddt {
	case TypeJSON:
		ms, ok = v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expect to be returned map[string]interface{} but not. v = %#v", v)
		}
		return ms, nil
	case TypeYAML:
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
		return nil, fmt.Errorf("invalid datatype for Marshal: %v", ddt)
	}
}
