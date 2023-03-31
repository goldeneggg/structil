package main

import (
	"fmt"
	"io"
	"os"

	"github.com/goldeneggg/structil/dynamicstruct/decoder"
)

func main() {
	// fmt.Println("NOT IMPLEMENTED YET")

	f, err := os.Open("demo.json")
	if err != nil {
		panic(err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	dec, err := decoder.FromJSON(b)
	if err != nil {
		panic(err)
	}

	nest := true
	useTag := true
	ds, err := dec.DynamicStruct(nest, useTag)
	if err != nil {
		panic(err)
	}

	// Print struct definition from DynamicStruct
	fmt.Println(ds.Definition())
}
