package structil_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/goldeneggg/structil"
)

func TestNewGetter(t *testing.T) {
	t.Parallel()

	testStructVal := newTestStruct()
	testStructPtr := newTestStructPtr()

	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "valid struct",
			args:    args{i: testStructVal},
			wantErr: false,
		},
		{
			name:    "valid struct ptr",
			args:    args{i: testStructPtr},
			wantErr: false,
		},
		{
			name:    "invalid (nil)",
			args:    args{i: nil},
			wantErr: true,
		},
		{
			name:    "invalid (struct nil)",
			args:    args{i: (*TestStruct)(nil)},
			wantErr: true,
		},
		{
			name:    "invalid (string)",
			args:    args{i: "abc"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGetter(tt.args.i)

			if err == nil {
				if _, ok := got.(Getter); !ok {
					t.Errorf("NewGetter() want Getter but got %+v", got)
				}
			} else if !tt.wantErr {
				t.Errorf("NewGetter() unexpected error [%v] occured. wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHas(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid name and it's type is string",
			args: args{name: "String"},
			want: true,
		},
		{
			name: "valid name and it's type is string (2nd)",
			args: args{name: "String"},
			want: true,
		},
		{
			name: "invalid name",
			args: args{name: "NonExist"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.Has(tt.args.name)
			if got != tt.want {
				t.Errorf("unexpected mismatch: got: %v, want: %v. args: %+v", got, tt.want, tt.args)
			}
		})
	}
}

func TestGetType(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() unexpected error [%v] occured.", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		want      reflect.Type
		wantPanic bool
	}{
		{
			name:      "valid name and it's type is bytes",
			args:      args{name: "Bytes"},
			want:      reflect.TypeOf(testStructPtr.Bytes),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is string",
			args:      args{name: "String"},
			want:      reflect.TypeOf(testStructPtr.String),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is string (2nd)",
			args:      args{name: "String"},
			want:      reflect.TypeOf(testStructPtr.String),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is int64",
			args:      args{name: "Int64"},
			want:      reflect.TypeOf(testStructPtr.Int64),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is uint64",
			args:      args{name: "Uint64"},
			want:      reflect.TypeOf(testStructPtr.Uint64),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is float32",
			args:      args{name: "Float32"},
			want:      reflect.TypeOf(testStructPtr.Float32),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is float64",
			args:      args{name: "Float64"},
			want:      reflect.TypeOf(testStructPtr.Float64),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is bool",
			args:      args{name: "Bool"},
			want:      reflect.TypeOf(testStructPtr.Bool),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is map",
			args:      args{name: "Map"},
			want:      reflect.TypeOf(testStructPtr.Map),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is func",
			args:      args{name: "Func"},
			want:      reflect.TypeOf(testStructPtr.Func),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is chan int",
			args:      args{name: "ChInt"},
			want:      reflect.TypeOf(testStructPtr.ChInt),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is struct",
			args:      args{name: "TestStruct2"},
			want:      reflect.TypeOf(testStructPtr.TestStruct2),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is struct ptr",
			args:      args{name: "TestStruct2Ptr"},
			want:      reflect.TypeOf(testStructPtr.TestStruct2Ptr),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is struct slice",
			args:      args{name: "TestStruct4Slice"},
			want:      reflect.TypeOf(testStructPtr.TestStruct4Slice),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is struct ptr slice",
			args:      args{name: "TestStruct4PtrSlice"},
			want:      reflect.TypeOf(testStructPtr.TestStruct4PtrSlice),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is string and unexported field",
			args:      args{name: "privateString"},
			want:      reflect.TypeOf(testStructPtr.privateString),
			wantPanic: false,
		},
		{
			name:      "name does not exist",
			args:      args{name: "XXX"},
			want:      reflect.TypeOf(nil),
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferPanic(t, tt.wantPanic, false, tt.args)

			got := a.GetType(tt.args.name)
			if d := cmp.Diff(got.String(), tt.want.String()); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestGetValue(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		want      reflect.Value
		wantPanic bool
	}{
		{
			name:      "valid name and it's type is bytes",
			args:      args{name: "Bytes"},
			want:      reflect.ValueOf(testStructPtr.Bytes),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is string",
			args:      args{name: "String"},
			want:      reflect.ValueOf(testStructPtr.String),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is string (2nd)",
			args:      args{name: "String"},
			want:      reflect.ValueOf(testStructPtr.String),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is int64",
			args:      args{name: "Int64"},
			want:      reflect.ValueOf(testStructPtr.Int64),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is uint64",
			args:      args{name: "Uint64"},
			want:      reflect.ValueOf(testStructPtr.Uint64),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is float32",
			args:      args{name: "Float32"},
			want:      reflect.ValueOf(testStructPtr.Float32),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is float64",
			args:      args{name: "Float64"},
			want:      reflect.ValueOf(testStructPtr.Float64),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is bool",
			args:      args{name: "Bool"},
			want:      reflect.ValueOf(testStructPtr.Bool),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is map",
			args:      args{name: "Map"},
			want:      reflect.ValueOf(testStructPtr.Map),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is func",
			args:      args{name: "Func"},
			want:      reflect.ValueOf(testStructPtr.Func),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is chan int",
			args:      args{name: "ChInt"},
			want:      reflect.ValueOf(testStructPtr.ChInt),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is struct",
			args:      args{name: "TestStruct2"},
			want:      reflect.ValueOf(testStructPtr.TestStruct2),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is struct ptr",
			args:      args{name: "TestStruct2Ptr"},
			want:      reflect.ValueOf(testStructPtr.TestStruct2), // is not ptr
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is struct slice",
			args:      args{name: "TestStruct4Slice"},
			want:      reflect.ValueOf(testStructPtr.TestStruct4Slice),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is struct ptr slice",
			args:      args{name: "TestStruct4PtrSlice"},
			want:      reflect.ValueOf(testStructPtr.TestStruct4PtrSlice),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is string and unexported field",
			args:      args{name: "privateString"},
			want:      reflect.ValueOf(testStructPtr.privateString),
			wantPanic: false,
		},
		{
			name:      "name does not exist",
			args:      args{name: "XXX"},
			want:      reflect.ValueOf(nil),
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferPanic(t, tt.wantPanic, false, tt.args)

			got := a.GetValue(tt.args.name)
			if d := cmp.Diff(got.String(), tt.want.String()); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		want      interface{}
		wantPanic bool
		cmpopts   []cmp.Option
	}{
		{
			name:      "valid name and it's type is bytes",
			args:      args{name: "Bytes"},
			want:      testStructPtr.Bytes,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is string",
			args:      args{name: "String"},
			want:      testStructPtr.String,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is string (2nd)",
			args:      args{name: "String"},
			want:      testStructPtr.String,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is int64",
			args:      args{name: "Int64"},
			want:      testStructPtr.Int64,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is uint64",
			args:      args{name: "Uint64"},
			want:      testStructPtr.Uint64,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is float32",
			args:      args{name: "Float32"},
			want:      testStructPtr.Float32,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is float64",
			args:      args{name: "Float64"},
			want:      testStructPtr.Float64,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is bool",
			args:      args{name: "Bool"},
			want:      testStructPtr.Bool,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is map",
			args:      args{name: "Map"},
			want:      testStructPtr.Map,
			wantPanic: false,
		},
		// TODO: test fail when func
		// {
		// 	name:      "valid name and it's type is func",
		// 	args:      args{name: "Func"},
		// 	want:      testStructPtr.Func,
		// 	wantPanic: false,
		// },
		{
			name:      "valid name and it's type is chan int",
			args:      args{name: "ChInt"},
			want:      testStructPtr.ChInt,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is struct",
			args:      args{name: "TestStruct2"},
			want:      testStructPtr.TestStruct2,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is struct ptr",
			args:      args{name: "TestStruct2Ptr"},
			want:      *testStructPtr.TestStruct2Ptr,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is struct slice",
			args:      args{name: "TestStruct4Slice"},
			want:      testStructPtr.TestStruct4Slice,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is struct ptr slice",
			args:      args{name: "TestStruct4PtrSlice"},
			want:      testStructPtr.TestStruct4PtrSlice,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is string and unexported field",
			args:      args{name: "privateString"},
			want:      nil, // unexported field is nil
			wantPanic: false,
		},
		{
			name:      "name does not exist",
			args:      args{name: "XXX"},
			want:      nil,
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferPanic(t, tt.wantPanic, false, tt.args)

			got := a.Get(tt.args.name)
			if d := cmp.Diff(got, tt.want, tt.cmpopts...); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestBytes(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		want      reflect.Value
		wantPanic bool
	}{
		{
			name:      "valid name and it's type is bytes",
			args:      args{name: "Bytes"},
			want:      reflect.ValueOf(testStructPtr.Bytes),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is string",
			args:      args{name: "String"},
			want:      reflect.ValueOf(testStructPtr.String),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is string (2nd)",
			args:      args{name: "String"},
			want:      reflect.ValueOf(testStructPtr.String),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is int64",
			args:      args{name: "Int64"},
			want:      reflect.ValueOf(testStructPtr.Int64),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is uint64",
			args:      args{name: "Uint64"},
			want:      reflect.ValueOf(testStructPtr.Uint64),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is float32",
			args:      args{name: "Float32"},
			want:      reflect.ValueOf(testStructPtr.Float32),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is float64",
			args:      args{name: "Float64"},
			want:      reflect.ValueOf(testStructPtr.Float64),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is bool",
			args:      args{name: "Bool"},
			want:      reflect.ValueOf(testStructPtr.Bool),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is map",
			args:      args{name: "Map"},
			want:      reflect.ValueOf(testStructPtr.Map),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is func",
			args:      args{name: "Func"},
			want:      reflect.ValueOf(testStructPtr.Func),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is chan int",
			args:      args{name: "ChInt"},
			want:      reflect.ValueOf(testStructPtr.ChInt),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is struct ptr",
			args:      args{name: "TestStruct2"},
			want:      reflect.Indirect(reflect.ValueOf(testStructPtr.TestStruct2)),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is struct ptr slice",
			args:      args{name: "TestStruct4PtrSlice"},
			want:      reflect.ValueOf(testStructPtr.TestStruct4PtrSlice),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is string and unexported field",
			args:      args{name: "privateString"},
			want:      reflect.ValueOf(testStructPtr.privateString),
			wantPanic: true,
		},
		{
			name:      "name does not exist",
			args:      args{name: "XXX"},
			want:      reflect.ValueOf(nil),
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isXXX := a.IsBytes(tt.args.name)
			defer deferPanic(t, tt.wantPanic, isXXX, tt.args)

			got := a.Bytes(tt.args.name)
			if d := cmp.Diff(got, tt.want.Bytes()); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, IsInt64: %v, (-got +want)\n%s", tt.args, isXXX, d)
			}
		})
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		want      string
		wantPanic bool
	}{
		{
			name:      "valid name and it's type is bytes",
			args:      args{name: "Bytes"},
			want:      reflect.ValueOf(testStructPtr.Bytes).String(),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is string",
			args:      args{name: "String"},
			want:      reflect.ValueOf(testStructPtr.String).String(),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is string (2nd)",
			args:      args{name: "String"},
			want:      reflect.ValueOf(testStructPtr.String).String(),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is int64",
			args:      args{name: "Int64"},
			want:      reflect.ValueOf(testStructPtr.Int64).String(),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is uint64",
			args:      args{name: "Uint64"},
			want:      reflect.ValueOf(testStructPtr.Uint64).String(),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is float32",
			args:      args{name: "Float32"},
			want:      reflect.ValueOf(testStructPtr.Float32).String(),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is float64",
			args:      args{name: "Float64"},
			want:      reflect.ValueOf(testStructPtr.Float64).String(),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is bool",
			args:      args{name: "Bool"},
			want:      reflect.ValueOf(testStructPtr.Bool).String(),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is map",
			args:      args{name: "Map"},
			want:      reflect.ValueOf(testStructPtr.Map).String(),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is func",
			args:      args{name: "Func"},
			want:      reflect.ValueOf(testStructPtr.Func).String(),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is chan int",
			args:      args{name: "ChInt"},
			want:      reflect.ValueOf(testStructPtr.ChInt).String(),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is struct ptr",
			args:      args{name: "TestStruct2"},
			want:      reflect.Indirect(reflect.ValueOf(testStructPtr.TestStruct2)).String(),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is struct ptr slice",
			args:      args{name: "TestStruct4PtrSlice"},
			want:      reflect.ValueOf(testStructPtr.TestStruct4PtrSlice).String(),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is string and unexported field",
			args:      args{name: "privateString"},
			want:      reflect.ValueOf(testStructPtr.privateString).String(),
			wantPanic: false,
		},
		{
			name:      "name does not exist",
			args:      args{name: "XXX"},
			want:      reflect.ValueOf(nil).String(),
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isXXX := a.IsString(tt.args.name)
			defer deferPanic(t, tt.wantPanic, isXXX, tt.args)

			got := a.String(tt.args.name)
			if d := cmp.Diff(got, tt.want); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, IsString: %v, (-got +want)\n%s", tt.args, isXXX, d)
			}
		})
	}
}

func TestInt64(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		want      reflect.Value
		wantPanic bool
	}{
		{
			name:      "valid name and it's type is bytes",
			args:      args{name: "Bytes"},
			want:      reflect.ValueOf(testStructPtr.Bytes),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is string",
			args:      args{name: "String"},
			want:      reflect.ValueOf(testStructPtr.String),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is string (2nd)",
			args:      args{name: "String"},
			want:      reflect.ValueOf(testStructPtr.String),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is int64",
			args:      args{name: "Int64"},
			want:      reflect.ValueOf(testStructPtr.Int64),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is uint64",
			args:      args{name: "Uint64"},
			want:      reflect.ValueOf(testStructPtr.Uint64),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is float32",
			args:      args{name: "Float32"},
			want:      reflect.ValueOf(testStructPtr.Float32),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is float64",
			args:      args{name: "Float64"},
			want:      reflect.ValueOf(testStructPtr.Float64),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is bool",
			args:      args{name: "Bool"},
			want:      reflect.ValueOf(testStructPtr.Bool),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is map",
			args:      args{name: "Map"},
			want:      reflect.ValueOf(testStructPtr.Map),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is func",
			args:      args{name: "Func"},
			want:      reflect.ValueOf(testStructPtr.Func),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is chan int",
			args:      args{name: "ChInt"},
			want:      reflect.ValueOf(testStructPtr.ChInt),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is struct ptr",
			args:      args{name: "TestStruct2"},
			want:      reflect.Indirect(reflect.ValueOf(testStructPtr.TestStruct2)),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is struct ptr slice",
			args:      args{name: "TestStruct4PtrSlice"},
			want:      reflect.ValueOf(testStructPtr.TestStruct4PtrSlice),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is string and unexported field",
			args:      args{name: "privateString"},
			want:      reflect.ValueOf(testStructPtr.privateString),
			wantPanic: true,
		},
		{
			name:      "name does not exist",
			args:      args{name: "XXX"},
			want:      reflect.ValueOf(nil),
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isXXX := a.IsInt64(tt.args.name)
			defer deferPanic(t, tt.wantPanic, isXXX, tt.args)

			got := a.Int64(tt.args.name)
			if d := cmp.Diff(got, tt.want.Int()); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, IsInt64: %v, (-got +want)\n%s", tt.args, isXXX, d)
			}
		})
	}
}

