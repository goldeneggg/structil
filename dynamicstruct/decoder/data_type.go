package decoder

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

// DataType is the type of original data format
type DataType int

const (
	// TypeJSON is the type sign of JSON
	TypeJSON DataType = iota

	// TypeYAML is the type sign of YAML
	TypeYAML
)

var formats = [...]string{
	TypeJSON: "json",
	TypeYAML: "yaml",
}

func (dt DataType) String() string {
	if dt >= 0 && int(dt) < len(formats) {
		return formats[dt]
	}
	return ""
}

func (dt DataType) Unmarshal(data []byte, ptr interface{}) (err error) {
	switch dt {
	case TypeJSON:
		err = json.Unmarshal(data, ptr)
	case TypeYAML:
		err = yaml.Unmarshal(data, ptr)
	default:
		err = fmt.Errorf("invalid datatype: %v", dt)
	}

	return
}
