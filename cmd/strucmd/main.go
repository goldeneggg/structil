package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

func main() {
	const exampleConfig = `
foo = "bar"
baz = "boop"
`
	// type Config struct {
	// 	Foo string `hcl:"foo"`
	// 	Baz string `hcl:"baz"`
	// }
	// var config Config

	// var intf interface{} // PANIC

	var m map[string]interface{}
	err := hclsimple.Decode("example.hcl", []byte(exampleConfig), nil, &m)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	fmt.Printf("Configuration is %#v\n", m)
	fmt.Printf("foo = %#v\n", m["foo"])
	fmt.Printf("baz = %#v\n", m["baz"])
}
