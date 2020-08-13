package dynamicstruct_test

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/goldeneggg/structil"

	"github.com/goldeneggg/structil/dynamicstruct"
	. "github.com/goldeneggg/structil/dynamicstruct"
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

const (
	stringFieldTag    = `json:"string_field_with_tag"`
	intFieldTag       = `json:"int_field_with_tag"`
	byteFieldTag      = `json:"byte_field_with_tag"`
	float32FieldTag   = `json:"float32_field_with_tag"`
	float64FieldTag   = `json:"float64_field_with_tag"`
	boolFieldTag      = `json:"bool_field_with_tag"`
	mapFieldTag       = `json:"map_field_with_tag"`
	funcFieldTag      = `json:"func_field_with_tag"`
	chanBothFieldTag  = `json:"chan_both_field_with_tag"`
	chanRecvFieldTag  = `json:"chan_recv_field_with_tag"`
	chanSendFieldTag  = `json:"chan_send_field_with_tag"`
	structFieldTag    = `json:"struct_field_with_tag"`
	structPtrFieldTag = `json:"struct_ptr_field_with_tag"`
	sliceFieldTag     = `json:"slice_field_with_tag"`
	interfaceFieldTag = `json:"interface_field_with_tag"`
)

var (
	dynamicTestString2 = "test name2"
	dynamicTestFunc    = func(s string) interface{} { return s + "-func" }
	//dynamicTestChan    = make(chan int)
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
		AddStringWithTag("StringFieldWithTag", stringFieldTag).
		AddInt("IntField").
		AddIntWithTag("IntFieldWithTag", intFieldTag).
		AddByte("ByteField").
		AddByteWithTag("ByteFieldWithTag", byteFieldTag).
		AddFloat32("Float32Field").
		AddFloat32WithTag("Float32FieldWithTag", float32FieldTag).
		AddFloat64("Float64Field").
		AddFloat64WithTag("Float64FieldWithTag", float64FieldTag).
		AddBool("BoolField").
		AddBoolWithTag("BoolFieldWithTag", boolFieldTag).
		AddMap("MapField", SampleString, SampleFloat32).
		AddMapWithTag("MapFieldWithTag", SampleString, SampleFloat32, mapFieldTag).
		AddFunc("FuncField", []interface{}{SampleInt, SampleInt}, []interface{}{SampleBool, ErrSample}).
		AddFuncWithTag("FuncFieldWithTag", []interface{}{SampleInt, SampleInt}, []interface{}{SampleBool, ErrSample}, funcFieldTag).
		AddChanBoth("ChanBothField", SampleInt).
		AddChanBothWithTag("ChanBothFieldWithTag", SampleInt, chanBothFieldTag).
		AddChanRecv("ChanRecvField", SampleInt).
		AddChanRecvWithTag("ChanRecvFieldWithTag", SampleInt, chanRecvFieldTag).
		AddChanSend("ChanSendField", SampleInt).
		AddChanSendWithTag("ChanSendFieldWithTag", SampleInt, chanSendFieldTag).
		AddStruct("StructField", newDynamicTestStruct(), false).
		AddStructWithTag("StructFieldWithTag", newDynamicTestStruct(), false, structFieldTag).
		AddStructPtr("StructPtrField", newDynamicTestStructPtr()).
		AddStructPtrWithTag("StructPtrFieldWithTag", newDynamicTestStructPtr(), structPtrFieldTag).
		AddSlice("SliceField", newDynamicTestStructPtr()).
		AddSliceWithTag("SliceFieldWithTag", newDynamicTestStructPtr(), sliceFieldTag).
		AddInterface("InterfaceField", false).
		AddInterfaceWithTag("InterfaceFieldWithTag", false, interfaceFieldTag).
		AddInterface("InterfacePtrField", true).
		AddInterfaceWithTag("InterfacePtrFieldWithTag", true, interfaceFieldTag)
}