func TestUint64(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		want      reflect.Value
		wantPanic bool
	}{
		{
			name:      "valid name and it's type is bytes",
			args:      args{name: "Bytes"},
			want:      reflect.ValueOf(testStructPtr.Bytes),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is string",
			args:      args{name: "String"},
			want:      reflect.ValueOf(testStructPtr.String),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is string (2nd)",
			args:      args{name: "String"},
			want:      reflect.ValueOf(testStructPtr.String),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is int64",
			args:      args{name: "Int64"},
			want:      reflect.ValueOf(testStructPtr.Int64),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is uint64",
			args:      args{name: "Uint64"},
			want:      reflect.ValueOf(testStructPtr.Uint64),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is float32",
			args:      args{name: "Float32"},
			want:      reflect.ValueOf(testStructPtr.Float32),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is float64",
			args:      args{name: "Float64"},
			want:      reflect.ValueOf(testStructPtr.Float64),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is bool",
			args:      args{name: "Bool"},
			want:      reflect.ValueOf(testStructPtr.Bool),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is map",
			args:      args{name: "Map"},
			want:      reflect.ValueOf(testStructPtr.Map),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is func",
			args:      args{name: "Func"},
			want:      reflect.ValueOf(testStructPtr.Func),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is chan int",
			args:      args{name: "ChInt"},
			want:      reflect.ValueOf(testStructPtr.ChInt),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is struct ptr",
			args:      args{name: "TestStruct2"},
			want:      reflect.Indirect(reflect.ValueOf(testStructPtr.TestStruct2)),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is struct ptr slice",
			args:      args{name: "TestStruct4PtrSlice"},
			want:      reflect.ValueOf(testStructPtr.TestStruct4PtrSlice),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is string and unexported field",
			args:      args{name: "privateString"},
			want:      reflect.ValueOf(testStructPtr.privateString),
			wantPanic: true,
		},
		{
			name:      "name does not exist",
			args:      args{name: "XXX"},
			want:      reflect.ValueOf(nil),
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isXXX := a.IsUint64(tt.args.name)
			defer deferPanic(t, tt.wantPanic, isXXX, tt.args)

			got := a.Uint64(tt.args.name)
			if d := cmp.Diff(got, tt.want.Uint()); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, IsUint64: %v, (-got +want)\n%s", tt.args, isXXX, d)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		want      reflect.Value
		wantPanic bool
	}{
		{
			name:      "valid name and it's type is bytes",
			args:      args{name: "Bytes"},
			want:      reflect.ValueOf(testStructPtr.Bytes),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is string",
			args:      args{name: "String"},
			want:      reflect.ValueOf(testStructPtr.String),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is string (2nd)",
			args:      args{name: "String"},
			want:      reflect.ValueOf(testStructPtr.String),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is int64",
			args:      args{name: "Int64"},
			want:      reflect.ValueOf(testStructPtr.Int64),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is uint64",
			args:      args{name: "Uint64"},
			want:      reflect.ValueOf(testStructPtr.Uint64),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is float32",
			args:      args{name: "Float32"},
			want:      reflect.ValueOf(testStructPtr.Float32),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is float64",
			args:      args{name: "Float64"},
			want:      reflect.ValueOf(testStructPtr.Float64),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is bool",
			args:      args{name: "Bool"},
			want:      reflect.ValueOf(testStructPtr.Bool),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is map",
			args:      args{name: "Map"},
			want:      reflect.ValueOf(testStructPtr.Map),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is func",
			args:      args{name: "Func"},
			want:      reflect.ValueOf(testStructPtr.Func),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is chan int",
			args:      args{name: "ChInt"},
			want:      reflect.ValueOf(testStructPtr.ChInt),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is struct ptr",
			args:      args{name: "TestStruct2"},
			want:      reflect.Indirect(reflect.ValueOf(testStructPtr.TestStruct2)),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is struct ptr slice",
			args:      args{name: "TestStruct4PtrSlice"},
			want:      reflect.ValueOf(testStructPtr.TestStruct4PtrSlice),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is string and unexported field",
			args:      args{name: "privateString"},
			want:      reflect.ValueOf(testStructPtr.privateString),
			wantPanic: true,
		},
		{
			name:      "name does not exist",
			args:      args{name: "XXX"},
			want:      reflect.ValueOf(nil),
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isXXX := a.IsFloat64(tt.args.name)
			defer deferPanic(t, tt.wantPanic, isXXX, tt.args)

			got := a.Float64(tt.args.name)
			if d := cmp.Diff(got, tt.want.Float()); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, IsFloat64: %v, (-got +want)\n%s", tt.args, isXXX, d)
			}
		})
	}
}

