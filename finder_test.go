package structil_test

import (
	"testing"

	. "github.com/goldeneggg/structil"
)

type toMapTest struct {
	name      string
	wantErr   bool
	wantPanic bool
	wantMap   map[string]interface{}
}

func TestNewFinder(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "NewFinder with valid struct",
			args:    args{i: newTestStruct()},
			wantErr: false,
		},
		{
			name:    "NewFinder with valid struct ptr",
			args:    args{i: newTestStructPtr()},
			wantErr: false,
		},
		{
			name:    "NewFinder with string",
			args:    args{i: "string"},
			wantErr: true,
		},
		{
			name:    "NewFinder with nil",
			args:    args{i: nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFinder(tt.args.i)

			if err == nil {
				if tt.wantErr {
					t.Errorf("NewFinder() error did not occur. got: %v", got)
					return
				}

				if _, ok := got.(Finder); !ok {
					t.Errorf("NewFinder() want Finder but got %+v", got)
				}
			} else if !tt.wantErr {
				t.Errorf("NewFinder() unexpected error [%v] occured. wantErr %v", err, tt.wantErr)
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
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "NewFinderWithGetterAndSep with valid sep",
			args:    args{g: g, sep: ":"},
			wantErr: false,
		},
		{
			name:    "NewFinderWithGetterAndSep with empty sep",
			args:    args{g: g, sep: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFinderWithGetterAndSep(tt.args.g, tt.args.sep)

			if err == nil {
				if tt.wantErr {
					t.Errorf("NewFinderWithGetterAndSep() error did not occur. got: %v", got)
					return
				}

				if _, ok := got.(Finder); !ok {
					t.Errorf("NewFinderWithGetterAndSep() want Finder but got %+v", got)
				}
			} else if !tt.wantErr {
				t.Errorf("NewFinderWithGetterAndSep() unexpected error [%v] occured. wantErr %v", err, tt.wantErr)
			}
		})
	}
}

