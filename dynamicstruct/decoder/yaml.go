package decoder

import (
	"encoding/json"

	"github.com/goldeneggg/structil"
)

// NewYAML returns a concrete Decoder for YAML.
func NewYAML(data []byte) (*Decoder, error) {
	return New(data, TypeYAML)
}

// YAMLToI returns a decoded interface from YAML via DynamicStruct.
// Note: The gopkg.in/yaml.v2 package returns an unmarshaled interface as "map[interface{}]interface{}" type.
func YAMLToI(data []byte) (interface{}, error) {
	d, err := NewYAML(data)
	if err != nil {
		return nil, err
	}

	return d.toDynamicStructI()
}

// YAMLToGetter returns a structil.Getter with a decoded YAML via DynamicStruct.
// FIXME:
// この実装でも未知のYAML→Getter の変換が意図通り機能している事は確認できているが、
// DynamicStructのSetter対応と両睨みで対応方針を決める
func YAMLToGetter(data []byte) (*structil.Getter, error) {
	intf, err := YAMLToI(data)
	if err != nil {
		return nil, err
	}

	// Note:
	// YAMLToI returns a map[interface{}]interface{} so via json.Marshal for generating []byte data.
	// FIXME: json.Marshal does not support "map[interface{}]interface{}", and is expected "map[string]interface{}"
	jsonD, err := json.Marshal(intf)
	if err != nil {
		return nil, err
	}

	return JSONToGetter(jsonD)
}
