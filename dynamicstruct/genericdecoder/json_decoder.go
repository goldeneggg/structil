package genericdecoder

import (
	"encoding/json"
)

// JSONGenericDecoder is the decoded object from unknown format JSON/YAML/and others.
type JSONGenericDecoder struct {
}

// NewFromJSON returns a concrete GenericDecoder for JSON
// data argument must be a byte array data of valid JSON.
// TODO: YAML and TOML support
func NewFromJSON(jsonData []byte) (GenericDecoder, error) {
	return &JSONGenericDecoder{}, nil
}

// Decode decodes JSON data to interface via DynamicStruct.DecodeMap.
func (jgd *JSONGenericDecoder) Decode(data []byte) (*DecodedResult, error) {
	// FIXME:
	// want to add json validation. but is json.Valid(data) too slow?
	// See: https://stackoverflow.com/questions/22128282/how-to-check-string-is-in-json-format

	var intf interface{}
	if err := json.Unmarshal(data, &intf); err != nil {
		return nil, err
	}

	dr, err := decode(intf, nil)
	if err != nil {
		return nil, err
	}

	return dr, nil
}
