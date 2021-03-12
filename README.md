structil [![PkgGoDev](https://pkg.go.dev/badge/github.com/goldeneggg/structil)](https://pkg.go.dev/github.com/goldeneggg/structil)
==========

[![Workflow Status](https://github.com/goldeneggg/structil/workflows/CI/badge.svg)](https://github.com/goldeneggg/structil/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/goldeneggg/structil)](https://goreportcard.com/report/github.com/goldeneggg/structil)
[![Codecov](https://codecov.io/github/goldeneggg/structil/coverage.svg?branch=master)](https://codecov.io/github/goldeneggg/structil?branch=master)
[![MIT License](http://img.shields.io/badge/license-MIT-lightgrey.svg)](https://github.com/goldeneggg/structil/blob/master/LICENSE)

struct + util = __structil__, for runtime and dynamic environment in Go.


## Why?

I'd like to ...

- conveniently handle and decode the known or unknown formatted JSON/YAML
- conveniently dive into the specific field in nested struct
- simply verify if a field with the specified name and type exists in object
- etc

with Go reflection package experimentally.

```
*** JSON and YAML format is known or unknown ***


JSON →→→→→→→→→→→→→→→→↓        →→ (known format)   struct  →→→→→→→→→→→↓→→→ (use struct directly)
                     ↓        ↑                                      ↓
                     ↓→→ map →→→ (unknown format) "DynamicStruct" →→→→→→ "Getter", "Finder"
                     ↑
YAML →→→→→→→→→→→→→→→→↑
                     ↑
(and other formats) →↑
```

Please see [my medium post](https://medium.com/@s0k0mata/dynamic-and-runtime-struct-utilities-in-go-go-golang-reflection-25c154335185) as well.

## Simple Usage

Try printing the struct definition from __the unknown formatted__ JSON decoding.

```go
package main

import (
	"fmt"

	"github.com/goldeneggg/structil/dynamicstruct/decoder"
)

func main() {
	unknownJSON := []byte(`
{
	"string_field":"かきくけこ",
	"int_field":45678,
	"bool_field":false,
	"object_field":{
		"id":12,
		"name":"the name",
		"nested_object_field": {
			"address": "Tokyo",
			"is_manager": true
		}
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
		}
	],
	"null_field":null
}
`)

	jsonDec, err := decoder.NewJSON(unknownJSON)
	if err != nil {
		panic(err)
	}

	nest := true
	useTag := true
	ds, err := jsonDec.DynamicStruct(nest, useTag)
	if err != nil {
		panic(err)
	}

	// Print struct definition from DynamicStruct
	fmt.Println(ds.Definition())
}
```

This program will print a Go struct definition string as follows.

```
type DynamicStruct struct {
        ArrayStringField []string `json:"array_string_field"`
        ArrayStructField []struct {
                Kkk string `json:"kkk"`
                Vvvv string `json:"vvvv"`
        }
        BoolField bool `json:"bool_field"`
        IntField float64 `json:"int_field"`
        NullField interface {} `json:"null_field"`
        ObjectField struct {
                Id float64 `json:"id"`
                Name string `json:"name"`
                NestedObjectField struct {
                        Address string `json:"address"`
                        IsManager bool `json:"is_manager"`
                }
        }
        StringField string `json:"string_field"`
}
```

- Type name is "DynamicStruct"
- Field names are automatically camelized from input json attribute names
- Fields are ordered by field name
- If `nest` is true, nested object attributes will be also decoded to struct recursively
- If `useTag` is true, JSON Struct tags are defined


## More Examples


### `DynamicStruct`

We can create the dynamic and runtime struct.

See [example code](/dynamicstruct/examples_test.go#L10)


#### JSON unmershal with `DynamicStruct`

A decoding example from JSON to `DynamicStruct` with `StructTag` using `json.Unmarshal([]byte)` as follows.
This example works correctly not only JSON but also YAML, TOML and more.

See [example code](/dynamicstruct/examples_test.go#L107)

### `Getter`

We can access a struct using field name string, like (typed) map.

See [example code](/examples_test.go#L7)


#### `MapGet` method

`MapGet` method provides the __Map__ collection function for slice of struct

See [example code](/examples_test.go#L56)


### `Finder`

We can access usefully nested struct fields using field name string.

See [example code](/examples_test.go#L115)


#### With config file? use `FinderKeys`

We can create a Finder from the configuration file that have some finding target keys. We support some file formats of configuration file such as `yaml`, `json`, `toml` and more.

See [example code](/examples_test.go#L189)

___Thanks for the awesome configuration management library [spf13/viper](https://github.com/spf13/viper).___


## Benchmark

See [this file](https://github.com/goldeneggg/structil/blob/bench-latest/BENCHMARK_LATEST.txt)

It's the latest benchmark result that is executed on GitHub Actions runner instance.
