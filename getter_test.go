package structil_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/goldeneggg/structil"
)

type (
	TestStruct struct {
		ExpBytes       []byte
		ExpInt64       int64
		ExpUint64      uint64
		ExpFloat32     float32
		ExpFloat64     float64
		ExpString      string
		ExpStringptr   *string
		ExpStringslice []string
		ExpBool        bool
		ExpMap         map[string]interface{}
		ExpFunc        func(string) interface{}
		ExpChInt       chan int
		uexpString     string
		TestStruct2
		TestStruct2Ptr     *TestStruct2
		TestStructSlice    []TestStruct4
		TestStructPtrSlice []*TestStruct4
	}

	TestStruct2 struct {
		ExpString string
		*TestStruct3
	}

	TestStruct3 struct {
		ExpString string
		ExpInt    int
	}

	TestStruct4 struct {
		ExpString  string
		ExpString2 string
	}
)

const (
	testString = "test name"
)

var (
	testString2 = "test name2"
	testFunc    = func(s string) interface{} { return s + "-func" }
	testChan    = make(chan int)

	deferPanic = func(t *testing.T, wantPanic bool, isXXX bool, args interface{}) {
		r := recover()
		if r != nil {
			if !wantPanic {
				t.Errorf("unexpected panic occured: isXXX: %v, args: %+v, %+v", isXXX, args, r)
			}
		} else {
			if wantPanic {
				t.Errorf("expect to occur panic but does not: isXXX: %v, args: %+v, %+v", isXXX, args, r)
			}
		}
	}
)

func newTestStruct() TestStruct {
	return TestStruct{
		ExpBytes:       []byte{0x00, 0xFF},
		ExpInt64:       int64(-1),
		ExpUint64:      uint64(1),
		ExpFloat32:     float32(-1.23),
		ExpFloat64:     float64(-3.45),
		ExpString:      testString,
		ExpStringptr:   &testString2,
		ExpStringslice: []string{"strslice1", "strslice2"},
		ExpBool:        true,
		ExpMap:         map[string]interface{}{"k1": "v1", "k2": 2},
		ExpFunc:        testFunc,
		ExpChInt:       testChan,
		uexpString:     "unexported string",
		TestStruct2: TestStruct2{
			ExpString: "struct2 string",
			TestStruct3: &TestStruct3{
				ExpString: "struct3 string",
				ExpInt:    -123,
			},
		},
		TestStruct2Ptr: &TestStruct2{
			ExpString: "struct2 string ptr",
			TestStruct3: &TestStruct3{
				ExpString: "struct3 string ptr",
				ExpInt:    -456,
			},
		},
		TestStructSlice: []TestStruct4{
			{
				ExpString:  "key100",
				ExpString2: "value100",
			},
			{
				ExpString:  "key200",
				ExpString2: "value200",
			},
		},
		TestStructPtrSlice: []*TestStruct4{
			{
				ExpString:  "key991",
				ExpString2: "value991",
			},
			{
				ExpString:  "key992",
				ExpString2: "value992",
			},
		},
	}
}