func TestBool(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		want      reflect.Value
		wantPanic bool
	}{
		{
			name:      "valid name and it's type is bytes",
			args:      args{name: "Bytes"},
			want:      reflect.ValueOf(testStructPtr.Bytes),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is string",
			args:      args{name: "String"},
			want:      reflect.ValueOf(testStructPtr.String),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is string (2nd)",
			args:      args{name: "String"},
			want:      reflect.ValueOf(testStructPtr.String),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is int64",
			args:      args{name: "Int64"},
			want:      reflect.ValueOf(testStructPtr.Int64),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is uint64",
			args:      args{name: "Uint64"},
			want:      reflect.ValueOf(testStructPtr.Uint64),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is float32",
			args:      args{name: "Float32"},
			want:      reflect.ValueOf(testStructPtr.Float32),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is float64",
			args:      args{name: "Float64"},
			want:      reflect.ValueOf(testStructPtr.Float64),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is bool",
			args:      args{name: "Bool"},
			want:      reflect.ValueOf(testStructPtr.Bool),
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is map",
			args:      args{name: "Map"},
			want:      reflect.ValueOf(testStructPtr.Map),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is func",
			args:      args{name: "Func"},
			want:      reflect.ValueOf(testStructPtr.Func),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is chan int",
			args:      args{name: "ChInt"},
			want:      reflect.ValueOf(testStructPtr.ChInt),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is struct ptr",
			args:      args{name: "TestStruct2"},
			want:      reflect.Indirect(reflect.ValueOf(testStructPtr.TestStruct2)),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is struct ptr slice",
			args:      args{name: "TestStruct4PtrSlice"},
			want:      reflect.ValueOf(testStructPtr.TestStruct4PtrSlice),
			wantPanic: true,
		},
		{
			name:      "valid name and it's type is string and unexported field",
			args:      args{name: "privateString"},
			want:      reflect.ValueOf(testStructPtr.privateString),
			wantPanic: true,
		},
		{
			name:      "name does not exist",
			args:      args{name: "XXX"},
			want:      reflect.ValueOf(nil),
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isXXX := a.IsBool(tt.args.name)
			defer deferPanic(t, tt.wantPanic, isXXX, tt.args)

			got := a.Bool(tt.args.name)
			if d := cmp.Diff(got, tt.want.Bool()); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, IsBool: %v, (-got +want)\n%s", tt.args, isXXX, d)
			}
		})
	}
}

