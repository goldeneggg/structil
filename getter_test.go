package structil_test

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"unsafe"

	"github.com/google/go-cmp/cmp"

	. "github.com/goldeneggg/structil"
)

type (
	GetterTestStruct struct {
		Byte          byte
		Bytes         []byte
		String        string
		Stringptr     *string
		Int           int
		Int64         int64
		Uint          uint
		Uint64        uint64
		Uintptr       uintptr
		Float32       float32
		Float64       float64
		Bool          bool
		Complex64     complex64
		Complex128    complex128
		Unsafeptr     unsafe.Pointer
		Stringslice   []string
		Stringarray   [2]string
		Map           map[string]interface{}
		Func          func(string) interface{}
		ChInt         chan int
		privateString string
		GetterTestStruct2
		GetterTestStruct2Ptr      *GetterTestStruct2
		GetterTestStruct4Slice    []GetterTestStruct4
		GetterTestStruct4PtrSlice []*GetterTestStruct4
	}

	GetterTestStruct2 struct {
		String string
		*GetterTestStruct3
	}

	GetterTestStruct3 struct {
		String string
		Int    int
	}

	GetterTestStruct4 struct {
		String  string
		String2 string
	}
)

var (
	getterTestString2 = "test name2"
	getterTestFunc    = func(s string) interface{} { return s + "-func" }
	getterTestChan    = make(chan int)
)

func newGetterTestStruct() GetterTestStruct {
	return GetterTestStruct{
		Byte:          0x61,
		Bytes:         []byte{0x00, 0xFF},
		String:        "test name",
		Stringptr:     &getterTestString2,
		Int:           int(-2),
		Int64:         int64(-1),
		Uint:          uint(2),
		Uint64:        uint64(1),
		Uintptr:       uintptr(100),
		Float32:       float32(-1.23),
		Float64:       float64(-3.45),
		Bool:          true,
		Complex64:     complex64(1),
		Complex128:    complex128(1),
		Unsafeptr:     unsafe.Pointer(new(int)),
		Stringslice:   []string{"strslice1", "strslice2"},
		Stringarray:   [2]string{"strarray1", "strarray2"},
		Map:           map[string]interface{}{"k1": "v1", "k2": 2},
		Func:          getterTestFunc,
		ChInt:         getterTestChan,
		privateString: "unexported string",
		GetterTestStruct2: GetterTestStruct2{
			String: "struct2 string",
			GetterTestStruct3: &GetterTestStruct3{
				String: "struct3 string",
				Int:    -123,
			},
		},
		GetterTestStruct2Ptr: &GetterTestStruct2{
			String: "struct2 string ptr",
			GetterTestStruct3: &GetterTestStruct3{
				String: "struct3 string ptr",
				Int:    -456,
			},
		},
		GetterTestStruct4Slice: []GetterTestStruct4{
			{
				String:  "key100",
				String2: "value100",
			},
			{
				String:  "key200",
				String2: "value200",
			},
		},
		GetterTestStruct4PtrSlice: []*GetterTestStruct4{
			{
				String:  "key991",
				String2: "value991",
			},
			{
				String:  "key992",
				String2: "value992",
			},
		},
	}
}

func newGetterTestStructPtr() *GetterTestStruct {
	ts := newGetterTestStruct()
	return &ts
}

func newTestGetter() (Getter, error) {
	return NewGetter(newGetterTestStructPtr())
}

func deferGetterTestPanic(t *testing.T, wantPanic bool, args interface{}) {
	r := recover()
	if r != nil {
		msg := fmt.Sprintf("\n%v\n", r)
		for d := 0; ; d++ {
			pc, file, line, ok := runtime.Caller(d)
			if !ok {
				break
			}

			msg = msg + fmt.Sprintf(" -> %d: %s: %s:%d\n", d, runtime.FuncForPC(pc).Name(), file, line)
		}

		if wantPanic {
			t.Logf("OK panic is expected: args: %+v, %s", args, msg)
		} else {
			t.Errorf("unexpected panic occured: args: %+v, %s", args, msg)
		}
	} else {
		if wantPanic {
			t.Errorf("expect to occur panic but does not: args: %+v, %+v", args, r)
		}
	}
}

type getterTestArgs struct {
	name  string
	mapfn func(int, Getter) (interface{}, error)
}

type getterTest struct {
	name      string
	args      *getterTestArgs
	wantBool  bool
	wantIntf  interface{}
	wantType  reflect.Type
	wantValue reflect.Value
	wantError bool
	wantPanic bool
}

