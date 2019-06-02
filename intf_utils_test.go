package structil

import (
	"io"
	"os"
	"reflect"
	"testing"
)

func TestIElemOf(t *testing.T) {
	istr := "str"
	iint := 123
	var iwriter io.Writer
	iwriter = os.Stdout

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
			want: reflect.Indirect(reflect.ValueOf(istr)),
		},
		{
			name: "i is int",
			args: args{i: iint},
			want: reflect.Indirect(reflect.ValueOf(iint)),
		},
		{
			name: "i is io.Writer(os.Stdout)",
			args: args{i: iwriter},
			want: reflect.Indirect(reflect.ValueOf(iwriter)),
		},
		{
			name: "i is nil",
			args: args{i: nil},
			want: reflect.Indirect(reflect.ValueOf(nil)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := IElemOf(tt.args.i); !reflect.DeepEqual(got, tt.want) {
			if tt.args.i == nil {
				if got := IElemOf(tt.args.i); got.Kind() != reflect.Invalid {
					t.Errorf("IElemOf() = %v is not invalid", got)
				}
			} else {
				if got := IElemOf(tt.args.i); got.Interface() != tt.want.Interface() {
					t.Errorf("IElemOf() = %v, wank %v", got, tt.want)
				}
			}
		})
	}
}
