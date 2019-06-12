package structil_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/goldeneggg/structil"
)

type toMapTest struct {
	name      string
	wantError bool
	wantPanic bool
	wantMap   map[string]interface{}
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
			args:      args{i: newTestStruct()},
			wantError: false,
		},
		{
			name:      "NewFinder with valid struct ptr",
			args:      args{i: newTestStructPtr()},
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

				if f, ok := got.(Finder); ok {
					if f.GetNameSeparator() != "." {
						t.Errorf("NewFinder() unexpected separator got: %s, want: '.'", f.GetNameSeparator())
					}
				} else {
					t.Errorf("NewFinder() want Finder but got %+v", got)
				}
			} else if !tt.wantError {
				t.Errorf("NewFinder() unexpected error [%v] occured. wantError: %v", err, tt.wantError)
			}
		})
	}
}

func TestNewFinderWithGetterAndSep(t *testing.T) {
	t.Parallel()

	g, err := NewGetter(newTestStructPtr())
	if err != nil {
		t.Errorf("NewGetter() error = %v", err)
	}

	type args struct {
		g   Getter
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

				if f, ok := got.(Finder); ok {
					if f.GetNameSeparator() != tt.args.sep {
						t.Errorf("NewFinderWithGetterAndSep() unexpected separator got: %s, want: %s", f.GetNameSeparator(), tt.args.sep)
					}
				} else {
					t.Errorf("NewFinderWithGetterAndSep() want Finder but got %+v", got)
				}
			} else if !tt.wantError {
				t.Errorf("NewFinderWithGetterAndSep() unexpected error [%v] occured. wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestToMap(t *testing.T) {
	t.Parallel()

	var f Finder
	var err error
	fs := make([]Finder, 10)

	for i := 0; i < len(fs); i++ {
		f, err = NewFinder(newTestStructPtr())
		if err != nil {
			t.Errorf("NewFinder() error = %v", err)
			return
		}

		fs[i] = f
	}

	fsep, err := NewFinderWithSep(newTestStructPtr(), ":")
	if err != nil {
		t.Errorf("NewFinderWithSep() error = %v", err)
		return
	}

	type args struct {
		chain Finder
	}
	tests := []struct {
		name            string
		args            args
		wantError       bool
		wantErrorString string
		wantPanic       bool
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
						"TestStruct2",
						"TestStruct2Ptr",
						"TestStruct4Slice",
						"TestStruct4PtrSlice",
					),
			},
			wantMap: map[string]interface{}{
				"Int64":       int64(-1),
				"Float64":     float64(-3.45),
				"String":      "test name",
				"Stringptr":   testString2,
				"Stringslice": []string{"strslice1", "strslice2"},
				"Bool":        true,
				"Map":         map[string]interface{}{"k1": "v1", "k2": 2},
				//"Func":        testFunc,  // TODO: func is fail
				"ChInt":         testChan,
				"privateString": nil, // unexported field is nil
				"TestStruct2": TestStruct2{
					String:      "struct2 string",
					TestStruct3: &TestStruct3{String: "struct3 string", Int: -123},
				},
				"TestStruct2Ptr": TestStruct2{ // not ptr
					String:      "struct2 string ptr",
					TestStruct3: &TestStruct3{String: "struct3 string ptr", Int: -456},
				},
				"TestStruct4Slice": []TestStruct4{
					{String: "key100", String2: "value100"},
					{String: "key200", String2: "value200"},
				},
				"TestStruct4PtrSlice": []*TestStruct4{
					{String: "key991", String2: "value991"},
					{String: "key992", String2: "value992"},
				},
			},
		},
		{
			name: "with single-nest chain",
			args: args{
				chain: fs[1].
					Into("TestStruct2").Find("String"),
			},
			wantMap: map[string]interface{}{
				"TestStruct2.String": "struct2 string",
			},
		},
		{
			name: "with two-nest chain",
			args: args{
				chain: fs[2].
					Into("TestStruct2Ptr", "TestStruct3").Find("String", "Int"),
			},
			wantMap: map[string]interface{}{
				"TestStruct2Ptr.TestStruct3.String": "struct3 string ptr",
				"TestStruct2Ptr.TestStruct3.Int":    int(-456),
			},
		},
		{
			name: "with multi nest chains",
			args: args{
				chain: fs[3].
					Into("TestStruct2").Find("String").
					Into("TestStruct2Ptr").Find("String").
					Into("TestStruct2Ptr", "TestStruct3").Find("String", "Int"),
			},
			wantMap: map[string]interface{}{
				"TestStruct2.String":                "struct2 string",
				"TestStruct2Ptr.String":             "struct2 string ptr",
				"TestStruct2Ptr.TestStruct3.String": "struct3 string ptr",
				"TestStruct2Ptr.TestStruct3.Int":    int(-456),
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
				chain: fs[7].Into("TestStruct2").Find("NonExist"),
			},
			wantError:       true,
			wantErrorString: "field name NonExist does not exist",
		},
		{
			name: "with Struct with existed and non-existed name and Find",
			args: args{
				chain: fs[8].
					Into("TestStruct2").Find("String").
					Into("TestStruct2", "NonExist").Find("String"),
			},
			wantError:       true,
			wantErrorString: "Error in name: NonExist, key: TestStruct2.NonExist. [name NonExist does not exist]",
		},
		{
			name: "with multi nest chains separated by assigned sep",
			args: args{
				chain: fsep.
					Into("TestStruct2").Find("String").
					Into("TestStruct2Ptr").Find("String").
					Into("TestStruct2Ptr", "TestStruct3").Find("String", "Int"),
			},
			wantMap: map[string]interface{}{
				"TestStruct2:String":                "struct2 string",
				"TestStruct2Ptr:String":             "struct2 string ptr",
				"TestStruct2Ptr:TestStruct3:String": "struct3 string ptr",
				"TestStruct2Ptr:TestStruct3:Int":    int(-456),
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
						"TestStruct2",
						"TestStruct2Ptr",
						"TestStruct4Slice",
						"TestStruct4PtrSlice",
					).
					Into("TestStruct2Ptr").Find("String").
					Into("TestStruct2Ptr", "TestStruct3").Find("String", "Int"),
			},
			wantMap: map[string]interface{}{
				"Int64":       int64(-1),
				"Float64":     float64(-3.45),
				"String":      "test name",
				"Stringptr":   testString2,
				"Stringslice": []string{"strslice1", "strslice2"},
				"Bool":        true,
				"Map":         map[string]interface{}{"k1": "v1", "k2": 2},
				//"Func":        testFunc,  // TODO: func is fail
				"ChInt":         testChan,
				"privateString": nil, // unexported field is nil
				"TestStruct2": TestStruct2{
					String:      "struct2 string",
					TestStruct3: &TestStruct3{String: "struct3 string", Int: -123},
				},
				"TestStruct2Ptr": TestStruct2{ // not ptr
					String:      "struct2 string ptr",
					TestStruct3: &TestStruct3{String: "struct3 string ptr", Int: -456},
				},
				"TestStruct4Slice": []TestStruct4{
					{String: "key100", String2: "value100"},
					{String: "key200", String2: "value200"},
				},
				"TestStruct4PtrSlice": []*TestStruct4{
					{String: "key991", String2: "value991"},
					{String: "key992", String2: "value992"},
				},
				"TestStruct2Ptr.String":             "struct2 string ptr",
				"TestStruct2Ptr.TestStruct3.String": "struct3 string ptr",
				"TestStruct2Ptr.TestStruct3.Int":    int(-456),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferPanic(t, tt.wantPanic, tt.args)

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
	var f Finder
	var fk *FinderKeys
	var err error
	fs := make([]Finder, 5)
	fks := make([]*FinderKeys, 5)

	for i := 0; i < len(fs); i++ {
		fk, err = NewFinderKeysFromConf("examples/finder_from_conf", fmt.Sprintf("ex_test%s_yml", strconv.Itoa(i+1)))
		if err != nil {
			t.Errorf("NewFinderKeysFromConf() error = %v", err)
			return
		}
		fks[i] = fk

		f, err = NewFinder(newTestStructPtr())
		if err != nil {
			t.Errorf("NewFinder() error = %v", err)
			return
		}
		fs[i] = f
	}

	type args struct {
		chain Finder
	}
	tests := []struct {
		name            string
		args            args
		wantError       bool
		wantErrorString string
		wantPanic       bool
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
				"Stringptr":     testString2,
				"Stringslice":   []string{"strslice1", "strslice2"},
				"Bool":          true,
				"Map":           map[string]interface{}{"k1": "v1", "k2": 2},
				"ChInt":         testChan,
				"privateString": nil, // unexported field is nil
				"TestStruct2": TestStruct2{
					String:      "struct2 string",
					TestStruct3: &TestStruct3{String: "struct3 string", Int: -123},
				},
				"TestStruct4Slice": []TestStruct4{
					{String: "key100", String2: "value100"},
					{String: "key200", String2: "value200"},
				},
				"TestStruct4PtrSlice": []*TestStruct4{
					{String: "key991", String2: "value991"},
					{String: "key992", String2: "value992"},
				},
				"TestStruct2Ptr.String":             "struct2 string ptr",
				"TestStruct2Ptr.TestStruct3.String": "struct3 string ptr",
				"TestStruct2Ptr.TestStruct3.Int":    int(-456),
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
			defer deferPanic(t, tt.wantPanic, tt.args)

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