/*
func TestToMap_ValidFindOnly(t *testing.T) {
	var f Finder
	var err error

	f, err = NewFinder(newTestStructPtr())
	if err != nil {
		t.Errorf("NewFinder() error = %v", err)
		return
	}

	tst := &toMapTest{
		name: "with non-nest chain",
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
	}

	t.Run(tst.name, func(t *testing.T) {
		defer deferPanic(t, tst.wantPanic, nil)

		got, err := f.Find(
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
		).ToMap()

		if err == nil {
			if got == nil {
				t.Errorf("ToMap() result is nil %v", got)
				return
			}

			for k, wv := range tst.wantMap {
				gv, ok := got[k]
				if ok {
					if d := cmp.Diff(gv, wv); d != "" {
						t.Errorf("ToMap() key: %s, gotMap: %+v, (-got +want)\n%s", k, got, d)
						return
					}
				} else {
					t.Errorf("ToMap() ok: %v, key: %s, gotValue: [%v], wantValue: [%v], gotMap: %+v, ", ok, k, gv, wv, got)
					return
				}
			}
		} else {
			t.Errorf("ToMap() error = %v", err)
			return
		}
	})
}

func TestToMap_1Struct1Find(t *testing.T) {
	var f Finder
	var err error

	f, err = NewFinder(newTestStructPtr())
	if err != nil {
		t.Errorf("NewFinder() error = %v", err)
		return
	}

	tst := &toMapTest{
		name: "with one-nest chain",
		wantMap: map[string]interface{}{
			"TestStruct2.String": "struct2 string",
		},
	}

	t.Run(tst.name, func(t *testing.T) {
		defer deferPanic(t, tst.wantPanic, nil)

		got, err := f.Struct("TestStruct2").Find("String").ToMap()

		if err == nil {
			if got == nil {
				t.Errorf("ToMap() result is nil %v", got)
				return
			}

			for k, wv := range tst.wantMap {
				gv, ok := got[k]
				if ok {
					if d := cmp.Diff(gv, wv); d != "" {
						t.Errorf("ToMap() key: %s, gotMap: %+v, (-got +want)\n%s", k, got, d)
						return
					}
				} else {
					t.Errorf("ToMap() ok: %v, key: %s, gotValue: [%v], wantValue: [%v], gotMap: %+v, ", ok, k, gv, wv, got)
					return
				}
			}
		} else {
			t.Errorf("ToMap() error = %v", err)
			return
		}
	})
}

func TestToMap_2Struct2Find(t *testing.T) {
	var f Finder
	var err error

	f, err = NewFinder(newTestStructPtr())
	if err != nil {
		t.Errorf("NewFinder() error = %v", err)
		return
	}

	tst := &toMapTest{
		name: "with two-nest chain",
		wantMap: map[string]interface{}{
			"TestStruct2Ptr.TestStruct3.String": "struct3 string ptr",
			"TestStruct2Ptr.TestStruct3.Int":    int(-456),
		},
	}

	t.Run(tst.name, func(t *testing.T) {
		defer deferPanic(t, tst.wantPanic, nil)

		got, err := f.Struct("TestStruct2Ptr", "TestStruct3").Find("String", "Int").ToMap()

		if err == nil {
			if got == nil {
				t.Errorf("ToMap() result is nil %v", got)
				return
			}

			for k, wv := range tst.wantMap {
				gv, ok := got[k]
				if ok {
					if d := cmp.Diff(gv, wv); d != "" {
						t.Errorf("ToMap() key: %s, gotMap: %+v, (-got +want)\n%s", k, got, d)
						return
					}
				} else {
					t.Errorf("ToMap() ok: %v, key: %s, gotValue: [%v], wantValue: [%v], gotMap: %+v, ", ok, k, gv, wv, got)
					return
				}
			}
		} else {
			t.Errorf("ToMap() error = %v", err)
			return
		}
	})
}

func TestToMap_MultiStructMultiFind(t *testing.T) {
	var f Finder
	var err error

	f, err = NewFinder(newTestStructPtr())
	if err != nil {
		t.Errorf("NewFinder() error = %v", err)
		return
	}

	tst := &toMapTest{
		name: "with multi-nest chain",
		wantMap: map[string]interface{}{
			"TestStruct2.String":                "struct2 string",
			"TestStruct2Ptr.String":             "struct2 string ptr",
			"TestStruct2Ptr.TestStruct3.String": "struct3 string ptr",
			"TestStruct2Ptr.TestStruct3.Int":    int(-456),
		},
	}

	t.Run(tst.name, func(t *testing.T) {
		defer deferPanic(t, tst.wantPanic, nil)

		got, err := f.
			Struct("TestStruct2").Find("String").
			Struct("TestStruct2Ptr").Find("String").
			Struct("TestStruct2Ptr", "TestStruct3").Find("String", "Int").
			ToMap()

		if err == nil {
			if got == nil {
				t.Errorf("ToMap() result is nil %v", got)
				return
			}

			for k, wv := range tst.wantMap {
				gv, ok := got[k]
				if ok {
					if d := cmp.Diff(gv, wv); d != "" {
						t.Errorf("ToMap() key: %s, gotMap: %+v, (-got +want)\n%s", k, got, d)
						return
					}
				} else {
					t.Errorf("ToMap() ok: %v, key: %s, gotValue: [%v], wantValue: [%v], gotMap: %+v, ", ok, k, gv, wv, got)
					return
				}
			}
		} else {
			t.Errorf("ToMap() error = %v", err)
			return
		}
	})
}

func TestToMap_ErrorFindByInvalidName(t *testing.T) {
	var f Finder
	var err error

	f, err = NewFinder(newTestStructPtr())
	if err != nil {
		t.Errorf("NewFinder() error = %v", err)
		return
	}

	tst := &toMapTest{
		name: "with Find with non-existed name",
	}

	t.Run(tst.name, func(t *testing.T) {
		defer deferPanic(t, tst.wantPanic, nil)

		got, err := f.Find("NonExist").ToMap()
		if err == nil {
			t.Errorf("ToMap() expected error did not occur. got: %v", got)
		}
	})
}

func TestToMap_ErrorStructByInvalidName(t *testing.T) {
	var f Finder
	var err error

	f, err = NewFinder(newTestStructPtr())
	if err != nil {
		t.Errorf("NewFinder() error = %v", err)
		return
	}

	tst := &toMapTest{
		name: "with Find with non-existed name",
	}

	t.Run(tst.name, func(t *testing.T) {
		defer deferPanic(t, tst.wantPanic, nil)

		got, err := f.Struct("NonExist").Find("String").ToMap()
		if err == nil {
			t.Errorf("ToMap() expected error did not occur. got: %v", got)
		}
	})
}

func TestToMap_ErrorStructFindByInvalidName(t *testing.T) {
	var f Finder
	var err error

	f, err = NewFinder(newTestStructPtr())
	if err != nil {
		t.Errorf("NewFinder() error = %v", err)
		return
	}

	tst := &toMapTest{
		name: "with Find with non-existed name",
	}

	t.Run(tst.name, func(t *testing.T) {
		defer deferPanic(t, tst.wantPanic, nil)

		got, err := f.Struct("TestStruct2").Find("NonExist").ToMap()
		if err == nil {
			t.Errorf("ToMap() expected error did not occur. got: %v", got)
		}
	})
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
		wantErr   bool
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
			wantPanic: true,
		},
		{
			name: "with Find with existed and non-existed names",
			args: args{
				chain: fs[5].Find("String", "NonExist"),
			},
			wantPanic: true,
		},
		{
			name: "with Struct with non-existed name",
			args: args{
				chain: fs[6].Struct("NonExist").Find("String"),
			},
			wantPanic: true,
		},
		{
			name: "with Struct with existed name and Find with non-existed name",
			args: args{
				chain: fs[7].Struct("TestStruct2").Find("NonExist"),
			},
			wantPanic: true,
		},
		{
			name: "with Struct with existed and non-existed name and Find",
			args: args{
				chain: fs[8].
					Struct("TestStruct2").Find("String").
					Struct("TestStruct2", "NonExist").Find("String"),
			},
			wantPanic: true,
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
				if tt.wantErr {
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
			} else if !tt.wantErr {
				t.Errorf("ToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
*/
