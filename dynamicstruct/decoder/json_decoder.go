package decoder

import (
	"encoding/json"
)

// JSONDecoder is the decoded object from unknown format JSON/YAML/and others.
type JSONDecoder struct {
}

// NewJSONDecoder returns a concrete Decoder for JSON
func NewJSONDecoder() *JSONDecoder {
	return &JSONDecoder{}
}

// Decode decodes JSON data to interface via DynamicStruct.DecodeMap.
// data argument must be a byte array data of valid JSON.
func (jgd *JSONDecoder) Decode(data []byte) (*DecodedResult, error) {
	// FIXME:
	// want to add json validation. but is json.Valid(data) too slow?
	// See: https://stackoverflow.com/questions/22128282/how-to-check-string-is-in-json-format

	var ui interface{}
	if err := json.Unmarshal(data, &ui); err != nil {
		return nil, err
	}

	dr, err := decode(ui, nil)
	if err != nil {
		return nil, err
	}

	return dr, nil
}
