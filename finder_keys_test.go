package structil_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/goldeneggg/structil"
)

func TestNewFinderKeysFromConf(t *testing.T) {
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
			args:      args{d: "examples/finder_from_conf", n: "ex_test1_yml"},
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
				"TestStruct2",
				"TestStruct4Slice",
				"TestStruct4PtrSlice",
				"TestStruct2Ptr.String",
				"TestStruct2Ptr.TestStruct3.String",
				"TestStruct2Ptr.TestStruct3.Int",
			},
		},
		{
			name:      "with valid json file",
			args:      args{d: "examples/finder_from_conf", n: "ex_json"},
			wantError: false,
			wantLen:   6,
			wantKeys: []string{
				"Company.Group.Name",
				"Company.Group.Boss",
				"Company.Address",
				"Company.Period",
				"Name",
				"Age",
			},
		},
		{
			name:      "with invalid conf file that Keys does not exist",
			args:      args{d: "examples/finder_from_conf", n: "ex_test_nonkeys_yml"},
			wantError: true,
		},
		{
			name:      "with invalid conf file that is empty",
			args:      args{d: "examples/finder_from_conf", n: "ex_test_empty_yml"},
			wantError: true,
		},
		{
			name:      "with invalid conf file",
			args:      args{d: "examples/finder_from_conf", n: "ex_test_invalid_yml"},
			wantError: true,
		},
		{
			name:      "with conf file does not exist",
			args:      args{d: "examples/finder_from_conf", n: "ex_test_notexist"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFinderKeysFromConf(tt.args.d, tt.args.n)

			if err == nil {
				if tt.wantError {
					t.Errorf("NewFinderKeysFromConf() error did not occur. got: %v", got)
					return
				}

				if got.Len() != tt.wantLen {
					t.Errorf("NewFinderKeysFromConf() unexpected len. got: %d, want: %d", got.Len(), tt.wantLen)
				}

				if d := cmp.Diff(got.Keys(), tt.wantKeys); d != "" {
					t.Errorf("NewFinderKeysFromConf() unexpected keys. (-got +want)\n%s", d)
				}

			} else if !tt.wantError {
				t.Errorf("NewFinderKeysFromConf() unexpected error [%v] occured. wantError: %v", err, tt.wantError)
			}
		})
	}
}
