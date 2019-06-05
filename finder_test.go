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
