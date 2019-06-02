package structil_test

import (
	"io"
	"os"
	"testing"

	. "github.com/goldeneggg/structil"
)

type tStr struct {
	ID       int64
	FloatVal float64
	Name     string
	NamePtr  *string
	StrArr   []string
	IsMan    bool
	*tStr2
	KVArr []*KeyValue
}

type tStr2 struct {
	Name   string
	Writer io.Writer
	*tStr3
}

type tStr3 struct {
	Name string
	Cnt  int
}

type KeyValue struct {
	Key   string
	Value string
}

var (
	defName  = "test name"
	defName2 = "test name2"

	tStrVal = tStr{
		ID:       1,
		FloatVal: 1.23,
		Name:     defName,
		NamePtr:  &defName2,
		StrArr:   []string{"strarr1", "strarr2"},
		IsMan:    true,
		tStr2: &tStr2{
			Name:   "tSt22 name",
			Writer: os.Stdout,
			tStr3: &tStr3{
				Name: "tStr3 name name",
				Cnt:  999,
			},
		},
		KVArr: []*KeyValue{
			{
				Key:   "key100",
				Value: "value100",
			},
			{
				Key:   "key200",
				Value: "value200",
			},
		},
	}

	tStrPtr = &tStr{
		ID:       1,
		FloatVal: 1.23,
		Name:     defName,
		NamePtr:  &defName2,
		StrArr:   []string{"strarr1", "strarr2"},
		IsMan:    true,
		tStr2: &tStr2{
			Name:   "tSt22 name",
			Writer: os.Stdout,
			tStr3: &tStr3{
				Name: "tStr3 name name",
				Cnt:  999,
			},
		},
		KVArr: []*KeyValue{
			{
				Key:   "key100",
				Value: "value100",
			},
			{
				Key:   "key200",
				Value: "value200",
			},
		},
	}
)

func TestNewAccessor(t *testing.T) {
	type args struct {
		st interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "arg st is valid struct",
			args:    args{st: tStrVal},
			wantErr: false,
		},
		{
			name:    "arg st is valid struct ptr",
			args:    args{st: tStrPtr},
			wantErr: false,
		},
		{
			name:    "arg st is invalid (nil)",
			args:    args{st: nil},
			wantErr: true,
		},
		{
			name:    "arg st is invalid (string)",
			args:    args{st: "abc"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccessor(tt.args.st)

			if err == nil {
				if _, ok := got.(Accessor); !ok {
					t.Errorf("NewAccessor() does not return Accessor: %+v", got)
				}
			} else {
				if !tt.wantErr {
					t.Errorf("NewAccessor() unexpected error %v occured. wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