func TestIsBytes(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid name and it's type is bytes",
			args: args{name: "Bytes"},
			want: true,
		},
		{
			name: "valid name and it's type is string",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is string (2nd)",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is int64",
			args: args{name: "Int64"},
			want: false,
		},
		{
			name: "valid name and it's type is uint64",
			args: args{name: "Uint64"},
			want: false,
		},
		{
			name: "valid name and it's type is float32",
			args: args{name: "Float32"},
			want: false,
		},
		{
			name: "valid name and it's type is float64",
			args: args{name: "Float64"},
			want: false,
		},
		{
			name: "valid name and it's type is bool",
			args: args{name: "Bool"},
			want: false,
		},
		{
			name: "valid name and it's type is map",
			args: args{name: "Map"},
			want: false,
		},
		{
			name: "valid name and it's type is func",
			args: args{name: "Func"},
			want: false,
		},
		{
			name: "valid name and it's type is chan int",
			args: args{name: "ChInt"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr slice",
			args: args{name: "TestStruct4PtrSlice"},
			want: false,
		},
		{
			name: "valid name and it's type is string and unexported field",
			args: args{name: "privateString"},
			want: false,
		},
		{
			name: "name does not exist",
			args: args{name: "XXX"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.IsBytes(tt.args.name)
			if got != tt.want {
				t.Errorf("unexpected mismatch: got: %v, want: %v. args: %+v", got, tt.want, tt.args)
			}
		})
	}
}