func newTestStructPtr() *TestStruct {
	ts := newTestStruct()
	return &ts
}

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
			name:    "arg i is valid struct",
			args:    args{i: testStructVal},
			wantErr: false,
		},
		{
			name:    "arg i is valid struct ptr",
			args:    args{i: testStructPtr},
			wantErr: false,
		},
		{
			name:    "arg i is invalid (nil)",
			args:    args{i: nil},
			wantErr: true,
		},
		{
			name:    "arg i is invalid (struct nil)",
			args:    args{i: (*TestStruct)(nil)},
			wantErr: true,
		},
		{
			name:    "arg i is invalid (string)",
			args:    args{i: "abc"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGetter(tt.args.i)

			if err == nil {
				if _, ok := got.(Getter); !ok {
					t.Errorf("NewGetter() does not return Getter: %+v", got)
				}
			} else if !tt.wantErr {
				t.Errorf("NewGetter() unexpected error %v occured. wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetRT(t *testing.T) {
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
		want      reflect.Type
		wantPanic bool
	}{
		{
			name:      "name exists in accessor and it's type is bytes",
			args:      args{name: "ExpBytes"},
			want:      reflect.TypeOf(testStructPtr.ExpBytes),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is string",
			args:      args{name: "ExpString"},
			want:      reflect.TypeOf(testStructPtr.ExpString),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is string (2nd)",
			args:      args{name: "ExpString"},
			want:      reflect.TypeOf(testStructPtr.ExpString),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is int64",
			args:      args{name: "ExpInt64"},
			want:      reflect.TypeOf(testStructPtr.ExpInt64),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is uint64",
			args:      args{name: "ExpUint64"},
			want:      reflect.TypeOf(testStructPtr.ExpUint64),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is float32",
			args:      args{name: "ExpFloat32"},
			want:      reflect.TypeOf(testStructPtr.ExpFloat32),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is float64",
			args:      args{name: "ExpFloat64"},
			want:      reflect.TypeOf(testStructPtr.ExpFloat64),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is bool",
			args:      args{name: "ExpBool"},
			want:      reflect.TypeOf(testStructPtr.ExpBool),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is map",
			args:      args{name: "ExpMap"},
			want:      reflect.TypeOf(testStructPtr.ExpMap),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is func",
			args:      args{name: "ExpFunc"},
			want:      reflect.TypeOf(testStructPtr.ExpFunc),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is chan int",
			args:      args{name: "ExpChInt"},
			want:      reflect.TypeOf(testStructPtr.ExpChInt),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is struct",
			args:      args{name: "TestStruct2"},
			want:      reflect.TypeOf(testStructPtr.TestStruct2),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is struct ptr",
			args:      args{name: "TestStruct2Ptr"},
			want:      reflect.TypeOf(testStructPtr.TestStruct2Ptr),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is struct slice",
			args:      args{name: "TestStructSlice"},
			want:      reflect.TypeOf(testStructPtr.TestStructSlice),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is struct ptr slice",
			args:      args{name: "TestStructPtrSlice"},
			want:      reflect.TypeOf(testStructPtr.TestStructPtrSlice),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is string and unexported field",
			args:      args{name: "uexpString"},
			want:      reflect.TypeOf(testStructPtr.uexpString),
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

			got := a.GetRT(tt.args.name)
			if d := cmp.Diff(got.String(), tt.want.String()); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
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
			name: "name exists in accessor and it's type is string",
			args: args{name: "ExpString"},
			want: true,
		},
		{
			name: "name exists in accessor and it's type is string (2nd)",
			args: args{name: "ExpString"},
			want: true,
		},
		{
			name: "name does not exist in accessor",
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
			name:      "name exists in accessor and it's type is bytes",
			args:      args{name: "ExpBytes"},
			want:      testStructPtr.ExpBytes,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is string",
			args:      args{name: "ExpString"},
			want:      testStructPtr.ExpString,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is string (2nd)",
			args:      args{name: "ExpString"},
			want:      testStructPtr.ExpString,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is int64",
			args:      args{name: "ExpInt64"},
			want:      testStructPtr.ExpInt64,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is uint64",
			args:      args{name: "ExpUint64"},
			want:      testStructPtr.ExpUint64,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is float32",
			args:      args{name: "ExpFloat32"},
			want:      testStructPtr.ExpFloat32,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is float64",
			args:      args{name: "ExpFloat64"},
			want:      testStructPtr.ExpFloat64,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is bool",
			args:      args{name: "ExpBool"},
			want:      testStructPtr.ExpBool,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is map",
			args:      args{name: "ExpMap"},
			want:      testStructPtr.ExpMap,
			wantPanic: false,
		},
		// TODO: test fail when func
		// {
		// 	name:      "name exists in accessor and it's type is func",
		// 	args:      args{name: "ExpFunc"},
		// 	want:      testStructPtr.ExpFunc,
		// 	wantPanic: false,
		// },
		{
			name:      "name exists in accessor and it's type is chan int",
			args:      args{name: "ExpChInt"},
			want:      testStructPtr.ExpChInt,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is struct",
			args:      args{name: "TestStruct2"},
			want:      testStructPtr.TestStruct2,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is struct ptr",
			args:      args{name: "TestStruct2Ptr"},
			want:      *testStructPtr.TestStruct2Ptr,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is struct slice",
			args:      args{name: "TestStructSlice"},
			want:      testStructPtr.TestStructSlice,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is struct ptr slice",
			args:      args{name: "TestStructPtrSlice"},
			want:      testStructPtr.TestStructPtrSlice,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is string and unexported field",
			args:      args{name: "uexpString"},
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

func TestGetBytes(t *testing.T) {
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
			name:      "name exists in accessor and it's type is bytes",
			args:      args{name: "ExpBytes"},
			want:      reflect.ValueOf(testStructPtr.ExpBytes),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is string",
			args:      args{name: "ExpString"},
			want:      reflect.ValueOf(testStructPtr.ExpString),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is string (2nd)",
			args:      args{name: "ExpString"},
			want:      reflect.ValueOf(testStructPtr.ExpString),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is int64",
			args:      args{name: "ExpInt64"},
			want:      reflect.ValueOf(testStructPtr.ExpInt64),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is uint64",
			args:      args{name: "ExpUint64"},
			want:      reflect.ValueOf(testStructPtr.ExpUint64),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is float32",
			args:      args{name: "ExpFloat32"},
			want:      reflect.ValueOf(testStructPtr.ExpFloat32),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is float64",
			args:      args{name: "ExpFloat64"},
			want:      reflect.ValueOf(testStructPtr.ExpFloat64),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is bool",
			args:      args{name: "ExpBool"},
			want:      reflect.ValueOf(testStructPtr.ExpBool),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is map",
			args:      args{name: "ExpMap"},
			want:      reflect.ValueOf(testStructPtr.ExpMap),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is func",
			args:      args{name: "ExpFunc"},
			want:      reflect.ValueOf(testStructPtr.ExpFunc),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is chan int",
			args:      args{name: "ExpChInt"},
			want:      reflect.ValueOf(testStructPtr.ExpChInt),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is struct ptr",
			args:      args{name: "TestStruct2"},
			want:      reflect.Indirect(reflect.ValueOf(testStructPtr.TestStruct2)),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is struct ptr slice",
			args:      args{name: "TestStructPtrSlice"},
			want:      reflect.ValueOf(testStructPtr.TestStructPtrSlice),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is string and unexported field",
			args:      args{name: "uexpString"},
			want:      reflect.ValueOf(testStructPtr.uexpString),
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

			got := a.GetBytes(tt.args.name)
			if d := cmp.Diff(got, tt.want.Bytes()); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, IsInt64: %v, (-got +want)\n%s", tt.args, isXXX, d)
			}
		})
	}
}

func TestGetString(t *testing.T) {
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
			name:      "name exists in accessor and it's type is bytes",
			args:      args{name: "ExpBytes"},
			want:      reflect.ValueOf(testStructPtr.ExpBytes).String(),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is string",
			args:      args{name: "ExpString"},
			want:      reflect.ValueOf(testStructPtr.ExpString).String(),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is string (2nd)",
			args:      args{name: "ExpString"},
			want:      reflect.ValueOf(testStructPtr.ExpString).String(),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is int64",
			args:      args{name: "ExpInt64"},
			want:      reflect.ValueOf(testStructPtr.ExpInt64).String(),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is uint64",
			args:      args{name: "ExpUint64"},
			want:      reflect.ValueOf(testStructPtr.ExpUint64).String(),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is float32",
			args:      args{name: "ExpFloat32"},
			want:      reflect.ValueOf(testStructPtr.ExpFloat32).String(),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is float64",
			args:      args{name: "ExpFloat64"},
			want:      reflect.ValueOf(testStructPtr.ExpFloat64).String(),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is bool",
			args:      args{name: "ExpBool"},
			want:      reflect.ValueOf(testStructPtr.ExpBool).String(),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is map",
			args:      args{name: "ExpMap"},
			want:      reflect.ValueOf(testStructPtr.ExpMap).String(),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is func",
			args:      args{name: "ExpFunc"},
			want:      reflect.ValueOf(testStructPtr.ExpFunc).String(),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is chan int",
			args:      args{name: "ExpChInt"},
			want:      reflect.ValueOf(testStructPtr.ExpChInt).String(),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is struct ptr",
			args:      args{name: "TestStruct2"},
			want:      reflect.Indirect(reflect.ValueOf(testStructPtr.TestStruct2)).String(),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is struct ptr slice",
			args:      args{name: "TestStructPtrSlice"},
			want:      reflect.ValueOf(testStructPtr.TestStructPtrSlice).String(),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is string and unexported field",
			args:      args{name: "uexpString"},
			want:      reflect.ValueOf(testStructPtr.uexpString).String(),
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

			got := a.GetString(tt.args.name)
			if d := cmp.Diff(got, tt.want); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, IsString: %v, (-got +want)\n%s", tt.args, isXXX, d)
			}
		})
	}
}

