package decoder_test

import (
	"testing"

	. "github.com/goldeneggg/structil/dynamicstruct/decoder"
)

func BenchmarkDynamicStructSingleJSON(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d, _ := New(singleJSON, TypeJSON)
		_, _ = d.DynamicStruct(false, false)
	}
}

func BenchmarkDynamicStructArrayJSON(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d, _ := New(arrayJSON, TypeJSON)
		_, _ = d.DynamicStruct(false, false)
	}
}

func BenchmarkDynamicStructSingleYAML(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d, _ := New(singleYAML, TypeYAML)
		_, _ = d.DynamicStruct(false, false)
	}
}

func BenchmarkDynamicStructArrayYAML(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d, _ := New(arrayYAML, TypeYAML)
		_, _ = d.DynamicStruct(false, false)
	}
}
