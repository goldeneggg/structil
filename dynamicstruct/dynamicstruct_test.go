package dynamicstruct_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/goldeneggg/structil"
	. "github.com/goldeneggg/structil/dynamicstruct"
	"github.com/google/go-cmp/cmp"
)

type (
	DynamicTestStruct struct {
		Byte        byte
		Bytes       []byte
		Int         int
		Int64       int64
		Uint        uint
		Uint64      uint64
		Float32     float32
		Float64     float64
		String      string
		Stringptr   *string
		Stringslice []string
		Bool        bool
		Map         map[string]interface{}
		Func        func(string) interface{}
		// ChInt       chan int  // Note: type chan is not supported by mapstructure
		DynamicTestStruct2
		DynamicTestStruct2Ptr      *DynamicTestStruct2
		DynamicTestStruct4Slice    []DynamicTestStruct4
		DynamicTestStruct4PtrSlice []*DynamicTestStruct4
	}

	DynamicTestStruct2 struct {
		String string
		*DynamicTestStruct3
	}

	DynamicTestStruct3 struct {
		String string
		Int    int
	}

	DynamicTestStruct4 struct {
		String  string
		String2 string
	}
)

var (
	dynamicTestString2 = "test name2"
	dynamicTestFunc    = func(s string) interface{} { return s + "-func" }
	dynamicTestChan    = make(chan int)
)

func newDynamicTestStruct() DynamicTestStruct {
	return DynamicTestStruct{
		Byte:        0x61,
		Bytes:       []byte{0x00, 0xFF},
		Int:         int(-2),
		Int64:       int64(-1),
		Uint:        uint(2),
		Uint64:      uint64(1),
		Float32:     float32(-1.23),
		Float64:     float64(-3.45),
		String:      "test name",
		Stringptr:   &dynamicTestString2,
		Stringslice: []string{"strslice1", "strslice2"},
		Bool:        true,
		Map:         map[string]interface{}{"k1": "v1", "k2": 2},
		Func:        dynamicTestFunc,
		// ChInt:       dynamicTestChan, // Note: type chan is not supported by mapstructure
		DynamicTestStruct2: DynamicTestStruct2{
			String: "struct2 string",
			DynamicTestStruct3: &DynamicTestStruct3{
				String: "struct3 string",
				Int:    -123,
			},
		},
		DynamicTestStruct2Ptr: &DynamicTestStruct2{
			String: "struct2 string ptr",
			DynamicTestStruct3: &DynamicTestStruct3{
				String: "struct3 string ptr",
				Int:    -456,
			},
		},
		DynamicTestStruct4Slice: []DynamicTestStruct4{
			{
				String:  "key100",
				String2: "value100",
			},
			{
				String:  "key200",
				String2: "value200",
			},
		},
		DynamicTestStruct4PtrSlice: []*DynamicTestStruct4{
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

func newDynamicTestStructPtr() *DynamicTestStruct {
	ts := newDynamicTestStruct()
	return &ts
}

func newDynamicTestBuilder() Builder {
	return NewBuilder().
		AddString("StringField").
		AddInt("IntField").
		AddFloat("FloatField").
		AddBool("BoolField").
		AddMap("MapField", SampleString, SampleFloat).
		AddFunc("FuncField", []interface{}{SampleInt, SampleInt}, []interface{}{SampleBool, SampleError}).
		AddChanBoth("ChanBothField", SampleInt).
		AddChanRecv("ChanRecvField", SampleInt).
		AddChanSend("ChanSendField", SampleInt).
		AddStruct("StructField", newDynamicTestStruct(), false).
		AddStructPtr("StructPtrField", newDynamicTestStructPtr()).
		AddSlice("SliceField", newDynamicTestStructPtr())
}

func deferDynamicTestPanic(t *testing.T, wantPanic bool, args interface{}) {
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

func TestBuilderAddRemoveExistsNumField(t *testing.T) {
	t.Parallel()

	type args struct {
		builder Builder
	}
	tests := []struct {
		name               string
		args               args
		wantExistsIntField bool
		wantNumField       int
		wantPanic          bool
	}{
		{
			name:               "have fields set by newDynamicTestBuilder()",
			args:               args{builder: newDynamicTestBuilder()},
			wantExistsIntField: true,
			wantNumField:       12, // See: newDynamicTestBuilder()
		},
		{
			name:               "have fields set by newDynamicTestBuilder() and Remove(IntField)",
			args:               args{builder: newDynamicTestBuilder().Remove("IntField")},
			wantExistsIntField: false,
			wantNumField:       11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferDynamicTestPanic(t, tt.wantPanic, tt.args)

			if tt.args.builder.Exists("IntField") != tt.wantExistsIntField {
				t.Errorf("result Exists(IntField) is unexpected. got: %v, want: %v", tt.args.builder.Exists("IntField"), tt.wantExistsIntField)
				return
			}

			if tt.args.builder.NumField() != tt.wantNumField {
				t.Errorf("result numfield is unexpected. got: %d, want: %d", tt.args.builder.NumField(), tt.wantNumField)
				return
			}
		})
	}
}

func TestBuilderAddStringWithEmptyName(t *testing.T) {
	t.Parallel()

	type args struct {
		builder Builder
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name:      "try to AddString with empty name",
			args:      args{builder: newDynamicTestBuilder()},
			wantPanic: false, // FIXME: or NOT panic but error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferDynamicTestPanic(t, tt.wantPanic, tt.args)
		})
	}
}

func TestBuilderAddMapWithNilKey(t *testing.T) {
	t.Parallel()

	type args struct {
		builder Builder
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name:      "try to AddMap with nil key",
			args:      args{builder: newDynamicTestBuilder()},
			wantPanic: true, // expect to occur panic  FIXME: is error better than panic?
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferDynamicTestPanic(t, tt.wantPanic, tt.args)

			tt.args.builder.AddMap("MapFieldWithNilKey", nil, SampleFloat)
		})
	}
}

