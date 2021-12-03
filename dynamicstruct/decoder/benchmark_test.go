package decoder_test

import (
	"testing"

	. "github.com/goldeneggg/structil/dynamicstruct/decoder"
)

func BenchmarkFromJSON_singleJSON(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FromJSON(singleJSON)
	}
}

func BenchmarkDynamicStruct_singleJSON_nonNest_nonUseTag(b *testing.B) {
	d, _ := FromJSON(singleJSON)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(false, false)
	}
}

func BenchmarkDynamicStruct_singleJSON_nest_nonUseTag(b *testing.B) {
	d, _ := FromJSON(singleJSON)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(true, false)
	}
}

func BenchmarkDynamicStruct_singleJSON_nest_useTag(b *testing.B) {
	d, _ := FromJSON(singleJSON)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(true, true)
	}
}

func BenchmarkFromJSON_arrayJSON(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FromJSON(arrayJSON)
	}
}

func BenchmarkDynamicStruct_arrayJSON_nonNest_nonUseTag(b *testing.B) {
	d, _ := FromJSON(arrayJSON)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(false, false)
	}
}

func BenchmarkDynamicStruct_arrayJSON_nest_nonUseTag(b *testing.B) {
	d, _ := FromJSON(arrayJSON)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(true, false)
	}
}

func BenchmarkDynamicStruct_arrayJSON_nest_useTag(b *testing.B) {
	d, _ := FromJSON(arrayJSON)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(true, true)
	}
}

func BenchmarkFromYAML_singleYAML(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FromYAML(singleYAML)
	}
}

func BenchmarkDynamicStruct_singleYAML_nonNest_nonUseTag(b *testing.B) {
	d, _ := FromYAML(singleYAML)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(false, false)
	}
}

func BenchmarkDynamicStruct_singleYAML_nest_nonUseTag(b *testing.B) {
	d, _ := FromYAML(singleYAML)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(true, false)
	}
}

func BenchmarkDynamicStruct_singleYAML_nest_useTag(b *testing.B) {
	d, _ := FromYAML(singleYAML)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(true, false)
	}
}

func BenchmarkFromYAML_arrayYAML(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FromYAML(arrayYAML)
	}
}

func BenchmarkDynamicStruct_arrayYAML_nonNest_nonUseTag(b *testing.B) {
	d, _ := FromYAML(arrayYAML)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(false, false)
	}
}

func BenchmarkDynamicStruct_arrayYAML_nest_nonUseTag(b *testing.B) {
	d, _ := FromYAML(arrayYAML)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(true, false)
	}
}

func BenchmarkDynamicStruct_arrayYAML_nest_useTag(b *testing.B) {
	d, _ := FromYAML(arrayYAML)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(true, true)
	}
}

func BenchmarkJSONToGetter_singleJSON_nonNest(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = JSONToGetter(singleJSON, false)
	}
}

func BenchmarkJSONToGetter_singleJSON_nest(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = JSONToGetter(singleJSON, true)
	}
}

func BenchmarkJSONToGetter_arrayJSON_nonNest(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = JSONToGetter(arrayJSON, false)
	}
}

func BenchmarkJSONToGetter_arrayJSON_nest(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = JSONToGetter(arrayJSON, true)
	}
}

func BenchmarkYAMLToGetter_singleYAML_nonNest(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = YAMLToGetter(singleYAML, false)
	}
}

func BenchmarkYAMLToGetter_singleYAML_nest(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = YAMLToGetter(singleYAML, true)
	}
}

func BenchmarkYAMLToGetter_arrayYAML_nonNest(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = YAMLToGetter(arrayYAML, false)
	}
}

func BenchmarkYAMLToGetter_arrayYAML_nest(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = YAMLToGetter(arrayYAML, true)
	}
}
