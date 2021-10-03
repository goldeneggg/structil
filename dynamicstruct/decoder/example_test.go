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

	ds, err := dec.DynamicStruct(false, true)
	if err != nil {
		panic(err)
	}

	// Print struct definition from DynamicStruct
	fmt.Println(ds.Definition())

	// Output:
	//type DynamicStruct struct {
	//	ArrayStringField []string `json:"array_string_field"`
	//	ArrayStructField []map[string]interface {} `json:"array_struct_field"`
	//	BoolField bool `json:"bool_field"`
	//	Float32Field float64 `json:"float32_field"`
	//	IntField float64 `json:"int_field"`
	//	NullField interface {} `json:"null_field"`
	//	StringField string `json:"string_field"`
	//	StructPtrField map[string]interface {} `json:"struct_ptr_field"`
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

	/*
		// Confirm decoded result using Getter with Interface
		// *When input map keys are NOT camelized*
		m, ok := dec.Interface().(map[string]interface{})
		if !ok {
			panic(fmt.Sprintf("dec.Interface() does not return map: %#v", dec.Interface()))
		}
		intf, err := ds.DecodeMapWithKeyCamelize(m)
		if err != nil {
			panic(err)
		}
		g, err := structil.NewGetter(intf)
	*/

	g, err := JSONToGetter(unknownFormatJSON)
	if err != nil {
		panic(err)
	}
	s, _ := g.String("StringField")   // field names of DynamicStruct are camelized original json field key
	i, _ := g.Float64("IntField")     // Note: type of unmarshalled number fields are float64. See: https://golang.org/pkg/encoding/json/#Unmarshal
	f, _ := g.Float64("Float32Field") // same as above
	b, _ := g.Bool("BoolField")
	strct, _ := g.Get("StructPtrField")
	arrS, _ := g.Get("ArrayStringField")
	arrStrct, _ := g.Get("ArrayStructField")
	null, _ := g.Get("NullField")
	fmt.Printf(
		"num of fields=%d\n'StringField'=%s\n'IntField'=%f\n'Float32Field'=%f\n'BoolField'=%t\n'StructPtrField'=%#v\n'ArrayStringField'=%+v\n'ArrayStructField'=%+v\n'NullField'=%+v",
		g.NumField(),
		s,
		i, // Note: type of unmarshalled number fields are float64. See: https://golang.org/pkg/encoding/json/#Unmarshal
		f, // same as above
		b,
		strct,
		arrS,
		arrStrct,
		null,
	)
	// Output:
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
