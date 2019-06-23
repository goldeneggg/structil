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
		AddFloat("FloatField").
		AddBool("BoolField").
		AddMap("MapField", dynamicstruct.SampleString, dynamicstruct.SampleFloat).
		AddChanBoth("ChanBothField", dynamicstruct.SampleInt).
		AddStructPtr("StructPtrField", hogePtr).
		AddSlice("SliceField", hogePtr)

	b = b.Remove("FloatField")

	ds := b.Build()

	// try mapstructure.Decode using dynamic struct
	input := map[string]interface{}{
		"StringField": "Test String Field",
		"IntField":    12345,
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
