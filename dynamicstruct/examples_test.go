package dynamicstruct

import (
	"fmt"

	"github.com/goldeneggg/structil"
)

func ExampleImpl_DecodeMap() {
	type Hoge struct {
		Key   string
		Value interface{}
	}

	var hogePtr *Hoge

	// Add struct fields using Builder
	b := NewBuilder().
		AddString("StringField").
		AddInt("IntField").
		AddFloat("FloatField").
		AddBool("BoolField").
		AddMap("MapField", SampleString, SampleFloat).
		AddStructPtr("StructPtrField", hogePtr).
		AddSlice("SliceField", hogePtr)

	// Remove one field
	b = b.Remove("FloatField")

	// Build returns a DynamicStruct
	ds := b.Build()

	// Decode to struct from map
	input := map[string]interface{}{
		"StringField": "Abc Def",
		"IntField":    12345,
		"BoolField":   true,
		"MapField":    map[string]float64{"mkey1": 1.23, "mkey2": 4.56},
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
	fmt.Printf(
		"NumField: %d, String: %s, Int: %d, Bool: %v, Map: %+v\n",
		ds.NumField(),
		g.String("StringField"),
		g.Int("IntField"),
		g.Bool("BoolField"),
		g.Get("MapField"),
	)
	// Output:
	// NumField: 6, String: Abc Def, Int: 12345, Bool: true, Map: map[mkey1:1.23 mkey2:4.56]
}
