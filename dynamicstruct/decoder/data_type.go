package decoder

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

type DataType interface {
	String() string
	Unmarshal([]byte, interface{}) error
}

// DefaultDataType is the type of original data format
type DefaultDataType int

const (
	// TypeJSON is the type sign of JSON
	TypeJSON DefaultDataType = iota

	// TypeYAML is the type sign of YAML
	TypeYAML
)

var formats = [...]string{
	TypeJSON: "json",
	TypeYAML: "yaml",
}

func (dt DefaultDataType) String() string {
	if dt >= 0 && int(dt) < len(formats) {
		return formats[dt]
	}
	return ""
}

func (dt DefaultDataType) Unmarshal(data []byte, ptr interface{}) (err error) {
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
