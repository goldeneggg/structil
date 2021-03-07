package deprecated_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/goldeneggg/structil/internal/deprecated"
)

type (
	FinderTestStruct struct {
		Byte          byte
		Bytes         []byte
		Int           int
		Int64         int64
		Uint          uint
		Uint64        uint64
		Float32       float32
		Float64       float64
		String        string
		Stringptr     *string
		Stringslice   []string
		Bool          bool
		Map           map[string]interface{}
		Func          func(string) interface{}
		ChInt         chan int
		privateString string
		FinderTestStruct2
		FinderTestStruct2Ptr      *FinderTestStruct2
		FinderTestStruct4Slice    []FinderTestStruct4
		FinderTestStruct4PtrSlice []*FinderTestStruct4
	}

	FinderTestStruct2 struct {
		String string
		*FinderTestStruct3
	}

	FinderTestStruct3 struct {
		String string
		Int    int
	}

	FinderTestStruct4 struct {
		String  string
		String2 string
	}
)

var (
	finderTestString2 = "test name2"
	finderTestFunc    = func(s string) interface{} { return s + "-func" }
	finderTestChan    = make(chan int)
)

func newFinderTestStruct() FinderTestStruct {
	return FinderTestStruct{
		Byte:          0x61,
		Bytes:         []byte{0x00, 0xFF},
		Int:           int(-2),
		Int64:         int64(-1),
		Uint:          uint(2),
		Uint64:        uint64(1),
		Float32:       float32(-1.23),
		Float64:       float64(-3.45),
		String:        "test name",
		Stringptr:     &finderTestString2,
		Stringslice:   []string{"strslice1", "strslice2"},
		Bool:          true,
		Map:           map[string]interface{}{"k1": "v1", "k2": 2},
		Func:          finderTestFunc,
		ChInt:         finderTestChan,
		privateString: "unexported string",
		FinderTestStruct2: FinderTestStruct2{
			String: "struct2 string",
			FinderTestStruct3: &FinderTestStruct3{
				String: "struct3 string",
				Int:    -123,
			},
		},
		FinderTestStruct2Ptr: &FinderTestStruct2{
			String: "struct2 string ptr",
			FinderTestStruct3: &FinderTestStruct3{
				String: "struct3 string ptr",
				Int:    -456,
			},
		},
		FinderTestStruct4Slice: []FinderTestStruct4{
			{
				String:  "key100",
				String2: "value100",
			},
			{
				String:  "key200",
				String2: "value200",
			},
		},
		FinderTestStruct4PtrSlice: []*FinderTestStruct4{
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

func newFinderTestStructPtr() *FinderTestStruct {
	ts := newFinderTestStruct()
	return &ts
}

func TestNewFinder(t *testing.T) {
	t.Parallel()

	type args struct {
		i interface{}
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
	}{
		{
			name:      "NewFinder with valid struct",
			args:      args{i: newFinderTestStruct()},
			wantError: false,
		},
		{
			name:      "NewFinder with valid struct ptr",
			args:      args{i: newFinderTestStructPtr()},
			wantError: false,
		},
		{
			name:      "NewFinder with string",
			args:      args{i: "string"},
			wantError: true,
		},
		{
			name:      "NewFinder with nil",
			args:      args{i: nil},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFinder(tt.args.i)

			if err == nil {
				if tt.wantError {
					t.Errorf("NewFinder() error did not occur. got: %v", got)
					return
				}
			} else if !tt.wantError {
				t.Errorf("NewFinder() unexpected error [%v] occured. wantError: %v", err, tt.wantError)
			}
		})
	}
}

func TestNewFinderWithGetterAndSep(t *testing.T) {
	t.Parallel()

	g, err := NewGetter(newFinderTestStructPtr())
	if err != nil {
		t.Errorf("NewGetter() error = %v", err)
	}

	type args struct {
		g   *Getter
		sep string
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
	}{
		{
			name:      "NewFinderWithGetterAndSep with valid sep",
			args:      args{g: g, sep: ":"},
			wantError: false,
		},
		{
			name:      "NewFinderWithGetterAndSep with empty sep",
			args:      args{g: g, sep: ""},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFinderWithGetterAndSep(tt.args.g, tt.args.sep)

			if err == nil {
				if tt.wantError {
					t.Errorf("NewFinderWithGetterAndSep() error did not occur. got: %v", got)
					return
				}
			} else if !tt.wantError {
				t.Errorf("NewFinderWithGetterAndSep() unexpected error [%v] occured. wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestToMap(t *testing.T) {
	t.Parallel()

	var f *Finder
	var err error
	fs := make([]*Finder, 10)

	for i := 0; i < len(fs); i++ {
		f, err = NewFinder(newFinderTestStructPtr())
		if err != nil {
			t.Errorf("NewFinder() error = %v", err)
			return
		}

		fs[i] = f
	}

	fsep, err := NewFinderWithSep(newFinderTestStructPtr(), ":")
	if err != nil {
		t.Errorf("NewFinderWithSep() error = %v", err)
		return
	}

	type args struct {
		chain *Finder
	}
	tests := []struct {
		name            string
		args            args
		wantError       bool
		wantErrorString string
		wantMap         map[string]interface{}
		cmpopts         []cmp.Option
	}{
		{
			name: "with toplevel find chain",
			args: args{
				chain: fs[0].
					Find(
						"Int64",
						"Float64",
						"String",
						"Stringptr",
						"Stringslice",
						"Bool",
						"Map",
						//"Func",
						"ChInt",
						"privateString",
						"FinderTestStruct2",
						"FinderTestStruct2Ptr",
						"FinderTestStruct4Slice",
						"FinderTestStruct4PtrSlice",
					),
			},
			wantMap: map[string]interface{}{
				"Int64":       int64(-1),
				"Float64":     float64(-3.45),
				"String":      "test name",
				"Stringptr":   finderTestString2,
				"Stringslice": []string{"strslice1", "strslice2"},
				"Bool":        true,
				"Map":         map[string]interface{}{"k1": "v1", "k2": 2},
				//"Func":        finderTestFunc,  // TODO: func is fail
				"ChInt":         finderTestChan,
				"privateString": nil, // unexported field is nil
				"FinderTestStruct2": FinderTestStruct2{
					String:            "struct2 string",
					FinderTestStruct3: &FinderTestStruct3{String: "struct3 string", Int: -123},
				},
				"FinderTestStruct2Ptr": FinderTestStruct2{ // not ptr
					String:            "struct2 string ptr",
					FinderTestStruct3: &FinderTestStruct3{String: "struct3 string ptr", Int: -456},
				},
				"FinderTestStruct4Slice": []FinderTestStruct4{
					{String: "key100", String2: "value100"},
					{String: "key200", String2: "value200"},
				},
				"FinderTestStruct4PtrSlice": []*FinderTestStruct4{
					{String: "key991", String2: "value991"},
					{String: "key992", String2: "value992"},
				},
			},
		},
		{
			name: "with single-nest chain",
			args: args{
				chain: fs[1].
					Into("FinderTestStruct2").Find("String"),
			},
			wantMap: map[string]interface{}{
				"FinderTestStruct2.String": "struct2 string",
			},
		},
		{
			name: "with two-nest chain",
			args: args{
				chain: fs[2].
					Into("FinderTestStruct2Ptr", "FinderTestStruct3").Find("String", "Int"),
			},
			wantMap: map[string]interface{}{
				"FinderTestStruct2Ptr.FinderTestStruct3.String": "struct3 string ptr",
				"FinderTestStruct2Ptr.FinderTestStruct3.Int":    int(-456),
			},
		},
		{
			name: "with multi nest chains",
			args: args{
				chain: fs[3].
					Into("FinderTestStruct2").Find("String").
					Into("FinderTestStruct2Ptr").Find("String").
					Into("FinderTestStruct2Ptr", "FinderTestStruct3").Find("String", "Int"),
			},
			wantMap: map[string]interface{}{
				"FinderTestStruct2.String":                      "struct2 string",
				"FinderTestStruct2Ptr.String":                   "struct2 string ptr",
				"FinderTestStruct2Ptr.FinderTestStruct3.String": "struct3 string ptr",
				"FinderTestStruct2Ptr.FinderTestStruct3.Int":    int(-456),
			},
		},
		{
			name: "with Find with non-existed name",
			args: args{
				chain: fs[4].Find("NonExist"),
			},
			wantError:       true,
			wantErrorString: "field name NonExist does not exist",
		},
		{
			name: "with Find with existed and non-existed names",
			args: args{
				chain: fs[5].Find("String", "NonExist"),
			},
			wantError:       true,
			wantErrorString: "field name NonExist does not exist",
		},
		{
			name: "with Struct with non-existed name",
			args: args{
				chain: fs[6].Into("NonExist").Find("String"),
			},
			wantError:       true,
			wantErrorString: "Error in name: NonExist, key: NonExist. [name NonExist does not exist]",
		},
		{
			name: "with Struct with existed name and Find with non-existed name",
			args: args{
				chain: fs[7].Into("FinderTestStruct2").Find("NonExist"),
			},
			wantError:       true,
			wantErrorString: "field name NonExist does not exist",
		},
		{
			name: "with Struct with existed and non-existed name and Find",
			args: args{
				chain: fs[8].
					Into("FinderTestStruct2").Find("String").
					Into("FinderTestStruct2", "NonExist").Find("String"),
			},
			wantError:       true,
			wantErrorString: "Error in name: NonExist, key: FinderTestStruct2.NonExist. [name NonExist does not exist]",
		},
		{
			name: "with multi nest chains separated by assigned sep",
			args: args{
				chain: fsep.
					Into("FinderTestStruct2").Find("String").
					Into("FinderTestStruct2Ptr").Find("String").
					Into("FinderTestStruct2Ptr", "FinderTestStruct3").Find("String", "Int"),
			},
			wantMap: map[string]interface{}{
				"FinderTestStruct2:String":                      "struct2 string",
				"FinderTestStruct2Ptr:String":                   "struct2 string ptr",
				"FinderTestStruct2Ptr:FinderTestStruct3:String": "struct3 string ptr",
				"FinderTestStruct2Ptr:FinderTestStruct3:Int":    int(-456),
			},
		},
		{
			name: "with toplevel and multi-nest find chain using FindTop",
			args: args{
				chain: fs[9].
					FindTop(
						"Int64",
						"Float64",
						"String",
						"Stringptr",
						"Stringslice",
						"Bool",
						"Map",
						//"Func",
						"ChInt",
						"privateString",
						"FinderTestStruct2",
						"FinderTestStruct2Ptr",
						"FinderTestStruct4Slice",
						"FinderTestStruct4PtrSlice",
					).
					Into("FinderTestStruct2Ptr").Find("String").
					Into("FinderTestStruct2Ptr", "FinderTestStruct3").Find("String", "Int"),
			},
			wantMap: map[string]interface{}{
				"Int64":       int64(-1),
				"Float64":     float64(-3.45),
				"String":      "test name",
				"Stringptr":   finderTestString2,
				"Stringslice": []string{"strslice1", "strslice2"},
				"Bool":        true,
				"Map":         map[string]interface{}{"k1": "v1", "k2": 2},
				//"Func":        finderTestFunc,  // TODO: func is fail
				"ChInt":         finderTestChan,
				"privateString": nil, // unexported field is nil
				"FinderTestStruct2": FinderTestStruct2{
					String:            "struct2 string",
					FinderTestStruct3: &FinderTestStruct3{String: "struct3 string", Int: -123},
				},
				"FinderTestStruct2Ptr": FinderTestStruct2{ // not ptr
					String:            "struct2 string ptr",
					FinderTestStruct3: &FinderTestStruct3{String: "struct3 string ptr", Int: -456},
				},
				"FinderTestStruct4Slice": []FinderTestStruct4{
					{String: "key100", String2: "value100"},
					{String: "key200", String2: "value200"},
				},
				"FinderTestStruct4PtrSlice": []*FinderTestStruct4{
					{String: "key991", String2: "value991"},
					{String: "key992", String2: "value992"},
				},
				"FinderTestStruct2Ptr.String":                   "struct2 string ptr",
				"FinderTestStruct2Ptr.FinderTestStruct3.String": "struct3 string ptr",
				"FinderTestStruct2Ptr.FinderTestStruct3.Int":    int(-456),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.chain.ToMap()

			if err == nil {
				if tt.wantError {
					t.Errorf("error does not occur. got: %v", got)
					return
				}

				if got == nil {
					t.Errorf("result is nil %v", got)
					return
				}

				for k, wv := range tt.wantMap {
					gv, ok := got[k]
					if ok {
						if d := cmp.Diff(gv, wv, tt.cmpopts...); d != "" {
							t.Errorf("key: %s, gotMap: %+v, (-got +want)\n%s", k, got, d)
							return
						}
					} else {
						t.Errorf("ok: %v, key: %s, gotValue: [%v], wantValue: [%v], gotMap: %+v, ", ok, k, gv, wv, got)
						return
					}
				}
			} else {
				if tt.args.chain.HasError() && tt.wantError {
					if d := cmp.Diff(err.Error(), tt.wantErrorString); d != "" {
						t.Errorf("error string is unmatch. (-got +want)\n%s", d)
						return
					}

					tt.args.chain.Reset()
					if tt.args.chain.HasError() {
						t.Errorf("Reset() does not work expectedly. Errors still remain.")
					}
				} else {
					t.Errorf("unexpected error = %v, HasError: %v, wantError: %v", err, tt.args.chain.HasError(), tt.wantError)
				}
			}
		})
	}
}

// This test should *NOT* be parallel
func TestFromKeys(t *testing.T) {
	var f *Finder
	var fk *FinderKeys
	var err error
	fs := make([]*Finder, 5)
	fks := make([]*FinderKeys, 5)

	for i := 0; i < len(fs); i++ {
		fk, err = NewFinderKeys("../../testdata/finder_from_conf", fmt.Sprintf("ex_test%s_yml", strconv.Itoa(i+1)))
		if err != nil {
			t.Errorf("NewFinderKeys() error = %v", err)
			return
		}
		fks[i] = fk

		f, err = NewFinder(newFinderTestStructPtr())
		if err != nil {
			t.Errorf("NewFinder() error = %v", err)
			return
		}
		fs[i] = f
	}

	type args struct {
		chain *Finder
	}
	tests := []struct {
		name            string
		args            args
		wantError       bool
		wantErrorString string
		wantMap         map[string]interface{}
		cmpopts         []cmp.Option
	}{
		{
			name: "with toplevel find chain",
			args: args{
				chain: fs[0].FromKeys(fks[0]),
			},
			wantMap: map[string]interface{}{
				"Int64":         int64(-1),
				"Float64":       float64(-3.45),
				"String":        "test name",
				"Stringptr":     finderTestString2,
				"Stringslice":   []string{"strslice1", "strslice2"},
				"Bool":          true,
				"Map":           map[string]interface{}{"k1": "v1", "k2": 2},
				"ChInt":         finderTestChan,
				"privateString": nil, // unexported field is nil
				"FinderTestStruct2": FinderTestStruct2{
					String:            "struct2 string",
					FinderTestStruct3: &FinderTestStruct3{String: "struct3 string", Int: -123},
				},
				"FinderTestStruct4Slice": []FinderTestStruct4{
					{String: "key100", String2: "value100"},
					{String: "key200", String2: "value200"},
				},
				"FinderTestStruct4PtrSlice": []*FinderTestStruct4{
					{String: "key991", String2: "value991"},
					{String: "key992", String2: "value992"},
				},
				"FinderTestStruct2Ptr.String":                   "struct2 string ptr",
				"FinderTestStruct2Ptr.FinderTestStruct3.String": "struct3 string ptr",
				"FinderTestStruct2Ptr.FinderTestStruct3.Int":    int(-456),
			},
		},
		{
			name: "with Find with non-existed name",
			args: args{
				chain: fs[1].FromKeys(fks[1]),
			},
			wantError:       true,
			wantErrorString: "field name NonExist does not exist",
		},
		{
			name: "with Find with existed and non-existed names",
			args: args{
				chain: fs[2].FromKeys(fks[2]),
			},
			wantError:       true,
			wantErrorString: "field name NonExist does not exist",
		},
		{
			name: "with Struct with non-existed name",
			args: args{
				chain: fs[3].FromKeys(fks[3]),
			},
			wantError:       true,
			wantErrorString: "Error in name: NonExist, key: NonExist. [name NonExist does not exist]",
		},
		{
			name: "with Struct with existed name and Find with non-existed name",
			args: args{
				chain: fs[4].FromKeys(fks[4]),
			},
			wantError:       true,
			wantErrorString: "field name NonExist does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.chain.ToMap()

			if err == nil {
				if tt.wantError {
					t.Errorf("error does not occur. got: %v", got)
					return
				}

				if got == nil {
					t.Errorf("result is nil %v", got)
					return
				}

				for k, wv := range tt.wantMap {
					gv, ok := got[k]
					if ok {
						if d := cmp.Diff(gv, wv, tt.cmpopts...); d != "" {
							t.Errorf("key: %s, gotMap: %+v, (-got +want)\n%s", k, got, d)
							return
						}
					} else {
						t.Errorf("ok: %v, key: %s, gotValue: [%v], wantValue: [%v], gotMap: %+v, ", ok, k, gv, wv, got)
						return
					}
				}
			} else {
				if tt.args.chain.HasError() && tt.wantError {
					if d := cmp.Diff(err.Error(), tt.wantErrorString); d != "" {
						t.Errorf("error string is unmatch. (-got +want)\n%s", d)
						return
					}

					tt.args.chain.Reset()
					if tt.args.chain.HasError() {
						t.Errorf("Reset() does not work expectedly. Errors still remain.")
					}
				} else {
					t.Errorf("unexpected error = %v, HasError: %v, wantError: %v", err, tt.args.chain.HasError(), tt.wantError)
				}
			}
		})
	}
}

func TestNewFinderKeys(t *testing.T) {
	t.Parallel()

	type args struct {
		d string
		n string
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
		wantLen   int
		wantKeys  []string
	}{
		{
			name:      "with valid yaml file",
			args:      args{d: "../../testdata/finder_from_conf", n: "ex_test1_yml"},
			wantError: false,
			wantLen:   15,
			wantKeys: []string{
				"Int64",
				"Float64",
				"String",
				"Stringptr",
				"Stringslice",
				"Bool",
				"Map",
				"ChInt",
				"privateString",
				"FinderTestStruct2",
				"FinderTestStruct4Slice",
				"FinderTestStruct4PtrSlice",
				"FinderTestStruct2Ptr.String",
				"FinderTestStruct2Ptr.FinderTestStruct3.String",
				"FinderTestStruct2Ptr.FinderTestStruct3.Int",
			},
		},
		{
			name:      "with valid json file",
			args:      args{d: "../../testdata/finder_from_conf", n: "ex_test1_json"},
			wantError: false,
			wantLen:   15,
			wantKeys: []string{
				"Int64",
				"Float64",
				"String",
				"Stringptr",
				"Stringslice",
				"Bool",
				"Map",
				"ChInt",
				"privateString",
				"FinderTestStruct2",
				"FinderTestStruct4Slice",
				"FinderTestStruct4PtrSlice",
				"FinderTestStruct2Ptr.String",
				"FinderTestStruct2Ptr.FinderTestStruct3.String",
				"FinderTestStruct2Ptr.FinderTestStruct3.Int",
			},
		},
		{
			name:      "with invalid conf file that Keys does not exist",
			args:      args{d: "../../testdata/finder_from_conf", n: "ex_test_nonkeys_yml"},
			wantError: true,
		},
		{
			name:      "with invalid conf file that is empty",
			args:      args{d: "../../testdata/finder_from_conf", n: "ex_test_empty_yml"},
			wantError: true,
		},
		{
			name:      "with invalid conf file",
			args:      args{d: "../../testdata/finder_from_conf", n: "ex_test_invalid_yml"},
			wantError: true,
		},
		{
			name:      "with conf file does not exist",
			args:      args{d: "../../testdata/finder_from_conf", n: "ex_test_notexist"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFinderKeys(tt.args.d, tt.args.n)

			if err == nil {
				if tt.wantError {
					t.Errorf("NewFinderKeys() error did not occur. got: %v", got)
					return
				}

				if got.Len() != tt.wantLen {
					t.Errorf("NewFinderKeys() unexpected len. got: %d, want: %d", got.Len(), tt.wantLen)
				}

				if d := cmp.Diff(got.Keys(), tt.wantKeys); d != "" {
					t.Errorf("NewFinderKeys() unexpected keys. (-got +want)\n%s", d)
				}

			} else if !tt.wantError {
				t.Errorf("NewFinderKeys() unexpected error [%v] occured. wantError: %v", err, tt.wantError)
			}
		})
	}
}

// benchmark tests

func BenchmarkNewFinder_Val(b *testing.B) {
	var f *Finder
	var e error

	testStructVal := newFinderTestStruct()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, e = NewFinder(testStructVal)
		if e == nil {
			_ = f
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", e)
		}
	}
}

func BenchmarkNewFinder_Ptr(b *testing.B) {
	var f *Finder
	var e error

	testStructPtr := newFinderTestStructPtr()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, e = NewFinder(testStructPtr)
		if e == nil {
			_ = f
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", e)
		}
	}
}

func BenchmarkToMap_1FindOnly(b *testing.B) {
	var m map[string]interface{}

	f, err := NewFinder(newFinderTestStructPtr())
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, err = f.Find("String").ToMap()
		if err == nil {
			_ = m
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkToMap_2FindOnly(b *testing.B) {
	var m map[string]interface{}

	f, err := NewFinder(newFinderTestStructPtr())
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, err = f.Find("String", "Int64").ToMap()
		if err == nil {
			_ = m
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkToMap_1Struct_1Find(b *testing.B) {
	var m map[string]interface{}

	f, err := NewFinder(newFinderTestStructPtr())
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, err = f.Into("FinderTestStruct2").Find("String").ToMap()
		if err == nil {
			_ = m
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkToMap_1Struct_1Find_2Pair(b *testing.B) {
	var m map[string]interface{}

	f, err := NewFinder(newFinderTestStructPtr())
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, err = f.Into("FinderTestStruct2").Find("String").Into("FinderTestStruct2Ptr").Find("String").ToMap()
		if err == nil {
			_ = m
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkToMap_2Struct_1Find(b *testing.B) {
	var m map[string]interface{}

	f, err := NewFinder(newFinderTestStructPtr())
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, err = f.Into("FinderTestStruct2", "FinderTestStruct3").Find("String").ToMap()
		if err == nil {
			_ = m
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkToMap_2Struct_2Find(b *testing.B) {
	var m map[string]interface{}

	f, err := NewFinder(newFinderTestStructPtr())
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, err = f.Into("FinderTestStruct2", "FinderTestStruct3").Find("String", "Int").ToMap()
		if err == nil {
			_ = m
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkNewFinderKeys_yml(b *testing.B) {
	f, err := NewFinder(newFinderTestStructPtr())
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fks, err := NewFinderKeys("../../testdata/finder_from_conf", "ex_test1_yml")
		if err == nil {
			_ = f.FromKeys(fks)
			f.Reset()
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkNewFinderKeys_json(b *testing.B) {
	f, err := NewFinder(newFinderTestStructPtr())
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fks, err := NewFinderKeys("../../testdata/finder_from_conf", "ex_test1_json")
		if err == nil {
			_ = f.FromKeys(fks)
			f.Reset()
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}
