package util_test

import (
	"errors"
	"reflect"
	"testing"

	. "github.com/goldeneggg/structil/util"
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

func BenchmarkElemTypeOf_String(b *testing.B) {
	benchmarkElemTypeOf("str", b)
}

func BenchmarkElemTypeOf_Int(b *testing.B) {
	benchmarkElemTypeOf(123, b)
}

func BenchmarkElemTypeOf_StructPtr(b *testing.B) {
	benchmarkElemTypeOf(testTstrPtr, b)
}

func BenchmarkElemTypeOf_Map(b *testing.B) {
	benchmarkElemTypeOf(testMap, b)
}

func BenchmarkElemTypeOf_Func(b *testing.B) {
	benchmarkElemTypeOf(testFunc, b)
}

func BenchmarkElemTypeOf_Chan(b *testing.B) {
	benchmarkElemTypeOf(testChan, b)
}

func BenchmarkElemTypeOf_Error(b *testing.B) {
	benchmarkElemTypeOf(errors.New("testerror"), b)
}

func BenchmarkElemTypeOf_Nil(b *testing.B) {
	benchmarkElemTypeOf(nil, b)
}

func benchmarkElemTypeOf(i interface{}, b *testing.B) {
	var typ reflect.Type

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		typ = ElemTypeOf(i)
		_ = typ
	}
}

func BenchmarkRecoverToError(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		recoverToErrorDemo()
	}
}

func recoverToErrorDemo() {
	defer func() {
		if r := recover(); r != nil {
			_ = RecoverToError(r)
		}
	}()
	panic("panic for benchmark")
}
