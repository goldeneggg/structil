// FIXME: Remove unnessesary table-driven tests (and simplifize tests)

package dynamicstruct_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

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

const expectedDefinition = `type DynamicStruct struct {
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
	SliceField []struct {
		Bool bool
		Byte uint8
		Bytes []uint8
		DynamicTestStruct2 struct {
			DynamicTestStruct3 struct {
				Int int
				String string
			}
			String string
		}
		DynamicTestStruct2Ptr struct {
			DynamicTestStruct3 struct {
				Int int
				String string
			}
			String string
		}
		DynamicTestStruct4PtrSlice []struct {
			String string
			String2 string
		}
		DynamicTestStruct4Slice []struct {
			String string
			String2 string
		}
		Float32 float32
		Float64 float64
		Func func(string) interface {}
		Int int
		Int64 int64
		Map map[string]interface {}
		String string
		Stringptr *string
		Stringslice []string
		Uint uint
		Uint64 uint64
	}
	SliceFieldWithTag []struct {
		Bool bool
		Byte uint8
		Bytes []uint8
		DynamicTestStruct2 struct {
			DynamicTestStruct3 struct {
				Int int
				String string
			}
			String string
		}
		DynamicTestStruct2Ptr struct {
			DynamicTestStruct3 struct {
				Int int
				String string
			}
			String string
		}
		DynamicTestStruct4PtrSlice []struct {
			String string
			String2 string
		}
		DynamicTestStruct4Slice []struct {
			String string
			String2 string
		}
		Float32 float32
		Float64 float64
		Func func(string) interface {}
		Int int
		Int64 int64
		Map map[string]interface {}
		String string
		Stringptr *string
		Stringslice []string
		Uint uint
		Uint64 uint64
	} ` + "`json:\"slice_field_with_tag\"`" + `
	StringField string
	StringFieldWithTag string ` + "`json:\"string_field_with_tag\"`" + `
	StructField struct {
		Bool bool
		Byte uint8
		Bytes []uint8
		DynamicTestStruct2 struct {
			DynamicTestStruct3 struct {
				Int int
				String string
			}
			String string
		}
		DynamicTestStruct2Ptr struct {
			DynamicTestStruct3 struct {
				Int int
				String string
			}
			String string
		}
		DynamicTestStruct4PtrSlice []struct {
			String string
			String2 string
		}
		DynamicTestStruct4Slice []struct {
			String string
			String2 string
		}
		Float32 float32
		Float64 float64
		Func func(string) interface {}
		Int int
		Int64 int64
		Map map[string]interface {}
		String string
		Stringptr *string
		Stringslice []string
		Uint uint
		Uint64 uint64
	}
	StructFieldWithTag struct {
		Bool bool
		Byte uint8
		Bytes []uint8
		DynamicTestStruct2 struct {
			DynamicTestStruct3 struct {
				Int int
				String string
			}
			String string
		}
		DynamicTestStruct2Ptr struct {
			DynamicTestStruct3 struct {
				Int int
				String string
			}
			String string
		}
		DynamicTestStruct4PtrSlice []struct {
			String string
			String2 string
		}
		DynamicTestStruct4Slice []struct {
			String string
			String2 string
		}
		Float32 float32
		Float64 float64
		Func func(string) interface {}
		Int int
		Int64 int64
		Map map[string]interface {}
		String string
		Stringptr *string
		Stringslice []string
		Uint uint
		Uint64 uint64
	} ` + "`json:\"struct_field_with_tag\"`" + `
	StructPtrField struct {
		Bool bool
		Byte uint8
		Bytes []uint8
		DynamicTestStruct2 struct {
			DynamicTestStruct3 struct {
				Int int
				String string
			}
			String string
		}
		DynamicTestStruct2Ptr struct {
			DynamicTestStruct3 struct {
				Int int
				String string
			}
			String string
		}
		DynamicTestStruct4PtrSlice []struct {
			String string
			String2 string
		}
		DynamicTestStruct4Slice []struct {
			String string
			String2 string
		}
		Float32 float32
		Float64 float64
		Func func(string) interface {}
		Int int
		Int64 int64
		Map map[string]interface {}
		String string
		Stringptr *string
		Stringslice []string
		Uint uint
		Uint64 uint64
	}
	StructPtrFieldWithTag struct {
		Bool bool
		Byte uint8
		Bytes []uint8
		DynamicTestStruct2 struct {
			DynamicTestStruct3 struct {
				Int int
				String string
			}
			String string
		}
		DynamicTestStruct2Ptr struct {
			DynamicTestStruct3 struct {
				Int int
				String string
			}
			String string
		}
		DynamicTestStruct4PtrSlice []struct {
			String string
			String2 string
		}
		DynamicTestStruct4Slice []struct {
			String string
			String2 string
		}
		Float32 float32
		Float64 float64
		Func func(string) interface {}
		Int int
		Int64 int64
		Map map[string]interface {}
		String string
		Stringptr *string
		Stringslice []string
		Uint uint
		Uint64 uint64
	} ` + "`json:\"struct_ptr_field_with_tag\"`" + `
}`

