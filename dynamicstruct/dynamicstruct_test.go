package dynamicstruct_test

import (
	"fmt"
	"runtime"
	"testing"

	. "github.com/goldeneggg/structil/dynamicstruct"
)

var (
	deferPanic = func(t *testing.T, wantPanic bool, args interface{}) {
		r := recover()
		if r != nil {
			msg := fmt.Sprintf("\n%v\n", r)
			for d := 0; ; d++ {
				pc, file, line, ok := runtime.Caller(d)
				if !ok {
					break
				}

				msg = msg + fmt.Sprintf(" -> %d: %s: %s:%d\n", d, runtime.FuncForPC(pc).Name(), file, line)
			}

			if wantPanic {
				t.Logf("OK panic is expected: args: %+v, %s", args, msg)
			} else {
				t.Errorf("unexpected panic occured: args: %+v, %s", args, msg)
			}
		} else {
			if wantPanic {
				t.Errorf("expect to occur panic but does not: args: %+v, %+v", args, r)
			}
		}
	}
)

func TestBuild(t *testing.T) {
	t.Parallel()

	bs := make([]Builder, 10)

	for i := 0; i < len(bs); i++ {
		bs[i] = NewBuilder()
	}

	type args struct {
		builder Builder
	}
	tests := []struct {
		name         string
		args         args
		wantNumField int
		wantPanic    bool
	}{
		{
			name: "with struct that have only primitive fields",
			args: args{
				builder: bs[0].
					AddString("StringField").
					AddInt("IntField").
					AddFloat("FloatField").
					AddBool("BoolField"),
			},
			wantNumField: 4,
		},
		/*
			{
				name: "with single-nest builder",
				args: args{
					builder: fs[1].
						Into("TestStruct2").Find("String"),
				},
				wantMap: map[string]interface{}{
					"TestStruct2.String": "struct2 string",
				},
			},
			{
				name: "with two-nest builder",
				args: args{
					builder: fs[2].
						Into("TestStruct2Ptr", "TestStruct3").Find("String", "Int"),
				},
				wantMap: map[string]interface{}{
					"TestStruct2Ptr.TestStruct3.String": "struct3 string ptr",
					"TestStruct2Ptr.TestStruct3.Int":    int(-456),
				},
			},
			{
				name: "with multi nest builders",
				args: args{
					builder: fs[3].
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
					builder: fs[4].Find("NonExist"),
				},
				wantError:       true,
				wantErrorString: "field name NonExist does not exist",
			},
			{
				name: "with Find with existed and non-existed names",
				args: args{
					builder: fs[5].Find("String", "NonExist"),
				},
				wantError:       true,
				wantErrorString: "field name NonExist does not exist",
			},
			{
				name: "with Struct with non-existed name",
				args: args{
					builder: fs[6].Into("NonExist").Find("String"),
				},
				wantError:       true,
				wantErrorString: "Error in name: NonExist, key: NonExist. [name NonExist does not exist]",
			},
			{
				name: "with Struct with existed name and Find with non-existed name",
				args: args{
					builder: fs[7].Into("TestStruct2").Find("NonExist"),
				},
				wantError:       true,
				wantErrorString: "field name NonExist does not exist",
			},
			{
				name: "with Struct with existed and non-existed name and Find",
				args: args{
					builder: fs[8].
						Into("TestStruct2").Find("String").
						Into("TestStruct2", "NonExist").Find("String"),
				},
				wantError:       true,
				wantErrorString: "Error in name: NonExist, key: TestStruct2.NonExist. [name NonExist does not exist]",
			},
			{
				name: "with multi nest builders separated by assigned sep",
				args: args{
					builder: fsep.
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
				name: "with toplevel and multi-nest find builder using FindTop",
				args: args{
					builder: fs[9].
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
		*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer deferPanic(t, tt.wantPanic, tt.args)

			got := tt.args.builder.Build()
			if got == nil {
				t.Errorf("result is nil %v", got)
				return
			}

			if got.NumField() != tt.wantNumField {
				t.Errorf("result numfield is unexpected. got: %d, want: %d", got.NumField(), tt.wantNumField)
				return
			}
		})
	}
}