func TestIsString(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid name and it's type is bytes",
			args: args{name: "Bytes"},
			want: false,
		},
		{
			name: "valid name and it's type is string",
			args: args{name: "String"},
			want: true,
		},
		{
			name: "valid name and it's type is string (2nd)",
			args: args{name: "String"},
			want: true,
		},
		{
			name: "valid name and it's type is int64",
			args: args{name: "Int64"},
			want: false,
		},
		{
			name: "valid name and it's type is uint64",
			args: args{name: "Uint64"},
			want: false,
		},
		{
			name: "valid name and it's type is float32",
			args: args{name: "Float32"},
			want: false,
		},
		{
			name: "valid name and it's type is float64",
			args: args{name: "Float64"},
			want: false,
		},
		{
			name: "valid name and it's type is bool",
			args: args{name: "Bool"},
			want: false,
		},
		{
			name: "valid name and it's type is map",
			args: args{name: "Map"},
			want: false,
		},
		{
			name: "valid name and it's type is func",
			args: args{name: "Func"},
			want: false,
		},
		{
			name: "valid name and it's type is chan int",
			args: args{name: "ChInt"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr slice",
			args: args{name: "TestStruct4PtrSlice"},
			want: false,
		},
		{
			name: "valid name and it's type is string and unexported field",
			args: args{name: "privateString"},
			want: true,
		},
		{
			name: "name does not exist",
			args: args{name: "XXX"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.IsString(tt.args.name)
			if got != tt.want {
				t.Errorf("unexpected mismatch: got: %v, want: %v. args: %+v", got, tt.want, tt.args)
			}
		})
	}
}

