package structil_test

import (
	"testing"

	. "github.com/goldeneggg/structil"
	"github.com/google/go-cmp/cmp"
)

type toMapTest struct {
	name      string
	wantError bool
	wantPanic bool
	wantMap   map[string]interface{}
}

func TestNewFinder(t *testing.T) {
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

				if _, ok := got.(Finder); !ok {
					t.Errorf("NewFinder() want Finder but got %+v", got)
				}
			} else if !tt.wantError {
				t.Errorf("NewFinder() unexpected error [%v] occured. wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestNewFinderWithGetterAndSep(t *testing.T) {
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

				if _, ok := got.(Finder); !ok {
					t.Errorf("NewFinderWithGetterAndSep() want Finder but got %+v", got)
				}
			} else if !tt.wantError {
				t.Errorf("NewFinderWithGetterAndSep() unexpected error [%v] occured. wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestToMap(t *testing.T) {
	var f Finder
	var fs []Finder
	var err error

	for i := 0; i < 9; i++ {
		f, err = NewFinder(newTestStructPtr())
		if err != nil {
			t.Errorf("NewFinder() error = %v", err)
			return
		}

		fs = append(fs, f)
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
		name      string
		args      args
		wantError bool
		wantPanic bool
		wantMap   map[string]interface{}
		cmpopts   []cmp.Option
	}{
		{
			name: "with non-nest chain",
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
				"Stringptr":   testString2, // TODO: if pointer, test is fail
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
					Struct("TestStruct2").Find("String"),
			},
			wantMap: map[string]interface{}{
				"TestStruct2.String": "struct2 string",
			},
		},
		{
			name: "with two-nest chain",
			args: args{
				chain: fs[2].
					Struct("TestStruct2Ptr", "TestStruct3").Find("String", "Int"),
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
					Struct("TestStruct2").Find("String").
					Struct("TestStruct2Ptr").Find("String").
					Struct("TestStruct2Ptr", "TestStruct3").Find("String", "Int"),
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
			wantError: true,
		},
		{
			name: "with Find with existed and non-existed names",
			args: args{
				chain: fs[5].Find("String", "NonExist"),
			},
			wantError: true,
		},
		{
			name: "with Struct with non-existed name",
			args: args{
				chain: fs[6].Struct("NonExist").Find("String"),
			},
			wantError: true,
		},
		{
			name: "with Struct with existed name and Find with non-existed name",
			args: args{
				chain: fs[7].Struct("TestStruct2").Find("NonExist"),
			},
			wantError: true,
		},
		{
			name: "with Struct with existed and non-existed name and Find",
			args: args{
				chain: fs[8].
					Struct("TestStruct2").Find("String").
					Struct("TestStruct2", "NonExist").Find("String"),
			},
			wantError: true,
		},
		{
			name: "with multi nest chains separated by assigned sep",
			args: args{
				chain: fsep.
					Struct("TestStruct2").Find("String").
					Struct("TestStruct2Ptr").Find("String").
					Struct("TestStruct2Ptr", "TestStruct3").Find("String", "Int"),
			},
			wantMap: map[string]interface{}{
				"TestStruct2:String":                "struct2 string",
				"TestStruct2Ptr:String":             "struct2 string ptr",
				"TestStruct2Ptr:TestStruct3:String": "struct3 string ptr",
				"TestStruct2Ptr:TestStruct3:Int":    int(-456),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferPanic(t, tt.wantPanic, tt.args)

			got, err := tt.args.chain.ToMap()

			if err == nil {
				if tt.wantError {
					t.Errorf("ToMap() error does not occur. got: %v", got)
					return
				}

				if got == nil {
					t.Errorf("ToMap() result is nil %v", got)
					return
				}

				for k, wv := range tt.wantMap {
					gv, ok := got[k]
					if ok {
						if d := cmp.Diff(gv, wv, tt.cmpopts...); d != "" {
							t.Errorf("ToMap() key: %s, gotMap: %+v, (-got +want)\n%s", k, got, d)
							return
						}
					} else {
						t.Errorf("ToMap() ok: %v, key: %s, gotValue: [%v], wantValue: [%v], gotMap: %+v, ", ok, k, gv, wv, got)
						return
					}
				}
			} else if !tt.wantError {
				t.Errorf("ToMap() error = %v, wantError %v", err, tt.wantError)
				return
			}
		})
	}
}
