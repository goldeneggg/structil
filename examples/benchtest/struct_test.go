package benchtest

import (
	"reflect"
	"testing"
)

func BenchmarkReflect_New(b *testing.B) {
	var s *Person
	sv := reflect.TypeOf(Person{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sn := reflect.New(sv)
		s, _ = sn.Interface().(*Person)
	}
	_ = s
}

func BenchmarkDirect_New(b *testing.B) {
	var s *Person
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = new(Person)
	}
	_ = s
}

func BenchmarkReflect_Set(b *testing.B) {
	var s *Person
	sv := reflect.TypeOf(Person{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sn := reflect.New(sv)
		s = sn.Interface().(*Person)
		s.Name = "Jerry"
		s.Age = 18
	}
}

func BenchmarkReflect_SetFieldByName(b *testing.B) {
	sv := reflect.TypeOf(Person{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sn := reflect.New(sv).Elem()
		sn.FieldByName("Name").SetString("Jerry")
		sn.FieldByName("Age").SetInt(18)
	}
}

func BenchmarkReflect_SetFieldByIndex(b *testing.B) {
	sv := reflect.TypeOf(Person{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sn := reflect.New(sv).Elem()
		sn.Field(0).SetString("Jerry")
		sn.Field(1).SetInt(18)
	}
}

func BenchmarkDirect_Set(b *testing.B) {
	var s *Person
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = new(Person)
		s.Name = "Jerry"
		s.Age = 18
	}
}
