package structil_test

import (
	"testing"

	. "github.com/goldeneggg/structil"
)

func BenchmarkNewFinder_Val(b *testing.B) {
	testStructVal := newTestStruct()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewFinder(testStructVal)
	}
}

func BenchmarkNewFinder_Ptr(b *testing.B) {
	testStructPtr := newTestStructPtr()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewFinder(testStructPtr)
	}
}

func BenchmarkToMap_1FindOnly(b *testing.B) {
	testStructPtr := newTestStructPtr()
	f, err := NewFinder(testStructPtr)
	if err != nil {
		b.Errorf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = f.Find("String").ToMap()
		if err != nil {
			b.Errorf("ToMap() occurs unexpected error: %v", err)
			return
		}
	}
}

func BenchmarkToMap_2FindOnly(b *testing.B) {
	testStructPtr := newTestStructPtr()
	f, err := NewFinder(testStructPtr)
	if err != nil {
		b.Errorf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = f.Find("String", "Int64").ToMap()
		if err != nil {
			b.Errorf("ToMap() occurs unexpected error: %v", err)
			return
		}
	}
}

func BenchmarkToMap_1Struct_1Find(b *testing.B) {
	testStructPtr := newTestStructPtr()
	f, err := NewFinder(testStructPtr)
	if err != nil {
		b.Errorf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = f.Struct("TestStruct2").Find("String").ToMap()
		if err != nil {
			b.Errorf("ToMap() occurs unexpected error: %v", err)
			return
		}
	}
}

func BenchmarkToMap_1Struct_1Find_2Pair(b *testing.B) {
	testStructPtr := newTestStructPtr()
	f, err := NewFinder(testStructPtr)
	if err != nil {
		b.Errorf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = f.Struct("TestStruct2").Find("String").Struct("TestStruct2Ptr").Find("String").ToMap()
		if err != nil {
			b.Errorf("ToMap() occurs unexpected error: %v", err)
			return
		}
	}
}

func BenchmarkToMap_2Struct_1Find(b *testing.B) {
	testStructPtr := newTestStructPtr()
	f, err := NewFinder(testStructPtr)
	if err != nil {
		b.Errorf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = f.Struct("TestStruct2", "TestStruct3").Find("String").ToMap()
		if err != nil {
			b.Errorf("ToMap() occurs unexpected error: %v", err)
			return
		}
	}
}

func BenchmarkToMap_2Struct_2Find(b *testing.B) {
	testStructPtr := newTestStructPtr()
	f, err := NewFinder(testStructPtr)
	if err != nil {
		b.Errorf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = f.Struct("TestStruct2", "TestStruct3").Find("String", "Int").ToMap()
		if err != nil {
			b.Errorf("ToMap() occurs unexpected error: %v", err)
			return
		}
	}
}