func TestGetInt64(t *testing.T) {
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
			name:      "name exists in accessor and it's type is bytes",
			args:      args{name: "ExpBytes"},
			want:      reflect.ValueOf(testStructPtr.ExpBytes),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is string",
			args:      args{name: "ExpString"},
			want:      reflect.ValueOf(testStructPtr.ExpString),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is string (2nd)",
			args:      args{name: "ExpString"},
			want:      reflect.ValueOf(testStructPtr.ExpString),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is int64",
			args:      args{name: "ExpInt64"},
			want:      reflect.ValueOf(testStructPtr.ExpInt64),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is uint64",
			args:      args{name: "ExpUint64"},
			want:      reflect.ValueOf(testStructPtr.ExpUint64),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is float32",
			args:      args{name: "ExpFloat32"},
			want:      reflect.ValueOf(testStructPtr.ExpFloat32),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is float64",
			args:      args{name: "ExpFloat64"},
			want:      reflect.ValueOf(testStructPtr.ExpFloat64),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is bool",
			args:      args{name: "ExpBool"},
			want:      reflect.ValueOf(testStructPtr.ExpBool),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is map",
			args:      args{name: "ExpMap"},
			want:      reflect.ValueOf(testStructPtr.ExpMap),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is func",
			args:      args{name: "ExpFunc"},
			want:      reflect.ValueOf(testStructPtr.ExpFunc),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is chan int",
			args:      args{name: "ExpChInt"},
			want:      reflect.ValueOf(testStructPtr.ExpChInt),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is struct ptr",
			args:      args{name: "TestStruct2"},
			want:      reflect.Indirect(reflect.ValueOf(testStructPtr.TestStruct2)),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is struct ptr slice",
			args:      args{name: "TestStructPtrSlice"},
			want:      reflect.ValueOf(testStructPtr.TestStructPtrSlice),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is string and unexported field",
			args:      args{name: "uexpString"},
			want:      reflect.ValueOf(testStructPtr.uexpString),
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

			got := a.GetInt64(tt.args.name)
			if d := cmp.Diff(got, tt.want.Int()); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, IsInt64: %v, (-got +want)\n%s", tt.args, isXXX, d)
			}
		})
	}
}

