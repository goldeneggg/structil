package structil_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"unsafe"

	. "github.com/goldeneggg/structil"
	"github.com/google/go-cmp/cmp"
)

type (
	GetterTestStruct struct {
		Byte          byte
		Bytes         []byte
		String        string
		Stringptr     *string
		Int           int
		Int8          int8
		Int16         int16
		Int32         int32
		Int64         int64
		Uint          uint
		Uint8         uint8
		Uint16        uint16
		Uint32        uint32
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
	// top level struct has 30 public fields and 1 private fields
	return GetterTestStruct{
		Byte:          math.MaxUint8,
		Bytes:         []byte{0x00, 0xFF},
		String:        "test name",
		Stringptr:     &getterTestString2,
		Int:           math.MinInt32,
		Int8:          math.MinInt8,
		Int16:         math.MinInt16,
		Int32:         math.MinInt32,
		Int64:         math.MinInt64,
		Uint:          math.MaxUint32,
		Uint8:         math.MaxUint8,
		Uint16:        math.MaxUint16,
		Uint32:        math.MaxUint32,
		Uint64:        math.MaxUint64,
		Uintptr:       uintptr(100),
		Float32:       math.SmallestNonzeroFloat32,
		Float64:       math.SmallestNonzeroFloat64,
		Bool:          true,
		Complex64:     complex(math.MaxFloat32, math.SmallestNonzeroFloat32),
		Complex128:    complex(math.MaxFloat64, math.SmallestNonzeroFloat64),
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

func newTestGetter() (*Getter, error) {
	return NewGetter(newGetterTestStructPtr())
}

type getterTestArgs struct {
	name  string
	mapfn func(int, *Getter) (interface{}, error)
}

type getterTest struct {
	name      string
	args      *getterTestArgs
	wantBool  bool
	wantIntf  interface{}
	wantType  reflect.Type
	wantValue reflect.Value
	wantError bool
	wantNotOK bool
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
			name: "Int8",
			args: &getterTestArgs{name: "Int8"},
		},
		{
			name: "Int16",
			args: &getterTestArgs{name: "Int16"},
		},
		{
			name: "Int32",
			args: &getterTestArgs{name: "Int32"},
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
			name: "Uint8",
			args: &getterTestArgs{name: "Uint8"},
		},
		{
			name: "Uint16",
			args: &getterTestArgs{name: "Uint16"},
		},
		{
			name: "Uint32",
			args: &getterTestArgs{name: "Uint32"},
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
			} else if !tt.wantErr {
				t.Errorf("NewGetter() unexpected error [%v] occurred. wantErr %v", err, tt.wantErr)
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
			want: 31,
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
				t.Errorf("NewGetter() unexpected error [%v] occurred", err)
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

func TestNames(t *testing.T) {
	t.Parallel()

	type args struct {
		i interface{}
	}
	tests := []struct {
		name                string
		args                args
		wantSecondFieldName string
		wantNamesLength     int
	}{
		{
			name:                "use GetterTestStruct",
			args:                args{i: &GetterTestStruct{}},
			wantSecondFieldName: "Bytes",
			wantNamesLength:     31,
		},
		{
			name:                "use GetterTestStruct2",
			args:                args{i: &GetterTestStruct2{}},
			wantSecondFieldName: "GetterTestStruct3",
			wantNamesLength:     2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGetter(tt.args.i)
			if err != nil {
				t.Errorf("NewGetter() unexpected error [%v] occurred", err)
			}

			names := g.Names()
			if len(names) != tt.wantNamesLength {
				t.Errorf("unmatch NumField. got: %d, want: %d", len(names), tt.wantNamesLength)
			}
			if names[1] != tt.wantSecondFieldName {
				t.Errorf("unmatch second field name. got: %s, want: %s", names[1], tt.wantSecondFieldName)
			}
		})
	}
}

func testGetSeries(t *testing.T, wantNotOK bool, wantError bool, fn func(*testing.T, *getterTest, *Getter)) {
	t.Parallel()

	testStructPtr := newGetterTestStructPtr()
	g, err := NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() unexpected error [%v] occurred.", err)
		return
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Byte":
				tt.wantType = reflect.TypeOf(testStructPtr.Byte)
				tt.wantValue = reflect.ValueOf(testStructPtr.Byte)
				tt.wantIntf = testStructPtr.Byte
			case "Bytes":
				tt.wantType = reflect.TypeOf(testStructPtr.Bytes)
				tt.wantValue = reflect.ValueOf(testStructPtr.Bytes)
				tt.wantIntf = testStructPtr.Bytes
			case "String":
				tt.wantType = reflect.TypeOf(testStructPtr.String)
				tt.wantValue = reflect.ValueOf(testStructPtr.String)
				tt.wantIntf = testStructPtr.String
			case "Int":
				tt.wantType = reflect.TypeOf(testStructPtr.Int)
				tt.wantValue = reflect.ValueOf(testStructPtr.Int)
				tt.wantIntf = testStructPtr.Int
			case "Int8":
				tt.wantType = reflect.TypeOf(testStructPtr.Int8)
				tt.wantValue = reflect.ValueOf(testStructPtr.Int8)
				tt.wantIntf = testStructPtr.Int8
			case "Int16":
				tt.wantType = reflect.TypeOf(testStructPtr.Int16)
				tt.wantValue = reflect.ValueOf(testStructPtr.Int16)
				tt.wantIntf = testStructPtr.Int16
			case "Int32":
				tt.wantType = reflect.TypeOf(testStructPtr.Int32)
				tt.wantValue = reflect.ValueOf(testStructPtr.Int32)
				tt.wantIntf = testStructPtr.Int32
			case "Int64":
				tt.wantType = reflect.TypeOf(testStructPtr.Int64)
				tt.wantValue = reflect.ValueOf(testStructPtr.Int64)
				tt.wantIntf = testStructPtr.Int64
			case "Uint":
				tt.wantType = reflect.TypeOf(testStructPtr.Uint)
				tt.wantValue = reflect.ValueOf(testStructPtr.Uint)
				tt.wantIntf = testStructPtr.Uint
			case "Uint8":
				tt.wantType = reflect.TypeOf(testStructPtr.Uint8)
				tt.wantValue = reflect.ValueOf(testStructPtr.Uint8)
				tt.wantIntf = testStructPtr.Uint8
			case "Uint16":
				tt.wantType = reflect.TypeOf(testStructPtr.Uint16)
				tt.wantValue = reflect.ValueOf(testStructPtr.Uint16)
				tt.wantIntf = testStructPtr.Uint16
			case "Uint32":
				tt.wantType = reflect.TypeOf(testStructPtr.Uint32)
				tt.wantValue = reflect.ValueOf(testStructPtr.Uint32)
				tt.wantIntf = testStructPtr.Uint32
			case "Uint64":
				tt.wantType = reflect.TypeOf(testStructPtr.Uint64)
				tt.wantValue = reflect.ValueOf(testStructPtr.Uint64)
				tt.wantIntf = testStructPtr.Uint64
			case "Uintptr":
				tt.wantType = reflect.TypeOf(testStructPtr.Uintptr)
				tt.wantValue = reflect.ValueOf(testStructPtr.Uintptr)
				tt.wantIntf = testStructPtr.Uintptr
			case "Float32":
				tt.wantType = reflect.TypeOf(testStructPtr.Float32)
				tt.wantValue = reflect.ValueOf(testStructPtr.Float32)
				tt.wantIntf = testStructPtr.Float32
			case "Float64":
				tt.wantType = reflect.TypeOf(testStructPtr.Float64)
				tt.wantValue = reflect.ValueOf(testStructPtr.Float64)
				tt.wantIntf = testStructPtr.Float64
			case "Bool":
				tt.wantType = reflect.TypeOf(testStructPtr.Bool)
				tt.wantValue = reflect.ValueOf(testStructPtr.Bool)
				tt.wantIntf = testStructPtr.Bool
			case "Complex64":
				tt.wantType = reflect.TypeOf(testStructPtr.Complex64)
				tt.wantValue = reflect.ValueOf(testStructPtr.Complex64)
				tt.wantIntf = testStructPtr.Complex64
			case "Complex128":
				tt.wantType = reflect.TypeOf(testStructPtr.Complex128)
				tt.wantValue = reflect.ValueOf(testStructPtr.Complex128)
				tt.wantIntf = testStructPtr.Complex128
			case "Unsafeptr":
				tt.wantType = reflect.TypeOf(testStructPtr.Unsafeptr)
				tt.wantValue = reflect.ValueOf(testStructPtr.Unsafeptr)
				tt.wantIntf = testStructPtr.Unsafeptr
			case "Map":
				tt.wantType = reflect.TypeOf(testStructPtr.Map)
				tt.wantValue = reflect.ValueOf(testStructPtr.Map)
				tt.wantIntf = testStructPtr.Map
			case "Func":
				tt.wantType = reflect.TypeOf(testStructPtr.Func)
				tt.wantValue = reflect.ValueOf(testStructPtr.Func)
				tt.wantIntf = testStructPtr.Func
			case "ChInt":
				tt.wantType = reflect.TypeOf(testStructPtr.ChInt)
				tt.wantValue = reflect.ValueOf(testStructPtr.ChInt)
				tt.wantIntf = testStructPtr.ChInt
			case "GetterTestStruct2":
				tt.wantType = reflect.TypeOf(testStructPtr.GetterTestStruct2)
				tt.wantValue = reflect.ValueOf(testStructPtr.GetterTestStruct2)
				tt.wantIntf = testStructPtr.GetterTestStruct2
			case "GetterTestStruct2Ptr":
				tt.wantType = reflect.TypeOf(testStructPtr.GetterTestStruct2Ptr)
				tt.wantValue = reflect.ValueOf(testStructPtr.GetterTestStruct2) // Note: *NOT* GetterTestStruct2Ptr
				tt.wantIntf = *testStructPtr.GetterTestStruct2Ptr               // Note: *NOT* testStructPtr.GetterTestStruct2Ptr
			case "GetterTestStruct4Slice":
				tt.wantType = reflect.TypeOf(testStructPtr.GetterTestStruct4Slice)
				tt.wantValue = reflect.ValueOf(testStructPtr.GetterTestStruct4Slice)
				tt.wantIntf = testStructPtr.GetterTestStruct4Slice
			case "GetterTestStruct4PtrSlice":
				tt.wantType = reflect.TypeOf(testStructPtr.GetterTestStruct4PtrSlice)
				tt.wantValue = reflect.ValueOf(testStructPtr.GetterTestStruct4PtrSlice)
				tt.wantIntf = testStructPtr.GetterTestStruct4PtrSlice
			case "Stringarray":
				tt.wantType = reflect.TypeOf(testStructPtr.Stringarray)
				tt.wantValue = reflect.ValueOf(testStructPtr.Stringarray)
				tt.wantIntf = testStructPtr.Stringarray
			case "privateString":
				tt.wantType = reflect.TypeOf(testStructPtr.privateString)
				tt.wantValue = reflect.ValueOf(testStructPtr.privateString)
				tt.wantIntf = nil // Note: unexported field is nil
			case "NotExist":
				tt.wantNotOK = wantNotOK
				tt.wantError = wantError
			}

			fn(t, tt, g)
		})
	}
}