var (
	dynamicTestString2 = "test name2"
	dynamicTestFunc    = func(s string) interface{} { return s + "-func" }
	//dynamicTestChan    = make(chan int)
)

func newTestDynamicStruct() DynamicTestStruct {
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

func newTestDynamicStructPtr() *DynamicTestStruct {
	ts := newTestDynamicStruct()
	return &ts
}

// See: "expectedDefinition" constant (this is the Definition of the Builder as follows)
func newTestBuilder() *Builder {
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
		AddStruct("StructField", newTestDynamicStruct(), false).
		AddStructWithTag("StructFieldWithTag", newTestDynamicStruct(), false, structFieldTag).
		AddStructPtr("StructPtrField", newTestDynamicStructPtr()).
		AddStructPtrWithTag("StructPtrFieldWithTag", newTestDynamicStructPtr(), structPtrFieldTag).
		AddSlice("SliceField", newTestDynamicStructPtr()).
		AddSliceWithTag("SliceFieldWithTag", newTestDynamicStructPtr(), sliceFieldTag).
		AddInterface("InterfaceField", false).
		AddInterfaceWithTag("InterfaceFieldWithTag", false, interfaceFieldTag).
		AddInterface("InterfacePtrField", true).
		AddInterfaceWithTag("InterfacePtrFieldWithTag", true, interfaceFieldTag)
}

func newTestBuilderWithStructName(name string) *Builder {
	b := newTestBuilder()
	b.SetStructName(name)
	return b
}

func TestBuilderAddRemoveExistsNumField(t *testing.T) {
	t.Parallel()

	type args struct {
		builder *Builder
	}
	tests := []struct {
		name               string
		args               args
		wantExistsIntField bool
		wantNumField       int
		wantStructName     string
	}{
		{
			name:               "have fields set by newTestBuilder()",
			args:               args{builder: newTestBuilder()},
			wantExistsIntField: true,
			wantNumField:       32, // See: newTestBuilder()
			wantStructName:     "DynamicStruct",
		},
		{
			name:               "have fields set by newTestBuilder() and Remove(IntField)",
			args:               args{builder: newTestBuilder().Remove("IntField")},
			wantExistsIntField: false,
			wantNumField:       31,
			wantStructName:     "DynamicStruct",
		},
		{
			name:               "have struct name by newTestBuilderWithStructName()",
			args:               args{builder: newTestBuilderWithStructName("Abc")},
			wantExistsIntField: true,
			wantNumField:       32,
			wantStructName:     "Abc",
		},
	}

	for _, tt := range tests {
		tt := tt // See: https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

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
		builder *Builder
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "try to AddString with empty name",
			args: args{builder: newTestBuilder()},
		},
	}

	for _, tt := range tests {
		tt := tt // See: https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.args.builder.AddString("").Build()
			if err == nil {
				t.Errorf("expect to occur error but does not: args: %+v", tt.args)
			}
		})
	}
}

func TestBuilderAddMapWithNilKey(t *testing.T) {
	t.Parallel()

	type args struct {
		builder *Builder
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "try to AddMap with nil key",
			args: args{builder: newTestBuilder()},
		},
	}

	for _, tt := range tests {
		tt := tt // See: https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.args.builder.AddMap("MapFieldWithNilKey", nil, SampleFloat32).Build()
			if err == nil {
				t.Errorf("expect to occur error but does not: args: %+v", tt.args)
			}
		})
	}
}

