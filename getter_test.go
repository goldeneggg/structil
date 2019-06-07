package structil_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/goldeneggg/structil"
	. "github.com/goldeneggg/structil"
)

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
			name: "Bytes",
			args: &getterTestArgs{name: "Bytes"},
		},
		{
			name: "String",
			args: &getterTestArgs{name: "String"},
		},
		{
			name: "Int64",
			args: &getterTestArgs{name: "Int64"},
		},
		{
			name: "Uint64",
			args: &getterTestArgs{name: "Uint64"},
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
			name: "TestStruct2",
			args: &getterTestArgs{name: "TestStruct2"},
		},
		{
			name: "TestStruct2Ptr",
			args: &getterTestArgs{name: "TestStruct2Ptr"},
		},
		{
			name: "TestStruct4Slice",
			args: &getterTestArgs{name: "TestStruct4Slice"},
		},
		{
			name: "TestStruct4PtrSlice",
			args: &getterTestArgs{name: "TestStruct4PtrSlice"},
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
			name:    "valid struct ptr nil",
			args:    args{i: (*TestStruct)(nil)},
			wantErr: false,
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

	testStructPtr := newTestStructPtr()
	g, err := structil.NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() unexpected error [%v] occured.", err)
		return
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Bytes":
				tt.wantType = reflect.TypeOf(testStructPtr.Bytes)
			case "String":
				tt.wantType = reflect.TypeOf(testStructPtr.String)
			case "Int64":
				tt.wantType = reflect.TypeOf(testStructPtr.Int64)
			case "Uint64":
				tt.wantType = reflect.TypeOf(testStructPtr.Uint64)
			case "Float32":
				tt.wantType = reflect.TypeOf(testStructPtr.Float32)
			case "Float64":
				tt.wantType = reflect.TypeOf(testStructPtr.Float64)
			case "Bool":
				tt.wantType = reflect.TypeOf(testStructPtr.Bool)
			case "Map":
				tt.wantType = reflect.TypeOf(testStructPtr.Map)
			case "Func":
				tt.wantType = reflect.TypeOf(testStructPtr.Func)
			case "ChInt":
				tt.wantType = reflect.TypeOf(testStructPtr.ChInt)
			case "TestStruct2":
				tt.wantType = reflect.TypeOf(testStructPtr.TestStruct2)
			case "TestStruct2Ptr":
				tt.wantType = reflect.TypeOf(testStructPtr.TestStruct2Ptr)
			case "TestStruct4Slice":
				tt.wantType = reflect.TypeOf(testStructPtr.TestStruct4Slice)
			case "TestStruct4PtrSlice":
				tt.wantType = reflect.TypeOf(testStructPtr.TestStruct4PtrSlice)
			case "privateString":
				tt.wantType = reflect.TypeOf(testStructPtr.privateString)
			case "NotExist":
				tt.wantPanic = true
			}

			defer deferPanic(t, tt.wantPanic, tt.args)

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

	testStructPtr := newTestStructPtr()
	g, err := structil.NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
		return
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Bytes":
				tt.wantValue = reflect.ValueOf(testStructPtr.Bytes)
			case "String":
				tt.wantValue = reflect.ValueOf(testStructPtr.String)
			case "Int64":
				tt.wantValue = reflect.ValueOf(testStructPtr.Int64)
			case "Uint64":
				tt.wantValue = reflect.ValueOf(testStructPtr.Uint64)
			case "Float32":
				tt.wantValue = reflect.ValueOf(testStructPtr.Float32)
			case "Float64":
				tt.wantValue = reflect.ValueOf(testStructPtr.Float64)
			case "Bool":
				tt.wantValue = reflect.ValueOf(testStructPtr.Bool)
			case "Map":
				tt.wantValue = reflect.ValueOf(testStructPtr.Map)
			case "Func":
				tt.wantValue = reflect.ValueOf(testStructPtr.Func)
			case "ChInt":
				tt.wantValue = reflect.ValueOf(testStructPtr.ChInt)
			case "TestStruct2":
				tt.wantValue = reflect.ValueOf(testStructPtr.TestStruct2)
			case "TestStruct2Ptr":
				tt.wantValue = reflect.ValueOf(testStructPtr.TestStruct2) // Note: *NOT* TestStruct2Ptr
			case "TestStruct4Slice":
				tt.wantValue = reflect.ValueOf(testStructPtr.TestStruct4Slice)
			case "TestStruct4PtrSlice":
				tt.wantValue = reflect.ValueOf(testStructPtr.TestStruct4PtrSlice)
			case "privateString":
				tt.wantValue = reflect.ValueOf(testStructPtr.privateString)
			case "NotExist":
				tt.wantPanic = true
			}

			defer deferPanic(t, tt.wantPanic, tt.args)

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

	testStructPtr := newTestStructPtr()
	g, err := structil.NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Bytes":
				tt.wantIntf = testStructPtr.Bytes
			case "String":
				tt.wantIntf = testStructPtr.String
			case "Int64":
				tt.wantIntf = testStructPtr.Int64
			case "Uint64":
				tt.wantIntf = testStructPtr.Uint64
			case "Float32":
				tt.wantIntf = testStructPtr.Float32
			case "Float64":
				tt.wantIntf = testStructPtr.Float64
			case "Bool":
				tt.wantIntf = testStructPtr.Bool
			case "Map":
				tt.wantIntf = testStructPtr.Map
			case "Func":
				tt.wantIntf = testStructPtr.Func
			case "ChInt":
				tt.wantIntf = testStructPtr.ChInt
			case "TestStruct2":
				tt.wantIntf = testStructPtr.TestStruct2
			case "TestStruct2Ptr":
				tt.wantIntf = *testStructPtr.TestStruct2Ptr // Note: *NOT* testStructPtr.TestStruct2Ptr
			case "TestStruct4Slice":
				tt.wantIntf = testStructPtr.TestStruct4Slice
			case "TestStruct4PtrSlice":
				tt.wantIntf = testStructPtr.TestStruct4PtrSlice
			case "privateString":
				tt.wantIntf = nil // Note: unexported field is nil
			case "NotExist":
				tt.wantPanic = true
			}

			defer deferPanic(t, tt.wantPanic, tt.args)

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

	testStructPtr := newTestStructPtr()
	g, err := structil.NewGetter(testStructPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	tests := newGetterTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Bytes":
				tt.wantIntf = testStructPtr.Bytes
			case "String":
				tt.wantIntf = testStructPtr.String
			case "Int64":
				tt.wantIntf = testStructPtr.Int64
			case "Uint64":
				tt.wantIntf = testStructPtr.Uint64
			case "Float32":
				tt.wantIntf = testStructPtr.Float32
			case "Float64":
				tt.wantIntf = testStructPtr.Float64
			case "Bool":
				tt.wantIntf = testStructPtr.Bool
			case "Map":
				tt.wantIntf = testStructPtr.Map
			case "Func":
				tt.wantIntf = testStructPtr.Func
			case "ChInt":
				tt.wantIntf = testStructPtr.ChInt
			case "TestStruct2":
				tt.wantIntf = testStructPtr.TestStruct2
			case "TestStruct2Ptr":
				tt.wantIntf = *testStructPtr.TestStruct2Ptr // Note: *NOT* testStructPtr.TestStruct2Ptr
			case "TestStruct4Slice":
				tt.wantIntf = testStructPtr.TestStruct4Slice
			case "TestStruct4PtrSlice":
				tt.wantIntf = testStructPtr.TestStruct4PtrSlice
			case "privateString":
				tt.wantIntf = nil // Note: unexported field is nil
			case "NotExist":
				tt.wantError = true
			}

			defer deferPanic(t, tt.wantPanic, tt.args)

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

