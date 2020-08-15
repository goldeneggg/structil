package genericdecoder

import (
	"fmt"

	"github.com/goldeneggg/structil"
)

func Example() {
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

	// Confirm decoded result using Getter with DecodedInterface
	g, err := structil.NewGetter(dr.DecodedInterface)
	if err != nil {
		panic(err)
	}
	fmt.Printf(
		"num of fields=%d\n'StringField'=%s\n'IntField'=%f\n'Float32Field'=%f\n'BoolField'=%t\n'StructPtrField'=%+v\n'ArrayStringField'=%+v\n'ArrayStructField'=%+v\n'NullField'=%+v",
		g.NumField(),
		g.String("StringField"),   // field names of DynamicStruct are camelized original json field key
		g.Float64("IntField"),     // Note: type of unmarshalled number fields are float64. See: https://golang.org/pkg/encoding/json/#Unmarshal
		g.Float64("Float32Field"), // same as above
		g.Bool("BoolField"),
		g.Get("StructPtrField"),
		g.Get("ArrayStringField"),
		g.Get("ArrayStructField"),
		g.Get("NullField"),
	)
	// Output:
	//type DynamicStruct struct {
	//	ArrayStringField []string `json:"array_string_field"`
	//	ArrayStructField []map[string]interface {} `json:"array_struct_field"`
	//	BoolField bool `json:"bool_field"`
	//	Float32Field float64 `json:"float32_field"`
	//	IntField float64 `json:"int_field"`
	//	NullField interface {} `json:"null_field"`
	//	StringField string `json:"string_field"`
	//	StructPtrField map[string]string `json:"struct_ptr_field"`
	//}
	// num of fields=8
	// 'StringField'=かきくけこ
	// 'IntField'=45678.000000
	// 'Float32Field'=9.876000
	// 'BoolField'=false
	// 'StructPtrField'=map[key:hugakey value:hugavalue]
	// 'ArrayStringField'=[array_str_1 array_str_2]
	// 'ArrayStructField'=[map[kkk:kkk1 vvvv:vvv1] map[kkk:kkk2 vvvv:vvv2] map[kkk:kkk3 vvvv:vvv3]]
	// 'NullField'=<nil>
}