func TestBuilderAddMapWithNilValue(t *testing.T) {
	t.Parallel()

	type args struct {
		builder *Builder
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "try to AddMap with nil key",
			args: args{builder: newTestBuilder()},
		},
	}

	for _, tt := range tests {
		tt := tt // See: https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.args.builder.AddMap("MapFieldWithNilKey", SampleString, nil).Build()
			// nil map value does NOT cause error
			if err != nil {
				t.Errorf("unexpected error occured %v: args: %+v", err, tt.args)
			}
		})
	}
}

func TestBuilderAddFuncWithNilArgs(t *testing.T) {
	t.Parallel()

	type args struct {
		builder *Builder
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "try to AddFunc with nil args",
			args: args{builder: newTestBuilder()},
		},
	}

	for _, tt := range tests {
		tt := tt // See: https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.args.builder.AddFunc("FuncFieldWithNilArgs", nil, []interface{}{SampleBool}).Build()
			if err != nil {
				t.Errorf("unexpected error occurred: args: %+v, %v", tt.args, err)
			}
		})
	}
}

func TestBuilderAddFuncWithNilReturns(t *testing.T) {
	t.Parallel()

	type args struct {
		builder *Builder
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "try to AddFunc with nil returns",
			args: args{builder: newTestBuilder()},
		},
	}

	for _, tt := range tests {
		tt := tt // See: https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.args.builder.AddFunc("FuncFieldWithNilReturns", []interface{}{SampleInt}, nil).Build()
			if err != nil {
				t.Errorf("unexpected error occurred: args: %+v, %v", tt.args, err)
			}
		})
	}
}

func TestBuilderAddChanBothWithNilElem(t *testing.T) {
	t.Parallel()

	type args struct {
		builder *Builder
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "try to AddChanBoth with nil elem",
			args: args{builder: newTestBuilder()},
		},
	}

	for _, tt := range tests {
		tt := tt // See: https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.args.builder.AddChanBoth("MapFieldWithNilKey", nil).Build()
			if err == nil {
				t.Errorf("expect to occur error but does not: args: %+v", tt.args)
			}
		})
	}
}

func TestBuilderAddChanRecvWithNilElem(t *testing.T) {
	t.Parallel()

	type args struct {
		builder *Builder
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
	}{
		{
			name:      "try to AddChanRecv with nil elem",
			args:      args{builder: newTestBuilder()},
			wantError: true,
		},
	}

	for _, tt := range tests {
		tt := tt // See: https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.args.builder.AddChanRecv("MapFieldWithNilKey", nil).Build()
			if err == nil {
				t.Errorf("expect to occur error but does not: args: %+v", tt.args)
			}
		})
	}
}

func TestBuilderAddChanSendWithNilElem(t *testing.T) {
	t.Parallel()

	type args struct {
		builder *Builder
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "try to AddChanSend with nil elem",
			args: args{builder: newTestBuilder()},
		},
	}

	for _, tt := range tests {
		tt := tt // See: https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.args.builder.AddChanSend("MapFieldWithNilKey", nil).Build()
			if err == nil {
				t.Errorf("expect to occur error but does not: args: %+v", tt.args)
			}
		})
	}
}

func TestBuilderAddStructWithNil(t *testing.T) {
	t.Parallel()

	type args struct {
		builder *Builder
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "try to AddStruct with nil",
			args: args{builder: newTestBuilder()},
		},
	}

	for _, tt := range tests {
		tt := tt // See: https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.args.builder.AddStruct("StructFieldWithNil", nil, false).Build()
			if err == nil {
				t.Errorf("expect to occur error but does not: args: %+v", tt.args)
			}
		})
	}
}

func TestBuilderAddStructPtrWithNil(t *testing.T) {
	t.Parallel()

	type args struct {
		builder *Builder
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "try to AddStructWith with nil",
			args: args{builder: newTestBuilder()},
		},
	}

	for _, tt := range tests {
		tt := tt // See: https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.args.builder.AddStructPtr("StructPtrFieldWithNil", nil).Build()
			if err == nil {
				t.Errorf("expect to occur error but does not: args: %+v", tt.args)
			}
		})
	}
}

