package main

import (
	"fmt"

	"github.com/goldeneggg/structil/dynamicstruct"
)

// Hoge is test struct
type Hoge struct {
	Key   string
	Value interface{}
}

var (
	hoge    Hoge
	hogePtr *Hoge
)

func main() {
	b := dynamicstruct.NewBuilder().
		AddString("StringField").
		AddInt("IntField").
		AddFloat("FloatField").
		AddBool("BoolField").
		AddMap("MapField", dynamicstruct.SampleString, dynamicstruct.SampleFloat).
		AddFunc("FuncField", []interface{}{dynamicstruct.SampleInt, dynamicstruct.SampleInt}, []interface{}{dynamicstruct.SampleBool}).
		AddChanBoth("ChanBothField", dynamicstruct.SampleInt).
		AddChanRecv("ChanRecvField", dynamicstruct.SampleInt).
		AddChanSend("ChanSendField", dynamicstruct.SampleInt).
		AddStruct("StructField", hoge, false).
		AddStructPtr("StructPtrField", hogePtr).
		AddSlice("SliceField", hogePtr)
	b = b.Remove("FloatField")
	ds := b.Build()
	fmt.Printf("ds: %#v\n", ds)

	// try mapstructure.Decode using dynamic struct
	input := map[string]interface{}{
		"StringField": "@@@!!!@@@",
		"IntField":    12345,
		"BoolField":   true,
		"extra": map[string]float64{
			"twitter": 3.14,
		},
	}

	// 2nd arg need to be a pointer of dynamic struct
	dec, err := ds.DecodeMap(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decoded intf: %#v\n", dec)
}
