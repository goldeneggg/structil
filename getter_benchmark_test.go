package structil_test

import (
	"testing"

	. "github.com/goldeneggg/structil"
)

func BenchmarkNewGetter_Val(b *testing.B) {
	testStructVal := newTestStruct()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewGetter(testStructVal)
	}
}

func BenchmarkNewGetter_Ptr(b *testing.B) {
	testStructPtr := newTestStructPtr()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewGetter(testStructPtr)
	}
}

func BenchmarkGetterGetType_String(b *testing.B) {
	g, err := newTestGetter()
	if err != nil {
		b.Errorf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.GetType("String")
	}
}

func BenchmarkGetterGetValue_String(b *testing.B) {
	g, err := newTestGetter()
	if err != nil {
		b.Errorf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.GetValue("String")
	}
}

func BenchmarkGetterHas_String(b *testing.B) {
	g, err := newTestGetter()
	if err != nil {
		b.Errorf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Has("String")
	}
}

func BenchmarkGetterGet_String(b *testing.B) {
	g, err := newTestGetter()
	if err != nil {
		b.Errorf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Get("String")
	}
}

func BenchmarkGetterEGet_String(b *testing.B) {
	g, err := newTestGetter()
	if err != nil {
		b.Errorf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.EGet("String")
	}
}

func BenchmarkGetterString(b *testing.B) {
	g, err := newTestGetter()
	if err != nil {
		b.Errorf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.String("String")
	}
}

func BenchmarkGetterMapGet(b *testing.B) {
	g, err := newTestGetter()
	if err != nil {
		b.Errorf("NewGetter() occurs unexpected error: %v", err)
		return
	}
	fn := func(i int, g Getter) interface{} {
		return g.String("String") + ":" + g.String("String2")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.MapGet("TestStruct4PtrSlice", fn)
	}
}
