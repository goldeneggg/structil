package structil_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/goldeneggg/structil"
)

type TestStruct struct {
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
	uexpString     string
	TestStruct2
	TestStruct2Ptr     *TestStruct2
	TestStructSlice    []TestStruct4
	TestStructPtrSlice []*TestStruct4
}

type TestStruct2 struct {
	ExpString string
	*TestStruct3
}

type TestStruct3 struct {
	ExpString string
	ExpInt    int
}

type TestStruct4 struct {
	ExpString  string
	ExpString2 string
}

const (
	testString = "test name"
)

var (
	testString2 = "test name2"

	testFunc = func(s string) interface{} { return s + "-func" }

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
			name:      "name exists in accessor and it's type is struct slice ptr",
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
			name:      "name exists in accessor and it's type is struct slice ptr",
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
		want      reflect.Value
		wantPanic bool
	}{
		{
			name:      "name exists in accessor and it's type is string",
			args:      args{name: "ExpString"},
			want:      reflect.ValueOf(testStructPtr.ExpString),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is string (2nd)",
			args:      args{name: "ExpString"},
			want:      reflect.ValueOf(testStructPtr.ExpString),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is int64",
			args:      args{name: "ExpInt64"},
			want:      reflect.ValueOf(testStructPtr.ExpInt64),
			wantPanic: false, // TODO: should be true?
		},
		{
			name:      "name exists in accessor and it's type is uint64",
			args:      args{name: "ExpUint64"},
			want:      reflect.ValueOf(testStructPtr.ExpUint64),
			wantPanic: false, // TODO: should be true?
		},
		{
			name:      "name exists in accessor and it's type is float32",
			args:      args{name: "ExpFloat32"},
			want:      reflect.ValueOf(testStructPtr.ExpFloat32),
			wantPanic: false, // TODO: should be true?
		},
		{
			name:      "name exists in accessor and it's type is float64",
			args:      args{name: "ExpFloat64"},
			want:      reflect.ValueOf(testStructPtr.ExpFloat64),
			wantPanic: false, // TODO: should be true?
		},
		{
			name:      "name exists in accessor and it's type is bool",
			args:      args{name: "ExpBool"},
			want:      reflect.ValueOf(testStructPtr.ExpBool),
			wantPanic: false, // TODO: should be true?
		},
		{
			name:      "name exists in accessor and it's type is struct ptr",
			args:      args{name: "TestStruct2"},
			want:      reflect.Indirect(reflect.ValueOf(testStructPtr.TestStruct2)),
			wantPanic: false, // TODO: should be true?
		},
		{
			name:      "name exists in accessor and it's type is struct slice ptr",
			args:      args{name: "TestStructPtrSlice"},
			want:      reflect.ValueOf(testStructPtr.TestStructPtrSlice),
			wantPanic: false, // TODO: should be true?
		},
		{
			name:      "name exists in accessor and it's type is string and unexported field",
			args:      args{name: "uexpString"},
			want:      reflect.ValueOf(testStructPtr.uexpString),
			wantPanic: false, // TODO: should be true?
		},
		{
			name:      "name does not exist",
			args:      args{name: "XXX"},
			want:      reflect.ValueOf(nil),
			wantPanic: false, // TODO: should be true?
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isXXX := a.IsString(tt.args.name)
			defer deferPanic(t, tt.wantPanic, isXXX, tt.args)

			got := a.GetString(tt.args.name)
			if d := cmp.Diff(got, tt.want.String()); d != "" {
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
			wantPanic: true, // TODO: why true?
		},
		{
			name:      "name exists in accessor and it's type is float32",
			args:      args{name: "ExpFloat32"},
			want:      reflect.ValueOf(testStructPtr.ExpFloat32),
			wantPanic: true, // TODO: why true?
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
			name:      "name exists in accessor and it's type is struct ptr",
			args:      args{name: "TestStruct2"},
			want:      reflect.Indirect(reflect.ValueOf(testStructPtr.TestStruct2)),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is struct slice ptr",
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
			wantPanic: false, // TODO: why true?
		},
		{
			name:      "name exists in accessor and it's type is float32",
			args:      args{name: "ExpFloat32"},
			want:      reflect.ValueOf(testStructPtr.ExpFloat32),
			wantPanic: true, // TODO: why true?
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
			name:      "name exists in accessor and it's type is struct ptr",
			args:      args{name: "TestStruct2"},
			want:      reflect.Indirect(reflect.ValueOf(testStructPtr.TestStruct2)),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is struct slice ptr",
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
			wantPanic: false, // TODO: why false?
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
			name:      "name exists in accessor and it's type is struct ptr",
			args:      args{name: "TestStruct2"},
			want:      reflect.Indirect(reflect.ValueOf(testStructPtr.TestStruct2)),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is struct slice ptr",
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
			name:      "name exists in accessor and it's type is struct ptr",
			args:      args{name: "TestStruct2"},
			want:      reflect.Indirect(reflect.ValueOf(testStructPtr.TestStruct2)),
			wantPanic: true,
		},
		{
			name:      "name exists in accessor and it's type is struct slice ptr",
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
