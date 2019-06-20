package dynamicstruct_test

import (
	"testing"

	. "github.com/goldeneggg/structil/dynamicstruct"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "Call New func",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New()
			if f, ok := got.(DynamicStruct); !ok {
				t.Errorf("New() got object is not DynamicStruct type. got: %#v", f)
			}
		})
	}
}

/*
func TestBuild(t *testing.T) {
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
*/
