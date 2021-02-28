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
