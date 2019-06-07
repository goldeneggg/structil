// This file is helper settings and functions for test

package structil_test

import (
	"testing"
)

type (
	TestStruct struct {
		ExpBytes       []byte
		ExpInt64       int64
		ExpUint64      uint64
		ExpFloat32     float32
		ExpFloat64     float64
		ExpString      string
		ExpStringptr   *string
		ExpStringslice []string
		ExpBool        bool
		ExpMap         map[string]interface{}
		ExpFunc        func(string) interface{}
		ExpChInt       chan int
		uexpString     string
		TestStruct2
		TestStruct2Ptr     *TestStruct2
		TestStructSlice    []TestStruct4
		TestStructPtrSlice []*TestStruct4
	}

	TestStruct2 struct {
		ExpString string
		*TestStruct3
	}

	TestStruct3 struct {
		ExpString string
		ExpInt    int
	}

	TestStruct4 struct {
		ExpString  string
		ExpString2 string
	}
)

var (
	testString2 = "test name2"
	testFunc    = func(s string) interface{} { return s + "-func" }
	testChan    = make(chan int)

	deferPanic = func(t *testing.T, wantPanic bool, isXXX bool, args interface{}) {
		r := recover()
		if r != nil {
			if !wantPanic {
				t.Errorf("unexpected panic occured: isXXX: %v, args: %+v, %+v", isXXX, args, r)
			}
		} else {
			if wantPanic {
				t.Errorf("expect to occur panic but does not: isXXX: %v, args: %+v, %+v", isXXX, args, r)
			}
		}
	}
)

func newTestStruct() TestStruct {
	return TestStruct{
		ExpBytes:       []byte{0x00, 0xFF},
		ExpInt64:       int64(-1),
		ExpUint64:      uint64(1),
		ExpFloat32:     float32(-1.23),
		ExpFloat64:     float64(-3.45),
		ExpString:      "test name",
		ExpStringptr:   &testString2,
		ExpStringslice: []string{"strslice1", "strslice2"},
		ExpBool:        true,
		ExpMap:         map[string]interface{}{"k1": "v1", "k2": 2},
		ExpFunc:        testFunc,
		ExpChInt:       testChan,
		uexpString:     "unexported string",
		TestStruct2: TestStruct2{
			ExpString: "struct2 string",
			TestStruct3: &TestStruct3{
				ExpString: "struct3 string",
				ExpInt:    -123,
			},
		},
		TestStruct2Ptr: &TestStruct2{
			ExpString: "struct2 string ptr",
			TestStruct3: &TestStruct3{
				ExpString: "struct3 string ptr",
				ExpInt:    -456,
			},
		},
		TestStructSlice: []TestStruct4{
			{
				ExpString:  "key100",
				ExpString2: "value100",
			},
			{
				ExpString:  "key200",
				ExpString2: "value200",
			},
		},
		TestStructPtrSlice: []*TestStruct4{
			{
				ExpString:  "key991",
				ExpString2: "value991",
			},
			{
				ExpString:  "key992",
				ExpString2: "value992",
			},
		},
	}
}

func newTestStructPtr() *TestStruct {
	ts := newTestStruct()
	return &ts
}
