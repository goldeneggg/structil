package decoder

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"gopkg.in/yaml.v3"
)

// dataType is the type of original data format
// This type provides an unified interface for marshal and unmarshal functions per data formats.
type dataType int

const (
	typeJSON dataType = iota
	typeYAML
	typeHCL

	// FIXME: futures as follows
	// TypeXML
	// TypeTOML
	// TypeCSV

	end // end of iota
)

var formats = [...]string{
	typeJSON: "json",
	typeYAML: "yaml",
	typeHCL:  "hcl",
}

func (dt dataType) string() string {
	if dt >= 0 && int(dt) < len(formats) {
		return formats[dt]
	}
	return ""
}

func (dt dataType) unmarshal(data []byte, iptr interface{}) error {
	var err error

	switch dt {
	case typeJSON:
		err = json.Unmarshal(data, iptr)
	case typeYAML:
		err = yaml.Unmarshal(data, iptr)
	case typeHCL:
		var i map[string]interface{}
		iptr = &i
		err = hclsimple.Decode("example.hcl", data, nil, iptr)
	default:
		err = fmt.Errorf("invalid datatype for Unmarshal: %v", dt)
	}

	return err
}

// TODO: add tests and examples
// func (dt dataType) marshal(v interface{}) (data []byte, err error) {
func (dt dataType) marshal(m map[string]interface{}) (data []byte, err error) {
	switch dt {
	case typeJSON:
		data, err = json.Marshal(m)
	case typeYAML:
		data, err = yaml.Marshal(m)
	default:
		err = fmt.Errorf("invalid datatype for Marshal: %v", dt)
	}

	return
}
