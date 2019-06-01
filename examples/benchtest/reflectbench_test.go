package benchtest

import (
	"reflect"
	"testing"
)

type Person struct {
	Name    string
	Age     int
	Hobbies []*Hobby
}

type Hobby struct {
	Name    string
	Subject string
}

// 3x slower than KindOnce...
func Benchmark_KindEach(b *testing.B) {
	frv := getFrv("Hobbies")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = frv.Kind() == reflect.Array || frv.Kind() == reflect.Slice
	}
}

func Benchmark_KindOnce(b *testing.B) {
	frv := getFrv("Hobbies")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var k reflect.Kind
		k = frv.Kind()
		_ = k == reflect.Array || k == reflect.Slice
	}
}

func getFrv(name string) reflect.Value {
	s := &Person{
		Hobbies: []*Hobby{
			{
				Name:    "X",
				Subject: "a",
			},
			{
				Name:    "Y",
				Subject: "B",
			},
		},
	}
	sv := reflect.Indirect(reflect.ValueOf(s))
	return sv.FieldByName(name)
}