func newGetterTests() []*getterTest {
	return []*getterTest{
		{
			name: "Byte",
			args: &getterTestArgs{name: "Byte"},
		},
		{
			name: "Bytes",
			args: &getterTestArgs{name: "Bytes"},
		},
		{
			name: "String",
			args: &getterTestArgs{name: "String"},
		},
		{
			name: "Int",
			args: &getterTestArgs{name: "Int"},
		},
		{
			name: "Int64",
			args: &getterTestArgs{name: "Int64"},
		},
		{
			name: "Uint",
			args: &getterTestArgs{name: "Uint"},
		},
		{
			name: "Uint64",
			args: &getterTestArgs{name: "Uint64"},
		},
		{
			name: "Uintptr",
			args: &getterTestArgs{name: "Uintptr"},
		},
		{
			name: "Float32",
			args: &getterTestArgs{name: "Float32"},
		},
		{
			name: "Float64",
			args: &getterTestArgs{name: "Float64"},
		},
		{
			name: "Bool",
			args: &getterTestArgs{name: "Bool"},
		},
		{
			name: "Complex64",
			args: &getterTestArgs{name: "Complex64"},
		},
		{
			name: "Complex128",
			args: &getterTestArgs{name: "Complex128"},
		},
		{
			name: "Unsafeptr",
			args: &getterTestArgs{name: "Unsafeptr"},
		},
		{
			name: "Map",
			args: &getterTestArgs{name: "Map"},
		},
		{
			name: "Func",
			args: &getterTestArgs{name: "Func"},
		},
		{
			name: "ChInt",
			args: &getterTestArgs{name: "ChInt"},
		},
		{
			name: "GetterTestStruct2",
			args: &getterTestArgs{name: "GetterTestStruct2"},
		},
		{
			name: "GetterTestStruct2Ptr",
			args: &getterTestArgs{name: "GetterTestStruct2Ptr"},
		},
		{
			name: "GetterTestStruct4Slice",
			args: &getterTestArgs{name: "GetterTestStruct4Slice"},
		},
		{
			name: "GetterTestStruct4PtrSlice",
			args: &getterTestArgs{name: "GetterTestStruct4PtrSlice"},
		},
		{
			name: "Stringarray",
			args: &getterTestArgs{name: "Stringarray"},
		},
		{
			name: "privateString",
			args: &getterTestArgs{name: "privateString"},
		},
		{
			name: "NotExist",
			args: &getterTestArgs{name: "NotExist"},
		},
	}
}

func TestNewGetter(t *testing.T) {
	t.Parallel()

	testStructVal := newGetterTestStruct()
	testStructPtr := newGetterTestStructPtr()

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
			name:    "valid struct ptr nil",
			args:    args{i: (*GetterTestStruct)(nil)},
			wantErr: true,
		},
		{
			name:    "invalid (nil)",
			args:    args{i: nil},
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
				if tt.wantErr {
					t.Errorf("NewGetter() error did not occur. got: %v", got)
					return
				}

				if _, ok := got.(Getter); !ok {
					t.Errorf("NewGetter() want Getter but got %+v", got)
				}
			} else if !tt.wantErr {
				t.Errorf("NewGetter() unexpected error [%v] occured. wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNumField(t *testing.T) {
	t.Parallel()

	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "use GetterTestStruct",
			args: args{i: &GetterTestStruct{}},
			want: 25,
		},
		{
			name: "use GetterTestStruct2",
			args: args{i: &GetterTestStruct2{}},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGetter(tt.args.i)
			if err != nil {
				t.Errorf("NewGetter() unexpected error [%v] occured", err)
			}

			nf := g.NumField()
			if nf != tt.want {
				t.Errorf("unmatch NumField. got: %d, want: %d", nf, tt.want)
			}
		})
	}
}

func TestHas(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	tests := newGetterTests()
	for _, tt := range tests {
		switch tt.name {
		case "NotExist":
			tt.wantBool = false
		default:
			tt.wantBool = true
		}

		t.Run(tt.name, func(t *testing.T) {
			got := g.Has(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v. args: %+v", got, tt.wantBool, tt.args)
			}
		})
	}
}