func TestGetUint64(t *testing.T) {
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
			name:      "name exists in accessor and it's type is bytes",
			args:      args{name: "ExpBytes"},
			want:      reflect.ValueOf(testStructPtr.ExpBytes),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is string",
			args:      args{name: "ExpString"},
			want:      reflect.ValueOf(testStructPtr.ExpString),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is string (2nd)",
			args:      args{name: "ExpString"},
			want:      reflect.ValueOf(testStructPtr.ExpString),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is int64",
			args:      args{name: "ExpInt64"},
			want:      reflect.ValueOf(testStructPtr.ExpInt64),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is uint64",
			args:      args{name: "ExpUint64"},
			want:      reflect.ValueOf(testStructPtr.ExpUint64),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is float32",
			args:      args{name: "ExpFloat32"},
			want:      reflect.ValueOf(testStructPtr.ExpFloat32),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is float64",
			args:      args{name: "ExpFloat64"},
			want:      reflect.ValueOf(testStructPtr.ExpFloat64),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is bool",
			args:      args{name: "ExpBool"},
			want:      reflect.ValueOf(testStructPtr.ExpBool),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is map",
			args:      args{name: "ExpMap"},
			want:      reflect.ValueOf(testStructPtr.ExpMap),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is func",
			args:      args{name: "ExpFunc"},
			want:      reflect.ValueOf(testStructPtr.ExpFunc),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is chan int",
			args:      args{name: "ExpChInt"},
			want:      reflect.ValueOf(testStructPtr.ExpChInt),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is struct ptr",
			args:      args{name: "TestStruct2"},
			want:      reflect.Indirect(reflect.ValueOf(testStructPtr.TestStruct2)),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is struct ptr slice",
			args:      args{name: "TestStructPtrSlice"},
			want:      reflect.ValueOf(testStructPtr.TestStructPtrSlice),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is string and unexported field",
			args:      args{name: "uexpString"},
			want:      reflect.ValueOf(testStructPtr.uexpString),
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

			got := a.GetUint64(tt.args.name)
			if d := cmp.Diff(got, tt.want.Uint()); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, IsUint64: %v, (-got +want)\n%s", tt.args, isXXX, d)
			}
		})
	}
}