func TestBuilderAddMapWithNilValue(t *testing.T) {
	t.Parallel()

	type args struct {
		builder Builder
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name:      "try to AddMap with nil key",
			args:      args{builder: newDynamicTestBuilder()},
			wantPanic: true, // expect to occur panic  FIXME: is error better than panic?
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferDynamicTestPanic(t, tt.wantPanic, tt.args)

			tt.args.builder.AddMap("MapFieldWithNilKey", SampleString, nil)
		})
	}
}

func TestBuilderAddFuncWithNilArgs(t *testing.T) {
	t.Parallel()

	type args struct {
		builder Builder
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name:      "try to AddFunc with nil args",
			args:      args{builder: newDynamicTestBuilder()},
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferDynamicTestPanic(t, tt.wantPanic, tt.args)

			tt.args.builder.AddFunc("FuncFieldWithNilArgs", nil, []interface{}{SampleBool})
		})
	}
}

func TestBuilderAddFuncWithNilReturns(t *testing.T) {
	t.Parallel()

	type args struct {
		builder Builder
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name:      "try to AddFunc with nil returns",
			args:      args{builder: newDynamicTestBuilder()},
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferDynamicTestPanic(t, tt.wantPanic, tt.args)

			tt.args.builder.AddFunc("FuncFieldWithNilReturns", []interface{}{SampleInt}, nil)
		})
	}
}

func TestBuilderAddChanBothWithNilElem(t *testing.T) {
	t.Parallel()

	type args struct {
		builder Builder
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name:      "try to AddChanBoth with nil elem",
			args:      args{builder: newDynamicTestBuilder()},
			wantPanic: true, // expect to occur panic  FIXME: is error better than panic?
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferDynamicTestPanic(t, tt.wantPanic, tt.args)

			tt.args.builder.AddChanBoth("MapFieldWithNilKey", nil)
		})
	}
}

func TestBuilderAddChanRecvWithNilElem(t *testing.T) {
	t.Parallel()

	type args struct {
		builder Builder
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name:      "try to AddChanRecv with nil elem",
			args:      args{builder: newDynamicTestBuilder()},
			wantPanic: true, // expect to occur panic  FIXME: is error better than panic?
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferDynamicTestPanic(t, tt.wantPanic, tt.args)

			tt.args.builder.AddChanRecv("MapFieldWithNilKey", nil)
		})
	}
}

func TestBuilderAddChanSendWithNilElem(t *testing.T) {
	t.Parallel()

	type args struct {
		builder Builder
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name:      "try to AddChanSend with nil elem",
			args:      args{builder: newDynamicTestBuilder()},
			wantPanic: true, // expect to occur panic  FIXME: is error better than panic?
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferDynamicTestPanic(t, tt.wantPanic, tt.args)

			tt.args.builder.AddChanSend("MapFieldWithNilKey", nil)
		})
	}
}

func TestBuilderAddStructWithNil(t *testing.T) {
	t.Parallel()

	type args struct {
		builder Builder
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name:      "try to AddStruct with nil",
			args:      args{builder: newDynamicTestBuilder()},
			wantPanic: true, // expect to occur panic  FIXME: is error better than panic?
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferDynamicTestPanic(t, tt.wantPanic, tt.args)

			tt.args.builder.AddStruct("StructFieldWithNil", nil, false)
		})
	}
}

