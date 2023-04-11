package decoder

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// dataType is the type of original data format
// This type provides an unified interface for marshal and unmarshal functions per data formats.
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

	end // end of iota
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

func (dt dataType) unmarshal(data []byte) (any, error) {
	var intf any
	err := dt.unmarshalWithIPtr(data, &intf)
	return intf, err
}

func (dt dataType) unmarshalWithIPtr(data []byte, iptr any) error {
	var err error

	switch dt {
	case typeJSON:
		// Note: iptr should be "map[string]any"
		err = json.Unmarshal(data, iptr)
	case typeYAML:
		// Note: iptr should be "map[any]any" using gopkg.in/yaml.v2 package
		err = yaml.Unmarshal(data, iptr)
	default:
		err = fmt.Errorf("invalid datatype for Unmarshal: %v", dt)
	}

	return err
}

// TODO: add tests and examples
// func (dt dataType) marshal(v any) (data []byte, err error) {
func (dt dataType) marshal(m map[string]any) (data []byte, err error) {
	switch dt {
	case typeJSON:
		// Note: v is expected to be "map[string]any"
		data, err = json.Marshal(m)
	case typeYAML:
		// Note: v is expected to be converted from "map[any]any" to "map[string]any"
		data, err = yaml.Marshal(m)
	default:
		err = fmt.Errorf("invalid datatype for Marshal: %v", dt)
	}

	return
}
