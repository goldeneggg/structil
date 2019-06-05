package structil_test

import (
	"reflect"
	"testing"

	. "github.com/goldeneggg/structil"
)

func TestNewFinder(t *testing.T) {
	tests := []struct {
		name string
		want Finder
	}{
		{
			name: "Default NewFinder",
			want: NewFinder(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFinder()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFinder() = %v, want %v", got, tt.want)
			}
			if got.GetNameSeparator() != "." {
				t.Errorf("separator %s, want '.'", got.GetNameSeparator())
			}
		})
	}
}

func TestNewFinderWithSep(t *testing.T) {
	type args struct {
		sep string
	}
	tests := []struct {
		name string
		args args
		want Finder
	}{
		{
			name: "Use non-default separator",
			args: args{sep: ":"},
			want: NewFinderWithSep(":"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFinderWithSep(tt.args.sep)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFinderWithSep() = %v, want %v", got, tt.want)
			}
			if got.GetNameSeparator() != tt.args.sep {
				t.Errorf("separator %s, want %s", got.GetNameSeparator(), tt.args.sep)
			}
		})
	}
}