func TestGetFloat64(t *testing.T) {
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
			name:      "name exists in accessor and it's type is bytes",
			args:      args{name: "ExpBytes"},
			want:      reflect.ValueOf(testStructPtr.ExpBytes),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is string",
			args:      args{name: "ExpString"},
			want:      reflect.ValueOf(testStructPtr.ExpString),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is string (2nd)",
			args:      args{name: "ExpString"},
			want:      reflect.ValueOf(testStructPtr.ExpString),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is int64",
			args:      args{name: "ExpInt64"},
			want:      reflect.ValueOf(testStructPtr.ExpInt64),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is uint64",
			args:      args{name: "ExpUint64"},
			want:      reflect.ValueOf(testStructPtr.ExpUint64),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is float32",
			args:      args{name: "ExpFloat32"},
			want:      reflect.ValueOf(testStructPtr.ExpFloat32),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is float64",
			args:      args{name: "ExpFloat64"},
			want:      reflect.ValueOf(testStructPtr.ExpFloat64),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is bool",
			args:      args{name: "ExpBool"},
			want:      reflect.ValueOf(testStructPtr.ExpBool),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is map",
			args:      args{name: "ExpMap"},
			want:      reflect.ValueOf(testStructPtr.ExpMap),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is func",
			args:      args{name: "ExpFunc"},
			want:      reflect.ValueOf(testStructPtr.ExpFunc),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is chan int",
			args:      args{name: "ExpChInt"},
			want:      reflect.ValueOf(testStructPtr.ExpChInt),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is struct ptr",
			args:      args{name: "TestStruct2"},
			want:      reflect.Indirect(reflect.ValueOf(testStructPtr.TestStruct2)),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is struct ptr slice",
			args:      args{name: "TestStructPtrSlice"},
			want:      reflect.ValueOf(testStructPtr.TestStructPtrSlice),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is string and unexported field",
			args:      args{name: "uexpString"},
			want:      reflect.ValueOf(testStructPtr.uexpString),
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

			got := a.GetFloat64(tt.args.name)
			if d := cmp.Diff(got, tt.want.Float()); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, IsFloat64: %v, (-got +want)\n%s", tt.args, isXXX, d)
			}
		})
	}
}