func TestGetType(t *testing.T) {
	assertionFunc := func(t *testing.T, tt *getterTest, g *Getter) {
		got, ok := g.GetType(tt.args.name)

		if ok {
			if tt.wantNotOK {
				t.Errorf("expected ok is false but true. args: %+v", tt.args)
			} else if d := cmp.Diff(got.String(), tt.wantType.String()); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		} else {
			if !tt.wantNotOK {
				t.Errorf("expected ok is true but false. args: %+v", tt.args)
			}
		}
	}

	testGetSeries(t, true, false, assertionFunc)
}

func TestGetValue(t *testing.T) {
	assertionFunc := func(t *testing.T, tt *getterTest, g *Getter) {
		got, ok := g.GetValue(tt.args.name)

		if ok {
			if tt.wantNotOK {
				t.Errorf("expected ok is false but true. args: %+v", tt.args)
			} else if d := cmp.Diff(got.String(), tt.wantValue.String()); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
			}
		} else {
			if !tt.wantNotOK {
				t.Errorf("expected ok is true but false. args: %+v", tt.args)
			}
		}
	}

	testGetSeries(t, true, false, assertionFunc)
}

func TestGet(t *testing.T) {
	assertionFunc := func(t *testing.T, tt *getterTest, g *Getter) {
		got, ok := g.Get(tt.args.name)

		if ok {
			if tt.wantNotOK {
				t.Errorf("expected ok is false but true. args: %+v", tt.args)
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
		} else {
			if !tt.wantNotOK {
				t.Errorf("expected ok is true but false. args: %+v", tt.args)
			}
		}
	}

	testGetSeries(t, true, false, assertionFunc)
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
			case "Uint8":
				tt.wantIntf = testStructPtr.Uint8
			default:
				tt.wantNotOK = true
			}

			got, ok := g.Byte(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
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
				tt.wantNotOK = true
			}

			got, ok := g.Bytes(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
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
				tt.wantNotOK = true
			}

			got, ok := g.String(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
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
				tt.wantNotOK = true
			}

			got, ok := g.Int(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
			}
		})
	}
}

