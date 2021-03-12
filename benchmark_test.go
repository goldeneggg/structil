package structil_test

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	. "github.com/goldeneggg/structil"
)

func BenchmarkNewGetter_Val(b *testing.B) {
	var g *Getter
	var e error

	testStructVal := newGetterTestStruct() // See: getter_test.go
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
	var g *Getter
	var e error

	testStructPtr := newGetterTestStructPtr() // See: getter_test.go
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

	g, err := newTestGetter() // See: getter_test.go
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t, _ = g.GetType("String")
		_ = t
	}
}

func BenchmarkGetterGetValue_String(b *testing.B) {
	var v reflect.Value

	g, err := newTestGetter() // See: getter_test.go
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v, _ = g.GetValue("String")
		_ = v
	}
}

func BenchmarkGetterHas_String(b *testing.B) {
	var bl bool

	g, err := newTestGetter() // See: getter_test.go
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

	g, err := newTestGetter() // See: getter_test.go
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it, _ = g.Get("String")
		_ = it
	}
}

func BenchmarkGetterString(b *testing.B) {
	var str string

	g, err := newTestGetter() // See: getter_test.go
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str, _ = g.String("String")
		_ = str
	}
}

func BenchmarkGetterUintptr(b *testing.B) {
	var up uintptr

	g, err := newTestGetter() // See: getter_test.go
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		up, _ = g.Uintptr("Uintptr")
		_ = up
	}
}

func BenchmarkGetterUnsafePointer(b *testing.B) {
	var up unsafe.Pointer

	g, err := newTestGetter() // See: getter_test.go
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		up, _ = g.UnsafePointer("Unsafeptr")
		_ = up
	}
}

func BenchmarkGetterIsStruct(b *testing.B) {
	var is bool

	g, err := newTestGetter() // See: getter_test.go
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		is = g.IsStruct("GetterTestStruct2")
		_ = is
	}
}

func BenchmarkGetterIsSlice_Bytes(b *testing.B) {
	var is bool

	g, err := newTestGetter() // See: getter_test.go
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		is = g.IsSlice("Bytes")
		_ = is
	}
}

func BenchmarkGetterIsSlice_StructSlice(b *testing.B) {
	var is bool

	g, err := newTestGetter() // See: getter_test.go
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		is = g.IsSlice("GetterTestStruct4Slice")
		_ = is
	}
}

func BenchmarkGetterIsSlice_StructPtrSlice(b *testing.B) {
	var is bool

	g, err := newTestGetter() // See: getter_test.go
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		is = g.IsSlice("GetterTestStruct4PtrSlice")
		_ = is
	}
}

func BenchmarkGetterMapGet(b *testing.B) {
	var ia []interface{}

	g, err := newTestGetter() // See: getter_test.go
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}
	fn := func(i int, g *Getter) (interface{}, error) {
		str, _ := g.String("String")
		str2, _ := g.String("String")
		return fmt.Sprintf("%s:%s", str, str2), nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ia, err = g.MapGet("GetterTestStruct4PtrSlice", fn)
		if err == nil {
			_ = ia
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkNewFinder_Val(b *testing.B) {
	var f *Finder
	var e error

	testStructVal := newFinderTestStruct() // See: getter_test.go
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, e = NewFinder(testStructVal)
		if e == nil {
			_ = f
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", e)
		}
	}
}

func BenchmarkNewFinder_Ptr(b *testing.B) {
	var f *Finder
	var e error

	testStructPtr := newFinderTestStructPtr() // See: getter_test.go
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, e = NewFinder(testStructPtr)
		if e == nil {
			_ = f
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", e)
		}
	}
}

func BenchmarkToMap_1FindOnly(b *testing.B) {
	var m map[string]interface{}

	f, err := NewFinder(newFinderTestStructPtr()) // See: finder_test.go
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, err = f.Find("String").ToMap()
		if err == nil {
			_ = m
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkToMap_2FindOnly(b *testing.B) {
	var m map[string]interface{}

	f, err := NewFinder(newFinderTestStructPtr()) // See: finder_test.go
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, err = f.Find("String", "Int64").ToMap()
		if err == nil {
			_ = m
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkToMap_1Struct_1Find(b *testing.B) {
	var m map[string]interface{}

	f, err := NewFinder(newFinderTestStructPtr()) // See: finder_test.go
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, err = f.Into("FinderTestStruct2").Find("String").ToMap()
		if err == nil {
			_ = m
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkToMap_1Struct_1Find_2Pair(b *testing.B) {
	var m map[string]interface{}

	f, err := NewFinder(newFinderTestStructPtr()) // See: finder_test.go
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, err = f.Into("FinderTestStruct2").Find("String").Into("FinderTestStruct2Ptr").Find("String").ToMap()
		if err == nil {
			_ = m
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkToMap_2Struct_1Find(b *testing.B) {
	var m map[string]interface{}

	f, err := NewFinder(newFinderTestStructPtr()) // See: finder_test.go
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, err = f.Into("FinderTestStruct2", "FinderTestStruct3").Find("String").ToMap()
		if err == nil {
			_ = m
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkToMap_2Struct_2Find(b *testing.B) {
	var m map[string]interface{}

	f, err := NewFinder(newFinderTestStructPtr()) // See: finder_test.go
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, err = f.Into("FinderTestStruct2", "FinderTestStruct3").Find("String", "Int").ToMap()
		if err == nil {
			_ = m
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkNewFinderKeys_yml(b *testing.B) {
	f, err := NewFinder(newFinderTestStructPtr()) // See: finder_test.go
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fks, err := NewFinderKeys("testdata/finder_from_conf", "ex_test1_yml")
		if err == nil {
			_ = f.FromKeys(fks)
			f.Reset()
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkNewFinderKeys_json(b *testing.B) {
	f, err := NewFinder(newFinderTestStructPtr()) // See: finder_test.go
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fks, err := NewFinderKeys("testdata/finder_from_conf", "ex_test1_json")
		if err == nil {
			_ = f.FromKeys(fks)
			f.Reset()
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}