func TestGetBool(t *testing.T) {
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
			name:      "name exists in accessor and it's type is bytes",
			args:      args{name: "ExpBytes"},
			want:      reflect.ValueOf(testStructPtr.ExpBytes),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is string",
			args:      args{name: "ExpString"},
			want:      reflect.ValueOf(testStructPtr.ExpString),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is string (2nd)",
			args:      args{name: "ExpString"},
			want:      reflect.ValueOf(testStructPtr.ExpString),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is int64",
			args:      args{name: "ExpInt64"},
			want:      reflect.ValueOf(testStructPtr.ExpInt64),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is uint64",
			args:      args{name: "ExpUint64"},
			want:      reflect.ValueOf(testStructPtr.ExpUint64),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is float32",
			args:      args{name: "ExpFloat32"},
			want:      reflect.ValueOf(testStructPtr.ExpFloat32),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is float64",
			args:      args{name: "ExpFloat64"},
			want:      reflect.ValueOf(testStructPtr.ExpFloat64),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is bool",
			args:      args{name: "ExpBool"},
			want:      reflect.ValueOf(testStructPtr.ExpBool),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is map",
			args:      args{name: "ExpMap"},
			want:      reflect.ValueOf(testStructPtr.ExpMap),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is func",
			args:      args{name: "ExpFunc"},
			want:      reflect.ValueOf(testStructPtr.ExpFunc),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is chan int",
			args:      args{name: "ExpChInt"},
			want:      reflect.ValueOf(testStructPtr.ExpChInt),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is struct ptr",
			args:      args{name: "TestStruct2"},
			want:      reflect.Indirect(reflect.ValueOf(testStructPtr.TestStruct2)),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is struct ptr slice",
			args:      args{name: "TestStructPtrSlice"},
			want:      reflect.ValueOf(testStructPtr.TestStructPtrSlice),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is string and unexported field",
			args:      args{name: "uexpString"},
			want:      reflect.ValueOf(testStructPtr.uexpString),
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

			got := a.GetBool(tt.args.name)
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
			name: "name exists in accessor and it's type is bytes",
			args: args{name: "ExpBytes"},
			want: true,
		},
		{
			name: "name exists in accessor and it's type is string",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string (2nd)",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is int64",
			args: args{name: "ExpInt64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is uint64",
			args: args{name: "ExpUint64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float32",
			args: args{name: "ExpFloat32"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float64",
			args: args{name: "ExpFloat64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is bool",
			args: args{name: "ExpBool"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is map",
			args: args{name: "ExpMap"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is func",
			args: args{name: "ExpFunc"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is chan int",
			args: args{name: "ExpChInt"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr slice",
			args: args{name: "TestStructPtrSlice"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string and unexported field",
			args: args{name: "uexpString"},
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
			name: "name exists in accessor and it's type is bytes",
			args: args{name: "ExpBytes"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string",
			args: args{name: "ExpString"},
			want: true,
		},
		{
			name: "name exists in accessor and it's type is string (2nd)",
			args: args{name: "ExpString"},
			want: true,
		},
		{
			name: "name exists in accessor and it's type is int64",
			args: args{name: "ExpInt64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is uint64",
			args: args{name: "ExpUint64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float32",
			args: args{name: "ExpFloat32"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float64",
			args: args{name: "ExpFloat64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is bool",
			args: args{name: "ExpBool"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is map",
			args: args{name: "ExpMap"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is func",
			args: args{name: "ExpFunc"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is chan int",
			args: args{name: "ExpChInt"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr slice",
			args: args{name: "TestStructPtrSlice"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string and unexported field",
			args: args{name: "uexpString"},
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
			name: "name exists in accessor and it's type is bytes",
			args: args{name: "ExpBytes"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string (2nd)",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is int64",
			args: args{name: "ExpInt64"},
			want: true,
		},
		{
			name: "name exists in accessor and it's type is uint64",
			args: args{name: "ExpUint64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float32",
			args: args{name: "ExpFloat32"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float64",
			args: args{name: "ExpFloat64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is bool",
			args: args{name: "ExpBool"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is map",
			args: args{name: "ExpMap"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is func",
			args: args{name: "ExpFunc"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is chan int",
			args: args{name: "ExpChInt"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr slice",
			args: args{name: "TestStructPtrSlice"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string and unexported field",
			args: args{name: "uexpString"},
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
			name: "name exists in accessor and it's type is bytes",
			args: args{name: "ExpBytes"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string (2nd)",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is int64",
			args: args{name: "ExpInt64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is uint64",
			args: args{name: "ExpUint64"},
			want: true,
		},
		{
			name: "name exists in accessor and it's type is float32",
			args: args{name: "ExpFloat32"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float64",
			args: args{name: "ExpFloat64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is bool",
			args: args{name: "ExpBool"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is map",
			args: args{name: "ExpMap"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is func",
			args: args{name: "ExpFunc"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is chan int",
			args: args{name: "ExpChInt"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr slice",
			args: args{name: "TestStructPtrSlice"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string and unexported field",
			args: args{name: "uexpString"},
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
			name: "name exists in accessor and it's type is bytes",
			args: args{name: "ExpBytes"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string (2nd)",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is int64",
			args: args{name: "ExpInt64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is uint64",
			args: args{name: "ExpUint64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float32",
			args: args{name: "ExpFloat32"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float64",
			args: args{name: "ExpFloat64"},
			want: true,
		},
		{
			name: "name exists in accessor and it's type is bool",
			args: args{name: "ExpBool"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is map",
			args: args{name: "ExpMap"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is func",
			args: args{name: "ExpFunc"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is chan int",
			args: args{name: "ExpChInt"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr slice",
			args: args{name: "TestStructPtrSlice"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string and unexported field",
			args: args{name: "uexpString"},
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
			name: "name exists in accessor and it's type is bytes",
			args: args{name: "ExpBytes"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string (2nd)",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is int64",
			args: args{name: "ExpInt64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is uint64",
			args: args{name: "ExpUint64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float32",
			args: args{name: "ExpFloat32"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float64",
			args: args{name: "ExpFloat64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is bool",
			args: args{name: "ExpBool"},
			want: true,
		},
		{
			name: "name exists in accessor and it's type is map",
			args: args{name: "ExpMap"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is func",
			args: args{name: "ExpFunc"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is chan int",
			args: args{name: "ExpChInt"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr slice",
			args: args{name: "TestStructPtrSlice"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string and unexported field",
			args: args{name: "uexpString"},
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
			name: "name exists in accessor and it's type is bytes",
			args: args{name: "ExpBytes"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string (2nd)",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is int64",
			args: args{name: "ExpInt64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is uint64",
			args: args{name: "ExpUint64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float32",
			args: args{name: "ExpFloat32"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float64",
			args: args{name: "ExpFloat64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is bool",
			args: args{name: "ExpBool"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is map",
			args: args{name: "ExpMap"},
			want: true,
		},
		{
			name: "name exists in accessor and it's type is func",
			args: args{name: "ExpFunc"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is chan int",
			args: args{name: "ExpChInt"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr slice",
			args: args{name: "TestStructPtrSlice"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string and unexported field",
			args: args{name: "uexpString"},
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
			name: "name exists in accessor and it's type is bytes",
			args: args{name: "ExpBytes"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string (2nd)",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is int64",
			args: args{name: "ExpInt64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is uint64",
			args: args{name: "ExpUint64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float32",
			args: args{name: "ExpFloat32"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float64",
			args: args{name: "ExpFloat64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is bool",
			args: args{name: "ExpBool"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is map",
			args: args{name: "ExpMap"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is func",
			args: args{name: "ExpFunc"},
			want: true,
		},
		{
			name: "name exists in accessor and it's type is chan int",
			args: args{name: "ExpChInt"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr slice",
			args: args{name: "TestStructPtrSlice"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string and unexported field",
			args: args{name: "uexpString"},
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
			name: "name exists in accessor and it's type is bytes",
			args: args{name: "ExpBytes"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string (2nd)",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is int64",
			args: args{name: "ExpInt64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is uint64",
			args: args{name: "ExpUint64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float32",
			args: args{name: "ExpFloat32"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float64",
			args: args{name: "ExpFloat64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is bool",
			args: args{name: "ExpBool"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is map",
			args: args{name: "ExpMap"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is func",
			args: args{name: "ExpFunc"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is chan int",
			args: args{name: "ExpChInt"},
			want: true,
		},
		{
			name: "name exists in accessor and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr slice",
			args: args{name: "TestStructPtrSlice"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string and unexported field",
			args: args{name: "uexpString"},
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
			name: "name exists in accessor and it's type is bytes",
			args: args{name: "ExpBytes"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string (2nd)",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is int64",
			args: args{name: "ExpInt64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is uint64",
			args: args{name: "ExpUint64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float32",
			args: args{name: "ExpFloat32"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float64",
			args: args{name: "ExpFloat64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is bool",
			args: args{name: "ExpBool"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is map",
			args: args{name: "ExpMap"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is func",
			args: args{name: "ExpFunc"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is chan int",
			args: args{name: "ExpChInt"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: true,
		},
		{
			name: "name exists in accessor and it's type is struct ptr slice",
			args: args{name: "TestStructPtrSlice"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string and unexported field",
			args: args{name: "uexpString"},
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
			name: "name exists in accessor and it's type is bytes",
			args: args{name: "ExpBytes"},
			want: true,
		},
		{
			name: "name exists in accessor and it's type is string",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is string (2nd)",
			args: args{name: "ExpString"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is int64",
			args: args{name: "ExpInt64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is uint64",
			args: args{name: "ExpUint64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float32",
			args: args{name: "ExpFloat32"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is float64",
			args: args{name: "ExpFloat64"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is bool",
			args: args{name: "ExpBool"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is map",
			args: args{name: "ExpMap"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is func",
			args: args{name: "ExpFunc"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is chan int",
			args: args{name: "ExpChInt"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr",
			args: args{name: "TestStruct2"},
			want: false,
		},
		{
			name: "name exists in accessor and it's type is struct ptr slice",
			args: args{name: "TestStructPtrSlice"},
			want: true,
		},
		{
			name: "name exists in accessor and it's type is string and unexported field",
			args: args{name: "uexpString"},
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
			name:      "name exists in accessor and it's type is string",
			args:      args{name: "ExpString"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is string (2nd)",
			args:      args{name: "ExpString"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is int64",
			args:      args{name: "ExpInt64"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is uint64",
			args:      args{name: "ExpUint64"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is float32",
			args:      args{name: "ExpFloat32"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is float64",
			args:      args{name: "ExpFloat64"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is bool",
			args:      args{name: "ExpBool"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is map",
			args:      args{name: "ExpMap"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is func",
			args:      args{name: "ExpFunc"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is chan int",
			args:      args{name: "ExpChInt"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is struct",
			args:      args{name: "TestStruct2"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is struct ptr",
			args:      args{name: "TestStruct2Ptr"},
			wantErr:   true,
			wantPanic: false,
		},
		{
			name: "name exists in accessor and it's type is struct slice",
			args: args{
				name: "TestStructSlice",
				fn: func(i int, g Getter) interface{} {
					return g.GetString("ExpString") + "=" + g.GetString("ExpString2")
				},
			},
			wantErr:   false,
			wantPanic: false,
			want:      []interface{}{string("key100=value100"), string("key200=value200")},
		},
		{
			name: "name exists in accessor and it's type is struct ptr slice",
			args: args{
				name: "TestStructPtrSlice",
				fn: func(i int, g Getter) interface{} {
					return g.GetString("ExpString") + ":" + g.GetString("ExpString2")
				},
			},
			wantErr:   false,
			wantPanic: false,
			want:      []interface{}{string("key991:value991"), string("key992:value992")},
		},
		{
			name:      "name exists in accessor and it's type is string and unexported field",
			args:      args{name: "uexpString"},
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
