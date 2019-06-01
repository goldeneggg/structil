package structil

import (
	"reflect"
	"testing"
)

func TestIElemOf(t *testing.T) {
	istr := "str"
	iint := 123

	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want reflect.Value
	}{
		{
			name: "i is string",
			args: args{i: istr},
			want: reflect.ValueOf(&istr).Elem(),
		},
		{
			name: "i is int",
			args: args{i: iint},
			want: reflect.ValueOf(&iint).Elem(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := IElemOf(tt.args.i); !reflect.DeepEqual(got, tt.want) {
			if got := IElemOf(tt.args.i); got.Interface() != tt.want.Interface() {
				t.Errorf("IElemOf() = %v, wank %v", got, tt.want)
			}
		})
	}
}
