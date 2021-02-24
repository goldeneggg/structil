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
type defaultDataType int

const (
	// TypeJSON is the type sign of JSON
	TypeJSON defaultDataType = iota

	// TypeYAML is the type sign of YAML
	TypeYAML
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
		err = yaml.Unmarshal(data, ptr)
	default:
		err = fmt.Errorf("invalid datatype: %v", ddt)
	}

	return
}
