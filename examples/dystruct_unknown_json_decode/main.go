package main

import (
	"fmt"

	"github.com/goldeneggg/structil"
	"github.com/goldeneggg/structil/dynamicstruct"
)

func main() {
	fmt.Println(">>>>>>>>>> singleJSON")
	singleJSON()
	fmt.Println(">>>>>>>>>> arrayJSON")
	arrayJSON()
}

func singleJSON() {
	unknownFormatJSON := []byte(`
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

	intf, err := dynamicstruct.JSONToDynamicStructInterface(unknownFormatJSON)
	if err != nil {
		panic(err)
	}

	convertToGetter(intf)
}

func arrayJSON() {
	unknownFormatJSON := []byte(`
[
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
		]
	},
	{
		"string_field":"さしすせそ",
		"int_field":7890,
		"float32_field":4.99,
		"bool_field":true,
		"struct_ptr_field":{
			"key":"hugakeyXXX",
			"value":"hugavalueXXX"
		},
		"array_string_field":[
			"array_str_111",
			"array_str_222"
		],
		"array_struct_field":[
			{
				"kkk":"kkk99",
				"vvvv":"vvv99"
			},
			{
				"kkk":"kkk999",
				"vvvv":"vvv999"
			},
			{
				"kkk":"kkk9999",
				"vvvv":"vvv9999"
			}
		]
	}
]
`)

	intf, err := dynamicstruct.JSONToDynamicStructInterface(unknownFormatJSON)
	if err != nil {
		panic(err)
	}

	// convert from interface{} to []interface{}
	intfArr, ok := intf.([]interface{})
	if !ok {
		panic(fmt.Errorf("intf can not convert to []interface{}: %#v", intf))
	}

	for _, elem := range intfArr {
		convertToGetter(elem)
	}
}

func convertToGetter(i interface{}) structil.Getter {
	g, err := structil.NewGetter(i)
	if err != nil {
		panic(err)
	}

	fmt.Printf("----- g = %#v\n", g)
	fmt.Printf("----- g.NumField() = %#v\n", g.NumField())
	for _, name := range g.Names() {
		fmt.Printf("----- g.Get(%s) = %#v\n", name, g.Get(name))
	}

	return g
}