func TestBuilderAddSliceWithNil(t *testing.T) {
	t.Parallel()

	type args struct {
		builder *Builder
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "try to AddStructWith with nil",
			args: args{builder: newTestBuilder()},
		},
	}

	for _, tt := range tests {
		tt := tt // See: https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.args.builder.AddSlice("SliceFieldWithNil", nil).Build()
			if err == nil {
				t.Errorf("expect to occur error but does not: args: %+v", tt.args)
			}
		})
	}
}

type buildArgs struct {
	builder *Builder
	isPtr   bool
}

type buildTest struct {
	name                string
	args                buildArgs
	wantIsPtr           bool
	wantStructName      string
	wantNumField        int
	wantDefinition      string
	camelizeKeys        bool
	tryAddDynamicStruct bool
}

func TestBuilderBuild(t *testing.T) {
	t.Parallel()

	tests := []buildTest{
		{
			name:                "Build() with valid Builder",
			args:                buildArgs{builder: newTestBuilder(), isPtr: true},
			wantIsPtr:           true,
			wantStructName:      "DynamicStruct",
			wantNumField:        32, // See: newTestBuilder()
			wantDefinition:      expectedDefinition,
			camelizeKeys:        true,
			tryAddDynamicStruct: true,
		},
		{
			name:                "Build() with valid Builder",
			args:                buildArgs{builder: newTestBuilder(), isPtr: true},
			wantIsPtr:           true,
			wantStructName:      "DynamicStruct",
			wantNumField:        32, // See: newTestBuilder()
			wantDefinition:      expectedDefinition,
			tryAddDynamicStruct: true,
		},
		{
			name:           "BuildNonPtr() with valid Builder",
			args:           buildArgs{builder: newTestBuilder(), isPtr: false},
			wantIsPtr:      false,
			wantStructName: "DynamicStruct",
			wantNumField:   32,
		},
		{
			name:           "Build() with valid Builder with struct name",
			args:           buildArgs{builder: newTestBuilderWithStructName("HogeHuga"), isPtr: true},
			wantIsPtr:      true,
			wantStructName: "HogeHuga",
			wantNumField:   32, // See: newTestBuilder()
			camelizeKeys:   true,
		},
	}

	var ds *DynamicStruct
	var err error

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			// FIXME: comment out t.Parallel() because of race condition in Build() method
			t.Parallel()

			if tt.args.isPtr {
				ds, err = tt.args.builder.Build()
			} else {
				ds, err = tt.args.builder.BuildNonPtr()
			}
			if err != nil {
				t.Errorf("unexpected error caused by DynamicStruct Build: %v", err)
			}

			if !testBuilderBuildWant(t, ds, tt) {
				return
			}

			if !testBuilderBuildTag(t, ds, tt) {
				return
			}

			if tt.tryAddDynamicStruct {
				if !testBuilderBuildAddDynamicStruct(t, ds, tt) {
					return
				}
			}
		})
	}
}

func testBuilderBuildWant(t *testing.T, ds *DynamicStruct, tt buildTest) bool {
	t.Helper()

	if ds.Name() != tt.wantStructName {
		t.Fatalf("result struct name is unexpected. got: %s, want: %s", ds.Name(), tt.wantStructName)
	}

	k := ds.Type().Kind()
	if k != reflect.Struct {
		t.Fatalf("result struct Type.Kind is unexpected. got: %s, want: Struct", k)
	}

	flds := ds.Fields()
	if len(flds) != tt.wantNumField {
		t.Fatalf("result Fields's length is unexpected. got: %d, want: %d", len(flds), tt.wantNumField)
	}

	if len(flds) > 0 {
		f := ds.Field(0)
		if flds[0].Name != f.Name {
			t.Fatalf("result Field(0) '%s' is unmatch with flds[0] '%s'", flds[0].Name, f.Name)
		}
	}

	if ds.NumField() != tt.wantNumField {
		t.Fatalf("result numfield is unexpected. got: %d, want: %d", ds.NumField(), tt.wantNumField)
	}

	if ds.IsPtr() != tt.wantIsPtr {
		t.Fatalf("unexpected pointer or not result. got: %v, want: %v", ds.IsPtr(), tt.wantIsPtr)
	}

	if tt.wantDefinition != "" {
		if d := cmp.Diff(ds.Definition(), tt.wantDefinition); d != "" {
			t.Fatalf("unexpected mismatch Definition: (-got +want)\n%s", d)
		}

		// 2nd call
		if d := cmp.Diff(ds.Definition(), tt.wantDefinition); d != "" {
			t.Fatalf("unexpected mismatch Definition(2nd call): (-got +want)\n%s", d)
		}
	}

	return true
}

