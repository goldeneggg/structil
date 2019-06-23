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
		AddFloatWithTag("FloatField", `json:"float_field"`).
		AddBoolWithTag("BoolField", `json:"bool_field"`).
		AddStructPtrWithTag("StructPtrField", hogePtr, `json:"struct_ptr_field"`)

	// get interface of DynamicStruct using Interface() method
	ds := b.Build()
	intf := ds.Interface()

	// try json unmarshal
	input := []byte(`
{
	"string_field":"あいうえお",
	"int_field":9876,
	"float_field":5.67,
	"bool_field":true,
	"struct_ptr_field":{
		"key":"hogekey",
		"value":"hogevalue"
	}
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
	fmt.Printf("String: %v, Float: %v, StructPtr: %+v\n", g.String("StringField"), g.Float64("FloatField"), g.Get("StructPtrField"))
}
