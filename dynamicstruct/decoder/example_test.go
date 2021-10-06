package decoder

import (
	"fmt"
)

func ExampleDecoder_DynamicStruct_json() {
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

	dec, err := FromJSON(unknownFormatJSON)
	if err != nil {
		panic(err)
	}

	ds, err := dec.DynamicStruct(true, true)
	if err != nil {
		panic(err)
	}

	// Print struct definition from DynamicStruct
	fmt.Println(ds.Definition())

	// Output:
	//type DynamicStruct struct {
	//	ArrayStringField []string `json:"array_string_field"`
	//	ArrayStructField []struct {
	//		Kkk string `json:"kkk"`
	//		Vvvv string `json:"vvvv"`
	//	} `json:"array_struct_field"`
	//	BoolField bool `json:"bool_field"`
	//	Float32Field float64 `json:"float32_field"`
	//	IntField float64 `json:"int_field"`
	//	NullField interface {} `json:"null_field"`
	//	StringField string `json:"string_field"`
	//	StructPtrField struct {
	//		Key string `json:"key"`
	//		Value string `json:"value"`
	//	} `json:"struct_ptr_field"`
	//}
}

func ExampleJSONToGetter() {
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

	g, err := JSONToGetter(unknownFormatJSON, true)
	if err != nil {
		panic(err)
	}
	s, _ := g.String("StringField")   // field names of DynamicStruct are camelized original json field key
	i, _ := g.Float64("IntField")     // Note: type of unmarshalled number fields are float64. See: https://golang.org/pkg/encoding/json/#Unmarshal
	f, _ := g.Float64("Float32Field") // same as above
	b, _ := g.Bool("BoolField")
	strct, _ := g.Get("StructPtrField")
	arrS, _ := g.Get("ArrayStringField")
	null, _ := g.Get("NullField")

	// FIXME: test Getter for array struct element
	// arrStrct, _ := g.Get("ArrayStructField")
	fmt.Printf("g.IsStruct(StructPtrField) = %v\n", g.IsStruct("StructPtrField"))
	fmt.Printf("g.IsSlice(ArrayStructField) = %v\n", g.IsSlice("ArrayStructField"))

	fmt.Printf(
		"num of fields=%d\n'StringField'=%s\n'IntField'=%f\n'Float32Field'=%f\n'BoolField'=%t\n'StructPtrField'=%#v\n'ArrayStringField'=%+v\n'NullField'=%+v",
		g.NumField(),
		s,
		i, // Note: type of unmarshalled number fields are float64. See: https://golang.org/pkg/encoding/json/#Unmarshal
		f, // same as above
		b,
		strct,
		arrS,
		null,
	)
	// Output:
	// g.IsStruct(StructPtrField) = true
	// g.IsSlice(ArrayStructField) = true
	// num of fields=8
	// 'StringField'=かきくけこ
	// 'IntField'=45678.000000
	// 'Float32Field'=9.876000
	// 'BoolField'=false
	// 'StructPtrField'=struct { Key string "json:\"key\""; Value string "json:\"value\"" }{Key:"hugakey", Value:"hugavalue"}
	// 'ArrayStringField'=[array_str_1 array_str_2]
	// 'NullField'=<nil>
}

func ExampleDecoder_DynamicStruct_yaml() {
	unknownFormatYAML := []byte(`
string_field: あああ
obj_field:
  id: 45
  name: Test Jiou
  boss: true
  objobj_field:
    user_id: 678
    status: progress
`)

	dec, err := FromYAML(unknownFormatYAML)
	if err != nil {
		panic(err)
	}

	ds, err := dec.DynamicStruct(false, true)
	if err != nil {
		panic(err)
	}

	// Print struct definition from DynamicStruct
	fmt.Println(ds.Definition())

	// Output:
	//type DynamicStruct struct {
	//	ObjField map[string]interface {} `yaml:"obj_field"`
	//	StringField string `yaml:"string_field"`
	//}
}
