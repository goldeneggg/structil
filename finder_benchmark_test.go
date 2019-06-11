package structil_test

import (
	"testing"

	. "github.com/goldeneggg/structil"
)

func BenchmarkNewFinder_Val(b *testing.B) {
	var f Finder
	var e error

	testStructVal := newTestStruct()
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
	var f Finder
	var e error

	testStructPtr := newTestStructPtr()
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

	f, err := NewFinder(newTestStructPtr())
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

	f, err := NewFinder(newTestStructPtr())
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

	f, err := NewFinder(newTestStructPtr())
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, err = f.Into("TestStruct2").Find("String").ToMap()
		if err == nil {
			_ = m
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkToMap_1Struct_1Find_2Pair(b *testing.B) {
	var m map[string]interface{}

	f, err := NewFinder(newTestStructPtr())
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, err = f.Into("TestStruct2").Find("String").Into("TestStruct2Ptr").Find("String").ToMap()
		if err == nil {
			_ = m
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkToMap_2Struct_1Find(b *testing.B) {
	var m map[string]interface{}

	f, err := NewFinder(newTestStructPtr())
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, err = f.Into("TestStruct2", "TestStruct3").Find("String").ToMap()
		if err == nil {
			_ = m
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkToMap_2Struct_2Find(b *testing.B) {
	var m map[string]interface{}

	f, err := NewFinder(newTestStructPtr())
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, err = f.Into("TestStruct2", "TestStruct3").Find("String", "Int").ToMap()
		if err == nil {
			_ = m
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkFromKeys_Yml(b *testing.B) {
	f, err := NewFinder(newTestStructPtr())
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fks, err := NewFinderKeysFromConf("examples/finder_from_conf", "ex_test1_yml")
		if err == nil {
			_ = f.FromKeys(fks)
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkFromKeys_Json(b *testing.B) {
	f, err := NewFinder(newTestStructPtr())
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fks, err := NewFinderKeysFromConf("examples/finder_from_conf", "ex_test1_json")
		if err == nil {
			_ = f.FromKeys(fks)
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}
