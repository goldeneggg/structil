package structil_test

import (
	"io"
	"os"
	"reflect"
	"testing"

	. "github.com/goldeneggg/structil"
	"github.com/google/go-cmp/cmp"
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

	assertPanic = func(t *testing.T, wantPanic bool) {
		r := recover()
		if r != nil {
			if !wantPanic {
				t.Errorf("unexpected panic occured: %+v", r)
			}
		} else {
			if wantPanic {
				t.Errorf("expect to occur panic but does not: %+v", r)
			}
		}
	}
)

func TestNewGetter(t *testing.T) {
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
			got, err := NewGetter(tt.args.st)

			if err == nil {
				if _, ok := got.(Getter); !ok {
					t.Errorf("NewGetter() does not return Getter: %+v", got)
				}
			} else {
				if !tt.wantErr {
					t.Errorf("NewGetter() unexpected error %v occured. wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestGetString(t *testing.T) {
	a, err := NewGetter(tStrPtr)
	if err != nil {
		t.Errorf("NewGetter() occurs unexpected error: %v", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		want      reflect.Value
		wantPanic bool
	}{
		{
			name:      "name exists in accessor and it's value is string",
			args:      args{name: "Name"},
			want:      reflect.ValueOf(tStrPtr.Name),
			wantPanic: false,
		},
		{
			name:      "name exists in accessor but it's value is not string",
			args:      args{name: "ID"},
			want:      reflect.ValueOf(tStrPtr.ID),
			wantPanic: false, // TODO: should be true?
		},
		{
			name:      "name does not exist",
			args:      args{name: "XXX"},
			want:      reflect.ValueOf(nil),
			wantPanic: false, // TODO: should be true?
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer assertPanic(t, tt.wantPanic)

			got := a.GetString(tt.args.name)
			if d := cmp.Diff(got, tt.want.String()); d != "" {
				t.Errorf("unexpected mismatch: (-got +want)\n%s", d)
			}
		})
	}
}