func TestInt8(t *testing.T) {
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
			case "Int8":
				tt.wantIntf = testStructPtr.Int8
			default:
				tt.wantNotOK = true
			}

			got, ok := g.Int8(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
			}
		})
	}
}

func TestInt16(t *testing.T) {
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
			case "Int16":
				tt.wantIntf = testStructPtr.Int16
			default:
				tt.wantNotOK = true
			}

			got, ok := g.Int16(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
			}
		})
	}
}

func TestInt32(t *testing.T) {
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
			case "Int32":
				tt.wantIntf = testStructPtr.Int32
			default:
				tt.wantNotOK = true
			}

			got, ok := g.Int32(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
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
				tt.wantNotOK = true
			}

			got, ok := g.Int64(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
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
				tt.wantNotOK = true
			}

			got, ok := g.Uint(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
			}
		})
	}
}

func TestUint8(t *testing.T) {
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
			case "Uint8":
				tt.wantIntf = testStructPtr.Uint8
			default:
				tt.wantNotOK = true
			}

			got, ok := g.Uint8(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
			}
		})
	}
}

func TestUint16(t *testing.T) {
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
			case "Uint16":
				tt.wantIntf = testStructPtr.Uint16
			default:
				tt.wantNotOK = true
			}

			got, ok := g.Uint16(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
			}
		})
	}
}

