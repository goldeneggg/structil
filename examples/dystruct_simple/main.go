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
	b := dynamicstruct.NewBuilder().
		AddString("StringField").
		AddInt("IntField").
		AddFloat32("Float32Field").
		AddBool("BoolField").
		AddMap("MapField", dynamicstruct.SampleString, dynamicstruct.SampleFloat32).
		AddChanBoth("ChanBothField", dynamicstruct.SampleInt).
		AddStructPtr("StructPtrField", hogePtr).
		AddSlice("SliceField", hogePtr)

	b = b.Remove("Float32Field")

	ds := b.Build()

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

	g, err := structil.NewGetter(dec)
	if err != nil {
		panic(err)
	}
	fmt.Printf("String: %v, Int: %v, Bool: %v\n", g.String("StringField"), g.Int("IntField"), g.Get("BoolField"))
}
