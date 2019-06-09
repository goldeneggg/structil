package structil_test

import (
	"reflect"
	"testing"

	. "github.com/goldeneggg/structil"
)

func BenchmarkNewGetter_Val(b *testing.B) {
	var g Getter
	var e error

	testStructVal := newTestStruct()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g, e = NewGetter(testStructVal)
		if e == nil {
			_ = g
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", e)
		}
	}
}

func BenchmarkNewGetter_Ptr(b *testing.B) {
	var g Getter
	var e error

	testStructPtr := newTestStructPtr()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g, e = NewGetter(testStructPtr)
		if e == nil {
			_ = g
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", e)
		}
	}
}

func BenchmarkGetterGetType_String(b *testing.B) {
	var t reflect.Type

	g, err := newTestGetter()
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t = g.GetType("String")
		_ = t
	}
}

func BenchmarkGetterGetValue_String(b *testing.B) {
	var v reflect.Value

	g, err := newTestGetter()
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v = g.GetValue("String")
		_ = v
	}
}

func BenchmarkGetterHas_String(b *testing.B) {
	var bl bool

	g, err := newTestGetter()
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bl = g.Has("String")
		_ = bl
	}
}

func BenchmarkGetterGet_String(b *testing.B) {
	var it interface{}

	g, err := newTestGetter()
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it = g.Get("String")
		_ = it
	}
}

func BenchmarkGetterEGet_String(b *testing.B) {
	var it interface{}

	g, err := newTestGetter()
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it, err = g.EGet("String")
		if err == nil {
			_ = it
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkGetterString(b *testing.B) {
	var str string

	g, err := newTestGetter()
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str = g.String("String")
		_ = str
	}
}

func BenchmarkGetterMapGet(b *testing.B) {
	var ia []interface{}

	g, err := newTestGetter()
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}
	fn := func(i int, g Getter) (interface{}, error) {
		return g.String("String") + ":" + g.String("String2"), nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ia, err = g.MapGet("TestStruct4PtrSlice", fn)
		if err == nil {
			_ = ia
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}
