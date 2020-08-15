structil [![PkgGoDev](https://pkg.go.dev/badge/github.com/goldeneggg/structil)](https://pkg.go.dev/github.com/goldeneggg/structil)
==========

[![Workflow Status](https://github.com/goldeneggg/structil/workflows/CI/badge.svg)](https://github.com/goldeneggg/structil/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/goldeneggg/structil)](https://goreportcard.com/report/github.com/goldeneggg/structil)
[![Codecov](https://codecov.io/github/goldeneggg/structil/coverage.svg?branch=master)](https://codecov.io/github/goldeneggg/structil?branch=master)
[![MIT License](http://img.shields.io/badge/license-MIT-lightgrey.svg)](https://github.com/goldeneggg/structil/blob/master/LICENSE)

struct + util = __structil__, for runtime and dynamic environment in Go.


## Why?

I'd like to ...

- conveniently dive into the specific field in nested struct
- simply verify if a field with the specified name and type exists in object
- conveniently handle known ___and unknown___ formatted JSON/YAML
- etc

```
*** JSON and YAML format is known or unknown ***


JSON →↓        →→ (known case) struct  →→→→→→→→→↓→→ (use struct directly)
      ↓        ↑                                ↓
      ↓→→ map →→→ (unknown case) DynamicStruct →→→ Getter, Finder
      ↑
YAML →↑
```

Please see [my medium post](https://medium.com/@s0k0mata/dynamic-and-runtime-struct-utilities-in-go-go-golang-reflection-25c154335185) as well.


## Examples

### `Finder`
We can access usefully nested struct fields using field name string.

See [example code](/examples_test.go)


#### With config file? use `FinderKeys`
We can create a Finder from the configuration file that have some finding target keys. We support some file formats of configuration file such as `yaml`, `json`, `toml` and more.

See [example code](/examples_test.go)

___Thanks for the awesome configuration management library [spf13/viper](https://github.com/spf13/viper).___


### `Getter`
We can access a struct using field name string, like map.

See [example code](/examples_test.go)


#### `MapGet` method
`MapGet` method provides the __Map__ collection function for slice of struct

See [example code](/examples_test.go)


### `DynamicStruct`
We can create the dynamic and runtime struct.

See [example code](/dynamicstruct/examples_test.go)


#### JSON unmershal with `DynamicStruct`
A decoding example from JSON to `DynamicStruct` with `StructTag` using `json.Unmarshal([]byte)` as follows.
This example works correctly not only JSON but also YAML, TOML and more.

See [example code](/dynamicstruct/examples_test.go)

### `GenericDecoder`
A decoding example from __unknown format__ JSON to interface of `DynamicStruct` with `JSONGenericDecoder.Decode` as follows.

See [example code](/dynamicstruct/genericdecoder/examples_test.go)


## Benchmark
See [this file](https://github.com/goldeneggg/structil/blob/bench-latest/BENCHMARK_LATEST.txt)

It's the latest benchmark result that is executed on GitHub Actions runner instance.
