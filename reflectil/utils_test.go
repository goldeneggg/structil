package reflectil_test

import (
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