func TestUint32(t *testing.T) {
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
			case "Uint32":
				tt.wantIntf = testStructPtr.Uint32
			default:
				tt.wantNotOK = true
			}

			got, ok := g.Uint32(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
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
				tt.wantNotOK = true
			}

			got, ok := g.Uint64(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
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
				tt.wantNotOK = true
			}

			got, ok := g.Uintptr(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
			}
		})
	}
}

func TestFloat32(t *testing.T) {
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
			case "Float32":
				tt.wantIntf = testStructPtr.Float32
			default:
				tt.wantNotOK = true
			}

			got, ok := g.Float32(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
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
				tt.wantNotOK = true
			}

			got, ok := g.Float64(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
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
				tt.wantNotOK = true
			}

			got, ok := g.Bool(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
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
				tt.wantNotOK = true
			}

			got, ok := g.Complex64(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
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
				tt.wantNotOK = true
			}

			got, ok := g.Complex128(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
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
				tt.wantNotOK = true
			}

			got, ok := g.UnsafePointer(tt.args.name)

			if ok {
				if tt.wantNotOK {
					t.Errorf("expected ok is false but true. args: %+v", tt.args)
				} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
					t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
				}
			} else {
				if !tt.wantNotOK {
					t.Errorf("expected ok is true but false. args: %+v", tt.args)
				}
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
			case "Uint8":
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

func TestIsInt8(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Int8":
				tt.wantBool = true
			}

			got := g.IsInt8(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsInt16(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Int16":
				tt.wantBool = true
			}

			got := g.IsInt16(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsInt32(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Int32":
				tt.wantBool = true
			}

			got := g.IsInt32(tt.args.name)
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

func TestIsUint8(t *testing.T) {
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
			case "Uint8":
				tt.wantBool = true
			}

			got := g.IsUint8(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsUint16(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Uint16":
				tt.wantBool = true
			}

			got := g.IsUint16(tt.args.name)
			if got != tt.wantBool {
				t.Errorf("unexpected mismatch: got: %v, want: %v", got, tt.wantBool)
			}
		})
	}
}

func TestIsUint32(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Uint32":
				tt.wantBool = true
			}

			got := g.IsUint32(tt.args.name)
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

func TestIsFloat32(t *testing.T) {
	t.Parallel()

	g, err := newTestGetter()
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Float32":
				tt.wantBool = true
			}

			got := g.IsFloat32(tt.args.name)
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
				tt.args.mapfn = func(i int, g *Getter) (interface{}, error) {
					str, _ := g.String("String")
					str2, _ := g.String("String2")
					return fmt.Sprintf("%s=%s", str, str2), nil
				}
				tt.wantIntf = []interface{}{string("key100=value100"), string("key200=value200")}
			case "GetterTestStruct4PtrSlice":
				tt.args.mapfn = func(i int, g *Getter) (interface{}, error) {
					str, _ := g.String("String")
					str2, _ := g.String("String2")
					return fmt.Sprintf("%s:%s", str, str2), nil
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
				t.Errorf("MapGet() unexpected error %v occurred. wantErr %v", err, tt.wantError)
			}
		})
	}
}