func testBuilderBuildTag(t *testing.T, ds *DynamicStruct, tt buildTest) bool {
	t.Helper()

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
		st, ok := ds.FieldByName(fName)
		if ok {
			if d := cmp.Diff(st.Tag, reflect.StructTag("")); d != "" {
				t.Fatalf("unexpected mismatch Tag: fName: %s, (-got +want)\n%s", fName, d)
			}
		} else {
			t.Fatalf("Field %s does not exist.", fName)
		}

		// test with tag fields
		fName = prefix + "FieldWithTag"
		sft, ok := ds.FieldByName(fName)
		if ok {
			if d := cmp.Diff(sft.Tag, reflect.StructTag(tagWithTag)); d != "" {
				t.Fatalf("unexpected mismatch WithTag.Tag: fName: %s, (-got +want)\n%s", fName, d)
			}
		} else {
			t.Fatalf("Field %s does not exist.", fName)
		}
	}

	return true
}

func testBuilderBuildAddDynamicStruct(t *testing.T, ds *DynamicStruct, tt buildTest) bool {
	t.Helper()

	builder := newTestBuilder()
	builder.AddDynamicStructWithTag("AdditionalDynamicStruct", ds, false, "json")
	builder.AddDynamicStructPtrWithTag("AdditionalDynamicStructPtr", ds, "json")
	newds, err := builder.Build()
	if err != nil {
		t.Fatalf("unexpected error occurred from Build: %v", err)
	}

	if newds.NumField() != tt.wantNumField+2 {
		t.Fatalf("result numfield is unexpected. got: %d, want: %d", newds.NumField(), tt.wantNumField+2)
	}

	_, ok := newds.FieldByName("AdditionalDynamicStruct")
	if !ok {
		t.Fatalf("additional DynamicStruct field does not exist")
	}
	_, ok = newds.FieldByName("AdditionalDynamicStructPtr")
	if !ok {
		t.Fatalf("additional DynamicStructPtr field does not exist")
	}

	// TODO:
	// wantDefinition := tt.wantDefinition + `
	// 	StructFieldWithTag struct { Byte uint8; Bytes []uint8; Int int; Int64 int64; Uint uint; Uint64 uint64; Float32 float32; Float64 float64; String string; Stringptr *string; Stringslice []string; Bool bool; Map map[string]interface {}; Func func(string) interface {}; DynamicTestStruct2 dynamicstruct_test.DynamicTestStruct2; DynamicTestStruct2Ptr *dynamicstruct_test.DynamicTestStruct2; DynamicTestStruct4Slice []dynamicstruct_test.DynamicTestStruct4; DynamicTestStruct4PtrSlice []*dynamicstruct_test.DynamicTestStruct4 } ` + "`json:\"struct_field_with_tag\"`" + `
	// 	StructPtrFieldWithTag *struct { Byte uint8; Bytes []uint8; Int int; Int64 int64; Uint uint; Uint64 uint64; Float32 float32; Float64 float64; String string; Stringptr *string; Stringslice []string; Bool bool; Map map[string]interface {}; Func func(string) interface {}; DynamicTestStruct2 dynamicstruct_test.DynamicTestStruct2; DynamicTestStruct2Ptr *dynamicstruct_test.DynamicTestStruct2; DynamicTestStruct4Slice []dynamicstruct_test.DynamicTestStruct4; DynamicTestStruct4PtrSlice []*dynamicstruct_test.DynamicTestStruct4 } ` + "`json:\"struct_ptr_field_with_tag\"`" + `
	// }`
	// if d := cmp.Diff(newds.Definition(), wantDefinition); d != "" {
	// 	t.Errorf("unexpected mismatch Definition: (-got +want)\n%s", d)
	// 	t.Logf("@@@@@ Entire newds.Definition = %s\n", newds.Definition())
	// 	return false
	// }

	newds.NewInterface()

	return true
}
