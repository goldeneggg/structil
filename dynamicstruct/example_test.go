package dynamicstruct

import (
	"encoding/json"
	"fmt"

	"github.com/goldeneggg/structil"
)

func Example() {
	type Hoge struct {
		Key   string
		Value interface{}
	}

	hogePtr := &Hoge{
		Key:   "keystr",
		Value: "valuestr",
	}

	// Add fields using Builder with AddXXX method chain
	b := NewBuilder().
		AddString("StringField").
		AddInt("IntField").
		AddFloat32("Float32Field").
		AddBool("BoolField").
		AddMap("MapField", SampleString, SampleFloat32).
		AddStructPtr("StructPtrField", hogePtr).
		AddSlice("SliceField", SampleInt).
		AddInterfaceWithTag("SomeObjectField", true, `json:"some_object_field"`)

	// Remove removes a field by assigned name
	b = b.Remove("Float32Field")

	// SetStructName sets the name of DynamicStruct
	// Note: Default struct name is "DynamicStruct"
	b.SetStructName("MyStruct")

	// Build returns a DynamicStruct
	ds, err := b.Build()
	if err != nil {
		panic(err)
	}

	// Print struct definition with Definition method
	// Struct fields are automatically orderd by field name
	fmt.Println(ds.Definition())

	// DecodeMap decodes from map to DynamicStruct
	input := map[string]interface{}{
		"StringField":     "Abc Def",
		"IntField":        int(123),
		"BoolField":       true,
		"MapField":        map[string]float32{"mkey1": float32(1.23), "mkey2": float32(4.56)},
		"StructPtrField":  hogePtr,
		"SliceField":      []int{111, 222},
		"SomeObjectField": nil,
	}
	dec, err := ds.DecodeMap(input)
	if err != nil {
		panic(err)
	}

	// Confirm decoded result using Getter
	g, err := structil.NewGetter(dec)
	if err != nil {
		panic(err)
	}
	s, _ := g.String("StringField")
	i, _ := g.Int("IntField")
	bl, _ := g.Bool("BoolField")
	m, _ := g.Get("MapField")
	strct, _ := g.Get("StructPtrField")
	sl, _ := g.Get("SliceField")
	obj, _ := g.Get("SomeObjectField")
	fmt.Printf(
		"num of fields=%d\n'StringField'=%s\n'IntField'=%d\n'BoolField'=%t\n'MapField'=%+v\n'StructPtrField'=%+v\n'SliceField'=%+v\n'SomeObjectField'=%+v",
		g.NumField(),
		s,
		i,
		bl,
		m,
		strct,
		sl,
		obj,
	)
	// Output:
	// type MyStruct struct {
	// 	BoolField bool
	// 	IntField int
	// 	MapField map[string]float32
	// 	SliceField []int
	// 	SomeObjectField *interface {} `json:"some_object_field"`
	// 	StringField string
	// 	StructPtrField struct {
	// 		Key string
	// 		Value interface {}
	// 	}
	// }
	// num of fields=7
	// 'StringField'=Abc Def
	// 'IntField'=123
	// 'BoolField'=true
	// 'MapField'=map[mkey1:1.23 mkey2:4.56]
	// 'StructPtrField'={Key:keystr Value:valuestr}
	// 'SliceField'=[111 222]
	// 'SomeObjectField'=<nil>
}

func Example_unmarshalJSON() {
	type Hoge struct {
		Key   string      `json:"key"`
		Value interface{} `json:"value"`
	}

	var hogePtr *Hoge

	b := NewBuilder().
		AddStringWithTag("StringField", `json:"string_field"`).
		AddIntWithTag("IntField", `json:"int_field"`).
		AddFloat32WithTag("Float32Field", `json:"float32_field"`).
		AddBoolWithTag("BoolField", `json:"bool_field"`).
		AddStructPtrWithTag("StructPtrField", hogePtr, `json:"struct_ptr_field"`).
		AddSliceWithTag("SliceField", "", `json:"slice_string_field"`)
	ds, err := b.Build()
	if err != nil {
		panic(err)
	}

	// prints Go struct definition of this DynamicStruct
	fmt.Println(ds.Definition())

	// try json unmarshal with NewInterface
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
	intf := ds.NewInterface() // returns a new interface of this DynamicStruct
	err = json.Unmarshal(input, &intf)
	if err != nil {
		panic(err)
	}

	g, err := structil.NewGetter(intf)
	if err != nil {
		panic(err)
	}
	s, _ := g.String("StringField")
	f, _ := g.Float32("Float32Field")
	strct, _ := g.Get("StructPtrField")
	sl, _ := g.Get("SliceField")
	fmt.Printf(
		"num of fields=%d\n'StringField'=%s\n'Float32Field'=%f\n'StructPtrField'=%+v\n'SliceField'=%+v",
		g.NumField(),
		s,
		f,
		strct,
		sl,
	)
	// Output:
	// type DynamicStruct struct {
	// 	BoolField bool `json:"bool_field"`
	// 	Float32Field float32 `json:"float32_field"`
	// 	IntField int `json:"int_field"`
	// 	SliceField []string `json:"slice_string_field"`
	// 	StringField string `json:"string_field"`
	// 	StructPtrField struct {
	// 		Key string `json:"key"`
	// 		Value interface {} `json:"value"`
	// 	} `json:"struct_ptr_field"`
	// }
	// num of fields=6
	// 'StringField'=あいうえお
	// 'Float32Field'=5.670000
	// 'StructPtrField'={Key:hogekey Value:hogevalue}
	// 'SliceField'=[a b]
}
