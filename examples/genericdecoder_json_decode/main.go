package main

import (
	"fmt"

	"github.com/goldeneggg/structil"
	"github.com/goldeneggg/structil/dynamicstruct/genericdecoder"
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

	dr, err := genericdecoder.NewJSONGenericDecoder().Decode(unknownFormatJSON)
	if err != nil {
		panic(err)
	}

	convertToGetter(dr.DecodedInterface)
	// Output:
	// ----- g.NumField() = 8
	// ----- g.Get(ArrayStringField) = []string{"array_str_1", "array_str_2"}
	// ----- g.Get(ArrayStructField) = []map[string]interface {}{map[string]interface {}{"kkk":"kkk1", "vvvv":"vvv1"}, map[string]interface {}{"kkk":"kkk2", "vvvv":"vvv2"}, map[string]interface {}{"kkk":"kkk3", "vvvv":"vvv3"}}
	// ----- g.Get(NullField) = <nil>
	// ----- g.Get(StringField) = "かきくけこ"
	// ----- g.Get(IntField) = 45678
	// ----- g.Get(Float32Field) = 9.876
	// ----- g.Get(BoolField) = false
	// ----- g.Get(StructPtrField) = map[string]string{"key":"hugakey", "value":"hugavalue"}
}

func arrayJSON() {
	unknownFormatJSON := []byte(`[
  {
    "string_field": "かきくけこ",
    "int_field": 45678,
    "float32_field": 9.876,
    "bool_field": false,
    "struct_ptr_field": {
      "key": "hugakey",
      "value": "hugavalue"
    },
    "array_string_field": [
      "array_str_1",
      "array_str_2"
    ],
    "array_struct_field": [
      {
        "kkk": "kkk1",
        "vvvv": "vvv1"
      },
      {
        "kkk": "kkk2",
        "vvvv": "vvv2"
      },
      {
        "kkk": "kkk3",
        "vvvv": "vvv3"
      }
    ]
  },
  {
    "string_field": "さしすせそ",
    "int_field": 7890,
    "float32_field": 4.99,
    "bool_field": true,
    "struct_ptr_field": {
      "key": "hugakeyXXX",
      "value": "hugavalueXXX"
    },
    "array_string_field": [
      "array_str_111",
      "array_str_222"
    ],
    "array_struct_field": [
      {
        "kkk": "kkk99",
        "vvvv": "vvv99"
      },
      {
        "kkk": "kkk999",
        "vvvv": "vvv999"
      },
      {
        "kkk": "kkk9999",
        "vvvv": "vvv9999"
      }
    ]
  }
]`)

	dr, err := genericdecoder.NewJSONGenericDecoder().Decode(unknownFormatJSON)
	if err != nil {
		panic(err)
	}

	// convert from interface{} to []interface{}
	intfArr, ok := dr.DecodedInterface.([]interface{})
	if !ok {
		panic(fmt.Errorf("intf can not convert to []interface{}: %#v", dr.DecodedInterface))
	}

	for _, elem := range intfArr {
		convertToGetter(elem)
	}
	// Output:
	// ----- g.NumField() = 7
	// ----- g.Get(ArrayStringField) = []string{"array_str_1", "array_str_2"}
	// ----- g.Get(ArrayStructField) = []map[string]interface {}{map[string]interface {}{"kkk":"kkk1", "vvvv":"vvv1"}, map[string]interface {}{"kkk":"kkk2", "vvvv":"vvv2"}, map[string]interface {}{"kkk":"kkk3", "vvvv":"vvv3"}}
	// ----- g.Get(StringField) = "かきくけこ"
	// ----- g.Get(IntField) = 45678
	// ----- g.Get(Float32Field) = 9.876
	// ----- g.Get(BoolField) = false
	// ----- g.Get(StructPtrField) = map[string]string{"key":"hugakey", "value":"hugavalue"}
	// ----- g.NumField() = 7
	// ----- g.Get(ArrayStringField) = []string{"array_str_111", "array_str_222"}
	// ----- g.Get(ArrayStructField) = []map[string]interface {}{map[string]interface {}{"kkk":"kkk99", "vvvv":"vvv99"}, map[string]interface {}{"kkk":"kkk999", "vvvv":"vvv999"}, map[string]interface {}{"kkk":"kkk9999", "vvvv":"vvv9999"}}
	// ----- g.Get(StringField) = "さしすせそ"
	// ----- g.Get(IntField) = 7890
	// ----- g.Get(Float32Field) = 4.99
	// ----- g.Get(BoolField) = true
	// ----- g.Get(StructPtrField) = map[string]string{"key":"hugakeyXXX", "value":"hugavalueXXX"}
}

func convertToGetter(i interface{}) structil.Getter {
	g, err := structil.NewGetter(i)
	if err != nil {
		panic(err)
	}

	fmt.Printf("----- g.NumField() = %#v\n", g.NumField())
	for _, name := range g.Names() {
		fmt.Printf("----- g.Get(%s) = %#v\n", name, g.Get(name))
	}

	return g
}