func TestIsInt64(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid name and it's type is bytes",
			args: args{name: "Bytes"},
			want: false,
		},
		{
			name: "valid name and it's type is string",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is string (2nd)",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is int64",
			args: args{name: "Int64"},
			want: true,
		},
		{
			name: "valid name and it's type is uint64",
			args: args{name: "Uint64"},
			want: false,
		},
		{
			name: "valid name and it's type is float32",
			args: args{name: "Float32"},
			want: false,
		},
		{
			name: "valid name and it's type is float64",
			args: args{name: "Float64"},
			want: false,
		},
		{
			name: "valid name and it's type is bool",
			args: args{name: "Bool"},
			want: false,
		},
		{
			name: "valid name and it's type is map",
			args: args{name: "Map"},
			want: false,
		},
		{
			name: "valid name and it's type is func",
			args: args{name: "Func"},
			want: false,
		},
		{
			name: "valid name and it's type is chan int",
			args: args{name: "ChInt"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr slice",
			args: args{name: "TestStruct4PtrSlice"},
			want: false,
		},
		{
			name: "valid name and it's type is string and unexported field",
			args: args{name: "privateString"},
			want: false,
		},
		{
			name: "name does not exist",
			args: args{name: "XXX"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.IsInt64(tt.args.name)
			if got != tt.want {
				t.Errorf("unexpected mismatch: got: %v, want: %v. args: %+v", got, tt.want, tt.args)
			}
		})
	}
}

func TestIsUint64(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid name and it's type is bytes",
			args: args{name: "Bytes"},
			want: false,
		},
		{
			name: "valid name and it's type is string",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is string (2nd)",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is int64",
			args: args{name: "Int64"},
			want: false,
		},
		{
			name: "valid name and it's type is uint64",
			args: args{name: "Uint64"},
			want: true,
		},
		{
			name: "valid name and it's type is float32",
			args: args{name: "Float32"},
			want: false,
		},
		{
			name: "valid name and it's type is float64",
			args: args{name: "Float64"},
			want: false,
		},
		{
			name: "valid name and it's type is bool",
			args: args{name: "Bool"},
			want: false,
		},
		{
			name: "valid name and it's type is map",
			args: args{name: "Map"},
			want: false,
		},
		{
			name: "valid name and it's type is func",
			args: args{name: "Func"},
			want: false,
		},
		{
			name: "valid name and it's type is chan int",
			args: args{name: "ChInt"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr slice",
			args: args{name: "TestStruct4PtrSlice"},
			want: false,
		},
		{
			name: "valid name and it's type is string and unexported field",
			args: args{name: "privateString"},
			want: false,
		},
		{
			name: "name does not exist",
			args: args{name: "XXX"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.IsUint64(tt.args.name)
			if got != tt.want {
				t.Errorf("unexpected mismatch: got: %v, want: %v. args: %+v", got, tt.want, tt.args)
			}
		})
	}
}

func TestIsFloat64(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid name and it's type is bytes",
			args: args{name: "Bytes"},
			want: false,
		},
		{
			name: "valid name and it's type is string",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is string (2nd)",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is int64",
			args: args{name: "Int64"},
			want: false,
		},
		{
			name: "valid name and it's type is uint64",
			args: args{name: "Uint64"},
			want: false,
		},
		{
			name: "valid name and it's type is float32",
			args: args{name: "Float32"},
			want: false,
		},
		{
			name: "valid name and it's type is float64",
			args: args{name: "Float64"},
			want: true,
		},
		{
			name: "valid name and it's type is bool",
			args: args{name: "Bool"},
			want: false,
		},
		{
			name: "valid name and it's type is map",
			args: args{name: "Map"},
			want: false,
		},
		{
			name: "valid name and it's type is func",
			args: args{name: "Func"},
			want: false,
		},
		{
			name: "valid name and it's type is chan int",
			args: args{name: "ChInt"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr slice",
			args: args{name: "TestStruct4PtrSlice"},
			want: false,
		},
		{
			name: "valid name and it's type is string and unexported field",
			args: args{name: "privateString"},
			want: false,
		},
		{
			name: "name does not exist",
			args: args{name: "XXX"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.IsFloat64(tt.args.name)
			if got != tt.want {
				t.Errorf("unexpected mismatch: got: %v, want: %v. args: %+v", got, tt.want, tt.args)
			}
		})
	}
}

func TestIsBool(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid name and it's type is bytes",
			args: args{name: "Bytes"},
			want: false,
		},
		{
			name: "valid name and it's type is string",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is string (2nd)",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is int64",
			args: args{name: "Int64"},
			want: false,
		},
		{
			name: "valid name and it's type is uint64",
			args: args{name: "Uint64"},
			want: false,
		},
		{
			name: "valid name and it's type is float32",
			args: args{name: "Float32"},
			want: false,
		},
		{
			name: "valid name and it's type is float64",
			args: args{name: "Float64"},
			want: false,
		},
		{
			name: "valid name and it's type is bool",
			args: args{name: "Bool"},
			want: true,
		},
		{
			name: "valid name and it's type is map",
			args: args{name: "Map"},
			want: false,
		},
		{
			name: "valid name and it's type is func",
			args: args{name: "Func"},
			want: false,
		},
		{
			name: "valid name and it's type is chan int",
			args: args{name: "ChInt"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr slice",
			args: args{name: "TestStruct4PtrSlice"},
			want: false,
		},
		{
			name: "valid name and it's type is string and unexported field",
			args: args{name: "privateString"},
			want: false,
		},
		{
			name: "name does not exist",
			args: args{name: "XXX"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.IsBool(tt.args.name)
			if got != tt.want {
				t.Errorf("unexpected mismatch: got: %v, want: %v. args: %+v", got, tt.want, tt.args)
			}
		})
	}
}

