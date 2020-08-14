package main

import (
	"encoding/json"
	"fmt"

	"github.com/goldeneggg/structil"
	"github.com/goldeneggg/structil/dynamicstruct"
)

// Hoge is test struct
type Hoge struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

var (
	hoge    Hoge
	hogePtr *Hoge
)

func main() {
	b := dynamicstruct.NewBuilder().
		AddStringWithTag("StringField", `json:"string_field"`).
		AddIntWithTag("IntField", `json:"int_field"`).
		AddFloat32WithTag("Float32Field", `json:"float32_field"`).
		AddBoolWithTag("BoolField", `json:"bool_field"`).
		AddStructPtrWithTag("StructPtrField", hogePtr, `json:"struct_ptr_field"`).
		AddSliceWithTag("SliceField", "", `json:"slice_string_field"`)

	// get interface of DynamicStruct using Interface() method
	ds := b.Build()
	fmt.Println(ds.Definition())
	intf := ds.NewInterface()

	// try json unmarshal
	input := []byte(`
{
	"string_field":"あいうえお",
	"int_field":9876,
	"float32_field":5.67,
	"bool_field":true,
	"struct_ptr_field":{
		"key":"hogekey",
		"value":"hogevalue"
	},
	"slice_string_field":[
		"a",
		"b"
	]
}
`)

	err := json.Unmarshal(input, &intf)
	if err != nil {
		panic(err)
	}

	g, err := structil.NewGetter(intf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("String: %v, Float: %v, StructPtr: %+v, Slice: %+v\n", g.String("StringField"), g.Float32("Float32Field"), g.Get("StructPtrField"), g.Get("SliceField"))
	// Output:
	// String: あいうえお, Float: 5.67, StructPtr: {Key:hogekey Value:hogevalue}, Slice: [a b]
}