func TestBuilderAddStructPtrWithNil(t *testing.T) {
	t.Parallel()

	type args struct {
		builder Builder
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name:      "try to AddStructWith with nil",
			args:      args{builder: newDynamicTestBuilder()},
			wantPanic: true, // expect to occur panic  FIXME: is error better than panic?
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferDynamicTestPanic(t, tt.wantPanic, tt.args)

			tt.args.builder.AddStructPtr("StructPtrFieldWithNil", nil)
		})
	}
}

func TestBuilderAddSliceWithNil(t *testing.T) {
	t.Parallel()

	type args struct {
		builder Builder
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name:      "try to AddStructWith with nil",
			args:      args{builder: newDynamicTestBuilder()},
			wantPanic: true, // expect to occur panic  FIXME: is error better than panic?
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferDynamicTestPanic(t, tt.wantPanic, tt.args)

			tt.args.builder.AddSlice("SliceFieldWithNil", nil)
		})
	}
}

func TestBuilderBuild(t *testing.T) {
	t.Parallel()

	testMap := map[string]interface{}{
		"StringField": "ABCDEFGH",
		"IntField":    987,
		"FloatField":  1.23,
		"BoolField":   true,
		"MapField":    map[string]float64{"mfkey": 4.56},
		//"FuncField":   func(i1 int, i2 int) (bool, error) { return true, nil },  // FIXME
		"StructField": DynamicTestStruct{String: "Hoge"},
		"SliceField":  []*DynamicTestStruct{{String: "Huga1"}, {String: "Huga2"}},
	}

	type args struct {
		builder Builder
		isPtr   bool
	}
	tests := []struct {
		name               string
		args               args
		wantIsPtr          bool
		wantNumField       int
		testMap            map[string]interface{}
		wantErrorDecodeMap bool
		wantPanic          bool
	}{
		{
			name:         "Build() with valid Builder",
			args:         args{builder: newDynamicTestBuilder(), isPtr: true},
			wantIsPtr:    true,
			wantNumField: 12, // See: newDynamicTestBuilder()
			testMap:      testMap,
		},
		{
			name:               "BuildNonPtr() with valid Builder",
			args:               args{builder: newDynamicTestBuilder().Remove("ChanBothField"), isPtr: false},
			wantIsPtr:          false,
			wantNumField:       11,
			testMap:            testMap,
			wantErrorDecodeMap: true, // Note: can't execute DecodeMap if dynamic struct is NOT pointer.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferDynamicTestPanic(t, tt.wantPanic, tt.args)

			var got DynamicStruct
			if tt.args.isPtr {
				got = tt.args.builder.Build()
			} else {
				got = tt.args.builder.BuildNonPtr()
			}

			if got.IsPtr() != tt.wantIsPtr {
				t.Errorf("unexpected pointer or not result. got: %v, want: %v", got.IsPtr(), tt.wantIsPtr)
				return
			}

			if got.NumField() != tt.wantNumField {
				t.Errorf("result numfield is unexpected. got: %d, want: %d", got.NumField(), tt.wantNumField)
				return
			}

			if tt.testMap != nil {
				dec, err := got.DecodeMap(tt.testMap)
				if err != nil {
					if !tt.wantErrorDecodeMap {
						t.Errorf("unexpected error occured from DecodeMap: %v", err)
					}
					return
				} else if tt.wantErrorDecodeMap {
					t.Errorf("expected error did not occur from DecodeMap. dec: %+v", dec)
					return
				}

				getter, err := structil.NewGetter(dec)
				if err != nil {
					t.Errorf("unexpected error occured from NewGetter: %v", err)
					return
				}

				for k, v := range tt.testMap {
					gotValue, err := getter.EGet(k)
					if err != nil {
						t.Errorf("unexpected error occured from Getter.EGet: %v. name: %s", err, k)
						return
					}

					switch k {
					case "StructField":
						getter, err := structil.NewGetter(gotValue)
						if err != nil {
							t.Errorf("unexpected error occured from NewGetter for StructField: %v", err)
							return
						}

						ds, _ := v.(DynamicTestStruct)
						if getter.Get("String") != ds.String {
							t.Errorf("unexpected mismatch Struct String field: got: %v, want: %v", getter.Get("String"), ds.String)
							return
						}
					default:
						if d := cmp.Diff(gotValue, v); d != "" {
							t.Errorf("unexpected mismatch: name: %s, (-got +want)\n%s", k, d)
							return
						}
					}
				}
			}
		})
	}
}
