package main

import (
	"fmt"

	"github.com/goldeneggg/structil/dynamicstruct"
	"github.com/mitchellh/mapstructure"
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
	ds := dynamicstruct.New().
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
	ds = ds.Remove("FloatField")
	sPtr := ds.Build()
	fmt.Printf("sPtr: %#v\n", sPtr)

	// try mapstructure.Decode using dynamic struct
	input := map[string]interface{}{
		"StringField": "Mitchell",
		"IntField":    91,
		"extra": map[string]float64{
			"twitter": 3.14,
		},
	}

	// 2nd arg need to be a pointer of dynamic struct
	err := mapstructure.Decode(input, &sPtr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decoded sPtr: %#v\n", sPtr)
}