func TestIsMap(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid name and it's type is bytes",
			args: args{name: "Bytes"},
			want: false,
		},
		{
			name: "valid name and it's type is string",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is string (2nd)",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is int64",
			args: args{name: "Int64"},
			want: false,
		},
		{
			name: "valid name and it's type is uint64",
			args: args{name: "Uint64"},
			want: false,
		},
		{
			name: "valid name and it's type is float32",
			args: args{name: "Float32"},
			want: false,
		},
		{
			name: "valid name and it's type is float64",
			args: args{name: "Float64"},
			want: false,
		},
		{
			name: "valid name and it's type is bool",
			args: args{name: "Bool"},
			want: false,
		},
		{
			name: "valid name and it's type is map",
			args: args{name: "Map"},
			want: true,
		},
		{
			name: "valid name and it's type is func",
			args: args{name: "Func"},
			want: false,
		},
		{
			name: "valid name and it's type is chan int",
			args: args{name: "ChInt"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr slice",
			args: args{name: "TestStruct4PtrSlice"},
			want: false,
		},
		{
			name: "valid name and it's type is string and unexported field",
			args: args{name: "privateString"},
			want: false,
		},
		{
			name: "name does not exist",
			args: args{name: "XXX"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.IsMap(tt.args.name)
			if got != tt.want {
				t.Errorf("unexpected mismatch: got: %v, want: %v. args: %+v", got, tt.want, tt.args)
			}
		})
	}
}

func TestIsFunc(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid name and it's type is bytes",
			args: args{name: "Bytes"},
			want: false,
		},
		{
			name: "valid name and it's type is string",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is string (2nd)",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is int64",
			args: args{name: "Int64"},
			want: false,
		},
		{
			name: "valid name and it's type is uint64",
			args: args{name: "Uint64"},
			want: false,
		},
		{
			name: "valid name and it's type is float32",
			args: args{name: "Float32"},
			want: false,
		},
		{
			name: "valid name and it's type is float64",
			args: args{name: "Float64"},
			want: false,
		},
		{
			name: "valid name and it's type is bool",
			args: args{name: "Bool"},
			want: false,
		},
		{
			name: "valid name and it's type is map",
			args: args{name: "Map"},
			want: false,
		},
		{
			name: "valid name and it's type is func",
			args: args{name: "Func"},
			want: true,
		},
		{
			name: "valid name and it's type is chan int",
			args: args{name: "ChInt"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr slice",
			args: args{name: "TestStruct4PtrSlice"},
			want: false,
		},
		{
			name: "valid name and it's type is string and unexported field",
			args: args{name: "privateString"},
			want: false,
		},
		{
			name: "name does not exist",
			args: args{name: "XXX"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.IsFunc(tt.args.name)
			if got != tt.want {
				t.Errorf("unexpected mismatch: got: %v, want: %v. args: %+v", got, tt.want, tt.args)
			}
		})
	}
}

func TestIsChan(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid name and it's type is bytes",
			args: args{name: "Bytes"},
			want: false,
		},
		{
			name: "valid name and it's type is string",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is string (2nd)",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is int64",
			args: args{name: "Int64"},
			want: false,
		},
		{
			name: "valid name and it's type is uint64",
			args: args{name: "Uint64"},
			want: false,
		},
		{
			name: "valid name and it's type is float32",
			args: args{name: "Float32"},
			want: false,
		},
		{
			name: "valid name and it's type is float64",
			args: args{name: "Float64"},
			want: false,
		},
		{
			name: "valid name and it's type is bool",
			args: args{name: "Bool"},
			want: false,
		},
		{
			name: "valid name and it's type is map",
			args: args{name: "Map"},
			want: false,
		},
		{
			name: "valid name and it's type is func",
			args: args{name: "Func"},
			want: false,
		},
		{
			name: "valid name and it's type is chan int",
			args: args{name: "ChInt"},
			want: true,
		},
		{
			name: "valid name and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr slice",
			args: args{name: "TestStruct4PtrSlice"},
			want: false,
		},
		{
			name: "valid name and it's type is string and unexported field",
			args: args{name: "privateString"},
			want: false,
		},
		{
			name: "name does not exist",
			args: args{name: "XXX"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.IsChan(tt.args.name)
			if got != tt.want {
				t.Errorf("unexpected mismatch: got: %v, want: %v. args: %+v", got, tt.want, tt.args)
			}
		})
	}
}

