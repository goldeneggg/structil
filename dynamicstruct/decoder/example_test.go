package decoder

import (
	"fmt"

	"github.com/goldeneggg/structil"
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
			"vvvv":12
		},
		{
			"kkk":"kkk2",
			"vvvv":23
		},
		{
			"kkk":"kkk3",
			"vvvv":34
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
	arrS, _ := g.Get("ArrayStringField")
	null, _ := g.Get("NullField")

	fmt.Printf("g.IsStruct(StructPtrField) = %v\n", g.IsStruct("StructPtrField"))
	strct, _ := g.Get("StructPtrField") // FIXME: fieldの並びが時々変わるflaky

	fmt.Printf("g.IsSlice(ArrayStructField) = %v\n", g.IsSlice("ArrayStructField"))
	arrStrct, _ := g.Slice("ArrayStructField")
	gArrZero, err := structil.NewGetter(arrStrct[0])
	if err != nil {
		panic(err)
	}
	sGArrZero, _ := gArrZero.String("Kkk")
	fGArrZero, _ := gArrZero.Float64("Vvvv")

	fmt.Printf(
		"num of fields=%d\n'StringField'=%s\n'IntField'=%f\n'Float32Field'=%f\n'BoolField'=%t\n'ArrayStringField'=%+v\n'NullField'=%+v\n'StructPtrField'=%#v\n'ArrayStructField[0].Kkk'=%s\n'ArrayStructField[0].Vvvv'=%f",
		g.NumField(),
		s,
		i, // Note: type of unmarshalled number fields are float64. See: https://golang.org/pkg/encoding/json/#Unmarshal
		f, // same as above
		b,
		arrS,
		null,
		strct,
		sGArrZero,
		fGArrZero,
	)
	// Output:
	// g.IsStruct(StructPtrField) = true
	// g.IsSlice(ArrayStructField) = true
	// num of fields=8
	// 'StringField'=かきくけこ
	// 'IntField'=45678.000000
	// 'Float32Field'=9.876000
	// 'BoolField'=false
	// 'ArrayStringField'=[array_str_1 array_str_2]
	// 'NullField'=<nil>
	// 'StructPtrField'=struct { Key string "json:\"key\""; Value string "json:\"value\"" }{Key:"hugakey", Value:"hugavalue"}
	// 'ArrayStructField[0].Kkk'=kkk1
	// 'ArrayStructField[0].Vvvv'=12.000000
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
arr_obj_field:
  - aid: 45
    aname: Test Mike
  - aid: 678
    aname: Test Davis
`)

	dec, err := FromYAML(unknownFormatYAML)
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
	//	ArrObjField []struct {
	//		Aid int `yaml:"aid"`
	//		Aname string `yaml:"aname"`
	//	} `yaml:"arr_obj_field"`
	//	ObjField struct {
	//		Boss bool `yaml:"boss"`
	//		Id int `yaml:"id"`
	//		Name string `yaml:"name"`
	//		ObjobjField struct {
	// 			Status string `yaml:"status"`
	// 			UserId int `yaml:"user_id"`
	//		} `yaml:"objobj_field"`
	//	} `yaml:"obj_field"`
	//	StringField string `yaml:"string_field"`
	//}
}