func TestGetType(t *testing.T) {
	t.Parallel()

	testStructPtr := newGetterTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() unexpected error [%v] occured.", err)
		return
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Byte":
				tt.wantType = reflect.TypeOf(testStructPtr.Byte)
			case "Bytes":
				tt.wantType = reflect.TypeOf(testStructPtr.Bytes)
			case "String":
				tt.wantType = reflect.TypeOf(testStructPtr.String)
			case "Int":
				tt.wantType = reflect.TypeOf(testStructPtr.Int)
			case "Int64":
				tt.wantType = reflect.TypeOf(testStructPtr.Int64)
			case "Uint":
				tt.wantType = reflect.TypeOf(testStructPtr.Uint)
			case "Uint64":
				tt.wantType = reflect.TypeOf(testStructPtr.Uint64)
			case "Uintptr":
				tt.wantType = reflect.TypeOf(testStructPtr.Uintptr)
			case "Float32":
				tt.wantType = reflect.TypeOf(testStructPtr.Float32)
			case "Float64":
				tt.wantType = reflect.TypeOf(testStructPtr.Float64)
			case "Bool":
				tt.wantType = reflect.TypeOf(testStructPtr.Bool)
			case "Complex64":
				tt.wantType = reflect.TypeOf(testStructPtr.Complex64)
			case "Complex128":
				tt.wantType = reflect.TypeOf(testStructPtr.Complex128)
			case "Unsafeptr":
				tt.wantType = reflect.TypeOf(testStructPtr.Unsafeptr)
			case "Map":
				tt.wantType = reflect.TypeOf(testStructPtr.Map)
			case "Func":
				tt.wantType = reflect.TypeOf(testStructPtr.Func)
			case "ChInt":
				tt.wantType = reflect.TypeOf(testStructPtr.ChInt)
			case "GetterTestStruct2":
				tt.wantType = reflect.TypeOf(testStructPtr.GetterTestStruct2)
			case "GetterTestStruct2Ptr":
				tt.wantType = reflect.TypeOf(testStructPtr.GetterTestStruct2Ptr)
			case "GetterTestStruct4Slice":
				tt.wantType = reflect.TypeOf(testStructPtr.GetterTestStruct4Slice)
			case "GetterTestStruct4PtrSlice":
				tt.wantType = reflect.TypeOf(testStructPtr.GetterTestStruct4PtrSlice)
			case "Stringarray":
				tt.wantType = reflect.TypeOf(testStructPtr.Stringarray)
			case "privateString":
				tt.wantType = reflect.TypeOf(testStructPtr.privateString)
			case "NotExist":
				tt.wantPanic = true
			}

			defer deferGetterTestPanic(t, tt.wantPanic, tt.args)

			got := g.GetType(tt.args.name)
			if tt.wantPanic {
				t.Errorf("expected panic did not occur. args: %+v", tt.args)
			} else if d := cmp.Diff(got.String(), tt.wantType.String()); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestGetValue(t *testing.T) {
	t.Parallel()

	testStructPtr := newGetterTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Byte":
				tt.wantValue = reflect.ValueOf(testStructPtr.Byte)
			case "Bytes":
				tt.wantValue = reflect.ValueOf(testStructPtr.Bytes)
			case "String":
				tt.wantValue = reflect.ValueOf(testStructPtr.String)
			case "Int":
				tt.wantValue = reflect.ValueOf(testStructPtr.Int)
			case "Int64":
				tt.wantValue = reflect.ValueOf(testStructPtr.Int64)
			case "Uint":
				tt.wantValue = reflect.ValueOf(testStructPtr.Uint)
			case "Uint64":
				tt.wantValue = reflect.ValueOf(testStructPtr.Uint64)
			case "Uintptr":
				tt.wantValue = reflect.ValueOf(testStructPtr.Uintptr)
			case "Float32":
				tt.wantValue = reflect.ValueOf(testStructPtr.Float32)
			case "Float64":
				tt.wantValue = reflect.ValueOf(testStructPtr.Float64)
			case "Bool":
				tt.wantValue = reflect.ValueOf(testStructPtr.Bool)
			case "Complex64":
				tt.wantValue = reflect.ValueOf(testStructPtr.Complex64)
			case "Complex128":
				tt.wantValue = reflect.ValueOf(testStructPtr.Complex128)
			case "Unsafeptr":
				tt.wantValue = reflect.ValueOf(testStructPtr.Unsafeptr)
			case "Map":
				tt.wantValue = reflect.ValueOf(testStructPtr.Map)
			case "Func":
				tt.wantValue = reflect.ValueOf(testStructPtr.Func)
			case "ChInt":
				tt.wantValue = reflect.ValueOf(testStructPtr.ChInt)
			case "GetterTestStruct2":
				tt.wantValue = reflect.ValueOf(testStructPtr.GetterTestStruct2)
			case "GetterTestStruct2Ptr":
				tt.wantValue = reflect.ValueOf(testStructPtr.GetterTestStruct2) // Note: *NOT* GetterTestStruct2Ptr
			case "GetterTestStruct4Slice":
				tt.wantValue = reflect.ValueOf(testStructPtr.GetterTestStruct4Slice)
			case "GetterTestStruct4PtrSlice":
				tt.wantValue = reflect.ValueOf(testStructPtr.GetterTestStruct4PtrSlice)
			case "Stringarray":
				tt.wantValue = reflect.ValueOf(testStructPtr.Stringarray)
			case "privateString":
				tt.wantValue = reflect.ValueOf(testStructPtr.privateString)
			case "NotExist":
				tt.wantPanic = true
			}

			defer deferGetterTestPanic(t, tt.wantPanic, tt.args)

			got := g.GetValue(tt.args.name)
			if tt.wantPanic {
				t.Errorf("expected panic did not occur. args: %+v", tt.args)
			} else if d := cmp.Diff(got.String(), tt.wantValue.String()); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	testStructPtr := newGetterTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Byte":
				tt.wantIntf = testStructPtr.Byte
			case "Bytes":
				tt.wantIntf = testStructPtr.Bytes
			case "String":
				tt.wantIntf = testStructPtr.String
			case "Int":
				tt.wantIntf = testStructPtr.Int
			case "Int64":
				tt.wantIntf = testStructPtr.Int64
			case "Uint":
				tt.wantIntf = testStructPtr.Uint
			case "Uint64":
				tt.wantIntf = testStructPtr.Uint64
			case "Uintptr":
				tt.wantIntf = testStructPtr.Uintptr
			case "Float32":
				tt.wantIntf = testStructPtr.Float32
			case "Float64":
				tt.wantIntf = testStructPtr.Float64
			case "Bool":
				tt.wantIntf = testStructPtr.Bool
			case "Complex64":
				tt.wantIntf = testStructPtr.Complex64
			case "Complex128":
				tt.wantIntf = testStructPtr.Complex128
			case "Unsafeptr":
				tt.wantIntf = testStructPtr.Unsafeptr
			case "Map":
				tt.wantIntf = testStructPtr.Map
			case "Func":
				tt.wantIntf = testStructPtr.Func
			case "ChInt":
				tt.wantIntf = testStructPtr.ChInt
			case "GetterTestStruct2":
				tt.wantIntf = testStructPtr.GetterTestStruct2
			case "GetterTestStruct2Ptr":
				tt.wantIntf = *testStructPtr.GetterTestStruct2Ptr // Note: *NOT* testStructPtr.GetterTestStruct2Ptr
			case "GetterTestStruct4Slice":
				tt.wantIntf = testStructPtr.GetterTestStruct4Slice
			case "GetterTestStruct4PtrSlice":
				tt.wantIntf = testStructPtr.GetterTestStruct4PtrSlice
			case "Stringarray":
				tt.wantIntf = testStructPtr.Stringarray
			case "privateString":
				tt.wantIntf = nil // Note: unexported field is nil
			case "NotExist":
				tt.wantPanic = true
			}

			defer deferGetterTestPanic(t, tt.wantPanic, tt.args)

			got := g.Get(tt.args.name)
			if tt.wantPanic {
				t.Errorf("expected panic did not occur. args: %+v", tt.args)
			} else if tt.args.name == "Func" {
				// Note: cmp.Diff does not support comparing func and func
				gp := reflect.ValueOf(got).Pointer()
				wp := reflect.ValueOf(tt.wantIntf).Pointer()
				if gp != wp {
					t.Errorf("unexpected mismatch func type: gp: %v, wp: %v", gp, wp)
				}
			} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestEGet(t *testing.T) {
	t.Parallel()

	testStructPtr := newGetterTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Byte":
				tt.wantIntf = testStructPtr.Byte
			case "Bytes":
				tt.wantIntf = testStructPtr.Bytes
			case "String":
				tt.wantIntf = testStructPtr.String
			case "Int":
				tt.wantIntf = testStructPtr.Int
			case "Int64":
				tt.wantIntf = testStructPtr.Int64
			case "Uint":
				tt.wantIntf = testStructPtr.Uint
			case "Uint64":
				tt.wantIntf = testStructPtr.Uint64
			case "Uintptr":
				tt.wantIntf = testStructPtr.Uintptr
			case "Float32":
				tt.wantIntf = testStructPtr.Float32
			case "Float64":
				tt.wantIntf = testStructPtr.Float64
			case "Bool":
				tt.wantIntf = testStructPtr.Bool
			case "Complex64":
				tt.wantIntf = testStructPtr.Complex64
			case "Complex128":
				tt.wantIntf = testStructPtr.Complex128
			case "Unsafeptr":
				tt.wantIntf = testStructPtr.Unsafeptr
			case "Map":
				tt.wantIntf = testStructPtr.Map
			case "Func":
				tt.wantIntf = testStructPtr.Func
			case "ChInt":
				tt.wantIntf = testStructPtr.ChInt
			case "GetterTestStruct2":
				tt.wantIntf = testStructPtr.GetterTestStruct2
			case "GetterTestStruct2Ptr":
				tt.wantIntf = *testStructPtr.GetterTestStruct2Ptr // Note: *NOT* testStructPtr.GetterTestStruct2Ptr
			case "GetterTestStruct4Slice":
				tt.wantIntf = testStructPtr.GetterTestStruct4Slice
			case "GetterTestStruct4PtrSlice":
				tt.wantIntf = testStructPtr.GetterTestStruct4PtrSlice
			case "Stringarray":
				tt.wantIntf = testStructPtr.Stringarray
			case "privateString":
				tt.wantIntf = nil // Note: unexported field is nil
			case "NotExist":
				tt.wantError = true
			}

			defer deferGetterTestPanic(t, tt.wantPanic, tt.args)

			got, err := g.EGet(tt.args.name)
			if tt.wantPanic {
				t.Errorf("expected panic did not occur. args: %+v", tt.args)
				return
			}

			if err == nil {
				if tt.wantError {
					t.Errorf("error did not occur. got: %v", got)
					return
				}

				if tt.args.name == "Func" {
					// Note: cmp.Diff does not support comparing func and func
					gp := reflect.ValueOf(got).Pointer()
					wp := reflect.ValueOf(tt.wantIntf).Pointer()
					if gp != wp {
						t.Errorf("unexpected mismatch func type: gp: %v, wp: %v", gp, wp)
					}
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else if !tt.wantError {
				t.Errorf("unexpected error occured. wantError %v, err: %v", tt.wantError, err)
			}
		})
	}
}

