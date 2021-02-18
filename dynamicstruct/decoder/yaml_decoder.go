package decoder

import (
	"gopkg.in/yaml.v2"
)

// YAMLDecoder is the decoded object from unknown format YAML/YAML/and others.
type YAMLDecoder struct {
}

// NewYAMLDecoder returns a concrete Decoder for YAML
func NewYAMLDecoder() *YAMLDecoder {
	return &YAMLDecoder{}
}

// Decode decodes YAML data to interface via DynamicStruct.DecodeMap.
// data argument must be a byte array data of valid YAML.
func (yd *YAMLDecoder) Decode(data []byte) (*DecodedResult, error) {
	var ui interface{}
	if err := yaml.Unmarshal(data, &ui); err != nil {
		return nil, err
	}

	dr, err := decode(ui, nil)
	if err != nil {
		return nil, err
	}

	return dr, nil
}
