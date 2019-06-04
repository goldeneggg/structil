package structil_test

import (
	"io"
	"os"
	"reflect"
	"testing"

	. "github.com/goldeneggg/structil"
	"github.com/google/go-cmp/cmp"
)

type TestStruct struct {
	Exp_int64       int64
	Exp_float64     float64
	Exp_string      string
	Exp_stringptr   *string
	Exp_stringslice []string
	Exp_bool        bool
	Exp_func        func(string) string
	uexp_string     string
	*TestStruct2
	TestStructPtrSlice []*TestStruct4
}

type TestStruct2 struct {
	Exp_string string
	Writer     io.Writer
	*TestStruct3
}

type TestStruct3 struct {
	Exp_string string
	Exp_int    int
}

type TestStruct4 struct {
	Exp_string  string
	Exp_string2 string
}

const (
	testString = "test name"
)

var (
	testString2 = "test name2"

	testFunc = func(s string) string { return s + "-func" }

	deferGetterPanic = func(t *testing.T, wantPanic bool) {
		r := recover()
		if r != nil {
			if !wantPanic {
				t.Errorf("unexpected panic occured: %+v", r)
			}
		} else {
			if wantPanic {
				t.Errorf("expect to occur panic but does not: %+v", r)
			}
		}
	}

	testStructVal = TestStruct{
		Exp_int64:       1,
		Exp_float64:     1.23,
		Exp_string:      testString,
		Exp_stringptr:   &testString2,
		Exp_stringslice: []string{"strslice1", "strslice2"},
		Exp_bool:        true,
		Exp_func:        testFunc,
		uexp_string:     "unexported string",
		TestStruct2: &TestStruct2{
			Exp_string: "struct2 string",
			Writer:     os.Stdout,
			TestStruct3: &TestStruct3{
				Exp_string: "struct3 string",
				Exp_int:    123,
			},
		},
		TestStructPtrSlice: []*TestStruct4{
			{
				Exp_string:  "key100",
				Exp_string2: "value100",
			},
			{
				Exp_string:  "key200",
				Exp_string2: "value200",
			},
		},
	}

	testStructPtr = &TestStruct{
		Exp_int64:       9,
		Exp_float64:     9.23,
		Exp_string:      testString,
		Exp_stringptr:   &testString2,
		Exp_stringslice: []string{"strslice991", "strslice992"},
		Exp_bool:        false,
		Exp_func:        testFunc,
		uexp_string:     "unexported string 999",
		TestStruct2: &TestStruct2{
			Exp_string: "struct2 string999",
			Writer:     os.Stdout,
			TestStruct3: &TestStruct3{
				Exp_string: "struct3 string999",
				Exp_int:    999,
			},
		},
		TestStructPtrSlice: []*TestStruct4{
			{
				Exp_string:  "key901",
				Exp_string2: "value901",
			},
			{
				Exp_string:  "key902",
				Exp_string2: "value902",
			},
		},
	}
)

func TestNewGetter(t *testing.T) {
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
			} else {
				if !tt.wantErr {
					t.Errorf("NewGetter() unexpected error %v occured. wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestGetString(t *testing.T) {
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
			args:      args{name: "Exp_string"},
			want:      reflect.ValueOf(testStructPtr.Exp_string),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor and it's type is string (2nd)",
			args:      args{name: "Exp_string"},
			want:      reflect.ValueOf(testStructPtr.Exp_string),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor but it's type is int64",
			args:      args{name: "Exp_int64"},
			want:      reflect.ValueOf(testStructPtr.Exp_int64),
			wantPanic: false, // TODO: should be true?
		},
		{
			name:      "name exists in accessor but it's type is float64",
			args:      args{name: "Exp_float64"},
			want:      reflect.ValueOf(testStructPtr.Exp_float64),
			wantPanic: false, // TODO: should be true?
		},
		{
			name:      "name exists in accessor but it's type is bool",
			args:      args{name: "Exp_bool"},
			want:      reflect.ValueOf(testStructPtr.Exp_bool),
			wantPanic: false, // TODO: should be true?
		},
		{
			name:      "name exists in accessor but it's type is struct ptr",
			args:      args{name: "TestStruct2"},
			want:      reflect.Indirect(reflect.ValueOf(testStructPtr.TestStruct2)),
			wantPanic: false, // TODO: should be true?
		},
		{
			name:      "name exists in accessor but it's type is struct slice ptr",
			args:      args{name: "TestStructPtrSlice"},
			want:      reflect.ValueOf(testStructPtr.TestStructPtrSlice),
			wantPanic: false, // TODO: should be true?
		},
		{
			name:      "name exists in accessor and it's type is string and unexported field",
			args:      args{name: "uexp_string"},
			want:      reflect.ValueOf(testStructPtr.uexp_string),
			wantPanic: false,
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
			defer deferGetterPanic(t, tt.wantPanic)

			got := a.GetString(tt.args.name)
			if d := cmp.Diff(got, tt.want.String()); d != "" {
				t.Errorf("unexpected mismatch: (-got +want)\n%s", d)
			}
		})
	}
}