func TestByte(t *testing.T) {
	t.Parallel()

	testStructPtr := newGetterTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Byte":
				tt.wantIntf = testStructPtr.Byte
			default:
				tt.wantPanic = true
			}

			defer deferGetterTestPanic(t, tt.wantPanic, tt.args)

			got := g.Byte(tt.args.name)
			if tt.wantPanic {
				t.Errorf("expected panic did not occur. args: %+v", tt.args)
			} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestBytes(t *testing.T) {
	t.Parallel()

	testStructPtr := newGetterTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Bytes":
				tt.wantIntf = testStructPtr.Bytes
			default:
				tt.wantPanic = true
			}

			defer deferGetterTestPanic(t, tt.wantPanic, tt.args)

			got := g.Bytes(tt.args.name)
			if tt.wantPanic {
				t.Errorf("expected panic did not occur. args: %+v", tt.args)
			} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	testStructPtr := newGetterTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "String":
				tt.wantIntf = testStructPtr.String
			default:
				tt.wantPanic = true
			}

			defer deferGetterTestPanic(t, tt.wantPanic, tt.args)

			got := g.String(tt.args.name)
			if tt.wantPanic {
				t.Errorf("expected panic did not occur. args: %+v", tt.args)
			} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestInt(t *testing.T) {
	t.Parallel()

	testStructPtr := newGetterTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Int":
				tt.wantIntf = testStructPtr.Int
			default:
				tt.wantPanic = true
			}

			defer deferGetterTestPanic(t, tt.wantPanic, tt.args)

			got := g.Int(tt.args.name)
			if tt.wantPanic {
				t.Errorf("expected panic did not occur. args: %+v", tt.args)
			} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestInt64(t *testing.T) {
	t.Parallel()

	testStructPtr := newGetterTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Int64":
				tt.wantIntf = testStructPtr.Int64
			default:
				tt.wantPanic = true
			}

			defer deferGetterTestPanic(t, tt.wantPanic, tt.args)

			got := g.Int64(tt.args.name)
			if tt.wantPanic {
				t.Errorf("expected panic did not occur. args: %+v", tt.args)
			} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestUint(t *testing.T) {
	t.Parallel()

	testStructPtr := newGetterTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Uint":
				tt.wantIntf = testStructPtr.Uint
			default:
				tt.wantPanic = true
			}

			defer deferGetterTestPanic(t, tt.wantPanic, tt.args)

			got := g.Uint(tt.args.name)
			if tt.wantPanic {
				t.Errorf("expected panic did not occur. args: %+v", tt.args)
			} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestUint64(t *testing.T) {
	t.Parallel()

	testStructPtr := newGetterTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Uint64":
				tt.wantIntf = testStructPtr.Uint64
			default:
				tt.wantPanic = true
			}

			defer deferGetterTestPanic(t, tt.wantPanic, tt.args)

			got := g.Uint64(tt.args.name)
			if tt.wantPanic {
				t.Errorf("expected panic did not occur. args: %+v", tt.args)
			} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestUintptr(t *testing.T) {
	t.Parallel()

	testStructPtr := newGetterTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Uintptr":
				tt.wantIntf = testStructPtr.Uintptr
			default:
				tt.wantPanic = true
			}

			defer deferGetterTestPanic(t, tt.wantPanic, tt.args)

			got := g.Uintptr(tt.args.name)
			if tt.wantPanic {
				t.Errorf("expected panic did not occur. args: %+v", tt.args)
			} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	t.Parallel()

	testStructPtr := newGetterTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Float64":
				tt.wantIntf = testStructPtr.Float64
			default:
				tt.wantPanic = true
			}

			defer deferGetterTestPanic(t, tt.wantPanic, tt.args)

			got := g.Float64(tt.args.name)
			if tt.wantPanic {
				t.Errorf("expected panic did not occur. args: %+v", tt.args)
			} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestBool(t *testing.T) {
	t.Parallel()

	testStructPtr := newGetterTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Bool":
				tt.wantIntf = testStructPtr.Bool
			default:
				tt.wantPanic = true
			}

			defer deferGetterTestPanic(t, tt.wantPanic, tt.args)

			got := g.Bool(tt.args.name)
			if tt.wantPanic {
				t.Errorf("expected panic did not occur. args: %+v", tt.args)
			} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestComplex64(t *testing.T) {
	t.Parallel()

	testStructPtr := newGetterTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Complex64":
				tt.wantIntf = testStructPtr.Complex64
			default:
				tt.wantPanic = true
			}

			defer deferGetterTestPanic(t, tt.wantPanic, tt.args)

			got := g.Complex64(tt.args.name)
			if tt.wantPanic {
				t.Errorf("expected panic did not occur. args: %+v", tt.args)
			} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestComplex128(t *testing.T) {
	t.Parallel()

	testStructPtr := newGetterTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Complex128":
				tt.wantIntf = testStructPtr.Complex128
			default:
				tt.wantPanic = true
			}

			defer deferGetterTestPanic(t, tt.wantPanic, tt.args)

			got := g.Complex128(tt.args.name)
			if tt.wantPanic {
				t.Errorf("expected panic did not occur. args: %+v", tt.args)
			} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestUnsafePointer(t *testing.T) {
	t.Parallel()

	testStructPtr := newGetterTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Unsafeptr":
				tt.wantIntf = testStructPtr.Unsafeptr
			default:
				tt.wantPanic = true
			}

			defer deferGetterTestPanic(t, tt.wantPanic, tt.args)

			got := g.UnsafePointer(tt.args.name)
			if tt.wantPanic {
				t.Errorf("expected panic did not occur. args: %+v", tt.args)
			} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		})
	}
}

