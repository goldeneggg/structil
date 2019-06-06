package structil_test

import (
	"testing"

	"github.com/goldeneggg/structil"
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
			got, err := structil.NewFinder(tt.args.i)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewFinder() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else if _, ok := got.(structil.Finder); !ok {
				t.Errorf("NewFinder() = %v, not Finder type", got)
			}
		})
	}
}

func TestNewFinderWithGetterAndSep(t *testing.T) {
	g, err := structil.NewGetter(newTestStructPtr())
	if err != nil {
		t.Errorf("NewGetter() error = %v", err)
	}

	type args struct {
		g   structil.Getter
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
			_, err := structil.NewFinderWithGetterAndSep(tt.args.g, tt.args.sep)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFinder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestToMap(t *testing.T) {
	f, err := structil.NewFinder(newTestStructPtr())
	if err != nil {
		t.Errorf("NewFinder() error = %v", err)
	}

	type args struct {
		chain structil.Finder
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantMap map[string]interface{}
	}{
		{
			name: "ToMap with non-nest chain",
			args: args{
				chain: f.Find("ExpInt64", "ExpString"),
			},
			wantErr: false,
			wantMap: map[string]interface{}{
				"ExpInt64":  int64(-1),
				"ExpString": testString,
			},
		},
		{
			name: "ToMap with a nest chain",
			args: args{
				chain: f.Find("ExpInt64", "ExpString").
					Struct("TestStruct2").Find("ExpString"),
			},
			wantErr: false,
			wantMap: map[string]interface{}{
				"ExpInt64":              int64(-1),
				"ExpString":             testString,
				"TestStruct2.ExpString": "struct2 string",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tt.args.chain.ToMap()

			if err == nil {
				if res == nil {
					t.Errorf("ToMap() result is nil %v", res)
					return
				}

				for k, v := range tt.wantMap {
					resV, ok := res[k]
					if ok {
						// Note: reflectDeepEqual does not work
						// if d := cmp.Diff(v, resV); d != "" {
						// 	t.Errorf("ToMap() key: %s, want: [%v], resV: [%v], resmap: %+v, diff: \n%s ", k, v, resV, res, d)
						// 	return
						// }
						if v != resV {
							t.Errorf("ToMap() key: %s, want: [%v], resV: [%v], resmap: %+v", k, v, resV, res)
							return
						}
					} else {
						t.Errorf("ToMap() ok: %v, key: %s, want: [%v], resV: [%v], resmap: %+v, ", ok, k, v, resV, res)
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
