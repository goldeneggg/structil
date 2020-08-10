package main

import (
	"fmt"

	"github.com/iancoleman/strcase"

	"github.com/goldeneggg/structil"
	"github.com/goldeneggg/structil/dynamicstruct"
)

func main() {
	jsonData := []byte(`
{
	"string_field":"かきくけこ",
	"int_field":45678,
	"float32_field":9.876,
	"bool_field":false,
	"struct_ptr_field":{
		"key":"hugakey",
		"value":"hugavalue"
	},
	"array_string_field":[
		"array_str_1",
		"array_str_2"
	],
	"array_struct_field":[
		{
			"kkk":"kkk1",
			"vvvv":"vvv1"
		},
		{
			"kkk":"kkk2",
			"vvvv":"vvv2"
		},
		{
			"kkk":"kkk3",
			"vvvv":"vvv3"
		}
	],
	"null_field":null
}
`)

	intf, err := dynamicstruct.JSONToDynamicStructInterface(jsonData)
	if err != nil {
		panic(err)
	}

	g, err := structil.NewGetter(intf)
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"string_field: %v, int_field: %v, float32_field: %v, bool_field: %v, struct_ptr_field: %v, array_string_field: %v, array_struct_field: %v",
		g.Get(strcase.ToCamel("string_field")),
		g.Get(strcase.ToCamel("int_field")),
		g.Get(strcase.ToCamel("float32_field")),
		g.Get(strcase.ToCamel("bool_field")),
		g.Get(strcase.ToCamel("struct_ptr_field")),
		g.Get(strcase.ToCamel("array_string_field")),
		g.Get(strcase.ToCamel("array_struct_field")),
	)
	// Output:
	// string_field: かきくけこ, int_field: 45678, float32_field: 9.876, bool_field: false, struct_ptr_field: map[key:hugakey value:hugavalue], array_string_field: [array_str_1 array_str_2], array_struct_field: [map[kkk:kkk1 vvvv:vvv1] map[kkk:kkk2 vvvv:vvv2] map[kkk:kkk3 vvvv:vvv3]]
}