func TestIsByte(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Byte":
				tt.wantBool = true
			}

			got := g.IsByte(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsBytes(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Bytes":
				tt.wantBool = true
			}

			got := g.IsBytes(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsString(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "String", "privateString": // Note: IsString can work for private string field
				tt.wantBool = true
			}

			got := g.IsString(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsInt(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Int":
				tt.wantBool = true
			}

			got := g.IsInt(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsInt64(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Int64":
				tt.wantBool = true
			}

			got := g.IsInt64(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsUint(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Uint":
				tt.wantBool = true
			}

			got := g.IsUint(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsUint64(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Uint64":
				tt.wantBool = true
			}

			got := g.IsUint64(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsUintptr(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Uintptr":
				tt.wantBool = true
			}

			got := g.IsUintptr(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsFloat64(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Float64":
				tt.wantBool = true
			}

			got := g.IsFloat64(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsBool(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Bool":
				tt.wantBool = true
			}

			got := g.IsBool(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsComplex64(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Complex64":
				tt.wantBool = true
			}

			got := g.IsComplex64(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsComplex128(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Complex128":
				tt.wantBool = true
			}

			got := g.IsComplex128(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsUnsafePointer(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Unsafeptr":
				tt.wantBool = true
			}

			got := g.IsUnsafePointer(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsMap(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Map":
				tt.wantBool = true
			}

			got := g.IsMap(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsFunc(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Func":
				tt.wantBool = true
			}

			got := g.IsFunc(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsChan(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "ChInt":
				tt.wantBool = true
			}

			got := g.IsChan(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsStruct(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "GetterTestStruct2", "GetterTestStruct2Ptr":
				tt.wantBool = true
			}

			got := g.IsStruct(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsSlice(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Bytes", "GetterTestStruct4Slice", "GetterTestStruct4PtrSlice":
				tt.wantBool = true
			}

			got := g.IsSlice(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsArray(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Stringarray":
				tt.wantBool = true
			}

			got := g.IsArray(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestMapGet(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "GetterTestStruct4Slice":
				tt.args.mapfn = func(i int, g Getter) (interface{}, error) {
					return g.String("String") + "=" + g.String("String2"), nil
				}
				tt.wantIntf = []interface{}{string("key100=value100"), string("key200=value200")}
			case "GetterTestStruct4PtrSlice":
				tt.args.mapfn = func(i int, g Getter) (interface{}, error) {
					return g.String("String") + ":" + g.String("String2"), nil
				}
				tt.wantIntf = []interface{}{string("key991:value991"), string("key992:value992")}
			default:
				tt.wantError = true
			}

			got, err := g.MapGet(tt.args.name, tt.args.mapfn)
			if err == nil {
				if tt.wantError {
					t.Errorf("error did not occur. got: %v", got)
					return
				}

				if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else if !tt.wantError {
				t.Errorf("MapGet() unexpected error %v occured. wantErr %v", err, tt.wantError)
			}
		})
	}
}

// benchmark tests

func BenchmarkNewGetter_Val(b *testing.B) {
	var g Getter
	var e error

	testStructVal := newGetterTestStruct()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g, e = NewGetter(testStructVal)
		if e == nil {
			_ = g
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", e)
		}
	}
}

func BenchmarkNewGetter_Ptr(b *testing.B) {
	var g Getter
	var e error

	testStructPtr := newGetterTestStructPtr()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g, e = NewGetter(testStructPtr)
		if e == nil {
			_ = g
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", e)
		}
	}
}

func BenchmarkGetterGetType_String(b *testing.B) {
	var t reflect.Type

	g, err := newTestGetter()
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t = g.GetType("String")
		_ = t
	}
}

func BenchmarkGetterGetValue_String(b *testing.B) {
	var v reflect.Value

	g, err := newTestGetter()
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v = g.GetValue("String")
		_ = v
	}
}

func BenchmarkGetterHas_String(b *testing.B) {
	var bl bool

	g, err := newTestGetter()
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bl = g.Has("String")
		_ = bl
	}
}

func BenchmarkGetterGet_String(b *testing.B) {
	var it interface{}

	g, err := newTestGetter()
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it = g.Get("String")
		_ = it
	}
}

func BenchmarkGetterEGet_String(b *testing.B) {
	var it interface{}

	g, err := newTestGetter()
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it, err = g.EGet("String")
		if err == nil {
			_ = it
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkGetterString(b *testing.B) {
	var str string

	g, err := newTestGetter()
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str = g.String("String")
		_ = str
	}
}

func BenchmarkGetterMapGet(b *testing.B) {
	var ia []interface{}

	g, err := newTestGetter()
	if err != nil {
		b.Fatalf("NewGetter() occurs unexpected error: %v", err)
		return
	}
	fn := func(i int, g Getter) (interface{}, error) {
		return g.String("String") + ":" + g.String("String2"), nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ia, err = g.MapGet("GetterTestStruct4PtrSlice", fn)
		if err == nil {
			_ = ia
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}
