package reflectil_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/goldeneggg/structil/reflectil"
)

type tStruct struct {
	ID   int
	Name string
}

var (
	testTstrPtr = &tStruct{10, "Name10"}
	testMap     = map[string]interface{}{"key1": "value1", "key2": 2}
	testFunc    = func(s string) interface{} { return s + "-func" }
	testChan    = make(chan int)
)

func TestToI(t *testing.T) {
	t.Parallel()

	type args struct {
		i reflect.Value
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "string",
			args: args{i: reflect.ValueOf("str")},
			want: "str",
		},
		{
			name: "int",
			args: args{i: reflect.ValueOf(123)},
			want: 123,
		},
		{
			name: "struct",
			args: args{i: reflect.ValueOf(tStruct{10, "Name10"})},
			want: tStruct{10, "Name10"},
		},
		{
			name: "struct ptr",
			args: args{i: reflect.ValueOf(testTstrPtr)},
			want: testTstrPtr,
		},
		{
			name: "struct slice",
			args: args{i: reflect.ValueOf([]tStruct{{10, "Name10"}, {20, "Name20"}})},
			want: []tStruct{{10, "Name10"}, {20, "Name20"}},
		},
		{
			name: "struct ptr slice",
			args: args{i: reflect.ValueOf([]*tStruct{&tStruct{30, "Name30"}, &tStruct{40, "Name40"}})},
			want: []*tStruct{&tStruct{30, "Name30"}, &tStruct{40, "Name40"}},
		},
		{
			name: "map",
			args: args{i: reflect.ValueOf(testMap)},
			want: testMap,
		},
		{
			name: "func",
			args: args{i: reflect.ValueOf(testFunc)},
			want: testFunc,
		},
		{
			name: "chan",
			args: args{i: reflect.ValueOf(testChan)},
			want: testChan,
		},
		{
			name: "nil",
			args: args{i: reflect.ValueOf(nil)},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToI(tt.args.i)

			if tt.name == "func" {
				// Note: cmp.Diff does not support comparing func and func
				gp := reflect.ValueOf(got).Pointer()
				wp := reflect.ValueOf(tt.want).Pointer()
				if gp != wp {
					t.Errorf("unexpected mismatch func type: gp: %v, wp: %v", gp, wp)
				}
			} else if d := cmp.Diff(got, tt.want); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestElemTypeOf(t *testing.T) {
	t.Parallel()

	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "string",
			args: args{i: "str"},
			want: reflect.TypeOf("str").String(),
		},
		{
			name: "int",
			args: args{i: 123},
			want: reflect.TypeOf(123).String(),
		},
		{
			name: "struct",
			args: args{i: tStruct{10, "Name10"}},
			want: reflect.TypeOf(tStruct{10, "Name10"}).String(),
		},
		{
			name: "struct ptr",
			args: args{i: testTstrPtr},
			want: reflect.TypeOf(testTstrPtr).Elem().String(),
		},
		{
			name: "struct slice",
			args: args{i: []tStruct{{10, "Name10"}, {20, "Name20"}}},
			want: reflect.TypeOf([]tStruct{{10, "Name10"}, {20, "Name20"}}).Elem().String(),
		},
		{
			name: "map",
			args: args{i: testMap},
			want: reflect.TypeOf(testMap).String(),
		},
		{
			name: "func",
			args: args{i: testFunc},
			want: reflect.TypeOf(testFunc).String(),
		},
		{
			name: "chan",
			args: args{i: testChan},
			want: reflect.TypeOf(testChan).Elem().String(),
		},
		{
			name: "error",
			args: args{i: errors.New("testerror")},
			want: reflect.TypeOf(errors.New("testerror")).Elem().String(),
		},
		{
			name: "(*error)(nil)",
			args: args{i: (*error)(nil)},
			want: reflect.TypeOf((*error)(nil)).Elem().String(),
		},
		{
			name: "nil",
			args: args{i: nil},
			want: "nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ElemTypeOf(tt.args.i)

			if got == nil {
				if tt.want != "nil" {
					t.Errorf("expected nil but got %+v is not nil", got)
				}
				return
			}

			gs := got.String()
			if gs != tt.want {
				t.Errorf("unexpected type: got: %s, want: %s", gs, tt.want)
				return
			}
		})
	}
}

// benchmark tests

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
