package structil_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/goldeneggg/structil"
)

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
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewFinder() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else if _, ok := got.(Finder); !ok {
				t.Errorf("NewFinder() = %v, not Finder type", got)
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
			_, err := NewFinderWithGetterAndSep(tt.args.g, tt.args.sep)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFinder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestToMap(t *testing.T) {
	var f Finder
	var fs []Finder
	var err error

	for i := 0; i < 5; i++ {
		f, err = NewFinder(newTestStructPtr())
		if err != nil {
			t.Errorf("NewFinder() error = %v", err)
			return
		}

		fs = append(fs, f)
	}

	type args struct {
		chain Finder
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantMap map[string]interface{}
		cmpopts []cmp.Option
	}{
		{
			name: "ToMap with non-nest chain",
			args: args{
				chain: fs[0].
					Find(
						"ExpInt64",
						"ExpFloat64",
						"ExpString",
						"ExpStringptr",
						"ExpStringslice",
						"ExpBool",
						"ExpMap",
						//"ExpFunc",
						"ExpChInt",
						"uexpString",
						"TestStruct2",
						"TestStruct2Ptr",
						"TestStructSlice",
						"TestStructPtrSlice",
					),
			},
			wantErr: false,
			wantMap: map[string]interface{}{
				"ExpInt64":       int64(-1),
				"ExpFloat64":     float64(-3.45),
				"ExpString":      testString,
				"ExpStringptr":   testString2, // TODO: if pointer, test is fail
				"ExpStringslice": []string{"strslice1", "strslice2"},
				"ExpBool":        true,
				"ExpMap":         map[string]interface{}{"k1": "v1", "k2": 2},
				//"ExpFunc":        testFunc,  // TODO: func is fail
				"ExpChInt":   testChan,
				"uexpString": nil, // unexported field is nil
				"TestStruct2": TestStruct2{
					ExpString:   "struct2 string",
					TestStruct3: &TestStruct3{ExpString: "struct3 string", ExpInt: -123},
				},
				"TestStruct2Ptr": TestStruct2{ // not ptr
					ExpString:   "struct2 string ptr",
					TestStruct3: &TestStruct3{ExpString: "struct3 string ptr", ExpInt: -456},
				},
				"TestStructSlice": []TestStruct4{
					{ExpString: "key100", ExpString2: "value100"},
					{ExpString: "key200", ExpString2: "value200"},
				},
				"TestStructPtrSlice": []*TestStruct4{
					{ExpString: "key991", ExpString2: "value991"},
					{ExpString: "key992", ExpString2: "value992"},
				},
			},
		},
		{
			name: "ToMap with single-nest chain",
			args: args{
				chain: fs[1].
					Struct("TestStruct2").Find("ExpString"),
			},
			wantErr: false,
			wantMap: map[string]interface{}{
				"TestStruct2.ExpString": "struct2 string",
			},
		},
		{
			name: "ToMap with two-nest chain",
			args: args{
				chain: fs[2].
					Struct("TestStruct2Ptr", "TestStruct3").Find("ExpString", "ExpInt"),
			},
			wantErr: false,
			wantMap: map[string]interface{}{
				"TestStruct2Ptr.TestStruct3.ExpString": "struct3 string ptr",
				"TestStruct2Ptr.TestStruct3.ExpInt":    int(-456),
			},
		},
		{
			name: "ToMap with multi nest chains",
			args: args{
				chain: fs[3].
					Struct("TestStruct2").Find("ExpString").
					Struct("TestStruct2Ptr").Find("ExpString").
					Struct("TestStruct2Ptr", "TestStruct3").Find("ExpString", "ExpInt"),
			},
			wantErr: false,
			wantMap: map[string]interface{}{
				"TestStruct2.ExpString":                "struct2 string",
				"TestStruct2Ptr.ExpString":             "struct2 string ptr",
				"TestStruct2Ptr.TestStruct3.ExpString": "struct3 string ptr",
				"TestStruct2Ptr.TestStruct3.ExpInt":    int(-456),
			},
		},
		{
			name: "ToMap with non-existed name",
			args: args{
				chain: fs[4].Find("NonExist"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.chain.ToMap()

			if err == nil {
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
