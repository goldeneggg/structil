package decoder

import (
	"github.com/goldeneggg/structil"
)

// NewJSON returns a concrete Decoder for JSON.
func NewJSON(data []byte) (*Decoder, error) {
	return newDecoder(data, typeJSON)
}

// JSONToI returns a decoded interface from JSON via DynamicStruct.
func JSONToI(data []byte) (interface{}, error) {
	d, err := NewJSON(data)
	if err != nil {
		return nil, err
	}

	return d.toDynamicStructI()
}

// JSONToGetter returns a structil.Getter with a decoded JSON via DynamicStruct.
// FIXME:
// この実装でも未知のJSON→Getter の変換が意図通り機能している事は確認できているが、
// DynamicStructのSetter対応と両睨みで対応方針を決める
func JSONToGetter(data []byte) (*structil.Getter, error) {
	intf, err := JSONToI(data)
	if err != nil {
		return nil, err
	}

	return structil.NewGetter(intf)
}
