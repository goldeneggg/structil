package main

import (
	"fmt"
	"reflect"

	"github.com/goldeneggg/structil/dynamicstruct"
)

func main() {
	ds := dynamicstruct.New().
		AddString("StringField").
		AddInt("IntField").
		AddFloat("FloatField").
		AddBool("BoolField").
		AddMap("MapField")
	fmt.Printf("ds: %#v\n", ds)
	dsi := ds.Build()
	fmt.Printf("dsi: %#v\n", dsi)

	var sfi reflect.StructField
	for i := 0; i < ds.NumBuiltField(); i++ {
		sfi = ds.BuiltField(i)
		fmt.Printf("i: %d, sfi: %+v\n", i, sfi)

	}
}