func newDynamicTestBuilderWithStructName(name string) Builder {
	b := newDynamicTestBuilder()
	b.SetStructName(name)
	return b
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
		wantStructName     string
		wantPanic          bool
	}{
		{
			name:               "have fields set by newDynamicTestBuilder()",
			args:               args{builder: newDynamicTestBuilder()},
			wantExistsIntField: true,
			wantNumField:       32, // See: newDynamicTestBuilder()
			wantStructName:     "DynamicStruct",
		},
		{
			name:               "have fields set by newDynamicTestBuilder() and Remove(IntField)",
			args:               args{builder: newDynamicTestBuilder().Remove("IntField")},
			wantExistsIntField: false,
			wantNumField:       31,
			wantStructName:     "DynamicStruct",
		},
		{
			name:               "have struct name by newDynamicTestBuilderWithStructName()",
			args:               args{builder: newDynamicTestBuilderWithStructName("Abc")},
			wantExistsIntField: true,
			wantNumField:       32,
			wantStructName:     "Abc",
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

			if tt.args.builder.GetStructName() != tt.wantStructName {
				t.Errorf("result structName is unexpected. got: %s, want: %s", tt.args.builder.GetStructName(), tt.wantStructName)
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

			tt.args.builder.AddString("")
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

			tt.args.builder.AddMap("MapFieldWithNilKey", nil, SampleFloat32)
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

type buildArgs struct {
	builder Builder
	isPtr   bool
}

type buildTest struct {
	name               string
	args               buildArgs
	wantIsPtr          bool
	wantStructName     string
	wantNumField       int
	wantDefinition     string
	testMap            map[string]interface{}
	wantErrorDecodeMap bool
	wantPanic          bool
}

func TestBuilderBuild(t *testing.T) {
	t.Parallel()

	testMap := map[string]interface{}{
		"StringField":  "ABCDEFGH",
		"IntField":     987,
		"ByteField":    byte(1),
		"Float32Field": float32(1.23),
		"Float64Field": float64(2.3),
		"BoolField":    true,
		"MapField":     map[string]float32{"mfkey": float32(4.56)},
		//"FuncField":   func(i1 int, i2 int) (bool, error) { return true, nil },  // FIXME
		"StructField": DynamicTestStruct{String: "Hoge"},
		"SliceField":  []*DynamicTestStruct{{String: "Huga1"}, {String: "Huga2"}},
	}

	tests := []buildTest{
		{
			name:           "Build() with valid Builder",
			args:           buildArgs{builder: newDynamicTestBuilder(), isPtr: true},
			wantIsPtr:      true,
			wantStructName: "DynamicStruct",
			wantNumField:   32, // See: newDynamicTestBuilder()
			wantDefinition: `type DynamicStruct struct {
	BoolField bool
	BoolFieldWithTag bool ` + "`json:\"bool_field_with_tag\"`" + `
	ByteField uint8
	ByteFieldWithTag uint8 ` + "`json:\"byte_field_with_tag\"`" + `
	ChanBothField chan int
	ChanBothFieldWithTag chan int ` + "`json:\"chan_both_field_with_tag\"`" + `
	ChanRecvField <-chan int
	ChanRecvFieldWithTag <-chan int ` + "`json:\"chan_recv_field_with_tag\"`" + `
	ChanSendField chan<- int
	ChanSendFieldWithTag chan<- int ` + "`json:\"chan_send_field_with_tag\"`" + `
	Float32Field float32
	Float32FieldWithTag float32 ` + "`json:\"float32_field_with_tag\"`" + `
	Float64Field float64
	Float64FieldWithTag float64 ` + "`json:\"float64_field_with_tag\"`" + `
	FuncField func(int, int) (bool, *errors.errorString)
	FuncFieldWithTag func(int, int) (bool, *errors.errorString) ` + "`json:\"func_field_with_tag\"`" + `
	IntField int
	IntFieldWithTag int ` + "`json:\"int_field_with_tag\"`" + `
	InterfaceField interface {}
	InterfaceFieldWithTag interface {} ` + "`json:\"interface_field_with_tag\"`" + `
	InterfacePtrField *interface {}
	InterfacePtrFieldWithTag *interface {} ` + "`json:\"interface_field_with_tag\"`" + `
	MapField map[string]float32
	MapFieldWithTag map[string]float32 ` + "`json:\"map_field_with_tag\"`" + `
	SliceField []*dynamicstruct_test.DynamicTestStruct
	SliceFieldWithTag []*dynamicstruct_test.DynamicTestStruct ` + "`json:\"slice_field_with_tag\"`" + `
	StringField string
	StringFieldWithTag string ` + "`json:\"string_field_with_tag\"`" + `
	StructField struct { Byte uint8; Bytes []uint8; Int int; Int64 int64; Uint uint; Uint64 uint64; Float32 float32; Float64 float64; String string; Stringptr *string; Stringslice []string; Bool bool; Map map[string]interface {}; Func func(string) interface {}; DynamicTestStruct2 dynamicstruct_test.DynamicTestStruct2; DynamicTestStruct2Ptr *dynamicstruct_test.DynamicTestStruct2; DynamicTestStruct4Slice []dynamicstruct_test.DynamicTestStruct4; DynamicTestStruct4PtrSlice []*dynamicstruct_test.DynamicTestStruct4 }
	StructFieldWithTag struct { Byte uint8; Bytes []uint8; Int int; Int64 int64; Uint uint; Uint64 uint64; Float32 float32; Float64 float64; String string; Stringptr *string; Stringslice []string; Bool bool; Map map[string]interface {}; Func func(string) interface {}; DynamicTestStruct2 dynamicstruct_test.DynamicTestStruct2; DynamicTestStruct2Ptr *dynamicstruct_test.DynamicTestStruct2; DynamicTestStruct4Slice []dynamicstruct_test.DynamicTestStruct4; DynamicTestStruct4PtrSlice []*dynamicstruct_test.DynamicTestStruct4 } ` + "`json:\"struct_field_with_tag\"`" + `
	StructPtrField *struct { Byte uint8; Bytes []uint8; Int int; Int64 int64; Uint uint; Uint64 uint64; Float32 float32; Float64 float64; String string; Stringptr *string; Stringslice []string; Bool bool; Map map[string]interface {}; Func func(string) interface {}; DynamicTestStruct2 dynamicstruct_test.DynamicTestStruct2; DynamicTestStruct2Ptr *dynamicstruct_test.DynamicTestStruct2; DynamicTestStruct4Slice []dynamicstruct_test.DynamicTestStruct4; DynamicTestStruct4PtrSlice []*dynamicstruct_test.DynamicTestStruct4 }
	StructPtrFieldWithTag *struct { Byte uint8; Bytes []uint8; Int int; Int64 int64; Uint uint; Uint64 uint64; Float32 float32; Float64 float64; String string; Stringptr *string; Stringslice []string; Bool bool; Map map[string]interface {}; Func func(string) interface {}; DynamicTestStruct2 dynamicstruct_test.DynamicTestStruct2; DynamicTestStruct2Ptr *dynamicstruct_test.DynamicTestStruct2; DynamicTestStruct4Slice []dynamicstruct_test.DynamicTestStruct4; DynamicTestStruct4PtrSlice []*dynamicstruct_test.DynamicTestStruct4 } ` + "`json:\"struct_ptr_field_with_tag\"`" + `
}`,
			testMap: testMap,
		},
		{
			name:               "BuildNonPtr() with valid Builder",
			args:               buildArgs{builder: newDynamicTestBuilder(), isPtr: false},
			wantIsPtr:          false,
			wantStructName:     "DynamicStruct",
			wantNumField:       32,
			testMap:            testMap,
			wantErrorDecodeMap: true, // Note: can't execute DecodeMap if dynamic struct is NOT pointer.
		},
		{
			name:           "Build() with valid Builder with struct name",
			args:           buildArgs{builder: newDynamicTestBuilderWithStructName("HogeHuga"), isPtr: true},
			wantIsPtr:      true,
			wantStructName: "HogeHuga",
			wantNumField:   32, // See: newDynamicTestBuilder()
			testMap:        testMap,
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

			if !testBuilderBuildWant(t, got, tt) {
				return
			}

			if !testBuilderBuildTag(t, got, tt) {
				return
			}

			if tt.testMap != nil {
				if !testBuilderBuildDecodeMap(t, got, tt) {
					return
				}
			}
		})
	}
}

func testBuilderBuildWant(t *testing.T, got DynamicStruct, tt buildTest) bool {
	if got.IsPtr() != tt.wantIsPtr {
		t.Errorf("unexpected pointer or not result. got: %v, want: %v", got.IsPtr(), tt.wantIsPtr)
		return false
	}

	if got.Name() != tt.wantStructName {
		t.Errorf("result struct name is unexpected. got: %s, want: %s", got.Name(), tt.wantStructName)
		return false
	}

	if got.NumField() != tt.wantNumField {
		t.Errorf("result numfield is unexpected. got: %d, want: %d", got.NumField(), tt.wantNumField)
		return false
	}

	if tt.wantDefinition != "" {
		if d := cmp.Diff(got.Definition(), tt.wantDefinition); d != "" {
			//t.Errorf("unexpected mismatch Definition. got:\n%s\n, want:\n%s\n", got.Definition(), tt.wantDefinition)
			t.Errorf("unexpected mismatch Definition: (-got +want)\n%s", d)
			return false
		}
	}

	return true
}

func testBuilderBuildTag(t *testing.T, got DynamicStruct, tt buildTest) bool {
	prefixes := map[string]string{
		"String":    stringFieldTag,
		"Int":       intFieldTag,
		"Byte":      byteFieldTag,
		"Float32":   float32FieldTag,
		"Float64":   float64FieldTag,
		"Bool":      boolFieldTag,
		"Map":       mapFieldTag,
		"Func":      funcFieldTag,
		"ChanBoth":  chanBothFieldTag,
		"ChanRecv":  chanRecvFieldTag,
		"ChanSend":  chanSendFieldTag,
		"Struct":    structFieldTag,
		"StructPtr": structPtrFieldTag,
		"Slice":     sliceFieldTag,
	}

	var fName string
	for prefix, tagWithTag := range prefixes {
		// test without tag fields
		fName = prefix + "Field"
		st, ok := got.FieldByName(fName)
		if ok {
			if d := cmp.Diff(st.Tag, reflect.StructTag("")); d != "" {
				t.Errorf("unexpected mismatch Tag: fName: %s, (-got +want)\n%s", fName, d)
				return false
			}
		} else {
			t.Errorf("Field %s does not exist.", fName)
			return false
		}

		// test with tag fields
		fName = prefix + "FieldWithTag"
		sft, ok := got.FieldByName(fName)
		if ok {
			if d := cmp.Diff(sft.Tag, reflect.StructTag(tagWithTag)); d != "" {
				t.Errorf("unexpected mismatch WithTag.Tag: fName: %s, (-got +want)\n%s", fName, d)
				return false
			}
		} else {
			t.Errorf("Field %s does not exist.", fName)
			return false
		}
	}

	return true
}

func testBuilderBuildDecodeMap(t *testing.T, got DynamicStruct, tt buildTest) bool {
	dec, err := got.DecodeMap(tt.testMap)
	if err != nil {
		if !tt.wantErrorDecodeMap {
			t.Errorf("unexpected error occured from DecodeMap: %v", err)
		}
		return false
	} else if tt.wantErrorDecodeMap {
		t.Errorf("expected error did not occur from DecodeMap. dec: %+v", dec)
		return false
	}

	getter, err := structil.NewGetter(dec)
	if err != nil {
		t.Errorf("unexpected error occured from NewGetter: %v", err)
		return false
	}

	for k, v := range tt.testMap {
		gotValue, err := getter.EGet(k)
		if err != nil {
			t.Errorf("unexpected error occured from Getter.EGet: %v. name: %s", err, k)
			return false
		}

		switch k {
		case "StructField":
			getter, err := structil.NewGetter(gotValue)
			if err != nil {
				t.Errorf("unexpected error occured from NewGetter for StructField: %v", err)
				return false
			}

			ds, _ := v.(DynamicTestStruct)
			if getter.Get("String") != ds.String {
				t.Errorf("unexpected mismatch Struct String field: got: %v, want: %v", getter.Get("String"), ds.String)
				return false
			}
		default:
			if d := cmp.Diff(gotValue, v); d != "" {
				t.Errorf("unexpected mismatch: name: %s, (-got +want)\n%s", k, d)
				return false
			}
		}
	}

	return true
}

func TestJSONToDynamicStructInterface(t *testing.T) {
	t.Parallel()

	type args struct {
		jsonData []byte
	}
	tests := []struct {
		name           string
		args           args
		wantError      bool
		numField       int
		hasStringField bool
	}{
		{
			name: "JSON does not have null field",
			args: args{
				jsonData: []byte(`
{
	"null_field":null,
	"string_field":"かきくけこ",
	"int_field":45678,
	"float32_field":9.876,
	"bool_field":false,
	"struct_ptr_field":{
		"key":"hugakey",
		"value":"hugavalue"
	},
	"array_string_field":[
		"array_str_1",
		"array_str_2"
	],
	"array_struct_field":[
		{
			"kkk":"kkk1",
			"vvvv":"vvv1"
		},
		{
			"kkk":"kkk2",
			"vvvv":"vvv2"
		},
		{
			"kkk":"kkk3",
			"vvvv":"vvv3"
		}
	]
}
`),
			},
			wantError:      false,
			numField:       7,
			hasStringField: true,
		},
		{
			name: "JSON is valid array",
			args: args{
				jsonData: []byte(`
		[
			{
				"null_field":null,
				"string_field":"かきくけこ",
				"int_field":45678,
				"float32_field":9.876,
				"bool_field":false,
				"struct_ptr_field":{
					"key":"hugakey",
					"value":"hugavalue"
				},
				"array_string_field":[
					"array_str_1",
					"array_str_2"
				],
				"array_struct_field":[
					{
						"kkk":"kkk1",
						"vvvv":"vvv1"
					},
					{
						"kkk":"kkk2",
						"vvvv":"vvv2"
					},
					{
						"kkk":"kkk3",
						"vvvv":"vvv3"
					}
				]
			},
			{
				"null_field":null,
				"string_field":"さしすせそ",
				"int_field":7890,
				"float32_field":4.99,
				"bool_field":true,
				"struct_ptr_field":{
					"key":"hugakeyXXX",
					"value":"hugavalueXXX"
				},
				"array_string_field":[
					"array_str_111",
					"array_str_222"
				],
				"array_struct_field":[
					{
						"kkk":"kkk99",
						"vvvv":"vvv99"
					},
					{
						"kkk":"kkk999",
						"vvvv":"vvv999"
					},
					{
						"kkk":"kkk9999",
						"vvvv":"vvv9999"
					}
				]
			}
		]
		`),
			},
			wantError:      false,
			numField:       1,
			hasStringField: false,
		},
		{
			name: "Only one null field",
			args: args{
				jsonData: []byte(`{"nullfield":null}`),
			},
			wantError:      false,
			numField:       0,
			hasStringField: false,
		},
		{
			name: "Empty JSON",
			args: args{
				jsonData: []byte(`{}`),
			},
			wantError:      false,
			numField:       0,
			hasStringField: false,
		},
		{
			name: "Empty array JSON",
			args: args{
				jsonData: []byte(`[]`),
			},
			wantError:      false,
			numField:       0,
			hasStringField: false,
		},
		{
			name: "empty",
			args: args{
				jsonData: []byte(``),
			},
			wantError:      true,
			numField:       0,
			hasStringField: false,
		},
		{
			name: "null",
			args: args{
				jsonData: []byte(`null`),
			},
			wantError:      true,
			numField:       0,
			hasStringField: false,
		},
		{
			name: "Invalid string",
			args: args{
				jsonData: []byte(`invalid`),
			},
			wantError:      true,
			numField:       0,
			hasStringField: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			intf, err := dynamicstruct.JSONToDynamicStructInterface(tt.args.jsonData)
			if err == nil {
				if tt.wantError {
					t.Errorf("error did not occur. intf: %#v", intf)
					return
				}

				// FIXME: other test method (e.g. json.Marshal(intf))
				/*
					g, err := structil.NewGetter(intf)
					if err != nil {
						t.Errorf("unexpected error occured in structil.NewGetter: %v", err)
					}

					if g.Has("StringField") != tt.hasStringField {
						t.Errorf("unexpected result of Has StringField. got: %v, want: %v", g.Has("StringField"), tt.hasStringField)
					}
				*/

			} else if !tt.wantError {
				t.Errorf("unexpected error occured. wantError %v, err: %v", tt.wantError, err)
			}
		})
	}
}

// benchmark tests

func BenchmarkAddString(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddString("StringField")
	}
}

