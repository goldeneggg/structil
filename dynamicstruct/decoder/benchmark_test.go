package decoder_test

import (
	"testing"

	. "github.com/goldeneggg/structil/dynamicstruct/decoder"
)

func BenchmarkDynamicStructSingleJSON(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d, _ := NewJSON(singleJSON)
		_, _ = d.DynamicStruct(false, false)
	}
}

func BenchmarkDynamicStructArrayJSON(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d, _ := NewJSON(arrayJSON)
		_, _ = d.DynamicStruct(false, false)
	}
}

func BenchmarkDynamicStructSingleYAML(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d, _ := NewYAML(singleYAML)
		_, _ = d.DynamicStruct(false, false)
	}
}

func BenchmarkDynamicStructArrayYAML(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d, _ := NewYAML(arrayYAML)
		_, _ = d.DynamicStruct(false, false)
	}
}
