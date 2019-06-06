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

func BenchmarkGetRT_String(b *testing.B) {
	testStructPtr := newTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		b.Errorf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.GetRT("ExpString")
	}
}

func BenchmarkGetRV_String(b *testing.B) {
	testStructPtr := newTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		b.Errorf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.GetRV("ExpString")
	}
}

func BenchmarkHas_String(b *testing.B) {
	testStructPtr := newTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		b.Errorf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Has("ExpString")
	}
}

func BenchmarkGet_String(b *testing.B) {
	testStructPtr := newTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		b.Errorf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Get("ExpString")
	}
}

func BenchmarkGetString(b *testing.B) {
	testStructPtr := newTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		b.Errorf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.GetString("ExpString")
	}
}

func BenchmarkMapGet(b *testing.B) {
	testStructPtr := newTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		b.Errorf("NewGetter() occurs unexpected error: %v", err)
		return
	}
	fn := func(i int, g Getter) interface{} {
		return g.GetString("ExpString") + ":" + g.GetString("ExpString2")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.MapGet("TestStructPtrSlice", fn)
	}
}