func BenchmarkAddStringWithTag(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddStringWithTag("StringField", stringFieldTag)
	}
}

func BenchmarkAddInt(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddInt("IntField")
	}
}

func BenchmarkAddFloa32(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddFloat32("Float32Field")
	}
}

func BenchmarkAddFloa64(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddFloat64("Float64Field")
	}
}

func BenchmarkAddBool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddBool("BoolField")
	}

}

func BenchmarkAddMap(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddMap("MapField", SampleString, SampleFloat32)
	}
}

func BenchmarkAddFunc(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddFunc("FuncField", []interface{}{SampleInt, SampleInt}, []interface{}{SampleBool, ErrSample})
	}
}

func BenchmarkAddChanBoth(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddChanBoth("ChanBothField", SampleInt)
	}
}

func BenchmarkAddChanRecv(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddChanRecv("ChanRecvField", SampleInt)
	}
}

func BenchmarkAddChanSend(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddChanSend("ChanSendField", SampleInt)
	}
}

func BenchmarkAddStruct(b *testing.B) {
	st := newDynamicTestStruct()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddStruct("StructField", st, false)
	}
}

func BenchmarkAddStructPtr(b *testing.B) {
	st := newDynamicTestStructPtr()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddStructPtr("StructPtrField", st)
	}
}

func BenchmarkAddSlice(b *testing.B) {
	st := newDynamicTestStructPtr()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddSlice("SliceField", st)
	}
}

func BenchmarkAddInterface(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddInterface("InterfaceField", false)
	}
}

func BenchmarkAddInterfacePtr(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddInterface("InterfacePtrField", true)
	}
}

func BenchmarkBuild(b *testing.B) {
	builder := newDynamicTestBuilder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.Build()
	}
}

func BenchmarkBuildNonPtr(b *testing.B) {
	builder := newDynamicTestBuilder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.BuildNonPtr()
	}
}

func BenchmarkDefinition(b *testing.B) {
	builder := newDynamicTestBuilder()
	ds := builder.Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ds.Definition()
	}
}
