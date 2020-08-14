package genericdecoder

import (
	"fmt"

	"github.com/goldeneggg/structil"
)

func ExampleJSONGenericDecoder_Decode() {
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

	dr, err := NewJSONGenericDecoder().Decode(unknownFormatJSON)
	if err != nil {
		panic(err)
	}

	fmt.Println(dr.DynamicStruct.Definition())

	g, err := structil.NewGetter(dr.DecodedInterface)
	if err != nil {
		panic(err)
	}

	fmt.Printf("g.NumField() = %d\n", g.NumField())
	fmt.Printf(
		"ArrayStringField = %#v\nArrayStructField = %#v\nBoolField = %v\nFloat32Field=%v\nIntField=%v\nNullField=%v\nStringField=%v\nStructPtrField=%v\n",
		g.Get("ArrayStringField"),
		g.Get("ArrayStructField"),
		g.Get("BoolField"),
		g.Get("Float32Field"),
		g.Get("IntField"),
		g.Get("NullField"),
		g.Get("StringField"),
		g.Get("StructPtrField"),
	)

	// Output:
	//type DynamicStruct struct {
	//	ArrayStringField []string `json:"ArrayStringField"`
	//	ArrayStructField []map[string]interface {} `json:"ArrayStructField"`
	//	BoolField bool `json:"BoolField"`
	//	Float32Field float64 `json:"Float32Field"`
	//	IntField float64 `json:"IntField"`
	//	NullField interface {} `json:"NullField"`
	//	StringField string `json:"StringField"`
	//	StructPtrField map[string]string `json:"StructPtrField"`
	//}
	// g.NumField() = 8
	// ArrayStringField = []string{"array_str_1", "array_str_2"}
	// ArrayStructField = []map[string]interface {}{map[string]interface {}{"kkk":"kkk1", "vvvv":"vvv1"}, map[string]interface {}{"kkk":"kkk2", "vvvv":"vvv2"}, map[string]interface {}{"kkk":"kkk3", "vvvv":"vvv3"}}
	// BoolField = false
	// Float32Field=9.876
	// IntField=45678
	// NullField=<nil>
	// StringField=かきくけこ
	// StructPtrField=map[key:hugakey value:hugavalue]
}