func TestIsStruct(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid name and it's type is bytes",
			args: args{name: "Bytes"},
			want: false,
		},
		{
			name: "valid name and it's type is string",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is string (2nd)",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is int64",
			args: args{name: "Int64"},
			want: false,
		},
		{
			name: "valid name and it's type is uint64",
			args: args{name: "Uint64"},
			want: false,
		},
		{
			name: "valid name and it's type is float32",
			args: args{name: "Float32"},
			want: false,
		},
		{
			name: "valid name and it's type is float64",
			args: args{name: "Float64"},
			want: false,
		},
		{
			name: "valid name and it's type is bool",
			args: args{name: "Bool"},
			want: false,
		},
		{
			name: "valid name and it's type is map",
			args: args{name: "Map"},
			want: false,
		},
		{
			name: "valid name and it's type is func",
			args: args{name: "Func"},
			want: false,
		},
		{
			name: "valid name and it's type is chan int",
			args: args{name: "ChInt"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: true,
		},
		{
			name: "valid name and it's type is struct ptr slice",
			args: args{name: "TestStruct4PtrSlice"},
			want: false,
		},
		{
			name: "valid name and it's type is string and unexported field",
			args: args{name: "privateString"},
			want: false,
		},
		{
			name: "name does not exist",
			args: args{name: "XXX"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.IsStruct(tt.args.name)
			if got != tt.want {
				t.Errorf("unexpected mismatch: got: %v, want: %v. args: %+v", got, tt.want, tt.args)
			}
		})
	}
}

func TestIsSlice(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid name and it's type is bytes",
			args: args{name: "Bytes"},
			want: true,
		},
		{
			name: "valid name and it's type is string",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is string (2nd)",
			args: args{name: "String"},
			want: false,
		},
		{
			name: "valid name and it's type is int64",
			args: args{name: "Int64"},
			want: false,
		},
		{
			name: "valid name and it's type is uint64",
			args: args{name: "Uint64"},
			want: false,
		},
		{
			name: "valid name and it's type is float32",
			args: args{name: "Float32"},
			want: false,
		},
		{
			name: "valid name and it's type is float64",
			args: args{name: "Float64"},
			want: false,
		},
		{
			name: "valid name and it's type is bool",
			args: args{name: "Bool"},
			want: false,
		},
		{
			name: "valid name and it's type is map",
			args: args{name: "Map"},
			want: false,
		},
		{
			name: "valid name and it's type is func",
			args: args{name: "Func"},
			want: false,
		},
		{
			name: "valid name and it's type is chan int",
			args: args{name: "ChInt"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "valid name and it's type is struct ptr slice",
			args: args{name: "TestStruct4PtrSlice"},
			want: true,
		},
		{
			name: "valid name and it's type is string and unexported field",
			args: args{name: "privateString"},
			want: false,
		},
		{
			name: "name does not exist",
			args: args{name: "XXX"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.IsSlice(tt.args.name)
			if got != tt.want {
				t.Errorf("unexpected mismatch: got: %v, want: %v. args: %+v", got, tt.want, tt.args)
			}
		})
	}
}

func TestMapGet(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()

	a, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
		fn   func(int, Getter) interface{}
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPanic bool
		want      interface{}
		cmpopts   []cmp.Option
	}{
		{
			name:      "valid name and it's type is string",
			args:      args{name: "String"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is string (2nd)",
			args:      args{name: "String"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is int64",
			args:      args{name: "Int64"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is uint64",
			args:      args{name: "Uint64"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is float32",
			args:      args{name: "Float32"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is float64",
			args:      args{name: "Float64"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is bool",
			args:      args{name: "Bool"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is map",
			args:      args{name: "Map"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is func",
			args:      args{name: "Func"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is chan int",
			args:      args{name: "ChInt"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is struct",
			args:      args{name: "TestStruct2"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "valid name and it's type is struct ptr",
			args:      args{name: "TestStruct2Ptr"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name: "valid name and it's type is struct slice",
			args: args{
				name: "TestStruct4Slice",
				fn: func(i int, g Getter) interface{} {
					return g.String("String") + "=" + g.String("String2")
				},
			},
			wantErr:   false,
			wantPanic: false,
			want:      []interface{}{string("key100=value100"), string("key200=value200")},
		},
		{
			name: "valid name and it's type is struct ptr slice",
			args: args{
				name: "TestStruct4PtrSlice",
				fn: func(i int, g Getter) interface{} {
					return g.String("String") + ":" + g.String("String2")
				},
			},
			wantErr:   false,
			wantPanic: false,
			want:      []interface{}{string("key991:value991"), string("key992:value992")},
		},
		{
			name:      "valid name and it's type is string and unexported field",
			args:      args{name: "privateString"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "name does not exist",
			args:      args{name: "XXX"},
			wantErr:   true,
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferPanic(t, tt.wantPanic, false, tt.args)

			got, err := a.MapGet(tt.args.name, tt.args.fn)
			if err == nil {
				if d := cmp.Diff(got, tt.want, tt.cmpopts...); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else if !tt.wantErr {
				t.Errorf("MapGet() unexpected error %v occured. wantErr %v", err, tt.wantErr)
			}
		})
	}
}
