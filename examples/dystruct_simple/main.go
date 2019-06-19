package main

import (
	"fmt"
	"reflect"

	"github.com/goldeneggg/structil/dynamicstruct"
)

type Hoge struct {
	Key   string
	Value interface{}
}

var (
	hoge    Hoge
	hogePtr *Hoge
)

func main() {
	ds := dynamicstruct.New().
		AddString("StringField").
		AddInt("IntField").
		AddFloat("FloatField").
		AddBool("BoolField").
		AddMap("MapField").
		AddFunc("FuncField").
		AddChanBoth("ChanBothField", dynamicstruct.SampleInt).
		AddChanRecv("ChanRecvField", dynamicstruct.SampleInt).
		AddChanSend("ChanSendField", dynamicstruct.SampleInt).
		AddStruct("StructField", hoge, false).
		AddStructPtr("StructPtrField", hogePtr).
		AddSlice("SliceField", hogePtr)
	fmt.Printf("ds: %#v\n", ds)
	fmt.Printf("ds.Exists(FloatField) BEFORE: %v\n", ds.Exists("FloatField"))
	ds = ds.Remove("FloatField")
	fmt.Printf("ds.Exists(FloatField) AFTER: %v\n", ds.Exists("FloatField"))
	sPtr := ds.Build()
	fmt.Printf("sPtr: %#v\n", sPtr)
	sVal := ds.BuildNonPtr()
	fmt.Printf("sVal: %#v\n", sVal)

	var sfi reflect.StructField
	for i := 0; i < ds.NumBuiltField(); i++ {
		sfi = ds.BuiltField(i)
		fmt.Printf("i: %d, sfi: %+v\n", i, sfi)
	}
}
