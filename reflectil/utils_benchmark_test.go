package reflectil_test

import (
	"reflect"
	"testing"

	. "github.com/goldeneggg/structil/reflectil"
)

func BenchmarkToI_String(b *testing.B) {
	benchmarkToI(reflect.ValueOf("str"), b)
}

func BenchmarkToI_Int(b *testing.B) {
	benchmarkToI(reflect.ValueOf(123), b)
}

func BenchmarkToI_StructPtr(b *testing.B) {
	benchmarkToI(reflect.ValueOf(testTstrPtr), b)
}

func BenchmarkToI_Map(b *testing.B) {
	benchmarkToI(reflect.ValueOf(testMap), b)
}

func BenchmarkToI_Func(b *testing.B) {
	benchmarkToI(reflect.ValueOf(testFunc), b)
}

func BenchmarkToI_Chan(b *testing.B) {
	benchmarkToI(reflect.ValueOf(testChan), b)
}

func BenchmarkToI_Nil(b *testing.B) {
	benchmarkToI(reflect.ValueOf(nil), b)
}

func benchmarkToI(v reflect.Value, b *testing.B) {
	var intf interface{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		intf = ToI(v)
		_ = intf
	}
}