func TestBytes(t *testing.T) {
	t.Parallel()

	testStructPtr := newTestStructPtr()
	g, err := structil.NewGetter(testStructPtr)
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

			defer deferPanic(t, tt.wantPanic, tt.args)

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

	testStructPtr := newTestStructPtr()
	g, err := structil.NewGetter(testStructPtr)
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

			defer deferPanic(t, tt.wantPanic, tt.args)

			got := g.String(tt.args.name)
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

	testStructPtr := newTestStructPtr()
	g, err := structil.NewGetter(testStructPtr)
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

			defer deferPanic(t, tt.wantPanic, tt.args)

			got := g.Int64(tt.args.name)
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

	testStructPtr := newTestStructPtr()
	g, err := structil.NewGetter(testStructPtr)
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

			defer deferPanic(t, tt.wantPanic, tt.args)

			got := g.Uint64(tt.args.name)
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

	testStructPtr := newTestStructPtr()
	g, err := structil.NewGetter(testStructPtr)
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

			defer deferPanic(t, tt.wantPanic, tt.args)

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

	testStructPtr := newTestStructPtr()
	g, err := structil.NewGetter(testStructPtr)
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

			defer deferPanic(t, tt.wantPanic, tt.args)

			got := g.Bool(tt.args.name)
			if tt.wantPanic {
				t.Errorf("expected panic did not occur. args: %+v", tt.args)
			} else if d := cmp.Diff(got, tt.wantIntf); d != "" {
				t.Errorf("unexpected mismatch: args: %+v, (-got +want)\n%s", tt.args, d)
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
			case "TestStruct2", "TestStruct2Ptr":
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
			case "Bytes", "TestStruct4Slice", "TestStruct4PtrSlice":
				tt.wantBool = true
			}

			got := g.IsSlice(tt.args.name)
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
			case "TestStruct4Slice":
				tt.args.mapfn = func(i int, g Getter) (interface{}, error) {
					return g.String("String") + "=" + g.String("String2"), nil
				}
				tt.wantIntf = []interface{}{string("key100=value100"), string("key200=value200")}
			case "TestStruct4PtrSlice":
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
