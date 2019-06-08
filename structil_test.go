// This file is helper settings and functions for test

package structil_test

import (
	"fmt"
	"runtime"
	"testing"

	. "github.com/goldeneggg/structil"
)

type (
	TestStruct struct {
		Byte          byte
		Bytes         []byte
		Int           int
		Int64         int64
		Uint          uint
		Uint64        uint64
		Float32       float32
		Float64       float64
		String        string
		Stringptr     *string
		Stringslice   []string
		Bool          bool
		Map           map[string]interface{}
		Func          func(string) interface{}
		ChInt         chan int
		privateString string
		TestStruct2
		TestStruct2Ptr      *TestStruct2
		TestStruct4Slice    []TestStruct4
		TestStruct4PtrSlice []*TestStruct4
	}

	TestStruct2 struct {
		String string
		*TestStruct3
	}

	TestStruct3 struct {
		String string
		Int    int
	}

	TestStruct4 struct {
		String  string
		String2 string
	}
)

var (
	testIntf interface{}
)

var (
	testString2 = "test name2"
	testFunc    = func(s string) interface{} { return s + "-func" }
	testChan    = make(chan int)

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

func newTestStruct() TestStruct {
	return TestStruct{
		Byte:          0x61,
		Bytes:         []byte{0x00, 0xFF},
		Int:           int(-2),
		Int64:         int64(-1),
		Uint:          uint(2),
		Uint64:        uint64(1),
		Float32:       float32(-1.23),
		Float64:       float64(-3.45),
		String:        "test name",
		Stringptr:     &testString2,
		Stringslice:   []string{"strslice1", "strslice2"},
		Bool:          true,
		Map:           map[string]interface{}{"k1": "v1", "k2": 2},
		Func:          testFunc,
		ChInt:         testChan,
		privateString: "unexported string",
		TestStruct2: TestStruct2{
			String: "struct2 string",
			TestStruct3: &TestStruct3{
				String: "struct3 string",
				Int:    -123,
			},
		},
		TestStruct2Ptr: &TestStruct2{
			String: "struct2 string ptr",
			TestStruct3: &TestStruct3{
				String: "struct3 string ptr",
				Int:    -456,
			},
		},
		TestStruct4Slice: []TestStruct4{
			{
				String:  "key100",
				String2: "value100",
			},
			{
				String:  "key200",
				String2: "value200",
			},
		},
		TestStruct4PtrSlice: []*TestStruct4{
			{
				String:  "key991",
				String2: "value991",
			},
			{
				String:  "key992",
				String2: "value992",
			},
		},
	}
}

func newTestStructPtr() *TestStruct {
	ts := newTestStruct()
	return &ts
}

func newTestGetter() (Getter, error) {
	return NewGetter(newTestStructPtr())
}
