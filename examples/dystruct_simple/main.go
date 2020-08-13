package main

import (
	"fmt"

	"github.com/goldeneggg/structil"
	"github.com/goldeneggg/structil/dynamicstruct"
)

// Hoge is test struct
type Hoge struct {
	Key   string
	Value interface{}
}

var (
	hoge    Hoge
	hogePtr *Hoge
)

func main() {
	// Add fields using Builder.
	// We can use AddXXX method chain.
	b := dynamicstruct.NewBuilder().
		AddString("StringField").
		AddInt("IntField").
		AddFloat32("Float32Field").
		AddBool("BoolField").
		AddMap("MapField", dynamicstruct.SampleString, dynamicstruct.SampleFloat32).
		AddChanBoth("ChanBothField", dynamicstruct.SampleInt).
		AddStructPtr("StructPtrField", hogePtr).
		AddSlice("SliceField", hogePtr)

	// Remove removes a field by assigned name.
	b = b.Remove("Float32Field")

	// Build generates a DynamicStruct
	ds := b.Build()

	fmt.Println(ds.Definition())

	// DecodeMap decodes from map to DynamicStruct
	// try mapstructure.Decode using dynamic struct
	input := map[string]interface{}{
		"StringField": "Test String Field",
		"IntField":    int(12),
		"BoolField":   true,
	}

	// 2nd arg need to be a pointer of dynamic struct
	dec, err := ds.DecodeMap(input)
	if err != nil {
		panic(err)
	}

	// Confirm decoded result using Getter
	g, err := structil.NewGetter(dec)
	if err != nil {
		panic(err)
	}

	fmt.Printf("String: %v, Int: %v, Bool: %v\n", g.String("StringField"), g.Int("IntField"), g.Get("BoolField"))
	// Output:
	// String: Test String Field, Int: 12, Bool: true
}
